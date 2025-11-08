#!/bin/bash
set -e

# This script runs only once when the database volume is first initialized

echo "Creating additional databases..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE customer_db;
    CREATE DATABASE inventory_db;
    CREATE DATABASE payment_db;
EOSQL

echo "✅ Databases created: catalog_db, customer_db, inventory_db, payment_db"