package rafthandler

import (
	"fmt"

	"github.com/hashicorp/raft"
	"github.com/prometheus/common/log"
)

// Handler struct handler
type Handler struct {
	raft *raft.Raft
}

func New(raft *raft.Raft) *Handler {
	return &Handler{
		raft: raft,
	}
}

func (h Handler) GetRaft() *raft.Raft {
	return h.raft
}

// Join handling join raft
func (h Handler) Join(port int, nodeId string) error {

	log.Info("Joining node % on host 127.0.0.1:%d", nodeId)
	if h.raft.State() != raft.Leader {
		log.Info("not a leader")
	}

	configFuture := h.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		log.Error(err)
	}

	bindAddress := fmt.Sprintf("127.0.0.1:%d", port)
	// This must be run on the leader or it will fail.
	f := h.raft.AddVoter(raft.ServerID(nodeId), raft.ServerAddress(bindAddress), 0, 0)
	if f.Error() != nil {
		log.Error("Invalid voter")
	}

	return nil
}
