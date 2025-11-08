#!/bin/bash
# E-Commerce Interservice Communication Manual Demo
# This script demonstrates the interservice calls step by step

echo "🏪 E-COMMERCE INTERSERVICE COMMUNICATION DEMO"
echo "=============================================="
echo ""

# Generate unique IDs for this demo
ORDER_ID=$((3000 + RANDOM % 1000))
SHIPMENT_ID=$((RANDOM))
TRACKING_NO="FS$SHIPMENT_ID"

echo "📋 Demo Setup:"
echo "   • Order ID: $ORDER_ID"
echo "   • Shipment ID: $SHIPMENT_ID"
echo "   • Tracking Number: $TRACKING_NO"
echo ""

# Step 1: Create Shipment
echo "🔄 STEP 1: Creating shipment record..."
echo "POST http://localhost:8001/shipments"

curl -X POST "http://localhost:8001/shipments" \
  -H "Content-Type: application/json" \
  -d "{
    \"shipment_id\": $SHIPMENT_ID,
    \"order_id\": $ORDER_ID,
    \"carrier\": \"FastShip Express\",
    \"status\": \"PREPARING\", 
    \"tracking_no\": \"$TRACKING_NO\",
    \"shipped_at\": null,
    \"delivered_at\": null
  }"

echo -e "\n"

# Step 2: Send Order Notification
echo "🔄 STEP 2: Sending order preparation notification..."
echo "POST http://localhost:8080/v1/notifications/email"

curl -X POST "http://localhost:8080/v1/notifications/email" \
  -H "Content-Type: application/json" \
  -d "{
    \"orderId\": $ORDER_ID,
    \"shipmentId\": $SHIPMENT_ID,
    \"type\": \"SHIPMENT\",
    \"channel\": \"EMAIL\",
    \"messageContent\": \"Your order #$ORDER_ID is being prepared for shipment\",
    \"customerEmail\": \"customer@example.com\",
    \"subject\": \"Order #$ORDER_ID - Preparing for Shipment\",
    \"message\": \"Great news! Your order #$ORDER_ID is being prepared for shipment with FastShip Express. Tracking: $TRACKING_NO\"
  }"

echo -e "\n"

# Step 3: Send Shipped Notification  
echo "🔄 STEP 3: Sending shipment notification..."
echo "POST http://localhost:8080/v1/notifications/email"

curl -X POST "http://localhost:8080/v1/notifications/email" \
  -H "Content-Type: application/json" \
  -d "{
    \"orderId\": $ORDER_ID,
    \"shipmentId\": $SHIPMENT_ID,
    \"type\": \"SHIPMENT\", 
    \"channel\": \"EMAIL\",
    \"messageContent\": \"Your order #$ORDER_ID has been shipped! Tracking: $TRACKING_NO\",
    \"customerEmail\": \"customer@example.com\",
    \"subject\": \"Order #$ORDER_ID - Shipped!\",
    \"message\": \"Excellent! Your order #$ORDER_ID has been shipped via FastShip Express. Track your package: $TRACKING_NO\"
  }"

echo -e "\n"

# Step 4: Send Delivery Confirmation
echo "🔄 STEP 4: Sending delivery confirmation..."
echo "POST http://localhost:8080/v1/notifications/email"

curl -X POST "http://localhost:8080/v1/notifications/email" \
  -H "Content-Type: application/json" \
  -d "{
    \"orderId\": $ORDER_ID,
    \"shipmentId\": $SHIPMENT_ID,
    \"type\": \"SHIPMENT\",
    \"channel\": \"EMAIL\", 
    \"messageContent\": \"Your order #$ORDER_ID has been delivered successfully!\",
    \"customerEmail\": \"customer@example.com\",
    \"subject\": \"Order #$ORDER_ID - Delivered Successfully!\",
    \"message\": \"🎉 Great news! Your order #$ORDER_ID has been delivered successfully. Thank you for choosing us!\"
  }"

echo -e "\n"

echo "✅ INTERSERVICE COMMUNICATION DEMO COMPLETED!"
echo ""
echo "📊 Summary:"
echo "   • Created shipment record in Shipment Service"
echo "   • Sent 3 notifications via Notification Service"
echo "   • Demonstrated order → shipment → notification flow"
echo "   • Total API calls: 4"
echo ""
echo "🔗 Services Used:"
echo "   • Shipment Service (Port 8001): 1 POST /shipments"
echo "   • Notification Service (Port 8080): 3 POST /v1/notifications/email"