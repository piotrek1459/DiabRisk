# Use Case Model Documentation – DiabRisk

**Project:** DiabRisk – Diabetes Risk Screener Platform  
**Date:** January 2026  
**Version:** 1.0  
**Phase:** Stage 2 - Use Case Modeling  

---

## 1. Actor Identification

### 1.1 Actor Candidates from Project Dictionary

Based on the analysis of the Project Dictionary (Section 2.4 - Persons), the following actor candidates were identified:

| Actor Candidate | Source Term | Selection Decision |
|-----------------|-------------|-------------------|
| End User | "End User" (PD 2.4) | ✅ Selected - Primary system actor |
| Clinician Reviewer | "Clinician Reviewer" (PD 2.4) | ✅ Selected - Secondary actor for exported data |
| Product Owner | "Product Owner" (PD 2.4) | ❌ Excluded - Project management role, not system user |
| Go Developer | "Go Developer" (PD 2.4) | ❌ Excluded - Development role, not end user |
| ML Engineer | "ML Engineer" (PD 2.4) | ❌ Excluded - Development role, not end user |
| Frontend Developer | "Frontend Developer" (PD 2.4) | ❌ Excluded - Development role, not end user |
| QA Engineer | "QA Engineer" (PD 2.4) | ❌ Excluded - Testing role, not end user |
| DevOps Engineer | "DevOps Engineer" (PD 2.4) | ✅ Selected as "System Administrator" |
| Course Instructor | "Course Instructor" (PD 2.4) | ❌ Excluded - Academic supervisor, not system user |

### 1.2 Actor Refinement and Specialization

**Primary Actors:**
- **Anonymous User**: Unauthenticated visitor who can submit one-off assessments
- **Registered User**: Authenticated user with account who can save and track assessments
  - *Generalization*: Registered User inherits all capabilities of Anonymous User

**Secondary Actors:**
- **Clinician Reviewer**: Healthcare professional who reviews exported patient assessment data
- **System Administrator**: Technical role managing system health and performance

### 1.3 Final Actor List

| Actor ID | Actor Name | Description | Primary/Secondary |
|----------|-----------|-------------|-------------------|
| A1 | Anonymous User | Performs one-time risk assessments without registration | Primary |
| A2 | Registered User | Authenticated user with persistent account and history | Primary |
| A3 | Clinician Reviewer | Reviews exported assessment data as supplementary information | Secondary |
| A4 | System Administrator | Monitors system health and manages technical operations | Secondary |

---

## 2. Use Case Identification

### 2.1 Use Case Candidates from Project Dictionary

**From Functions (Section 2.2):**
- Register / Sign In ✅
- Submit Assessment ✅
- Predict Risk ✅ (internal, included in Submit Assessment)
- Explain Prediction ✅
- Generate Report ✅
- Save Assessment ✅
- List Assessments ✅
- Export Assessments ✅
- Delete Assessment ✅
- Delete Account & Data ✅
- List Model Versions ✅
- Health Check ✅

**From Activities (Section 2.3):**
- Onboarding ✅ (merged with Register/Sign In)
- Data Entry ✅ (part of Submit Assessment)
- Risk Evaluation ✅ (internal process)
- Report Generation ✅
- Monitoring & Logging ✅ (admin function)

**From Prioritized Use Cases (Part IV):**
- All 15 use cases (UC-1 through UC-15) identified and included

### 2.2 Use Case Categorization

**Core Assessment (Anonymous + Registered):**
- UC-1: Submit Anonymous Assessment
- UC-2: View Risk Result with Visualization
- UC-15: Acknowledge Legal Disclaimer

**Authentication:**
- UC-4: Register / Sign In via OAuth

**Account Management (Registered only):**
- UC-5: Save Assessment to Account
- UC-6: List Past Assessments
- UC-7: Delete Individual Assessment
- UC-12: Delete Account & All Data (GDPR)

**Reports & Export:**
- UC-3: Download PDF Report
- UC-8: Export All Assessments (CSV/JSON)

**Advanced Analysis:**
- UC-9: View Model Version and Metrics
- UC-10: Generate SHAP Explanations
- UC-11: View Calibration Plot

**Administrative:**
- UC-13: Admin Health Check Endpoint
- UC-14: Rate Limit Handling

---

## 3. Relationship Identification

### 3.1 Actor Generalizations

```
Registered User --|> Anonymous User
```

**Rationale**: A registered user can perform all actions available to anonymous users, plus additional authenticated features. This is a classic "is-a" relationship.

### 3.2 <<include>> Relationships (Mandatory)

| Base Use Case | Included Use Case | Justification |
|---------------|-------------------|---------------|
| UC-1: Submit Anonymous Assessment | UC-15: Acknowledge Legal Disclaimer | Every assessment submission requires disclaimer acknowledgment |
| UC-1: Submit Anonymous Assessment | Validate Input Data | All submissions must validate 21 health parameters |
| UC-1: Submit Anonymous Assessment | Execute ML Inference | Core functionality - assessment requires prediction |
| UC-2: View Risk Result | Execute ML Inference | Displaying results requires inference completion |
| UC-3: Download PDF Report | Generate Report Content | PDF creation requires content generation |
| UC-5: Save Assessment to Account | UC-1: Submit Assessment | Saving requires an assessment to exist |

### 3.3 <<extend>> Relationships (Optional)

| Base Use Case | Extending Use Case | Extension Point | Condition |
|---------------|-------------------|-----------------|-----------|
| UC-2: View Risk Result | UC-10: Generate SHAP Explanations | After displaying result | User requests detailed explanation |
| UC-2: View Risk Result | UC-11: View Calibration Plot | After displaying result | User requests model reliability info |
| UC-3: Download PDF Report | UC-10: Generate SHAP Explanations | During report generation | User opted for detailed report |

**Rationale**: SHAP explanations and calibration plots are optional enhancements that provide additional insight but are not required for core functionality.

### 3.4 Use Case Dependencies

**Prerequisite Relationships:**
- UC-5, UC-6, UC-7, UC-12 require UC-4 (user must be authenticated)
- UC-3 requires UC-2 (must have a result to generate report)
- UC-7 requires UC-6 (must list assessments before deleting one)

---

## 4. Use Case Diagram

The use case diagram is provided in PlantUML format in `use_case_diagram.puml`.

**Key Features:**
- Actor hierarchy showing generalization (Registered User inherits from Anonymous User)
- System boundary clearly delineated
- 15 prioritized use cases organized by functional area
- <<include>> relationships showing mandatory dependencies
- <<extend>> relationships showing optional enhancements
- Association lines connecting actors to their use cases

**To generate the diagram:**
```bash
# Using PlantUML CLI
plantuml use_case_diagram.puml

# Using VS Code PlantUML extension
# Open file and use preview command
```

---

## 5. Use Case Prioritization Matrix

| Use Case | Priority | Actor | Complexity | Phase | Dependencies |
|----------|----------|-------|------------|-------|--------------|
| UC-1 | Critical | Anonymous User | Medium | 1 | UC-15 |
| UC-2 | Critical | Anonymous User | Low | 1 | UC-1 |
| UC-15 | Critical | Anonymous User | Low | 1 | None |
| UC-4 | High | Registered User | High | 1 | None |
| UC-5 | High | Registered User | Medium | 1 | UC-1, UC-4 |
| UC-6 | High | Registered User | Low | 1 | UC-4 |
| UC-3 | High | Registered User | High | 1 | UC-2, UC-4 |
| UC-12 | High | Registered User | Medium | 1 | UC-4 |
| UC-7 | Medium | Registered User | Low | 1 | UC-4, UC-6 |
| UC-9 | Medium | Anonymous User | Low | 1 | None |
| UC-14 | Medium | System | Medium | 1 | None |
| UC-8 | Medium | Registered User | Medium | 2 | UC-4, UC-6 |
| UC-10 | Medium | Registered User | High | 2 | UC-2 |
| UC-11 | Medium | Registered User | High | 2 | UC-2 |
| UC-13 | Low | Admin | Low | 1 | None |

---

## 6. Actor-Use Case Matrix

| Use Case | Anonymous User | Registered User | Clinician | Admin |
|----------|----------------|-----------------|-----------|-------|
| UC-1: Submit Assessment | ✓ | ✓ (inherited) | | |
| UC-2: View Result | ✓ | ✓ (inherited) | | |
| UC-3: Download Report | | ✓ | | |
| UC-4: Sign In | | ✓ | | |
| UC-5: Save Assessment | | ✓ | | |
| UC-6: List Assessments | | ✓ | | |
| UC-7: Delete Assessment | | ✓ | | |
| UC-8: Export Data | | ✓ | ✓ | |
| UC-9: View Model Info | ✓ | ✓ (inherited) | | |
| UC-10: SHAP Explanations | | ✓ | | |
| UC-11: Calibration Plot | | ✓ | | |
| UC-12: Delete Account | | ✓ | | |
| UC-13: Health Check | | | | ✓ |
| UC-14: Rate Limiting | System | System | System | ✓ |
| UC-15: Disclaimer | ✓ | ✓ (inherited) | | |

---

## 7. Modeling Decisions and Rationale

### 7.1 Why Anonymous User and Registered User are Separate Actors

Despite the generalization relationship, we model them separately because:
1. **Different authentication states** lead to different system behaviors
2. **Distinct use case sets** - many features require authentication
3. **Clear business rules** - anonymous assessments are ephemeral, registered ones persist
4. **Security boundaries** - different authorization levels

### 7.2 Why Clinician Reviewer is Included

Although secondary, the clinician actor is included because:
- Explicitly mentioned in System Vision (Section 1.1)
- Has a specific use case (UC-8: reviewing exported data)
- Represents a different usage pattern than end users
- Important for complete system understanding

### 7.3 Internal Use Cases (Not Shown)

Some technical processes are **not** modeled as use cases:
- **Validate Input Data**: Technical subprocess, shown as included
- **Execute ML Inference**: Internal system operation, not user-initiated
- **Generate Report Content**: Technical subprocess of UC-3

**Rationale**: Use case diagrams show user-visible functionality. Internal technical steps are documented in sequence diagrams and activity diagrams (future stages).

### 7.4 Rate Limit Handling (UC-14)

Modeled as a system-level use case because:
- Transparent to normal users (automatic)
- Only visible when limits exceeded
- Managed by administrators
- Cross-cutting concern affecting all use cases

---

## 8. Traceability to System Dictionary

| Use Case | Dictionary Functions | Dictionary Activities |
|----------|---------------------|---------------------|
| UC-1 | Submit Assessment | Data Entry, Risk Evaluation |
| UC-2 | Predict Risk | Risk Evaluation |
| UC-3 | Generate Report | Report Generation |
| UC-4 | Register / Sign In | Onboarding |
| UC-5 | Save Assessment | - |
| UC-6 | List Assessments | - |
| UC-7 | Delete Assessment | - |
| UC-8 | Export Assessments | - |
| UC-9 | List Model Versions | - |
| UC-10 | Explain Prediction | Risk Evaluation |
| UC-11 | (Calibration Curve entity) | - |
| UC-12 | Delete Account & Data | - |
| UC-13 | Health Check | Monitoring & Logging |
| UC-14 | (Rate limiting) | Monitoring & Logging |
| UC-15 | (Legal disclaimer) | Onboarding |

---

## 9. Summary Statistics

- **Total Actors**: 4 (2 primary, 2 secondary)
- **Total Use Cases**: 15
- **Actor Generalizations**: 1
- **<<include>> Relationships**: 6
- **<<extend>> Relationships**: 3
- **Phase 1 Use Cases**: 11
- **Phase 2 Use Cases**: 4

---

**Document Prepared By:** Development Team  
**Review Status:** Ready for Stage 2 Submission  
**Next Steps:** Detailed use case specifications, sequence diagrams, activity diagrams

---

**End of Use Case Model Documentation**
