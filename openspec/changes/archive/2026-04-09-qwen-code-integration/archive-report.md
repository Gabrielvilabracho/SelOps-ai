# Archive Report: qwen-code-integration

**Status**: SHIPPED  
**Archived**: 2026-06-02  
**Mode**: Standard (strict_tdd: true in config)  
**Verify verdict**: PASS  

---

## What was built

Full integration of Qwen Code (Alibaba) as a supported AI coding agent in the Gentle AI ecosystem, mirroring the Gemini CLI adapter pattern.

### Delivered

| Area | Details |
|------|---------|
| Adapter | `internal/agents/qwen/` — full `agents.Adapter` implementation (21 methods) |
| Detection | Binary lookup + `~/.qwen/` config dir stat |
| Installation | `npm install -g @qwen-code/qwen-code@latest` with sudo logic per platform |
| Config paths | `~/.qwen/` (GlobalConfigDir, SystemPromptDir), `QWEN.md`, `settings.json` |
| Strategies | `StrategyFileReplace` (system prompt) + `StrategyMergeIntoSettings` (MCP) |
| Capabilities | `SupportsSlashCommands: true`, `SupportsSkills: true`, `SupportsMCP: true` |
| SDD orchestrator | `internal/assets/qwen/sdd-orchestrator.md` (Qwen-specific paths) |
| Permissions | `auto_edit` overlay |
| Engram setup | Slug `"qwen-code"` |
| CLI/TUI wiring | Cases in `validate.go`, `model.go`, `config_scan.go` |

---

## PR merged

| PR | Description |
|----|-------------|
| #263 (upstream gentle-ai) | `feat: add Qwen Code agent integration` |

Merged to upstream `gentleman-programming/gentle-ai` main. Also reflected in this fork via `git log ef601b7` (engram fix) and `c72f7e8` (e2e hardening).

---

## Test evidence

```
go build ./...              → clean
go vet ./...                → clean
go test ./...               → 17 packages pass, 0 fail
internal/agents/qwen        → ok (coverage 82.9%)
```

**Spec compliance**: 40/40 scenarios — full PASS.  
**Design coherence**: all 10 decisions followed, 0 deviations.

---

## Known limitations / follow-ups (deferred)

| # | Item |
|---|------|
| S-01 | Integration test for `gentle-ai install --agent qwen-code --dry-run` in test harness |
| S-02 | Coverage from 82.9% → 90%+ (test `defaultStat()` error path, `OutputStyleDir()`) |
| W-01 | `tasks.md` checkboxes left unchecked (tasks created retroactively — documentation debt only) |
