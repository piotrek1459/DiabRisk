# API Reference

## Routing Overview

The browser reaches the system through Traefik ingress on
`http://localhost`.

| Public Path | Implemented In | Upstream Target | Auth Required |
|------------|----------------|-----------------|---------------|
| `/` | frontend | static files from Nginx | No |
| `/auth/*` | api-gateway | proxied to `auth-svc` | No |
| `/api/*` | api-gateway | handled in gateway, may call `ml-api` | Yes |

`/healthz` endpoints exist on individual services, but they are not routed
through the current ingress manifest.

## Authentication Endpoints

All auth routes are exposed through `api-gateway` and forwarded to
`auth-svc` as-is.

### POST /auth/register

Creates a user account and an authenticated session.

**Request**

```json
{
  "email": "user@example.com",
  "password": "min6chars",
  "full_name": "John Doe"
}
```

**Response (201)**

```json
{
  "message": "User created successfully",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "full_name": "John Doe",
    "role": "registered",
    "created_at": "2026-04-10T20:00:00Z"
  }
}
```

**Notes**

- sets `session_token` as an HttpOnly cookie
- password validation is enforced in `auth-svc` with `min=6`

### POST /auth/login

Authenticates an existing user and creates a fresh session.

**Request**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200)**

```json
{
  "message": "Logged in successfully",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "full_name": "John Doe",
    "role": "registered",
    "created_at": "2026-04-10T20:00:00Z"
  }
}
```

### POST /auth/logout

Revokes the current session if the cookie is present.

**Response (200)**

```json
{
  "message": "Logged out successfully"
}
```

If the browser has no session cookie, the current implementation still
returns:

```json
{
  "message": "Already logged out"
}
```

### GET /auth/session

Validates the current session and returns the authenticated user object.

**Response (200)**

```json
{
  "id": "uuid",
  "email": "user@example.com",
  "full_name": "John Doe",
  "role": "registered",
  "created_at": "2026-04-10T20:00:00Z"
}
```

**Response (401)**

```json
{
  "error": "Not authenticated"
}
```

or

```json
{
  "error": "Invalid or expired session"
}
```

## Protected Gateway Endpoints

All `/api/*` endpoints require a valid `session_token` cookie. The gateway
validates the cookie by calling `GET /auth/session` on `auth-svc`.

### POST /api/risk

Submits health features for diabetes risk prediction.

The gateway accepts either of these request shapes.

**Canonical request**

```json
{
  "features": {
    "HighBP": 1,
    "HighChol": 0,
    "CholCheck": 1,
    "BMI": 27.5,
    "Smoker": 0,
    "Stroke": 0,
    "HeartDiseaseorAttack": 0,
    "PhysActivity": 1,
    "Fruits": 1,
    "Veggies": 1,
    "HvyAlcoholConsump": 0,
    "AnyHealthcare": 1,
    "NoDocbcCost": 0,
    "GenHlth": 2,
    "MentHlth": 0,
    "PhysHlth": 0,
    "DiffWalk": 0,
    "Sex": 1,
    "Age": 7,
    "Education": 5,
    "Income": 6
  }
}
```

**Also accepted by the current gateway**

```json
{
  "HighBP": 1,
  "HighChol": 0,
  "CholCheck": 1,
  "BMI": 27.5,
  "Smoker": 0,
  "Stroke": 0,
  "HeartDiseaseorAttack": 0,
  "PhysActivity": 1,
  "Fruits": 1,
  "Veggies": 1,
  "HvyAlcoholConsump": 0,
  "AnyHealthcare": 1,
  "NoDocbcCost": 0,
  "GenHlth": 2,
  "MentHlth": 0,
  "PhysHlth": 0,
  "DiffWalk": 0,
  "Sex": 1,
  "Age": 7,
  "Education": 5,
  "Income": 6
}
```

**Response (200)**

```json
{
  "RiskPercent": 0.23,
  "Category": "low",
  "Message": "Low risk detected. No immediate action required."
}
```

**Notes**

- `RiskPercent` is returned as a `0..1` float
- the frontend formats the value as percent for display
- `api-gateway` does not persist results to PostgreSQL
- the current gateway forwards the parsed ML response body and does not
  preserve the upstream ML status code

### GET /api/features

Returns the feature list expected by `ml-api`.

**Response (200)**

```json
{
  "feature_names": [
    "HighBP",
    "HighChol",
    "CholCheck",
    "BMI",
    "Smoker",
    "Stroke",
    "HeartDiseaseorAttack",
    "PhysActivity",
    "Fruits",
    "Veggies",
    "HvyAlcoholConsump",
    "AnyHealthcare",
    "NoDocbcCost",
    "GenHlth",
    "MentHlth",
    "PhysHlth",
    "DiffWalk",
    "Sex",
    "Age",
    "Education",
    "Income"
  ],
  "count": 21
}
```

## Direct Service Endpoints

These endpoints are useful when port-forwarding internal services for
debugging.

### api-gateway

`GET /healthz` returns:

```json
{
  "status": "ok"
}
```

### auth-svc

`GET /healthz` returns a simple health payload from the auth service.

### ml-api

`GET /healthz` returns:

```json
{
  "status": "ok"
}
```

`GET /features` returns the same payload shape as `GET /api/features`.

`POST /predict` accepts the canonical `{"features": {...}}` body and
returns:

```json
{
  "RiskPercent": 0.23,
  "Category": "low",
  "Message": "Low risk detected. No immediate action required."
}
```

## Accessing Internal Health Endpoints

Example:

```bash
kubectl port-forward svc/api-gateway 8080:8080
curl http://localhost:8080/healthz
```

Use the same pattern for `auth-svc` on `8081` and `ml-api` on `8000`.
