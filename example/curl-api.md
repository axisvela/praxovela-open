# Calling the axond API with curl

Minimal, accurate curl patterns against the **axond** (AXON Core) HTTP API. Full contract: [`../docs/api.md`](../docs/api.md).

**Base URL**: `http://localhost:8420` (configurable via `AXON_PORT`). When `AXON_API_KEY` is set, add `-H "Authorization: Bearer <token>"`.

## 1. Health

```bash
curl http://localhost:8420/health
```

```json
{"status":"ok","runtime":"AXON Core","version":"0.10.0"}
```

## 2. Run Center — list runs

```bash
curl "http://localhost:8420/v1/runcenter/runs?status=running"
```

```json
{"runs":[{"id":"...","workspace_id":"...","title":"...","category":"interactive","status":"running","actor_id":"human:user","created_at":"...","updated_at":"...","needs_review":false}]}
```

Filter with `category` (`interactive|autonomous|scheduled|watchdog|process`), `status` (`pending|running|completed|failed|cancelled|paused`), `workspace_id`, and `needs_review`.

## 3. Knowledge / memory search

```bash
curl "http://localhost:8420/v1/knowledge/search?q=agent+runtime"
```

```json
{"results":[{"id":"...","title":"...","snippet":"..."}]}
```

> Memory facts are also reachable per workspace: `GET /v1/workspaces/{id}/facts`. Index documents first with the `knowledge.index` tool or `POST /v1/knowledge`.

## 4. Bonus — create a session and start a run

```bash
curl -X POST http://localhost:8420/v1/sessions \
  -H "Content-Type: application/json" \
  -d '{"workspace_id":"..."}'

curl -X POST http://localhost:8420/v1/sessions/{id}/runs \
  -H "Content-Type: application/json" \
  -d '{"message":"Hello PRAXOVELA"}'
```

Stream run events over SSE: `GET /v1/sessions/{id}/events`.
