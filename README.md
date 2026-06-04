<div align="center">

<h1>SelOps-ai</h1>

<p><strong>AI agent configuration installer for enterprise teams.</strong></p>

<p>
<img src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white" alt="Go 1.22+">
<img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT">
<img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey" alt="Platform">
<img src="https://img.shields.io/badge/TDD-Strict-brightgreen" alt="Strict TDD">
</p>

</div>

---

## What It Does

**SelOps-ai** is a Go CLI tool that configures AI agents (Claude, OpenCode, Kimi) for enterprise teams. It resolves a component graph and installs role-specific profiles, skills, MCPs, and personas in a single command.

```bash
selops-ai install --agent opencode --preset selops-operational
```

---

## Profiles

Two fully isolated profiles ship out of the box:

| Preset | Role | What it installs |
|--------|------|-----------------|
| `full-gentleman` | Developer | SDD workflow, coding skills, DEV persona |
| `selops-operational` | AI Operator | OPS skills, MCP wiring, Operator persona |

DEV/OPS separation is enforced by byte-for-byte regression tests — an OPS install never mutates DEV configuration.

---

## Operational Profile — 6 Domains

The `selops-operational` preset covers six domains of enterprise AI operations:

| # | Domain | Description |
|---|--------|-------------|
| 1 | Standard Documentation | READMEs, ADRs, API contracts — written before implementation |
| 2 | Modular Architecture | Module boundaries, published interfaces, data store ownership |
| 3 | Data Contracts | Versioned schemas between producers and consumers |
| 4 | Governance | Approval workflows, compliance checkpoints, immutable audit trails |
| 5 | Observability | Metrics, distributed tracing, and logging standards |
| 6 | Graduated Autonomy | Permission levels for AI agents (Suggest → Supervised → Autonomous) |

---

## Quick Start

```bash
# Install the CLI
go install github.com/Gabrielvilabracho/SelOps-ai/cmd/selops-ai@latest

# Install the operational profile for OpenCode
selops-ai install --agent opencode --preset selops-operational

# Install the developer profile
selops-ai install --agent opencode --preset full-gentleman

# Dry run — see what would be installed
selops-ai install --agent opencode --preset selops-operational --dry-run
```

---

## Supported Agents

| Agent | System Prompt | MCP | Skills | Commands |
|-------|-------------|-----|--------|----------|
| Claude (claude.ai) | ✅ | ✅ | ✅ | ✅ |
| OpenCode | ✅ | ✅ | ✅ | ✅ |
| Kimi Code | ✅ | ✅ | ✅ | ✅ |
| Gemini CLI | ✅ | ✅ | ✅ | ✅ |
| Cursor | ✅ | ✅ | ✅ | ✅ |
| Windsurf | ✅ | ✅ | ✅ | — |
| VS Code (Copilot) | ✅ | ✅ | — | — |
| Qwen Code | ✅ | ✅ | ✅ | ✅ |

---

## Architecture

SelOps-ai is built around a **component graph** — a directed acyclic graph where each node is an installable unit with explicit dependencies and ordering.

```
Preset → Component resolution → Ordered install plan → Injectors → Filesystem
```

Each component has:
- An **injector** — pure function that writes to the agent's config directory
- A **strategy** — how to write (replace, merge, append)
- A **post-check** — validates the written output before returning

---

## Development

```bash
# Run tests
go test ./...

# Run full CI check
go test ./... && go vet ./... && gofmt -l .

# Build binary
go build -o selops-ai ./cmd/selops-ai
```

**50+ tested packages** · Strict TDD · Golden tests per adapter and strategy · Chained PR releases

---

## Author

**Gabriel Vila Bracho**  
[github.com/Gabrielvilabracho](https://github.com/Gabrielvilabracho)

---

## License

MIT
