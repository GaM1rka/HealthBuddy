FROM golang:1.24-alpine AS builder

# Install Git (if needed), ca-certificates for go:generate or TLS
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy sources and build
COPY . .
# Build a static binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o profile_service ./cmd/app/main.go


### Final stage
FROM alpine:3.18

# Install CA certificates for HTTPS
RUN apk add --no-cache ca-certificates && update-ca-certificates

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /home/appuser

# Copy binary from builder
COPY --from=builder /app/profile_service ./profile_service


# Run the service
ENTRYPOINT ["./profile_service"]