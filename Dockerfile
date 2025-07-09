# Build stage
FROM golang:1.23.9 AS builder

WORKDIR /app

# Pass the GitHub token as a build argument
ARG GITHUB_TOKEN

# Configure Git to use the token for authentication
RUN echo "machine github.com login ${GITHUB_TOKEN}" > ~/.netrc

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Install swag for Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
RUN swag init

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/result-service main.go

# Final stage
FROM alpine:3.18

# Install ca-certificates for HTTPS support
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/result-service .

# Copy the Swagger docs directory
COPY --from=builder /app/docs ./docs

RUN mkdir -p /app/config

# Copy the migration files
COPY --from=builder /app/db /app/db

# Expose the default API port
EXPOSE 8903

# Set the entrypoint to the binary
ENTRYPOINT ["/app/result-service"]

# Default command to run the API server
CMD ["serve-http"]