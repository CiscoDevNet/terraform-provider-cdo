#!/bin/bash
# Backoff in an geometric progression for a total of 1,022 seconds; calculated coz sum = 2^(n+1) - 1 - 2^0
for EXPONENTIAL_BACKOFF in {1..9}; do
    nc -w 1 -z $1 443 && break;
    DELAY=$((2**$EXPONENTIAL_BACKOFF))
    echo "ASAv not yet available, sleeping for $DELAY seconds"
    sleep $DELAY
done