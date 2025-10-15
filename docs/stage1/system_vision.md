# System Vision – DiabRisk: Diabetes Risk Screener Platform

## 1  Overview and Purpose
**DiabRisk** is a web-based application that estimates an individual’s risk of developing **Type 2 Diabetes Mellitus (T2DM)** using routine, non-invasive health data.  
Its goal is to provide students, early-career clinicians, and health-conscious individuals with a **transparent, data-driven risk score** that encourages preventive action long before clinical onset.

Unlike commercial black-box calculators, DiabRisk is **open, reproducible, and educational**.  
It publishes model cards, allows users to explore feature contributions, and stores results privately under the user’s own account.

---

## 2  Problem Statement
Worldwide, over 400 million adults have diabetes, yet many remain undiagnosed until complications appear.  
Early identification of risk through simple metrics such as **age, BMI, blood pressure, and fasting glucose** can guide lifestyle interventions that significantly delay or prevent T2DM.  
Existing solutions either:
- require integration with closed electronic-health-record systems,  
- hide model details, or  
- neglect user privacy.

**DiabRisk** addresses these gaps with a minimal, privacy-respecting web app powered by an interpretable ML model trained on public datasets.

---

## 3  Product Positioning
| Aspect | DiabRisk |
|---------|-----------|
| **Target users** | Students, health-care trainees, and individuals tracking metabolic health |
| **Need** | Quick estimation of 3-year diabetes risk from easily available metrics |
| **Product type** | Web application (Svelte frontend + Go/Python microservices backend) |
| **Value proposition** | Open, explainable ML model with personal report history |
| **Primary benefit** | Early awareness and self-education—not diagnosis |

---

## 4  System Scope
The system covers:
1. Data entry of personal health parameters.  
2. ML inference that produces a probability of developing T2DM.  
3. Generation of a PDF report containing score, calibration chart, and explanation.  
4. Optional user authentication to store past assessments and export data.

Excluded from Phase 1:
- EHR integration or real medical certification,
- mobile apps,
- multilingual UI beyond English.

---

## 5  Machine-Learning Component
**Objective** – Binary classification: predict current or near-future T2DM status.

**Datasets**
- **NHANES (1999–2020)** — comprehensive U.S. population health survey.  
  Labels derived from *HbA1c ≥ 6.5 %* or *fasting glucose ≥ 126 mg/dL*.
- **UCI Pima Indians Diabetes** — baseline for prototyping.  
- **CDC BRFSS** — large self-reported dataset for external validation.

**Feature set**
Demographics (age, sex), anthropometrics (BMI, waist), vitals (SBP, DBP), basic labs (glucose, HbA1c, HDL, LDL, TG), lifestyle flags (smoking, activity).

**Models**
- Logistic Regression (baseline, interpretable)
- Gradient-Boosted Trees (XGBoost/LightGBM)
- Calibration via Platt or Isotonic regression  
  → export to ONNX for serving

**Evaluation**
AUROC, AUPRC, Brier score, Expected Calibration Error;  
subgroup metrics (sex, age, BMI bands).  
Each training run produces a signed **model card** stored with metrics and hyper-parameters.

---

## 6  Architecture Summary
**Frontend:** Svelte SPA  
– routes: `/new`, `/result/:id`, `/history`, `/login`  
– communicates via fetch API with credentials included  

**Backend:** Go microservices (+ Python for ML)
api-gateway (Go) – routing, rate-limit, CORS
auth-svc (Go) – OAuth / magic-link login, session JWT
data-svc (Go) – store assessments (Postgres)
risk-svc (Python) – load ONNX model and predict
report-svc (Go) – render PDF reports
event-bus (NATS) – async events (report requested → generate)

**Persistence:** PostgreSQL + object storage (S3-compatible).  
**Deployment:** Docker Compose → Kubernetes (stage 2).  
**Observability:** Prometheus metrics and Grafana dashboards.

---

## 7  Authentication and User Flow
- Unauthenticated users may perform one-off assessments (ephemeral).  
- Authenticated users (via Google OAuth or email link) obtain an **HttpOnly session cookie**.  
- Their assessments are stored in the `assessments` table and visible on `/history`.  
- Each report belongs to exactly one user; GDPR deletion removes all associated data.

---

## 8  Functional Goals (≤ 20 use cases)
1. Submit Assessment (anonymous)  
2. View Risk Result  
3. Download Report (PDF)  
4. Register / Sign In  
5. Save Assessment to Account  
6. List Past Assessments  
7. Delete Assessment  
8. Export All Assessments (CSV/JSON)  
9. View Model Version and Metrics  
10. Generate Explanations  
11. View Calibration Plot  
12. Delete Account & Data  
13. Health Check (Admin)  
14. Rate-limit Handling  
15. Legal Disclaimer Acknowledgement  

---

## 9  Non-Functional Requirements
| Category | Requirement |
|-----------|-------------|
| **Performance** | Prediction ≤ 800 ms P50 |
| **Scalability** | ≤ 1 000 concurrent users (Phase 1) |
| **Security** | All cookies HttpOnly + SameSite Lax; HTTPS only |
| **Privacy** | Data encryption at rest + deletion on request |
| **Reliability** | 99 % availability per semester |
| **Maintainability** | Service boundaries ≤ 300 LOC each avg |
| **Test coverage** | ≥ 60 % unit/integration Go services |

---

## 10  Stakeholders and Roles
| Role | Description |
|------|-------------|
| **End User** | Enters health data, views and downloads reports |
| **Clinician Reviewer** | Optionally reviews exported reports |
| **Product Owner** | Defines backlog, ensures scope ≤ 20 UCs |
| **Go Developer** | Implements gateway, auth-svc, data-svc |
| **ML Engineer** | Trains & deploys models, writes model card |
| **QA Engineer** | Tests endpoints and UI flows |
| **DevOps Engineer** | CI/CD, monitoring, infrastructure |
| **Instructor** | Academic supervisor, reviews artifacts |

---

## 11  Constraints and Assumptions
- The system is **educational**, not a certified medical device.  
- Only publicly licensed datasets are used for training.  
- Internet access is required; offline mode is not supported in Phase 1.  
- English UI only; Polish labels may be added later.  
- Team size 3–5 students; project duration ≈ 1 semester.

---

## 12  Success Metrics (Phase 1)
- AUROC ≥ 0.80 on held-out NHANES-like cohort.  
- 100 % reproducible training pipeline (hash-matched artifacts).  
- < 5 % bug escape after system testing.  
- User satisfaction survey ≥ 4/5 average.

---

## 13  Risks and Mitigations
| Risk | Mitigation |
|-------|-------------|
| **Data heterogeneity** | Write dataset-specific adapters + tests |
| **Model overfitting** | Cross-validation + early stopping |
| **Regulatory concerns** | Prominent “educational only” disclaimer |
| **Polyglot complexity** | Containerize each service; docker-compose dev |
| **Scope creep** | Limit to ≤ 20 UCs and 3 sprints |

---

## 14  Phase 1 Deliverables
1. **System Vision** (this document)  
2. **System Dictionary** of entities, functions, activities, persons  
3. **GitHub repository:** `https://github.com/piotrek1459/diabrisk`  
4. **Service skeletons:** auth-svc, risk-svc, data-svc, report-svc  
5. **Model Card Template** + example training notebook  
6. **Docker Compose** file for local integration tests  

---

**Summary:**  
DiabRisk combines interpretable machine learning with cloud-native microservices to provide a secure, educational diabetes-risk screener.  
Phase 1 establishes a vertical slice from data entry to PDF report generation, a reproducible ML pipeline, and a privacy-safe user login enabling personal report history.
