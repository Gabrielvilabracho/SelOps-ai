---
name: ops-modular-architecture
description: "SelOps modular architecture patterns and boundaries. Trigger: When designing or reviewing module boundaries, service contracts, or data store ownership."
---

# Modular Architecture

## Purpose

Guide the design of loosely coupled, independently deployable modules within the SelOps platform.

## Guidelines

- Define explicit module boundaries with published interfaces.
- Avoid direct cross-module database access; use service contracts instead.
- Each module owns its own data store and migration lifecycle.

<!-- Content research is DEFERRED — this is a layout/placeholder stub. -->
