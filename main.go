package main

import (
	"log"

	"github.com/balchua/balsa/pkg/config"
	"github.com/balchua/balsa/pkg/server"
	"github.com/spf13/viper"
)

var configKeys = []string{
	"HTTP_SERVER_PORT",
	"RAFT_SERVER_PORT",
	"RAFT_NODE_ID",
	"RAFT_VOLUME_DIR",
}

func main() {
	var v = viper.New()
	v.AutomaticEnv()
	if err := v.BindEnv(configKeys...); err != nil {
		log.Fatal(err)
		return
	}
	conf := &config.Config{
		HttpPort:   v.GetInt("HTTP_SERVER_PORT"),
		RaftPort:   v.GetInt("RAFT_SERVER_PORT"),
		RaftNodeId: v.GetString("RAFT_NODE_ID"),
		RaftVolume: v.GetString("RAFT_VOLUME_DIR"),
	}

	s := server.New(conf)
	s.Start()
}
