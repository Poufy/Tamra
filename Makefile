MIGRATION_DIR = ./migrations/
DOCKER_COMPOSE_FILE = ./deployments/docker-compose.yml
MAIN_DIR = ./cmd/tamra
HANDLERS_DIR = ./internal/app/tamra/handlers

# Test database variables
TEST_DB_NAME = tamra-postgis-test
TEST_DB_USER = postgres
TEST_DB_PASSWORD = mysecretpassword
TEST_DB_PORT = 5433
TEST_DB_CONTAINER_NAME = tamra-postgis-test

# Build and run Docker Compose
docker-up:
	@echo "Building and starting Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

docker-up-db:
	@echo "Starting the database container..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up db 

docker-down-db:
	@echo "Stopping the database container..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down db
	
migrate-up:
	@echo "Running migrations..."
	migrate -path $(MIGRATION_DIR) -database "postgresql://postgres:mysecretpassword@localhost:5432/tamra-postgis?sslmode=disable" up

migrate-down:
	@echo "Running migrations..."
	migrate -path $(MIGRATION_DIR) -database "postgresql://postgres:mysecretpassword@localhost:5432/tamra-postgis?sslmode=disable" down

swagger:
	@echo "Generating Swagger documentation..."
	swag fmt && swag init -d $(MAIN_DIR),$(HANDLERS_DIR) -g main.go --parseInternal --parseDependency -o docs

local:
	@echo "Running the application locally..."
	go build -o ./bin/tamra ./cmd/tamra/  && ./bin/tamra -port=8080 -db=postgres://postgres:mysecretpassword@localhost:5432/tamra-postgis?sslmode=disable

all : swagger local

create-test-db:
	@echo "Creating test database..."
	docker run --name $(TEST_DB_CONTAINER_NAME) -e POSTGRES_DB=$(TEST_DB_NAME) -e POSTGRES_USER=$(TEST_DB_USER) -e POSTGRES_PASSWORD=$(TEST_DB_PASSWORD) -p $(TEST_DB_PORT):5432 -d postgis/postgis
	until docker exec $(TEST_DB_CONTAINER_NAME) pg_isready; do sleep 7; done

migrate-test-db:
	@echo "Running migrations for test database..."
	migrate -path $(MIGRATION_DIR) -database "postgresql://$(TEST_DB_USER):$(TEST_DB_PASSWORD)@localhost:$(TEST_DB_PORT)/$(TEST_DB_NAME)?sslmode=disable" up

migrate-test-db-down:
	@echo "Running migrations for test database..."
	migrate -path $(MIGRATION_DIR) -database "postgresql://$(TEST_DB_USER):$(TEST_DB_PASSWORD)@localhost:$(TEST_DB_PORT)/$(TEST_DB_NAME)?sslmode=disable" down

# We set PGPASSWORD as an environment variable to avoid the password prompt
seed-test-db:
	@echo "Seeding test database..."
	PGPASSWORD=$(TEST_DB_PASSWORD) psql -h localhost -p $(TEST_DB_PORT) -U $(TEST_DB_USER) -d $(TEST_DB_NAME) -a -f ./testdata/seed.sql

delete-test-db:
	@echo "Deleting test database..."
	docker stop $(TEST_DB_CONTAINER_NAME) && docker rm $(TEST_DB_CONTAINER_NAME)

run-tests:
	@echo "Running tests..."
	go test -v ./...

test: create-test-db migrate-test-db run-tests delete-test-db

.PHONY: migrate-up migrate-down docker-up docker-down swagger local docker-up-db docker-down-db create-test-db migrate-test-db seed-test-db run-tests test
