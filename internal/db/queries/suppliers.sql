-- name: CreateSupplier :one
INSERT INTO suppliers (
    code, name, contact_person, email, phone, address, 
    tax_id, payment_terms, lead_time_days, rating, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetSupplier :one
SELECT * FROM suppliers 
WHERE supplier_id = $1;

-- name: GetSupplierByCode :one
SELECT * FROM suppliers 
WHERE code = $1;

-- name: ListSuppliers :many
SELECT * FROM suppliers 
WHERE is_active = true
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: ListAllSuppliers :many
SELECT * FROM suppliers 
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdateSupplier :one
UPDATE suppliers 
SET 
    name = COALESCE($2, name),
    contact_person = COALESCE($3, contact_person),
    email = COALESCE($4, email),
    phone = COALESCE($5, phone),
    address = COALESCE($6, address),
    tax_id = COALESCE($7, tax_id),
    payment_terms = COALESCE($8, payment_terms),
    lead_time_days = COALESCE($9, lead_time_days),
    rating = COALESCE($10, rating),
    is_active = COALESCE($11, is_active)
WHERE supplier_id = $1
RETURNING *;

-- name: DeactivateSupplier :exec
UPDATE suppliers 
SET is_active = false 
WHERE supplier_id = $1;

-- name: ActivateSupplier :exec
UPDATE suppliers 
SET is_active = true 
WHERE supplier_id = $1;

-- name: SearchSuppliers :many
SELECT * FROM suppliers 
WHERE is_active = true 
    AND (
        name ILIKE '%' || $1 || '%' 
        OR code ILIKE '%' || $1 || '%' 
        OR contact_person ILIKE '%' || $1 || '%'
    )
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: ListActiveSuppliers :many
SELECT * FROM suppliers 
WHERE is_active = true 
ORDER BY name;

-- name: GetSupplierProducts :many
SELECT p.* 
FROM products p
INNER JOIN product_suppliers ps ON p.product_id = ps.product_id
WHERE ps.supplier_id = $1 
    AND ps.is_active = true
ORDER BY p.name
LIMIT $2 OFFSET $3;

-- name: GetSupplierPerformance :one
SELECT 
    COUNT(DISTINCT po.po_id) as total_orders,
    COUNT(DISTINCT po.product_id) as unique_products,
    AVG(po.unit_price) as avg_unit_price,
    MAX(po.order_date) as last_order_date
FROM purchase_order_items po
INNER JOIN purchase_orders p ON po.po_id = p.po_id
WHERE p.supplier_id = $1;