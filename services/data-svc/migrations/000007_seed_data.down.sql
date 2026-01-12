-- Remove seed data (careful in production!)

-- Remove test users
DELETE FROM users WHERE email IN ('admin@diabrisk.local', 'user@diabrisk.local');

-- Remove default model version
DELETE FROM model_versions WHERE version = 'v1.0.0';
