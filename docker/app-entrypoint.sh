#!/bin/sh
set -eu

DB_HOST="${DB_HOST:-postgres}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-hospital_middleware}"
DB_SSLMODE="${DB_SSLMODE:-disable}"
DB_WAIT_TIMEOUT="${DB_WAIT_TIMEOUT:-60}"

echo "Waiting for postgres at ${DB_HOST}:${DB_PORT}..."

i=0
until PGPASSWORD="${DB_PASSWORD}" pg_isready -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" >/dev/null 2>&1; do
  i=$((i + 1))
  if [ "${i}" -ge "${DB_WAIT_TIMEOUT}" ]; then
    echo "Postgres did not become ready in time"
    exit 1
  fi
  sleep 1
done

if [ "${RUN_MIGRATIONS:-true}" = "true" ]; then
  echo "Running migrations..."
  set +e
  migrate_output="$(migrate -path /app/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" up 2>&1)"
  migrate_status=$?
  set -e
  printf '%s\n' "${migrate_output}"
  if [ "${migrate_status}" -ne 0 ] && ! printf '%s' "${migrate_output}" | grep -q "no change"; then
    exit "${migrate_status}"
  fi
fi

exec /app/api
