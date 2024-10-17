# Dockerfile

# Stage 1: Build the Go application
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o exchange-rate-service cmd/main.go

# Stage 2: Run the application
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/exchange-rate-service .

# Copy configuration files
COPY --from=builder /app/configs/config.yaml ./configs/config.yaml

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./exchange-rate-service"]
