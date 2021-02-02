package fsm

import (
	"github.com/hashicorp/raft"
)

type FsmSnapshot struct {
}

// Persist persist to disk. Return nil on success, otherwise return error.
func (s FsmSnapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (s FsmSnapshot) Release() {}

func newDoNothingSnapshot() (raft.FSMSnapshot, error) {
	return &FsmSnapshot{}, nil
}
