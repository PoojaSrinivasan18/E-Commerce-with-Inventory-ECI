from fastapi import FastAPI, HTTPException, Header, Depends
from db_utils import get_connection
from pydantic import BaseModel
from datetime import datetime
from typing import Optional, List
import uvicorn
import requests
import hashlib
import json
import time
import uuid

app = FastAPI(title="ECI Orders API", version="1.0")

# ----------- MODELS ------------

class OrderItem(BaseModel):
    product_id: int
    sku: str
    quantity: int

class PlaceOrderRequest(BaseModel):
    customer_id: int
    items: List[OrderItem]
    idempotency_key: Optional[str] = None

class OrderResponse(BaseModel):
    order_id: int
    customer_id: int
    order_status: str
    payment_status: str
    order_total: float
    created_at: datetime
    items: List[dict]

# ----------- CONFIGURATION ------------
INVENTORY_SERVICE_URL = "http://inventoryservice:3000/v1"
PAYMENT_SERVICE_URL = "http://payment_service:8002/v1"
NOTIFICATION_SERVICE_URL = "http://notification_service:8080/v1"
SHIPMENT_SERVICE_URL = "http://shipment_service:8001/v1"

TAX_RATE = 0.05  # 5% tax
SHIPPING_COST = 10.00

# ----------- HEALTH CHECK ------------

@app.get("/health")
def health_check():
    return {"status": "healthy", "service": "order"}

# ----------- ORDER ENDPOINTS WITH /V1 VERSIONING ------------

@app.get("/v1/orders")
def get_orders(limit: int = 10):
    conn = get_connection("order_db")
    cur = conn.cursor(dictionary=True)
    try:
        cur.execute("SELECT * FROM Orders ORDER BY created_at DESC LIMIT %s", (limit,))
        orders = cur.fetchall()
        return {"orders": orders, "count": len(orders)}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    finally:
        cur.close()
        conn.close()

@app.get("/v1/orders/{order_id}")
def get_order_by_id(order_id: int):
    conn = get_connection("order_db")
    cur = conn.cursor(dictionary=True)
    try:
        cur.execute("SELECT * FROM Orders WHERE order_id = %s", (order_id,))
        order = cur.fetchone()
        if not order:
            raise HTTPException(status_code=404, detail="Order not found")
        
        cur.execute("SELECT * FROM Order_Items WHERE order_id = %s", (order_id,))
        order["items"] = cur.fetchall()
        return order
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    finally:
        cur.close()
        conn.close()

@app.post("/v1/orders")
def place_order(order_request: PlaceOrderRequest, idempotency_key: str = Header(None, alias="Idempotency-Key")):
    """
    Place Order Workflow: Reserve → Pay → Ship
    1. Validate idempotency key
    2. Reserve inventory for each item
    3. Calculate totals with tax and shipping
    4. Process payment
    5. Confirm order or rollback on failure
    6. Send notifications
    """
    
    # Use provided idempotency key or generate one
    if not idempotency_key:
        if order_request.idempotency_key:
            idempotency_key = order_request.idempotency_key
        else:
            idempotency_key = str(uuid.uuid4())
    
    conn = get_connection("order_db")
    cur = conn.cursor(dictionary=True)
    
    try:
        # Check for existing order with same idempotency key
        cur.execute("SELECT * FROM Orders WHERE idempotency_key = %s", (idempotency_key,))
        existing_order = cur.fetchone()
        if existing_order:
            cur.execute("SELECT * FROM Order_Items WHERE order_id = %s", (existing_order['order_id'],))
            existing_order["items"] = cur.fetchall()
            return {"message": "Order already exists", "order": existing_order, "idempotent": True}
        
        # Generate order ID
        order_id = int(time.time() * 1000) % 2147483647  # Max MySQL INT
        
        # Step 1: Reserve inventory for all items
        reservations = []
        total_amount = 0.0
        order_items = []
        
        for item in order_request.items:
            # Get current price from catalog (would be implemented)
            # For now, using a mock price
            unit_price = 29.99  # Mock price
            
            # Reserve inventory
            reservation_data = {
                "product_id": item.product_id,
                "quantity": item.quantity,
                "idempotency_key": f"{idempotency_key}_{item.product_id}",
                "order_id": str(order_id)
            }
            
            try:
                response = requests.post(
                    f"{INVENTORY_SERVICE_URL}/inventory/reserve",
                    json=reservation_data,
                    timeout=10
                )
                
                if response.status_code == 200:
                    reservation_result = response.json()
                    reservations.append({
                        "idempotency_key": reservation_data["idempotency_key"],
                        "order_id": str(order_id)
                    })
                    
                    order_items.append({
                        "product_id": item.product_id,
                        "sku": item.sku,
                        "quantity": item.quantity,
                        "unit_price": unit_price
                    })
                    
                    total_amount += unit_price * item.quantity
                    
                else:
                    # Rollback previous reservations
                    for prev_reservation in reservations:
                        try:
                            requests.post(
                                f"{INVENTORY_SERVICE_URL}/inventory/release",
                                json=prev_reservation,
                                timeout=5
                            )
                        except:
                            pass
                    
                    raise HTTPException(
                        status_code=409,
                        detail=f"Failed to reserve inventory for product {item.product_id}: {response.text}"
                    )
                    
            except requests.RequestException as e:
                # Rollback previous reservations
                for prev_reservation in reservations:
                    try:
                        requests.post(
                            f"{INVENTORY_SERVICE_URL}/inventory/release",
                            json=prev_reservation,
                            timeout=5
                        )
                    except:
                        pass
                
                raise HTTPException(status_code=503, detail=f"Inventory service unavailable: {str(e)}")
        
        # Step 2: Calculate totals
        subtotal = total_amount
        tax_amount = subtotal * TAX_RATE
        shipping_amount = SHIPPING_COST
        total_with_tax_shipping = subtotal + tax_amount + shipping_amount
        
        # Create totals signature (hash for tamper detection)
        totals_data = {
            "subtotal": subtotal,
            "tax": tax_amount,
            "shipping": shipping_amount,
            "total": total_with_tax_shipping
        }
        totals_signature = hashlib.sha256(json.dumps(totals_data, sort_keys=True).encode()).hexdigest()
        
        # Step 3: Process Payment
        payment_data = {
            "order_id": order_id,
            "amount": total_with_tax_shipping,
            "customer_id": order_request.customer_id,
            "idempotency_key": f"payment_{idempotency_key}"
        }
        
        payment_success = False
        payment_id = None
        
        try:
            # Mock payment processing (replace with actual payment service call)
            payment_success = True
            payment_id = f"pay_{int(time.time())}"
            
        except Exception as e:
            # Payment failed - release all reservations
            for reservation in reservations:
                try:
                    requests.post(
                        f"{INVENTORY_SERVICE_URL}/inventory/release",
                        json=reservation,
                        timeout=5
                    )
                except:
                    pass
            
            raise HTTPException(status_code=402, detail=f"Payment processing failed: {str(e)}")
        
        if not payment_success:
            # Payment failed - release all reservations
            for reservation in reservations:
                try:
                    requests.post(
                        f"{INVENTORY_SERVICE_URL}/inventory/release",
                        json=reservation,
                        timeout=5
                    )
                except:
                    pass
            
            raise HTTPException(status_code=402, detail="Payment was declined")
        
        # Step 4: Create order record
        cur.execute("""
            INSERT INTO Orders (order_id, customer_id, order_status, payment_status, 
                              order_total, subtotal, tax_amount, shipping_amount,
                              totals_signature, idempotency_key, created_at)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """, (order_id, order_request.customer_id, "CONFIRMED", "PAID",
              total_with_tax_shipping, subtotal, tax_amount, shipping_amount,
              totals_signature, idempotency_key, datetime.now()))
        
        # Create order items
        for item in order_items:
            cur.execute("""
                INSERT INTO Order_Items (order_id, product_id, sku, quantity, unit_price)
                VALUES (%s, %s, %s, %s, %s)
            """, (order_id, item["product_id"], item["sku"], 
                  item["quantity"], item["unit_price"]))
        
        conn.commit()
        
        # Step 5: Send notifications (async)
        try:
            notification_data = {
                "type": "ORDER_CONFIRMED",
                "customer_id": order_request.customer_id,
                "order_id": order_id,
                "message": f"Your order #{order_id} has been confirmed and payment processed successfully."
            }
            requests.post(
                f"{NOTIFICATION_SERVICE_URL}/notifications",
                json=notification_data,
                timeout=5
            )
        except:
            pass  # Don't fail order if notification fails
        
        # Step 6: Create shipment (optional)
        try:
            shipment_data = {
                "order_id": order_id,
                "customer_id": order_request.customer_id,
                "items": order_items
            }
            requests.post(
                f"{SHIPMENT_SERVICE_URL}/shipments",
                json=shipment_data,
                timeout=5
            )
        except:
            pass  # Don't fail order if shipment creation fails
        
        # Return successful order
        return {
            "message": "Order placed successfully",
            "order": {
                "order_id": order_id,
                "customer_id": order_request.customer_id,
                "order_status": "CONFIRMED",
                "payment_status": "PAID",
                "subtotal": subtotal,
                "tax_amount": tax_amount,
                "shipping_amount": shipping_amount,
                "order_total": total_with_tax_shipping,
                "payment_id": payment_id,
                "idempotency_key": idempotency_key,
                "created_at": datetime.now().isoformat(),
                "items": order_items
            }
        }
        
    except HTTPException:
        raise
    except Exception as e:
        # Rollback reservations on any error
        for reservation in reservations:
            try:
                requests.post(
                    f"{INVENTORY_SERVICE_URL}/inventory/release",
                    json=reservation,
                    timeout=5
                )
            except:
                pass
        
        conn.rollback()
        raise HTTPException(status_code=500, detail=f"Order processing failed: {str(e)}")
    
    finally:
        cur.close()
        conn.close()

@app.post("/v1/orders/{order_id}/cancel")
def cancel_order(order_id: int):
    """Cancel an order and release reservations"""
    conn = get_connection("order_db")
    cur = conn.cursor(dictionary=True)
    
    try:
        # Get order details
        cur.execute("SELECT * FROM Orders WHERE order_id = %s", (order_id,))
        order = cur.fetchone()
        
        if not order:
            raise HTTPException(status_code=404, detail="Order not found")
        
        if order["order_status"] in ["CANCELLED", "SHIPPED", "DELIVERED"]:
            raise HTTPException(status_code=400, detail=f"Cannot cancel order with status: {order['order_status']}")
        
        # Release inventory reservations
        if order.get("idempotency_key"):
            cur.execute("SELECT * FROM Order_Items WHERE order_id = %s", (order_id,))
            items = cur.fetchall()
            
            for item in items:
                try:
                    release_data = {
                        "idempotency_key": f"{order['idempotency_key']}_{item['product_id']}",
                        "order_id": str(order_id)
                    }
                    requests.post(
                        f"{INVENTORY_SERVICE_URL}/inventory/release",
                        json=release_data,
                        timeout=5
                    )
                except:
                    pass  # Continue with cancellation even if release fails
        
        # Update order status
        cur.execute(
            "UPDATE Orders SET order_status = %s WHERE order_id = %s",
            ("CANCELLED", order_id)
        )
        conn.commit()
        
        return {"message": "Order cancelled successfully", "order_id": order_id}
        
    except HTTPException:
        raise
    except Exception as e:
        conn.rollback()
        raise HTTPException(status_code=500, detail=str(e))
    finally:
        cur.close()
        conn.close()

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)