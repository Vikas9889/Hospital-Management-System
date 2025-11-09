# ---------- Stage 1: Build ----------
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install git (for go mod)
RUN apk add --no-cache git

# Copy go mod and sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source and run tidy (tidy needs source to find imports)
COPY . .
# Ensure module files are consistent inside the build environment
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /prescription-service ./cmd

# ---------- Stage 2: Run ----------
FROM alpine:latest
RUN apk add --no-cache ca-certificates netcat-openbsd
WORKDIR /app

# Copy binary and wait-for script
COPY --from=builder /prescription-service /prescription-service
COPY wait-for.sh /wait-for.sh

RUN chmod +x /wait-for.sh

EXPOSE 8087

# Wait for Postgres (appointments) then start
CMD ["/wait-for.sh", "postgres-appointments:5432", "/prescription-service"]
