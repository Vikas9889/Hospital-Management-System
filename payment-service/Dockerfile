# ---------- Stage 1: Build ----------
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install git (for go mod)
RUN apk add --no-cache git

# Copy go mod and sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /payment-service ./cmd

# ---------- Stage 2: Run ----------
FROM alpine:latest
RUN apk add --no-cache ca-certificates netcat-openbsd
WORKDIR /app

# Copy binary and wait-for script
COPY --from=builder /payment-service /payment-service
COPY wait-for.sh /wait-for.sh

RUN chmod +x /wait-for.sh

EXPOSE 8088

# Wait for Postgres (billing) then start
CMD ["/wait-for.sh", "postgres-billing:5432", "/payment-service"]
