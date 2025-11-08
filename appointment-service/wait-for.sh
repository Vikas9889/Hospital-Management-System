#!/bin/sh
set -e
host="postgres-appointments"
port=5432
echo "Waiting for Postgres to be ready on $host:$port..."
until nc -z $host $port; do
  sleep 1
done
echo "Postgres is up - starting application"
exec "$@"
