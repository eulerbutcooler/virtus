#!/bin/sh
set -e

echo "Seeding the database..."

CONTAINER_ID=$(docker compose ps -q postgres)

if [ -z "$CONTAINER_ID" ]; then
  echo "Error: Postgres container is not running. Start it with 'docker compose up -d' first."
  exit 1
fi

docker exec -i $CONTAINER_ID psql -U virtus -d virtus <<'SQL'
INSERT INTO pools (id, balance, total_in, total_out, currency, updated_at)
VALUES ('00000000-0000-0000-0000-000000000001', 0.00, 0.00, 0.00, 'USD', NOW())
ON CONFLICT (id) DO NOTHING;
SQL

echo "Database seeded successfully!"
