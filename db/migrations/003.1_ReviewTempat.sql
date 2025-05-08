CREATE TABLE IF NOT EXISTS review_tempat (
    id VARCHAR(50) PRIMARY KEY,
    place_id VARCHAR(255),
    users_id VARCHAR(50) NULL,
    author VARCHAR(255),
    review_created VARCHAR(50),
    text TEXT,
    rating int, 
    isfrom_google Boolean,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT  fk_review_place FOREIGN KEY (place_id)
    REFERENCES tempat_pariwisata(place_id)
    ON DELETE CASCADE,

    CONSTRAINT  fk_review_user FOREIGN KEY (users_id)
    REFERENCES users(id) 
    ON DELETE SET NULL
);