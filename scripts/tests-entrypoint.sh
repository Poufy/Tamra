#!/bin/sh

# First, wait for the database to be ready
/usr/local/bin/wait-for-it.sh ${TEST_DB_HOST}:${TEST_DB_PORT} -t 30  -- echo "Database is up"

# Now run the migrations
echo "Running migrations..."
migrate -path ./migrations -database "postgres://${TEST_DB_USER}:${TEST_DB_PW}@${TEST_DB_HOST}:${TEST_DB_PORT}/${TEST_DB_NAME}?sslmode=disable" up

# Run seeding
echo "Running seed scripts..."
PGPASSWORD=${TEST_DB_PW} psql -h ${TEST_DB_HOST} -U ${TEST_DB_USER} -d ${TEST_DB_NAME} -a -f ./seeds/seed.sql

# Run the tests
echo "Running tests..."
go test -v ./...