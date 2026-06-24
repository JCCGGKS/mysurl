#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="${SCRIPT_DIR}/.."
WRK_DIR="${SCRIPT_DIR}"

BASE_URL="${1:-http://127.0.0.1:8888}"
SHORT_CODE="${SHORT_CODE:-${2:-2sw4YUefiX6}}"
THREADS="${THREADS:-4}"
CONNS="${CONNS:-100}"
DURATION="${DURATION:-30s}"
CONFIG_FILE="${CONFIG_FILE:-etc/mysurl1-api.yaml}"
STARTUP_WAIT_SECONDS="${STARTUP_WAIT_SECONDS:-20}"
SERVER_LOG="${SERVER_LOG:-wrk/run_get.log}"
PORT="${PORT:-8888}"
PROVIDER=""

SERVER_PID=""
START_TS="$(date +%s)"

cleanup() {
  local pids

  pids="$(lsof -tiTCP:${PORT} 2>/dev/null || true)"
  if [[ -n "${pids}" ]]; then
    echo "stopping processes on port ${PORT}: ${pids}"
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
  fi

  if [[ -n "${SERVER_PID}" ]]; then
    wait "${SERVER_PID}" 2>/dev/null || true
  fi

  echo "service stopped"
}

trap cleanup EXIT

cd "${ROOT_DIR}"

PROVIDER="$(sed -n '/^Short:/,/^[^ ]/p' "${CONFIG_FILE}" | sed -n 's/^[[:space:]]*Provider:[[:space:]]*//p' | head -n 1)"

cleanup

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
echo "short_code: ${SHORT_CODE}"

cd "${WRK_DIR}"

echo "running: SHORT_CODE=${SHORT_CODE} wrk -t${THREADS} -c${CONNS} -d${DURATION} -s get_code.lua ${BASE_URL}"

SHORT_CODE="${SHORT_CODE}" wrk \
  -t"${THREADS}" \
  -c"${CONNS}" \
  -d"${DURATION}" \
  -s get_code.lua \
  "${BASE_URL}"

cleanup
trap - EXIT
