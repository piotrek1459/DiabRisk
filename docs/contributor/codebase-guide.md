# Codebase Guide

## Purpose

This document helps a developer quickly locate the right place to make
changes and understand which files must stay in sync.

## Repository Map

| Area | Primary Paths | Change Here When |
|------|---------------|------------------|
| Frontend UI | `frontend/src/App.svelte`, `frontend/src/lib/` | changing login/register views, risk form, or result rendering |
| API gateway | `services/api-gateway/main.go` | adding browser-facing routes, auth middleware, or proxy behavior |
| Authentication | `services/auth-svc/main.go` | changing registration, login, logout, session validation, or admin seeding |
| Database schema | `services/data-svc/migrations/` | adding tables, columns, seed data, or constraints |
| Migration runner | `services/data-svc/main.go` | changing startup checks, migration execution, or schema verification |
| ML inference | `src/FastAPI/ml_api.py` | changing feature validation, model loading, or prediction output |
| ML training | `src/Ml/main.py` | retraining models or changing training-time preprocessing |
| Deployment | `deploy/k8s/`, `scripts/install-local-k3d.*`, `Dockerfile.*` | changing local infrastructure, image build rules, manifests, or secrets |
| Documentation | `docs/` | changing setup, runtime assumptions, API contracts, or user guidance |

## Source of Truth by Concern

| Concern | Primary Source |
|--------|----------------|
| public browser routes | `services/api-gateway/main.go` and `deploy/k8s/ingress.yaml` |
| session cookie lifecycle | `services/auth-svc/main.go` |
| expected ML feature names | model artifact loaded by `src/FastAPI/ml_api.py` |
| frontend request payload shape | `frontend/src/App.svelte` |
| database schema | `services/data-svc/migrations/*.sql` |
| local deployment behavior | `scripts/install-local-k3d.ps1` and `scripts/install-local-k3d.sh` |
| runtime image build inputs | `Dockerfile.*` files |

## Common Change Scenarios

### Adding or Renaming a Risk Feature

Keep these places aligned:

1. `frontend/src/App.svelte`
2. `src/Ml/main.py`
3. `src/FastAPI/ml_api.py`
4. `docs/contributor/api-reference.md`
5. `data/raw/raw_dataset_description.md` and `data/processed/dataset_description.md` if the documented contract changes

### Changing Authentication or Session Handling

Review all of:

1. `services/auth-svc/main.go`
2. `services/api-gateway/main.go`
3. `docs/contributor/api-reference.md`
4. `docs/contributor/architecture.md`
5. `services/auth-svc/main_test.go`

### Changing the Prediction Contract

Review all of:

1. `src/Ml/main.py`
2. `models/`
3. `src/FastAPI/ml_api.py`
4. `frontend/src/App.svelte`
5. `deploy/k8s/ml-api.yaml`
6. `docs/contributor/api-reference.md`

### Adding a Database Table or Column

Use a new numbered migration pair in `services/data-svc/migrations/`.
Then verify:

1. readers or writers in `auth-svc`
2. seed data if bootstrap rows are required
3. `services/data-svc/main.go` schema verification
4. documentation if the runtime behavior changes

## Current Boundaries and Caveats

- `data-svc` currently runs migrations and schema verification; it is not a CRUD API
- `assessments`, `reports`, and `audit_logs` exist in the schema but are not part of the active browser flow
- the browser auth path is local email/password plus a session cookie
- `frontend/src/App.svelte.old` is a legacy file and should not be treated as the source of truth
- `deploy/k8s/configmap.yaml` is not used by the active install scripts

## Documentation Update Checklist

When you change the system, update docs if the change affects:

- route ownership or request/response payloads
- setup steps or required tools
- runtime topology or service responsibilities
- test commands or CI coverage
- user-visible behavior in the browser
