#!/bin/bash

set -e
echo "Running db migrations..."

docker run --rm \
-v $(pwd) /migrations:/migrations \
--network virtus_default
migrate/migrate \
-path=/migrations/ \
-database "postgres://virtus:virtus@postgres:5432/virtus?sslmode=disable" \
up

echo "Migrations completed successfully"
