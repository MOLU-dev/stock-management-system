#!/bin/bash

BASE_URL="http://localhost:8080/api/v1/categories"

echo "=== TESTING CATEGORY API ==="

# 1. Create a root category
echo "1. Creating root category..."
ROOT_CAT=$(curl -s -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "category_code": "CAT-TECH",
    "name": "Technology",
    "description": "All tech stuff"
  }')
echo "$ROOT_CAT" | jq .
ROOT_ID=$(echo "$ROOT_CAT" | jq .category_id)

# 2. Create a subcategory
echo -e "\n2. Creating subcategory..."
SUB_CAT=$(curl -s -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -d "{
    \"category_code\": \"CAT-LAPTOP\",
    \"name\": \"Laptops\",
    \"parent_category_id\": $ROOT_ID,
    \"description\": \"All portable computers\"
  }")
echo "$SUB_CAT" | jq .
SUB_ID=$(echo "$SUB_CAT" | jq .category_id)

# 3. List all categories
echo -e "\n3. Listing all categories..."
curl -s "$BASE_URL" | jq .

# 4. List root categories
echo -e "\n4. Listing root categories..."
curl -s "$BASE_URL/root" | jq .

# 5. List subcategories
echo -e "\n5. Listing subcategories for ID $ROOT_ID..."
curl -s "$BASE_URL/$ROOT_ID/subcategories" | jq .

# 6. Get category by ID
echo -e "\n6. Getting category ID $SUB_ID..."
curl -s "$BASE_URL/$SUB_ID" | jq .

# 7. Get category by code
echo -e "\n7. Getting category code CAT-TECH..."
curl -s "$BASE_URL/code/CAT-TECH" | jq .

# 8. Update category
echo -e "\n8. Updating category ID $SUB_ID..."
curl -s -X PUT "$BASE_URL/$SUB_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptops & Notebooks",
    "description": "Portables, notebooks, and ultrabooks"
  }' | jq .

# 9. Delete category (cleanup)
# echo -e "\n9. Deleting categories..."
# curl -s -X DELETE "$BASE_URL/$SUB_ID" | jq .
# curl -s -X DELETE "$BASE_URL/$ROOT_ID" | jq .

echo -e "\n=== CATEGORY API TEST COMPLETE ==="
