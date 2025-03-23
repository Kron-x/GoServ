# Compilating
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod main.go .

ENV CGO_ENABLED = 0

RUN go build -o main .

# running
FROM alpine:latest

WORKDIR /app

# Копируем скомпилированный бинарник из стадии builder
COPY --from=builder /app/main .

COPY config.json images/ .

RUN mkdir -p ./logs/

EXPOSE 5000

CMD ["./main"]