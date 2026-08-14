#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

source .env

MIGRATE="$PROJECT_DIR/bin/migrate"
MIGRATIONS_DIR="file://$PROJECT_DIR/src/database/migrations"
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

echo "==> Rolling back migrations..."

if [ ! -x "$MIGRATE" ]; then
    echo "Error: migrate binary not found at $MIGRATE"
    echo "Run: make install-tools"
    exit 1
fi

ACTION="${1:-last}"

if [ "$ACTION" = "all" ]; then
    "$MIGRATE" \
        -source "$MIGRATIONS_DIR" \
        -database "$DATABASE_URL" \
        down
else
    "$MIGRATE" \
        -source "$MIGRATIONS_DIR" \
        -database "$DATABASE_URL" \
        down 1
fi

echo "==> Rollback completed successfully."
