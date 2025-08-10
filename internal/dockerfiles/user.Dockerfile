# Stage 1: Build everything from repo root
FROM golang:1.24 AS builder

WORKDIR /app

# Copy whole repo (build context is repo root)
COPY . .

# Change to the module directory for the service
WORKDIR /app/internal/services/user

# Download dependencies (this can now resolve local replace ../../database)
ENV GOPROXY=https://goproxy.io,direct
RUN go mod download

# Build the service binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/user .

# Stage 2: Runtime
FROM alpine:3.20
RUN apk add --no-cache ca-certificates

WORKDIR /root
COPY --from=builder /app/user .

EXPOSE 50051
CMD ["./user"]