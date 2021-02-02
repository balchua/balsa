#!/usr/bin/env bash

echo "Starting node2"
echo "Cleaning the raft volumes."

rm -rf /tmp/node2
mkdir -p /tmp/node2

export HTTP_SERVER_PORT=8001
export RAFT_SERVER_PORT=9001
export RAFT_NODE_ID="n2"
export RAFT_VOLUME_DIR="/tmp/node2"

./balsa