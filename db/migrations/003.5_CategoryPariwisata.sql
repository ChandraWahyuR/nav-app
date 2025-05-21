CREATE TABLE IF NOT EXISTS category_pariwisata (
    place_id VARCHAR(255),
    category_code VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (place_id, category_code),
    FOREIGN KEY (place_id) REFERENCES tempat_pariwisata(place_id) ON DELETE CASCADE,
    FOREIGN KEY (category_code) REFERENCES master_category(code) ON DELETE CASCADE
);
