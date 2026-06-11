#!/bin/sh
set -e

# Load DB credentials from .env if present
if [ -f .env ]; then
  export $(grep -v '^#' .env | grep -E '^DB_' | xargs)
fi

DB_USER="${DB_USER:-virtus}"
DB_PASSWORD="${DB_PASSWORD:-virtus}"
DB_NAME="${DB_NAME:-virtusdb}"
DB_PORT="${DB_PORT:-5432}"
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo "Running database migrations..."

if command -v migrate > /dev/null 2>&1; then
  migrate -path ./migrations -database "${DB_URL}" up
else
  echo "  'migrate' CLI not found, falling back to Docker..."
  docker run --rm \
    -v "$(pwd)/migrations:/migrations" \
    --network host \
    migrate/migrate \
    -path=/migrations/ \
    -database "${DB_URL}" up
fi

echo "Migrations completed successfully!"
