# selops-operational-mcp-wiring Specification

## Purpose

Define the operational MCP wiring boundary so the fork writes external MCP connection entries without becoming a runtime, retrieval engine, or secret carrier.

## Requirements

### Requirement: External MCP wiring boundary

The operational MCP component MUST only write MCP connection entries into agent configuration, MUST NOT contain RAG or tool runtime logic, and MUST preserve the existing DEV profile unchanged.

#### Scenario: Wiring writes connection entries only
Given external operational MCP servers exist outside this repository
When the operational MCP component is applied
Then it MUST write only connection entries such as command or URL plus env references into agent MCP config
And it MUST NOT add retrieval, storage, server, or execution logic to this fork

### Requirement: Existing merge path and Context7 safety

The operational MCP component MUST reuse `filemerge.MergeJSONObjects` for config merging and MUST NOT modify the existing Context7 MCP component.

#### Scenario: Wiring reuses the established JSON merge mechanism
Given agent MCP config already merges component JSON into a shared file
When operational MCP entries are added
Then the merge behavior MUST reuse `filemerge.MergeJSONObjects`
And the change MUST NOT require a new merge algorithm or write path

#### Scenario: Context7 remains untouched
Given the existing Context7 MCP component is already supported
When operational MCP wiring is introduced
Then the Context7 component MUST remain behaviorally unchanged
And operational wiring MUST be additive rather than a modification of Context7 assets

### Requirement: Graceful degradation and secret handling

The operational MCP component MUST degrade gracefully when external servers are not configured, and connection secrets MUST follow existing security conventions.

#### Scenario: Missing external MCP configuration does not fail install
Given an external operational MCP server has not yet been configured
When the operational MCP component is installed
Then installation MUST degrade gracefully by writing a documented placeholder or disabled entry
And the install MUST NOT fail solely because the external server is absent

#### Scenario: Secrets are never embedded in assets
Given MCP connections may require environment variables or credentials
When operational MCP assets are defined
Then secrets MUST be referenced through existing env and permission conventions
And no secret value MAY be embedded directly in fork assets
