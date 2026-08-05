# PRAXOVELA — Governed Agent Execution Environment

**PRAXOVELA** (formerly **AxisRobo Agent**) is a local-first, governable desktop agent runtime for AI execution. PRAXOVELA's core principle is **deny by default**: the agent requests capability, a governance plane decides, high-risk operations run sandboxed, and everything is audited.

This repository is the **public open-source face** of the PRAXOVELA ecosystem — the place third parties integrate.

---

## 这是什么仓库

**PRAXOVELA-open** 是 PRAXOVELA（AxisRobo Agent）生态的**公开开源门面**，以 Apache-2.0 协议发布，面向所有希望与 PRAXOVELA 集成的第三方开发者、集成商与构建方。这里承载：

- **SDK/API 文档** —— `docs/api.md`：axond 本地代理 API 参考（会话、运行中心、工作区、模型、知识库、审批、审计等客户端契约）。
- **集成契约（Protocols）** —— `protocols/`：四份稳定契约文档：
  - `protocols/axislink/` —— AxisLink UI 协议（L6 UI ↔ L5 Agent）
  - `protocols/capability/` —— 能力契约（L5 Agent ↔ L4 治理）
  - `protocols/connector/` —— 连接器协议（L5/L3 ↔ 企业系统）
  - `protocols/sandbox/` —— 沙箱提供者接口（L3 ↔ L2）
- **示例** —— `example/`：面向集成者的调用示例（含 axond API 的 curl 用例）。
- **构建工具** —— `scripts/`：构建/发布/测试脚本与 SBOM、release manifest 生成器。

## 这仓库不是什么

**PRAXOVELA-open 不包含任何核心运行时源码（AGPL），也不包含任何企业版特性（专有）。**

- ❌ 没有 `packages/` —— AXON Core、Janus Gateway、Vulcan Forge 等核心 Go 源码，**绝不**在本仓库。
- ❌ 没有架构文档 / 路线图 / 实施计划 —— `architecture.md`、`ROADMAP.md`、`ecosystem-boundaries.md`、规划文档等留在私有主仓库。
- ❌ 没有桌面应用（Tauri shell）与插件 SDK 源码。

> 架构类内容属私有仓库（见下表 `praxovela`），仅供授权维护者访问。核心源码与架构文档的权威来源是私有主仓库，**永不复制到本仓库**。

## 三仓库生态

| Repo | License | Visibility | Contents |
|------|---------|------------|----------|
| `praxovela` (AxisRobo Agent) | AGPL | Private (may go public later) | Main runtime source: AXON core, Janus, Vulcan, Run Center, AxisLink, adapters |
| `praxovela-ee` | Proprietary | Private | Enterprise: Archon Guard, Argus Trace, runcenter-ext, PDP/SIEM/OIDC/Mneme |
| `praxovela-open` | Apache-2.0 | Public | SDK/API docs, integration contracts, examples, build tooling |

## 集成者快速开始

1. **启动运行时** —— 安装并启动 AxisRobo Agent 桌面应用，或直接构建/运行 `axond`（AXON Core），默认监听 `http://localhost:8420`（`AXON_PORT` 可配置）。
2. **健康检查** —— 确认运行时在线：

   ```bash
   curl http://localhost:8420/health
   # → {"status":"ok","runtime":"AXON Core","version":"0.3.0"}
   ```

3. **读取 API 文档** —— 完整端点契约见 [`docs/api.md`](docs/api.md)（会话 /v1/sessions、运行中心 /v1/runcenter/runs、工作区 /v1/workspaces、模型 /v1/models、知识库 /v1/knowledge、审批 /v1/approvals、SSE 事件流等）。
4. **按协议集成** —— 需要与 AXON 内部平面对接时，阅读 `protocols/` 下对应契约。
5. **调用示例** —— 见 [`example/curl-api.md`](example/curl-api.md)。

> **注意**：`scripts/` 下的构建脚本参考的是 monorepo 布局（`packages/*`、`apps/desktop/*`），它们是构建/发布流程的示例参考，在公开仓库中**不可直接运行**。

## License & Contribution

- 本仓库以 **Apache License 2.0** 发布，全文见 [`LICENSE`](LICENSE)。
- **贡献即同意**：所有贡献均按 Apache-2.0 同一许可条款提交（Apache-2.0 §5）。
- 欢迎提交 issues、PR 与集成示例；若你的贡献涉及核心运行时或架构文档，请转向私有主仓库流程，或先与我们联系。

---

# PRAXOVELA — Governed Agent Execution Environment

**PRAXOVELA** is a local-first, governable desktop agent runtime for AI execution. Deny by default: agents request capability, governance decides, high-risk operations are sandboxed, and everything is audited.

This repository is the **public, Apache-2.0 open-source face** of the PRAXOVELA ecosystem. It carries only public-facing assets:

| Path | Contents |
|------|----------|
| `docs/api.md` | axond local agent API reference (client contract) |
| `protocols/axislink/` | AxisLink UI protocol (L6 UI ↔ L5 Agent) |
| `protocols/capability/` | Capability contract (L5 Agent ↔ L4 Governance) |
| `protocols/connector/` | Connector protocol (L5/L3 ↔ enterprise systems) |
| `protocols/sandbox/` | Sandbox provider interface (L3 ↔ L2) |
| `example/` | Integration examples (curl against the axond API) |
| `scripts/` | Build/release tooling (references the monorepo layout) |
| `version.json` | Product version single source |

## What this repo is NOT

No core runtime source (AGPL), no enterprise features (proprietary), no architecture/roadmap/planning docs. The source of truth for core code and architecture is the **private main repository** — it is never copied here.

## Ecosystem

| Repo | License | Visibility | Contents |
|------|---------|------------|----------|
| `praxovela` (AxisRobo Agent) | AGPL | Private (may go public later) | Main runtime source: AXON core, Janus, Vulcan, Run Center, AxisLink, adapters |
| `praxovela-ee` | Proprietary | Private | Enterprise: Archon Guard, Argus Trace, runcenter-ext, PDP/SIEM/OIDC/Mneme |
| `praxovela-open` | Apache-2.0 | Public | SDK/API docs, integration contracts, examples, build tooling |

## Quick start for integrators

1. Run the runtime: install the AxisRobo Agent desktop app, or build/run `axond` (AXON Core). Default endpoint `http://localhost:8420` (override with `AXON_PORT`).
2. Verify it is alive: `curl http://localhost:8420/health`
3. Read the client-facing contract in [`docs/api.md`](docs/api.md).
4. For internal-plane integration, read the matching contract under `protocols/`.
5. Copy the curl patterns from [`example/curl-api.md`](example/curl-api.md).

> **Note**: `scripts/` build tooling references the monorepo layout (`packages/*`, `apps/desktop/*`) and is kept as a reference of the build/release process — it is not runnable in this repo.

## License & contributions

Licensed under the **Apache License 2.0** (see [`LICENSE`](LICENSE)). Contributions are accepted under the same license terms (Apache-2.0 §5). PRs, issues, and integration examples are welcome.
