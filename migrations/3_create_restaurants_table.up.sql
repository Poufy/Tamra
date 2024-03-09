CREATE TABLE restaurants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    logo_url TEXT,
	location GEOGRAPHY(Point, 4326) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fb_user_id VARCHAR(255) NOT NULL UNIQUE
);

CREATE INDEX restaurants_uid_index ON restaurants (fb_user_id);
