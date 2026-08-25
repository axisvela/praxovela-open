// +build ignore

// S7 Lifecycle Example: fetch → validate → load
// Demonstrates the complete fetch/validate/load lifecycle via HTTP API.
//
// This shell script demonstrates:
// 1. POST /fetch - Create a FETCHED instance
// 2. POST /validate - Transition to VALIDATED state
// 3. POST /load - Transition to ACTIVE state
//
// Usage:
//   ./s7-lifecycle.sh
//
// Requires:
// - Syndovela facade running on http://localhost:8765/syndovela
// - curl command available

#!/bin/bash set -e

SYNDOVELA_URL="http://localhost:8765/syndovela"

echo "=== S7 Lifecycle: fetch → validate → load ==="
echo ""

# Step 1: Fetch a bundle
echo "1. POST /fetch (create FETCHED instance)"
FETCH_RESPONSE=$(curl -s -X POST "$SYNDOVELA_URL/fetch" \
  -H "Content-Type: application/json" \
  -d '{"bundleId":"my-skill"}')

FETCH_INSTANCE_ID=$(echo "$FETCH_RESPONSE" | jq -r '.instanceId')
FETCH_STATE=$(echo "$FETCH_RESPONSE" | jq -r '.state')

echo "   Instance ID: $FETCH_INSTANCE_ID"
echo "   State: $FETCH_STATE"
echo ""

# Step 2: Validate the instance
echo "2. POST /validate (FETCHED → VALIDATED)"
VALIDATE_RESPONSE=$(curl -s -X POST "$SYNDOVELA_URL/validate" \
  -H "Content-Type: application/json" \
  -d "{\"instanceId\":\"$FETCH_INSTANCE_ID\"}")

VALIDATE_STATE=$(echo "$VALIDATE_RESPONSE" | jq -r '.state')

echo "   State: $VALIDATE_STATE"
echo ""

# Step 3: Load the instance
echo "3. POST /load (VALIDATED → LOADED → ACTIVE)"
LOAD_RESPONSE=$(curl -s -X POST "$SYNDOVELA_URL/load" \
  -H "Content-Type: application/json" \
  -d "{\"instanceId\":\"$FETCH_INSTANCE_ID\"}")

LOAD_STATE=$(echo "$LOAD_RESPONSE" | jq -r '.state')
LOAD_HEALTH=$(echo "$LOAD_RESPONSE" | jq -r '.health')

echo "   State: $LOAD_STATE"
echo "   Health: $LOAD_HEALTH"
echo ""

# Verification
echo "=== Verification ==="
if [ "$LOAD_STATE" = "ACTIVE" ]; then
    echo "✅ Lifecycle complete! Instance is now ACTIVE"
else
    echo "❌ Unexpected state: $LOAD_STATE"
    exit 1
fi

echo ""
echo "S7 lifecycle demonstration complete!"