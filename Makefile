.PHONY: build run test clean get-token docker-build docker-run docker-compose-up docker-compose-down

# Переменные
BINARY_NAME=todo-server
DOCKER_IMAGE=todo-server
PORT=7540

# Локальная разработка
build:
	@echo "Building binary..."
	go build -o $(BINARY_NAME) .

run:
	@echo "Starting server..."
	go run main.go

test:
	@echo "Running tests..."
	cd tests && go test -v

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -f *.db
	rm -rf data/

get-token:
	@echo "Getting authentication token..."
	cd scripts && go run get_token.go

# Docker команды
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	@echo "Running Docker container..."
	docker run -p $(PORT):$(PORT) \
		-e TODO_PASSWORD=mysecretpassword \
		-e TODO_DBFILE=/app/data/scheduler.db \
		-v $(PWD)/data:/app/data \
		$(DOCKER_IMAGE)

docker-compose-up:
	@echo "Starting with Docker Compose..."
	docker-compose up -d

docker-compose-down:
	@echo "Stopping Docker Compose..."
	docker-compose down

docker-logs:
	docker-compose logs -f

# Полные workflow
dev: build run

docker-dev: docker-build docker-run

production: docker-build docker-compose-up

# Вспомогательные команды
setup-env:
	@echo "Setting environment variables..."
	export TODO_PORT=7540
	export TODO_DBFILE="scheduler.db"
	export TODO_PASSWORD="mysecretpassword"

help:
	@echo "Available commands:"
	@echo "  make build        - Build binary"
	@echo "  make run          - Run server locally"
	@echo "  make test         - Run tests"
	@echo "  make get-token    - Get authentication token"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo "  make dev          - Build and run locally"
	@echo "  make docker-dev   - Build and run with Docker"