dependency:
	@echo ">> Downloading Dependencies"
	@go mod download

swag-init:
	@echo ">> Running swagger init"
	@swag init

run-api: dependency swag-init
	@echo ">> Running Result API Server"
	@go run main.go serve-http

run-consumer: dependency swag-init
	@echo ">> Running Result Kafka Consumer"
	@go run main.go run-consumer

# ClickHouse commands
clickhouse-up:
	@echo ">> Starting ClickHouse with Docker"
	@docker run -d \
		--name clickhouse-server \
		--ulimit nofile=262144:262144 \
		-p 8123:8123 \
		-p 9000:9000 \
		-p 9009:9009 \
		clickhouse/clickhouse-server

clickhouse-down:
	@echo ">> Stopping ClickHouse"
	@docker stop clickhouse-server || true
	@docker rm clickhouse-server || true

clickhouse-logs:
	@echo ">> Showing ClickHouse logs"
	@docker logs clickhouse-server -f

# Database migration commands
migrate-up:
	@echo ">> Running ClickHouse Migration Up"
	@clickhouse-client --host localhost --port 9000 --multiquery < db/migrations/000001_init_vote_results.up.sql

migrate-down:
	@echo ">> Running ClickHouse Migration Down"
	@clickhouse-client --host localhost --port 9000 --multiquery < db/migrations/000001_init_vote_results.down.sql

# Create database
create-db:
	@echo ">> Creating ClickHouse Database"
	@clickhouse-client --host localhost --port 9000 --query "CREATE DATABASE IF NOT EXISTS vote_results"

# Development setup
dev-setup: clickhouse-up
	@echo ">> Waiting for ClickHouse to start..."
	@sleep 10
	@make create-db
	@make migrate-up
	@echo ">> Development environment ready"

# Clean up development environment
dev-clean: clickhouse-down
	@echo ">> Development environment cleaned"

# Test commands
test:
	@echo ">> Running tests"
	@go test -v ./...

test-coverage:
	@echo ">> Running tests with coverage"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Build commands
build:
	@echo ">> Building Result Service"
	@go build -o bin/result-service main.go

build-docker:
	@echo ">> Building Docker image"
	@docker build -t result-service:latest .

# Mock generation
remock:
	@echo ">> Generating mocks"
	@mockery --all --recursive --dir ./internal/domain/repository --output ./internal/domain/repository/mocks_repository --outpkg mocks_repository
	@mockery --all --dir ./internal/usecases --output ./internal/usecases/mocks_usecases --outpkg mocks_usecases
	@mockery --all --recursive --dir ./internal/interfaces --output ./internal/interfaces/mocks_interfaces --outpkg mocks_interfaces

# Linting
lint:
	@echo ">> Running golangci-lint"
	@golangci-lint run

# Format code
fmt:
	@echo ">> Formatting code"
	@go fmt ./...

# Generate documentation
docs:
	@echo ">> Generating API documentation"
	@swag init
	@echo ">> Documentation available at http://localhost:8904/docs/"

# All-in-one development commands
dev-run-api: dev-setup run-api

dev-run-consumer: dev-setup run-consumer

# Production build
prod-build:
	@echo ">> Building for production"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/result-service main.go

# Clean build artifacts
clean:
	@echo ">> Cleaning build artifacts"
	@rm -rf bin/
	@rm -rf coverage.out coverage.html

.PHONY: dependency swag-init run-api run-consumer clickhouse-up clickhouse-down clickhouse-logs migrate-up migrate-down create-db dev-setup dev-clean test test-coverage build build-docker remock lint fmt docs dev-run-api dev-run-consumer prod-build clean