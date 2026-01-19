-- name: CreateWarehouse :one
INSERT INTO warehouses (
  code, name, address, contact_person, contact_phone, contact_email
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetWarehouse :one
SELECT * FROM warehouses WHERE warehouse_id = $1;

-- name: GetWarehouseByCode :one
SELECT * FROM warehouses WHERE code = $1;

-- name: ListWarehouses :many
SELECT * FROM warehouses
WHERE is_active = true
ORDER BY name;

-- name: ListAllWarehouses :many
SELECT * FROM warehouses
ORDER BY name;

-- name: UpdateWarehouse :one
UPDATE warehouses
SET name = $2,
    address = $3,
    contact_person = $4,
    contact_phone = $5,
    contact_email = $6
WHERE warehouse_id = $1
RETURNING *;

-- name: DeactivateWarehouse :exec
UPDATE warehouses
SET is_active = false
WHERE warehouse_id = $1;

-- name: CreateLocation :one
INSERT INTO locations (
  warehouse_id, location_code, aisle, shelf, bin, max_capacity
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetLocation :one
SELECT * FROM locations WHERE location_id = $1;

-- name: GetLocationByCode :one
SELECT * FROM locations
WHERE warehouse_id = $1 AND location_code = $2;

-- name: ListLocationsByWarehouse :many
SELECT * FROM locations
WHERE warehouse_id = $1 AND is_active = true
ORDER BY location_code;

-- name: UpdateLocation :one
UPDATE locations
SET aisle = $2,
    shelf = $3,
    bin = $4,
    max_capacity = $5
WHERE location_id = $1
RETURNING *;

-- name: DeactivateLocation :exec
UPDATE locations
SET is_active = false
WHERE location_id = $1;

-- name: GetWarehouseInventorySummary :many
SELECT 
  w.warehouse_id,
  w.name as warehouse_name,
  COUNT(DISTINCT i.product_id) as unique_products,
  COALESCE(SUM(i.quantity), 0) as total_items,
  COALESCE(SUM(i.reserved_quantity), 0) as reserved_items
FROM warehouses w
LEFT JOIN inventory i ON w.warehouse_id = i.warehouse_id
WHERE w.is_active = true
GROUP BY w.warehouse_id, w.name
ORDER BY w.name;