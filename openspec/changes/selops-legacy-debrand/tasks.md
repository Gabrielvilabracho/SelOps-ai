# Tasks: selops-legacy-debrand

## Bucket A — Docs & Repo Meta

- [x] A-01 README.md — h1, alt text, tagline, all Gentleman-Programming/gentle-ai URLs, prose mentions
- [x] A-02 CONTRIBUTING.md — title, prose, GitHub URLs
- [x] A-03 CONTRIBUTORS.md — prose, GitHub URL
- [x] A-04 docs/CODEBASE-GUIDE.md — title, prose
- [x] A-05 docs/opencode-profiles.md — prose
- [x] A-06 docs/intended-usage.md — prose
- [x] A-07 docs/quickstart.md — prose
- [x] A-08 docs/architecture.md — prose
- [x] A-09 docs/kiro.md — prose
- [x] A-10 docs/codebase/*.md — prose (6 files)
- [x] A-11 docs/skill-registry.md — prose
- [x] A-12 docs/gga-powershell-shim.md — prose
- [x] A-13 docs/antigravity-sdd-workaround.md — prose
- [x] A-14 docs/usage.md — CLI command examples
- [x] A-15 docs/non-interactive.md — CLI command examples
- [x] A-16 docs/platforms.md — prose
- [x] A-17 docs/components.md — prose
- [x] A-18 docs/prd-opencode-profiles.md — prose
- [x] A-19 docs/AGENTS.md — prose
- [x] A-20 PRD.md — title, prose (~50 occurrences)
- [x] A-21 PRD-AGENT-BUILDER.md — prose (~20 occurrences)

## Bucket B — Code Copy, Embedded Assets, Install Scripts, npm Meta

- [x] B-01 internal/app/help.go — banner, Remove SelOps line, doc URL
- [x] B-02 internal/agentbuilder/prompt.go — systemPromptBase ecosystem reference
- [x] B-03 internal/assets/opencode/ops-orchestrator.md — title
- [x] B-04 internal/assets/opencode/sdd-orchestrator.md — title
- [x] B-05 internal/assets/opencode/sdd-overlay-multi.json — description
- [x] B-06 internal/assets/opencode/sdd-overlay-single.json — description
- [x] B-07 internal/assets/claude/sdd-orchestrator.md — line ~177 prose
- [x] B-08 internal/assets/claude/commands/sdd-continue.md — prose
- [x] B-09 internal/assets/claude/commands/sdd-ff.md — prose
- [x] B-10 internal/assets/claude/commands/sdd-new.md — prose
- [x] B-11 internal/assets/kiro/sdd-orchestrator.md — line 277 prose
- [x] B-12 internal/assets/windsurf/sdd-orchestrator.md — line 97 prose
- [x] B-13 internal/assets/cursor/sdd-orchestrator.md — line 11 prose
- [x] B-14 package.json — name, description, repository, bugs, homepage URLs
- [x] B-15 scripts/install.sh — banner, GITHUB_OWNER, GITHUB_REPO, footer text
- [x] B-16 scripts/install.ps1 — synopsis, description, GITHUB_OWNER, GITHUB_REPO, banner text
- [x] B-17 internal/app/upgrade_test.go — fixture URL updated

## Bucket C — Repo-Root Skill Files

- [x] C-01 skills/issue-creation/SKILL.md — name, description, author, all prose + GitHub URLs
- [x] C-02 skills/branch-pr/SKILL.md — name, description, author, all prose + GitHub URLs
- [x] C-03 skills/chained-pr/SKILL.md — name, author frontmatter
- [x] C-04 skills/comment-writer/SKILL.md — author frontmatter
- [x] C-05 skills/cognitive-doc-design/SKILL.md — author frontmatter
- [x] C-06 skills/work-unit-commits/SKILL.md — author frontmatter
- [x] C-07 AGENTS.md (repo root) — title, naming convention, skill table names

## Bucket C-path — Path Rename

- [x] CP-01 internal/tui/model.go — ~/.config/gentle-ai/ → ~/.config/selops-ai/
- [x] CP-02 openspec/changes/agent-builder/design.md — path in table
- [x] CP-03 openspec/changes/agent-builder/proposal.md — path in body + rollback
- [x] CP-04 openspec/changes/agent-builder/tasks.md — path in T-07
- [x] CP-05 openspec/changes/agent-builder/spec.md — path in requirement
- [x] CP-06 PRD-AGENT-BUILDER.md — path in 3 locations

## Goldens Updated

- [x] G-01 testdata/golden/sdd-opencode-multi-settings.golden — description field (not tested by suite)
- [x] G-02 testdata/golden/sdd-vscode-instructions.golden — name frontmatter (not tested by suite)
