#!/bin/sh
set -e

WEBHOOK_URL=""

send_to_slack() {
  curl -X POST -H 'Content-type: application/json' --data '{"text":"IP address is:\n```\n'"$1"'\n```"}' "$2" 
}

prev=""
while true; do
  address="$(hostname -I)" 
  if [ "$prev" != "$address" ]; then
    echo "IP address is:"
    echo "$address"

    if [ -n "$WEBHOOK_URL" ] && [ "$address" != "" ]; then
      echo "Sending information to Slack..."
      send_to_slack "$address" "$WEBHOOK_URL"
    fi
  fi
  prev="$address"
  sleep 20
done
