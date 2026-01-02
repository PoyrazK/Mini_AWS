# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
# CGO_ENABLED=0 for static binary
RUN CGO_ENABLED=0 go build -o /bin/api cmd/api/main.go

# Production stage
FROM alpine:3.19

# Install runtime dependencies (ca-certificates for HTTPS, docker-cli optional if needed inside)
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /bin/api /app/api
COPY --from=builder /app/internal/repositories/postgres/migrations /app/migrations
# Copy docs for swagger if needed
COPY --from=builder /app/docs /app/docs

# Create data directory
RUN mkdir -p /app/thecloud-data

# Expose API port
EXPOSE 8080

CMD ["/app/api"]
