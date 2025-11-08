# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

# Copy the rest of the source code
COPY . .

RUN apk add --no-cache git
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service ./cmd

# Stage 2: Run
FROM alpine:latest
RUN apk add --no-cache ca-certificates netcat-openbsd
WORKDIR /app
COPY --from=builder /user-service /user-service
COPY .env .env
COPY wait-for.sh /wait-for.sh
RUN chmod +x /wait-for.sh
EXPOSE 8081
ENTRYPOINT ["/wait-for.sh", "/user-service"]
