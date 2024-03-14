-- Add phone_number and location_description columns to the restaurants table
ALTER TABLE restaurants
ADD COLUMN phone_number TEXT,
ADD COLUMN location_description TEXT NOT NULL;