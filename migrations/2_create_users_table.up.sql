CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    location GEOGRAPHY(Point, 4326),
    is_active BOOLEAN,
    fcm_token TEXT UNIQUE,
    phone VARCHAR(255) UNIQUE,
    radius INT,
    last_order_received TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);