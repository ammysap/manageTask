# Stage 1: Build from repo root
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

WORKDIR /app/internal/services/taskmanager

ENV GOPROXY=https://goproxy.io,direct
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/taskmanager .

# Stage 2: Runtime
FROM alpine:3.20
RUN apk add --no-cache ca-certificates

WORKDIR /root
COPY --from=builder /app/taskmanager .

EXPOSE 8080
CMD ["./taskmanager"]