CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    code VARCHAR(255) UNIQUE,
    description TEXT,
    state VARCHAR(100),
    restaurant_id INT REFERENCES restaurants(id) ,
    user_id INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);