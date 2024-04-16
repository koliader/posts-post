# Build stage
FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go
RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN chmod +x start.sh

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY internal/db/migration ./migration

# Set GIN_MODE to release
EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
