FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache git
RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /doctor-service ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /doctor-service /doctor-service
COPY .env .env
COPY migrations /migrations
EXPOSE 8085
CMD ["/doctor-service"]
