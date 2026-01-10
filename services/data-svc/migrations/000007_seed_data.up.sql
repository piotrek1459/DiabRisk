-- Seed initial data for development and production

-- Insert default ML model version
INSERT INTO model_versions (
    version,
    trained_at,
    auroc,
    auprc,
    calibration_error,
    artifact_path,
    calibration_data,
    is_active
) VALUES (
    'v1.0.0',
    '2015-01-01',
    0.8562,
    0.4321,
    0.0234,
    '/models/v1.0.0/model.pkl',
    '{
        "calibration_curve": [[0.0, 0.0], [0.2, 0.18], [0.4, 0.39], [0.6, 0.62], [0.8, 0.81], [1.0, 1.0]],
        "calibration_method": "isotonic_regression",
        "thresholds": {
            "prediabetes": 0.35,
            "diabetes": 0.65
        },
        "performance_metrics": {
            "accuracy": 0.8234,
            "precision": 0.7512,
            "recall": 0.6891,
            "f1_score": 0.7189
        }
    }'::JSONB,
    TRUE
) ON CONFLICT (version) DO NOTHING;

-- Insert test admin user (for development only)
INSERT INTO users (
    email,
    role
) VALUES (
    'admin@diabrisk.local',
    'admin'
) ON CONFLICT (email) DO NOTHING;

-- Insert test regular user (for development only)
INSERT INTO users (
    email,
    role
) VALUES (
    'user@diabrisk.local',
    'registered'
) ON CONFLICT (email) DO NOTHING;

COMMENT ON TABLE model_versions IS 'NOTE: Default model v1.0.0 has basic calibration data for diabetes risk thresholds';
