#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="${SCRIPT_DIR}/.."
WRK_DIR="${SCRIPT_DIR}"

BASE_URL="${1:-http://127.0.0.1:8888}"
THREADS="${THREADS:-4}"
CONNS="${CONNS:-100}"
DURATION="${DURATION:-30s}"
CONFIG_FILE="${CONFIG_FILE:-etc/mysurl1-api.yaml}"
STARTUP_WAIT_SECONDS="${STARTUP_WAIT_SECONDS:-20}"
SERVER_LOG="${SERVER_LOG:-wrk/run_batch_create.log}"
PORT="${PORT:-8888}"
BATCH_SIZE="${BATCH_SIZE:-20}"
URL_PREFIX="${URL_PREFIX:-https://github.com/JCCGGKS/mysurl-batch-}"
AUTH_TOKEN="${AUTH_TOKEN:-}"
PROVIDER=""

SERVER_PID=""
START_TS="$(date +%s)"

cleanup_port_processes() {
  local action="${1:-cleaning}"
  local pids
  pids="$(lsof -tiTCP:${PORT} 2>/dev/null || true)"
  if [[ -z "${pids}" ]]; then
    return
  fi

  echo "${action} processes on port ${PORT}: ${pids}"
  while IFS= read -r pid; do
    [[ -z "${pid}" ]] && continue
    kill -TERM "${pid}" 2>/dev/null || true
  done <<< "${pids}"

  sleep 1

  pids="$(lsof -tiTCP:${PORT} 2>/dev/null || true)"
  if [[ -n "${pids}" ]]; then
    while IFS= read -r pid; do
      [[ -z "${pid}" ]] && continue
      kill -INT "${pid}" 2>/dev/null || true
    done <<< "${pids}"
  fi
}

cleanup() {
  cleanup_port_processes "stopping"

  if [[ -n "${SERVER_PID}" ]]; then
    wait "${SERVER_PID}" 2>/dev/null || true
  fi

  echo "service stopped"
}

trap cleanup EXIT

if [[ -z "${AUTH_TOKEN}" ]]; then
  echo "AUTH_TOKEN is required for /api/v1/links/batch benchmark" >&2
  exit 1
fi

cd "${ROOT_DIR}"

PROVIDER="$(sed -n '/^Short:/,/^[^ ]/p' "${CONFIG_FILE}" | sed -n 's/^[[:space:]]*Provider:[[:space:]]*//p' | head -n 1)"

cleanup_port_processes "cleaning"

go run mysurl1.go -f "${CONFIG_FILE}" >"${SERVER_LOG}" 2>&1 &
SERVER_PID=$!

for ((i=1; i<=STARTUP_WAIT_SECONDS; i++)); do
  if ! kill -0 "${SERVER_PID}" 2>/dev/null; then
    echo "service exited early, log: ${SERVER_LOG}" >&2
    exit 1
  fi

  if curl -sS -o /dev/null "${BASE_URL}"; then
    break
  fi

  sleep 1
done

if ! kill -0 "${SERVER_PID}" 2>/dev/null; then
  echo "service exited before ready check completed, log: ${SERVER_LOG}" >&2
  exit 1
fi

if ! curl -sS -o /dev/null "${BASE_URL}"; then
  echo "service did not become ready within ${STARTUP_WAIT_SECONDS}s, log: ${SERVER_LOG}" >&2
  exit 1
fi

READY_TS="$(date +%s)"
echo "service ready in $((READY_TS - START_TS))s, log: ${SERVER_LOG}"
echo "provider: ${PROVIDER:-unknown}"
echo "batch_size: ${BATCH_SIZE}"

cd "${WRK_DIR}"

echo "running: AUTH_TOKEN=*** BATCH_SIZE=${BATCH_SIZE} wrk -t${THREADS} -c${CONNS} -d${DURATION} -H 'Authorization: Bearer ***' -s post_batch_create.lua ${BASE_URL}"

AUTH_TOKEN="${AUTH_TOKEN}" \
BATCH_SIZE="${BATCH_SIZE}" \
URL_PREFIX="${URL_PREFIX}" \
wrk \
  -t"${THREADS}" \
  -c"${CONNS}" \
  -d"${DURATION}" \
  -H "Authorization: Bearer ${AUTH_TOKEN}" \
  -s post_batch_create.lua \
  "${BASE_URL}"

cleanup
trap - EXIT
