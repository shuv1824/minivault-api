# Project configuration
MODEL_NAME ?= llama2

.PHONY: all build run stop restart logs pull-model run-bin build-bin

# Build Go API Docker image
build:
	docker-compose build

# Start services (Ollama + API)
run: run-model
	docker-compose up -d

# Pull and run the model into the Ollama container
run-model:
	docker-compose up -d ollama
	@echo "Waiting for Ollama to be ready..."
	@sleep 5
	docker exec ollama ollama run $(MODEL_NAME)

# Stop all containers
stop:
	docker-compose down

# Restart everything
restart: stop build run

# Show logs
logs:
	docker-compose logs -f

# Build Go CLI binary and run when ollama is running in docker and
run-cli: build-cli
	OLLAMA_URL=http://localhost:11434 MODE=CLI ./bin/minivault

build-cli:
	@go build -o bin/minivault .
