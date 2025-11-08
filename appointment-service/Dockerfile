FROM golang:1.22-alpine
WORKDIR /app
COPY . .
RUN go mod init temp && go mod tidy
RUN go build -o /main ./cmd
CMD ["/main"]
