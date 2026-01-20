#!/bin/bash
set -e

# Build the binary
echo "Building storage-node..."
go build -o bin/storage-node cmd/storage-node/main.go

# Create data directories
mkdir -p data/node-1 data/node-2 data/node-3

# Function to kill child processes on exit
cleanup() {
    echo "Stopping storage nodes..."
    pkill -P $$
}
trap cleanup EXIT

echo "Starting 3 storage nodes..."

# Start Node 1
./bin/storage-node --port 9101 --data-dir ./data/node-1 &
PID1=$!
echo "Node 1 started (PID $PID1) on port 9101"

# Start Node 2
./bin/storage-node --port 9102 --data-dir ./data/node-2 &
PID2=$!
echo "Node 2 started (PID $PID2) on port 9102"

# Start Node 3
./bin/storage-node --port 9103 --data-dir ./data/node-3 &
PID3=$!
echo "Node 3 started (PID $PID3) on port 9103"

echo "Cluster is running. Press Ctrl+C to stop."
wait
