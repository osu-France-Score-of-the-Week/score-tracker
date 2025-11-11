#!/bin/bash
set -euo pipefail

# This script runs during Postgres container initialization. It uses environment
# variables (from docker-compose/.env) to create the DB user and DB if they do not exist.

PSQL_CMD="psql -v ON_ERROR_STOP=1 --username \"$POSTGRES_USER\""

echo "Running DB init script with DB_NAME=${DB_NAME}, DB_USER=${DB_USER}"

# Create user if not exists
${PSQL_CMD} <<-EOSQL
DO
\$do\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = '${DB_USER}') THEN
        CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';
    END IF;
END
\$do\$;
EOSQL

# Create database if not exists and set owner, then grant privileges
${PSQL_CMD} <<-EOSQL
SELECT 'CREATE DATABASE ${DB_NAME} OWNER ${DB_USER}'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '${DB_NAME}') \gexec
GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};
EOSQL

echo "DB init script finished"
