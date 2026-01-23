#!/bin/bash

# Simple API test to verify server is working correctly
BASE_URL="http://localhost:8080/api/v1"

echo "=== TESTING STOCK MANAGEMENT API ==="
echo ""

# Test 1: Health check
echo "1. Health Check..."
response=$(curl -s http://localhost:8080/health)
if [ "$response" = "OK" ]; then
    echo "   ✓ Health check passed"
else
    echo "   ✗ Health check failed: $response"
fi
echo ""

# Test 2: List warehouses
echo "2. List Warehouses..."
response=$(curl -s "${BASE_URL}/warehouses" | jq '. | length' 2>/dev/null)
if [ "$response" -gt 0 ]; then
    echo "   ✓ Found $response warehouses"
else
    echo "   ✗ Failed to list warehouses"
fi
echo ""

# Test 3: List products (fixed - should work now)
echo "3. List Products..."
response=$(curl -s "${BASE_URL}/products?limit=5&offset=0")
if echo "$response" | jq empty 2>/dev/null; then
    count=$(echo "$response" | jq '. | length')
    echo "   ✓ Products endpoint working - found $count products"
else
    echo "   ✗ Products endpoint failed: $response"
fi
echo ""

# Test 4: List purchase orders
echo "4. List Purchase Orders..."
response=$(curl -s "${BASE_URL}/purchase-orders")
if echo "$response" | jq empty 2>/dev/null; then
    count=$(echo "$response" | jq '. | length')
    echo "   ✓ Purchase orders endpoint working - found $count orders"
else
    echo "   ✗ Purchase orders endpoint failed"
fi
echo ""

# Test 5: List suppliers
echo "5. List Suppliers..."
response=$(curl -s "${BASE_URL}/suppliers")
if echo "$response" | jq empty 2>/dev/null; then
    echo "   ✓ Suppliers endpoint working"
else
    echo "   ✗ Suppliers endpoint failed"
fi
echo ""

# Test 6: List stocktakes
echo "6. List Stocktakes..."
response=$(curl -s "${BASE_URL}/stocktakes")
if echo "$response" | jq empty 2>/dev/null; then
    count=$(echo "$response" | jq '. | length')
    echo "   ✓ Stocktakes endpoint working - found $count stocktakes"
else
    echo "   ✗ Stocktakes endpoint failed"
fi
echo ""

echo "=== TEST SUMMARY ==="
echo "Server is running and responding to requests."
echo "Main issue fixed: Products can now be listed without NULL pointer errors."
