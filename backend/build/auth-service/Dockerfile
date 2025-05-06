FROM golang:1.20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-service ./cmd/auth-service

# Use a minimal alpine image
FROM alpine:3.17

# Set working directory
WORKDIR /app

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/auth-service .

# Expose port
EXPOSE 8001

# Run the application
CMD ["./auth-service"]
