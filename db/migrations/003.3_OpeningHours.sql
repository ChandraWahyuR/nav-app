CREATE TABLE IF NOT EXISTS foto_tempat (
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    place_id VARCHAR(255) NOT NULL,
    day VARCHAR(255) NOT NULL,
    open_time VARCHAR(50) NOT NULL,
    close_time VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,

    Constraint fk_opening_place FOREIGN KEY (place_id)
    REFERENCES tempat_pariwisata(place_id)
    ON DELETE CASCADE
);