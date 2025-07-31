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
	@echo "Pulling model: $(MODEL_NAME)"
	docker exec ollama ollama run $(MODEL_NAME)

# Stop all containers
stop:
	docker-compose down

# Restart everything
restart: stop build run

# Show logs
logs:
	docker-compose logs -f

# Build Go API binary and run without docker and
# with locally installed ollama service
run-bin: build-bin
	@./bin/minivault

build-bin:
	OLLAMA_URL=http://localhost:11434 go build -o bin/minivault .
