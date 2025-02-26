FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bot ./cmd/bot

FROM debian

RUN apt update && apt install -y libc6 ca-certificates postgresql-client && update-ca-certificates

WORKDIR /app

COPY --from=builder /app/bot /app/bot
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/database/migrations /app/database/migrations
COPY --from=builder /app/bot-entrypoint.sh /app/bot-entrypoint.sh

COPY ./goose /usr/local/bin/goose

RUN chmod +x /app/bot /usr/local/bin/goose /app/bot-entrypoint.sh

ENTRYPOINT ["/app/bot-entrypoint.sh"]