#!/bin/bash

# Stock Management System API Testing Commands
# Base URL
BASE_URL="http://localhost:8080/api/v1"

# ============================================================================
# DATABASE SETUP (Run this first if starting fresh)
# ============================================================================

# echo "=== DATABASE SETUP ==="
# echo "Note: Run these SQL commands first if starting with a fresh database:"
# echo ""
# echo "1. Insert categories:"
# echo "   INSERT INTO categories (category_code, name) VALUES ('CAT-001', 'Electronics');"
# echo ""
# echo "2. Insert users:"
# echo "   INSERT INTO users (username, email, password_hash, full_name, role)"
# echo "   VALUES ('testuser', 'test@example.com', 'hashedpassword', 'Test User', 'admin');"
# echo ""
# echo "Press any key to continue or Ctrl+C to abort..."
# read -n 1

# ============================================================================
# HEALTH CHECK
# ============================================================================

echo "=== HEALTH CHECK ==="
curl -X GET http://localhost:8080/health
echo -e "\n"

# ============================================================================
# WAREHOUSES (Create first for other operations)
# ============================================================================

echo "=== WAREHOUSES ==="

# List warehouses first to see what exists
echo "Listing warehouses..."
curl -X GET "${BASE_URL}/warehouses"
echo -e "\n"

# Create main warehouse (only if it doesn't exist)
echo "Creating main warehouse (if not exists)..."
curl -X POST "${BASE_URL}/warehouses" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "WH-TEST-001",
    "name": "Test Main Warehouse",
    "address": "123 Test Street, Lagos",
    "contact_person": "Test Manager",
    "contact_phone": "+234-111-222-3333",
    "contact_email": "test@warehouse.com"
  }'
echo -e "\n"

# Create secondary warehouse
echo "Creating secondary warehouse..."
curl -X POST "${BASE_URL}/warehouses" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "WH-TEST-002",
    "name": "Test Secondary Warehouse",
    "address": "456 Test Avenue, Abuja",
    "contact_person": "Test Assistant",
    "contact_phone": "+234-444-555-6666",
    "contact_email": "test2@warehouse.com"
  }'
echo -e "\n"

# ============================================================================
# CATEGORIES (Need categories before products)
# ============================================================================

echo "=== CATEGORIES ==="

echo "Creating test category..."
curl -X POST "${BASE_URL}/categories" \
  -H "Content-Type: application/json" \
  -d '{
    "category_code": "CAT-TEST",
    "name": "Test Electronics",
    "description": "Test electronics category"
  }'
echo -e "\n"

# ============================================================================
# LOCATIONS (Create locations within warehouses)
# ============================================================================

echo "=== LOCATIONS ==="

# Get warehouse IDs first
echo "Getting warehouse IDs..."
WAREHOUSE1_ID=1
WAREHOUSE2_ID=2

# Create locations in main warehouse
echo "Creating locations in warehouse $WAREHOUSE1_ID..."

curl -X POST "${BASE_URL}/warehouses/$WAREHOUSE1_ID/locations" \
  -H "Content-Type: application/json" \
  -d '{
    "location_code": "TEST-A-01-01",
    "aisle": "A",
    "shelf": "01",
    "bin": "01",
    "max_capacity": 100
  }'
echo -e "\n"

curl -X POST "${BASE_URL}/warehouses/$WAREHOUSE1_ID/locations" \
  -H "Content-Type: application/json" \
  -d '{
    "location_code": "TEST-A-01-02",
    "aisle": "A",
    "shelf": "01",
    "bin": "02",
    "max_capacity": 100
  }'
echo -e "\n"

# Create locations in secondary warehouse
curl -X POST "${BASE_URL}/warehouses/$WAREHOUSE2_ID/locations" \
  -H "Content-Type: application/json" \
  -d '{
    "location_code": "TEST-B-01-01",
    "aisle": "B",
    "shelf": "01",
    "bin": "01",
    "max_capacity": 80
  }'
echo -e "\n"

# List locations by warehouse
echo "Listing locations in warehouse $WAREHOUSE1_ID..."
curl -X GET "${BASE_URL}/warehouses/$WAREHOUSE1_ID/locations"
echo -e "\n"

# ============================================================================
# SUPPLIERS (Create before products)
# ============================================================================

echo "=== SUPPLIERS ==="

# Create first supplier
echo "Creating first supplier..."
curl -X POST "${BASE_URL}/suppliers" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "SUP-TEST-001",
    "name": "Test ElectroParts Inc.",
    "contact_person": "Test Supplier",
    "email": "test@electroparts.com",
    "phone": "+234-777-888-9999",
    "address": "123 Test Supplier Street",
    "tax_id": "TAX-TEST-001",
    "payment_terms": "Net 30",
    "lead_time_days": 7,
    "rating": 4.5,
    "is_active": true
  }'
echo -e "\n"

# Create second supplier
echo "Creating second supplier..."
curl -X POST "${BASE_URL}/suppliers" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "SUP-TEST-002",
    "name": "Test Global Components Ltd.",
    "contact_person": "Test Wholesaler",
    "email": "test@globalcomponents.com",
    "phone": "+234-999-888-7777",
    "address": "456 Test Industrial Area",
    "tax_id": "TAX-TEST-002",
    "payment_terms": "Net 45",
    "lead_time_days": 14,
    "rating": 4.2,
    "is_active": true
  }'
echo -e "\n"

# List suppliers
echo "Listing all suppliers..."
curl -X GET "${BASE_URL}/suppliers"
echo -e "\n"

# ============================================================================
# PRODUCTS (Require suppliers and categories)
# ============================================================================

echo "=== PRODUCTS ==="

# Get category ID (assuming first category is ID 1)
CATEGORY_ID=1

# Create first product
echo "Creating first product..."
curl -X POST "${BASE_URL}/products" \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "TEST-ELEC-001",
    "name": "Test Smartphone X",
    "description": "Test smartphone with advanced features",
    "category_id": 1,
    "unit_price": "599.99",
    "cost_price": "350.00",
    "barcode": "TEST-1234567890123",
    "weight": "0.3",
    "dimensions": "15x7x0.8",
    "supplier_id": 1,
    "min_stock_level": 10,
    "max_stock_level": 200,
    "reorder_point": 25,
    "safety_stock": 15,
    "lead_time_days": 7,
    "auto_reorder": true,
    "is_active": true
  }'
echo -e "\n"

# Create second product
echo "Creating second product..."
curl -X POST "${BASE_URL}/products" \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "TEST-ELEC-002",
    "name": "Test Laptop Pro",
    "description": "Test high-performance laptop",
    "category_id": 1,
    "unit_price": "1299.99",
    "cost_price": "800.00",
    "barcode": "TEST-9876543210987",
    "weight": "1.8",
    "dimensions": "35x24x2",
    "supplier_id": 2,
    "min_stock_level": 5,
    "max_stock_level": 50,
    "reorder_point": 10,
    "safety_stock": 5,
    "lead_time_days": 14,
    "auto_reorder": true,
    "is_active": true
  }'
echo -e "\n"

# Create third product
echo "Creating third product..."
curl -X POST "${BASE_URL}/products" \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "TEST-ACC-001",
    "name": "Test Wireless Headphones",
    "description": "Test noise-cancelling headphones",
    "category_id": 1,
    "unit_price": "199.99",
    "cost_price": "120.00",
    "barcode": "TEST-4567891234567",
    "weight": "0.25",
    "dimensions": "18x15x8",
    "supplier_id": 1,
    "min_stock_level": 20,
    "max_stock_level": 150,
    "reorder_point": 30,
    "safety_stock": 20,
    "lead_time_days": 5,
    "auto_reorder": true,
    "is_active": true
  }'
echo -e "\n"

# List products
echo "Listing products..."
curl -X GET "${BASE_URL}/products?limit=50&offset=0"
echo -e "\n"

# ============================================================================
# INVENTORY
# ============================================================================

echo "=== INVENTORY ==="

# List all inventory
echo "Listing all inventory..."
curl -X GET "${BASE_URL}/inventory"
echo -e "\n"

# ============================================================================
# PURCHASE ORDERS
# ============================================================================

echo "=== PURCHASE ORDERS ==="

# Create purchase order for supplier 1
echo "Creating purchase order..."
curl -X POST "${BASE_URL}/purchase-orders" \
  -H "Content-Type: application/json" \
  -d '{
    "po_number": "PO-TEST-001",
    "supplier_id": 1,
    "order_date": "2026-01-19T10:00:00Z",
    "expected_delivery_date": "2026-01-26T10:00:00Z",
    "status": "pending",
    "total_amount": "1000.00",
    "notes": "Test urgent order",
    "created_by": 1
  }'
echo -e "\n"

# Create another purchase order
echo "Creating second purchase order..."
curl -X POST "${BASE_URL}/purchase-orders" \
  -H "Content-Type: application/json" \
  -d '{
    "po_number": "PO-TEST-002",
    "supplier_id": 2,
    "order_date": "2026-01-19T11:00:00Z",
    "expected_delivery_date": "2026-02-02T10:00:00Z",
    "status": "pending",
    "total_amount": "2500.00",
    "notes": "Test regular order",
    "created_by": 1
  }'
echo -e "\n"

# List purchase orders
echo "Listing purchase orders..."
curl -X GET "${BASE_URL}/purchase-orders"
echo -e "\n"

# ============================================================================
# PURCHASE ORDER ITEMS
# ============================================================================

echo "=== PURCHASE ORDER ITEMS ==="

# Get latest PO IDs by checking the response above or assuming sequential
PO1_ID=1
PO2_ID=2

# Add items to first purchase order
echo "Adding items to purchase order $PO1_ID..."

curl -X POST "${BASE_URL}/purchase-orders/$PO1_ID/items" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity_ordered": 50,
    "quantity_received": 0,
    "unit_price": "350.00",
    "total_price": "17500.00"
  }'
echo -e "\n"

curl -X POST "${BASE_URL}/purchase-orders/$PO1_ID/items" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 3,
    "quantity_ordered": 100,
    "quantity_received": 0,
    "unit_price": "120.00",
    "total_price": "12000.00"
  }'
echo -e "\n"

# Add items to second purchase order
echo "Adding items to purchase order $PO2_ID..."

curl -X POST "${BASE_URL}/purchase-orders/$PO2_ID/items" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 2,
    "quantity_ordered": 20,
    "quantity_received": 0,
    "unit_price": "800.00",
    "total_price": "16000.00"
  }'
echo -e "\n"

# Get purchase order items
echo "Getting items for purchase order $PO1_ID..."
curl -X GET "${BASE_URL}/purchase-orders/$PO1_ID/items"
echo -e "\n"

# ============================================================================
# RECEIVE PURCHASE ORDER ITEMS (Creates inventory)
# ============================================================================

echo "=== RECEIVING PURCHASE ORDER ITEMS ==="

# Assuming PO item IDs are 1, 2, 3 based on creation order

# Receive items from first PO
echo "Receiving items from PO 1..."

curl -X POST "${BASE_URL}/purchase-orders/items/1/receive" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 50,
    "location_id": 1
  }'
echo -e "\n"

curl -X POST "${BASE_URL}/purchase-orders/items/2/receive" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 100,
    "location_id": 2
  }'
echo -e "\n"

# Receive partial shipment from second PO
echo "Receiving partial shipment from PO 2..."

curl -X POST "${BASE_URL}/purchase-orders/items/3/receive" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 10,
    "location_id": 3
  }'
echo -e "\n"

# ============================================================================
# INVENTORY CHECK (After receiving items)
# ============================================================================

echo "=== INVENTORY CHECK ==="

# List all inventory
echo "Listing all inventory..."
curl -X GET "${BASE_URL}/inventory"
echo -e "\n"

# List inventory by warehouse
echo "Listing inventory in warehouse 1..."
curl -X GET "${BASE_URL}/inventory/warehouse/1"
echo -e "\n"

# ============================================================================
# STOCK MOVEMENTS
# ============================================================================

echo "=== STOCK MOVEMENTS ==="

# List stock movements
echo "Listing stock movements..."
curl -X GET "${BASE_URL}/stock-movements"
echo -e "\n"

# ============================================================================
# STOCK ADJUSTMENTS
# ============================================================================

echo "=== STOCK ADJUSTMENTS ==="

# Create stock adjustment
echo "Creating stock adjustment..."
curl -X POST "${BASE_URL}/stock-adjustments" \
  -H "Content-Type: application/json" \
  -d '{
    "adjustment_number": "ADJ-TEST-001",
    "warehouse_id": 1,
    "adjustment_date": "2026-01-19T14:00:00Z",
    "reason": "damage",
    "status": "pending",
    "notes": "Test damaged items",
    "created_by": 1
  }'
echo -e "\n"

# Create adjustment item
echo "Creating adjustment item..."
curl -X POST "${BASE_URL}/stock-adjustments/1/items" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity_before": 50,
    "quantity_adjusted": -3,
    "cost_price": "350.00",
    "reason": "Test screen damage"
  }'
echo -e "\n"

# ============================================================================
# STOCK TRANSFERS
# ============================================================================

echo "=== STOCK TRANSFERS ==="

# Create stock transfer between warehouses
echo "Creating stock transfer..."
curl -X POST "${BASE_URL}/stock-transfers" \
  -H "Content-Type: application/json" \
  -d '{
    "transfer_number": "TRF-TEST-001",
    "from_warehouse_id": 1,
    "to_warehouse_id": 2,
    "status": "pending",
    "transfer_date": "2026-01-19T15:00:00Z",
    "expected_completion_date": "2026-01-21T10:00:00Z",
    "notes": "Test stock rebalancing",
    "created_by": 1
  }'
echo -e "\n"

# ============================================================================
# STOCKTAKES
# ============================================================================

echo "=== STOCKTAKES ==="

# Create stocktake
echo "Creating stocktake..."
curl -X POST "${BASE_URL}/stocktakes" \
  -H "Content-Type: application/json" \
  -d '{
    "stocktake_number": "ST-TEST-001",
    "warehouse_id": 1,
    "start_date": "2026-01-19T16:00:00Z",
    "status": "in_progress",
    "notes": "Test monthly stocktake",
    "created_by": 1
  }'
echo -e "\n"

# ============================================================================
# SUMMARY AND CLEANUP
# ============================================================================

echo "=== TEST SUMMARY ==="
echo ""
echo "Test operations completed."
echo "Check the responses above for any errors."
echo ""
echo "If you see 'Failed to create' errors, the records may already exist."
echo "If you see SQL errors, check your database schema and data."
echo ""
echo "To clean up test data, you can run:"
echo "DELETE FROM stocktakes WHERE stocktake_number LIKE 'ST-TEST%';"
echo "DELETE FROM stock_transfers WHERE transfer_number LIKE 'TRF-TEST%';"
echo "DELETE FROM stock_adjustments WHERE adjustment_number LIKE 'ADJ-TEST%';"
echo "DELETE FROM purchase_order_items WHERE purchase_order_id IN (SELECT po_id FROM purchase_orders WHERE po_number LIKE 'PO-TEST%');"
echo "DELETE FROM purchase_orders WHERE po_number LIKE 'PO-TEST%';"
echo "DELETE FROM inventory_items WHERE product_id IN (SELECT product_id FROM products WHERE sku LIKE 'TEST-%');"
echo "DELETE FROM products WHERE sku LIKE 'TEST-%';"
echo "DELETE FROM suppliers WHERE code LIKE 'SUP-TEST%';"
echo "DELETE FROM locations WHERE location_code LIKE 'TEST-%';"
echo "DELETE FROM warehouses WHERE code LIKE 'WH-TEST%';"
echo "DELETE FROM categories WHERE category_code = 'CAT-TEST';"