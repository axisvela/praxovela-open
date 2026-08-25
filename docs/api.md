# PRAXOVELA — API Reference

**Base URL**: `http://localhost:8420` (configurable via `AXON_PORT`)
**Version**: 1.9.0

## HTTP API

### Health
```
GET /health     → {"status":"ok","runtime":"AXON Core","version":"1.9.0"}
GET /api/status → alias of /health
```

### Sessions
```
POST /v1/sessions                → Create session (optional workspace field)
GET  /v1/sessions                → List all sessions
GET  /v1/sessions/{id}           → Get session details
PUT  /v1/sessions/{id}           → Update session
PUT  /v1/sessions/{id}/model     → Persist selected session model
GET  /v1/sessions/{id}/messages  → List session messages
POST /v1/sessions/{id}/messages  → Send message; @agent-id creates child session
POST /v1/sessions/{id}/runs      → Start agent run
GET  /v1/sessions/{id}/events    → SSE event stream
GET  /v1/sessions/{id}/children  → List child sessions
POST /v1/runs/{id}/cancel        → Cancel run
POST /v1/runs/{id}/recover       → Restore a run from its latest checkpoint
```

**POST /v1/runs/{id}/recover** — restores a run (and its owning session) from the
latest checkpoint via the effect-cursor-verified recovery path. Query param
`session_id` is optional; when omitted the owning session is resolved from the
checkpoint. Tenant-scoped requests for a run owned by another tenant are
reported as 404. Recovery is best-effort: a missing checkpoint or a failed
effect-cursor verification returns 200 with `restored: false` and a `reason`.

```json
// Response (restored):
{ "checkpoint_id": "cp-...", "session_id": "s_123", "run_id": "r_123",
  "restored": true, "message_count": 4, "cursor_verified": true,
  "cursor": { "run_id": "r_123", "last_effect_id": "fx-...", "last_created_at": 0 } }
// Response (nothing to restore):
{ "session_id": "", "run_id": "unknown", "restored": false,
  "message_count": 0, "cursor_verified": false, "reason": "no checkpoint found" }
```

`GET /v1/sessions/{id}/events` emits Server-Sent Events. The stable AxisLink target envelope for this remediation is:

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

Current OSS SSE events include `tenant_id`, `workspace_id`, `run_id`, and `timestamp` when available. Tenant/workspace/walkspace/run context is populated from the session during runtime propagation. JWT-authenticated Enterprise and remote mode must scope requests by token claims and request context; local OSS anonymous or API-key mode remains local-admin behavior for the local runtime.

**POST /v1/sessions/{id}/messages** — @mention dispatch:

```json
// Request: {"message":"@coding-review Review the diff"}
// Response when mention resolves:
{ "child_session_id": "...", "run_id": "...", "agent_id": "coding-review" }
// Response when mention fails: 422 {"error": "agent not found or not invokable"}
// Normal message (no @) starts a parent run as with /runs.
```

**POST /v1/sessions/{id}/runs** request body:

```json
{
  "message": "总结这个论文",
  "images": [
    { "name": "diagram.png", "data": "data:image/png;base64,...", "mime_type": "image/png" }
  ],
  "files": [
    { "name": "paper.docx", "data": "data:application/vnd.openxmlformats-officedocument.wordprocessingml.document;base64,...", "mime_type": "application/vnd.openxmlformats-officedocument.wordprocessingml.document" }
  ]
}
```

Binary `files[]` are decoded by AXON Core, saved under `<workspace>/.praxovela/uploads/`, and appended to the user message as a `doc.parse`-ready saved path. Image attachments are validated by count, MIME type, and base64 size before being sent as model content blocks.

### Agents
```
GET /v1/agents                  → List agents (?mode=primary|subagent|all&include_hidden=true)
GET /v1/agents/{id}             → Get agent details
POST /v1/agents/generate        → Generate a new agent
POST /v1/agent/tools/execute    → Execute a tool directly (SUT driver entry point, Member)
```

### Workspaces
```
GET    /v1/workspaces             → List workspaces
POST   /v1/workspaces             → Create (name, profile, path, description, persona)
GET    /v1/workspaces/{id}        → Get workspace
PUT    /v1/workspaces/{id}        → Update workspace name/profile/path
DELETE /v1/workspaces/{id}        → Delete workspace
GET    /v1/workspaces/{id}/sessions → Sessions in workspace
GET    /v1/workspaces/{id}/files    → Workspace file tree (?depth=3)
GET    /v1/workspaces/{id}/provider → Workspace's LLM provider
```

### Workspace Sub-resources

#### Snapshots
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/workspaces/{id}/snapshots` | List workspace snapshots |
| POST | `/v1/workspaces/{id}/snapshots` | Save a snapshot |
| GET | `/v1/workspaces/{id}/snapshots/{snapId}` | Get a snapshot |
| DELETE | `/v1/workspaces/{id}/snapshots/{snapId}` | Delete a snapshot |

#### Plans
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/workspaces/{id}/plans` | List plans |
| POST | `/v1/workspaces/{id}/plans` | Save a plan |
| GET | `/v1/workspaces/{id}/plans/{planId}` | Get a plan |
| DELETE | `/v1/workspaces/{id}/plans/{planId}` | Delete a plan |

#### Storage
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/workspaces/{id}/storage` | Get workspace storage |
| POST | `/v1/workspaces/{id}/storage` | Set workspace storage |
| GET | `/v1/workspaces/{id}/storage/keys` | List storage keys |

#### Todos
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/workspaces/{id}/todos` | Get todos |
| POST | `/v1/workspaces/{id}/todos` | Set todos |

#### Sync
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/workspaces/{id}/sync` | Get sync state |
| POST | `/v1/workspaces/{id}/sync` | Set sync state |
| PUT | `/v1/workspaces/{id}/sync/merge` | Merge sync state |

### Provider Config
```
GET  /v1/config                    → System config (port, remote, auth)
GET  /v1/config/provider           → Current provider + available providers
POST /v1/config/provider           → Switch provider (provider, api_key, model)
GET  /v1/config/providers          → List saved providers (masked keys)
GET  /v1/config/files              → List config files (masked secrets)
GET  /v1/config/profiles           → Get current profile
POST /v1/config/profiles           → Save profile
```

### Provider & Model Management

#### Providers

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/providers` | List all providers |
| POST | `/v1/providers` | Add custom provider |
| POST | `/v1/providers/check` | Check provider connectivity |
| POST | `/v1/providers/{id}/verify` | Verify saved provider API key against `/models` |
| PUT | `/v1/providers/{id}` | Update provider config |
| DELETE | `/v1/providers/{id}` | Delete provider |

**POST /v1/providers** request body:
```json
{
  "id": "custom-provider",
  "name": "My Provider",
  "type": "openai-compatible",
  "baseURL": "https://api.example.com/v1",
  "api_key": "sk-...",
  "headers": {},
  "authMethod": "api-key"
}
```

#### Models

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/models` | List all models with enable status |
| GET | `/v1/models/usable` | List models enabled and backed by verified providers |
| POST | `/v1/models/check` | Check model connectivity |
| PUT | `/v1/models/{id}/enable` | Enable or disable a model |
| PUT | `/v1/models/{p}/{m}/enable` | Enable or disable a model by provider/model path |
| PUT | `/v1/models/{id}/options` | Update model default options |

**PUT /v1/models/{id}/enable** request body:
```json
{ "enabled": true }
```

**PUT /v1/models/{id}/options** request body:
```json
{
  "options": { "temperature": 0.7, "reasoningEffort": "medium" },
  "variants": { "high": { "reasoningEffort": "high" }, "low": { "reasoningEffort": "low" } }
}
```

#### Workspace Model Selection

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/workspaces/{id}/model` | Get selected model |
| PUT | `/v1/workspaces/{id}/model` | Switch workspace model |

**PUT /v1/workspaces/{id}/model** request body:
```json
{ "model": "openai/gpt-4o" }
```

### Knowledge Base

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/knowledge` | List all knowledge entries |
| POST | `/v1/knowledge` | Create knowledge entry |
| DELETE | `/v1/knowledge/{id}` | Delete knowledge entry |

### Knowledge (New)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/knowledge/documents` | List indexed documents |
| GET | `/v1/knowledge/search?q=...` | Search indexed knowledge |
| GET | `/v1/knowledge/graph?path=...` | Get document tree graph |

### Knowledge Tools (New)

| Tool | Description |
|------|-------------|
| `knowledge.index` | Index a document file for retrieval (`path` or `file_path`) |
| `knowledge.search` | Search indexed documents (`query`, `top_k`) |
| `knowledge.list` | List all indexed documents |

### Document Tools

| Tool | Description |
|------|-------------|
| `doc.parse` | Parse `.docx`, `.pptx`, `.pdf`, `.xlsx`, `.md`, or `.txt` into structured Markdown using the parser chain. |
| `doc.toMarkdown` | Parse a document and optionally write the Markdown result to a workspace-scoped file. |
| `doc.generate` | Generate a document from Markdown into `docx`, `pdf`, `pptx`, `md`, or `txt` when the backend can truthfully produce that format. |
| `doc.edit` | Regenerate a document from fully edited Markdown. |
| `doc.new` | Create a document from a built-in template (`resume`, `report`, `letter`). |
| `doc.render` | One-step workspace-aware document render into `<workspace>/.praxovela/outputs/`. |
| `doc.open` | Open an existing file with the OS default application. |

`doc.parse` accepts either `path` or `file_path` for compatibility with generated tool schemas. Tool declarations such as `<function_declaration>` are metadata only; executable calls should use JSON, plain `tool.name { ... }`, DeepSeek DSML, or XML-style `<tool_call name="...">...</tool_call>`.

`doc.render` / `doc.generate` require `pandoc` for real `.docx`, `.pdf`, and `.pptx` output. PDF output also requires a LaTeX engine such as `xelatex`. If those dependencies are missing, the tool returns an explicit error rather than writing fake binary output.

### Workspace Personality & Goals

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/workspaces/{id}/personality` | Get workspace personality |
| PUT | `/v1/workspaces/{id}/personality` | Set workspace personality |
| GET | `/v1/workspaces/{id}/goals` | List workspace goals |
| POST | `/v1/workspaces/{id}/goals` | Add goal |
| PUT | `/v1/workspaces/{id}/goals/{goalId}` | Update goal |
| DELETE | `/v1/workspaces/{id}/goals/{goalId}` | Delete goal |

### Workspace Memory

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/workspaces/{id}/facts` | List memory facts |
| POST | `/v1/workspaces/{id}/facts` | Add memory fact |

**POST /v1/workspaces/{id}/facts** request body:
```json
{ "content": "This project uses React 18 with TypeScript" }
```

### Profiles
```
GET /v1/profiles → List all profile types
```

### Run Center

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/runcenter/runs` | List runs (filtered) |
| GET | `/v1/runcenter/runs/{id}` | Get run detail + events |
| GET | `/v1/runcenter/runs/{id}/events` | List events for a run |
| GET | `/v1/runcenter/runs/{id}/events/stream` | SSE event stream (cursor-polling, tenant-scoped) |
| POST | `/v1/runcenter/runs/{id}/cancel` | Cancel a run |
| GET | `/v1/evidence/runs/{id}` | Run evidence reference graph |

**GET /v1/runcenter/runs** query parameters:

| Parameter | Type | Description |
|-----------|------|-------------|
| `category` | string | Filter by category (`interactive`, `autonomous`, `scheduled`, `watchdog`, `process`) |
| `status` | string | Filter by status (`pending`, `running`, `completed`, `failed`, `cancelled`, `paused`) |
| `workspace_id` | string | Filter by workspace ID |
| `needs_review` | bool | Filter by review flag (`true` or `false`) |

**Response shapes:**

`GET /v1/runcenter/runs`:
```json
{ "runs": [{ "id": "...", "workspace_id": "...", "title": "...", "category": "interactive", "status": "running", "actor_id": "human:user", "created_at": "...", "updated_at": "...", "needs_review": false }] }
```

`GET /v1/runcenter/runs/{id}`:
```json
{ "run": { "id": "...", "status": "running", ... }, "events": [{ "id": "...", "run_id": "...", "type": "run.created", ... }] }
```

`GET /v1/runcenter/runs/{id}/events`:
```json
{ "events": [{ "id": "...", "run_id": "...", "sequence": 1, "type": "run.created", "payload": {}, "created_at": "..." }] }
```

`POST /v1/runcenter/runs/{id}/cancel`:
```json
{ "status": "cancelled" }
```

### Approvals
```
GET  /v1/approvals                 → Pending approvals
POST /v1/approvals/{id}/response   → Respond (approved/denied)
```

### Audit & Trace
```
GET /v1/trace/public-key       → Ed25519 public key + chain integrity
GET /v1/trace/sessions/{id}    → Durable trace for a session
```

### Grant / Capability Revocation

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/v1/grants/{handle}/revoke` | Admin | Revoke a grant by handle (fail-closed, sticky deny) |
| POST | `/v1/capabilities/{id}/revoke` | Admin | Revoke a capability across sessions (persistent) |

### Effect Reconciliation

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/v1/effects/reconcile` | Admin | List unknown-effect review queue |
| POST | `/v1/effects/{id}/unknown` | Admin | Mark an effect unknown (human review) |
| POST | `/v1/effects/{id}/compensate` | Admin | Compensate a committed effect (idempotent) |

### Auth & Users

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/v1/auth/register` | None | Register user |
| POST | `/v1/auth/login` | None | Login, returns JWT |
| POST | `/v1/auth/refresh` | JWT | Refresh token |
| GET | `/v1/auth/oidc/login` | None | OIDC SSO redirect (EE only) |
| GET | `/v1/auth/oidc/callback` | None | OIDC callback → JWT (EE only) |
| GET | `/v1/users/me` | JWT | Current user info |
| GET | `/v1/users` | Admin | List users |
| PATCH | `/v1/users/{id}/role` | Admin | Update user role |
| GET | `/v1/tenants` | - | List tenants |
| GET | `/v1/tenants/{id}` | - | Get tenant |

**OIDC (EE only):** Enable via `AXOND_OIDC_ENABLED=true` + `AXOND_OIDC_ISSUER`, `AXOND_OIDC_CLIENT_ID`.

Enterprise/governed mode requires protected endpoints to use `Authorization: Bearer <jwt>` and is intended to scope remote requests to the authenticated tenant and workspace context. Personal local OSS mode allows anonymous access with `user_id=anonymous` or API-key access as local-admin behavior.

### Scheduler

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/scheduler/schedules` | List scheduled tasks |
| POST | `/v1/scheduler/schedules` | Create schedule `{"id":"...","cron_expr":"...","message":"...","workspace_id":"..."}` |
| POST | `/v1/scheduler/schedules/{id}/toggle` | Enable/disable schedule `{"enabled":true}` |

### Autonomy

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/autonomy/goals` | Create autonomy goal |
| GET | `/v1/autonomy/goals` | List all autonomy goals |
| GET | `/v1/autonomy/goals/{id}` | Get single goal detail |
| POST | `/v1/autonomy/goals/{id}/cancel` | Cancel goal (marks failed + stops runner) |
| POST | `/v1/autonomy/goals/{id}/pause` | Pause goal |
| POST | `/v1/autonomy/goals/{id}/resume` | Resume paused goal |
| GET | `/v1/autonomy/twins?user_id=...` | List user's digital twins |
| PUT | `/v1/autonomy/twins/{id}` | Update twin (persona, max_risk, name) |

**POST /v1/autonomy/goals** request body:
```json
{
  "tenant_id": "...",
  "workspace_id": "...",
  "user_id": "...",
  "twin_id": "...",
  "title": "Goal title",
  "prompt": "Goal prompt"
}
```

**PUT /v1/autonomy/twins/{id}** request body:
```json
{
  "persona": "Agent persona description",
  "max_risk": "medium",
  "name": "Twin Name"
}
```

### Background Jobs

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/jobs` | List background jobs |
| GET | `/v1/jobs/{id}` | Get job status |
| POST | `/v1/jobs/{id}/cancel` | Cancel job |

### Commands & Project
```
GET  /v1/commands          → List available commands
GET  /v1/commands/{name}   → Get a command
POST /v1/project/detect    → Detect project type
```

### Format
```
GET  /v1/format            → List formatters
POST /v1/format            → Format a file
```

### Share
```
POST /v1/share/export      → Export a share
POST /v1/share/import      → Import a share
GET  /v1/share             → List shares
```

### Web & Shell
```
POST /v1/web/fetch   → Fetch a URL
POST /v1/web/search  → Web search
POST /v1/shell/run   → Run a shell command (member-porous local operator)
```

### Environment
```
GET /v1/env   → List environment information
```

### Git
```
GET  /v1/git/status    → Git status
GET  /v1/git/diff-base → Diff against base branch
POST /v1/git/pr        → Create a pull request (Admin)
```

### LSP
```
GET  /v1/lsp/servers      → List LSP servers
POST /v1/lsp/diagnostics  → Run diagnostics
```

### Image
```
POST /v1/image/generate   → Generate an image
```

### Schema
```
POST /v1/schema/validate  → Validate a schema
```

### Message Bus
```
POST /v1/bus/emit         → Emit a bus event (Admin)
GET  /v1/bus/history      → Bus history
GET  /v1/bus/health       → Bus health
POST /v1/bus/enqueue      → Enqueue a message (Member)
GET  /v1/bus/dequeue      → Dequeue a message (Member)
```

### Plugins
```
GET  /v1/plugins               → List plugins
POST /v1/plugins               → Register plugin (Admin)
DELETE /v1/plugins/{name}      → Delete plugin (Admin)
PUT  /v1/plugins/{name}/toggle

### Plugin Governance (SYNDOVELA skill plane)
```
GET  /v1/governance/plugins                     → unified catalog (id/version/kind/source/state/invocable/digest)
POST /v1/governance/plugins/{id}/disable        → close mediation gate (Member)
POST /v1/governance/plugins/{id}/enable         → restore (Member); QUARANTINED → 409 approval required
POST /v1/governance/plugins/{id}/quarantine     → mark QUARANTINED (Admin)
``` → Toggle plugin (Admin)
```

### External Directory
```
POST /v1/extdir/check     → Check external directory
POST /v1/extdir/rules     → Add external directory rule (Admin)
```

### Tasks & Delegation
```
POST /v1/task/dispatch    → Dispatch a task (Admin)
POST /v1/task/delegate    → Delegate a task (Member)
GET  /v1/task/delegations → List delegations (Member)
```

### Worktrees
```
GET    /v1/worktrees        → List worktrees
POST   /v1/worktrees        → Create a worktree
DELETE /v1/worktrees        → Remove a worktree
```

### Patch
```
POST /v1/patch/apply    → Apply a patch (Admin)
```

### Install & Server Info
```
GET  /v1/server/info      → Server information
GET  /v1/install          → Install info
POST /v1/install/init     → Initialize install (Admin)
GET  /v1/install/check    → Check install
POST /v1/install/cleanup  → Clean up install (Admin)
GET  /v1/id/generate      → Generate an ID
POST /v1/truncate         → Truncate store (Admin)
POST /v1/echo             → Echo request
```

### CLI
```
GET  /v1/cli/help    → CLI help text
POST /v1/cli/parse   → Parse a CLI command
```

### Doom (session health guard)
```
GET  /v1/doom/check    → Check doom status
POST /v1/doom/record   → Record a doom event (Admin)
POST /v1/doom/reset    → Reset doom state (Admin)
```

### Tools & Memory
```
GET /v1/tools                → List available tools
GET /v1/memories/search      → Search memory
```

### Notifications (SMTP)
```
GET  /v1/notifications/smtp/config → Get SMTP config
POST /v1/notifications/smtp/config → Set SMTP config (Admin)
POST /v1/notifications/smtp/test   → Test SMTP config (Admin)
```

### CopilotKit
```
POST /api/copilotkit/chat/completions → OpenAI-compatible streaming
POST /chat/completions                → OpenAI-compatible streaming (alias)
GET  /api/copilotkit/tools            → Available tools
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `AXON_PORT` | 8420 | HTTP port |
| `AXON_ALLOW_REMOTE` | `""` | Set to `1` to listen on `0.0.0.0` |
| `AXON_API_KEY` | `""` | Bearer token auth |
| `AXON_SYSTEM_PROMPT` | `""` | Custom system prompt |
| `AXON_DB_PATH` | ./data/axon.db | SQLite path |
| `AXON_OPENAI_API_KEY` | - | OpenAI key |
| `AXON_DEEPSEEK_API_KEY` | - | DeepSeek key |
| `AXON_QWEN_API_KEY` | - | Qwen key |
| `AXON_ZHIPU_API_KEY` | - | ZhipuAI key |
| `AXON_KIMI_API_KEY` | - | Kimi key |
| `AXON_MINIMAX_API_KEY` | - | MiniMax key |
| `AXON_MODEL_NAME` | per-provider | Override model name |
| `AXON_CUSTOM_BASE_URL` | per-provider | Override API URL |
| `AXON_MCP_SERVERS` | - | MCP server config (format: `name:cmd:args`) |
| `AXON_JWT_SECRET` | - | JWT signing secret (required for enterprise mode) |
| `MNEME_API_URL` | `http://127.0.0.1:8000` | MNEME memory platform URL |

## Tauri Commands

```typescript
invoke("get_agent_url") → "http://localhost:8420"
invoke("get_health")    → {"status":"ok","runtime":"AXON Core"}
invoke("get_active_window") → "Code|main.go"
```
