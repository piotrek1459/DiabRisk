# DiabRisk

DiabRisk is an educational web application for estimating diabetes risk
from 21 health indicators derived from the BRFSS 2015 dataset. The current
system is a small Kubernetes-based stack with a Svelte frontend, Go
services for gateway and authentication, a FastAPI ML service, and
PostgreSQL.

## Current Application Flow

1. The user registers or logs in with email and password.
2. `auth-svc` creates a `session_token` HttpOnly cookie.
3. After login, the frontend shows the risk form with 21 fields.
4. The frontend submits the form to `POST /api/risk`.
5. `api-gateway` validates the session through `auth-svc`.
6. `ml-api` loads the model artifact and returns:
   - `RiskPercent`
   - `Category`
   - `Message`

## Quick Start

### Requirements

- Docker
- k3d
- kubectl

### Windows

```powershell
.\scripts\install-local-k3d.ps1
```

### Linux / macOS

```bash
./scripts/install-local-k3d.sh
```

Then open `http://localhost`.

### Default Local Credentials

The setup script seeds a default admin account:

- Email: `admin@diabrisk.local`
- Password: `default_admin_password`

You can override both values with `ADMIN_EMAIL` and `ADMIN_PASSWORD`
before running the install script.

## Runtime Components

| Component | Technology | Responsibility |
|----------|------------|----------------|
| frontend | Svelte 5 + Vite + Nginx | login/register UI, risk form, result display |
| api-gateway | Go + Gin | browser-facing backend entrypoint, auth middleware, proxying |
| auth-svc | Go + Gin | registration, login, logout, session validation |
| ml-api | Python 3.12 + FastAPI | feature validation, model loading, prediction |
| data-svc | Go binary | waits for PostgreSQL, runs migrations, verifies schema |
| postgres | PostgreSQL 16 | stores users, sessions, and seeded schema tables |
| ingress | Traefik on k3d | exposes `/`, `/auth/*`, and `/api/*` on `localhost` |

## Repository Map

| Path | Purpose |
|------|---------|
| `frontend/` | Svelte application |
| `services/api-gateway/` | gateway service and auth middleware |
| `services/auth-svc/` | local account authentication and sessions |
| `services/data-svc/` | migration runner and schema verification |
| `src/FastAPI/` | ML inference API |
| `src/Ml/` | model training scripts |
| `models/` | current joblib model artifact |
| `deploy/k8s/` | Kubernetes manifests used by local deployment scripts |
| `docs/` | user, contributor, and design documentation |
| `data/` | raw dataset and precomputed processed splits |

## Current Scope

The current user-facing application provides:

- account registration
- account login and logout
- session restoration after refresh
- a protected risk-assessment form
- immediate display of prediction results

The current application does not provide a user-facing flow for:

- assessment history
- report download
- CSV or JSON export
- account deletion
- feature-level explanation charts

## Documentation

| Folder | Audience |
|--------|----------|
| `docs/user/` | people using the application |
| `docs/contributor/` | developers extending the repository |
| `docs/design/` | implementation-focused design documentation |

## License

Educational and non-commercial use only.

DiabRisk is not a certified medical device and must not be used for
clinical diagnosis or treatment.
