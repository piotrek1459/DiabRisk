# Use Case Model Documentation - DiabRisk

**Scope:** current implementation in the repository

## 1. Actors

| Actor ID | Actor | Description |
|----------|-------|-------------|
| A1 | Visitor | opens the application and can register or log in |
| A2 | Authenticated User | logged-in user who can access the risk form and submit predictions |
| A3 | Developer / Operator | runs the local cluster and inspects service health |

## 2. Current Use Cases

| Use Case ID | Use Case | Primary Actor |
|-------------|----------|---------------|
| UC-1 | Register account | Visitor |
| UC-2 | Log in | Visitor |
| UC-3 | Restore session | Authenticated User |
| UC-4 | Submit risk assessment | Authenticated User |
| UC-5 | View prediction result | Authenticated User |
| UC-6 | Log out | Authenticated User |
| UC-7 | Read feature contract | Authenticated User / Developer |
| UC-8 | Check internal service health | Developer / Operator |
| UC-9 | Deploy local environment | Developer / Operator |

## 3. Use Case Notes

### UC-1 Register Account

- user submits email, password, and optional full name
- system creates a user and initial session

### UC-2 Log In

- user submits email and password
- system validates credentials and returns a fresh session

### UC-3 Restore Session

- frontend calls `GET /auth/session`
- system validates `session_token` and returns the current user

### UC-4 Submit Risk Assessment

- authenticated user submits the 21-field form
- gateway validates the session and forwards the payload to the ML service

### UC-5 View Prediction Result

- frontend displays the returned percentage, category, and message

### UC-6 Log Out

- frontend calls `POST /auth/logout`
- system revokes the current session

### UC-7 Read Feature Contract

- frontend or developer reads `GET /api/features` or `GET /features`
- system returns the 21 expected feature names

### UC-8 Check Internal Service Health

- operator checks service-local `/healthz`
- this is an internal maintenance flow, not a public browser route

### UC-9 Deploy Local Environment

- operator runs `scripts/install-local-k3d.ps1` or `.sh`
- system builds images, applies manifests, and initializes the stack

## 4. Relationships

### Actor Relationship

`Authenticated User` is the post-login state of `Visitor`.

### Mandatory Internal Dependencies

- UC-4 includes session validation before ML inference
- UC-5 depends on UC-4
- UC-6 depends on an existing session

## 5. Actor-Use Case Matrix

| Use Case | Visitor | Authenticated User | Developer / Operator |
|----------|---------|--------------------|----------------------|
| UC-1 Register account | yes | | |
| UC-2 Log in | yes | | |
| UC-3 Restore session | | yes | |
| UC-4 Submit risk assessment | | yes | |
| UC-5 View prediction result | | yes | |
| UC-6 Log out | | yes | |
| UC-7 Read feature contract | | yes | yes |
| UC-8 Check internal service health | | | yes |
| UC-9 Deploy local environment | | | yes |

## 6. Out-of-Scope Use Cases for the Current Runtime

The current repository does not implement browser-facing use cases for:

- anonymous prediction
- assessment history
- report download
- data export
- account deletion
- Google OAuth login

## 7. Summary

The current use-case model is intentionally narrow. It centers on one
working authenticated flow: register or log in, submit the risk form, read
the result, and log out.
