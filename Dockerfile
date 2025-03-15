# Stage 1: Dependencies
FROM golang:1.24.1-alpine AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Stage 2: Builder
FROM deps AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o textindex cmd/main.go

# Stage 3: Security scan (optional)
FROM aquasec/trivy:latest AS security-scan
COPY --from=builder /app /app
RUN trivy filesystem --no-progress --exit-code 1 --severity HIGH,CRITICAL /app

# Stage 4: Final minimal image
FROM alpine:3.19 AS production
RUN apk add --no-cache \
    poppler-utils=~23.11 \
    pandoc=~3.1 \
    && rm -rf /var/cache/apk/*

WORKDIR /app

# Copy only the necessary binary
COPY --from=builder /app/textindex .

# Create non-root user
RUN addgroup -S jamtext && \
    adduser -S jamtext -G jamtext && \
    chown -R jamtext:jamtext /app

# Create and set permissions for index directory
RUN mkdir -p /data/indexes && \
    chown -R jamtext:jamtext /data/indexes

# Switch to non-root user
USER jamtext

# Configure environment
ENV INDEX_DIR=/data/indexes \
    CHUNK_SIZE=4096 \
    MAX_SHARD_SIZE=100000

# Expose volume for persistent storage
VOLUME ["/data/indexes"]

# Health check
HEALTHCHECK --interval=30s --timeout=3s \
    CMD ./textindex -c stats || exit 1

ENTRYPOINT ["/app/textindex"]