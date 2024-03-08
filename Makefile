MIGRATION_DIR = ./migrations/
DOCKER_COMPOSE_FILE = ./deployments/docker-compose.yml
MAIN_DIR = ./cmd/tamra
HANDLERS_DIR = ./internal/app/tamra/handlers
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
	go build -o ./bin/tamra ./cmd/tamra/  && ./bin/tamra -port=8080 -db=postgres://postgres:mysecretpassword@localhost:5432/tamra-postgis?sslmode=disable -firebase-config=firebaseConfig.json

all : swagger local
.PHONY: migrate-up migrate-down docker-up docker-down swagger local docker-up-db docker-down-db
