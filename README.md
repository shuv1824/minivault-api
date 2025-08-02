# MiniVault

A lightweight local REST API that allows you to use a local LLM to respond to user prompts — completely **offline**.

## Project Overview

MiniVault provides two interfaces for interacting with local LLMs:

1. **REST API Server** - Exposes a `POST /generate` endpoint for programmatic access
2. **CLI Tool** - Interactive command-line interface for direct conversation

All interactions are logged to disk for audit and analysis purposes.

## Features

- **REST API** with `POST /generate` endpoint
- **Interactive CLI** for real-time conversations
- **System Status** endpoint (`GET /status`) with performance metrics
- **Automatic Logging** - All prompts/responses saved to `logs/chat.log`
- **Docker Support** - Containerized deployment with Docker Compose
- **Configurable Models** - Support for different Ollama models
- **Performance Insights** - Memory usage, uptime, and system metrics

## Setup Instructions

### Prerequisites

- **Docker & Docker Compose** (recommended)
- **Go 1.23+** (for local development)
- **Make** (for convenience commands)

### Option 1: Docker Deployment (Recommended)

1. **Clone the repository:**

   ```bash
   git clone <repository-url>
   cd minivault
   ```

2. **Start the services:**

   ```bash
   make build
   make run
   ```

   This will:

   - Build the Go API container
   - Pull and start Ollama container
   - Download the default model (`llama2`)
   - Start the MiniVault API
   - Expose API on `http://localhost:8080`

### Option 2: Local Development

1. **Install Ollama locally:**

   ```bash
   # On macOS
   brew install ollama

   # On Linux
   curl -fsSL https://ollama.ai/install.sh | sh
   ```

2. **Start Ollama and pull a model:**

   ```bash
   ollama serve &
   ollama run llama2
   ```

3. **Or run the CLI:**

   ```bash
   make run-cli
   ```

## How to Run

### Using Make Commands (Docker)

```bash
# Build Docker image
make build

# Start everything
make run

# Stop services
make stop

# View logs
make logs

# Restart everything
make restart

# Run CLI tool (with Ollama in Docker)
make run-cli
```

### Manual Docker Commands

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f
```

### Environment Variables

- `MODE`: Set to `"API"` for server mode or `"CLI"` for interactive mode
- `OLLAMA_URL`: Ollama server URL (default: `http://localhost:11434`)
- `MODEL`: Model name to use (default: `llama2`)

## 📡 API Usage

### Generate Response

**Endpoint:** `POST /generate`

**Request:**

```json
{
  "prompt": "What is the capital of France?"
}
```

**Response:**

```json
{
  "response": "The capital of France is Paris. It is located in the north-central part of the country and serves as the political, economic, and cultural center of France."
}
```

**Example with curl:**

```bash
curl -X POST http://localhost:8080/generate \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Hello, how are you?"}'
```

### System Status

**Endpoint:** `GET /status`

**Response:**

```json
{
  "uptime": "2h30m15s",
  "memory_used_mb": 45.2,
  "num_goroutine": 8,
  "num_cpu": 8
}
```

## CLI Usage

Start the interactive CLI:

```bash
make run-cli
```

Example session:

```
Running CLI against model: llama2 (type 'exit' to quit)

Prompt: What is machine learning?
=> Machine learning is a subset of artificial intelligence that enables computers to learn and improve from experience without being explicitly programmed...

Prompt: exit
Exiting...
```

## Logging

All interactions are automatically logged to `logs/log.jsonl` file

## Project Structure

```
minivault/
├── main.go              # Application entry point
├── server/
│   ├── api.go          # REST API handlers and server logic
│   └── cli.go          # CLI interface implementation
├── logs/               # Request/response logs (auto-created)
├── bin/                # Compiled binaries
├── Dockerfile          # Container definition
├── docker-compose.yaml # Multi-service orchestration
├── Makefile           # Build and run commands
└── README.md          # This file
```

## 🔧 Development

### Building

```bash
# Build Docker image
make build

# Build local binary
go build -o bin/minivault .
```

### Testing the API

```bash
# Test generate endpoint
curl -X POST http://localhost:8080/generate \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Test prompt"}'

# Test status endpoint
curl http://localhost:8080/status
```

## What to Expect

### First Run

1. **Model Download**: The first run will download the specified LLM model (can take several minutes)
2. **Service Startup**: Ollama and MiniVault containers will start
3. **Ready State**: API will be available at `http://localhost:8080`

### Performance

- **Response Time**: Varies by model size and prompt complexity (typically 1-30 seconds)
- **Memory Usage**: Depends on model size (2-8GB+ for larger models)
- **Disk Usage**: Models are cached locally, logs accumulate over time

### Supported Models

Any model available in Ollama's library:

- `llama2` (default)
- `mistral`
- `codellama`
- `phi`
- `gemma`
- And many more...

## Docker Details

- **Base Image**: `golang:1.23-alpine`
- **Exposed Port**: `8080`
- **Volumes**:
  - `./logs:/app/logs` (log persistence)
  - `ollama:/root/.ollama` (model storage)

## Troubleshooting

### Common Issues

1. **"Failed to contact Ollama"**

   - Ensure Ollama is running and accessible
   - Check `OLLAMA_URL` environment variable
   - Verify model is downloaded: `docker exec ollama ollama list`

2. **Port conflicts**

   - Default ports: `8080` (API), `11434` (Ollama)
   - Modify `docker-compose.yaml` to use different ports

3. **Model not found**
   - Pull the model manually: `docker exec ollama ollama pull <model-name>`
   - Check available models: `docker exec ollama ollama list`

### Logs

```bash
# View all service logs
make logs

# View specific service logs
docker-compose logs -f minivault-api
docker-compose logs -f ollama
```

## Performance Insights

The `/status` endpoint provides real-time performance metrics:

- **Uptime**: How long the service has been running
- **Memory Usage**: Current memory consumption in MB
- **Goroutines**: Number of active goroutines (Go concurrency)
- **CPU Count**: Available CPU cores

This information helps monitor resource usage and optimize deployment.

## Future Enhancements

- [ ] Request rate limiting
- [ ] Model switching API endpoint
- [ ] Response streaming for real-time output
- [ ] Authentication and API keys
- [ ] Metrics and observability (Prometheus)
- [ ] Response caching
- [ ] Multi-model support

---

**Happy prompting!**
