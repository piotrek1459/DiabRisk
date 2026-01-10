-- Create audit_logs table (security and compliance tracking)
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- NULL for anonymous actions
    action VARCHAR(100) NOT NULL, -- e.g., 'assessment_created', 'report_generated', 'login', 'logout'
    resource_type VARCHAR(50), -- e.g., 'assessment', 'report', 'session'
    resource_id UUID, -- ID of the affected resource
    ip_address INET,
    user_agent TEXT,
    request_metadata JSONB, -- Request details: endpoint, method, query params, etc.
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for audit queries
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id, timestamp DESC);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp DESC);
CREATE INDEX idx_audit_logs_action ON audit_logs(action, timestamp DESC);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);

-- GIN index for JSONB metadata searches
CREATE INDEX idx_audit_logs_metadata ON audit_logs USING gin(request_metadata);

COMMENT ON TABLE audit_logs IS 'Comprehensive audit trail for security, compliance, and debugging';
COMMENT ON COLUMN audit_logs.action IS 'Action performed: assessment_created, report_generated, login_success, login_failed, data_exported, etc.';
COMMENT ON COLUMN audit_logs.request_metadata IS 'Additional context: HTTP method, endpoint, query params, headers, response status, duration_ms, etc.';
