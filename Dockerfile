FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go build -o bot ./cmd/bot

FROM debian:bookworm

RUN apt update && apt install -y libc6
RUN apt update && apt install -y ca-certificates && update-ca-certificates

WORKDIR /app

COPY --from=builder /app/bot /app/bot
COPY --from=builder /app/database/migrations /app/database/migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose

RUN chmod +x /app/bot

CMD ["/app/bot"]