CREATE TABLE restaurants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    image_url TEXT,
	location GEOGRAPHY(Point, 4326)
);
