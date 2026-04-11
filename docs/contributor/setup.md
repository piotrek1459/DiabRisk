# Development Setup

## Prerequisites

| Tool | Recommended Use |
|------|------------------|
| Docker | Builds runtime images for all services |
| k3d | Creates the local Kubernetes cluster used by this repo |
| kubectl | Applies manifests, inspects pods, reads logs |
| Node.js 22+ | Frontend development in `frontend/` |
| Go 1.25 | Local development for `api-gateway`, `auth-svc`, `data-svc` |
| Python 3.12 | Local work on `src/FastAPI/` and `src/Ml/` |

## Primary Local Workflow

The supported local environment is k3d plus Kubernetes manifests from `deploy/k8s/`.

### Windows

```powershell
.\scripts\install-local-k3d.ps1
```

### Linux / macOS

```bash
./scripts/install-local-k3d.sh
```

## What the Setup Script Does

The install script is the operational source of truth for local deployment. It:

1. creates a fresh k3d cluster named `diabrisk`
2. creates Kubernetes secrets for PostgreSQL and admin credentials
3. builds Docker images for frontend, gateway, auth, data-svc, and ml-api
4. imports those images into the cluster
5. applies manifests from `deploy/k8s/`
6. waits for rollout of all deployments
7. prints URLs, credentials, and basic log commands

## Default Credentials

| Setting | Value |
|--------|-------|
| Admin email | `admin@diabrisk.local` |
| Admin password | `default_admin_password` |

Override them with `ADMIN_EMAIL` and `ADMIN_PASSWORD` before running the setup script.

## Frontend Development Loop

The usual flow is:

1. start the cluster with the install script
2. keep backend services inside k3d
3. run the frontend locally with Vite

```bash
cd frontend
npm install
npm run dev
```

Vite serves the app on `http://localhost:5173` and proxies `/auth` and `/api` to `http://localhost`.

## Backend Development Loop

If you change a backend service and want to deploy only that service again, rebuild and roll out only the affected image.

### Image Names

| Component | Docker Image Tag | Dockerfile |
|----------|-------------------|------------|
| frontend | `diabrisk-frontend:dev` | `Dockerfile.frontend` |
| api-gateway | `diabrisk-api:dev` | `Dockerfile.api-gateway` |
| auth-svc | `diabrisk-auth:dev` | `Dockerfile.auth-svc` |
| data-svc | `diabrisk-data:dev` | `Dockerfile.data-svc` |
| ml-api | `diabrisk-ml:dev` | `Dockerfile.ml-api` |

### Example: Rebuild and Restart `auth-svc`

```bash
docker build -f Dockerfile.auth-svc -t diabrisk-auth:dev services/auth-svc
k3d image import --cluster diabrisk diabrisk-auth:dev
kubectl rollout restart deployment/auth-svc
kubectl rollout status deployment/auth-svc
```

Use the same pattern for `api-gateway`, `frontend`, `data-svc`, and `ml-api`.

## Useful Local Commands

### Read Logs

```bash
kubectl logs -f -l app=api-gateway
kubectl logs -f -l app=auth-svc
kubectl logs -f -l app=ml-api
kubectl logs -f -l app=data-svc
kubectl logs -f -l app=postgres
```

### Port Forward Internal Services

```bash
kubectl port-forward svc/auth-svc 8081:8081
kubectl port-forward svc/ml-api 8000:8000
kubectl port-forward svc/postgres 5432:5432
```

This is useful for debugging services directly without going through ingress.

## Environment Variables

The repo ships `.env.example` with the main settings:

| Variable | Used By | Purpose |
|---------|---------|---------|
| `DATABASE_URL` | `auth-svc`, `data-svc` | PostgreSQL connection string |
| `POSTGRES_USER` | local/dev docs | Postgres bootstrap user |
| `POSTGRES_PASSWORD` | local/dev docs | Postgres bootstrap password |
| `POSTGRES_DB` | local/dev docs | Postgres database name |
| `PORT` | `auth-svc` | HTTP port for auth service |
| `ADMIN_EMAIL` | setup script, `auth-svc` | Seeded admin account email |
| `ADMIN_PASSWORD` | setup script, `auth-svc` | Seeded admin account password |
| `AUTH_SERVICE_URL` | `api-gateway` | Internal URL used by auth middleware and auth proxy |
| `ML_SERVICE_URL` | `api-gateway` | Internal URL used for prediction and features calls |
| `MODEL_PATH` | `ml-api` | Joblib artifact path loaded by FastAPI |

## Repo-Specific Notes

- `data-svc` currently behaves as a migration runner and schema verifier. It does not expose user-facing CRUD endpoints.
- The frontend uses relative URLs. If you run backend services outside Kubernetes, update environment or proxy configuration accordingly.
- The model artifact expected by `ml-api` lives under `models/` in the repo and is copied into the image built from `Dockerfile.ml-api`.

## Troubleshooting

### Reset the Local Cluster

If the cluster state is inconsistent, rerun the install script - it deletes and recreates the cluster before deploying.

### Check Rollout Status

```bash
kubectl rollout status deployment/data-svc
kubectl rollout status deployment/auth-svc
kubectl rollout status deployment/ml-api
kubectl rollout status deployment/api-gateway
kubectl rollout status deployment/frontend
```

### Inspect Current Resources

```bash
kubectl get pods
kubectl get svc
kubectl get ingress
```
