.PHONY: build up down logs clean help

help:
	@echo "Available commands:"
	@echo "  make build    - Build the Docker images"
	@echo "  make up       - Start the services in detached mode"
	@echo "  make down     - Stop the services"
	@echo "  make logs     - Show service logs"
	@echo "  make clean    - Stop services and remove volumes (clean data)"
	@echo "  make dev      - Start services with logs (development)"

# Build the Docker images
build:
	docker-compose build

# Start all services in detached mode
up:
	docker-compose up -d

# Start services with logs (development)
dev:
	docker-compose up

# Stop all services
down:
	docker-compose down

# Show service logs
logs:
	docker-compose logs -f

# Stop services and remove volumes (clean all data)
clean:
	docker-compose down -v
	docker system prune -f

# Check service status
status:
	docker-compose ps

# Health check
health:
	curl -f http://localhost:8080/health || echo "Service is not healthy"

# Default command (build and start)
default: build up