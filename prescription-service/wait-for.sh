#!/bin/sh
set -e

host="$1"
shift
cmd="$@"

echo "⏳ Waiting for $host to be ready..."
until nc -z ${host%:*} ${host#*:}; do
  sleep 1
done

echo "✅ $host is up!"
exec $cmd
