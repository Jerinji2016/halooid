# Halooid Backend

This directory contains the backend services for the Halooid platform.

## Directory Structure

- `cmd/`: Application entry points
- `internal/`: Private application code
- `pkg/`: Reusable packages
- `api/`: API definitions
- `scripts/`: Build and deployment scripts
- `configs/`: Configuration files

## Getting Started

### Prerequisites

- Go 1.20 or later
- PostgreSQL 14 or later
- Redis 6 or later
- Docker and Docker Compose (for local development)

### Setup

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Start the development environment:
   ```bash
   cd ..
   docker-compose up -d
   ```

3. Run the API Gateway:
   ```bash
   go run cmd/api-gateway/main.go
   ```

## Development

### Adding a New Service

1. Create a new directory in `cmd/` for the service entry point
2. Create appropriate packages in `internal/` for the service implementation
3. Add the service to the Docker Compose file if needed

### Testing

Run tests with:
```bash
go test ./...
```

## API Documentation

API documentation will be available at `/api/docs` when the API Gateway is running.
