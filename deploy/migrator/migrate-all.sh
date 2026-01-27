#!/usr/bin/env bash
set -euo pipefail

function wait_for_db() {
  local host=$1
  local port=$2
  echo "Waiting for database $host:$port ..."
  until pg_isready -h "$host" -p "$port" >/dev/null 2>&1; do
    sleep 1
  done
}

wait_for_db "$POSTGRES_HOST" "$POSTGRES_PORT"

function run_migrations() {
  local db_name=$1
  local migrations_dir=$2
  local dsn="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${db_name}?sslmode=disable"
  echo "Applying migrations for ${db_name} from ${migrations_dir}"
  /usr/local/bin/migrate -path "$migrations_dir" -database "$dsn" up
}

run_migrations "$AUTH_DB" "/migrations/auth"
run_migrations "$GATEWAY_DB" "/migrations/gateway"
run_migrations "$REQUESTS_DB" "/migrations/requests"
run_migrations "$TECH_DB" "/migrations/tech"

echo "Migrations completed"