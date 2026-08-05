# AxisRobo Agent — API Reference

**Base URL**: `http://localhost:8420` (configurable via `AXON_PORT`)
**Version**: 0.5.0

## HTTP API

### Health
```
GET /health → {"status":"ok","runtime":"AXON Core","version":"0.3.0"}
```

### Sessions
```
POST /v1/sessions                → Create session (optional workspace field)
GET  /v1/sessions                → List all sessions
GET  /v1/sessions/{id}           → Get session details
PUT  /v1/sessions/{id}/model     → Persist selected session model
POST /v1/sessions/{id}/messages  → Send message; @agent-id creates child session
POST /v1/sessions/{id}/runs      → Start agent run
GET  /v1/sessions/{id}/events    → SSE event stream
GET  /v1/sessions/{id}/children  → List child sessions
POST /v1/runs/{id}/cancel        → Cancel run
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

### Agents (v0.3.0)
```
GET  /v1/agents                  → List agents (?mode=primary|subagent|all&include_hidden=true)
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

Binary `files[]` are decoded by AXON Core, saved under `<workspace>/.axisrobo/uploads/`, and appended to the user message as a `doc.parse`-ready saved path. Image attachments are validated by count, MIME type, and base64 size before being sent as model content blocks.

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

### Provider Config
```
GET  /v1/config                    → System config (port, remote, auth)
GET  /v1/config/provider           → Current provider + available providers
POST /v1/config/provider           → Switch provider (provider, api_key, model)
GET  /v1/config/providers          → List saved providers (masked keys)
```

### Provider & Model Management

#### Providers

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/providers` | List all providers |
| POST | `/v1/providers` | Add custom provider |
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
| PUT | `/v1/models/{id}/enable` | Enable or disable a model |
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

### Agents
```
GET /v1/agents      → List agents (build, plan)
GET /v1/agents/{id} → Get agent details
```

### Document Tools

| Tool | Description |
|------|-------------|
| `doc.parse` | Parse `.docx`, `.pptx`, `.pdf`, `.xlsx`, `.md`, or `.txt` into structured Markdown using the parser chain. |
| `doc.toMarkdown` | Parse a document and optionally write the Markdown result to a workspace-scoped file. |
| `doc.generate` | Generate a document from Markdown into `docx`, `pdf`, `pptx`, `md`, or `txt` when the backend can truthfully produce that format. |
| `doc.edit` | Regenerate a document from fully edited Markdown. |
| `doc.new` | Create a document from a built-in template (`resume`, `report`, `letter`). |
| `doc.render` | One-step workspace-aware document render into `<workspace>/.axisrobo/outputs/`. |
| `doc.open` | Open an existing file with the OS default application. |

`doc.parse` accepts either `path` or `file_path` for compatibility with generated tool schemas. Tool declarations such as `<function_declaration>` are metadata only; executable calls should use JSON, plain `tool.name { ... }`, DeepSeek DSML, or XML-style `<tool_call name="...">...</tool_call>`.

`doc.render` / `doc.generate` require `pandoc` for real `.docx`, `.pdf`, and `.pptx` output. PDF output also requires a LaTeX engine such as `xelatex`. If those dependencies are missing, the tool returns an explicit error rather than writing fake binary output.

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
| GET | `/v1/runcenter/runs/{id}/events/stream` | SSE event stream (stub, returns 501) |
| POST | `/v1/runcenter/runs/{id}/cancel` | Cancel a run |

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

### Audit
```
GET /v1/trace/public-key → Ed25519 public key + chain integrity
```

### CopilotKit
```
POST /api/copilotkit/chat/completions → OpenAI-compatible streaming
GET  /api/copilotkit/tools             → Available tools
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

### Workspace File Tree (New)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/workspaces/{id}/files?depth=3` | Workspace file tree |

### Auth & Users

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/v1/auth/register` | None | Register user |
| POST | `/v1/auth/login` | None | Login, returns JWT |
| POST | `/v1/auth/refresh` | JWT | Refresh token |
| GET | `/v1/auth/oidc/login` | None | OIDC SSO redirect (EE only) |
| GET | `/v1/auth/oidc/callback` | None | OIDC callback → JWT (EE only) |
| GET | `/v1/users/me` | JWT | Current user info |

**OIDC (EE only):** Enable via `AXOND_OIDC_ENABLED=true` + `AXOND_OIDC_ISSUER`, `AXOND_OIDC_CLIENT_ID`.

Enterprise/governed mode requires protected endpoints to use `Authorization: Bearer <jwt>` and is intended to scope remote requests to the authenticated tenant and workspace context. Personal local OSS mode allows anonymous access with `user_id=anonymous` or API-key access as local-admin behavior.

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

### Background Jobs (New)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/jobs` | List background jobs |
| GET | `/v1/jobs/{id}` | Get job status |
| POST | `/v1/jobs/{id}/cancel` | Cancel job |

## Tauri Commands

```typescript
invoke("get_agent_url") → "http://localhost:8420"
invoke("get_health")    → {"status":"ok","runtime":"AXON Core"}
invoke("get_active_window") → "Code|main.go"
```
