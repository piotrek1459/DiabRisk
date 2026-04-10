# System Vision - DiabRisk

## 1. Overview and Purpose

DiabRisk is a web application that estimates diabetes risk from routine,
non-invasive health indicators. The current implementation is an
educational screening tool: it accepts 21 BRFSS-style inputs, runs a model
through a backend ML service, and returns a percentage, category, and short
message.

The system is intended to demonstrate a complete end-to-end flow from user
input to ML inference inside a small multi-service deployment.

## 2. Product Positioning

| Aspect | Current Position |
|--------|------------------|
| Target users | students, project reviewers, and users exploring an educational screening tool |
| Main need | quick risk estimation from easily available health indicators |
| Product type | Svelte frontend with Go and Python backend services |
| Primary value | clear, local, reproducible implementation of an ML-backed web application |
| Medical status | educational only, not diagnostic |

## 3. Current Scope

The current application includes:

- account registration
- account login and logout
- session restoration with an HttpOnly cookie
- an authenticated risk-assessment form
- prediction through the ML API
- result display in the browser
- local deployment on k3d/Kubernetes

The current application does not include a user-facing flow for:

- assessment history
- report generation or download
- CSV or JSON export
- account deletion
- Google OAuth

## 4. Machine-Learning Component

The current ML flow is based on the BRFSS 2015 diabetes-indicators dataset
kept in `data/raw/`.

### Training Path

- training entrypoint: `src/Ml/main.py`
- raw input: `data/raw/diabetes_012_health_indicators_BRFSS2015.csv`
- resampling: `SMOTE`
- model family in current script: random forest cascade
- saved artifact: `models/diabrisk_screening.joblib`

### Runtime Path

- runtime API: `src/FastAPI/ml_api.py`
- service framework: FastAPI
- prediction input: `{"features": {...}}`
- prediction output:
  - `RiskPercent`
  - `Category`
  - `Message`

## 5. Current Architecture Summary

At runtime, the application consists of:

- a Svelte frontend served by Nginx
- `api-gateway` in Go
- `auth-svc` in Go
- `ml-api` in Python/FastAPI
- PostgreSQL
- `data-svc` as migration runner

Traefik ingress on `http://localhost` exposes:

- `/`
- `/auth/*`
- `/api/*`

## 6. Current User Flow

1. A user opens the frontend on `http://localhost`.
2. The user registers or logs in with email and password.
3. `auth-svc` creates a session and sets `session_token`.
4. The frontend restores the authenticated session through `/auth/session`.
5. The user fills the 21-field form and submits it to `/api/risk`.
6. The gateway validates the session and forwards the request to `ml-api`.
7. The frontend displays the returned result.

## 7. Constraints

- the system is educational and non-diagnostic
- the current browser flow requires login before prediction
- prediction results are not persisted by the active request path
- health endpoints exist inside services, but they are not exposed publicly through ingress
- the supported local deployment path is k3d, not Docker Compose

## 8. Summary

The current DiabRisk implementation focuses on a narrow, working vertical
slice: local account authentication, session-protected access, ML
prediction, and result presentation. The design goal today is clarity of
implementation rather than breadth of product features.
