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
	
migrate-up:
	@echo "Running migrations..."
	migrate -path $(MIGRATION_DIR) -database "postgresql://postgres:mysecretpassword@localhost:5432/tamra-postgis?sslmode=disable" up

migrate-down:
	@echo "Running migrations..."
	migrate -path $(MIGRATION_DIR) -database "postgresql://postgres:mysecretpassword@localhost:5432/tamra-postgis?sslmode=disable" down

swagger:
	@echo "Generating Swagger documentation..."
	rm -rf public/docs && swag init -d $(MAIN_DIR),$(HANDLERS_DIR) -g main.go --parseInternal --parseDependency -o public/docs

.PHONY: migrate-up migrate-down docker-up docker-down swagger
