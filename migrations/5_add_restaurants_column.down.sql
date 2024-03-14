-- remove phone_number and location_description columns to the restaurants table
ALTER TABLE restaurants
DROP COLUMN phone_number,
DROP COLUMN location_description;