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
    poppler-utils \
    pandoc \
    && rm -rf /var/cache/apk/*

WORKDIR /app

# Copy only the necessary binary
COPY --from=builder /app/textindex .

# Create non-root user and set up directories
RUN addgroup -S jamtext && \
    adduser -S jamtext -G jamtext && \
    mkdir -p /data/indexes && \
    chown -R jamtext:jamtext /app /data/indexes && \
    chmod 755 /app/textindex

USER jamtext

# Configure environment
ENV INDEX_DIR=/data/indexes \
    CHUNK_SIZE=4096 \
    MAX_SHARD_SIZE=100000

VOLUME ["/data/indexes"]

# Health check
HEALTHCHECK --interval=30s --timeout=3s \
    CMD ./textindex -c stats || exit 1

ENTRYPOINT ["/app/textindex"]