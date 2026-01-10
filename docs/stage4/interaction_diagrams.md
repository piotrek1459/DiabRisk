# Interaction Diagrams

## Overview
This document contains the interaction diagrams for the two most complex use cases in the DiabRisk system. These diagrams are modeled as sequence diagrams using PlantUML syntax to illustrate the flow of interactions between system components.

---

## Use Case 1: Google OAuth Authentication Flow
### Description:
This use case models the complete authentication flow when a user signs in with their Google account. It involves multiple redirects, token exchanges, user creation/retrieval, session creation with secure token hashing, and cookie-based session management across multiple microservices.

### Sequence Diagram (PlantUML):
```plantuml
@startuml
actor User
participant "Frontend\n(Svelte SPA)" as Frontend
participant "API Gateway\n(Go)" as Gateway
participant "Auth Service\n(Go)" as AuthSvc
participant "Google OAuth\nServers" as Google
database "PostgreSQL\n(users, auth_sessions)" as DB

== Authentication Initiation ==
User -> Frontend: Click "Sign in with Google"
Frontend -> Frontend: Redirect to /auth/google/login

Frontend -> Gateway: GET /auth/google/login
Gateway -> AuthSvc: Forward GET /auth/google/login
AuthSvc -> AuthSvc: Generate random state (32 bytes)
AuthSvc -> AuthSvc: Set oauth_state cookie (10 min expiry)
AuthSvc -> Gateway: Redirect to Google OAuth URL\n(with state, client_id, redirect_uri)
Gateway -> Frontend: HTTP 307 Redirect
Frontend -> User: Browser redirects to Google

== Google OAuth Consent ==
User -> Google: Authenticate & grant consent
Google -> Google: Generate authorization code
Google -> Frontend: Redirect to /auth/google/callback\n(with code, state)

== Token Exchange & Session Creation ==
Frontend -> Gateway: GET /auth/google/callback?code=xxx&state=yyy
Gateway -> AuthSvc: Forward callback request (with cookies)
AuthSvc -> AuthSvc: Validate state matches oauth_state cookie
AuthSvc -> AuthSvc: Clear oauth_state cookie

AuthSvc -> Google: Exchange code for access token
Google -> AuthSvc: Return access_token

AuthSvc -> Google: GET /oauth2/v2/userinfo\n(with access_token)
Google -> AuthSvc: Return user info (email, name, picture, id)

AuthSvc -> DB: SELECT user by google_id
alt User exists
    DB -> AuthSvc: Return existing user
else User does not exist
    AuthSvc -> DB: INSERT new user\n(email, google_id, full_name, picture_url)
    DB -> AuthSvc: Return created user
end

AuthSvc -> AuthSvc: Generate session token (64 bytes random)
AuthSvc -> AuthSvc: Hash token with SHA-512
AuthSvc -> DB: INSERT auth_session\n(user_id, token_hash, expires_at: 7 days)
DB -> AuthSvc: Session created

AuthSvc -> Gateway: Set session_token cookie\n(HttpOnly, 7 days expiry)
Gateway -> Frontend: HTTP 307 Redirect to "/"
Frontend -> User: Show authenticated homepage

== Session Validation (Subsequent Requests) ==
User -> Frontend: Access protected resource
Frontend -> Gateway: Request with session_token cookie
Gateway -> AuthSvc: GET /auth/session (validate)
AuthSvc -> AuthSvc: Hash session_token from cookie
AuthSvc -> DB: SELECT session by token_hash\n(WHERE expires_at > NOW AND is_revoked = FALSE)
DB -> AuthSvc: Return session & user data
AuthSvc -> Gateway: Return user info (JSON)
Gateway -> Frontend: User authenticated
Frontend -> User: Display protected content

@enduml
```

---

## Use Case 2: Diabetes Risk Prediction with ML Service
### Description:
This use case models the complete flow of submitting health assessment data, authentication validation, ML model inference on an external server, and result display. It demonstrates the interaction between frontend, API gateway authentication middleware, and the remote ML prediction service.

### Sequence Diagram (PlantUML):
```plantuml
@startuml
actor User
participant "Frontend\n(Svelte SPA)" as Frontend
participant "API Gateway\n(Go)" as Gateway
participant "Auth Service\n(Go)" as AuthSvc
participant "ML Service\n(FastAPI)\n65.109.169.137:8000" as MLService

== User Fills Assessment Form ==
User -> Frontend: Enter health parameters\n(BMI, Age, HighBP, HighChol, etc.)
Frontend -> Frontend: Validate input ranges\n(BMI: 15-60, Age: 18-100)

User -> Frontend: Click "Submit Risk Assessment"

== Authentication Check ==
Frontend -> Gateway: POST /api/risk\n(with session_token cookie)\nBody: {21 health features}
Gateway -> Gateway: Extract session_token from cookie
Gateway -> AuthSvc: GET /auth/session\n(validate token)
AuthSvc -> AuthSvc: Hash session_token
AuthSvc -> AuthSvc: Query database for valid session
alt Session valid
    AuthSvc -> Gateway: 200 OK + user info JSON
    Gateway -> Gateway: Set user context in request
else Session invalid/expired
    AuthSvc -> Gateway: 401 Unauthorized
    Gateway -> Frontend: 401 Unauthorized
    Frontend -> User: Redirect to login
    note right: Flow terminates here if unauthenticated
end

== Risk Prediction ==
Gateway -> Gateway: Marshal health features to JSON
Gateway -> MLService: POST /predict\nContent-Type: application/json\nBody: {features}

MLService -> MLService: Load model artifact (model.joblib)
MLService -> MLService: Validate feature names
MLService -> MLService: Build feature vector (21 features)
MLService -> MLService: Standardize/normalize features
MLService -> MLService: Execute ML inference\n(Random Forest / XGBoost)
MLService -> MLService: Apply probability calibration
MLService -> MLService: Classify risk category\n(Low: <5%, Moderate: 5-15%, High: >15%)

MLService -> Gateway: 200 OK\nBody: {\n  "probability": 0.08,\n  "risk_category": "Moderate",\n  "model_version": "1.0.0"\n}

== Response Processing ==
Gateway -> Gateway: Parse ML response JSON
Gateway -> Frontend: 200 OK\nBody: {prediction data}

Frontend -> Frontend: Store result in component state
Frontend -> Frontend: Format probability as percentage
Frontend -> User: Display risk assessment result\n• Risk Category: Moderate\n• Probability: 8.0%\n• Visual: colored indicator

alt User wants to save result
    User -> Frontend: Click "Save to History"
    note right: Future: Save to data-svc\n(assessments table)
end

@enduml
```

---

## Diagram Generation Instructions

To generate diagrams from these PlantUML specifications:

1. **Online**: Copy the code block and paste into [PlantUML Online Editor](http://www.plantuml.com/plantuml)
2. **VS Code**: Install the "PlantUML" extension and preview the markdown file
3. **CLI**: Use `plantuml interaction_diagrams.md` to generate PNG/SVG files
4. **Docker**: `docker run -v $(pwd):/data plantuml/plantuml interaction_diagrams.md`

---

## Complexity Analysis

### Use Case 1 Complexity Factors:
- **7 system components** involved (User, Frontend, Gateway, AuthSvc, Google, DB, Cookies)
- **State management** across multiple redirects (oauth_state cookie)
- **Security operations** (token generation, SHA-512 hashing, session validation)
- **Error handling** for invalid states, expired tokens, database failures
- **Asynchronous redirects** through browser navigation
- **Cookie domain/expiry configuration** for session persistence

### Use Case 2 Complexity Factors:
- **Authentication middleware** integration (session validation before processing)
- **External service communication** with remote ML API (network latency, timeouts)
- **Data transformation** (JSON marshaling/unmarshaling, feature vector construction)
- **ML model inference** (loading artifacts, preprocessing, prediction, calibration)
- **Error propagation** from ML service to frontend (connection failures, validation errors)
- **Stateful client-side** result display and formatting