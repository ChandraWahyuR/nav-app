CREATE TABLE IF NOT EXISTS otp (
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    email VARCHAR(255) NOT NULL,
    otp_number INTEGER NOT NULL,
    valid_until TIMESTAMPTZ NOT NULL,
    status Boolean,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

