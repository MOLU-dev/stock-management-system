-- name: CreateStockAdjustment :one
INSERT INTO stock_adjustments (
    adjustment_number, warehouse_id, adjustment_date,
    reason, status, total_value, notes, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: ApproveStockAdjustment :one
UPDATE stock_adjustments 
SET 
    status = 'approved',
    approved_by = $2,
    approved_at = CURRENT_TIMESTAMP
WHERE adjustment_id = $1 AND status = 'pending'
RETURNING *;

-- name: CreateStockAdjustmentItem :one
INSERT INTO stock_adjustment_items (
    adjustment_id, product_id, quantity_before,
    quantity_adjusted, cost_price, reason
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING adjustment_item_id,
    quantity_before + quantity_adjusted as quantity_after,
    quantity_adjusted * cost_price as adjustment_value;