# selops-operational-assets Specification

## Purpose

Define the fork-private asset contract for SelOps operational personas, skills, and SDD-OPS domain assets while preserving upstream mergeability and DEV profile safety.

## Requirements

### Requirement: Fork-private naming and embedding

All new operational assets MUST use `ops-*` or `selops-*` prefixes, MUST be embedded through the existing `go:embed` asset mechanism, and MUST remain distinct from upstream persona identifiers.

#### Scenario: New operational assets are namespaced
Given the fork already contains upstream-oriented asset identifiers
When a new operational persona, skill, or MCP asset is added
Then its identifier and file naming MUST use an `ops-*` or `selops-*` prefix
And the new asset MUST NOT collide with an upstream gentle-ai identifier

#### Scenario: Operational assets use the existing embedding channel
Given the installer already distributes static assets through embedded files
When operational assets are shipped
Then they MUST be embedded through the existing `go:embed` mechanism
And the change MUST NOT introduce a second distribution channel

#### Scenario: Operator persona is fork-private
Given the system already supports gentleman, neutral, and custom personas
When the operational persona asset is installed
Then it MUST be a fork-private operator persona distinct from those personas
And the existing DEV personas MUST remain unchanged

### Requirement: SDD-OPS domain asset shape

The SDD-OPS component MUST ship one namespaced operational skill asset per domain and MUST define only the asset shape and presence for this change.

#### Scenario: Standard documentation domain asset ships
Given the SDD-OPS component is present
When its assets are enumerated
Then it MUST ship a namespaced `ops-*` skill asset for standard documentation

#### Scenario: Modular architecture domain asset ships
Given the SDD-OPS component is present
When its assets are enumerated
Then it MUST ship a namespaced `ops-*` skill asset for modular architecture

#### Scenario: Data contracts domain asset ships
Given the SDD-OPS component is present
When its assets are enumerated
Then it MUST ship a namespaced `ops-*` skill asset for data contracts

#### Scenario: Governance domain asset ships
Given the SDD-OPS component is present
When its assets are enumerated
Then it MUST ship a namespaced `ops-*` skill asset for governance

#### Scenario: Observability domain asset ships
Given the SDD-OPS component is present
When its assets are enumerated
Then it MUST ship a namespaced `ops-*` skill asset for observability

#### Scenario: Graduated autonomy domain asset ships
Given the SDD-OPS component is present
When its assets are enumerated
Then it MUST ship a namespaced `ops-*` skill asset for graduated autonomy
