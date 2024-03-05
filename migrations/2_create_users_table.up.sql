CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    location GEOGRAPHY(Point, 4326) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    fcm_token TEXT UNIQUE NOT NULL,
    phone VARCHAR(255) UNIQUE,
    radius INT NOT NULL,
    last_order_received TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);