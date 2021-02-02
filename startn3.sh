#!/usr/bin/env bash

echo "Starting node3"
echo "Cleaning the raft volumes."
rm -rf /tmp/node3
mkdir -p /tmp/node3

export HTTP_SERVER_PORT=8002
export RAFT_SERVER_PORT=9002
export RAFT_NODE_ID="n3"
export RAFT_VOLUME_DIR="/tmp/node3"

./balsa