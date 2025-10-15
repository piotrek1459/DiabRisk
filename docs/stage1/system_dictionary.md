# System Dictionary – DiabRisk: Diabetes Risk Screener Platform

> The System Dictionary defines all terms used in the vision and requirements phase.  
> Each entry specifies its **type** — *Entity*, *Function*, *Activity*, or *Person* — followed by a short, precise definition.

---

## 1  Core Entities

| Term | Type | Definition | Notes / Synonyms |
|------|------|-------------|------------------|
| **User** | Entity | A registered account owner who can log in, submit assessments, and view saved reports. | Actor, Account |
| **Assessment** | Entity | A single request for diabetes-risk estimation containing input features, model version, and resulting probability. | Evaluation, Risk request |
| **Feature** | Entity | Individual input parameter used by the model, e.g. `age`, `bmi`, `glucose`. | Variable |
| **Model Version** | Entity | Frozen ML artifact (weights, preprocessor, metrics) identified by semantic version (e.g., v0.1). | Model artifact |
| **Explanation** | Entity | Per-feature importance values (e.g., SHAP) explaining a specific prediction. | Interpretation |
| **Calibration Curve** | Entity | Data used to show reliability of probabilities (expected vs. observed). | Reliability diagram |
| **Report** | Entity | A generated document (PDF/PNG) summarizing an assessment, its probability, calibration plot, and explanation chart. | Result summary |
| **Dataset Source** | Entity | Reference to public dataset used for model training (e.g., NHANES, Pima). | Data provenance |
| **Auth Session** | Entity | Temporary secure cookie-based session that identifies a logged-in user. | Session token |
| **Event** | Entity | Message passed between services for asynchronous actions (e.g., report generation). | Message |
| **Audit Log** | Entity | Record of critical operations such as exports, deletions, or logins. | Trace |

---

## 2  Supporting Entities (Infrastructure-Level)

| Term | Type | Definition |
|------|------|-------------|
| **Service** | Entity | A microservice providing isolated functionality (auth, data, risk, report). |
| **API Gateway** | Entity | Entry point for external requests; routes to internal services with authentication and rate limits. |
| **Database** | Entity | PostgreSQL schema storing users, assessments, and model metadata. |
| **Object Storage** | Entity | Repository for large artifacts such as ONNX models or generated PDFs. |

---

## 3  Functions (System Capabilities)

| Function | Description | Input / Output |
|-----------|-------------|----------------|
| **Register / Sign In** | Authenticate a user via OAuth or email magic link and create session. | `email` → session cookie |
| **Submit Assessment** | Accept user-provided features, validate ranges, trigger inference. | `features{}` → `assessment_id`, `probability` |
| **Predict Risk** | Core ML inference executed by `risk-svc`; returns probability and confidence band. | `features{}` → `probability` |
| **Explain Prediction** | Generate per-feature contribution values for interpretability. | `assessment_id` → `explanation[]` |
| **Generate Report** | Produce PDF summary with probability, calibration, and explanation plots. | `assessment_id` → `report.pdf` |
| **Save Assessment** | Store assessment results for authenticated user in `data-svc`. | Persisted row |
| **List Assessments** | Retrieve user’s past assessments for `/history` view. | `[assessment]` |
| **Export Assessments** | Create CSV/JSON export for a user. | Download file |
| **Delete Assessment** | Permanently remove an assessment from storage. | — |
| **Delete Account & Data** | Remove user and all related records (GDPR). | — |
| **List Model Versions** | Show available deployed models and metrics. | `[model_version]` |
| **Health Check** | Verify service uptime (used by monitoring). | Status 200 |
| **Generate Audit Entry** | Append a record of user or system actions to audit log. | — |

---

## 4  Activities (User or Background Processes)

| Activity | Description | Performed by |
|-----------|-------------|---------------|
| **Onboarding** | First-time user login; optional tutorial and first assessment creation. | User |
| **Data Entry** | Filling the form with age, BMI, glucose, etc. | User |
| **Risk Evaluation** | ML inference pipeline triggered by submission. | System |
| **Report Generation** | PDF rendering service that assembles charts and text. | System |
| **Daily Backup** | Nightly snapshot of Postgres database and object storage. | System |
| **Training Pipeline** | Offline ML workflow: fetch data → preprocess → train → evaluate → export model → register new version. | ML Engineer |
| **Monitoring & Logging** | Continuous health and metrics tracking by Prometheus/Grafana. | DevOps |
| **User Feedback Collection** | Optional post-assessment survey for educational use. | User |

---

## 5  Persons (Actors and Roles)

| Person | Description |
|--------|-------------|
| **End User** | Main actor using the app anonymously or after login. |
| **Clinician Reviewer** | Secondary actor reviewing exported reports. |
| **Product Owner** | Responsible for backlog and scope control. |
| **Go Developer** | Implements Go microservices and integration tests. |
| **ML Engineer** | Designs, trains, and deploys the ML model and pipelines. |
| **QA Engineer** | Tests functionality, performance, and security aspects. |
| **DevOps Engineer** | Handles CI/CD pipelines, container orchestration, and monitoring. |
| **Instructor** | Course supervisor verifying artifacts and documentation. |

---

## 6  Relationships

User ──┬── submits ──▶ Assessment ──▶ Report
│
└── owns ──▶ [Assessments]
Assessment ── uses ──▶ ModelVersion
Assessment ── generates ──▶ Explanation + CalibrationCurve
Report ── stored in ──▶ ObjectStorage
AuthSession ── belongs to ──▶ User
Event ── triggers ──▶ ReportGeneration


---

## 7  Naming Conventions

| Layer | Convention | Example |
|-------|-------------|----------|
| **Entities (domain)** | PascalCase | `Assessment`, `ModelVersion` |
| **Database tables** | snake_case plural | `assessments`, `model_versions` |
| **API routes** | kebab-case plural | `/api/v1/assessments`, `/api/v1/models` |
| **Events** | dot-separated | `assessment.created`, `report.generated` |
| **Files / Repos** | lowercase | `diabrisk`, `auth-svc`, `risk-svc` |

---

## 8  Glossary (selected)

| Term | Definition |
|------|-------------|
| **AUROC** | Area Under Receiver Operating Characteristic curve; measures model discrimination. |
| **AUPRC** | Area Under Precision-Recall Curve; robust under class imbalance. |
| **SHAP** | SHapley Additive exPlanations; technique for feature attribution. |
| **ONNX** | Open Neural Network Exchange format used for portable model serving. |
| **GDPR** | General Data Protection Regulation (EU); governs user data rights. |
| **CI/CD** | Continuous Integration / Continuous Deployment pipelines. |
| **Microservice** | Independent deployable component exposing a small API. |

---

## 9  Cross-References to Use Cases
| UC ID | Function(s) involved | Primary Entity |
|-------|----------------------|----------------|
| UC-1 | Submit Assessment | Assessment |
| UC-2 | View Risk Result | Assessment |
| UC-3 | Download Report | Report |
| UC-4 | Register / Sign In | User, AuthSession |
| UC-5 | Save Assessment | Assessment |
| UC-6 | List Past Assessments | Assessment |
| UC-7 | Export Data | Assessment |
| UC-10 | Generate Explanations | Explanation |
| UC-11 | View Model Version | ModelVersion |
| UC-15 | Delete Account & Data | User |

---

**Summary:**  
This dictionary defines the key entities, functions, activities, and actors within the DiabRisk system.  
It ensures consistent terminology across design, implementation, and testing stages of the project and serves as a reference for the upcoming UML and use-case specification phase.
