# Technology Stack - DiabRisk

**Scope:** technologies actively used by the current repository

## Frontend

| Area | Technology | Current Role |
|------|------------|--------------|
| UI | Svelte 5 | components and state |
| Build tooling | Vite 7 | dev server and production build |
| Production serving | Nginx 1.27 Alpine | serves compiled frontend assets |

## Go Services

| Area | Technology | Current Role |
|------|------------|--------------|
| Language | Go 1.25 | `api-gateway`, `auth-svc`, `data-svc` |
| HTTP framework | Gin | routing, middleware, JSON handling |
| Database access | pgx v5 | PostgreSQL access in Go services |
| SQL migrations | golang-migrate | migration execution in `data-svc` |
| Password hashing | bcrypt | local account authentication |
| UUID generation | `github.com/google/uuid` | identifiers for users and sessions |

## Python and ML

| Area | Technology | Current Role |
|------|------------|--------------|
| Language | Python 3.12 | runtime for ML service and training scripts |
| API framework | FastAPI | ML HTTP API |
| Validation | Pydantic | request and response models |
| Numerical data | NumPy | feature arrays and inference input |
| Tabular data | pandas | dataset loading in training |
| ML library | scikit-learn | random forest training and inference |
| Resampling | imbalanced-learn | `SMOTE` in the training script |
| Model persistence | joblib | saved artifact in `models/` |

## Data and Infrastructure

| Area | Technology | Current Role |
|------|------------|--------------|
| Database | PostgreSQL 16 | persistent storage for users, sessions, and seeded schema tables |
| Containers | Docker | builds runtime images |
| Orchestration | Kubernetes | local deployment target |
| Local cluster | k3d | lightweight Kubernetes for development |
| Ingress | Traefik | routes `/`, `/auth/*`, and `/api/*` on `localhost` |

## Runtime Configuration

Current runtime configuration is driven by:

- environment variables in service containers
- Kubernetes secrets for database and admin credentials
- Kubernetes manifests in `deploy/k8s/`
- Vite proxy settings in `frontend/vite.config.js`

## Model Artifact Format

The active ML service loads a joblib artifact:

- current artifact path in the repo: `models/diabrisk_screening.joblib`
- deployment path in `ml-api`: `/app/models/diabrisk_screening.joblib`

The current runtime does not use ONNX.

## Technologies Not in Active Runtime Use

The current repository does not actively use these technologies in the
running path:

- Google OAuth
- Docker Compose as the main local workflow
- NATS
- object storage
- Prometheus
- Grafana
- a separate report-generation service

Some legacy files and earlier notes still mention them, but they are not
part of the current execution path.
