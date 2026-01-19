-- name: CreatePurchaseOrder :one
INSERT INTO purchase_orders (
    po_number, supplier_id, order_date, expected_delivery_date,
    status, total_amount, notes, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetPurchaseOrder :one
SELECT po.*, s.name as supplier_name, s.code as supplier_code,
       u.full_name as creator_name
FROM purchase_orders po
LEFT JOIN suppliers s ON po.supplier_id = s.supplier_id
LEFT JOIN users u ON po.created_by = u.user_id
WHERE po.po_id = $1;

-- name: ListPurchaseOrders :many
SELECT po.*, s.name as supplier_name
FROM purchase_orders po
LEFT JOIN suppliers s ON po.supplier_id = s.supplier_id
ORDER BY po.order_date DESC
LIMIT $1 OFFSET $2;

-- name: ListPurchaseOrdersByStatus :many
SELECT po.*, s.name as supplier_name
FROM purchase_orders po
LEFT JOIN suppliers s ON po.supplier_id = s.supplier_id
WHERE po.status = $1
ORDER BY po.order_date DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePurchaseOrderStatus :one
UPDATE purchase_orders 
SET 
    status = $2,
    total_amount = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE po_id = $1
RETURNING *;

-- name: CreatePurchaseOrderItem :one
INSERT INTO purchase_order_items (
    po_id, product_id, quantity_ordered, quantity_received,
    unit_price, total_price
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetPurchaseOrderItems :many
SELECT poi.*, p.name as product_name, p.sku
FROM purchase_order_items poi
JOIN products p ON poi.product_id = p.product_id
WHERE poi.po_id = $1
ORDER BY poi.po_item_id;

-- name: UpdatePurchaseOrderItemReceivedQty :one
UPDATE purchase_order_items 
SET 
    quantity_received = quantity_received + $2,
    updated_at = CURRENT_TIMESTAMP
WHERE po_item_id = $1
RETURNING *;