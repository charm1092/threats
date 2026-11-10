CREATE TABLE IF NOT EXISTS threats (
    md5 VARCHAR(32) PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    severity TEXT NOT NULL CHECK (
    severity IN ('safe', 'low', 'medium', 'high')
    )
);