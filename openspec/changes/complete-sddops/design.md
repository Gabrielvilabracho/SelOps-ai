# Design: Complete sddops

## Technical Approach

Extend `internal/components/sddops/inject.go` with three new capability-gated injection helpers (`injectOpsSubAgents`, `injectOpsSlashCommands`, `injectOpsOpenCodeOverlay`) plus a package-private `mergeJSONFile`. Add `assets.OpsCommandsAssetDir(agent)` and the missing asset files (sub-agents, ops-commands, ops-overlay.json). All gating is driven by `adapter.SupportsSubAgents()` / `SupportsSlashCommands()` / `Agent()==OpenCode|Kilocode`. The current 164-line `Inject()` keeps its skills + system-prompt section as the first two steps; the three new calls slot in afterward and contribute to the aggregate `InjectionResult`.

## Inject() Call Flow

```
Inject(targetDir, adapter, opts)
  │
  ├─ 1. skills.InjectWithCapability       (existing)        → result
  ├─ 2. post-check ops skill SKILL.md >=100              (existing)
  ├─ 3. ops-orchestrator system-prompt section            (existing)
  │
  ├─ 4. injectOpsSubAgents(homeDir, adapter, opts)        (NEW)
  │       gate: adapter.SupportsSubAgents()
  │
  ├─ 5. injectOpsSlashCommands(homeDir, adapter)          (NEW)
  │       gate: adapter.SupportsSlashCommands()
  │
  └─ 6. injectOpsOpenCodeOverlay(homeDir, adapter)        (NEW)
          gate: adapter.Agent() == OpenCode || Kilocode

  Each helper returns (changed bool, files []string, error)
  Aggregated into result.Changed (OR) and result.Files (append).
```

`homeDir` is resolved once at the top of `Inject()` via the same helper the old `sdd` package used (read from the adapter or `os.UserHomeDir()`). The existing `targetDir` (workspace) is preserved for the system-prompt step; sub-agents and commands write under `homeDir` (the user-level adapter dir).

## Function Signatures

```go
// injectOpsSubAgents copies ops-*.md (and ops-*.yaml for Kimi) from
// adapter.EmbeddedSubAgentsDir() into adapter.SubAgentsDir(homeDir),
// resolving {{CLAUDE_MODEL}} / {{KIRO_MODEL}} placeholders for adapters
// that implement claudeModelResolver / kiroModelResolver.
// Post-check: at least one of {ops-brief, ops-deliver} present as .md
// or .yaml with Size() >= 10 bytes.
func injectOpsSubAgents(homeDir string, adapter agents.Adapter, opts InjectOptions) (changed bool, files []string, err error)

// injectOpsSlashCommands copies every file from assets.OpsCommandsAssetDir(adapter.Agent())
// into adapter.CommandsDir(homeDir). No size post-check (matches pre-strip).
func injectOpsSlashCommands(homeDir string, adapter agents.Adapter) (changed bool, files []string, err error)

// injectOpsOpenCodeOverlay reads internal/assets/opencode/ops-overlay.json,
// inlines the ops-orchestrator prompt from opencode/ops-orchestrator.md,
// then merges into adapter.SettingsPath(homeDir) via mergeJSONFile.
// Post-check (semantic, against merged in-memory bytes): ops-orchestrator
// AND ops-brief keys present under "agent". No byte threshold.
func injectOpsOpenCodeOverlay(homeDir string, adapter agents.Adapter) (changed bool, files []string, err error)

// mergeJSONFile is package-private to sddops (copied, not shared).
// Mirrors the engram pattern: read base (nil if missing), deep-merge
// via filemerge.MergeJSONObjects, write atomically. Returns the merged
// bytes so the post-check validates the in-memory result (Windows/WSL2
// rename-visibility safety).
func mergeJSONFile(path string, overlay []byte) (writeResult filemerge.WriteResult, merged []byte, err error)
```

## Placeholder Resolution

Reuse the proven type-assertion pattern from `git show f599511^:internal/components/sdd/inject.go` (lines 574–600). Inside `injectOpsSubAgents`, after reading each entry:

```go
if kmr, ok := adapter.(kiroModelResolver); ok {
    alias := resolveKiroAlias(opts, phase) // sonnet default
    content = strings.ReplaceAll(content, "{{KIRO_MODEL}}", kmr.KiroModelID(alias))
}
if cmr, ok := adapter.(claudeModelResolver); ok {
    alias := resolveClaudeModelAlias(opts.ClaudeModelAssignments, phase)
    content = strings.ReplaceAll(content, "{{CLAUDE_MODEL}}", cmr.ClaudeModelID(alias))
}
```

Cursor implements neither interface → its asset files contain no placeholder. Kimi implements neither → `.md` and `.yaml` ship literal content. `claudeModelResolver` / `kiroModelResolver` interfaces already exist in the `sdd` package; sddops MUST declare its own local copies (private) to avoid an inter-package dependency.

`InjectOptions` gains optional `ClaudeModelAssignments map[string]model.ClaudeModelAlias` and `KiroModelAssignments map[string]model.ClaudeModelAlias` (zero-value safe). MVP callers pass nil → safe-default alias `sonnet`.

## Asset File Inventory

All paths relative to `internal/assets/`.

| Path | Count | Purpose |
|------|-------|---------|
| `claude/agents/ops-{brief,structure,produce,review,deliver}.md` | 5 | Claude sub-agents, contain `{{CLAUDE_MODEL}}` |
| `cursor/agents/ops-{brief,structure,produce,review,deliver}.md` | 5 | Cursor sub-agents, no placeholder |
| `kimi/agents/ops-{brief,structure,produce,review,deliver}.md` | 5 | Kimi MD form |
| `kimi/agents/ops-{brief,structure,produce,review,deliver}.yaml` | 5 | Kimi YAML form |
| `kiro/agents/ops-{brief,structure,produce,review,deliver}.md` | 5 | Kiro sub-agents, contain `{{KIRO_MODEL}}` |
| `claude/ops-commands/ops-{brief,structure,produce,review,deliver}.md` | 5 | Claude slash commands (NEW dir) |
| `opencode/ops-commands/ops-{brief,structure,produce,review,deliver}.md` | 5 | OpenCode/Kilocode slash commands (NEW dir) |
| `opencode/ops-overlay.json` | 1 | Single JSON overlay |

### Skeleton: Claude sub-agent (`claude/agents/ops-brief.md`)

```markdown
---
name: ops-brief
description: >
  OPS phase 1 — capture the operational brief. Use when ops-orchestrator launches
  the brief phase. Reads inputs, persists brief artifact.
model: {{CLAUDE_MODEL}}
tools: Read, Edit, Write, Glob, Grep, Bash, mcp__plugin_engram_engram__mem_search, mcp__plugin_engram_engram__mem_get_observation, mcp__plugin_engram_engram__mem_save
---

You are the OPS **brief** executor. Do this phase's work yourself. Do NOT delegate further.
You are not the orchestrator. Do NOT call the Task tool. Do NOT launch sub-agents.

Read the skill file at `~/.claude/skills/ops-brief/SKILL.md` and follow it exactly.
```

### Skeleton: Kimi YAML (`kimi/agents/ops-brief.yaml`)

```yaml
name: ops-brief
description: OPS phase 1 — capture the operational brief
prompt_file: ops-brief.md
```

### Skeleton: Claude command (`claude/ops-commands/ops-brief.md`)

```markdown
---
description: Run the OPS brief phase for a change
---

Follow the OPS orchestrator workflow for running the brief phase on "$ARGUMENTS".
Launch `ops-brief` to capture the operational brief.
```

### Skeleton: OpenCode command (`opencode/ops-commands/ops-brief.md`)

```markdown
---
description: Run the OPS brief phase for a change
agent: ops-orchestrator
subtask: true
---

Run the OPS brief phase for the change named "$ARGUMENTS".
```

## ops-overlay.json Shape

```json
{
  "agent": {
    "ops-orchestrator": {
      "mode": "primary",
      "description": "SelOps OPS Orchestrator — coordinates the 5-phase OPS pipeline",
      "prompt": "<INLINED_AT_INJECT_TIME from opencode/ops-orchestrator.md>",
      "permission": {
        "task": {
          "__replace__": {
            "*": "deny",
            "ops-brief": "allow",
            "ops-structure": "allow",
            "ops-produce": "allow",
            "ops-review": "allow",
            "ops-deliver": "allow"
          }
        }
      },
      "tools": { "read": true, "write": true, "edit": true, "bash": true,
                 "delegate": true, "delegation_read": true, "delegation_list": true }
    },
    "ops-brief":     { "mode": "subagent", "hidden": true,
                       "description": "Capture the operational brief",
                       "prompt": "You are an OPS executor for the brief phase... Read ~/.config/opencode/skills/ops-brief/SKILL.md and follow it.",
                       "tools": { "read": true, "write": true, "edit": true, "bash": true } },
    "ops-structure": { "mode": "subagent", "hidden": true, "...": "same shape" },
    "ops-produce":   { "mode": "subagent", "hidden": true, "...": "same shape" },
    "ops-review":    { "mode": "subagent", "hidden": true, "...": "same shape" },
    "ops-deliver":   { "mode": "subagent", "hidden": true, "...": "same shape" }
  }
}
```

**Prompt inlining rule**: `ops-orchestrator.prompt` is read at inject time from `internal/assets/opencode/ops-orchestrator.md` (already exists) and string-replaced into the overlay BEFORE `mergeJSONFile` runs. The 5 sub-agent prompts are short literals inside the JSON (no file references). This satisfies "no absolute paths" — the merged `opencode.json` contains only inline strings.

Source asset uses sentinel `"prompt": "{{OPS_ORCHESTRATOR_PROMPT}}"`; `injectOpsOpenCodeOverlay` does `bytes.Replace(overlay, []byte("{{OPS_ORCHESTRATOR_PROMPT}}"), jsonEscape(orchMD), 1)` (json-escape the raw markdown so the result is valid JSON).

## OpsCommandsAssetDir Helper

Added to `internal/assets/commands.go`, mirroring `SDDCommandsAssetDir`:

```go
// OpsCommandsAssetDir returns the embedded OPS slash-command asset directory
// for an agent. Claude uses Claude-native frontmatter under claude/ops-commands;
// agents without a dedicated command set fall back to OpenCode-compatible assets.
func OpsCommandsAssetDir(agent model.AgentID) string {
    switch agent {
    case model.AgentClaudeCode:
        return "claude/ops-commands"
    default:
        return "opencode/ops-commands"
    }
}
```

## mergeJSONFile Semantics

- **Read base**: `os.ReadFile(path)`; treat `os.IsNotExist` as empty (`nil`) — bootstrapping a missing `opencode.json` is supported.
- **Deep merge**: delegate to `internal/components/filemerge.MergeJSONObjects(base, overlay)` (the same helper used by engram/mcp/persona/permissions/theme). It performs non-destructive recursive object merging — overlay keys win on leaf collision, existing siblings preserved.
- **Atomic write**: `filemerge.WriteFileAtomic(path, merged, 0o644)` → temp file + rename.
- **Return merged bytes** in addition to `WriteResult` so the post-check validates the in-memory result. The post-check unmarshals `merged` into `map[string]any`, asserts `root["agent"]["ops-orchestrator"]` and `root["agent"]["ops-brief"]` both exist. No re-read from disk on the happy path; disk re-read is the fallback if in-memory bytes are empty (mirrors `f599511^:inject.go:646–656`).
- **Private copy**: sddops gets its own `mergeJSONFile` body. No extraction. Same shape as engram's (see `internal/components/engram/inject.go:506`).

## Architecture Decisions

### Decision: Private `mergeJSONFile` per package
**Choice**: Copy the function into `internal/components/sddops/`.
**Alternatives**: Extract to `filemerge.MergeJSONFile(path, overlay)`.
**Rationale**: The codebase already has SIX private copies (engram, mcp, persona, permissions, theme, old sdd). Extracting now expands scope, churns five other packages, and breaks the proposal's "Out of Scope" line. Honor the established convention; revisit in a dedicated refactor.

### Decision: Dedicated `ops-commands/` asset directories (Approach A)
**Choice**: New `claude/ops-commands/` and `opencode/ops-commands/`.
**Alternatives**: Co-locate ops-*.md inside existing `claude/commands/` and `opencode/commands/`.
**Rationale**: SDD and OPS are separable components — a future "OPS only, no SDD" preset must not pull SDD command files. Co-location forces filename prefix filtering at copy time and couples the two surfaces. Dedicated dirs + dedicated `OpsCommandsAssetDir` helper is parallel to the existing `SDDCommandsAssetDir` pattern.

### Decision: Semantic JSON post-check, no byte threshold
**Choice**: Unmarshal merged bytes, assert two keys present.
**Alternatives**: Stat `opencode.json` and require N bytes.
**Rationale**: Replicates pre-strip behavior exactly (`f599511^:inject.go:658,667`). Byte thresholds give false confidence — a 50-KB `opencode.json` with the wrong agent key still fails at runtime. Key-presence is the only check that proves the overlay landed.

### Decision: 3-PR slice mapped to functions
- **PR1 (assets only)**: every file in the inventory + `OpsCommandsAssetDir` helper. `size:exception` documented. Adds embedded files only — no behavior change.
- **PR2 (inject features 1+2)**: `injectOpsSubAgents` + `injectOpsSlashCommands` + tests. Depends on PR1. ~220 lines.
- **PR3 (inject feature 3)**: `injectOpsOpenCodeOverlay` + private `mergeJSONFile` + tests. Depends on PR1. ~220 lines.
**Rationale**: Each PR builds and tests independently. PR1 alone is a no-op asset drop (existing `Inject()` unchanged). PR2 and PR3 can theoretically merge in either order because they touch disjoint code paths in `Inject()` (sub-agents vs settings file).

## Data Flow

```
adapter.SupportsSubAgents()  ──┐
adapter.EmbeddedSubAgentsDir() ┴─→ assets.FS.ReadDir → for each ops-*
                                  ├─ resolve placeholders ─→ filemerge.WriteFileAtomic
                                  └─ → adapter.SubAgentsDir(homeDir)/

adapter.SupportsSlashCommands()──→ assets.OpsCommandsAssetDir(agent)
                                  → assets.FS.ReadDir → for each ops-*.md
                                  → filemerge.WriteFileAtomic → CommandsDir/

adapter.Agent()∈{OpenCode,Kilocode}
  ├─ assets.MustRead("opencode/ops-overlay.json")
  ├─ assets.MustRead("opencode/ops-orchestrator.md")
  ├─ inline orchestrator prompt into overlay bytes
  ├─ mergeJSONFile(settingsPath, overlay) ──→ (writeResult, merged)
  └─ semantic post-check on merged bytes
```

## Test Surface (spec scenario → function → table-test)

| # | Scenario | Function | Test idea |
|---|----------|----------|-----------|
| 1 | Sub-agent capability gates native assets | `injectOpsSubAgents` | table: adapters {claude:true, opencode:false} × assert files presence |
| 2 | Slash-command capability gates command assets | `injectOpsSlashCommands` | table: adapters {claude:true, kimi:false} × assert files |
| 3 | Enabled feature writes are atomic and post-checked | all three | inject into `t.TempDir()`, stat files, assert sizes per pre-strip rules |
| 4 | Second identical inject reports Changed=false | all three | run `Inject` twice, assert `result.Changed==false` second time |
| 5 | Claude receives the five OPS sub-agents | `injectOpsSubAgents` | stat 5 files in claude agents dir, each `>=10` |
| 6 | Claude/Kiro placeholders resolved before write | `injectOpsSubAgents` | grep `{{CLAUDE_MODEL}}` / `{{KIRO_MODEL}}` absent from output |
| 7 | Kimi dual `.md` + `.yaml` | `injectOpsSubAgents` | stat both extensions for every phase |
| 8 | OpenCode/Kilocode no native sub-agent files | `injectOpsSubAgents` | assert sub-agent dir empty for those adapters |
| 9 | Native sub-agent injection idempotent | `injectOpsSubAgents` | run twice, mtime/content diff |
| 10 | Claude five OPS slash commands | `injectOpsSlashCommands` | stat 5 files in `claude commands dir` |
| 11 | Claude command frontmatter has description, no `agent:` | `injectOpsSlashCommands` | regex on file content |
| 12 | OpenCode/Kilocode command has `agent: ops-orchestrator` + `subtask: true` | `injectOpsSlashCommands` | regex on file content |
| 13 | Cursor/Kimi/Kiro no OPS slash commands | `injectOpsSlashCommands` | assert commands dir empty |
| 14 | Slash command injection idempotent | `injectOpsSlashCommands` | run twice |
| 15 | OpenCode overlay registers orchestrator + 5 pipeline agents | `injectOpsOpenCodeOverlay` | unmarshal merged, assert 6 keys |
| 16 | Merge preserves existing user keys | `injectOpsOpenCodeOverlay` | seed `opencode.json` with `{"foo":"bar"}`, assert preserved |
| 17 | Prompts inlined, no absolute path refs | `injectOpsOpenCodeOverlay` | grep merged JSON for `{file:` and `/Users/`/`$HOME` — assert none |
| 18 | Non-OpenCode adapters skip overlay | `injectOpsOpenCodeOverlay` | run on claude, assert settings file untouched |
| 19 | Overlay injection idempotent | `injectOpsOpenCodeOverlay` | run twice, byte-equal merged result |

All tests use `t.TempDir()`, `filemerge.WriteFileAtomic` produces real files. Tests live in `internal/components/sddops/inject_test.go` (already exists for current MVP) — extend, do not replace.

## Migration / Rollout

No migration required. Additive and capability-gated. Re-applying remains idempotent (existing behavior preserved + new feature outputs are content-equal on second run, so atomic writes report `Changed=false`).

## Open Questions

None blocking. All locked decisions came from the proposal, spec, and the post-check thresholds Engram entry (#2650).
