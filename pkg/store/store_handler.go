package store

import (
	"encoding/json"
	"log"
	"time"

	"github.com/balchua/balsa/pkg/fsm"
	"github.com/hashicorp/raft"
)

// StoreHandler struct handler
type StoreHandler struct {
	raft *raft.Raft
}

func New(raft *raft.Raft) *StoreHandler {
	return &StoreHandler{
		raft: raft,
	}
}

func (h StoreHandler) Store(key string, value string) error {

	//Check if current is leader
	if h.raft.State() != raft.Leader {
		log.Print("Not a leader, do nothing")
		return nil
	}
	command := fsm.CommandOperation{
		Operation: "SET",
		Key:       key,
		Value:     value,
	}

	data, err := json.Marshal(command)

	if err != nil {
		log.Fatal("Invalid Command")
	}

	_ = h.raft.Apply(data, 500*time.Millisecond)
	return nil

}
