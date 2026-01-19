package router

import (
	"net/http"

	"github.com/gorilla/mux"
	db "github.com/molu/stock-management-system/internal/db/sqlc"
	"github.com/molu/stock-management-system/internal/handlers"
	"github.com/molu/stock-management-system/internal/middleware"
)

func New(queries *db.Queries, jwtSecret string) http.Handler {
	r := mux.NewRouter()

	// Initialize handlers
	productHandler := handlers.NewProductHandler(queries)
	inventoryHandler := handlers.NewInventoryHandler(queries)
	stockMovementHandler := handlers.NewStockMovementHandler(queries)
	purchaseOrderHandler := handlers.NewPurchaseOrderHandler(queries)
	stockAdjustmentHandler := handlers.NewStockAdjustmentHandler(queries)
	transferHandler := handlers.NewTransferHandler(queries)
	stocktakeHandler := handlers.NewStocktakeHandler(queries)
	warehouseHandler := handlers.NewWarehouseHandler(queries)
	supplierHandler := handlers.NewSupplierHandler(queries) // Add this line

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.CORS)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Products
	products := api.PathPrefix("/products").Subrouter()
	products.HandleFunc("", productHandler.List).Methods("GET")
	products.HandleFunc("", productHandler.Create).Methods("POST")
	products.HandleFunc("/{id}", productHandler.Get).Methods("GET")
	products.HandleFunc("/{id}", productHandler.Update).Methods("PUT")
	products.HandleFunc("/{id}", productHandler.Delete).Methods("DELETE")
	products.HandleFunc("/sku/{sku}", productHandler.GetBySKU).Methods("GET")
	products.HandleFunc("/category/{categoryId}", productHandler.ListByCategory).Methods("GET")
	products.HandleFunc("/reorder/below-point", productHandler.ListBelowReorderPoint).Methods("GET")

	// Inventory
	inventory := api.PathPrefix("/inventory").Subrouter()
	inventory.HandleFunc("", inventoryHandler.List).Methods("GET")
	inventory.HandleFunc("/{id}", inventoryHandler.Get).Methods("GET")
	inventory.HandleFunc("/product/{productId}/warehouse/{warehouseId}", inventoryHandler.GetByProductWarehouse).Methods("GET")
	inventory.HandleFunc("/warehouse/{warehouseId}", inventoryHandler.ListByWarehouse).Methods("GET")
	inventory.HandleFunc("/product/{productId}", inventoryHandler.ListByProduct).Methods("GET")
	inventory.HandleFunc("/expiring", inventoryHandler.ListExpiring).Methods("GET")
	inventory.HandleFunc("/{id}/quantity", inventoryHandler.UpdateQuantity).Methods("PUT")
	inventory.HandleFunc("/{id}/reserve", inventoryHandler.Reserve).Methods("POST")
	inventory.HandleFunc("/{id}/release", inventoryHandler.Release).Methods("POST")
	inventory.HandleFunc("/{id}/status", inventoryHandler.UpdateStatus).Methods("PUT")

	// Stock Movements
	movements := api.PathPrefix("/stock-movements").Subrouter()
	movements.HandleFunc("", stockMovementHandler.Create).Methods("POST")
	movements.HandleFunc("/{id}", stockMovementHandler.Get).Methods("GET")
	movements.HandleFunc("/product/{productId}", stockMovementHandler.ListByProduct).Methods("GET")
	movements.HandleFunc("/warehouse/{warehouseId}", stockMovementHandler.ListByWarehouse).Methods("GET")
	movements.HandleFunc("/type/{type}", stockMovementHandler.ListByType).Methods("GET")
	movements.HandleFunc("/history/product/{productId}/warehouse/{warehouseId}", stockMovementHandler.GetProductHistory).Methods("GET")

	// Purchase Orders
	purchaseOrders := api.PathPrefix("/purchase-orders").Subrouter()
	purchaseOrders.HandleFunc("", purchaseOrderHandler.List).Methods("GET")
	purchaseOrders.HandleFunc("", purchaseOrderHandler.Create).Methods("POST")
	purchaseOrders.HandleFunc("/{id}", purchaseOrderHandler.Get).Methods("GET")
	purchaseOrders.HandleFunc("/{id}/status", purchaseOrderHandler.UpdateStatus).Methods("PUT")
	purchaseOrders.HandleFunc("/status/{status}", purchaseOrderHandler.ListByStatus).Methods("GET")
	purchaseOrders.HandleFunc("/{id}/items", purchaseOrderHandler.GetItems).Methods("GET")
	purchaseOrders.HandleFunc("/{id}/items", purchaseOrderHandler.CreateItem).Methods("POST")
	purchaseOrders.HandleFunc("/items/{itemId}/receive", purchaseOrderHandler.ReceiveItem).Methods("POST")

	// Stock Adjustments
	adjustments := api.PathPrefix("/stock-adjustments").Subrouter()
	adjustments.HandleFunc("", stockAdjustmentHandler.Create).Methods("POST")
	adjustments.HandleFunc("/{id}/approve", stockAdjustmentHandler.Approve).Methods("POST")
	adjustments.HandleFunc("/{id}/items", stockAdjustmentHandler.CreateItem).Methods("POST")

	// Stock Transfers
	transfers := api.PathPrefix("/stock-transfers").Subrouter()
	transfers.HandleFunc("", transferHandler.Create).Methods("POST")
	transfers.HandleFunc("/{id}/status", transferHandler.UpdateStatus).Methods("PUT")
	transfers.HandleFunc("/{id}/items", transferHandler.CreateItem).Methods("POST")
	transfers.HandleFunc("/items/{itemId}/quantities", transferHandler.UpdateItemQuantities).Methods("PUT")

	// Stocktakes
	stocktakes := api.PathPrefix("/stocktakes").Subrouter()
	stocktakes.HandleFunc("", stocktakeHandler.List).Methods("GET")
	stocktakes.HandleFunc("", stocktakeHandler.Create).Methods("POST")
	stocktakes.HandleFunc("/{id}", stocktakeHandler.Get).Methods("GET")
	stocktakes.HandleFunc("/{id}/status", stocktakeHandler.UpdateStatus).Methods("PUT")
	stocktakes.HandleFunc("/{id}/items", stocktakeHandler.GetItems).Methods("GET")
	stocktakes.HandleFunc("/{id}/items", stocktakeHandler.CreateItem).Methods("POST")
	stocktakes.HandleFunc("/items/{itemId}/count", stocktakeHandler.UpdateItemCount).Methods("PUT")
	stocktakes.HandleFunc("/{id}/variances", stocktakeHandler.GetVariances).Methods("GET")
	stocktakes.HandleFunc("/active", stocktakeHandler.GetActive).Methods("GET")
	stocktakes.HandleFunc("/warehouse/{warehouseId}", stocktakeHandler.ListByWarehouse).Methods("GET")

	// Warehouses
	warehouses := api.PathPrefix("/warehouses").Subrouter()
	warehouses.HandleFunc("", warehouseHandler.List).Methods("GET")
	warehouses.HandleFunc("", warehouseHandler.Create).Methods("POST")
	warehouses.HandleFunc("/{id}", warehouseHandler.Get).Methods("GET")
	warehouses.HandleFunc("/{id}", warehouseHandler.Update).Methods("PUT")
	warehouses.HandleFunc("/{id}", warehouseHandler.Deactivate).Methods("DELETE")
	warehouses.HandleFunc("/code/{code}", warehouseHandler.GetByCode).Methods("GET")
	warehouses.HandleFunc("/{id}/locations", warehouseHandler.ListLocations).Methods("GET")
	warehouses.HandleFunc("/{id}/locations", warehouseHandler.CreateLocation).Methods("POST")
	warehouses.HandleFunc("/summary", warehouseHandler.GetInventorySummary).Methods("GET")

	// Locations
	locations := api.PathPrefix("/locations").Subrouter()
	locations.HandleFunc("/{id}", warehouseHandler.GetLocation).Methods("GET")
	locations.HandleFunc("/{id}", warehouseHandler.UpdateLocation).Methods("PUT")
	locations.HandleFunc("/{id}", warehouseHandler.DeactivateLocation).Methods("DELETE")

	suppliers := api.PathPrefix("/suppliers").Subrouter()
	suppliers.HandleFunc("", supplierHandler.List).Methods("GET")
	suppliers.HandleFunc("", supplierHandler.Create).Methods("POST")
	suppliers.HandleFunc("/active", supplierHandler.ListActive).Methods("GET") // MUST COME BEFORE /{id}
	suppliers.HandleFunc("/code/{code}", supplierHandler.GetByCode).Methods("GET")
	suppliers.HandleFunc("/search", supplierHandler.Search).Methods("GET")
	suppliers.HandleFunc("/{id}", supplierHandler.Get).Methods("GET")
	suppliers.HandleFunc("/{id}", supplierHandler.Update).Methods("PUT")
	suppliers.HandleFunc("/{id}/deactivate", supplierHandler.Deactivate).Methods("POST")
	suppliers.HandleFunc("/{id}/activate", supplierHandler.Activate).Methods("POST")
	suppliers.HandleFunc("/{id}/products", supplierHandler.GetProducts).Methods("GET")
	suppliers.HandleFunc("/{id}/performance", supplierHandler.GetPerformance).Methods("GET")

	return r
}
