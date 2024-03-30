CREATE TABLE restaurants (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    logo_url TEXT,
    location GEOGRAPHY(Point, 4326) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
