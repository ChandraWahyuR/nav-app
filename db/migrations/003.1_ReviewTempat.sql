CREATE TABLE IF NOT EXISTS review_tempat (
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    place_id VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    time_reviewed Float NOT NULL,
    text TEXT NOT NULL,
    address VARCHAR(255),
    rating int, 
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,

    Constraint fk_review_place FOREIGN KEY (place_id)
    REFERENCES tempat_pariwisata(place_id)
    ON DELETE CASCADE
);