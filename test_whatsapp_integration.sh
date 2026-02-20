#!/bin/bash

echo "🧪 Testing Facebook WhatsApp Business API Integration"
echo "=================================================="

# Test the WhatsApp API endpoint
echo "1️⃣ Testing WhatsApp API endpoint..."
curl -X POST http://localhost:3001/api/whatsapp \
  -H "Content-Type: application/json" \
  -d '{
    "action": "test_facebook_api",
    "phone_number_id": "YOUR_PHONE_NUMBER_ID",
    "access_token": "YOUR_ACCESS_TOKEN",
    "recipient": "RECIPIENT_PHONE_NUMBER",
    "template_name": "hello_world",
    "language_code": "en_US"
  }' \
  -w "\nHTTP Status: %{http_code}\n"

echo ""
echo "2️⃣ Testing WhatsApp status..."
curl -X GET http://localhost:3001/api/whatsapp/status

echo ""
echo "3️⃣ Testing admin login..."
curl -X POST http://localhost:3001/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "picoclaw123"
  }' \
  -w "\nHTTP Status: %{http_code}\n"

echo ""
echo "✅ Test completed! Check the server logs for detailed results."