# Build stage
FROM golang:1.21-alpine AS builder

# Install git for go mod download
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY web-ui/go.mod ./web-ui/
WORKDIR /app/web-ui

# Download dependencies (if any)
RUN go mod download

# Copy source code
WORKDIR /app
COPY . .

# Build the application
WORKDIR /app/web-ui
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o web-ui .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests (needed for AI services) and wget for health checks
RUN apk --no-cache add ca-certificates git wget

# Create app user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary from builder stage to web-ui subdirectory
COPY --from=builder /app/web-ui/web-ui ./web-ui/web-ui

# Copy all challenge and package data to root level (so web-ui can find ../challenge-*)
COPY --chown=appuser:appgroup challenge-* ./
COPY --chown=appuser:appgroup packages ./packages/
COPY --chown=appuser:appgroup badges ./badges/
COPY --chown=appuser:appgroup scripts ./scripts/
COPY --chown=appuser:appgroup docs ./docs/
COPY --chown=appuser:appgroup images ./images/
COPY --chown=appuser:appgroup *.md ./
COPY --chown=appuser:appgroup *.sh ./

# Make scripts executable (before changing to web-ui directory)
RUN chmod +x *.sh 2>/dev/null || true
RUN find /app -name "*.sh" -exec chmod +x {} \; 2>/dev/null || true

# Set working directory to web-ui for execution
WORKDIR /app/web-ui

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./web-ui"]
