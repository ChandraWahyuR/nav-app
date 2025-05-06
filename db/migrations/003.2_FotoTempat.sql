CREATE TABLE IF NOT EXISTS foto_tempat (
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    place_id VARCHAR(255) NOT NULL,
    photo_reference TEXT NOT NULL,
    width_px Float NOT NULL,
    height_px TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,

    Constraint fk_photo_place FOREIGN KEY (place_id)
    REFERENCES tempat_pariwisata(place_id)
    ON DELETE CASCADE
);