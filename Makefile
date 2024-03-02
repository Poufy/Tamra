MIGRATION_DIR = ./migrations/
DOCKER_COMPOSE_FILE = ./deployments/docker-compose.yml

# Build and run Docker Compose
docker-up:
	@echo "Building and starting Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
	
# Run migrations
migrate:
	@echo "Running migrations..."
	migrate -path $(MIGRATION_DIR) -database "postgresql://postgres:mysecretpassword@localhost:5432/tamra-postgis?sslmode=disable" up

.PHONY: migrate docker-up docker-down
