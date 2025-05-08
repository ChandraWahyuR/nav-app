CREATE TABLE IF NOT EXISTS foto_tempat (
    id VARCHAR(50) PRIMARY KEY,
    place_id VARCHAR(255),
    review_id VARCHAR(50),
    users_id VARCHAR(50),
    photo_reference TEXT,
    width_px int,
    height_px int,
    isfrom_google Boolean,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT  fk_photo_place FOREIGN KEY (place_id)
    REFERENCES tempat_pariwisata(place_id)
    ON DELETE CASCADE,

    CONSTRAINT  fk_photo_review FOREIGN KEY (review_id)
    REFERENCES review_tempat(id)
    ON DELETE SET NULL,

    CONSTRAINT  fk_photo_user FOREIGN KEY (users_id)
    REFERENCES users(id)
    ON DELETE SET NULL
);