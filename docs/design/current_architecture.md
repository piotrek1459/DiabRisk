# Current System Architecture - DiabRisk

**Scope:** current implementation in the repository

## Purpose

This document describes the architecture that is actually implemented in
the current repository.

## Architectural Style

The system is a small web application deployed locally on Kubernetes
through k3d. At runtime it consists of:

- a Svelte frontend served by Nginx
- an API gateway written in Go
- an authentication service written in Go
- an ML inference service written in Python with FastAPI
- a PostgreSQL database
- a migration runner service named `data-svc`

The public entry point is Traefik ingress on `http://localhost`.

## Runtime Topology

```text
Browser
  |
  v
Traefik Ingress (:80, host: localhost)
  /         -> frontend (:80)
  /auth/*   -> api-gateway (:8080) -> auth-svc (:8081) -> PostgreSQL (:5432)
  /api/*    -> api-gateway (:8080) -> ml-api (:8000)
```

## Service Inventory

| Component | Implementation | Port | Responsibility |
|----------|----------------|------|----------------|
| frontend | Svelte app built with Vite, served by Nginx | 80 | login/register UI, form, result display |
| api-gateway | Go + Gin | 8080 | public backend entrypoint, auth middleware, proxy to auth-svc and ml-api |
| auth-svc | Go + Gin | 8081 | registration, login, logout, session validation, admin seeding |
| ml-api | Python + FastAPI | 8000 | feature validation and prediction |
| data-svc | Go binary | no public HTTP API | waits for DB, runs migrations, verifies schema |
| postgres | PostgreSQL 16 | 5432 | users, sessions, seeded schema tables |

`api-gateway`, `auth-svc`, and `ml-api` each implement `/healthz`, but the
current ingress does not expose `/healthz` publicly.

## Public and Internal Boundaries

### Publicly Reachable Through Ingress

| Path | Backing Component |
|------|-------------------|
| `/` | frontend |
| `/auth/register` | api-gateway -> auth-svc |
| `/auth/login` | api-gateway -> auth-svc |
| `/auth/logout` | api-gateway -> auth-svc |
| `/auth/session` | api-gateway -> auth-svc |
| `/api/risk` | api-gateway -> ml-api |
| `/api/features` | api-gateway -> ml-api |

### Internal-Only Services

- `auth-svc`
- `ml-api`
- `postgres`
- `data-svc`

## Request Flow

### Registration and Login

1. The frontend sends `POST /auth/register` or `POST /auth/login`.
2. Traefik forwards the request to `api-gateway`.
3. `api-gateway` proxies the request to `auth-svc`.
4. `auth-svc` writes user and session state to PostgreSQL.
5. `auth-svc` returns a `session_token` HttpOnly cookie.

### Session Validation

1. The frontend calls `GET /auth/session`.
2. `api-gateway` proxies the request to `auth-svc`.
3. `auth-svc` reads the `session_token` cookie.
4. The token is hashed and checked against `auth_sessions`.
5. On success, the current user object is returned as JSON.

### Risk Prediction

1. The frontend submits the health form to `POST /api/risk`.
2. `api-gateway` checks the session by calling `GET /auth/session`.
3. The gateway forwards the feature payload to `ml-api`.
4. `ml-api` validates features, builds the input vector, loads the joblib artifact, and returns a prediction payload.
5. The frontend renders `RiskPercent`, `Category`, and `Message`.

## Data and Persistence

| Table or Asset | Current Usage |
|---------------|---------------|
| `users` | actively used by `auth-svc` |
| `auth_sessions` | actively used by `auth-svc` |
| `model_versions` | seeded by migrations, not read by the browser flow |
| `assessments` | schema present, not written by `/api/risk` |
| `reports` | schema present, not used by the current runtime flow |
| `audit_logs` | schema present, not used by the current runtime flow |
| `models/diabrisk_screening.joblib` | loaded by `ml-api` |

## Deployment Shape

The supported local deployment path is Kubernetes on k3d.

- images are built from `Dockerfile.api-gateway`, `Dockerfile.auth-svc`,
  `Dockerfile.data-svc`, `Dockerfile.frontend`, and `Dockerfile.ml-api`
- manifests live under `deploy/k8s/`
- setup is automated by `scripts/install-local-k3d.ps1` and
  `scripts/install-local-k3d.sh`

## Current Architectural Constraints

- the current browser auth flow is email/password only
- prediction results are not persisted after `/api/risk`
- `data-svc` is an infrastructure workload, not a CRUD API
- the runtime path does not include a separate report service
