# selops-operational-profile Specification

## Purpose

Define a parallel SelOps operational preset that installs operator-facing assets without changing the existing DEV profile.

## Requirements

### Requirement: Operational preset selection

The system MUST allow `PresetSelOpsOperational` to install the operational profile as an additive preset that bundles `PersonaOperator`, operational skills, the SDD-OPS component, and operational MCP wiring.

#### Scenario: Selecting the operational preset installs the operational bundle
Given the installer supports the existing DEV preset
When a user selects `PresetSelOpsOperational`
Then the plan MUST include operator persona, operational skills, SDD-OPS, and operational MCP wiring
And the plan MUST treat them as additive operational components rather than DEV replacements

### Requirement: DEV profile safety contract

The existing DEV profile MUST remain BYTE-FOR-BYTE unchanged when the operational preset is installed, planned, or registered.

#### Scenario: Installing ops does not mutate DEV assets
Given DEV persona, skills, SDD, and MCP assets already exist
When the operational preset is installed
Then the DEV assets MUST remain BYTE-FOR-BYTE unchanged
And no DEV component identifier, asset path, or preset membership MAY be rewritten

#### Scenario: Adding ops does not change engine switch behavior
Given existing engine switch cases handle current DEV behavior
When operational preset support is added
Then existing switch case behavior MUST remain unchanged for pre-existing cases
And any operational behavior MUST be introduced only through additive registrations or new cases

### Requirement: Profile coexistence and rollback

The system MUST allow DEV and operational profiles to coexist in separate namespaces, and rollback to DEV-only MUST remain clean.

#### Scenario: DEV and ops profiles coexist without cross-mutation
Given both DEV and operational profiles are installed in the same fork
When either profile is listed, planned, or reapplied
Then each profile MUST retain separate component namespaces and separate skill directories
And installing one profile MUST NOT remove or alter the other

#### Scenario: DEV-only state is restored by not selecting or uninstalling ops
Given a fork can run in DEV-only mode
When the operational preset is not selected or its components are uninstalled
Then the resulting state MUST restore clean DEV-only behavior
And no residual operational requirement MAY remain on the DEV profile
