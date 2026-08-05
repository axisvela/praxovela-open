# Connector Protocol

Stable contract between AXON Core / Janus Gateway and external enterprise services such as CRM, ERP, ticketing, knowledge bases, document systems, MCP servers, REST APIs, gRPC services, and vendor SDK adapters.

## Principle

AXON Core should not embed enterprise-system-specific code. It requests named connector operations through Janus Gateway. Janus Gateway enforces capability tokens, workspace scope, tenant identity, and gateway routing before calling enterprise services.

## Deployment

Connector calls may run in two common enterprise topologies:

1. Remote Core: Desktop AxisView talks to AXON Core in the data center. AXON Core/Janus calls enterprise services inside the data center.
2. Local Core + Enterprise Services: Desktop AxisView talks to local AXON Core. Local AXON Core/Janus calls enterprise services through the enterprise API Gateway.

## Transport

Priority order:

1. HTTPS JSON for MVP and gateway compatibility.
2. MCP for tool-compatible enterprise services.
3. gRPC for high-throughput internal enterprise services.
4. Vendor SDK behind Janus Gateway adapter when no stable remote protocol exists.

## Request Schema

```json
{
  "request_id": "string",
  "session_id": "string",
  "workspace_id": "string",
  "tenant_id": "string",
  "connector_id": "crm.salesforce",
  "operation": "account.search",
  "capability_token": "string",
  "input": {},
  "trace_context": {
    "trace_id": "string",
    "span_id": "string"
  }
}
```

## Response Schema

```json
{
  "request_id": "string",
  "status": "ok|denied|not_found|rate_limited|error",
  "output": {},
  "artifacts": [
    {
      "type": "json|markdown|file|link",
      "name": "string",
      "uri": "string",
      "hash": "string"
    }
  ],
  "audit": {
    "service": "string",
    "operation": "string",
    "policy_decision_id": "string"
  },
  "error": {
    "code": "string",
    "message": "string"
  }
}
```

## Gateway Requirements

Enterprise API gateways may terminate TLS, authenticate users and service accounts, route tenant traffic, apply DLP rules, enforce rate limits, and emit centralized audit events. Connectors must treat gateway decisions as external constraints and must not bypass them.

## Security Rules

- Every connector request must carry a capability token from Archon Guard.
- Connectors must pass tenant identity explicitly.
- Secrets must be resolved by the connector runtime, not by the model or prompt.
- Connector outputs must be tagged with provenance for audit and UI display.
- High-risk write operations require explicit policy approval before execution.
