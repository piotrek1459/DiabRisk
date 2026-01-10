-- Drop auth_sessions table
DROP INDEX IF EXISTS idx_auth_sessions_last_activity;
DROP INDEX IF EXISTS idx_auth_sessions_expires;
DROP INDEX IF EXISTS idx_auth_sessions_token;
DROP INDEX IF EXISTS idx_auth_sessions_user;
DROP TABLE IF EXISTS auth_sessions;
