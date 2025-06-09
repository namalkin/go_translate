#!/bin/sh
set -e

echo "Waiting for MongoDB to be ready..."

# Ждём доступности порта MongoDB
until nc -z mongo 27017; do
  echo "Waiting for MongoDB..."
  sleep 2
done

echo "MongoDB is up. Running migrations..."
migrate-mongo up

echo "Migrations done. Exiting migrator."
