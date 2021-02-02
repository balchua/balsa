package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/balchua/balsa/pkg/config"
	"github.com/balchua/balsa/pkg/fsm"
	rafthandler "github.com/balchua/balsa/pkg/raft"
	"github.com/balchua/balsa/pkg/store"
	"github.com/gorilla/mux"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

const (
	maxPool = 3

	tcpTimeout = 10 * time.Second

	// The `retain` parameter controls how many
	// snapshots are retained. Must be at least 1.
	raftSnapShotRetain = 2

	// raftLogCacheSize is the maximum number of logs to cache in-memory.
	// This is used to reduce disk I/O for the recently committed entries.
	raftLogCacheSize = 3
)

// Server struct handling server
type WebServer struct {
	listenAddress string
	r             *mux.Router
	raftHandler   *rafthandler.Handler
	store         *store.StoreHandler
}

//JoinHandler handles requests to join the cluster
func (s *WebServer) JoinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Join handler called")
	nodeID := r.FormValue("nodeId")
	log.Println(fmt.Sprintf("node is %s", nodeID))
	raftPort := r.FormValue("raftPort")
	port, err := strconv.Atoi(raftPort)
	if err != nil {
		log.Fatal("Invalid port %e", err)
	}
	s.raftHandler.Join(port, nodeID)
}

//SetHandler handles setting values to the store (fsm)
func (s *WebServer) SetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Set handler called")
	key := r.FormValue("key")
	value := r.FormValue("value")
	s.store.Store(key, value)
}

//GetHandler handles getting values from the store (fsm)
func (s *WebServer) GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get handler called")
}

//Start will start the http server
func (s *WebServer) Start() {
	srv := &http.Server{
		Handler: s.r,
		Addr:    s.listenAddress,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func initializeRaft(config *config.Config) *rafthandler.Handler {
	var raftBinAddr = fmt.Sprintf("127.0.0.1:%d", config.RaftPort)

	raftConf := raft.DefaultConfig()
	raftConf.LocalID = raft.ServerID(config.RaftNodeId)
	raftConf.SnapshotThreshold = 1024

	store, err := raftboltdb.NewBoltStore(filepath.Join(config.RaftVolume, "raft.dataRepo"))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Wrap the store in a LogCache to improve performance.
	cacheStore, err := raft.NewLogCache(raftLogCacheSize, store)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	snapshotStore, err := raft.NewFileSnapshotStore(config.RaftVolume, raftSnapShotRetain, os.Stdout)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", raftBinAddr)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	transport, err := raft.NewTCPTransport(raftBinAddr, tcpAddr, maxPool, tcpTimeout, os.Stdout)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	fsm := fsm.NewKvFSM()
	raftServer, err := raft.NewRaft(raftConf, fsm, cacheStore, store, snapshotStore, transport)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// always start single server as a leader
	configuration := raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      raft.ServerID(config.RaftNodeId),
				Address: transport.LocalAddr(),
			},
		},
	}

	raftServer.BootstrapCluster(configuration)
	return rafthandler.New(raftServer)
}

//New instantiates the handlers
func New(config *config.Config) *WebServer {

	handler := initializeRaft(config)
	listenAddress := fmt.Sprintf("%s:%d", "127.0.0.1", config.HttpPort)
	storeHandler := store.New(handler.GetRaft())

	s := &WebServer{
		r:             mux.NewRouter(),
		listenAddress: listenAddress,
		raftHandler:   handler,
		store:         storeHandler,
	}

	s.r.HandleFunc("/join", s.JoinHandler)
	s.r.HandleFunc("/set", s.SetHandler)
	s.r.HandleFunc("/get", s.GetHandler)

	return s
}
