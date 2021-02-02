package fsm

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	sysout "log"
	"os"
	"strings"

	"github.com/hashicorp/raft"
)

type CommandOperation struct {
	Operation string `json:"op,omitempty"`
	Key       string `json:"key,omitempty"`
	Value     string `json:"value,omitempty"`
}

// KvFSM raft.FSM implementation
type kvFSM struct {
}

func (k kvFSM) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		var logEntry = CommandOperation{}
		if err := json.Unmarshal(log.Data, &logEntry); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error marshalling store payload %s\n", err.Error())
			return nil
		}

		op := strings.ToUpper(strings.TrimSpace(logEntry.Operation))
		switch op {
		case "SET":
			key := logEntry.Key
			value := logEntry.Value
			sysout.Print(fmt.Sprintf("Create key %s with value %s", key, value))
			return k.applyValue(logEntry.Key, logEntry.Value)
		case "GET":
			return nil

		case "DELETE":
			return nil
		}
	}

	_, _ = fmt.Fprintf(os.Stderr, "not raft log command type\n")
	return nil
}

func (k kvFSM) Snapshot() (raft.FSMSnapshot, error) {
	sysout.Print("calling snapshot")
	return newDoNothingSnapshot()
}

func (k kvFSM) Restore(rClose io.ReadCloser) error {
	return nil
}

func (k *kvFSM) applyValue(key, value string) error {
	log.Println("setting key=  ", key)
	return nil
}

// NewKvFSM raft.FSM implementation
func NewKvFSM() raft.FSM {
	return &kvFSM{}
}
