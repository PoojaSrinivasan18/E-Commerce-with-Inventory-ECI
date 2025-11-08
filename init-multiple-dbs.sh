#!/bin/bash
set -e

# This script runs only once when the database volume is first initialized

echo "🔧 Creating additional databases..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE customer_db;
    CREATE DATABASE inventory_db;
    CREATE DATABASE payment_db;
    CREATE DATABASE notifications;
    
    -- Grant all privileges to the user on all databases
    GRANT ALL PRIVILEGES ON DATABASE catalog_db TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON DATABASE customer_db TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON DATABASE inventory_db TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON DATABASE payment_db TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON DATABASE notifications TO $POSTGRES_USER;
EOSQL

echo "✅ Databases created successfully: catalog_db, customer_db, inventory_db, payment_db, notifications"
echo "✅ Privileges granted to user: $POSTGRES_USER"