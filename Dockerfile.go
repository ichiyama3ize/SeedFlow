# Dockerfile for Go main application
FROM golang:1.22-alpine AS builder

# Install dependencies including build tools
RUN apk add --no-cache git ca-certificates sqlite build-base

# Set working directory
WORKDIR /app

# Copy all source code
COPY go-app/ ./

# Tidy and download dependencies
RUN go mod tidy
RUN go mod download

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o knowledge-app .

# Final stage
FROM alpine:latest

# Install SQLite, CA certificates, and wget for health checks
RUN apk --no-cache add ca-certificates sqlite wget

# Create app directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/knowledge-app .

# Copy static assets
COPY --from=builder /app/static ./static/
COPY --from=builder /app/templates ./templates/

# Create necessary directories
RUN mkdir -p data logs exports backups config

# Set permissions
RUN chmod +x knowledge-app

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the application
CMD ["./knowledge-app"]