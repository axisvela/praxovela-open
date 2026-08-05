#!/bin/bash
set -e

echo "=========================================="
echo "  AxisRobo Agent - Development Launcher"
echo "=========================================="

# Build AXON Core
echo "[BUILD] Building axond..."
go build -o apps/desktop/src-tauri/binaries/axond ./packages/axon-core/cmd/axond/
echo "[OK] axond built"

# Start AXON Core
echo "[START] Starting AXON Core on :8420..."
./apps/desktop/src-tauri/binaries/axond &
AXON_PID=$!
echo "[OK] AXON Core PID: $AXON_PID"

# Wait for health
echo "[WAIT] Waiting for health check..."
for i in $(seq 1 15); do
    if curl -s http://localhost:8420/health > /dev/null 2>&1; then
        echo "[OK] AXON Core healthy"
        break
    fi
    sleep 1
done

echo ""
echo "=========================================="
echo "  AXON Core running on http://localhost:8420"
echo "  Health: http://localhost:8420/health"
echo "  Stop:   kill $AXON_PID"
echo "=========================================="

# Trap to kill axond on exit
trap "echo 'Stopping AXON Core...'; kill $AXON_PID 2>/dev/null" EXIT

# Keep running
wait $AXON_PID
