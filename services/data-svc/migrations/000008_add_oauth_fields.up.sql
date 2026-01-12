-- Add OAuth fields to users table
ALTER TABLE users 
    ADD COLUMN google_id VARCHAR(255) UNIQUE,
    ADD COLUMN picture_url TEXT,
    ADD COLUMN full_name VARCHAR(255);

-- Create index on google_id for OAuth lookups
CREATE INDEX idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;

COMMENT ON COLUMN users.google_id IS 'Google OAuth unique identifier';
COMMENT ON COLUMN users.picture_url IS 'User profile picture URL from Google';
COMMENT ON COLUMN users.full_name IS 'User full name from Google profile';
