#!/bin/bash

# Set the environment variables
export DB_HOST=db
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=mysecretpassword
export DB_NAME=tamra-postgis-test

# wait-for-db.sh
# Wait for the database to be ready before starting the application
echo "Waiting for database to be ready..."
while ! pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER; do
  sleep 7
done

echo "Database is ready. Running migrations..."

migrate -path=/app/migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up

echo "Migrations complete. Running seed data..."

PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f /app/seeds/seed.sql

echo "Seed data complete. Running the tests..."