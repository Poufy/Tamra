CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    location GEOGRAPHY(Point, 4326),
    is_active BOOLEAN,
    fcm_token TEXT,
    phone_number VARCHAR(255),
    search_radius INT,
    last_order_received_time TIMESTAMP
);