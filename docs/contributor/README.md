# Contributor Documentation

This section is intended for developers maintaining or extending DiabRisk.
It documents the current codebase, runtime shape, and development
workflow.

## Recommended Reading Order

1. `setup.md`
2. `architecture.md`
3. `codebase-guide.md`
4. `api-reference.md`
5. `testing.md`

## Document Map

| File | Purpose |
|------|---------|
| `setup.md` | local environment, k3d workflow, rebuild loop, config variables |
| `architecture.md` | current runtime topology, request flows, state ownership |
| `codebase-guide.md` | source-of-truth map for the repo and common change scenarios |
| `api-reference.md` | current backend contracts and route ownership |
| `testing.md` | current automated test scope and missing coverage |

## Scope Note

Use the files in this folder as the operational source of truth for the
running system.

`docs/design/` complements this folder with current design-oriented
descriptions of the same implementation: system vision, domain vocabulary,
use cases, interaction flows, and conceptual models.
