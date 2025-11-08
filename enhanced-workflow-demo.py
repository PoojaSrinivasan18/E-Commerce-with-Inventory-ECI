#!/usr/bin/env python3
"""
Enhanced E-Commerce Interservice Workflow Demonstration
Implements the complete Place Order workflow: Reserve → Pay → Ship
As per Problem Statement 4 requirements
"""

import requests
import json
import time
import uuid
from datetime import datetime

# Service URLs (updated for /v1 versioning)
SERVICES = {
    "catalog": "http://localhost:3000/v1",
    "inventory": "http://localhost:3002/v1", 
    "order": "http://localhost:8000/v1",
    "payment": "http://localhost:8002/v1",
    "notification": "http://localhost:8080/v1",
    "shipment": "http://localhost:8001/v1"
}

def print_step(step_num, title, emoji="🔄"):
    print(f"\n{emoji} STEP {step_num}: {title}")
    print("=" * 60)

def make_request(method, url, **kwargs):
    """Make HTTP request with error handling"""
    try:
        response = requests.request(method, url, timeout=10, **kwargs)
        return response
    except requests.RequestException as e:
        print(f"❌ Request failed: {e}")
        return None

def test_service_health():
    """Test all service health endpoints"""
    print("🏥 TESTING SERVICE HEALTH")
    print("=" * 60)
    
    health_urls = {
        "Catalog": "http://localhost:3000/health",
        "Inventory": "http://localhost:3002/health", 
        "Order": "http://localhost:8000/health",
        "Payment": "http://localhost:8002/health",
        "Notification": "http://localhost:8080/health",
        "Shipment": "http://localhost:8001/health"
    }
    
    for service, url in health_urls.items():
        response = make_request("GET", url)
        if response and response.status_code == 200:
            print(f"✅ {service} Service: HEALTHY")
        else:
            print(f"❌ {service} Service: UNHEALTHY")

def demonstrate_place_order_workflow():
    """
    Complete Place Order Workflow Implementation:
    1. Check inventory availability 
    2. Reserve inventory with idempotency
    3. Calculate totals with tax and shipping
    4. Process payment with idempotency
    5. Confirm order or rollback on failure
    6. Send notifications
    7. Create shipment
    """
    
    print("\n🛒 COMPLETE E-COMMERCE ORDER WORKFLOW")
    print("=" * 60)
    
    # Generate unique identifiers
    idempotency_key = str(uuid.uuid4())
    customer_id = 1
    order_items = [
        {"product_id": 1, "sku": "LAPTOP001", "quantity": 2},
        {"product_id": 2, "sku": "MOUSE001", "quantity": 1}
    ]
    
    print(f"🔑 Order ID: {idempotency_key}")
    print(f"👤 Customer ID: {customer_id}")
    print(f"📦 Items: {len(order_items)} products")
    
    # Step 1: Check Inventory Availability
    print_step(1, "Check Inventory Availability", "📊")
    
    available_items = []
    for item in order_items:
        url = f"{SERVICES['inventory']}/inventory/availability/{item['product_id']}"
        response = make_request("GET", url)
        
        if response and response.status_code == 200:
            availability = response.json()
            print(f"📦 Product {item['product_id']}: {availability['total_available']} available")
            
            if availability['total_available'] >= item['quantity']:
                available_items.append(item)
                print(f"✅ Sufficient stock for {item['quantity']} units")
            else:
                print(f"❌ Insufficient stock: need {item['quantity']}, have {availability['total_available']}")
        else:
            print(f"❌ Could not check availability for product {item['product_id']}")
    
    if not available_items:
        print("❌ No items available for order")
        return None
    
    # Step 2: Reserve Inventory (Atomic with Idempotency)
    print_step(2, "Reserve Inventory with TTL", "🔒")
    
    reservations = []
    for item in available_items:
        reservation_data = {
            "product_id": item["product_id"],
            "quantity": item["quantity"], 
            "idempotency_key": f"{idempotency_key}_product_{item['product_id']}",
            "order_id": idempotency_key
        }
        
        url = f"{SERVICES['inventory']}/inventory/reserve"
        response = make_request("POST", url, json=reservation_data)
        
        if response and response.status_code == 200:
            reservation = response.json()
            reservations.append(reservation_data)
            print(f"✅ Reserved {item['quantity']} units of product {item['product_id']}")
            print(f"⏰ Expires: {reservation.get('expires_at', 'N/A')}")
        else:
            print(f"❌ Failed to reserve product {item['product_id']}")
            # Rollback previous reservations
            rollback_reservations(reservations)
            return None
    
    # Step 3: Calculate Order Totals
    print_step(3, "Calculate Order Totals with Tax & Shipping", "💰")
    
    subtotal = sum(29.99 * item['quantity'] for item in available_items)  # Mock pricing
    tax_rate = 0.05  # 5% tax
    tax_amount = subtotal * tax_rate
    shipping_amount = 10.00
    total_amount = subtotal + tax_amount + shipping_amount
    
    print(f"💵 Subtotal: ${subtotal:.2f}")
    print(f"📈 Tax (5%): ${tax_amount:.2f}")
    print(f"🚚 Shipping: ${shipping_amount:.2f}")
    print(f"💳 Total: ${total_amount:.2f}")
    
    # Step 4: Process Payment with Idempotency
    print_step(4, "Process Payment with Idempotency", "💳")
    
    payment_data = {
        "order_id": hash(idempotency_key) % 2147483647,  # Convert to int
        "amount": total_amount,
        "customer_id": customer_id,
        "method": "CREDIT_CARD",
        "idempotency_key": f"payment_{idempotency_key}"
    }
    
    url = f"{SERVICES['payment']}/payments/charge"
    response = make_request("POST", url, json=payment_data)
    
    if response and response.status_code == 200:
        payment_result = response.json()
        print(f"✅ Payment processed successfully")
        print(f"💳 Payment ID: {payment_result.get('payment', {}).get('payment_id', 'N/A')}")
        print(f"🔗 Reference: {payment_result.get('payment', {}).get('reference', 'N/A')}")
    else:
        print("❌ Payment failed - rolling back reservations")
        rollback_reservations(reservations)
        return None
    
    # Step 5: Create Order Record 
    print_step(5, "Create Order Record", "📄")
    
    order_data = {
        "customer_id": customer_id,
        "items": available_items,
        "idempotency_key": idempotency_key
    }
    
    # Add idempotency header
    headers = {"Idempotency-Key": idempotency_key}
    url = f"{SERVICES['order']}/orders"
    response = make_request("POST", url, json=order_data, headers=headers)
    
    if response and response.status_code == 200:
        order_result = response.json()
        order_id = order_result.get('order', {}).get('order_id')
        print(f"✅ Order created successfully")
        print(f"📋 Order ID: {order_id}")
        print(f"📊 Status: {order_result.get('order', {}).get('order_status', 'N/A')}")
    else:
        print("❌ Order creation failed")
        return None
    
    # Step 6: Send Order Confirmation Notification
    print_step(6, "Send Order Confirmation", "📧")
    
    notification_data = {
        "type": "ORDER_CONFIRMED",
        "customer_id": customer_id,
        "order_id": order_id,
        "message": f"Your order #{order_id} has been confirmed. Total: ${total_amount:.2f}"
    }
    
    url = f"{SERVICES['notification']}/notifications"
    response = make_request("POST", url, json=notification_data)
    
    if response and response.status_code in [200, 201]:
        print("✅ Order confirmation sent")
    else:
        print("⚠️ Notification failed (non-critical)")
    
    # Step 7: Create Shipment
    print_step(7, "Create Shipment Record", "📦")
    
    shipment_data = {
        "order_id": order_id,
        "customer_id": customer_id,
        "items": available_items,
        "carrier": "UPS",
        "service_type": "GROUND"
    }
    
    url = f"{SERVICES['shipment']}/shipments"
    response = make_request("POST", url, json=shipment_data)
    
    if response and response.status_code in [200, 201]:
        shipment_result = response.json()
        tracking_no = shipment_result.get('shipment', {}).get('tracking_no', 'N/A')
        print(f"✅ Shipment created")
        print(f"📋 Tracking: {tracking_no}")
    else:
        print("⚠️ Shipment creation failed (non-critical)")
    
    # Step 8: Mark Inventory as Shipped
    print_step(8, "Mark Inventory as Shipped", "🚚")
    
    for reservation in reservations:
        ship_data = {
            "idempotency_key": reservation["idempotency_key"],
            "order_id": reservation["order_id"]
        }
        
        url = f"{SERVICES['inventory']}/inventory/ship"
        response = make_request("POST", url, json=ship_data)
        
        if response and response.status_code == 200:
            print(f"✅ Shipped product {reservation.get('product_id')}")
        else:
            print(f"⚠️ Failed to mark product {reservation.get('product_id')} as shipped")
    
    # Final Summary
    print("\n🎉 ORDER WORKFLOW COMPLETED SUCCESSFULLY!")
    print("=" * 60)
    print(f"📋 Order ID: {order_id}")
    print(f"💳 Payment: ${total_amount:.2f}")
    print(f"📦 Items: {len(available_items)} products")
    print(f"⏰ Processing Time: Complete")
    print(f"🔑 Idempotency Key: {idempotency_key}")
    
    return {
        "order_id": order_id,
        "total_amount": total_amount,
        "items_count": len(available_items),
        "idempotency_key": idempotency_key
    }

def rollback_reservations(reservations):
    """Rollback inventory reservations on failure"""
    print("🔄 Rolling back reservations...")
    
    for reservation in reservations:
        release_data = {
            "idempotency_key": reservation["idempotency_key"],
            "order_id": reservation["order_id"]
        }
        
        url = f"{SERVICES['inventory']}/inventory/release"
        response = make_request("POST", url, json=release_data)
        
        if response and response.status_code == 200:
            print(f"✅ Released reservation for product {reservation.get('product_id')}")
        else:
            print(f"❌ Failed to release reservation for product {reservation.get('product_id')}")

def demonstrate_error_scenarios():
    """Demonstrate error handling and rollback scenarios"""
    
    print("\n⚠️  ERROR HANDLING & ROLLBACK SCENARIOS")
    print("=" * 60)
    
    # Test insufficient inventory
    print("\n📊 Testing Insufficient Inventory Scenario")
    print("-" * 40)
    
    large_order = {
        "product_id": 1,
        "quantity": 1000,  # Intentionally large
        "idempotency_key": f"test_{uuid.uuid4()}",
        "order_id": "test_insufficient"
    }
    
    url = f"{SERVICES['inventory']}/inventory/reserve"
    response = make_request("POST", url, json=large_order)
    
    if response and response.status_code == 409:
        print("✅ Correctly rejected insufficient inventory")
    else:
        print("⚠️ Insufficient inventory handling needs review")
    
    # Test idempotency
    print("\n🔑 Testing Idempotency")
    print("-" * 40)
    
    idem_key = str(uuid.uuid4())
    reservation_data = {
        "product_id": 1,
        "quantity": 1,
        "idempotency_key": idem_key,
        "order_id": "test_idem"
    }
    
    # First request
    url = f"{SERVICES['inventory']}/inventory/reserve"
    response1 = make_request("POST", url, json=reservation_data)
    
    # Second request with same key
    response2 = make_request("POST", url, json=reservation_data)
    
    if (response1 and response1.status_code == 200 and 
        response2 and response2.status_code == 200):
        result2 = response2.json()
        if result2.get('idempotent'):
            print("✅ Idempotency working correctly")
        else:
            print("⚠️ Idempotency needs review")
    else:
        print("⚠️ Idempotency test failed")

def check_reservation_status():
    """Check reservation cleanup job status"""
    print("\n⏰ RESERVATION CLEANUP STATUS")
    print("=" * 60)
    
    url = f"{SERVICES['inventory']}/inventory/reservations/status"
    response = make_request("GET", url)
    
    if response and response.status_code == 200:
        status = response.json()
        print("📊 Reservation Statistics:")
        print(f"   Active: {status.get('reservation_stats', {}).get('active_reservations', 0)}")
        print(f"   Expired: {status.get('reservation_stats', {}).get('expired_reservations', 0)}")
        print(f"   Shipped: {status.get('reservation_stats', {}).get('shipped_reservations', 0)}")
        print(f"   Released: {status.get('reservation_stats', {}).get('released_reservations', 0)}")
        print(f"⏰ Expiring in 1 hour: {status.get('expiring_in_1_hour', 0)}")
        print(f"🔧 Cleanup job: {'Active' if status.get('cleanup_active') else 'Inactive'}")
    else:
        print("❌ Could not retrieve reservation status")

def main():
    """Main demonstration function"""
    print("🚀 E-COMMERCE MICROSERVICES DEMONSTRATION")
    print("Problem Statement 4 - Complete Implementation")
    print("=" * 60)
    print(f"⏰ Started at: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    
    # Test service health
    test_service_health()
    
    # Check reservation cleanup status
    check_reservation_status()
    
    # Demonstrate complete order workflow
    result = demonstrate_place_order_workflow()
    
    if result:
        # Demonstrate error scenarios
        demonstrate_error_scenarios()
        
        print(f"\n✅ DEMONSTRATION COMPLETED SUCCESSFULLY")
        print(f"📋 Final Order: {result['order_id']}")
        print(f"💰 Total Processed: ${result['total_amount']:.2f}")
    else:
        print("\n❌ DEMONSTRATION FAILED")
    
    print(f"⏰ Completed at: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")

if __name__ == "__main__":
    main()