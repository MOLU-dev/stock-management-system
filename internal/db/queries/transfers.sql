-- name: CreateStockTransfer :one
INSERT INTO stock_transfers (
    transfer_number, from_warehouse_id, to_warehouse_id,
    status, transfer_date, expected_completion_date,
    notes, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: UpdateStockTransferStatus :one
UPDATE stock_transfers 
SET 
    status = $2,
    transfer_date = CASE 
        WHEN $2 = 'completed' AND transfer_date IS NULL THEN CURRENT_DATE 
        ELSE transfer_date 
    END,
    updated_at = CURRENT_TIMESTAMP
WHERE transfer_id = $1
RETURNING *;

-- name: CreateStockTransferItem :one
INSERT INTO stock_transfer_items (
    transfer_id, product_id, quantity,
    from_location_id, to_location_id
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateStockTransferItemQuantities :one
UPDATE stock_transfer_items 
SET 
    quantity_sent = $2,
    quantity_received = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE transfer_item_id = $1
RETURNING *;