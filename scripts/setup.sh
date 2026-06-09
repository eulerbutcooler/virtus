#!/bin/bash

set -e

echo "Setting up Virtus local environment..."

echo "1/4 Starting Docker Compose..."
docker compose up -d

echo "2/4 Waiting for Postgres to be healthy..."
sleep 5

echo "3/4 Running migrations..."
./scripts/migrate.sh

echo "4/4 Seeding the database..."
./scripts/seed.sh

echo "Setup complete. The API is running on http://localhost:8080"
