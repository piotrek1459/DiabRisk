-- Drop model_versions table
DROP INDEX IF EXISTS idx_model_versions_version;
DROP INDEX IF EXISTS idx_model_versions_active;
DROP TABLE IF EXISTS model_versions;
