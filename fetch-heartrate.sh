#!/bin/zsh

FILE="heartrate.csv"

if [[ -z "${PULSOID_TOKEN:-}" ]]; then
  echo "Error: PULSOID_TOKEN environment variable not set"
  exit 1
fi

if [ ! -f "$FILE" ]; then
  echo "fetch_ts_ms,measured_at_ms,heart_rate" > "$FILE"
fi

while true; do
  FETCH_TS=$(($(date +%s) * 1000))

  RESP=$(curl -s -H "Authorization: Bearer $PULSOID_TOKEN" \
    "https://dev.pulsoid.net/api/v1/data/heart_rate/latest?scope=data:heart_rate:read")

  MEASURED_AT=$(echo "$RESP" | jq '.measured_at')
  HR=$(echo "$RESP" | jq '.data.heart_rate')

  echo "$FETCH_TS,$MEASURED_AT,$HR" >> "$FILE"

  sleep 1
done
