# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-s -w -X main.Version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev') \
    -X main.Commit=$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown') \
    -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" \
    -o gurl ./cmd/gurl

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/gurl .

# Create non-root user
RUN addgroup -g 1000 gurl && \
    adduser -D -u 1000 -G gurl gurl && \
    chown -R gurl:gurl /app

USER gurl

ENTRYPOINT ["/app/gurl"]
CMD ["--help"]
