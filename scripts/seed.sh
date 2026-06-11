#!/bin/sh
set -e

# Load DB credentials from .env if present
if [ -f .env ]; then
  export $(grep -v '^#' .env | grep -E '^DB_' | xargs)
fi

DB_USER="${DB_USER:-virtus}"
DB_NAME="${DB_NAME:-virtusdb}"

echo "Seeding the database..."

CONTAINER_ID=$(docker compose ps -q postgres)

if [ -z "$CONTAINER_ID" ]; then
  echo "Error: Postgres container is not running. Start it with 'docker compose up -d' first."
  exit 1
fi

docker exec -i "$CONTAINER_ID" psql -U "$DB_USER" -d "$DB_NAME" <<'SQL'
-- Seed the global pool
INSERT INTO pools (id, balance, total_in, total_out, currency, updated_at)
VALUES ('00000000-0000-0000-0000-000000000001', 0.00, 0.00, 0.00, 'USD', NOW())
ON CONFLICT (id) DO NOTHING;

-- Seed a dev admin account
-- Credentials: admin@virtus.dev / admin123
INSERT INTO users (id, email, name, password_hash, role, verified, joined_at, updated_at)
VALUES (
  '00000000-0000-0000-0000-000000000002',
  'admin@virtus.dev',
  'Admin',
  '$2a$10$7WURZP3TFKRFY9Pfw1.DXuiibIJXLtpbcz2QOGetLEqvamAZ/GT82',
  'admin',
  true,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;
SQL

echo "Database seeded successfully!"
