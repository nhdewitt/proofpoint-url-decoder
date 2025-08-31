# Multi-stage build for smaller final image
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o proofpoint-url-decoder .

# Final state
FROM alpine:latest

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/proofpoint-url-decoder .

# Copy templates and static files
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Change ownership
RUN chown appuser:appgroup /app/proofpoint-url-decoder

# Switch to non-root user
USER appuser

# Expose port 8089 for server mode
EXPOSE 8089

# Default to server mode
CMD ["./proofpoint-url-decoder", "-s"]