from fastapi import FastAPI, HTTPException, Query
from db_utils import get_connection
from pydantic import BaseModel
from datetime import datetime
import uvicorn

app = FastAPI(title="ECI Shipments API", version="1.0")

class Shipment(BaseModel):
    shipment_id: int
    order_id: int
    carrier: str
    status: str
    tracking_no: str
    shipped_at: datetime | None = None
    delivered_at: datetime | None = None


@app.get("/shipments")
def get_shipments(limit: int = Query(default=10, gt=0, description="Limit number of shipments to fetch")):
    conn = get_connection("shipping_db")
    cur = conn.cursor(dictionary=True)
    try:
        cur.execute(f"SELECT * FROM Shipments LIMIT {limit}")
        shipments = cur.fetchall()
        return shipments
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    finally:
        conn.close()


@app.post("/shipments")
def add_shipment(shipment: Shipment):
    conn = get_connection("shipping_db")
    cur = conn.cursor()
    try:
        cur.execute("""
            INSERT INTO Shipments (shipment_id, order_id, carrier, status, tracking_no, shipped_at, delivered_at)
            VALUES (%s, %s, %s, %s, %s, %s, %s)
        """, (shipment.shipment_id, shipment.order_id, shipment.carrier, shipment.status,
              shipment.tracking_no, shipment.shipped_at, shipment.delivered_at))
        conn.commit()
        return {"message": "✅ Shipment inserted successfully"}
    except Exception as e:
        conn.rollback()
        raise HTTPException(status_code=500, detail=str(e))
    finally:
        cur.close()
        conn.close()


@app.get("/shipments/{shipment_id}")
def get_shipment_by_id(shipment_id: int):
    conn = get_connection("shipping_db")
    cur = conn.cursor(dictionary=True)
    try:
        cur.execute("SELECT * FROM Shipments WHERE shipment_id = %s", (shipment_id,))
        shipment = cur.fetchone()
        if not shipment:
            raise HTTPException(status_code=404, detail="Shipment not found")
        return shipment
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    finally:
        conn.close()

@app.put("/shipments/{shipment_id}")
def update_shipment(shipment_id: int, updated_shipment: Shipment):
    """Update a shipment."""
    conn = get_connection("shipping_db")
    cur = conn.cursor()
    try:
        cur.execute("""
            UPDATE Shipments 
            SET order_id=%s, carrier=%s, status=%s, tracking_no=%s, shipped_at=%s, delivered_at=%s 
            WHERE shipment_id=%s
        """, (updated_shipment.order_id, updated_shipment.carrier, updated_shipment.status,
              updated_shipment.tracking_no, updated_shipment.shipped_at, updated_shipment.delivered_at, shipment_id))
        conn.commit()

        if cur.rowcount == 0:
            raise HTTPException(status_code=404, detail="Shipment not found")

        return {"message": "✅ Shipment updated successfully"}
    except Exception as e:
        conn.rollback()
        raise HTTPException(status_code=500, detail=str(e))
    finally:
        conn.close()

@app.delete("/shipments/{shipment_id}")
def delete_shipment(shipment_id: int):
    """Delete a shipment."""
    conn = get_connection("shipping_db")
    cur = conn.cursor()
    try:
        cur.execute("DELETE FROM Shipments WHERE shipment_id = %s", (shipment_id,))
        conn.commit()

        if cur.rowcount == 0:
            raise HTTPException(status_code=404, detail="Shipment not found")

        return {"message": f"🗑️ Shipment {shipment_id} deleted successfully"}
    except Exception as e:
        conn.rollback()
        raise HTTPException(status_code=500, detail=str(e))
    finally:
        conn.close()

if __name__ == "__main__":
    uvicorn.run("main:app", host="127.0.0.1", port=8001, reload=True)
