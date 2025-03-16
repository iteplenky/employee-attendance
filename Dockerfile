FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bot ./cmd/bot

FROM debian

RUN apt update && apt install -y libc6 ca-certificates postgresql-client curl && update-ca-certificates

RUN curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/migrate \
    && chmod +x /usr/local/bin/migrate

WORKDIR /app

COPY --from=builder /app/bot /app/bot
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/database/migrations /app/database/migrations
COPY --from=builder /app/bot-entrypoint.sh /app/bot-entrypoint.sh

RUN chmod +x /app/bot /app/bot-entrypoint.sh

ENTRYPOINT ["/app/bot-entrypoint.sh"]