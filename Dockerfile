# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod ./

# Remove any existing go.sum and download dependencies fresh
RUN rm -f go.sum && go mod download

# Copy source code
COPY . .

# Fix dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o /docger ./cmd/main.go

# Build worker
RUN CGO_ENABLED=1 GOOS=linux go build -o /worker ./cmd/worker/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /docger .
COPY --from=builder /worker .

# Copy database migration script
COPY database/migration.sql ./database/migration.sql

# Expose port
EXPOSE 8080

# Run the application
CMD ["./docger"]
