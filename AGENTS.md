---
doc_id: sdk-go-agent-entry
tier: T1
status: active
freshness_score: 100
last_verified: 2026-08-29
code_refs: []
spec_refs: []
compaction_level: L1
compacted_from: null
tags: [agent-guidance, entry, sdk-go]
---

# sdk-go Agent Entry

This repository is the HotelByte official Go SDK (`github.com/hotelbyte-com/sdk-go`).
It is a separate Git root. When this tree is opened alone, use this file; do not assume a hotel-be parent instruction chain is loaded.

## Scope

- Owned: public client, auth, hotels, room mapping, transport, protocol helpers, and `examples/`.
- Out of scope here: hotel-be monorepo workflow, UAT deploy, portal RBAC, schema migrations, and Lookout ops.

## Working rules

- Prefer small reversible diffs; keep exported APIs backward compatible unless the PR explicitly documents a breaking change.
- Follow existing package layout (`client.go`, `auth.go`, `hotels.go`, `transport.go`, `protocol/`).
- Do not commit live credentials, production tokens, or customer data in tests or examples.

## Verification

Smallest proving check for touched packages:

```bash
go test ./...
```

When only one package changed, prefer a focused run (`go test -run TestName ./`).

## Code Review Rules

- Flag breaking exported signatures without an explicit compatibility note.
- Flag examples or tests that embed live credentials or hit production endpoints by default.
- Flag generated/source drift when protocol helpers are regenerated from an upstream schema.
