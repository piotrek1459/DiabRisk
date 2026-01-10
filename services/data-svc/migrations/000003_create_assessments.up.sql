-- Create assessments table (core aggregate root)
CREATE TABLE assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    model_version_id UUID NOT NULL REFERENCES model_versions(id),
    features JSONB NOT NULL, -- Stores the 21 health indicator features
    raw_score DECIMAL(5,4) NOT NULL, -- Model output probability [0,1]
    calibrated_score DECIMAL(5,4), -- Calibrated probability after adjustments
    risk_level VARCHAR(20) NOT NULL CHECK (risk_level IN ('no_diabetes', 'prediabetes', 'diabetes')),
    confidence_level VARCHAR(20) CHECK (confidence_level IN ('high', 'medium', 'low')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure features contains all required 21 fields
    CONSTRAINT features_valid CHECK (
        features ?& ARRAY[
            'HighBP', 'HighChol', 'CholCheck', 'BMI', 'Smoker', 
            'Stroke', 'HeartDiseaseorAttack', 'PhysActivity', 'Fruits', 
            'Veggies', 'HvyAlcoholConsump', 'AnyHealthcare', 'NoDocbcCost', 
            'GenHlth', 'MentHlth', 'PhysHlth', 'DiffWalk', 'Sex', 
            'Age', 'Education', 'Income'
        ]
    )
);

-- Indexes for common queries
CREATE INDEX idx_assessments_user ON assessments(user_id, created_at DESC);
CREATE INDEX idx_assessments_created ON assessments(created_at DESC);
CREATE INDEX idx_assessments_risk_level ON assessments(risk_level);
CREATE INDEX idx_assessments_model_version ON assessments(model_version_id);

-- GIN index for JSONB features (enables efficient feature-based queries)
CREATE INDEX idx_assessments_features ON assessments USING gin(features);

COMMENT ON TABLE assessments IS 'Core aggregate: diabetes risk assessments with ML predictions';
COMMENT ON COLUMN assessments.features IS '21 health indicator features as JSONB: HighBP, HighChol, CholCheck, BMI, Smoker, Stroke, HeartDiseaseorAttack, PhysActivity, Fruits, Veggies, HvyAlcoholConsump, AnyHealthcare, NoDocbcCost, GenHlth, MentHlth, PhysHlth, DiffWalk, Sex, Age, Education, Income';
COMMENT ON COLUMN assessments.raw_score IS 'Raw ML model output probability before calibration [0,1]';
COMMENT ON COLUMN assessments.calibrated_score IS 'Probability after applying calibration curve adjustments';
COMMENT ON COLUMN assessments.confidence_level IS 'Model confidence: high (>0.8), medium (0.5-0.8), low (<0.5)';
