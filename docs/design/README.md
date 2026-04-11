# Design Documentation

This folder contains design-oriented documentation for the current DiabRisk
implementation.

## Recommended Reading Order

1. `system_vision.md`
2. `current_architecture.md`
3. `technology_stack.md`
4. `system_dictionary.md`
5. `use_case_model.md`
6. `interaction_diagrams.md`
7. `class_model.md`
8. `initial_system_specification.md`

## Document Map

| File | Purpose |
|------|---------|
| `system_vision.md` | product-level description of the current application |
| `current_architecture.md` | runtime topology, boundaries, and request flow |
| `technology_stack.md` | technologies actively used in the repository |
| `system_dictionary.md` | current domain and technical vocabulary |
| `use_case_model.md` | current actors and use cases |
| `interaction_diagrams.md` | sequence-style descriptions of key runtime flows |
| `class_model.md` | conceptual model of current runtime data and objects |
| `initial_system_specification.md` | compact current-state system specification |

## Source of Truth

The Markdown files in this folder describe the current implementation.
When documentation and code diverge, treat the code in `frontend/`,
`services/`, `src/FastAPI/`, and `src/Ml/` as the final source of truth.
