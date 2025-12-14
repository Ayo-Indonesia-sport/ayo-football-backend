# Build stage
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates (for fetching dependencies)
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/main ./cmd/api/main.go

# Final stage
FROM alpine:3.19

# Install ca-certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy .env.example as reference
COPY --from=builder /app/.env.example .env.example

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
