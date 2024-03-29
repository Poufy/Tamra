MIGRATION_DIR = ./migrations/
DOCKER_COMPOSE_FILE = ./deployments/docker-compose.yml
MAIN_DIR = ./cmd/tamra
HANDLERS_DIR = ./internal/app/tamra/handlers

# Test database variables
TEST_DB_NAME = tamra-postgis-test
TEST_DB_USER = postgres
TEST_DB_PASSWORD = mysecretpassword
TEST_DB_PORT = 5432
TEST_DB_CONTAINER_NAME = tamra-postgis-test
DOCKER_COMPOSE_TEST_FILE = ./deployments/test.docker-compose.yml

# Build and run Docker Compose for local development
dev-db-up:
	@echo "Starting the database container..."
	docker run --name $(TEST_DB_CONTAINER_NAME) -e POSTGRES_USER=$(TEST_DB_USER) -e POSTGRES_PASSWORD=$(TEST_DB_PASSWORD) -e POSTGRES_DB=$(TEST_DB_NAME) -p $(TEST_DB_PORT):$(TEST_DB_PORT) postgis/postgis

dev-db-down:
	@echo "Stopping the database container..."
	docker stop $(TEST_DB_CONTAINER_NAME) && docker rm $(TEST_DB_CONTAINER_NAME)
	
migrate-dev-db-up:
	@echo "Running migrations..."
	migrate -path $(MIGRATION_DIR) -database "postgresql://$(TEST_DB_USER):$(TEST_DB_PASSWORD)@localhost:$(TEST_DB_PORT)/$(TEST_DB_NAME)?sslmode=disable" up

migrate-dev-db-down:
	@echo "Running migrations..."
	migrate -path $(MIGRATION_DIR) -database "postgresql://$(TEST_DB_USER):$(TEST_DB_PASSWORD)@localhost:$(TEST_DB_PORT)/$(TEST_DB_NAME)?sslmode=disable" down

seed-dev-db:
	@echo "Seeding the database..."
	PGPASSWORD=$(TEST_DB_PASSWORD) psql -h localhost -U $(TEST_DB_USER) -d $(TEST_DB_NAME) -f ./seeds/seed.sql

swagger:
	@echo "Generating Swagger documentation..."
	swag fmt && swag init -d $(MAIN_DIR),$(HANDLERS_DIR) -g main.go --parseInternal --parseDependency -o docs

start:
	@echo "Running the application locally..."
	go build -o ./bin/tamra ./cmd/tamra/  && ./bin/tamra -stage=local -port=8080 -db=postgres://$(TEST_DB_USER):$(TEST_DB_PASSWORD)@localhost:$(TEST_DB_PORT)/$(TEST_DB_NAME)?sslmode=disable

dev : swagger start

# Variable is used to set the database connection string for the test database
export TEST_DB_CONNECTION_STRING=postgres://$(TEST_DB_USER):$(TEST_DB_PASSWORD)@localhost:$(TEST_DB_PORT)/$(TEST_DB_NAME)?sslmode=disable
run-unit-tests:
	@echo "Running tests..."
	go test -v ./...

test: migrate-dev-db-up seed-dev-db run-unit-tests migrate-dev-db-down
# Run tests in Docker. This is the process that will be used in CI/CD on AWS CodeBuild
docker-test:
	@echo "Running tests..."
	docker-compose -f $(DOCKER_COMPOSE_TEST_FILE) up --no-deps --build db app 

.PHONY: dev-db-up dev-db-down migrate-dev-db-up migrate-dev-db-down swagger start dev run-tests docker-test