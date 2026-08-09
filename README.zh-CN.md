# PRAXOVELA

**[English](./README.md)**

> 一个可治理（governable）、本地优先（local-first）的桌面级 AI 执行 Agent 运行时。

**PRAXOVELA**（原名 **AxisRobo Agent**，运行时：**AXON Core**）将**治理**视为一等公民原语：Agent 请求能力、策略面决策、高风险操作进入沙箱执行、每一步操作都追加进可审计、可重放的事件链——且这一切都发生在本机，数据默认不离开你的机器。

本仓库是 PRAXOVELA 生态的**公开开源门面**，以 **Apache License 2.0** 发布，面向所有希望与 PRAXOVELA 集成的第三方开发者、集成商与构建方。这里承载：面向客户端的 API 文档、稳定的集成契约（Protocols）、示例与构建工具。

---

## 为什么选择 PRAXOVELA

主流 Agent 框架把工具交给 Agent，却给了太多权力、太少问责。PRAXOVELA 解决的是 Agent 开始承担真实任务时会暴露的问题：

- **不可问责的工具调用** —— 每一次工具调用都经过默认拒绝（deny-by-default）的能力门禁，写入持久化 trace，崩溃或重启后可以重放。
- **不可恢复的 Run** —— 效果账本（effect ledger）+ 检查点存储（checkpoint store）让被打断的 Run 精确续跑；已提交的效果绝不重复执行。
- **敏感上下文泄漏到云端** —— 上下文默认本地驻留；远程执行是显式、受治理的选择。
- **运行不可信模型输出** —— skill 与 adapter 在 WASM / 容器 / OS 沙箱中执行，工件使用前强制 digest 校验。

## 本仓库包含什么

| 路径 | 内容 |
|------|------|
| [`docs/api.md`](docs/api.md) | `axond` 本地代理 API 参考（客户端契约） |
| [`protocols/axislink/`](protocols/axislink/README.md) | AxisLink UI 协议（L6 UI ↔ L5 Agent） |
| [`protocols/capability/`](protocols/capability/README.md) | 能力契约（L5 Agent ↔ L4 治理） |
| [`protocols/connector/`](protocols/connector/README.md) | 连接器协议（L5/L3 ↔ 企业系统） |
| [`protocols/sandbox/`](protocols/sandbox/README.md) | 沙箱提供者接口（L3 ↔ L2） |
| [`example/`](example/README.md) | 集成示例（针对 axond API 的 curl 用例） |
| `scripts/` | 构建/发布工具（引用 monorepo 布局） |
| `version.json` | 产品版本单一来源 |

## 本仓库不是什么

**不包含任何核心运行时源码（AGPL）、任何企业版特性（专有），也不包含架构/路线图/规划文档。** 核心代码与架构文档的权威来源是**私有主仓库**——永不复制到本仓库。

## 生态（三仓库）

| 仓库 | License | 可见性 | 内容 |
|------|---------|--------|------|
| `praxovela` | AGPL | 私有 | 主运行时源码：AXON core、Janus Gateway、Vulcan Forge、Run Center、AxisLink、adapters |
| `praxovela-ee` | 专有 | 私有 | 企业版：Archon Guard、Argus Trace、runcenter-ext、PDP/SIEM/OIDC/Mneme |
| `praxovela-open` | Apache-2.0 | **公开** | SDK/API 文档、集成契约、示例、构建工具 |

## 集成者快速开始

1. **启动运行时** —— 安装 PRAXOVELA 桌面应用，或构建/运行 `axond`（AXON Core）。默认监听 `http://localhost:8420`（可用 `AXON_PORT` 覆盖）。
2. **确认在线**：
   ```bash
   curl http://localhost:8420/health
   # → {"status":"ok","runtime":"AXON Core","version":"0.2.0-beta.1"}
   ```
3. **阅读客户端契约** —— [`docs/api.md`](docs/api.md)：会话 `/v1/sessions`、运行中心 `/v1/runcenter/runs`、工作区 `/v1/workspaces`、模型 `/v1/models`、知识库 `/v1/knowledge`、审批 `/v1/approvals`、SSE 事件流、恢复 `/v1/runs/{id}/recover`。
4. **对接内部平面** —— 阅读 [`protocols/`](protocols/) 下对应契约。
5. **复制调用示例** —— [`example/curl-api.md`](example/curl-api.md)。

> **注意**：`scripts/` 下的构建脚本引用的是 monorepo 布局（`packages/*`、`apps/desktop/*`），是构建/发布流程的参考实现，在本仓库中**不可直接运行**。

## License 与贡献

本仓库以 **Apache License 2.0** 发布（见 [`LICENSE`](LICENSE)）。所有贡献均按 Apache-2.0 同一许可条款提交（Apache-2.0 §5）。欢迎提交 issues、PR 与集成示例；若你的贡献涉及核心运行时或架构文档，请遵循私有主仓库流程，或先与我们联系。
