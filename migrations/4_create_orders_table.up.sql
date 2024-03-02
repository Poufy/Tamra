CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    code VARCHAR(255),
    description TEXT,
    state VARCHAR(100),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    restaurant_id INT REFERENCES restaurants(id),
    user_id INT REFERENCES users(id)
);