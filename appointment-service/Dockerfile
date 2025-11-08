# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache git
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /appointment-service ./cmd

# Stage 2: Run
FROM alpine:latest
RUN apk add --no-cache ca-certificates netcat-openbsd
WORKDIR /app
COPY --from=builder /appointment-service /appointment-service
COPY .env .env
COPY wait-for.sh /wait-for.sh
RUN chmod +x /wait-for.sh
EXPOSE 8082
ENTRYPOINT ["/wait-for.sh", "/appointment-service"]
