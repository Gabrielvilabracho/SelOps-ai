# Exploration: complete-sddops

**Change**: `complete-sddops`
**Date**: 2026-06-04
**Artifact store**: openspec + Engram (hybrid)

---

## Current State

`internal/components/sddops/inject.go` (164 lines) is an install-only MVP.
It does exactly two things:

1. **Skills** — calls `skills.InjectWithCapability` for all 11 ops-* skill IDs (6 domain + 5 pipeline) into `adapter.SkillsDir`.
2. **System prompt** — injects the `ops-orchestrator` markdown section into `adapter.SystemPromptFile` via `filemerge.InjectMarkdownSection`.

Three features from the old `sdd` component are **missing**:

| Feature | Old sdd behaviour | sddops current state |
|---|---|---|
| Native sub-agents | Copies `{adapter}/agents/sdd-*.md` files into `adapter.SubAgentsDir()` for every adapter where `SupportsSubAgents()==true` (claude, cursor, kimi, kiro) | NOT implemented |
| Slash commands | Copies `{adapter}/commands/sdd-*.md` into `adapter.CommandsDir()` for adapters where `SupportsSlashCommands()==true` (claude, opencode, kilocode) | NOT implemented |
| OpenCode JSON overlay | Merges a JSON overlay into `opencode.json` registering `ops-orchestrator` as an agent | NOT implemented |

---

## Affected Areas

- `internal/components/sddops/inject.go` — main injector; all three features added here
- `internal/components/sddops/inject_test.go` — new tests for each feature (TDD: RED first)
- `internal/assets/commands.go` — add `OpsCommandsAssetDir(agent)` helper (mirrors `SDDCommandsAssetDir`)
- `internal/assets/claude/agents/` — create `ops-brief.md`, `ops-structure.md`, `ops-produce.md`, `ops-review.md`, `ops-deliver.md`
- `internal/assets/cursor/agents/` — same 5 files (cursor sub-agents use `{{CLAUDE_MODEL}}` too? See note below)
- `internal/assets/kimi/agents/` — same 5 files (kimi also has `.yaml` variants — see note)
- `internal/assets/kiro/agents/` — same 5 files (kiro uses `{{KIRO_MODEL}}` placeholder)
- `internal/assets/claude/commands/` — create `ops-brief.md`, `ops-structure.md`, `ops-produce.md`, `ops-review.md`, `ops-deliver.md` (claude-native frontmatter, no `agent:` field)
- `internal/assets/opencode/commands/` — same 5 files (with `agent: ops-orchestrator` frontmatter)
- `internal/assets/opencode/ops-overlay.json` — NEW: single-mode JSON overlay registering `ops-orchestrator`

**No changes needed in**:
- `internal/model/types.go` — `ComponentSDDOps` already registered
- `internal/catalog/components.go` — `ComponentSDDOps` already in catalog
- `internal/planner/graph.go` — deps already correct
- `internal/cli/validate.go` — already handles `ComponentSDDOps`
- `internal/cli/run.go` — already calls `sddops.Inject`
- `internal/components/skills/presets.go` — skills are injected directly by sddops, not via preset resolution

---

## Feature 1: Native Sub-Agents

### Adapter Capability Matrix (current codebase)

| Adapter | `SupportsSubAgents()` | `EmbeddedSubAgentsDir()` | `SubAgentsDir(home)` |
|---|---|---|---|
| claude | `true` | `"claude/agents"` | `~/.claude/agents/` |
| opencode | `false` | `""` | (unused) |
| kilocode | `false` | `""` | (unused) |
| cursor | `true` | `"cursor/agents"` | `~/.cursor/rules/` (check actual) |
| kimi | `true` | `"kimi/agents"` | depends on adapter |
| kiro | `true` | `"kiro/agents"` | depends on adapter |

**Key finding**: OpenCode and Kilocode both return `SupportsSubAgents()==false`. They receive the `ops-orchestrator` through the **JSON overlay** (feature 3), not as a native sub-agent file. This is the same split as the old `sdd` component.

### OPS Sub-Agent Set — Correct Mapping

The old `sdd` wrote one sub-agent file per SDD phase (sdd-init, sdd-explore, ..., sdd-archive, sdd-onboard — 10 phases + JD agents). For OPS, the correct sub-agent set is the **5 pipeline phase agents** only:

| Old SDD sub-agent | OPS equivalent sub-agent |
|---|---|
| `sdd-init` | (no init phase; ops uses preflight in orchestrator) |
| `sdd-explore` | `ops-brief` (intake + framing) |
| `sdd-propose` | `ops-structure` (execution plan) |
| `sdd-spec` | `ops-produce` (execution) |
| `sdd-design` | `ops-review` (review + close) |
| `sdd-tasks` | `ops-deliver` (delivery + retrospective) |
| `sdd-apply` / `sdd-verify` / `sdd-archive` / `sdd-onboard` | (collapsed into the 5-phase model) |
| `jd-fix-agent` / `jd-judge-a` / `jd-judge-b` | (judgment-day — not part of OPS scope) |

**OPS sub-agents: 5 files** per adapter (ops-brief, ops-structure, ops-produce, ops-review, ops-deliver).

### Asset Files To Create

For each of the 4 adapters that support sub-agents:

**`internal/assets/claude/agents/` — 5 new files**
Format: YAML frontmatter with `name:`, `description:`, `model: {{CLAUDE_MODEL}}`, `tools:`, body referencing the skill at `~/.claude/skills/ops-{phase}/SKILL.md`.

**`internal/assets/cursor/agents/` — 5 new files**
Cursor uses no `model:` placeholder (cursor adapter does NOT implement `claudeModelResolver`). Simple markdown body like the existing sdd-apply.md for cursor.

**`internal/assets/kimi/agents/` — 5 new `.md` files (+ optionally 5 `.yaml` files)**
Kimi has both `.md` and `.yaml` variants for each phase. The `.yaml` is the Kimi-native agent config; the `.md` is the markdown prompt. Must check whether kimi requires YAML for OPS too. **Finding**: the old pattern writes both; create both to stay consistent.

**`internal/assets/kiro/agents/` — 5 new files**
Kiro uses `{{KIRO_MODEL}}` placeholder. Format identical to claude except the model sentinel.

### Model Placeholder Resolution

The old `inject.go` sub-agent loop resolves placeholders inline during copy:
- `{{CLAUDE_MODEL}}` → via `claudeModelResolver` interface (claude adapter implements it)
- `{{KIRO_MODEL}}` → via `kiroModelResolver` interface (kiro adapter implements it)
- Cursor and Kimi: no model placeholder needed

For OPS, the `sddops.Inject` function signature currently has NO model assignment options (no `ClaudeModelAssignments`, no `KiroModelAssignments` in `InjectOptions`). **Decision needed**: resolve OPS sub-agent models at inject time.

**Finding**: Since `sddops.Inject` currently takes `InjectOptions{Capability, WorkspaceDir}` only, and there's no existing TUI picker for OPS model assignments, the simplest approach for the MVP is to **default to a single hardcoded alias** (e.g. `model.ClaudeModelSonnet`) rather than adding full model assignment UI. This matches how the current system prompt injection works — it does not do per-phase model selection. Add `ClaudeModelAlias` and `KiroModelAlias` fields to `InjectOptions` for future expansion, defaulting to `sonnet`.

---

## Feature 2: Slash Commands

### Adapter Capability Matrix

| Adapter | `SupportsSlashCommands()` | `CommandsDir(home)` |
|---|---|---|
| claude | `true` | `~/.claude/commands/` |
| opencode | `true` | `~/.config/opencode/commands/` |
| kilocode | `true` | `~/.config/kilo/commands/` |
| cursor | `false` | (unused) |
| kimi | `false` | (unused) |
| kiro | `false` | (unused) |

### OPS Command Set

The old `sdd` had 9 commands: sdd-init, sdd-explore, sdd-onboard, sdd-apply, sdd-verify, sdd-archive, sdd-continue, sdd-ff, sdd-new.

For OPS, the minimal command set maps to the 5 pipeline phases:

| Old SDD command | OPS command |
|---|---|
| `sdd-init` | (no separate init; preflight is in orchestrator) |
| `sdd-explore` | `ops-brief` |
| `sdd-apply` | `ops-produce` |
| `sdd-verify` | `ops-review` |
| `sdd-archive` | `ops-deliver` |
| `sdd-onboard` | `ops-structure` (or omit for MVP) |
| `sdd-continue`, `sdd-ff`, `sdd-new` | (not applicable to OPS pipeline) |

**Proposed OPS command set: 5 files** (ops-brief, ops-structure, ops-produce, ops-review, ops-deliver).

### Frontmatter Format

**claude/commands/ops-{phase}.md** — Claude-native, NO `agent:` field:
```
---
description: {phase description}
---
Read your skill file at `~/.claude/skills/ops-{phase}/SKILL.md` and follow it exactly.
```

**opencode/commands/ops-{phase}.md** — OpenCode format, WITH `agent: ops-orchestrator` and `subtask: true`:
```
---
description: {phase description}
agent: ops-orchestrator
subtask: true
---
You are the `ops-orchestrator`. Route this to the ops-{phase} sub-agent or execute inline.
```

**Kilocode reuses opencode commands** — the `SDDCommandsAssetDir` pattern switches on `AgentClaudeCode` only, everything else falls back to opencode format. Same for OPS: `OpsCommandsAssetDir(agent)` returns `"claude/commands"` for claude, `"opencode/commands"` for all others (opencode + kilocode).

### New Helper Needed: `OpsCommandsAssetDir`

```go
// internal/assets/commands.go — add:
func OpsCommandsAssetDir(agent model.AgentID) string {
    switch agent {
    case model.AgentClaudeCode:
        return "claude/ops-commands"  // OR "claude/commands" with ops-* prefix filter
    default:
        return "opencode/ops-commands"
    }
}
```

**Alternative**: store OPS commands in the same directory as SDD commands (`claude/commands/`, `opencode/commands/`) — they already have `ops-` prefix so no collision. The injection loop already filters by directory listing.

**Recommendation**: **same directory, ops- prefix** — no new directories needed, no asset directory restructuring. The loop in inject.go reads all files from `CommandsDir` asset dir and copies them. Since OPS injects only OPS commands, we need a dedicated OPS commands subdirectory OR we filter by prefix in the loop.

**Better approach**: dedicated subdirectories `claude/ops-commands/` and `opencode/ops-commands/`. Cleaner separation from SDD commands.

---

## Feature 3: OpenCode JSON Overlay

### What the Old SDD Did

The old `sdd` merged `opencode/sdd-overlay-single.json` (or `sdd-overlay-multi.json`) into `~/.config/opencode/opencode.json`. The single overlay registered:
- `agent.gentle-orchestrator` — mode: primary, prompt from AGENTS.md, task permissions
- `agent.sdd-init` through `agent.sdd-onboard` — 10 sub-agents, mode: subagent, hidden: true

For **OpenCode** (`SupportsSubAgents()==false`), the overlay is the ONLY way to register sub-agents. The orchestrator is also registered here.

### OPS Equivalent

For OPS, we need `opencode/ops-overlay.json`:
- `agent.ops-orchestrator` — mode: primary, prompt from opencode/ops-orchestrator.md (inlined), no task permissions (OPS doesn't use delegate/task tool — it routes inline or via pipeline)
- `agent.ops-brief` through `agent.ops-deliver` — 5 sub-agents, mode: subagent, hidden: true

**Single overlay only** — OPS has no "multi-mode" equivalent. One overlay file suffices.

### Prompt Injection Strategy

The old `sdd` inlined the orchestrator prompt during `inlineOpenCodeSDDPrompts()` — a complex function that handles `{file:}` references, placeholder substitution, model assignments, and profile generation. 

For OPS, **much simpler approach**: inline the `assets.MustRead("opencode/ops-orchestrator.md")` directly into the JSON overlay, the same way the current `sddops.Inject` does it for the system prompt section. No file references, no placeholders, no profiles.

### mergeJSONFile — Already Available

The `mergeJSONFile` function exists in `internal/components/engram/inject.go`, `internal/components/persona/inject.go`, etc. — but it's package-private in each. The pattern across the codebase is to copy the ~15-line function into each package that needs it. `sddops/inject.go` must add its own copy (or extract to `internal/components/filemerge` — but that's scope creep).

### Post-Check

After merging, `sddops.Inject` must verify that `opencode.json` contains `"ops-orchestrator"` agent key. Use the same in-memory-first, disk-fallback pattern as the old sdd.

### Adapters That Get the Overlay

- opencode: `agent.Agent() == model.AgentOpenCode`
- kilocode: `agent.Agent() == model.AgentKilocode` (uses `~/.config/kilo/opencode.json`)

Both already have `SettingsPath()` returning the correct path.

---

## inject.go Changes — Exact New Functions

### 1. `mergeJSONFile` (private, copied from engram pattern)

```go
func mergeJSONFile(path string, overlay []byte) (filemerge.WriteResult, error) {
    base, err := osReadFile(path)
    ...
    merged, err := filemerge.MergeJSONObjects(base, overlay)
    ...
    return filemerge.WriteFileAtomic(path, merged, 0o644)
}
var osReadFile = func(path string) ([]byte, error) { ... }
```

### 2. `injectOpsSubAgents` (new)

```go
func injectOpsSubAgents(targetDir string, adapter agents.Adapter, opts InjectOptions) (InjectionResult, error) {
    if !adapter.SupportsSubAgents() { return InjectionResult{}, nil }
    agentsDir := adapter.SubAgentsDir(targetDir)
    os.MkdirAll(agentsDir, 0o755)
    embeddedDir := adapter.EmbeddedSubAgentsDir()
    // only read ops-*.md entries from embeddedDir
    entries, _ := assets.FS.ReadDir(embeddedDir)
    for _, entry := range entries {
        if !strings.HasPrefix(entry.Name(), "ops-") { continue }
        content := assets.MustRead(embeddedDir + "/" + entry.Name())
        // resolve {{CLAUDE_MODEL}} via claudeModelResolver
        // resolve {{KIRO_MODEL}} via kiroModelResolver
        outPath := filepath.Join(agentsDir, entry.Name())
        filemerge.WriteFileAtomic(outPath, []byte(content), 0o644)
    }
    // post-check ops-brief + ops-deliver exist
}
```

### 3. `injectOpsSlashCommands` (new)

```go
func injectOpsSlashCommands(targetDir string, adapter agents.Adapter) (InjectionResult, error) {
    if !adapter.SupportsSlashCommands() { return InjectionResult{}, nil }
    commandsDir := adapter.CommandsDir(targetDir)
    assetDir := assets.OpsCommandsAssetDir(adapter.Agent())
    entries, _ := fs.ReadDir(assets.FS, assetDir)
    for _, entry := range entries {
        content := assets.MustRead(assetDir + "/" + entry.Name())
        filemerge.WriteFileAtomic(filepath.Join(commandsDir, entry.Name()), ...)
    }
}
```

### 4. `injectOpsOpenCodeOverlay` (new)

```go
func injectOpsOpenCodeOverlay(targetDir string, adapter agents.Adapter) (InjectionResult, error) {
    if adapter.Agent() != model.AgentOpenCode && adapter.Agent() != model.AgentKilocode { return InjectionResult{}, nil }
    settingsPath := adapter.SettingsPath(targetDir)
    overlayContent := assets.MustRead("opencode/ops-overlay.json")
    // inline ops-orchestrator prompt
    overlayBytes := inlineOpsOrchestratorPrompt([]byte(overlayContent), adapter.Agent())
    result, _ := mergeJSONFile(settingsPath, overlayBytes)
    // post-check: settingsPath contains "ops-orchestrator"
}
```

### 5. `Inject` — extend the main function

Add calls to the three new functions after the existing skill+orchestrator injection.

---

## New Asset Files To Create

### Sub-agent assets (per adapter)

Each file follows the same pattern as existing `sdd-apply.md` equivalents but references `ops-{phase}/SKILL.md`:

**claude/agents/** (5 files):
- `ops-brief.md`, `ops-structure.md`, `ops-produce.md`, `ops-review.md`, `ops-deliver.md`
- Format: YAML frontmatter (`name:`, `description:`, `model: {{CLAUDE_MODEL}}`, `tools:`), body instructs executor role + skill path

**cursor/agents/** (5 files):
- Same 5 names, no `model:` placeholder (cursor has no model resolver)

**kimi/agents/** (10 files: 5 `.md` + 5 `.yaml`):
- `.md`: markdown prompt body (no model placeholder)
- `.yaml`: Kimi YAML agent config with `provider`, `model`, `system` fields

**kiro/agents/** (5 files):
- Same 5 names, `model: {{KIRO_MODEL}}` placeholder

### Command assets

**claude/ops-commands/** (5 files):
- `ops-brief.md`, `ops-structure.md`, `ops-produce.md`, `ops-review.md`, `ops-deliver.md`
- Claude-native: `---\ndescription: ...\n---`

**opencode/ops-commands/** (5 files):
- Same 5 names with `agent: ops-orchestrator\nsubtask: true` in frontmatter

### JSON overlay

**opencode/ops-overlay.json** (1 file):
```json
{
  "agent": {
    "ops-orchestrator": {
      "mode": "primary",
      "description": "SelOps OPS Orchestrator",
      "prompt": "__OPS_ORCHESTRATOR_PROMPT__"
    },
    "ops-brief":     { "mode": "subagent", "hidden": true, "description": "...", "prompt": "..." },
    "ops-structure": { "mode": "subagent", "hidden": true, "description": "...", "prompt": "..." },
    "ops-produce":   { "mode": "subagent", "hidden": true, "description": "...", "prompt": "..." },
    "ops-review":    { "mode": "subagent", "hidden": true, "description": "...", "prompt": "..." },
    "ops-deliver":   { "mode": "subagent", "hidden": true, "description": "...", "prompt": "..." }
  }
}
```
(Placeholder replaced at inject time)

---

## Registries — What's Already Done vs. What's Needed

| Registry | Status |
|---|---|
| `model.ComponentSDDOps` | ✅ Already registered |
| `catalog/components.go` | ✅ `ComponentSDDOps` already in catalog |
| `planner/graph.go` deps | ✅ Already correct (`{ComponentSDDOps: {ComponentEngram}}`) |
| `cli/validate.go` | ✅ Already handles `ComponentSDDOps` |
| `cli/run.go` | ✅ Already calls `sddops.Inject` |
| `assets/commands.go` | 🔲 Add `OpsCommandsAssetDir` helper |
| `assets.go` embed directive | ✅ Already embeds `claude`, `opencode`, `cursor`, `kimi`, `kiro` |
| `sdd-overlay-*.json` | ✅ Existing SDD overlays unaffected; OPS gets its own `ops-overlay.json` |

**No new catalog entries, model types, planner deps, or CLI wiring needed.**

---

## TDD Plan Sketch

### Feature 1: Native Sub-Agents

**Failing test (RED)**:
```go
func TestInjectWritesOpsSubAgentFilesForClaude(t *testing.T) {
    home := t.TempDir()
    adapter := claude.NewAdapter()
    Inject(home, adapter, InjectOptions{})
    
    agentsDir := adapter.SubAgentsDir(home)
    for _, phase := range []string{"ops-brief", "ops-structure", "ops-produce", "ops-review", "ops-deliver"} {
        path := filepath.Join(agentsDir, phase+".md")
        info, err := os.Stat(path)
        // FAILS: file doesn't exist yet
        if err != nil { t.Errorf("sub-agent %q not written: %v", phase, err) }
        if info.Size() < 50 { t.Errorf(...) }
    }
}
```

Additional table tests:
- `TestInjectSubAgentsWriteModelPlaceholderResolved` — verify `{{CLAUDE_MODEL}}` is replaced (not literal) in claude sub-agents
- `TestInjectSubAgentsSkippedForOpenCode` — opencode gets NO sub-agent files (SupportsSubAgents==false)
- `TestInjectSubAgentsIdempotent` — second call doesn't change files

### Feature 2: Slash Commands

**Failing test (RED)**:
```go
func TestInjectWritesOpsSlashCommandsForClaude(t *testing.T) {
    home := t.TempDir()
    adapter := claude.NewAdapter()
    Inject(home, adapter, InjectOptions{})
    
    commandsDir := adapter.CommandsDir(home)
    for _, phase := range []string{"ops-brief", "ops-structure", "ops-produce", "ops-review", "ops-deliver"} {
        path := filepath.Join(commandsDir, phase+".md")
        // FAILS: file doesn't exist
        ...
    }
}

func TestInjectWritesOpsSlashCommandsForOpenCode(t *testing.T) {
    // opencode gets ops commands with agent: ops-orchestrator frontmatter
    ...
    content, _ := os.ReadFile(path)
    if !strings.Contains(string(content), "agent: ops-orchestrator") { t.Errorf(...) }
}
```

Additional tests:
- `TestInjectCommandsFrontmatterClaudeHasNoAgentField` — claude commands must NOT have `agent:` field
- `TestInjectCommandsSkippedForCursor` — cursor has no slash commands

### Feature 3: OpenCode JSON Overlay

**Failing test (RED)**:
```go
func TestInjectWritesOpsOrchestratorToOpenCodeJSON(t *testing.T) {
    home := t.TempDir()
    adapter := opencode.NewAdapter()
    Inject(home, adapter, InjectOptions{})
    
    settingsPath := adapter.SettingsPath(home)
    data, _ := os.ReadFile(settingsPath)
    // FAILS: file doesn't exist / key not present
    if !strings.Contains(string(data), `"ops-orchestrator"`) { t.Errorf(...) }
}

func TestInjectOpsOverlayContainsAllPipelineSubAgents(t *testing.T) {
    // verify ops-brief..ops-deliver all present in opencode.json
}

func TestInjectOpsOverlayIsIdempotent(t *testing.T) {
    // two calls → same opencode.json content
}

func TestInjectOpsOverlaySkippedForClaude(t *testing.T) {
    // claude does NOT write to opencode.json
}

func TestInjectOpsOverlayOrchestratorPromptInlined(t *testing.T) {
    // the "ops-orchestrator" agent's "prompt" field contains the orchestrator text, not a placeholder
    // verify it contains "OPS Routing Threshold" or similar known content
}
```

---

## Size / Slicing Forecast

| Item | Est. Lines |
|---|---|
| `sddops/inject.go` — 3 new functions + Inject changes | ~120 |
| `sddops/inject_test.go` — ~12 new tests | ~200 |
| `assets/commands.go` — `OpsCommandsAssetDir` helper | ~12 |
| **Sub-agent assets** (5 × 4 adapters × ~40 lines) = 20 files | ~800 |
| **Kimi YAML** (5 × ~20 lines) | ~100 |
| **Command assets** claude (5 × ~15 lines) + opencode (5 × ~20 lines) | ~175 |
| **ops-overlay.json** | ~80 |
| **Total** | **~1,490 lines** |

**Budget risk**: HIGH. Exceeds 400-line review budget by ~3.7×.

### Recommended PR Slices

| PR | Content | Est. Lines |
|---|---|---|
| PR1 | Asset files only (sub-agents + commands + overlay JSON + `commands.go` helper) — no logic | ~1,170 |
| PR2 | `inject.go` feature 1+2 (sub-agents + slash commands) + tests | ~220 |
| PR3 | `inject.go` feature 3 (JSON overlay + `mergeJSONFile`) + tests | ~220 |

PR1 is still oversized (~1,170 lines) due to the volume of asset files. Consider:
- **Option A**: Split PR1 into asset-creation slices by adapter (claude assets | cursor+kimi+kiro assets | commands | overlay) — 4 PRs but very small each
- **Option B**: Accept PR1 as asset-only (no logic risk) and waive the 400-line budget for pure-content files

**Recommendation**: Option B — asset files are low-risk content additions with no logic. Apply the 400-line budget strictly to logic PRs (PR2, PR3). This gives 3 total PRs:
- PR1: All asset files (content only, no logic, ~1,170 lines) — budget waived
- PR2: inject.go sub-agents + slash commands + tests (~220 lines) ✅
- PR3: inject.go JSON overlay + mergeJSONFile + tests (~220 lines) ✅

---

## Approaches

### Approach 1 (Recommended): Minimal delta — OPS-only directories

- Create `claude/ops-commands/`, `opencode/ops-commands/` as separate asset subdirectories
- Add `OpsCommandsAssetDir()` helper in `assets/commands.go`
- In each adapter's `{adapter}/agents/`, only add `ops-*.md` files (filter by prefix in loop)
- Pros: clean separation from SDD assets; zero risk of breaking existing SDD injection; easy to extend
- Cons: slightly more directories
- Effort: Medium

### Approach 2: Shared command directories (no new subdirs)

- Add `ops-*.md` files directly into `claude/commands/` and `opencode/commands/` alongside `sdd-*.md`
- The injection loop for sddops must filter by `ops-` prefix
- Pros: fewer directories, simpler
- Cons: SDD and OPS commands coexist in same dir → potential coupling/confusion; harder to diff what belongs to whom
- Effort: Low-Medium

**Recommendation: Approach 1** — dedicated `ops-commands/` subdirectories. Clean, future-proof, no coupling to SDD asset paths.

---

## Gotchas and Risks

1. **Kimi YAML agents**: Kimi's `EmbeddedSubAgentsDir()` returns `"kimi/agents"` which contains BOTH `.md` and `.yaml` files. The sub-agent loop must handle both. For OPS, we need `.yaml` variants too — otherwise Kimi will get markdown prompts but won't register the agent in its native format.

2. **Cursor has no model resolver**: Cursor's adapter does not implement `claudeModelResolver` or `kiroModelResolver`. The `ops-*.md` files for cursor must not contain `{{CLAUDE_MODEL}}` or `{{KIRO_MODEL}}` placeholders. Use a fixed text like "claude-sonnet-4-5" or omit the model field entirely.

3. **mergeJSONFile is private per package**: The function must be copied into `sddops/inject.go`. This is the established pattern (engram, persona, mcp, permissions all have their own copy). Do NOT create a shared exported version — that's an architectural change outside this scope.

4. **opencode.json prompt inlining**: The ops-overlay.json must NOT use `{file:}` references for the sub-agent prompts — OpenCode file references require absolute paths which are not known at asset embed time. Inline the skill content directly (or use the skill path pattern with a relative reference like `~/.config/opencode/skills/ops-brief/SKILL.md`). The opencode sdd overlay inlines the orchestrator but uses `{file:}` for phase sub-agents. For OPS MVP, inline everything — simpler and avoids the prompt-file dependency.

5. **`osReadFile` variable**: The `mergeJSONFile` copy needs the `osReadFile` var for testability (the existing packages use this pattern for mocking in tests). Include it.

6. **Sub-agent post-check**: After writing sub-agents, verify at minimum `ops-brief` and `ops-deliver` exist with size ≥ 50 bytes (mirrors the old sdd check for `sdd-apply` and `sdd-verify`).

7. **Kilocode commands directory**: Kilocode uses `~/.config/kilo/commands/` — different path from opencode but same asset format. The `OpsCommandsAssetDir` helper correctly returns `"opencode/ops-commands"` for kilocode since it's not `AgentClaudeCode`.

---

## Recommendation

**Implement approach 1 in 3 PRs** (asset-only PR1, logic+test PR2, logic+test PR3). The implementation is well-understood — the old `sdd` component is the proven reference, and the OPS adaptation is structurally identical with a smaller sub-agent set (5 vs 10 phases). No registry changes needed. No new model types needed. The work is additive.

### Ready for Proposal
Yes — scope is fully bounded. The proposal should specify:
1. The 3-PR split
2. The OPS sub-agent set (5 pipeline phases)
3. The single-overlay-only decision (no OPS "multi-mode")
4. The inline-prompt decision for opencode sub-agents (no {file:} references in ops-overlay.json sub-agents)
5. TDD RED/GREEN/REFACTOR for PRs 2 and 3 (PR1 has no logic to test-first)
