-- Create model_versions table
CREATE TABLE IF NOT EXISTS model_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version VARCHAR(20) UNIQUE NOT NULL,
    trained_at TIMESTAMPTZ NOT NULL,
    auroc REAL NOT NULL CHECK (auroc BETWEEN 0.5 AND 1.0),
    auprc REAL NOT NULL,
    calibration_error REAL NOT NULL,
    artifact_path TEXT NOT NULL,
    calibration_data JSONB,
    is_active BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_model_versions_active ON model_versions(is_active) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_model_versions_version ON model_versions(version);

-- Add comments
COMMENT ON TABLE model_versions IS 'ML model versions with performance metrics';
COMMENT ON COLUMN model_versions.auroc IS 'Area Under ROC Curve (0.5-1.0)';
COMMENT ON COLUMN model_versions.calibration_data IS 'Pre-computed calibration bins as JSONB';
