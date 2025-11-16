# 🩺 DiabRisk — Diabetes Risk Screener Platform

**DiabRisk** is an educational web application that estimates the risk of developing **Type 2 Diabetes (T2DM)** using routine, non-invasive health metrics such as age, BMI, blood pressure, and glucose.  
The project demonstrates how interpretable **machine learning** can be integrated into a **cloud-native microservices** architecture.

---

## 🚀 Overview

| Component | Technology | Purpose |
|------------|-------------|----------|
| **Frontend** | [Svelte](https://svelte.dev/) | Minimal SPA for data entry and result visualization |
| **Backend** | [Go (Gin/Fiber)](https://go.dev/) | API gateway, auth, data, and reporting services |
| **ML Service** | [Python + FastAPI](https://fastapi.tiangolo.com/) | Inference using ONNX model |
| **Database** | PostgreSQL | Persistent storage for users and assessments |
| **Architecture** | Microservices + REST + Docker | Scalable, modular design |
| **Deployment** | Docker Compose → Kubernetes | Cloud-ready setup for later stages |

---

## 🎯 Core Features
- ⚙️ **ML-powered risk prediction** for Type 2 diabetes  
- 📊 **Explainable output** (per-feature importance, calibration chart)  
- 🧾 **PDF report generation** and download  
- 🔐 **Optional login** to save personal assessment history  
- 📂 **Data export** (CSV/JSON) and account deletion for GDPR compliance  
- 🧠 **Transparent model card** documenting dataset, metrics, and bias checks 
- 💌 **Predictions** providing possible life changes to lover a given prediction 

---

## 🧱 Repository Structure

```
diabrisk/
├── docs/                   # System vision, dictionary, model cards
├── services/
│   ├── api-gateway/        # Entry point for all requests
│   ├── auth-svc/           # OAuth/magic-link authentication (Go)
│   ├── data-svc/           # CRUD operations and persistence (Go + Postgres)
│   ├── report-svc/         # PDF generation (Go)
│   └── risk-svc/           # ML inference API (Python FastAPI + ONNX)
├── ml/                     # Training pipeline and dataset adapters
│   ├── data_ingest/
│   ├── features/
│   ├── train/
│   ├── eval/
│   └── notebooks/
├── frontend/               # Svelte single-page app
├── deploy/                 # Docker Compose / K8s manifests
└── .github/workflows/      # CI/CD configuration
```

---

## ⚡ Quick Start (Development)

**Requirements:** Docker, Docker Compose, and Python ≥3.10

```bash
# Clone repository
git clone https://github.com/piotrek1459/diabrisk.git
cd diabrisk

# Start all services
docker compose up --build
```

Then open **http://localhost:5173** (Svelte frontend).  
The backend gateway runs on **http://localhost:8080**.

---

## 🧬 Machine Learning

- **Training data:**  BRFSS2015
- **Models:** Logistic Regression & Gradient Boosted Trees (XGBoost)  
- **Evaluation metrics:** AUROC, AUPRC, Brier score, calibration error  
- **Serving format:** ONNX  
- **Explainability:** SHAP-like local feature importance  

Model artifacts and reproducible training scripts are located in `ml/`.  
Each model version is accompanied by a **model card** documenting performance and dataset provenance.

---

## 🧠 Architecture (Simplified)

```
Svelte Frontend
   │
   ▼
API Gateway (Go)
 ├──→ Auth Service (login/session)
 ├──→ Risk Service (Python ML inference)
 ├──→ Data Service (store assessments)
 └──→ Report Service (PDF generator)
```

Session-based authentication (OAuth / magic link) allows users to view and manage their report history securely.

---

## 🧩 Phase 1 Deliverables
- ✅ System Vision (docs/system_vision.md)  
- ✅ System Dictionary (docs/system_dictionary.md)  
- ✅ Repo scaffold with service folders and documentation  
- ✅ Model training plan & dataset references  

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
