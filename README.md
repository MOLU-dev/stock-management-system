# stock-management-system
# Stock Management System API Handlers

## Overview
This package contains HTTP handlers for a comprehensive stock management system, built with Go and PostgreSQL. The handlers provide RESTful APIs for managing inventory, products, purchase orders, stock movements, stocktakes, transfers, warehouses, and more.

## Handler Categories

### 1. Inventory Handler (`inventory.go`)
Manages inventory items including quantity tracking, reservations, and status updates.

**Key Endpoints:**
- `GET /inventory` - List all inventory (placeholder)
- `GET /inventory/{id}` - Get specific inventory item
- `GET /inventory/product/{productId}/warehouse/{warehouseId}` - Get inventory by product and warehouse
- `GET /inventory/warehouse/{warehouseId}` - List inventory by warehouse
- `GET /inventory/product/{productId}` - List inventory by product
- `GET /inventory/expiring` - List expiring inventory
- `PUT /inventory/{id}/quantity` - Update inventory quantity
- `POST /inventory/{id}/reserve` - Reserve inventory
- `POST /inventory/{id}/release` - Release reserved inventory
- `PUT /inventory/{id}/status` - Update inventory status

### 2. Product Handler (`product.go`)
Manages product catalog and product information.

**Key Endpoints:**
- `GET /products` - List products with pagination
- `GET /products/{id}` - Get product by ID
- `GET /products/sku/{sku}` - Get product by SKU
- `POST /products` - Create new product
- `PUT /products/{id}` - Update product
- `DELETE /products/{id}` - Soft delete product
- `GET /products/category/{categoryId}` - List products by category
- `GET /products/below-reorder-point` - List products below reorder point

### 3. Purchase Order Handler (`purchase_order.go`)
Handles purchase order creation, management, and item receiving.

**Key Endpoints:**
- `POST /purchase-orders` - Create purchase order
- `GET /purchase-orders/{id}` - Get purchase order
- `GET /purchase-orders` - List purchase orders
- `GET /purchase-orders/status/{status}` - List by status
- `PUT /purchase-orders/{id}/status` - Update status
- `GET /purchase-orders/{id}/items` - Get PO items
- `POST /purchase-orders/{id}/items` - Create PO item
- `POST /purchase-orders/items/{itemId}/receive` - Receive PO item

### 4. Stock Adjustment Handler (`stock_adjustment.go`)
Manages stock adjustments for inventory corrections.

**Key Endpoints:**
- `POST /stock-adjustments` - Create adjustment
- `POST /stock-adjustments/{id}/approve` - Approve adjustment
- `POST /stock-adjustments/{id}/items` - Create adjustment item

### 5. Stock Movement Handler (`stock_movement.go`)
Tracks all stock movements and history.

**Key Endpoints:**
- `POST /stock-movements` - Create movement record
- `GET /stock-movements/{id}` - Get movement
- `GET /stock-movements/product/{productId}` - List by product
- `GET /stock-movements/warehouse/{warehouseId}` - List by warehouse
- `GET /stock-movements/type/{type}` - List by type
- `GET /stock-movements/product/{productId}/warehouse/{warehouseId}/history` - Get product history

### 6. Stocktake Handler (`stocktake.go`)
Handles physical stock counting and variance tracking.

**Key Endpoints:**
- `POST /stocktakes` - Create stocktake
- `GET /stocktakes/{id}` - Get stocktake
- `GET /stocktakes` - List stocktakes
- `GET /stocktakes/warehouse/{warehouseId}` - List by warehouse
- `PUT /stocktakes/{id}/status` - Update status
- `POST /stocktakes/{id}/items` - Create stocktake item
- `GET /stocktakes/{id}/items` - Get stocktake items
- `PUT /stocktakes/items/{itemId}/count` - Update item count
- `GET /stocktakes/{id}/variances` - Get variances
- `GET /stocktakes/active` - Get active stocktakes

### 7. Transfer Handler (`transfer.go`)
Manages stock transfers between warehouses.

**Key Endpoints:**
- `POST /transfers` - Create transfer
- `PUT /transfers/{id}/status` - Update transfer status
- `POST /transfers/{id}/items` - Create transfer item
- `PUT /transfers/items/{itemId}/quantities` - Update transfer quantities

### 8. Warehouse Handler (`warehouse.go`)
Manages warehouses and storage locations.

**Key Endpoints:**
- `POST /warehouses` - Create warehouse
- `GET /warehouses/{id}` - Get warehouse
- `GET /warehouses/code/{code}` - Get by code
- `GET /warehouses` - List warehouses
- `PUT /warehouses/{id}` - Update warehouse
- `DELETE /warehouses/{id}` - Deactivate warehouse
- `POST /warehouses/{id}/locations` - Create location
- `GET /locations/{id}` - Get location
- `GET /warehouses/{id}/locations` - List locations
- `PUT /locations/{id}` - Update location
- `DELETE /locations/{id}` - Deactivate location
- `GET /warehouses/inventory-summary` - Get inventory summary

## Utility Functions

The package includes several helper functions for type conversion:

### Type Conversion Helpers:
- `NullString(s *string) sql.NullString` - Convert string pointer to SQL nullable string
- `toNullString(s *string) sql.NullString` - Convert string pointer to nullable string
- `toNullInt32FromValue(i int64) sql.NullInt32` - Convert int64 to nullable int32
- `toNullInt32FromInt32(i *int32) sql.NullInt32` - Convert int32 pointer to nullable int32
- `toNullInt32FromInt64(i *int64) sql.NullInt32` - Convert int64 pointer to nullable int32
- `toTimeOrZero(t *time.Time) time.Time` - Convert time pointer to time.Time (zero if nil)

### Response Helpers:
- `respondJSON(w http.ResponseWriter, status int, data interface{})` - Send JSON response
- `respondError(w http.ResponseWriter, status int, message string)` - Send error response

## Data Models

### Product Request/Response:
```go
type CreateProductRequest struct {
    SKU           string           `json:"sku"`
    Name          string           `json:"name"`
    Description   *string          `json:"description"`
    CategoryID    int64            `json:"category_id"`
    UnitPrice     decimal.Decimal  `json:"unit_price"`
    CostPrice     decimal.Decimal  `json:"cost_price"`
    // ... other fields
}
```

### Inventory Management:
```go
type UpdateQuantityRequest struct {
    Quantity         int32 `json:"quantity"`
    ReservedQuantity int32 `json:"reserved_quantity"`
}

type ReserveRequest struct {
    Quantity int32 `json:"quantity"`
}

type UpdateStatusRequest struct {
    Status string `json:"status"`
}
```

### Purchase Order:
```go
type CreatePurchaseOrderRequest struct {
    PONumber             string          `json:"po_number"`
    SupplierID           int64           `json:"supplier_id"`
    OrderDate            time.Time       `json:"order_date"`
    ExpectedDeliveryDate time.Time       `json:"expected_delivery_date"`
    Status               string          `json:"status"`
    TotalAmount          decimal.Decimal `json:"total_amount"`
    // ... other fields
}
```

## Database Integration

All handlers use the `db.SingleDb` interface from the internal database layer (`github.com/molu/stock-management-system/internal/db/sqlc`). This provides a clean separation between HTTP handling and data access.

## Dependencies

- **Database**: PostgreSQL with `sqlc` for type-safe queries
- **Routing**: `github.com/gorilla/mux` for HTTP routing
- **Decimal**: `github.com/shopspring/decimal` for precise monetary calculations
- **Standard Library**: `database/sql`, `encoding/json`, `net/http`, `time`

## Error Handling

All endpoints follow consistent error handling patterns:
1. Validation errors return `400 Bad Request`
2. Not found errors return `404 Not Found`
3. Database/processing errors return `500 Internal Server Error`
4. Success responses return appropriate status codes (`200 OK`, `201 Created`, `204 No Content`)

## Pagination

List endpoints support pagination via query parameters:
- `limit` - Number of items per page (default: 50)
- `offset` - Pagination offset (default: 0)

Example: `GET /products?limit=20&offset=40`

## Status Enums

The system uses several status enums defined in the database:
- `InventoryStatus` - For inventory items
- `PurchaseOrderStatus` - For purchase orders
- `AdjustmentStatus` - For stock adjustments
- `TransferStatus` - For stock transfers
- `StocktakeStatus` - For stocktakes

## Setup

1. Initialize handlers with database connection:
```go
queries := db.New(dbConn)
inventoryHandler := handlers.NewInventoryHandler(queries)
productHandler := handlers.NewProductHandler(queries)
// ... initialize other handlers
```

2. Register routes with your router:
```go
router.HandleFunc("/inventory", inventoryHandler.List).Methods("GET")
router.HandleFunc("/inventory/{id}", inventoryHandler.Get).Methods("GET")
// ... register other routes
```

## Notes

- All ID parameters are expected as integers
- Decimal values use `shopspring/decimal` for precision
- Soft delete is implemented for products
- Inventory reservations use optimistic locking
- Stock movements provide audit trail for all inventory changes
- Warehouse locations support hierarchical storage (aisle/shelf/bin)