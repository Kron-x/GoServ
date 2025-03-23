# Compilating
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod main.go .

RUN go build -o main main.go

# running
FROM alpine:latest

WORKDIR /app

# Копируем скомпилированный бинарник из стадии builder
COPY --from=builder /app/main .

COPY config.json images/ server.log .

EXPOSE 5000

CMD ["./main"]