#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

source .env

MIGRATE="$PROJECT_DIR/bin/migrate"
MIGRATIONS_DIR="file://$PROJECT_DIR/src/database/migrations"
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

echo "==> Running migrations..."
echo "Migrations dir: $MIGRATIONS_DIR"
echo "DB: ${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME}"

if [ ! -x "$MIGRATE" ]; then
    echo "Error: migrate binary not found at $MIGRATE"
    echo "Run: make install-tools"
    exit 1
fi

if [ -z "$1" ]; then
    "$MIGRATE" \
        -source "$MIGRATIONS_DIR" \
        -database "$DATABASE_URL" \
        up
else
    "$MIGRATE" \
        -source "$MIGRATIONS_DIR" \
        -database "$DATABASE_URL" \
        "$@"
fi

echo "==> Migrations completed successfully."
