dependency:
	@echo ">> Downloading Dependencies"
	@go mod download

swag-init:
	@echo ">> Running swagger init"
	@swag init

run-api: dependency swag-init
	@echo ">> Running API Server"
	@go run main.go serve-http

run-consumer: dependency swag-init
	@echo ">> Running Kafka Consumer"
	@go run main.go run-consumer

migrate-up:
	@echo ">> Running Migrate Up"
	@migrate -path db/migrations -database "postgres://postgres:1235813@localhost:5433/vote?sslmode=disable" up

migrate-down:
	@echo ">> Running Migrate down"
	@migrate -path db/migrations -database "postgres://postgres:1235813@localhost:5433/vote?sslmode=disable" down

remock:
	#https://github.com/vektra/mockery
	@echo ">> Mock Repositories"
	@mockery --all --recursive --dir ./internal/domain/repository --output ./internal/domain/repository/mocks_repository --outpkg mocks_repository

	@echo ">> Mock Usecases"
	@mockery --all --dir ./internal/usecases --output ./internal/usecases/mocks_usecases --outpkg mocks_usecases

	@echo ">> Mock Interfaces"
	@mockery --all --recursive --dir ./internal/interfaces --output ./internal/interfaces/mocks_interfaces --outpkg mocks_interfaces

	@echo ">> Mock Infra"
	@mockery --all --recursive --dir ./internal/infrastructures --output ./internal/infrastructures/mocks_infrastructures --outpkg mocks_infrastructures