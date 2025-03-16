#!/bin/sh
set -e

echo "Waiting for PostgreSQL to be ready..."
until pg_isready -h postgres -p 5432 -U user; do
  sleep 2
done

echo "Running migrations..."
migrate -database "$DATABASE_URL" -path /app/database/migrations up

echo "Starting bot..."
exec /app/bot