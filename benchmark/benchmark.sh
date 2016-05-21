#!/usr/bin/env bash
#
# Simple benchmark test suite
#
# You must have installed vegeta:
# go get github.com/tsenart/vegeta
#

# Default port to listen
port=8080

# Compile test server
go build server.go
chmod +x server

# Run test server
./server & > /dev/null
server=$!

suite() {
  sleep 2

  echo "------------------------------------------"
  echo "Running '$1' benchmark with concurrency $2"

  # Start the server
  go run ./$1/$1.go & > /dev/null
  pid=$!

  echo "GET http://localhost:$port" | vegeta attack \
    -duration=20s \
    -rate=$2 | vegeta report

  # Kill proxy server process
  kill -9 $pid
}

# Run suites
suite "simple" "50"

# Kill test server process
kill -9 $server 2> /dev/null
