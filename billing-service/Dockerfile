# ---------- Stage 1: Build ----------
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install git (for go mod)
RUN apk add --no-cache git

# Copy go mod and sum first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /billing-service ./cmd

# ---------- Stage 2: Run ----------
FROM alpine:latest
RUN apk add --no-cache ca-certificates netcat-openbsd
WORKDIR /app

# Copy binary and environment files
COPY --from=builder /billing-service /billing-service
COPY .env .env
COPY wait-for.sh /wait-for.sh

# Make wait-for script executable
RUN chmod +x /wait-for.sh

# Expose app port
EXPOSE 8083

# ✅ Correct CMD — wait for postgres before running service
CMD ["/wait-for.sh", "postgres-billing:5432", "/billing-service"]
