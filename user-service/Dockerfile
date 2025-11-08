# Stage 1: Build
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache git
RUN go mod download
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
