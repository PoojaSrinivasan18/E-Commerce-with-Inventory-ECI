import os
import pandas as pd
from db_utils import get_connection
from pathlib import Path

SERVICE_DIR = Path(__file__).resolve().parent
CSV_DIR = Path(os.getenv("CSV_DIR", str(SERVICE_DIR / "csv_files")))

def create_database_and_table():
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("CREATE DATABASE IF NOT EXISTS shipping_db")
    conn.commit()
    cur.close()
    conn.close()

    conn = get_connection("shipping_db")
    cur = conn.cursor()
    cur.execute("""
        CREATE TABLE IF NOT EXISTS Shipments (
            shipment_id INT PRIMARY KEY,
            order_id INT,
            carrier VARCHAR(100),
            status VARCHAR(50),
            tracking_no VARCHAR(100),
            shipped_at DATETIME,
            delivered_at DATETIME
        )
    """)
    conn.commit()
    conn.close()
    print("✅ shipping_db and table created")


# --- OPTIONAL: teardown SQL (commented) ---
# Use only when you intentionally want to drop data. Keep commented to avoid accidental execution.
#
# SQL to drop table:
# DROP TABLE IF EXISTS shipping_db.Shipments;
#
# SQL to drop database:
# DROP DATABASE IF EXISTS shipping_db;


def load_csv():
    file_path = CSV_DIR / "Shipments.csv"
    df = pd.read_csv(str(file_path))
    print(f"📥 Loading Shipments.csv ({len(df)} rows) into shipping_db.Shipments")
    conn = get_connection("shipping_db")
    cur = conn.cursor()
    cols = ", ".join(df.columns)
    placeholders = ", ".join(["%s"] * len(df.columns))
    for _, row in df.iterrows():
        cur.execute(f"INSERT INTO Shipments ({cols}) VALUES ({placeholders})", tuple(row))
    conn.commit()
    conn.close()
    print(f"✅ Inserted {len(df)} rows into shipping_db.Shipments")


if __name__ == "__main__":
    create_database_and_table()
    load_csv()
