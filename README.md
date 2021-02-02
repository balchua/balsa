# Purpose of this project

To learn raft concept.

To learn how to use hashicorp/raft 

Learn how to use gorilla/mux

The project currently doesnt keep the data, it is all just printing it to the console.

These are the endpoints exposed to interact with the raft cluster:

[x] `join` - join nodes to the cluster.

[x] `set` - sets some key/value pair.  currently they are not stored.

[ ] `get` - TODO

[ ] `remove` - TODO


TODO, snapshotting.
## Starting a 3 node raft cluster.

### First node contains the following environment variables

```shell
export HTTP_SERVER_PORT=8000
export RAFT_SERVER_PORT=9000
export RAFT_NODE_ID="n1"
export RAFT_VOLUME_DIR="/tmp/node1"
```
The raft data is stored in `RAFT_VOLUME_DIR`.


```
$ ./startn1.sh
```


### Second node contains the following environment variables

```shell
export HTTP_SERVER_PORT=8001
export RAFT_SERVER_PORT=9001
export RAFT_NODE_ID="n2"
export RAFT_VOLUME_DIR="/tmp/node2"
```
The raft data is stored in `RAFT_VOLUME_DIR`.


```
$ ./startn2.sh
```


### Third node contains the following environment variables

```shell
export HTTP_SERVER_PORT=8002
export RAFT_SERVER_PORT=9002
export RAFT_NODE_ID="n3"
export RAFT_VOLUME_DIR="/tmp/node3"
```
The raft data is stored in `RAFT_VOLUME_DIR`.


```
$ ./startn3.sh
```

## Joining a node

Join node 2 to node 1

`curl "http://localhost:8000/join?nodeId=n2&raftPort=9001"`

Join node 3 to node 1

`curl "http://localhost:8000/join?nodeId=n3&raftPort=9002"`

You should now have a 3 node raft cluster.

node 1 logs:

```console
$ ./startn1.sh 
Starting node1
Cleaning the raft volumes.
2021-02-02T15:16:18.228+0800 [INFO]  raft: initial configuration: index=0 servers=[]
2021-02-02T15:16:18.228+0800 [INFO]  raft: entering follower state: follower="Node at 127.0.0.1:9000 [Follower]" leader=
2021-02-02T15:16:20.170+0800 [WARN]  raft: heartbeat timeout reached, starting election: last-leader=
2021-02-02T15:16:20.170+0800 [INFO]  raft: entering candidate state: node="Node at 127.0.0.1:9000 [Candidate]" term=2
2021-02-02T15:16:20.192+0800 [DEBUG] raft: votes: needed=1
2021-02-02T15:16:20.192+0800 [DEBUG] raft: vote granted: from=n1 term=2 tally=1
2021-02-02T15:16:20.192+0800 [INFO]  raft: election won: tally=1
2021-02-02T15:16:20.192+0800 [INFO]  raft: entering leader state: leader="Node at 127.0.0.1:9000 [Leader]"
Join handler called
2021/02/02 15:16:21 node is n2
INFO[0003] Joining node % on host 127.0.0.1:%dn2         source="handler.go:28"
2021-02-02T15:16:21.544+0800 [INFO]  raft: updating configuration: command=AddStaging server-id=n2 server-addr=127.0.0.1:9001 servers="[{Suffrage:Voter ID:n1 Address:127.0.0.1:9000} {Suffrage:Voter ID:n2 Address:127.0.0.1:9001}]"
2021-02-02T15:16:21.560+0800 [INFO]  raft: added peer, starting replication: peer=n2
2021-02-02T15:16:21.564+0800 [WARN]  raft: appendEntries rejected, sending older logs: peer="{Voter n2 127.0.0.1:9001}" next=3
2021-02-02T15:16:21.567+0800 [INFO]  raft: pipelining replication: peer="{Voter n2 127.0.0.1:9001}"
Join handler called
2021/02/02 15:16:42 node is n3
INFO[0024] Joining node % on host 127.0.0.1:%dn3         source="handler.go:28"
2021-02-02T15:16:42.495+0800 [INFO]  raft: updating configuration: command=AddStaging server-id=n3 server-addr=127.0.0.1:9002 servers="[{Suffrage:Voter ID:n1 Address:127.0.0.1:9000} {Suffrage:Voter ID:n2 Address:127.0.0.1:9001} {Suffrage:Voter ID:n3 Address:127.0.0.1:9002}]"
2021-02-02T15:16:42.500+0800 [INFO]  raft: added peer, starting replication: peer=n3
2021-02-02T15:16:42.506+0800 [WARN]  raft: appendEntries rejected, sending older logs: peer="{Voter n3 127.0.0.1:9002}" next=3
2021-02-02T15:16:42.511+0800 [INFO]  raft: pipelining replication: peer="{Voter n3 127.0.0.1:9002}"
```

Logs from Node 2

```console
$ ./startn2.sh 
Starting node2
Cleaning the raft volumes.
2021-02-02T15:14:29.063+0800 [INFO]  raft: initial configuration: index=0 servers=[]
2021-02-02T15:14:29.064+0800 [INFO]  raft: entering follower state: follower="Node at 127.0.0.1:9001 [Follower]" leader=
2021-02-02T15:14:30.136+0800 [WARN]  raft: heartbeat timeout reached, starting election: last-leader=
2021-02-02T15:14:30.136+0800 [INFO]  raft: entering candidate state: node="Node at 127.0.0.1:9001 [Candidate]" term=2
2021-02-02T15:14:30.155+0800 [DEBUG] raft: votes: needed=1
2021-02-02T15:14:30.155+0800 [DEBUG] raft: vote granted: from=n2 term=2 tally=1
2021-02-02T15:14:30.155+0800 [INFO]  raft: election won: tally=1
2021-02-02T15:14:30.155+0800 [INFO]  raft: entering leader state: leader="Node at 127.0.0.1:9001 [Leader]"
2021-02-02T15:16:21.563+0800 [WARN]  raft: failed to get previous log: previous-index=3 last-index=2 error="log not found"
2021-02-02T15:16:21.564+0800 [INFO]  raft: entering follower state: follower="Node at 127.0.0.1:9001 [Follower]" leader=127.0.0.1:9000


```

Logs from node 3

```console
$ ./startn3.sh 
Starting node3
Cleaning the raft volumes.
2021-02-02T15:14:32.185+0800 [INFO]  raft: initial configuration: index=0 servers=[]
2021-02-02T15:14:32.185+0800 [INFO]  raft: entering follower state: follower="Node at 127.0.0.1:9002 [Follower]" leader=
2021-02-02T15:14:33.796+0800 [WARN]  raft: heartbeat timeout reached, starting election: last-leader=
2021-02-02T15:14:33.796+0800 [INFO]  raft: entering candidate state: node="Node at 127.0.0.1:9002 [Candidate]" term=2
2021-02-02T15:14:33.807+0800 [DEBUG] raft: votes: needed=1
2021-02-02T15:14:33.807+0800 [DEBUG] raft: vote granted: from=n3 term=2 tally=1
2021-02-02T15:14:33.807+0800 [INFO]  raft: election won: tally=1
2021-02-02T15:14:33.807+0800 [INFO]  raft: entering leader state: leader="Node at 127.0.0.1:9002 [Leader]"
2021-02-02T15:16:42.506+0800 [WARN]  raft: failed to get previous log: previous-index=4 last-index=2 error="log not found"
2021-02-02T15:16:42.506+0800 [INFO]  raft: entering follower state: follower="Node at 127.0.0.1:9002 [Follower]" leader=127.0.0.1:9000
```

### Setting values

```
$ ./pump.sh

```

This will send several curl requests to `/set?key=$i&value=$i`.  Setting values must only be done on the leader.

If you send a set request to the follower, it will log 

```console
Set handler called
2021/02/02 16:11:48 Not a leader, do nothing
```