-- name: CreateStockMovement :one
INSERT INTO stock_movements (
    reference_number, product_id, warehouse_id, location_id,
    movement_type, quantity_before, quantity_change, quantity_after,
    reference_id, reference_table, notes, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: GetStockMovement :one
SELECT * FROM stock_movements 
WHERE movement_id = $1;

-- name: ListStockMovementsByProduct :many
SELECT sm.*, p.name as product_name, p.sku, w.name as warehouse_name
FROM stock_movements sm
JOIN products p ON sm.product_id = p.product_id
JOIN warehouses w ON sm.warehouse_id = w.warehouse_id
WHERE sm.product_id = $1
ORDER BY sm.movement_date DESC
LIMIT $2 OFFSET $3;

-- name: ListStockMovementsByWarehouse :many
SELECT sm.*, p.name as product_name, p.sku
FROM stock_movements sm
JOIN products p ON sm.product_id = p.product_id
WHERE sm.warehouse_id = $1
ORDER BY sm.movement_date DESC
LIMIT $2 OFFSET $3;

-- name: ListStockMovementsByType :many
SELECT sm.*, p.name as product_name, p.sku, w.name as warehouse_name
FROM stock_movements sm
JOIN products p ON sm.product_id = p.product_id
JOIN warehouses w ON sm.warehouse_id = w.warehouse_id
WHERE sm.movement_type = $1
ORDER BY sm.movement_date DESC
LIMIT $2 OFFSET $3;

-- name: GetProductMovementHistory :many
SELECT sm.*, w.name as warehouse_name
FROM stock_movements sm
JOIN warehouses w ON sm.warehouse_id = w.warehouse_id
WHERE sm.product_id = $1 AND sm.warehouse_id = $2
ORDER BY sm.movement_date DESC
LIMIT $3 OFFSET $4;