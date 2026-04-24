# Architecture

## System Overview



```text
Traefik Ingress (:80, host: localhost)
  /         -> frontend (Nginx :80)
  /auth/*   -> api-gateway (Go/Gin :8080) -> auth-svc (Go/Gin :8081) -> PostgreSQL (:5432)
  /api/*    -> api-gateway (Go/Gin :8080) -> ml-api (FastAPI :8000)
                                           -> data-svc (Go :8082) -> PostgreSQL (:5432)
```

## Runtime Topology

| Service | Source | Runtime Port | Exposure | Responsibility |
|---------|--------|--------------|----------|----------------|
| frontend | `frontend/` | 80 | ingress `/` | login/register UI, risk form, result rendering |
| api-gateway | `services/api-gateway/` | 8080 | ingress `/auth/*`, `/api/*` | browser-facing backend entrypoint, auth middleware, proxy logic |
| auth-svc | `services/auth-svc/` | 8081 | internal service | registration, login, logout, session validation |
| ml-api | `src/FastAPI/` | 8000 | internal service | model loading, feature validation, prediction endpoints |
| data-svc | `services/data-svc/` | 8082 | internal service | runs migrations, verifies schema, stores assessment history |
| postgres | `deploy/k8s/postgres.yaml` | 5432 | internal service | users, sessions, seeded schema tables |

Service-local health endpoints exist on `api-gateway`, `auth-svc`, and
`ml-api`, but the current ingress manifest does not publish `/healthz`
publicly.

## Request Flows

### Authentication and Session Flow

1. The frontend calls `POST /auth/register` or `POST /auth/login`.
2. Traefik forwards `/auth/*` to `api-gateway`.
3. `api-gateway` proxies the request to `auth-svc` without rewriting the path.
4. `auth-svc` creates or validates the user in PostgreSQL.
5. `auth-svc` creates an `auth_sessions` record and returns a `session_token` HttpOnly cookie.
6. The frontend restores the session later with `GET /auth/session`.

### Risk Prediction Flow

1. The frontend submits the form to `POST /api/risk`.
2. `api-gateway` reads the `session_token` cookie.
3. The gateway calls `GET /auth/session` on `auth-svc`.
4. If the session is valid, the gateway forwards the features to `ml-api`.
5. `ml-api` validates the feature set, loads the joblib artifact, and returns a prediction payload.
6. `api-gateway` sends the assessment candidate to `data-svc`.
7. `data-svc` writes the assessment to PostgreSQL and returns the saved record.
8. The frontend renders the returned `RiskPercent`, `Category`, and `Message`.

### Assessment History Flow

1. The frontend calls `GET /api/history`.
2. `api-gateway` validates the session through `auth-svc`.
3. `api-gateway` calls `GET /internal/users/:userId/assessments` on `data-svc`.
4. `data-svc` loads the user's assessments from PostgreSQL.
5. The frontend renders the returned list.

### Startup Flow

1. `scripts/install-local-k3d.ps1` or `scripts/install-local-k3d.sh` recreates the k3d cluster.
2. The script creates `postgres-secret` and `auth-secret`.
3. Docker images are built and imported into the cluster.
4. Manifests from `deploy/k8s/` are applied.
5. `data-svc` waits for PostgreSQL, runs migrations, verifies the schema, and then starts its internal HTTP API.

## State and Persistence

| Asset | Owner | Current Usage |
|------|-------|---------------|
| `session_token` cookie | `auth-svc` | browser session identifier |
| `users` table | `auth-svc` | registered users and seeded admin account |
| `auth_sessions` table | `auth-svc` | hashed session tokens and expiry metadata |
| `model_versions` table | migrations | seeded metadata row, not read by the main browser flow |
| `assessments` table | `data-svc` | written on each successful `/api/risk` call and read by `/api/history` |
| `reports` table | schema only | present in the DB, not used by the current runtime path |
| `audit_logs` table | schema only | present in the DB, not written by the current runtime path |
| `models/diabrisk_screening.joblib` | ML training pipeline | artifact loaded by `ml-api` |

## Code-Level Boundaries

- `api-gateway` owns all browser-facing backend routes.
- `auth-svc` owns credentials, sessions, and cookie lifecycle.
- `ml-api` owns the feature contract and prediction response shape.
- `data-svc` owns migration execution, assessment persistence, and history reads.
- the frontend owns form defaults, request serialization, result rendering, and history display.

## Important Current Constraints

- the active browser auth flow is email/password plus a session cookie
- the runtime stack does not include a separate report service
- the current ingress exposes `/`, `/auth/*`, and `/api/*`, but not `/healthz`
- `deploy/k8s/configmap.yaml` exists in the repo but is not part of the active install path
