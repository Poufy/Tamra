-- Seed data for users table
INSERT INTO users (id, location, is_active, fcm_token, phone, radius, last_order_received)
VALUES
    ('user1', ST_SetSRID(ST_MakePoint(-77.0364, 38.8951), 4326), true, 'token1', '+09055234232', 10, CURRENT_TIMESTAMP),
    ('user2', ST_SetSRID(ST_MakePoint(-77.0364, 38.8951), 4326), true, 'token2', '+09055234234', 20, CURRENT_TIMESTAMP);
    

-- Seed data for restaurants table
INSERT INTO restaurants (id, name, logo_url, location, phone_number, location_description)
VALUES
    ('restaurant1', 'restaurant1', 'https://www.google.com', ST_SetSRID(ST_MakePoint(-75.0364, 38.8951), 4326), '+09055234232', 'restaurant1 location'),
    ('restaurant2', 'restaurant2', 'https://www.google.com', ST_SetSRID(ST_MakePoint(-74.0364, 38.8951), 4326), '+09055234234', 'restaurant2 location');

-- Seed data for orders table
INSERT INTO orders (id, code, description, state, user_id, restaurant_id)
VALUES
    (1, 'XYSDFS', 'order1 description', 'PENDING', 'user1', 'restaurant1'),
    (2, 'FSHWRY', 'order2 description', 'PENDING', 'user2', 'restaurant2');
