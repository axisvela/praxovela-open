# PRAXOVELA

**[中文版](./README.zh-CN.md)**

> A governable, local-first desktop agent runtime for AI execution.

**PRAXOVELA** (formerly **AxisRobo Agent**, runtime: **AXON Core**) treats
governance as a first-class primitive. Agents ask for capabilities, a policy
plane decides, high-risk operations run sandboxed, and every action is appended
to an auditable, replayable trace — before anything leaves your machine.

This repository is the **public open-source face** of the PRAXOVELA ecosystem,
published under the **Apache License 2.0**. It is where third parties integrate:
client-facing API docs, stable integration contracts, examples, and build
tooling.

---

## Why PRAXOVELA

Modern agent frameworks give agents tools but too much power and too little
accountability. PRAXOVELA solves the problems that surface when agents start
doing real work:

- **Unaccountable tool use** — every tool call passes a deny-by-default
  capability gate, is logged to a durable trace, and can be replayed after a
  crash or restart.
- **Unrecoverable runs** — the effect ledger + checkpoint store let an
  interrupted run resume exactly where it left off; committed effects are never
  re-executed.
- **Sensitive context leaking to the cloud** — context stays local by default;
  remote execution is an explicit, governed choice.
- **Running untrusted model output** — skills and adapters execute in WASM /
  container / OS sandboxes, and artifacts are digest-verified before use.

## What This Repo Contains

| Path | Contents |
|------|----------|
| [`docs/api.md`](docs/api.md) | `axond` local agent API reference (client contract) |
| [`protocols/axislink/`](protocols/axislink/README.md) | AxisLink UI protocol (L6 UI ↔ L5 Agent) |
| [`protocols/capability/`](protocols/capability/README.md) | Capability contract (L5 Agent ↔ L4 Governance) |
| [`protocols/connector/`](protocols/connector/README.md) | Connector protocol (L5/L3 ↔ enterprise systems) |
| [`protocols/sandbox/`](protocols/sandbox/README.md) | Sandbox provider interface (L3 ↔ L2) |
| [`example/`](example/README.md) | Integration examples (curl against the axond API) |
| `scripts/` | Build/release tooling (references the monorepo layout) |
| `version.json` | Product version single source |

## What This Repo is NOT

No core runtime source (AGPL), no enterprise features (proprietary), and no
architecture / roadmap / planning documents. The source of truth for core code
and architecture is the **private main repository** — it is never copied here.

## Ecosystem

| Repo | License | Visibility | Contents |
|------|---------|------------|----------|
| `praxovela` | AGPL | Private | Main runtime source: AXON core, Janus Gateway, Vulcan Forge, Run Center, AxisLink, adapters |
| `praxovela-ee` | Proprietary | Private | Enterprise: Archon Guard, Argus Trace, runcenter-ext, PDP/SIEM/OIDC/Mneme |
| `praxovela-open` | Apache-2.0 | **Public** | SDK/API docs, integration contracts, examples, build tooling |

## Quick Start for Integrators

1. **Run the runtime** — install the PRAXOVELA desktop app, or build/run
   `axond` (AXON Core). Default endpoint `http://localhost:8420` (override with
   `AXON_PORT`).
2. **Verify it is alive**:
   ```bash
   curl http://localhost:8420/health
    # → {"status":"ok","runtime":"AXON Core","version":"1.4.0"}
   ```
3. **Read the client contract** in [`docs/api.md`](docs/api.md) (sessions
   `/v1/sessions`, run center `/v1/runcenter/runs`, workspaces
   `/v1/workspaces`, models `/v1/models`, knowledge `/v1/knowledge`, approvals
   `/v1/approvals`, SSE event stream, recovery `/v1/runs/{id}/recover`).
4. **Integrate against an internal plane** — read the matching contract under
   [`protocols/`](protocols/).
5. **Copy the curl patterns** from [`example/curl-api.md`](example/curl-api.md).

> **Note**: `scripts/` build tooling references the monorepo layout
> (`packages/*`, `apps/desktop/*`) and is kept as a reference of the
> build/release process — it is not runnable in this repo.

## License & Contributions

Licensed under the **Apache License 2.0** (see [`LICENSE`](LICENSE)).
Contributions are accepted under the same license terms (Apache-2.0 §5). PRs,
issues, and integration examples are welcome. For contributions involving core
runtime source or architecture documents, please follow the private main repo
process or contact us first.
