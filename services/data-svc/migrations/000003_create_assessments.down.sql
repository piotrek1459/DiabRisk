-- Drop assessments table
DROP INDEX IF EXISTS idx_assessments_features;
DROP INDEX IF EXISTS idx_assessments_model_version;
DROP INDEX IF EXISTS idx_assessments_risk_level;
DROP INDEX IF EXISTS idx_assessments_created;
DROP INDEX IF EXISTS idx_assessments_user;
DROP TABLE IF EXISTS assessments;
