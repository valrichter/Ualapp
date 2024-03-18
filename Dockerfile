# Build stage
FROM golang:1.22.1-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY db/migrations ./db/migrations
RUN ls -la

EXPOSE 8080
CMD ["/app/main"]