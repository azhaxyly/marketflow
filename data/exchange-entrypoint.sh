#!/bin/bash

PORT=$1
EXCHANGE_NAME=$2

echo "Starting exchange server $EXCHANGE_NAME on port $PORT..."

while true; do
  { 
    while true; do
      for PAIR in BTCUSDT ETHUSDT DOGEUSDT TONUSDT SOLUSDT; do
        PRICE=$(awk -v min=100 -v max=1000 'BEGIN{srand(); print min+rand()*(max-min)}')
        TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
        JSON="{\"exchange\":\"$EXCHANGE_NAME\",\"pair\":\"$PAIR\",\"price\":$PRICE,\"time\":\"$TIMESTAMP\"}"
        echo "$JSON"
      done
      sleep 1
    done
  } | nc -lk -p $PORT
done
