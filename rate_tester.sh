#!/bin/bash

TARGET_URL="http://localhost:8080" # The target URL for the GET requests
REQUESTS_PER_SECOND=${1:-1}        # Default to 1 request per second if no argument is provided

if ! [[ "$REQUESTS_PER_SECOND" =~ ^[0-9]+$ ]] || [ "$REQUESTS_PER_SECOND" -le 0 ]; then
  echo "Error: REQUESTS_PER_SECOND must be a positive integer."
  echo "Usage: $0 <requests_per_second>"
  exit 1
fi

DELAY_MICROSECONDS=$(awk "BEGIN { printf \"%.0f\", 1000000 / $REQUESTS_PER_SECOND }")

echo "Sending $REQUESTS_PER_SECOND GET requests per second to $TARGET_URL"
echo "Press Ctrl+C to stop."
echo "---"

while true; do
  START_TIME_MICROSECONDS=$(date +%s%N | cut -b1-16) # Current time in microseconds

  for ((i = 0; i < $REQUESTS_PER_SECOND; i++)); do
    STATUS_CODE=$(curl -o /dev/null -s -w "%{http_code}\n" "$TARGET_URL" -H "X-User-Id: 123")
    echo "Status Code: $STATUS_CODE"
  done

  END_TIME_MICROSECONDS=$(date +%s%N | cut -b1-16) # Time after sending N requests

  ELAPSED_MICROSECONDS=$((END_TIME_MICROSECONDS - START_TIME_MICROSECONDS))

  SLEEP_DURATION_MICROSECONDS=$((1000000 - ELAPSED_MICROSECONDS))

  if [ "$SLEEP_DURATION_MICROSECONDS" -gt 0 ]; then
    SLEEP_SECONDS=$(awk "BEGIN { printf \"%.6f\", $SLEEP_DURATION_MICROSECONDS / 1000000 }")
    sleep "$SLEEP_SECONDS"
  fi
done
