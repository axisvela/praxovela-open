# AxisLink — UI Protocol

Stable contract between L6 AxisView (UI) and L5 AXON Core (Agent Runtime).

## Transport

Local-first and enterprise transport (priority order):
1. HTTP + SSE (MVP)
2. WebSocket
3. Unix domain socket / Windows named pipe (Enterprise)

AxisLink is the only protocol AxisView needs to talk to AXON Core. The target may be a local sidecar (`127.0.0.1`), a remote AXON Core service, or an enterprise API Gateway that forwards to AXON Core.

## Deployment Modes

| Mode | AxisView Location | AXON Core Location | Gateway |
|------|-------------------|--------------------|---------|
| Personal Local-First | Desktop | Desktop sidecar | none by default |
| Enterprise Remote Core | Desktop | Enterprise data center | optional between AxisView and AXON Core |
| Enterprise Local Core + Enterprise Services | Desktop | Desktop sidecar | between AXON/Janus and enterprise services |

## Event Envelope

The stable AxisLink target envelope for this remediation is:

```json
{
  "type": "run.started",
  "tenant_id": "tenant-1",
  "workspace_id": "workspace-1",
  "session_id": "s_123",
  "run_id": "r_123",
  "data": {},
  "timestamp": "2026-06-20T00:00:00Z"
}
```

Current OSS SSE events include these top-level fields when available. Tenant/workspace/run context is populated from the session during runtime propagation. Local anonymous/API-key mode may use the default tenant.

## Event Types

- `connected` — SSE connection established
- `run.started` — Agent run begins
- `status.changed` — Run or session status changed
- `activity.step` — Activity timeline step emitted
- `plan.generated` — Agent emits a plan
- `thought` — Agent reasoning/thought update
- `tool.call_requested` — Agent requests tool invocation
- `tool.call_completed` — Tool execution completes
- `tool.verification_failed` — Tool result verification failed
- `policy.check` — Capability or policy check completed
- `approval.required` — User approval needed
- `sandbox.created` — Sandbox instance created
- `sandbox.stdout` — Sandbox stdout
- `sandbox.stderr` — Sandbox stderr
- `sandbox.exited` — Sandbox execution completed
- `message` — User or assistant message update
- `error` — Error event
- `run.completed` — Successful terminal event emitted by current agent runs
- `run.finished` — Generic/compatibility finish event that may carry final status
- `run.failed` — Lifecycle/compatibility name for failed runs; may appear in Run Center or future/compatibility streams, not every current session SSE path
- `run.cancelled` — Lifecycle/compatibility name for cancelled runs; may appear in Run Center or future/compatibility streams, not every current session SSE path
- `done` — Bridge compatibility event that may follow server-terminal events to close the stream
