# Syndovela API Reference (Open)

## Overview
The Syndovela SBRP (Skill-Based Resource Protocol) runtime facade provides HTTP endpoints for
runtime-control plane interaction. This documentation covers the open (OSS) compatible endpoints.

## Base Path
`/sbrp` (mounted on HTTP mux or standalone server)

## Authentication
None required for basic operations. All endpoints are stateless and idempotent where possible.

## Endpoints

### GET /describe
Retrieve the runtime's self-description (capability negotiation).

**Request**: `GET /sbrp/describe`
**Response**: `RuntimeDescriptor` object

**Response Codes**:
- `200`: Successfully retrieved descriptor
- `4xx`: Client error (bad request, etc.)

**Example**:
```bash
curl -X GET http://localhost:8765/sbrp/describe
```

### POST /apply
Apply a bundle binding to the runtime.

**Request**: `POST /sbrp/apply`
**Body**: `BundleBinding` JSON object

**Response Codes**:
- `200`: Successfully applied, instance recorded ACTIVE
- `400`: Invalid bundle binding
- `404`: Unknown resolutionLockRef
- `409`: Digest mismatch with resolution lock

**Example**:
```bash
curl -X POST http://localhost:8765/sbrp/apply \
  -H "Content-Type: application/json" \
  -d '{"bundleId":"my-bundle","version":"1.0.0","digest":"sha256:abc","resolutionLockRef":"lock-1"}'
```

### GET /report
Retrieve the actual state report (honest reporting per spec rule #6).

**Request**: `GET /sbrp/report`
**Response**: `ActualStateReport` object

**Response Codes**:
- `200`: Successfully retrieved report

**Example**:
```bash
curl -X GET http://localhost:8765/sbrp/report
```

### POST /revoke
Revoke a running instance (stop new invocations, then drain).

**Request**: `POST /sbrp/revoke`
**Body**: `{"instanceId": "inst-1"}`

**Response Codes**:
- `200`: Successfully revoked
- `404`: Instance not found

**Example**:
```bash
curl -X POST http://localhost:8765/sbrp/revoke \
  -H "Content-Type: application/json" \
  -d '{"instanceId":"inst-1"}'
```

### POST /fetch (S7)
Create a FETCHED instance (beginning of fetch/validate/load lifecycle).

**Request**: `POST /sbrp/fetch`
**Body**: `{"bundleId": "bundle-a"}`

**Response Codes**:
- `200`: Successfully fetched
- `400`: Invalid fetch request

### POST /validate (S7)
Transition from FETCHED to VALIDATED state.

**Request**: `POST /sbrp/validate`
**Body**: `{"instanceId": "inst-1"}`

**Response Codes**:
- `200`: Successfully validated
- `409`: Instance not in FETCHED state

### POST /load (S7)
Transition from VALIDATED to ACTIVE (full lifecycle completion).

**Request**: `POST /sbrp/load`
**Body**: `{"instanceId": "inst-1"}`

**Response Codes**:
- `200`: Successfully loaded, instance is now ACTIVE
- `409`: Instance not in VALIDATED state

## State Machine

### S6 States: ACTIVE ↔ DRAINING ↔ STOPPED

```
ACTIVE  ──► DRAINING ──► STOPPED
     ▲              │
     └──── 收回 ──────┘
     (revoke, no in-flight invocations → immediate STOPPED)
```

### S7 States: FETCHED → VALIDATED → LOADED → ACTIVE

```
FETCHED ──► VALIDATED ──► LOADED ──► ACTIVE
     ▲              │              │
     └──── 收回 ──────────────┘
     (from FETCHED/VALIDATED → immediate STOPPED)
```

## Error Handling

All endpoints follow fail-closed semantics:
- Digest mismatch → `409 Conflict`
- Unknown lock ref → `404 Not Found`
- Invalid state transitions → `409 Conflict`
- Wrong HTTP method → `405 Method Not Allowed`

## Compatibility

- **Protocol Version**: `sbrp/v1`
- **Minimal Client**: Any HTTP client (curl, browser, test frameworks)
- **No Authentication**: Open mode (standalone or wired client)

---

## Related Documentation

- [Syndovela Architecture](https://github.com/axisrobo/praxovela/blob/v1.3.0/docs/architecture/sbrp-architecture.md)
- [SBRP Specification](https://github.com/axisrobo/praxovela/blob/v1.3.0/protocols/sbrp/)
- [Client Contract](packages/adapters/sbrp/contract.go)