# API Reference

## Authentication Endpoints

All auth endpoints are proxied through the API Gateway at `/auth/*`.

### POST /auth/register

Register a new user account.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "min6chars",
  "full_name": "John Doe"
}
```

**Response (201):**
```json
{
  "message": "User created successfully",
  "user": { "id": "uuid", "email": "user@example.com", "full_name": "John Doe", "role": "registered" }
}
```

Sets `session_token` HttpOnly cookie (7-day expiry).

### POST /auth/login

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "message": "Logged in successfully",
  "user": { "id": "uuid", "email": "user@example.com", "role": "registered" }
}
```

Sets `session_token` HttpOnly cookie (7-day expiry).

### POST /auth/logout

Revokes the current session. No request body needed.

### GET /auth/session

Returns the current user if session is valid. Returns 401 if not authenticated.

---

## Protected Endpoints

All `/api/*` endpoints require a valid `session_token` cookie.

### POST /api/risk

Submit health features for diabetes risk prediction.

**Request:**
```json
{
  "features": {
    "HighBP": 1, "HighChol": 0, "CholCheck": 1, "BMI": 27.5,
    "Smoker": 0, "Stroke": 0, "HeartDiseaseorAttack": 0,
    "PhysActivity": 1, "Fruits": 1, "Veggies": 1,
    "HvyAlcoholConsump": 0, "AnyHealthcare": 1, "NoDocbcCost": 0,
    "GenHlth": 2, "MentHlth": 0, "PhysHlth": 0, "DiffWalk": 0,
    "Sex": 1, "Age": 7, "Education": 5, "Income": 6
  }
}
```

**Response (200):**
```json
{
  "risk_percent": 23.5,
  "category": "Low",
  "message": "Based on the provided health indicators..."
}
```

### GET /api/features

Returns the list of expected feature names and count.

---

## Health Check

### GET /healthz

Available on all services. Returns `{"status": "ok"}`.
