CREATE TABLE IF NOT EXISTS master_category (
    code VARCHAR(100) PRIMARY KEY, -- e.g. 'tourist_attraction', 'amusement_park'
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);