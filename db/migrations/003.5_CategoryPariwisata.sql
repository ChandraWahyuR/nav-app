CREATE TABLE IF NOT EXISTS category_pariwisata (
    place_id VARCHAR(255),
    category_code VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    PRIMARY KEY (place_id, category_code),
    FOREIGN KEY (place_id) REFERENCES tempat_pariwisata(place_id) ON DELETE CASCADE,
    FOREIGN KEY (category_code) REFERENCES master_category(code) ON DELETE CASCADE
);
CREATE UNIQUE INDEX uniq_place_category ON category_pariwisata(place_id, category_code)
