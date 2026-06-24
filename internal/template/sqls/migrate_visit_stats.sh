#!/usr/bin/env bash

set -euo pipefail

MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-root}"
MYSQL_DATABASE="${MYSQL_DATABASE:-short}"
BATCH_SIZE="${BATCH_SIZE:-10}"
START_ID="${START_ID:-0}"
END_ID="${END_ID:-0}"
ONLY_POSITIVE_VISIT="${ONLY_POSITIVE_VISIT:-1}"

if ! command -v mysql >/dev/null 2>&1; then
  echo "mysql command not found"
  exit 1
fi

if ! [[ "$BATCH_SIZE" =~ ^[0-9]+$ ]] || [ "$BATCH_SIZE" -le 0 ]; then
  echo "BATCH_SIZE must be a positive integer"
  exit 1
fi

if ! [[ "$START_ID" =~ ^[0-9]+$ ]]; then
  echo "START_ID must be a non-negative integer"
  exit 1
fi

if ! [[ "$END_ID" =~ ^[0-9]+$ ]]; then
  echo "END_ID must be a non-negative integer"
  exit 1
fi

export MYSQL_PWD="$MYSQL_PASSWORD"

mysql_exec() {
  mysql \
    --host="$MYSQL_HOST" \
    --port="$MYSQL_PORT" \
    --user="$MYSQL_USER" \
    --batch \
    --skip-column-names \
    "$MYSQL_DATABASE" \
    -e "$1"
}

echo "ensuring visit_stats exists"
mysql_exec "
CREATE TABLE IF NOT EXISTS visit_stats (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  short_link_id bigint unsigned NOT NULL COMMENT '短链ID',
  visit_count bigint unsigned NOT NULL DEFAULT 0 COMMENT '访问次数',
  created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_short_link_id (short_link_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='短链访问统计表';
"

max_id_sql="SELECT COALESCE(MAX(id), 0) FROM short_links"
if [ "$END_ID" -gt 0 ]; then
  max_id="$END_ID"
else
  max_id="$(mysql_exec "$max_id_sql")"
fi

if [ -z "$max_id" ] || [ "$max_id" -eq 0 ]; then
  echo "no rows found in short_links"
  exit 0
fi

current_start="$START_ID"
if [ "$current_start" -eq 0 ]; then
  current_start=1
fi

echo "migrating visit_count from short_links to visit_stats"
echo "range: ${current_start}-${max_id}, batch_size=${BATCH_SIZE}, only_positive_visit=${ONLY_POSITIVE_VISIT}"

total_affected=0
batch_no=0
total_batches=$(( (max_id - current_start) / BATCH_SIZE + 1 ))

while [ "$current_start" -le "$max_id" ]; do
  batch_no=$((batch_no + 1))
  current_end=$((current_start + BATCH_SIZE - 1))
  if [ "$current_end" -gt "$max_id" ]; then
    current_end="$max_id"
  fi

  where_clause="id BETWEEN ${current_start} AND ${current_end}"
  if [ "$ONLY_POSITIVE_VISIT" = "1" ]; then
    where_clause="${where_clause} AND visit_count > 0"
  fi

  affected="$(mysql_exec "
INSERT INTO visit_stats (short_link_id, visit_count)
SELECT id, visit_count
FROM short_links
WHERE ${where_clause}
ON DUPLICATE KEY UPDATE
  visit_count = VALUES(visit_count);
SELECT ROW_COUNT();
")"

  affected="$(printf '%s\n' "$affected" | tail -n 1)"
  if [ -z "$affected" ]; then
    affected=0
  fi

  total_affected=$((total_affected + affected))
  echo "batch ${batch_no}/${total_batches}: id_range=${current_start}-${current_end}, affected_rows=${affected}, total_affected_rows=${total_affected}"

  current_start=$((current_end + 1))
done

final_rows="$(mysql_exec "SELECT COUNT(*) FROM visit_stats")"
echo "migration complete, total_affected_rows=${total_affected}, visit_stats_rows=${final_rows}"
