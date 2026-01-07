#!/bin/bash

# Configuration
API_URL="http://localhost:8080"
EMAIL="monitor-test@example.com"
PASSWORD="SecurePassword123!"

echo "--- Generating Activity for Monitoring Demo ---"

# 1. Register User
echo "1. Registering User..."
AUTH_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\", \"name\": \"Monitor User\"}")

# Extract Token (API Key) - in this MVP register returns the key? No, Login does.
# Let's Login to get the key.
echo "2. Logging In..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")

# Extract API Key (assuming simple JSON response with "token" or similar, or checking implementation)
# Based on auth.go: using `return user, key.Key, nil` -> likely JSON response with APIKey if using standard handlers.
# Let's check the handler... wait, I can just use a known key if I had one, or parse.
# For simplicity, I will use 'grep' or 'sed' to extract the key if it's in the body.
# Actually, the handler returns `c.JSON(http.StatusOK, gin.H{"user": user, "api_key": key})`
API_KEY=$(echo $LOGIN_RESPONSE | grep -o '"api_key":"[^"]*' | cut -d'"' -f4)

if [ -z "$API_KEY" ]; then
    echo "Failed to get API Key. Login Response: $LOGIN_RESPONSE"
    exit 1
fi
echo "   Got API Key: $API_KEY"

# 2. Launch Instances (3 of them)
echo "3. Launching 3 Instances..."
for i in {1..3}; do
    curl -s -X POST "$API_URL/instances" \
      -H "X-API-Key: $API_KEY" \
      -H "Content-Type: application/json" \
      -d "{\"name\": \"monitor-instance-$i\", \"image\": \"alpine:latest\", \"memory\": 512, \"cpu\": 1}" > /dev/null
    echo "   Launched monitor-instance-$i"
    sleep 1 # Stagger slightly
done

# 3. Stop one instance
echo "4. Stopping 1 Instance..."
# Get ID of first instance
INSTANCES=$(curl -s -H "X-API-Key: $API_KEY" "$API_URL/instances")
FIRST_ID=$(echo $INSTANCES | grep -o '"id":"[^"]*' | head -n 1 | cut -d'"' -f4)

if [ ! -z "$FIRST_ID" ]; then
    curl -s -X POST "$API_URL/instances/$FIRST_ID/stop" -H "X-API-Key: $API_KEY" > /dev/null
    echo "   Stopped instance $FIRST_ID"
fi

# 4. Create a Volume
echo "5. Creating a Volume..."
curl -s -X POST "$API_URL/volumes" \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"monitor-vol-1\", \"size\": 10}" > /dev/null
echo "   Volume created."

echo "--- Activity Generation Complete ---"
echo "Check Grafana in 30 seconds (scrape interval)."
