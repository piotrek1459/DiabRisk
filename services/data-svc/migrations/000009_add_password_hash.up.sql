-- Add password hash field to users table for local authentication
ALTER TABLE users
    ADD COLUMN password_hash VARCHAR(255);

-- Create index on email for login lookups (if not exists)
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

COMMENT ON COLUMN users.password_hash IS 'Hashed password for local authentication';