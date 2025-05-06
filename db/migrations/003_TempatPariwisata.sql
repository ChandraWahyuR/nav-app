CREATE TABLE IF NOT EXISTS tempat_pariwisata (
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    place_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    latitude Float NOT NULL,
    longtitude Float NOT NULL,
    address VARCHAR(255),
    icon VARCHAR(255), 
    business_status VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);