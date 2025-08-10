# Stage 1: Build everything from repo root (so local replace paths resolve)
FROM golang:1.24 AS builder

WORKDIR /app

# Copy whole repo (build context is repo root)
COPY . .

# Change to the module directory for the service
WORKDIR /app/internal/services/taskmanager

# Download dependencies (this can now resolve local replace ../../database)
RUN go mod download

# Build the service binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/taskmanager .

# Stage 2: Runtime
FROM alpine:3.20
RUN apk add --no-cache ca-certificates

WORKDIR /root
COPY --from=builder /app/taskmanager .

EXPOSE 8080
CMD ["./taskmanager"]
