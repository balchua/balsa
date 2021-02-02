package config

type Config struct {
	HttpPort   int    `mapstructure:"httpport"`
	RaftPort   int    `mapstructure:"raftport"`
	RaftVolume string `mapstructure:"volume_dir"`
	RaftNodeId string `mapstructure:raft_node_id"`
}
