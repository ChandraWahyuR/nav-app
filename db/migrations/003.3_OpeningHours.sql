CREATE TABLE IF NOT EXISTS opening_hours (
    id VARCHAR(50) PRIMARY KEY,
    place_id VARCHAR(255),
    day VARCHAR(255),
    open_time VARCHAR(50),
    close_time VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT  fk_opening_place FOREIGN KEY (place_id)
    REFERENCES tempat_pariwisata(place_id)
    ON DELETE CASCADE
);