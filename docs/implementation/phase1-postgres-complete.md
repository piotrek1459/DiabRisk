# Phase 1 Implementation Complete: PostgreSQL Setup

## Summary

Successfully implemented PostgreSQL database infrastructure for DiabRisk application with complete schema, migrations, and deployment to k3d cluster.

## What Was Implemented

### 1. Database Migrations (7 files)
Created SQL migration files in `/services/data-svc/migrations/`:

- **000001_create_users**: User accounts with email-based authentication
- **000002_create_model_versions**: ML model versioning with performance tracking
- **000003_create_assessments**: Core assessment aggregate with 21-feature JSONB storage
- **000004_create_reports**: PDF/HTML report generation tracking
- **000005_create_auth_sessions**: Session management with token hashing
- **000006_create_audit_logs**: Comprehensive audit trail for security/compliance
- **000007_seed_data**: Initial seed data (1 model version, 2 test users)

Each migration includes both `.up.sql` and `.down.sql` for rollback capability.

### 2. Database Schema Highlights

**Users Table:**
- UUID primary key
- Email with regex validation
- Role-based access (anonymous, registered, clinician, admin)
- Indexed for fast lookups

**Model Versions Table:**
- Version string (e.g., "v1.0.0")
- Performance metrics (AUROC, AUPRC, calibration_error)
- JSONB calibration data (curves, thresholds)
- Active model flag for deployment management

**Assessments Table:**
- Links user + model version
- 21 health indicator features in JSONB
- Raw and calibrated risk scores
- Risk level classification
- GIN index on JSONB for efficient feature queries

**Auth Sessions Table:**
- SHA-512 token hashing (never stores raw tokens)
- IP tracking and user agent logging
- Expiration and revocation support
- Last activity tracking for idle timeout

**Audit Logs Table:**
- Comprehensive action tracking
- JSONB metadata for request details
- Support for anonymous actions (user_id NULL)
- Indexed by user, timestamp, action, resource

### 3. Kubernetes Deployment

**PostgreSQL StatefulSet** ([postgres.yaml](../deploy/k8s/postgres.yaml)):
- PostgreSQL 16-alpine image
- 10Gi PersistentVolumeClaim for durability
- Health checks (readiness + liveness probes)
- Resource limits (256Mi memory, 500m CPU)
- ClusterIP service on port 5432

**Data Service** ([data-svc.yaml](../deploy/k8s/data-svc.yaml)):
- Go-based migration runner and database manager
- Automatic migration execution on startup
- Database connection validation
- Schema verification with table existence checks
- Seed data confirmation

**Secrets and Configuration**:
- **postgres-secret**: Database credentials (username, password, connection URL)
- **oauth-secret**: Google OAuth client credentials (placeholder for Phase 2)
- **app-config**: ConfigMap with service URLs and session settings

### 4. Data Service Implementation

**Main Features** ([services/data-svc/main.go](../services/data-svc/main.go)):
- PostgreSQL connection pooling with pgx/v5
- golang-migrate/v4 integration for migrations
- Wait-for-database startup logic (30 retries with backoff)
- Schema verification on startup
- Seed data validation
- Comprehensive logging

**Build & Deploy**:
- Multi-stage Docker build (Dockerfile.data-svc)
- Alpine-based minimal image
- Migrations bundled into container
- Deployed to k3d cluster with service discovery

## Current State

### ✅ Deployed and Running

```
NAME                           READY   STATUS    RESTARTS   AGE
api-gateway-6f87dd9586-hxnmm   1/1     Running   0          58m
data-svc-596b68f8d5-wxtqn      1/1     Running   0          3m
frontend-5c767bcc8b-v5z4j      1/1     Running   0          56m
postgres-0                     1/1     Running   0          5m
```

### Database Contents

- **7 tables created**: users, model_versions, assessments, reports, auth_sessions, audit_logs, schema_migrations
- **2 test users**: admin@diabrisk.local (admin role), user@diabrisk.local (registered role)
- **1 model version**: v1.0.0 with AUROC=0.8562, calibration thresholds at 0.35 (prediabetes) and 0.65 (diabetes)

### Schema Features

- ✅ UUID primary keys for all entities
- ✅ JSONB columns for flexible data (features, calibration, metadata)
- ✅ GIN indexes on JSONB for fast queries
- ✅ Foreign key constraints with ON DELETE cascades
- ✅ CHECK constraints for data validation
- ✅ Unique constraints on email, version strings
- ✅ Timestamps with timezone (TIMESTAMPTZ)
- ✅ Comprehensive table and column comments

## Verification Commands

Test database connectivity:
```bash
kubectl exec postgres-0 -- psql -U diabrisk -d diabrisk -c "\dt"
```

View users:
```bash
kubectl exec postgres-0 -- psql -U diabrisk -d diabrisk -c "SELECT * FROM users;"
```

View model versions:
```bash
kubectl exec postgres-0 -- psql -U diabrisk -d diabrisk -c "SELECT version, auroc, auprc, is_active FROM model_versions;"
```

Check data-svc logs:
```bash
kubectl logs -l app=data-svc
```

## Next Steps (Phase 2: Authentication Service)

1. **Implement auth-svc microservice**:
   - Google OAuth2 integration
   - Session creation and validation
   - JWT token generation (optional)
   - Password-based auth (future)

2. **Update API Gateway**:
   - Add authentication middleware
   - Extract user from session/token
   - Pass user context to downstream services
   - Public vs protected routes

3. **Frontend Integration**:
   - Google Sign-In button
   - Session cookie management
   - User profile display
   - Logout functionality

4. **Testing**:
   - End-to-end OAuth flow
   - Session expiration handling
   - Role-based access control

## Technical Notes

### Migration Challenges & Solutions

**Issue**: Migration system left database in "dirty" state during development.
**Solution**: Manual cleanup of schema_migrations table, set version to 7 (all migrations complete).

**Issue**: Seed data SQL didn't match actual table schema.
**Solution**: Simplified schema in migrations (removed google_id, password_hash from users table - these will be added by auth-svc in Phase 2).

**Issue**: Container credentials mismatch between secret and environment variable.
**Solution**: Updated data-svc deployment to use `secretKeyRef` instead of hardcoded values.

### Architecture Decisions

1. **JSONB for Features**: Assessments store 21 health indicators as JSONB instead of 21 columns:
   - Flexibility for future feature additions
   - GIN indexing enables efficient feature-based queries
   - Simplifies migration when adding new indicators

2. **Minimal Users Table**: Current users table only has email + role:
   - Authentication fields (password, OAuth IDs) will be added by auth-svc
   - Keeps data-svc focused on data management
   - Follows microservice separation of concerns

3. **StatefulSet for PostgreSQL**: Using StatefulSet instead of Deployment:
   - Predictable pod name (postgres-0)
   - Stable network identity
   - Ordered deployment and scaling
   - Proper for stateful workloads

## Files Created

```
services/data-svc/
├── go.mod                                          # Go dependencies
├── main.go                                         # Migration runner service
└── migrations/
    ├── 000001_create_users.up.sql
    ├── 000001_create_users.down.sql
    ├── 000002_create_model_versions.up.sql
    ├── 000002_create_model_versions.down.sql
    ├── 000003_create_assessments.up.sql
    ├── 000003_create_assessments.down.sql
    ├── 000004_create_reports.up.sql
    ├── 000004_create_reports.down.sql
    ├── 000005_create_auth_sessions.up.sql
    ├── 000005_create_auth_sessions.down.sql
    ├── 000006_create_audit_logs.up.sql
    ├── 000006_create_audit_logs.down.sql
    ├── 000007_seed_data.up.sql
    └── 000007_seed_data.down.sql

deploy/k8s/
├── postgres.yaml                                   # PostgreSQL StatefulSet
├── data-svc.yaml                                   # Data service deployment
├── secrets.yaml                                    # Database & OAuth secrets
└── configmap.yaml                                  # Application configuration

Dockerfile.data-svc                                 # Multi-stage Go build
```

## Performance Considerations

- **Connection Pooling**: pgx/v5 provides efficient connection pooling
- **Indexes**: Strategic indexes on frequently queried columns
- **GIN Indexes**: JSONB columns have GIN indexes for feature queries
- **Resource Limits**: PostgreSQL limited to 256Mi RAM, 500m CPU (adjust for production)
- **PVC Storage**: 10Gi allocated (monitor usage in production)

## Security Notes

- ✅ Database credentials stored in Kubernetes Secret
- ✅ Session tokens stored as SHA-512 hashes (never plaintext)
- ✅ Email format validation with regex constraint
- ✅ Role-based access control via CHECK constraint
- ⚠️  Test users created without authentication (auth-svc will handle this)
- ⚠️  Currently using `sslmode=disable` (enable SSL in production)

---

**Status**: Phase 1 Complete ✅  
**Next**: Begin Phase 2 (Authentication Service)  
**Date**: 2026-01-10
