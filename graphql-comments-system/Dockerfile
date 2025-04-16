FROM golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./server.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/internal /app/internal
COPY .env .

EXPOSE 8080

CMD ["./server"]
