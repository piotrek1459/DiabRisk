-- Create reports table (generated PDF/HTML reports)
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    format VARCHAR(10) NOT NULL CHECK (format IN ('pdf', 'html')),
    content_url TEXT NOT NULL, -- URL to stored report (S3, local file, etc.)
    metadata JSONB, -- Report generation metadata (template version, charts included, etc.)
    generated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- One report per format per assessment
    UNIQUE(assessment_id, format)
);

-- Indexes for lookups
CREATE INDEX idx_reports_assessment ON reports(assessment_id);
CREATE INDEX idx_reports_generated ON reports(generated_at DESC);

COMMENT ON TABLE reports IS 'Generated assessment reports (PDF/HTML) with download links';
COMMENT ON COLUMN reports.content_url IS 'URL to stored report file (e.g., s3://bucket/reports/uuid.pdf or file:///path/to/report.pdf)';
COMMENT ON COLUMN reports.metadata IS 'Report generation details: template_version, chart_types, language, page_count, etc.';
