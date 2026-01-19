-- name: GetProduct :one
SELECT * FROM products 
WHERE product_id = $1;

-- name: GetProductBySKU :one
SELECT * FROM products 
WHERE sku = $1;

-- name: ListProducts :many
SELECT * FROM products 
WHERE is_active = true
ORDER BY product_id
LIMIT $1 OFFSET $2;

-- name: ListProductsByCategory :many
SELECT * FROM products 
WHERE category_id = $1 AND is_active = true
ORDER BY product_id
LIMIT $2 OFFSET $3;

-- name: ListProductsBelowReorderPoint :many
SELECT p.*,
       COALESCE(SUM(i.quantity - i.reserved_quantity), 0) as available_qty
FROM products p
LEFT JOIN inventory i ON p.product_id = i.product_id
WHERE p.is_active = true
GROUP BY p.product_id
HAVING COALESCE(SUM(i.quantity - i.reserved_quantity), 0) <= p.reorder_point
ORDER BY p.product_id;

-- name: CreateProduct :one
INSERT INTO products (
    sku, name, description, category_id, unit_price, cost_price,
    barcode, weight, dimensions, supplier_id, min_stock_level,
    max_stock_level, reorder_point, safety_stock, lead_time_days,
    auto_reorder, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
) RETURNING *;

-- name: UpdateProduct :one
UPDATE products 
SET 
    name = $2,
    description = $3,
    category_id = $4,
    unit_price = $5,
    cost_price = $6,
    barcode = $7,
    weight = $8,
    dimensions = $9,
    supplier_id = $10,
    min_stock_level = $11,
    max_stock_level = $12,
    reorder_point = $13,
    safety_stock = $14,
    lead_time_days = $15,
    auto_reorder = $16,
    is_active = $17,
    updated_at = CURRENT_TIMESTAMP
WHERE product_id = $1
RETURNING *;

-- name: SoftDeleteProduct :exec
UPDATE products 
SET is_active = false, updated_at = CURRENT_TIMESTAMP
WHERE product_id = $1;

-- name: UpdateLastReorderDate :exec
UPDATE products 
SET last_reorder_date = CURRENT_DATE, updated_at = CURRENT_TIMESTAMP
WHERE product_id = $1;