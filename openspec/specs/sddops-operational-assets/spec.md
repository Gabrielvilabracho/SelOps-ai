# SDDOPS Operational Assets Specification

## Purpose

Defines install-time behavior for OPS operational assets so supported adapters receive the five-phase OPS pipeline through native sub-agents, slash commands, or the OpenCode JSON overlay.

## Requirements

### Requirement: Capability-Driven Operational Injection

`sddops.Inject` MUST derive OPS operational asset behavior from adapter capability methods and target-specific adapter paths. It MUST use atomic writes for every written artifact and MUST verify each enabled feature produced its critical output before reporting success, replicating the pre-strip `sdd` post-check pattern (git ref `f599511^`):
- Sub-agents: post-check MUST confirm at least one critical phase file exists as `.md` or `.yaml` with `Size() >= 10` bytes (detects empty/truncated writes; the `>= 10` byte floor and `.md`-or-`.yaml` acceptance match the pre-strip loop).
- OpenCode overlay: post-check MUST be semantic — confirm the `ops-orchestrator` key and at least one pipeline sub-agent key (`ops-brief`) are present in the merged JSON. It MUST NOT use a byte-size threshold.
- Slash commands: no size post-check is required (matches the pre-strip behavior, which wrote command files without a size guard).

Repeated runs with identical inputs MUST be idempotent: the second run SHALL report `Changed=false` for the feature outputs it re-evaluates and SHALL NOT duplicate or rewrite equivalent content.

#### Scenario: Sub-agent capability gates native assets
- GIVEN one adapter reports `SupportsSubAgents()` true and another reports false
- WHEN `sddops.Inject` runs for both adapters
- THEN only the capable adapter receives OPS sub-agent files in `SubAgentsDir(homeDir)`

#### Scenario: Slash-command capability gates command assets
- GIVEN one adapter reports `SupportsSlashCommands()` true and another reports false
- WHEN `sddops.Inject` runs for both adapters
- THEN only the capable adapter receives OPS command files in `CommandsDir(homeDir)`

#### Scenario: Enabled feature writes are atomic and post-checked
- GIVEN a feature is enabled for an adapter
- WHEN `sddops.Inject` writes its OPS artifacts
- THEN each artifact write uses the atomic file-write path already used by `filemerge.WriteFileAtomic`
- AND the feature verifies its critical output per the pre-strip post-check pattern: sub-agents require a critical phase file (`.md` or `.yaml`) of `Size() >= 10` bytes; the overlay requires the `ops-orchestrator` and `ops-brief` keys present in the merged JSON; slash commands require no size post-check

#### Scenario: A second identical inject reports no changes
- GIVEN an adapter target has already received all OPS artifacts for its enabled features
- WHEN `sddops.Inject` runs again with the same inputs
- THEN the second run reports `Changed=false`
- AND no duplicate artifact content is produced

---

### Requirement: Native OPS Sub-Agent Injection

For adapters where `SupportsSubAgents()` is true, the system MUST install exactly five OPS phase sub-agents: `ops-brief`, `ops-structure`, `ops-produce`, `ops-review`, and `ops-deliver`. Files MUST be copied from `adapter.EmbeddedSubAgentsDir()` into `adapter.SubAgentsDir(homeDir)`. Adapters where `SupportsSubAgents()` is false MUST NOT receive native OPS sub-agent files.

#### Scenario: Claude receives the five OPS sub-agents
- GIVEN the Claude adapter reports `SupportsSubAgents()` true
- WHEN `sddops.Inject` runs
- THEN `ops-brief.md` through `ops-deliver.md` exist in Claude's sub-agent directory
- AND the post-check confirms a critical phase file (`.md` or `.yaml`) of `Size() >= 10` bytes

#### Scenario: Claude and Kiro placeholders are resolved before write
- GIVEN Claude or Kiro OPS sub-agent assets contain model placeholders
- WHEN `sddops.Inject` writes those assets
- THEN the written files do not contain literal `{{CLAUDE_MODEL}}` or `{{KIRO_MODEL}}`

#### Scenario: Kimi receives dual-format sub-agent assets
- GIVEN the Kimi adapter reports `SupportsSubAgents()` true
- WHEN `sddops.Inject` runs
- THEN each OPS phase is installed in both `.md` and `.yaml` form in Kimi's sub-agent directory

#### Scenario: OpenCode and Kilocode do not receive native sub-agent files
- GIVEN the OpenCode or Kilocode adapter reports `SupportsSubAgents()` false
- WHEN `sddops.Inject` runs
- THEN no OPS sub-agent files are created in a native sub-agent directory

#### Scenario: Native sub-agent injection is idempotent
- GIVEN a supported sub-agent adapter has already been injected once
- WHEN `sddops.Inject` runs a second time with the same inputs
- THEN the second run reports no sub-agent changes
- AND the written sub-agent content is unchanged

---

### Requirement: OPS Slash Command Injection

For adapters where `SupportsSlashCommands()` is true, the system MUST install exactly five OPS slash commands: `ops-brief`, `ops-structure`, `ops-produce`, `ops-review`, and `ops-deliver`. Command assets MUST come from `OpsCommandsAssetDir(agent)` and be written to `adapter.CommandsDir(homeDir)`. Adapters where `SupportsSlashCommands()` is false MUST NOT receive OPS command files.

#### Scenario: Claude receives the five OPS slash commands
- GIVEN the Claude adapter reports `SupportsSlashCommands()` true
- WHEN `sddops.Inject` runs
- THEN `ops-brief.md` through `ops-deliver.md` exist in Claude's commands directory

#### Scenario: Claude command frontmatter stays Claude-native
- GIVEN Claude OPS command files are written
- WHEN a written command file is read
- THEN its frontmatter contains `description`
- AND its frontmatter does not contain `agent: ops-orchestrator`

#### Scenario: OpenCode and Kilocode command frontmatter targets the orchestrator
- GIVEN the OpenCode or Kilocode adapter reports `SupportsSlashCommands()` true
- WHEN `sddops.Inject` writes OPS command files
- THEN each written command contains `description`, `agent: ops-orchestrator`, and `subtask: true`

#### Scenario: Cursor, Kimi, and Kiro do not receive OPS slash commands
- GIVEN an adapter reports `SupportsSlashCommands()` false
- WHEN `sddops.Inject` runs
- THEN no OPS command files are created for that adapter

#### Scenario: OPS slash command injection is idempotent
- GIVEN a supported slash-command adapter has already been injected once
- WHEN `sddops.Inject` runs a second time with the same inputs
- THEN the second run reports no command changes
- AND the written command content is unchanged

---

### Requirement: OpenCode OPS Overlay Injection

For OpenCode and Kilocode, the system MUST merge `internal/assets/opencode/ops-overlay.json` into `opencode.json` as a non-destructive overlay. The merged result MUST register `agent.ops-orchestrator` in `primary` mode and the five OPS phase agents in `subagent` mode with `hidden: true`. Overlay prompts MUST be inlined in the merged JSON and MUST NOT depend on absolute file paths. Adapters other than OpenCode and Kilocode MUST NOT receive this overlay.

#### Scenario: OpenCode overlay registers orchestrator and pipeline agents
- GIVEN the OpenCode adapter target has no prior OPS overlay entries
- WHEN `sddops.Inject` runs
- THEN `opencode.json` contains `ops-orchestrator`
- AND it contains `ops-brief`, `ops-structure`, `ops-produce`, `ops-review`, and `ops-deliver`

#### Scenario: Overlay merge preserves existing user keys
- GIVEN `opencode.json` already contains unrelated user-defined keys
- WHEN `sddops.Inject` merges the OPS overlay
- THEN the unrelated pre-existing keys remain present and unchanged

#### Scenario: Overlay prompts are inlined without absolute path references
- GIVEN the OPS overlay is merged for OpenCode or Kilocode
- WHEN the resulting `opencode.json` is read
- THEN the OPS agent prompt fields contain inline prompt content
- AND the prompt fields do not contain absolute path references

#### Scenario: Non-OpenCode adapters do not receive the JSON overlay
- GIVEN an adapter is not OpenCode or Kilocode
- WHEN `sddops.Inject` runs
- THEN no OPS overlay merge is attempted for that adapter's settings file

#### Scenario: OpenCode overlay injection is idempotent
- GIVEN OpenCode or Kilocode has already been injected once
- WHEN `sddops.Inject` runs a second time with the same inputs
- THEN the second run reports no overlay changes
- AND the resulting `opencode.json` content is unchanged
