#!/bin/bash

set -e

# Always run from repo root
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR/.."

# Cleanup: kill background processes on exit or Ctrl+C
cleanup() {
  echo ""
  echo "Shutting down..."
  kill "$API_PID" "$WORKER_PID" 2>/dev/null
  wait "$API_PID" "$WORKER_PID" 2>/dev/null
  echo "Done."
}
trap cleanup INT TERM

echo "Virtus — Starting API + Worker"

echo ""
echo "Starting API server (http://localhost:8080)..."
go run ./cmd/api &
API_PID=$!

echo "Starting background worker..."
go run ./cmd/worker &
WORKER_PID=$!

echo ""
echo "Both processes are running."
echo "  API PID:    $API_PID"
echo "  Worker PID: $WORKER_PID"
echo ""
echo "Press Ctrl+C to stop both."
echo ""

# Wait for either process to exit
wait -n "$API_PID" "$WORKER_PID"

# If one exits unexpectedly, kill the other
echo "One process exited. Stopping all..."
cleanup
