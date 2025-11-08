# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache git
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /notification-service ./cmd

# Stage 2: Run
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /notification-service /notification-service
COPY .env .env
EXPOSE 8084
ENTRYPOINT ["/notification-service"]
