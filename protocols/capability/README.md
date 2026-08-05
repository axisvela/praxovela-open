# Capability Contract

Stable contract between L5 AXON Core (Agent) and L4 Archon Guard (Governance).

## Principle

Agent does not decide its own authority. Agent requests capability; Gateway approves, denies, or constrains.

## Capability Request

```json
{
  "request_id": "string",
  "tenant_id": "string",
  "workspace_id": "string",
  "user_id": "string",
  "user_role": "owner|admin|member|viewer|service",
  "actor_id": "string",
  "session_id": "string",
  "agent_id": "string",
  "tool": "string",
  "resource": "string",
  "operation": "string",
  "risk_hint": "low|medium|high|critical",
  "reason": "string",
  "ttl_seconds": 600
}
```

## Capability Decision

```json
{
  "request_id": "string",
  "decision": "approve|approve_with_confirmation|deny|downgrade",
  "capability_token": "string",
  "constraints": {
    "sandbox": "required|optional|none",
    "network": "disabled|allowlisted|full",
    "paths": ["glob patterns"],
    "max_bytes": 1048576,
    "ttl_seconds": 600,
    "requires_audit": true
  },
  "reason": "string"
}
```

## Risk Levels

| Level | Examples | Default Action |
|-------|----------|---------------|
| L0 | summarize visible text | auto-approve |
| L1 | read allowed project file | log, auto-approve |
| L2 | write draft file | approval required |
| L3 | run test command | approval + sandbox |
| L4 | install package | explicit approval + sandbox |
| L5 | delete files, access credentials | deny by default |
