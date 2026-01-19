-- name: GetInventory :one
SELECT * FROM inventory 
WHERE inventory_id = $1;

-- name: GetInventoryByProductWarehouse :one
SELECT * FROM inventory 
WHERE product_id = $1 AND warehouse_id = $2;

-- name: GetInventoryByLocation :one
SELECT * FROM inventory 
WHERE product_id = $1 AND warehouse_id = $2 AND location_id = $3;

-- name: ListInventoryByWarehouse :many
SELECT i.*, p.name as product_name, p.sku
FROM inventory i
JOIN products p ON i.product_id = p.product_id
WHERE i.warehouse_id = $1 AND p.is_active = true
ORDER BY i.product_id
LIMIT $2 OFFSET $3;

-- name: ListInventoryByProduct :many
SELECT i.*, w.name as warehouse_name, w.code as warehouse_code
FROM inventory i
JOIN warehouses w ON i.warehouse_id = w.warehouse_id
WHERE i.product_id = $1 AND w.is_active = true
ORDER BY i.warehouse_id;

-- name: ListExpiringInventory :many
SELECT i.*, p.name as product_name, p.sku, 
       w.name as warehouse_name, w.code as warehouse_code
FROM inventory i
JOIN products p ON i.product_id = p.product_id
JOIN warehouses w ON i.warehouse_id = w.warehouse_id
WHERE i.expiry_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '30 days'
  AND i.status = 'in_stock'
  AND p.is_active = true
  AND w.is_active = true
ORDER BY i.expiry_date;

-- name: UpdateInventoryQuantity :one
UPDATE inventory 
SET 
    quantity = $2,
    reserved_quantity = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE inventory_id = $1
RETURNING *;

-- name: ReserveInventory :one
UPDATE inventory 
SET 
    reserved_quantity = reserved_quantity + $2,
    updated_at = CURRENT_TIMESTAMP
WHERE inventory_id = $1 AND (quantity - reserved_quantity) >= $2
RETURNING *;

-- name: ReleaseInventoryReservation :one
UPDATE inventory 
SET 
    reserved_quantity = reserved_quantity - $2,
    updated_at = CURRENT_TIMESTAMP
WHERE inventory_id = $1 AND reserved_quantity >= $2
RETURNING *;

-- name: UpdateInventoryStatus :one
UPDATE inventory 
SET 
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE inventory_id = $1
RETURNING *;