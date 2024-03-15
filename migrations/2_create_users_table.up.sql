CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    location GEOGRAPHY(Point, 4326) NOT NULL,
    is_active BOOLEAN NOT NULL,
    fcm_token TEXT UNIQUE NOT NULL,
    phone VARCHAR(255) UNIQUE,
    radius INT NOT NULL,
    last_order_received TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX users_last_order_received_index ON users (last_order_received);