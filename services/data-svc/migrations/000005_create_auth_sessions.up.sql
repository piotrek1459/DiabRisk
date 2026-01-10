-- Create auth_sessions table (user session management)
CREATE TABLE auth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(128) NOT NULL UNIQUE, -- SHA-512 hash of session token
    ip_address INET, -- Client IP address
    user_agent TEXT, -- Browser user agent string
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    last_activity TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Sessions must have future expiration
    CONSTRAINT valid_expiration CHECK (expires_at > created_at)
);

-- Indexes for session lookups and cleanup
CREATE INDEX idx_auth_sessions_user ON auth_sessions(user_id) WHERE is_revoked = FALSE;
CREATE INDEX idx_auth_sessions_token ON auth_sessions(token_hash) WHERE is_revoked = FALSE;
CREATE INDEX idx_auth_sessions_expires ON auth_sessions(expires_at) WHERE is_revoked = FALSE;
CREATE INDEX idx_auth_sessions_last_activity ON auth_sessions(last_activity DESC);

COMMENT ON TABLE auth_sessions IS 'User authentication sessions with token management and revocation';
COMMENT ON COLUMN auth_sessions.token_hash IS 'SHA-512 hash of the session token sent to client (never store raw tokens)';
COMMENT ON COLUMN auth_sessions.is_revoked IS 'Manually revoked sessions (logout, security breach, etc.)';
COMMENT ON COLUMN auth_sessions.last_activity IS 'Last request timestamp using this session (for idle timeout)';
