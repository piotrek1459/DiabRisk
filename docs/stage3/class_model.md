# Class Modeling Documentation – DiabRisk

**Project:** DiabRisk – Diabetes Risk Screener Platform  
**Date:** January 2026  
**Version:** 1.0  
**Phase:** Stage 3 - Class Modeling  

---

## 1. Introduction to Class Modeling

### 1.1 What is a CASE Tool?

**CASE (Computer-Aided Software Engineering) Tool** is specialized software that supports software development activities, particularly system modeling and design.

**Types of CASE Tools:**

| Category | Purpose | Examples |
|----------|---------|----------|
| **Upper CASE** | Requirements analysis, design modeling | Enterprise Architect, Rational Rose, Visual Paradigm |
| **Lower CASE** | Code generation, testing | Eclipse, IntelliJ IDEA with UML plugins |
| **Integrated CASE** | Full lifecycle support | IBM Rational Suite, Microsoft Visio |
| **Lightweight CASE** | Simple modeling | PlantUML, draw.io, StarUML |

**Common Features:**
- UML diagram creation (class, sequence, use case, etc.)
- Model validation and consistency checking
- Code generation from models
- Reverse engineering (code to models)
- Version control integration
- Documentation generation

**For this project**, we use **PlantUML** - a lightweight, text-based CASE tool that:
- Uses simple text syntax to define diagrams
- Integrates with version control (Git)
- Generates professional diagrams automatically
- Supports all UML 2.x diagram types
- Enables collaborative modeling through text files

---

## 2. Class Identification Process

### 2.1 Candidates from Project Dictionary

Based on the Initial System Specification (Part II: Project Dictionary), we identify class candidates from three categories:

#### Category 1: Subjects (Domain Objects)

| Term | Source Section | Class Candidate | Rationale |
|------|---------------|-----------------|-----------|
| Assessment | PD 2.1 Core Entities | ✅ **Assessment** | Primary business object - represents single risk evaluation |
| Report | PD 2.1 Core Entities | ✅ **Report** | Concrete artifact - PDF document with results |
| Model Version | PD 2.1 Core Entities | ✅ **ModelVersion** | Versionable artifact - ML model with metadata |
| Feature | PD 2.1 Core Entities | ✅ **Feature** | Value object - health parameter measurement |

#### Category 2: Conceptual Entities (Abstract Concepts)

| Term | Source Section | Class Candidate | Rationale |
|------|---------------|-----------------|-----------|
| Risk Score | PD 2.1 Core Entities | ✅ **RiskScore** | Value object - encapsulates probability and category |
| Explanation | PD 2.1 Core Entities | ✅ **Explanation** | Complex structure - SHAP values per feature |
| Calibration Curve | PD 2.1 Core Entities | ✅ **CalibrationCurve** | Data structure - reliability visualization |
| Auth Session | PD 2.1 Core Entities | ✅ **AuthSession** | Temporal entity - user authentication state |
| Audit Log | PD 2.1 Core Entities | ✅ **AuditLog** | Event record - security and compliance |

#### Category 3: Persons (Actors)

| Term | Source Section | Class Candidate | Rationale |
|------|---------------|-----------------|-----------|
| End User | PD 2.4 Persons | ✅ **User** | Primary actor with account and data |
| Clinician Reviewer | PD 2.4 Persons | ⚠️ **Specialization of User** | Could be User role/type, not separate class |
| System Administrator | PD 2.4 Persons | ⚠️ **Specialization of User** | Could be User role/type, not separate class |

**Decision**: Model as single `User` class with `role` attribute (anonymous, registered, clinician, admin) rather than inheritance hierarchy. Rationale: Simple role-based access control, no behavioral differences.

### 2.2 Classes Excluded from Model

| Term | Why Excluded |
|------|--------------|
| Dataset Source | External entity, not managed by system |
| Service | Implementation detail, not domain concept |
| API Gateway | Infrastructure component |
| Database | Technical infrastructure |
| Object Storage | Technical infrastructure |
| Event | Messaging implementation detail |

---

## 3. Class Diagram Overview

The class diagram represents the **domain model** - core business concepts and their relationships.

### 3.1 Final Class List

| Class ID | Class Name | Stereotype | Description |
|----------|-----------|------------|-------------|
| C1 | User | Entity | Person with account, can have assessments |
| C2 | Assessment | Entity | Single risk evaluation request with results |
| C3 | Feature | Value Object | Individual health parameter measurement |
| C4 | RiskScore | Value Object | Probability, category, and message |
| C5 | ModelVersion | Entity | Versioned ML model artifact |
| C6 | Report | Entity | Generated PDF document |
| C7 | Explanation | Value Object | Per-feature SHAP importance values |
| C8 | CalibrationCurve | Value Object | Reliability data for visualization |
| C9 | AuthSession | Entity | Temporary authentication state |
| C10 | AuditLog | Entity | Immutable event record |

**Stereotypes:**
- **Entity**: Has identity, lifecycle, persistence
- **Value Object**: No identity, immutable, compared by value

---

## 4. Class Attributes and Methods

### 4.1 User Class

**Attributes:**
- `id: UUID` - Unique identifier
- `email: String` - User email address
- `role: UserRole` - Enum: anonymous, registered, clinician, admin
- `createdAt: DateTime` - Account creation timestamp
- `lastLoginAt: DateTime` - Last successful login

**Methods:**
- `register(email: String): User` - Create new user account
- `authenticate(credentials): AuthSession` - Login and create session
- `deleteAccount(): void` - GDPR-compliant deletion
- `exportData(): File` - Export all user data

**Responsibilities:** Manages user identity, authentication state, and ownership of assessments.

### 4.2 Assessment Class

**Attributes:**
- `id: UUID` - Unique identifier
- `userId: UUID` - Foreign key to User (nullable for anonymous)
- `modelVersionId: UUID` - Model used for prediction
- `createdAt: DateTime` - Submission timestamp
- `isPersisted: Boolean` - Whether saved to database

**Methods:**
- `submit(features: Feature[]): Assessment` - Create new assessment
- `calculateRisk(): RiskScore` - Execute ML inference
- `save(user: User): void` - Persist to account
- `generateReport(): Report` - Create PDF summary
- `explain(): Explanation` - Compute SHAP values
- `delete(): void` - Remove from storage

**Responsibilities:** Central aggregate root - coordinates assessment workflow, contains features and results.

### 4.3 Feature Class (Value Object)

**Attributes:**
- `name: String` - Feature identifier (e.g., "BMI", "Age")
- `value: Float` - Measured value
- `unit: String` - Measurement unit (optional)
- `isValid: Boolean` - Validation status

**Methods:**
- `validate(): Boolean` - Check value within acceptable range
- `normalize(): Float` - Apply feature scaling

**Responsibilities:** Encapsulates single health parameter with validation logic.

### 4.4 RiskScore Class (Value Object)

**Attributes:**
- `probability: Float` - Risk probability [0.0, 1.0]
- `category: RiskCategory` - Enum: low, medium, high
- `message: String` - Interpretive text
- `confidenceInterval: Tuple<Float, Float>` - 95% CI bounds

**Methods:**
- `categorize(probability: Float): RiskCategory` - Map probability to category
- `format(): String` - Display-friendly representation

**Responsibilities:** Immutable result value with business logic for categorization.

### 4.5 ModelVersion Class

**Attributes:**
- `id: UUID` - Unique identifier
- `version: String` - Semantic version (e.g., "v0.1.0")
- `trainedAt: DateTime` - Training completion timestamp
- `auroc: Float` - Model performance metric
- `auprc: Float` - Precision-recall metric
- `calibrationError: Float` - Expected calibration error
- `artifactPath: String` - ONNX file location
- `isActive: Boolean` - Currently deployed

**Methods:**
- `load(): Model` - Load ONNX artifact into memory
- `predict(features: Feature[]): RiskScore` - Execute inference
- `getMetrics(): Map<String, Float>` - Return performance statistics

**Responsibilities:** Represents frozen model artifact with versioning and metadata.

### 4.6 Report Class

**Attributes:**
- `id: UUID` - Unique identifier
- `assessmentId: UUID` - Associated assessment
- `generatedAt: DateTime` - Creation timestamp
- `format: ReportFormat` - Enum: PDF, PNG
- `storagePath: String` - Object storage location
- `fileSize: Integer` - Document size in bytes

**Methods:**
- `generate(assessment: Assessment): Report` - Create report from assessment
- `download(): File` - Retrieve from storage
- `delete(): void` - Remove from storage

**Responsibilities:** Manages report lifecycle and storage.

### 4.7 Explanation Class (Value Object)

**Attributes:**
- `featureImportances: Map<String, Float>` - SHAP values per feature
- `baseValue: Float` - Model baseline prediction
- `method: String` - Explanation algorithm used

**Methods:**
- `compute(model: ModelVersion, features: Feature[]): Explanation` - Calculate SHAP
- `getTopFeatures(n: Integer): List<Tuple>` - Most influential features
- `visualize(): Chart` - Generate importance bar chart

**Responsibilities:** Encapsulates feature attribution for transparency.

### 4.8 CalibrationCurve Class (Value Object)

**Attributes:**
- `bins: List<CalibrationBin>` - Probability bins with observed frequencies
- `numBins: Integer` - Number of bins (typically 10)
- `ece: Float` - Expected calibration error

**CalibrationBin:**
- `predictedProbability: Float` - Bin center
- `observedFrequency: Float` - Actual diabetes rate
- `count: Integer` - Number of samples

**Methods:**
- `compute(predictions: List<Float>, actuals: List<Boolean>): CalibrationCurve`
- `plot(): Chart` - Generate calibration plot

**Responsibilities:** Stores reliability data for model trustworthiness visualization.

### 4.9 AuthSession Class

**Attributes:**
- `id: UUID` - Session identifier
- `userId: UUID` - Associated user
- `token: String` - Encrypted session token
- `createdAt: DateTime` - Session start
- `expiresAt: DateTime` - Expiration timestamp
- `ipAddress: String` - Client IP for security
- `userAgent: String` - Browser information

**Methods:**
- `create(user: User): AuthSession` - Initialize session
- `validate(): Boolean` - Check if still valid
- `renew(): void` - Extend expiration
- `destroy(): void` - Logout

**Responsibilities:** Manages stateful authentication between requests.

### 4.10 AuditLog Class

**Attributes:**
- `id: UUID` - Unique identifier
- `userId: UUID` - Actor (nullable for system events)
- `action: AuditAction` - Enum: login, assessment_created, data_exported, account_deleted
- `timestamp: DateTime` - Event occurrence
- `ipAddress: String` - Source IP
- `metadata: JSON` - Additional context

**Methods:**
- `log(user: User, action: AuditAction, metadata: Map): AuditLog` - Create entry
- `query(filters: AuditFilter): List<AuditLog>` - Compliance reporting

**Responsibilities:** Immutable event log for security auditing and compliance (GDPR).

---

## 5. Relationships and Associations

### 5.1 Association Table

| ID | From Class | To Class | Relation Name | Cardinality | Reading Direction | Type |
|----|-----------|----------|---------------|-------------|-------------------|------|
| R1 | User | Assessment | owns | 1 → 0..* | User owns many Assessments | Aggregation |
| R2 | User | AuthSession | has | 1 → 0..1 | User has at most one active AuthSession | Composition |
| R3 | User | AuditLog | triggers | 1 → 0..* | User triggers many AuditLog entries | Association |
| R4 | Assessment | Feature | contains | 1 → 21 | Assessment contains exactly 21 Features | Composition |
| R5 | Assessment | RiskScore | produces | 1 → 1 | Assessment produces one RiskScore | Composition |
| R6 | Assessment | ModelVersion | uses | * → 1 | Many Assessments use one ModelVersion | Association |
| R7 | Assessment | Report | generates | 1 → 0..1 | Assessment generates optional Report | Aggregation |
| R8 | Assessment | Explanation | explains | 1 → 0..1 | Assessment has optional Explanation | Composition |
| R9 | Report | CalibrationCurve | includes | 1 → 1 | Report includes one CalibrationCurve | Composition |

### 5.2 Relationship Descriptions

#### R1: User owns Assessment (1 → 0..*)
- **Type**: Aggregation (hollow diamond on User side)
- **Cardinality**: One user can have zero or many assessments
- **Direction**: User → Assessment
- **Lifecycle**: Assessments can exist without User (anonymous), but if owned, deletion of User removes Assessments
- **Navigability**: Bidirectional - User can access assessments, Assessment knows its owner

#### R2: User has AuthSession (1 → 0..1)
- **Type**: Composition (solid diamond on User side)
- **Cardinality**: One user has at most one active session
- **Direction**: User → AuthSession
- **Lifecycle**: Session cannot exist without User; destroyed when User logs out or deleted
- **Navigability**: Bidirectional

#### R3: User triggers AuditLog (1 → 0..*)
- **Type**: Association (no diamond)
- **Cardinality**: One user triggers many audit events
- **Direction**: User → AuditLog
- **Lifecycle**: AuditLogs persist even after User deletion (anonymized)
- **Navigability**: Unidirectional - AuditLog references User, but User doesn't need to access logs

#### R4: Assessment contains Feature (1 → 21)
- **Type**: Composition (solid diamond on Assessment side)
- **Cardinality**: Exactly 21 features per assessment (fixed by ML model)
- **Direction**: Assessment → Feature
- **Lifecycle**: Features destroyed when Assessment deleted
- **Navigability**: Unidirectional - Assessment accesses Features

#### R5: Assessment produces RiskScore (1 → 1)
- **Type**: Composition
- **Cardinality**: Exactly one risk score per assessment
- **Direction**: Assessment → RiskScore
- **Lifecycle**: RiskScore lifecycle tied to Assessment
- **Navigability**: Unidirectional

#### R6: Assessment uses ModelVersion (* → 1)
- **Type**: Association
- **Cardinality**: Many assessments use same model version
- **Direction**: Assessment → ModelVersion
- **Lifecycle**: Independent - ModelVersion exists regardless of Assessments
- **Navigability**: Unidirectional - Assessment references model

#### R7: Assessment generates Report (1 → 0..1)
- **Type**: Aggregation
- **Cardinality**: Assessment may have zero or one report
- **Direction**: Assessment → Report
- **Lifecycle**: Report can outlive Assessment (stored separately)
- **Navigability**: Bidirectional

#### R8: Assessment explains Explanation (1 → 0..1)
- **Type**: Composition
- **Cardinality**: Optional explanation (computed on demand)
- **Direction**: Assessment → Explanation
- **Lifecycle**: Explanation destroyed with Assessment
- **Navigability**: Unidirectional

#### R9: Report includes CalibrationCurve (1 → 1)
- **Type**: Composition
- **Cardinality**: Each report includes one calibration curve
- **Direction**: Report → CalibrationCurve
- **Lifecycle**: CalibrationCurve embedded in Report
- **Navigability**: Unidirectional

### 5.3 Association Constraints

**Constraints and Business Rules:**

1. **Anonymous Assessments**: 
   - `Assessment.userId` can be NULL
   - Constraint: `Assessment.isPersisted = false WHERE userId IS NULL`

2. **Active Session Limit**:
   - Constraint: `COUNT(AuthSession WHERE userId = X AND expiresAt > NOW()) ≤ 1`

3. **Feature Count**:
   - Constraint: `COUNT(Feature WHERE assessmentId = X) = 21`

4. **Model Version Activation**:
   - Constraint: `COUNT(ModelVersion WHERE isActive = true) = 1`

5. **Report Generation Precondition**:
   - Constraint: `Report.assessmentId → Assessment.id AND Assessment.riskScore IS NOT NULL`

---

## 6. Generalizations (Inheritance)

### 6.1 Considered Generalizations

| Potential Hierarchy | Decision | Rationale |
|---------------------|----------|-----------|
| User ← AnonymousUser, RegisteredUser | ❌ Rejected | Simple role attribute sufficient; no behavioral differences |
| User ← Clinician, Admin | ❌ Rejected | Role-based access control via `role` enum cleaner |
| Report ← PDFReport, PNGReport | ❌ Rejected | Format as attribute; same structure and behavior |
| Feature ← NumericFeature, BinaryFeature | ❌ Rejected | All features numeric in current model |

**Conclusion**: No inheritance hierarchies in current domain model. Favor composition and attributes over inheritance for flexibility.

### 6.2 Future Extensibility

If system evolves to support multiple assessment types (diabetes, cardiovascular, etc.):

```
Assessment (abstract)
  ├─ DiabetesAssessment
  ├─ CardiovascularAssessment
  └─ HypertensionAssessment
```

Currently not needed per YAGNI (You Aren't Gonna Need It) principle.

---

## 7. Enumerations

### 7.1 Enumeration Types

```java
enum UserRole {
    ANONYMOUS,      // Unauthenticated user
    REGISTERED,     // Authenticated user with account
    CLINICIAN,      // Healthcare professional
    ADMIN           // System administrator
}

enum RiskCategory {
    LOW,            // probability < 0.3
    MEDIUM,         // 0.3 ≤ probability < 0.6
    HIGH            // probability ≥ 0.6
}

enum ReportFormat {
    PDF,            // Adobe PDF document
    PNG             // Image format (for embedding)
}

enum AuditAction {
    USER_REGISTERED,
    USER_LOGGED_IN,
    USER_LOGGED_OUT,
    ASSESSMENT_SUBMITTED,
    ASSESSMENT_SAVED,
    ASSESSMENT_DELETED,
    REPORT_GENERATED,
    DATA_EXPORTED,
    ACCOUNT_DELETED
}
```

---

## 8. Class Diagram

The class diagram is provided in PlantUML format in `class_diagram.puml`.

**Diagram Features:**
- 10 domain classes with attributes and key methods
- 9 associations with cardinalities and relation names
- Composition and aggregation relationships clearly marked
- Enumerations for categorical attributes
- Stereotypes (<<entity>>, <<value object>>)

**To generate the diagram:**
```bash
# Using PlantUML CLI
plantuml class_diagram.puml

# Using VS Code PlantUML extension
# Preview: Alt+D (Windows) or Cmd+D (Mac)

# Using online editor
# Copy content to: http://www.plantuml.com/plantuml/
```

---

## 9. Design Patterns Applied

### 9.1 Aggregate Pattern (Domain-Driven Design)

**Assessment as Aggregate Root:**
- Controls access to Features, RiskScore, Explanation
- Ensures invariants (e.g., exactly 21 features)
- Single entry point for operations

### 9.2 Value Object Pattern

**RiskScore, Feature, Explanation, CalibrationCurve:**
- Immutable after creation
- No identity (compared by value)
- Can be safely shared
- Reduce coupling

### 9.3 Repository Pattern (Implied)

Classes marked as `<<entity>>` will have corresponding repositories:
- `UserRepository`
- `AssessmentRepository`
- `ModelVersionRepository`
- `ReportRepository`
- `AuditLogRepository`

---

## 10. Traceability Matrix

### 10.1 Classes to Dictionary Entities

| Class | Project Dictionary Reference | Section |
|-------|------------------------------|---------|
| User | User | PD 2.1 |
| Assessment | Assessment | PD 2.1 |
| Feature | Feature | PD 2.1 |
| RiskScore | Risk Score | PD 2.1 |
| ModelVersion | Model Version | PD 2.1 |
| Report | Report | PD 2.1 |
| Explanation | Explanation | PD 2.1 |
| CalibrationCurve | Calibration Curve | PD 2.1 |
| AuthSession | Auth Session | PD 2.1 |
| AuditLog | Audit Log | PD 2.1 |

### 10.2 Classes to Use Cases

| Class | Primary Use Cases |
|-------|-------------------|
| User | UC-4 (Register/Sign In), UC-12 (Delete Account) |
| Assessment | UC-1 (Submit), UC-2 (View Result), UC-5 (Save), UC-6 (List), UC-7 (Delete) |
| Report | UC-3 (Download PDF), UC-8 (Export) |
| Explanation | UC-10 (Generate SHAP) |
| CalibrationCurve | UC-11 (View Calibration Plot) |
| ModelVersion | UC-9 (View Model Info) |
| AuthSession | UC-4 (Authentication) |
| AuditLog | UC-12 (GDPR Deletion), UC-13 (Admin Health Check) |

---

## 11. Implementation Considerations

### 11.1 Persistence Mapping

**PostgreSQL Tables:**
- `users` (C1)
- `assessments` (C2)
- `features` (C3) - JSONB column in assessments or separate table
- `risk_scores` (C4) - embedded in assessments table
- `model_versions` (C5)
- `reports` (C6) - metadata only, files in object storage
- `auth_sessions` (C9)
- `audit_logs` (C10)

**Object Storage (S3):**
- Report PDFs (actual files)
- ONNX model artifacts

### 11.2 API Endpoints to Classes

| Endpoint | Primary Class | Operation |
|----------|--------------|-----------|
| `POST /api/assessments` | Assessment | create() |
| `GET /api/assessments/:id` | Assessment | findById() |
| `GET /api/assessments` | Assessment | findByUser() |
| `DELETE /api/assessments/:id` | Assessment | delete() |
| `POST /api/auth/login` | User, AuthSession | authenticate() |
| `GET /api/reports/:id` | Report | download() |
| `GET /api/models` | ModelVersion | listActive() |

---

## 12. Summary Statistics

- **Total Classes**: 10
- **Entity Classes**: 6 (User, Assessment, ModelVersion, Report, AuthSession, AuditLog)
- **Value Object Classes**: 4 (Feature, RiskScore, Explanation, CalibrationCurve)
- **Associations**: 9
- **Enumerations**: 4
- **Aggregations**: 2
- **Compositions**: 6
- **Simple Associations**: 1

---

**Document Prepared By:** Development Team  
**CASE Tool Used:** PlantUML (Text-based UML generator)  
**Review Status:** Ready for Stage 3 Submission  
**Next Steps:** Sequence diagrams for key scenarios, state machine diagrams

---

**End of Class Modeling Documentation**
