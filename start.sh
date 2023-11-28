#!/bin/sh

# apabila output bukan 0, maka error
set -e

echo "run db migration"
# membuat dev.env berada pada lingkungan variabel
source /app/dev.env
# jalankan migrasi
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
# jalankan program dengan parameter yang ada yaitu /app/main
exec "$@"