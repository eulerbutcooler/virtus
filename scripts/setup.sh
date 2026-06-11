#!/bin/bash

set -e

# Always run from the repo root so .env loading works in child scripts
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR/.."

echo "Setting up Virtus local environment..."

echo "1/4 Starting Docker Compose..."
docker compose up -d

echo "2/4 Waiting for Postgres to be healthy..."
sleep 5

echo "3/4 Running migrations..."
bash scripts/migrate.sh

echo "4/4 Seeding the database..."
bash scripts/seed.sh

echo "Setup complete. Run 'go run ./cmd/api' to start the API on http://localhost:8080"
