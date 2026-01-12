# 🩺 DiabRisk — Diabetes Risk Screener Platform

**DiabRisk** is an educational web application that estimates the risk of developing **Type 2 Diabetes (T2DM)** using routine, non-invasive health metrics such as age, BMI, blood pressure, and glucose.  
The project demonstrates how interpretable **machine learning** can be integrated into a **cloud-native microservices** architecture.

---
## 🚀 Quick Start (Local Kubernetes)

**Requirements:** Docker, k3d, kubectl

```bash
# Deploy all services (builds images, applies manifests, runs migrations)
./scripts/install-local-k3d.sh
```

Then open **http://localhost** to access the application with Google OAuth authentication.

---

## 🏗️ Current Architecture

| Component | Technology | Status | Purpose |
|------------|-------------|---------|----------|
| **Frontend** | Svelte + Vite | ✅ Deployed | SPA with Google Sign-In, risk assessment form, user profile |
| **API Gateway** | Go (Gin) | ✅ Deployed | Routes requests, authentication middleware, CORS handling |
| **Auth Service** | Go (Gin) | ✅ Deployed | Google OAuth 2.0 flow, session management with secure cookies |
| **Data Service** | Go (Gin) | ✅ Deployed | Database migrations, CRUD operations, PostgreSQL integration |
| **ML Service** | Python (FastAPI) | 🔧 Deployed on server | Risk prediction using trained Random Forest model |
| **Database** | PostgreSQL 16 | ✅ Deployed | User profiles, assessments, sessions, audit logs |
| **Deployment** | Kubernetes (k3d) | ✅ Working | Microservices with Traefik ingress on localhost |

---

## 🎯 Current Features (Phase 2 Complete)
- 🔐 **Google OAuth authentication** with secure session management (SHA-512 hashed tokens, 7-day expiry)
- 👤 **User profiles** displaying name and avatar from Google account
- 📊 **Risk assessment form** for collecting health metrics (age, BMI, blood pressure, glucose, etc.)
- 💾 **Persistent storage** of user data and assessments in PostgreSQL
- 🔄 **Database migrations** with automatic schema management
- 🛡️ **Protected routes** requiring authentication to access assessment tool
- 📝 **Audit logging** for GDPR compliance (tracks login, assessment creation, data access)
- 🏗️ **Microservices architecture** with separate auth, data, and gateway services

### 🚧 Upcoming Features (Phase 3)
- ⚙️ **ML-powered risk prediction** integration with deployed FastAPI service
- 📊 **Explainable output** (per-feature importance, calibration chart)
- 🧾 **PDF report generation** and download
- 📂 **Data export** (CSV/JSON) and account deletion UI
- 🧠 **Model card** documenting dataset, metrics, and bias checks
- 💡 **Personalized recommendations** for lifestyle changes 

---

## 🧱 Repository Structure

```
diabrisk/
├── docs/                   # System vision, dictionary, model cards
├── services/
│   ├── api-gateway/        # Entry point, auth middleware, service routing (Go + Gin)
│   ├── auth-svc/           # Google OAuth 2.0, session management (Go + Gin)
│   ├── data-svc/           # Database migrations, CRUD operations (Go + Gin + pgx)
│   ├── report-svc/         # PDF generation (planned)
│   └── ml-api/             # Risk prediction inference (deployed separately)
├── data/                   # ML datasets (processed and raw)
│   ├── processed/          # X_train/test, y_train/test CSVs
│   └── raw/                # BRFSS 2015 dataset
├── src/
│   ├── Ml/                 # Model training scripts
│   └── FastAPI/            # ML inference API (deployed on server)
├── frontend/               # Svelte + Vite SPA with OAuth UI
├── deploy/k8s/             # Kubernetes manifests (postgres, services, ingress)
├── scripts/                # install-local-k3d.sh deployment script
└── Dockerfile.*            # Multi-stage builds for each service
```

---

## 🗄️ Database Schema

The PostgreSQL database stores all application data with the following tables:

| Table | Purpose | Key Relationships |
|-------|---------|-------------------|
| **users** | User profiles from Google OAuth | Primary key for auth_sessions, assessments |
| **auth_sessions** | Secure session tokens (SHA-512 hashed) | Foreign key to users |
| **model_versions** | ML model metadata and performance metrics | Referenced by assessments |
| **assessments** | User health data and risk predictions | Foreign keys to users and model_versions |
| **reports** | Generated PDF reports | Foreign key to assessments |
| **audit_logs** | User actions for GDPR compliance | Foreign key to users |

**Why we need the database:**
- Track assessment history over time for each user
- Enable GDPR compliance (data export, right to be forgotten)
- Model versioning and reproducibility
- Secure session management without JWTs
- Audit trail for regulatory requirements

---

## 🧬 Machine Learning

- **Training data:** BRFSS 2015 — curated diabetes health indicators dataset  
- **Model used:** Random Forest Classifier  
- **Prediction goals:**  
  - Multiclass diabetes status (0 = healthy, 1 = prediabetes, 2 = diabetes)  
  - Estimated risk presented also as a **percentage probability**  
- **Preprocessing:**  
  - Raw features preserved for recommendation logic  
  - Scaled numerical features (BMI, MentHlth, PhysHlth) used for model training   
- **Recommendation engine:**  
  - Generates simple lifestyle-oriented suggestions based on raw patient data  
  - Uses rule-based logic (e.g., high BMI, no physical activity, poor diet)
- **Model training setup:**  
  - Stratified train/test split  
  - Processed datasets stored as CSV (`X_train_processed`, `X_test_processed`, etc.)

Model artifacts and reproducible training scripts are located in `ml/`.  
Each model version iCurrent Implementation)

```
User Browser
   │
   ├──→ http://localhost (Traefik Ingress)
   │
   ▼
Svelte Frontend (Port 80)
   │
   ├──→ /auth/* → API Gateway → Auth Service
   │                             ├─→ Google OAuth 2.0
   │                             └─→ PostgreSQL (auth_sessions)
   │
   ├──→ /api/* → API Gateway (with auth middleware)
   │              │
   │              ├──→ Data Service → PostgreSQL
   │              │     (users, assessments, audit_logs)
   │              │
   │              └──→ ML Service (deployed on server)
   │                    (risk prediction)
   📋 Project Status

### ✅ Phase 1 (Complete)
- System Vision and Dictionary documentation
- Repository structure with microservices
- ML model training pipeline
- Dataset preparation (BRFSS 2015)

### ✅ Phase 2 (Complete)
- PostgreSQL database with migrations (6 tables, seed data)
- Google OAuth 2.0 authentication service
- Secure session management (SHA-512 tokens, 7-day expiry)
- API Gateway with authentication middleware
- Data service for CRUD operations
- Svelte frontend with OAuth UI
- Kubernetes deployment (k3d)
- One-command setup script

### 🚧 Phase 3 (In Progress)
- ML service integration (deployed separately)
- Risk assessment prediction flow
- PDF report generation
- Explainability features (SHAP/LIME)
- Data export and account deletion UI
- Model card documentation
├── auth-svc (port 8081)
├── data-svc (port 8082, runs migrations)
├── api-gateway (port 8080)
└── frontend (port 80)
```

**Authentication Flow:**
1. User clicks "Sign in with Google"
2. Redirected to Google OAuth consent screen
3. Callback to `/auth/google/callback` with code
4. Auth service exchanges code for user info
5. Creates/finds user in database
6. Generates secure session token (SHA-512 hash)
7. Sets HttpOnly cookie with domain=localhost
8. User can access protected `/api/risk` endpoint
 └──→ Report Service (PDF generator)
```

Session-based authentication (OAuth / magic link) allows users to view and manage their report history securely.

---



## 📜 License
Educational and non-commercial use only.  
This software is **not a certified medical device** and must not be used for clinical diagnosis or treatment.

---

## 👥 Authors
- **Team:** Software Engineering Project (Silesian University of Technology)  
- **Contact:** <piotrek1459@gmail.com>,
               <jeremiszcotka7@gmail.com>
---

**DiabRisk** — empowering early awareness through open, interpretable machine learning.
