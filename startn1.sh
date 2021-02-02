#!/usr/bin/env bash

echo "Starting node1"
echo "Cleaning the raft volumes."

rm -rf /tmp/node1
mkdir -p /tmp/node1

export HTTP_SERVER_PORT=8000
export RAFT_SERVER_PORT=9000
export RAFT_NODE_ID="n1"
export RAFT_VOLUME_DIR="/tmp/node1"

./balsa