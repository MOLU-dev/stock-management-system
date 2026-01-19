-- name: CreateStocktake :one
INSERT INTO stock_takes (
  stocktake_number, warehouse_id, start_date, end_date, status, notes, created_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetStocktake :one
SELECT st.*, w.name as warehouse_name
FROM stock_takes st
JOIN warehouses w ON st.warehouse_id = w.warehouse_id
WHERE st.stocktake_id = $1;

-- name: ListStocktakes :many
SELECT st.*, w.name as warehouse_name
FROM stock_takes st
JOIN warehouses w ON st.warehouse_id = w.warehouse_id
ORDER BY st.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListStocktakesByWarehouse :many
SELECT * FROM stock_takes
WHERE warehouse_id = $1
ORDER BY created_at DESC;

-- name: UpdateStocktakeStatus :one
UPDATE stock_takes
SET status = $2
WHERE stocktake_id = $1
RETURNING *;

-- name: CreateStocktakeItem :one
INSERT INTO stocktake_items (
  stocktake_id, product_id, location_id, system_quantity,
  counted_quantity, variance, counted_by, counted_at, notes
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetStocktakeItems :many
SELECT si.*, p.name as product_name, p.sku, l.location_code
FROM stocktake_items si
JOIN products p ON si.product_id = p.product_id
LEFT JOIN locations l ON si.location_id = l.location_id
WHERE si.stocktake_id = $1
ORDER BY p.name;

-- name: UpdateStocktakeItemCount :one
UPDATE stocktake_items
SET counted_quantity = $2,
    variance = $2 - system_quantity,
    counted_by = $3,
    counted_at = CURRENT_TIMESTAMP
WHERE stocktake_item_id = $1
RETURNING *;

-- name: GetStocktakeVariances :many
SELECT si.*, p.name as product_name, p.sku, l.location_code
FROM stocktake_items si
JOIN products p ON si.product_id = p.product_id
LEFT JOIN locations l ON si.location_id = l.location_id
WHERE si.stocktake_id = $1
  AND si.variance != 0
ORDER BY ABS(si.variance) DESC;

-- name: GetActiveStocktakes :many
SELECT st.*, w.name as warehouse_name
FROM stock_takes st
JOIN warehouses w ON st.warehouse_id = w.warehouse_id
WHERE st.status IN ('planned', 'in_progress')
ORDER BY st.start_date ASC;