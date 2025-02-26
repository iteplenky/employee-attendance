#!/bin/sh
set -e

echo "Waiting for PostgreSQL to be ready..."
until pg_isready -h postgres -p 5432 -U user; do
  sleep 2
done

echo "Running migrations..."
goose -dir /app/database/migrations postgres "$DATABASE_URL" up

echo "Starting bot..."
exec /app/bot