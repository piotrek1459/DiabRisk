# Architecture

## System Overview

DiabRisk is built with a microservices architecture deployed on Kubernetes (k3d).

```
Traefik Ingress (:80)
  /         -> frontend (Svelte/Nginx :80)
  /auth/*   -> api-gateway (Go/Gin :8080) -> auth-svc (Go/Gin :8081) -> PostgreSQL (:5432)
  /api/*    -> api-gateway (Go/Gin :8080) -> ml-api (Python/FastAPI :8000)
```

## Services

| Service | Language | Port | Purpose |
|---------|----------|------|---------|
| frontend | Svelte 5 / Nginx | 80 | User interface |
| api-gateway | Go / Gin | 8080 | Request routing, auth validation |
| auth-svc | Go / Gin | 8081 | Authentication, sessions |
| ml-api | Python / FastAPI | 8000 | ML model inference |
| data-svc | Go / migrate | - | Database migrations |
| postgres | PostgreSQL 16 | 5432 | Data storage |

## Database Schema

6 main tables: `users`, `auth_sessions`, `model_versions`, `assessments`, `reports`, `audit_logs`.

Migrations are managed by `data-svc` using golang-migrate (see `services/data-svc/migrations/`).

## ML Model

Two-stage cascade Random Forest classifier trained on BRFSS 2015 data:
1. **Model 1 (Screening):** Healthy vs At-Risk
2. **Model 2 (Severity):** Prediabetes vs Diabetes (only for at-risk samples)

Training code: `src/Ml/main.py`
Inference API: `src/FastAPI/ml_api.py`
