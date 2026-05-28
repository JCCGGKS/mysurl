#!/usr/bin/env bash

set -eu

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-root}"
MYSQL_DATABASE="${MYSQL_DATABASE:-short}"

if ! command -v mysql >/dev/null 2>&1; then
  echo "mysql command not found"
  exit 1
fi

SQL_FILES="$(find "$SCRIPT_DIR" -maxdepth 1 -type f -name '*.sql' | sort)"

if [ -z "$SQL_FILES" ]; then
  echo "no sql files found in $SCRIPT_DIR"
  exit 1
fi

export MYSQL_PWD="$MYSQL_PASSWORD"

mysql \
  --host="$MYSQL_HOST" \
  --port="$MYSQL_PORT" \
  --user="$MYSQL_USER" \
  --execute="CREATE DATABASE IF NOT EXISTS \`$MYSQL_DATABASE\` CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;"

for sql_file in $SQL_FILES; do
  echo "executing $sql_file"
  mysql \
    --host="$MYSQL_HOST" \
    --port="$MYSQL_PORT" \
    --user="$MYSQL_USER" \
    "$MYSQL_DATABASE" < "$sql_file"
done

echo "initialized $MYSQL_DATABASE from sql files in $SCRIPT_DIR"
