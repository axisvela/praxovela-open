# PRAXOVELA-open — Repository Rules

This is the **public open-source repo** (Apache-2.0) of the PRAXOVELA ecosystem. Keep it focused on public-facing assets only.

- **Source of truth**: The private AGPL main repo (`praxovela` / AxisRobo Agent) owns all core runtime source and architecture docs. Enterprise features live in the private proprietary repo (`praxovela-ee`).
- **Never copy** core source code, architecture docs, roadmap/planning docs, or enterprise features into this repo.
- **Allowed content**: API docs, integration contracts (protocols/), examples, build/release tooling, and client-facing metadata (version.json).
- When updating files, keep them in sync with the private repo where applicable (e.g. `docs/api.md`, `version.json`, `protocols/`).
- Release tags must stay aligned with the core repo: `PRAXOVELA-open` and `PRAXOVELA` share the exact same clean `vX.Y.Z` tag. `PRAXOVELA-EE` may tag independently, but all release tags still use clean `major.minor.patch` semantic versions unless the user explicitly overrides this rule.
- This repo is bilingual: `README.md` is English-primary with a link to the
  Chinese parallel file `README.zh-CN.md`; keep both in sync. Other docs match
  the ecosystem style (English headers + Chinese detail where applicable).

Verification commands for this repo (no Go/Rust build here):

```sh
git status            # clean working tree
git log -1 --stat     # review last commit
```
