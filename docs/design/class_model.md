# Class Modeling Documentation - DiabRisk

**Scope:** conceptual model of the current implementation

## 1. Purpose

This document describes the main runtime objects and persistent records that
exist in the current DiabRisk implementation. It is a conceptual model, not
an object-oriented description of the Go or Python code.

## 2. Current Class Set

| ID | Class | Type | Status |
|----|-------|------|--------|
| C1 | User | entity | active in runtime |
| C2 | AuthSession | entity | active in runtime |
| C3 | FeatureSet | value object | active in runtime |
| C4 | PredictionResult | value object | active in runtime |
| C5 | ModelArtifact | entity | active in runtime |
| C6 | ModelVersionRecord | entity | present in schema |
| C7 | AssessmentRecord | entity | present in schema, not written by active browser flow |
| C8 | ReportRecord | entity | present in schema, not used by active browser flow |
| C9 | AuditLogRecord | entity | present in schema, not used by active browser flow |

## 3. Class Summaries

### C1 User

- stored in `users`
- identified by UUID
- contains email, password hash, optional full name, role, timestamps
- owned by `auth-svc`

### C2 AuthSession

- stored in `auth_sessions`
- linked to one user
- contains hashed token, expiry, and revocation state
- used by `auth-svc` and validated by `api-gateway`

### C3 FeatureSet

- 21 health indicators submitted from the frontend
- serialized as JSON
- validated in `ml-api`

### C4 PredictionResult

- returned by `ml-api`
- consists of:
  - `RiskPercent`
  - `Category`
  - `Message`

### C5 ModelArtifact

- current artifact file: `models/diabrisk_screening.joblib`
- loaded in memory by `ml-api`
- contains feature names and trained model objects

### C6 ModelVersionRecord

- stored in `model_versions`
- seeded by migrations
- not used by the active browser flow

### C7 AssessmentRecord

- stored in `assessments`
- schema exists, but the current `/api/risk` path does not write records

### C8 ReportRecord

- stored in `reports`
- schema exists, but the current runtime does not generate reports

### C9 AuditLogRecord

- stored in `audit_logs`
- schema exists, but the current runtime path does not write audit events

## 4. Relationships

| From | To | Relationship | Notes |
|------|----|--------------|-------|
| User | AuthSession | one-to-many | one user may have multiple session records over time |
| FeatureSet | PredictionResult | produces | one submitted feature set produces one runtime result |
| PredictionResult | ModelArtifact | depends on | result is computed using the loaded model artifact |
| User | AssessmentRecord | one-to-many | schema-level relation only; not active in current browser flow |
| AssessmentRecord | ReportRecord | one-to-many | schema-level relation only; not active in current browser flow |
| User | AuditLogRecord | one-to-many | schema-level relation only; not active in current browser flow |
| AssessmentRecord | ModelVersionRecord | many-to-one | schema-level relation through model metadata |

## 5. Active Runtime Path

The classes that participate directly in the current browser flow are:

1. `User`
2. `AuthSession`
3. `FeatureSet`
4. `PredictionResult`
5. `ModelArtifact`

The remaining classes are currently database-schema concepts rather than
runtime objects used by the browser path.

## 6. Current Constraints

- the system does not persist prediction results after `POST /api/risk`
- there is no active report-generation workflow
- there is no active audit-log writing workflow
- the active authentication model is session-based, not token-based JWT auth

## 7. Summary

The conceptual model of the current implementation is centered on local
users, session records, submitted feature sets, and returned prediction
results. Additional database tables already exist, but they are not yet
part of the active runtime path.
