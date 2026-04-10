# Current System Specification - DiabRisk

**Scope:** current implementation in the repository

## 1. Purpose

DiabRisk is an educational diabetes-risk screening application. The current
system collects 21 health indicators, authenticates the user with a local
account, and returns a prediction produced by the ML service.

## 2. External Actors

| Actor | Current Role |
|------|--------------|
| Visitor | opens the app, registers, logs in |
| Authenticated User | submits the risk form and reads the result |
| Developer / Operator | runs the local cluster, inspects services, and updates the code |

## 3. Implemented Processes

### 3.1 Register Account

1. The visitor submits `POST /auth/register`.
2. `api-gateway` proxies the request to `auth-svc`.
3. `auth-svc` creates the user in PostgreSQL.
4. `auth-svc` creates a session and sets `session_token`.
5. The frontend transitions into the authenticated state.

### 3.2 Log In and Restore Session

1. The visitor submits `POST /auth/login`.
2. `auth-svc` validates the password hash stored in `users`.
3. A new session row is written to `auth_sessions`.
4. The frontend later calls `GET /auth/session` to restore the current user.

### 3.3 Submit Risk Assessment

1. The authenticated user fills the 21-field form.
2. The frontend sends `POST /api/risk`.
3. `api-gateway` validates the session through `auth-svc`.
4. The gateway forwards the feature payload to `ml-api`.
5. `ml-api` loads the model artifact and returns a prediction result.
6. The frontend displays the percentage, category, and message.

### 3.4 Local Deployment and Initialization

1. The install script recreates the k3d cluster.
2. Kubernetes secrets are created for PostgreSQL and admin credentials.
3. Images are built and imported to k3d.
4. Manifests are applied.
5. `data-svc` runs migrations and verifies that the required tables exist.

## 4. Architecture Summary

| Layer | Current Implementation |
|------|------------------------|
| Frontend | Svelte SPA built with Vite and served by Nginx |
| Public backend entrypoint | Go `api-gateway` with Gin |
| Authentication | Go `auth-svc` with Gin |
| ML inference | Python `ml-api` with FastAPI |
| Database | PostgreSQL 16 |
| Deployment | Kubernetes on k3d |

## 5. Public Interfaces

### Public Through Ingress

- `GET /`
- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/logout`
- `GET /auth/session`
- `POST /api/risk`
- `GET /api/features`

### Internal or Port-Forwarded

- `GET /healthz` on `api-gateway`
- `GET /healthz` on `auth-svc`
- `GET /healthz` on `ml-api`
- `GET /features` and `POST /predict` directly on `ml-api`

## 6. Data Ownership

| Data | Current Owner | Notes |
|------|---------------|-------|
| user accounts | `auth-svc` | stored in `users` |
| session records | `auth-svc` | stored in `auth_sessions` |
| model metadata row | migrations | stored in `model_versions` |
| prediction artifact | ML training pipeline | stored in `models/diabrisk_screening.joblib` |
| assessments schema | database only | table exists but current browser flow does not write to it |
| reports schema | database only | table exists but runtime does not use it |
| audit schema | database only | table exists but runtime does not use it |

## 7. Current Constraints

- educational only, not a medical device
- browser prediction flow requires authentication
- predictions are not persisted
- current local workflow is Kubernetes on k3d
- Google OAuth is not part of the active runtime path
- report download and export are not part of the current UI

## 8. Acceptance Points for the Current Implementation

The current implementation can be considered present and working when:

1. the cluster starts through `scripts/install-local-k3d.*`
2. the frontend is reachable on `http://localhost`
3. a user can register or log in
4. `/auth/session` restores the current user
5. `/api/risk` returns a prediction payload for a valid session

## 9. Summary

The repository currently implements a narrow but complete vertical slice:
account registration and login, session-protected access, ML inference, and
result display in a local Kubernetes environment.
