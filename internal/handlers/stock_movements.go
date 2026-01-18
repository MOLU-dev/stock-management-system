package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	db "github.com/molu/stock-management-system/internal/db/sqlc"
)

type StockMovementHandler struct {
	queries *db.Queries
}

func NewStockMovementHandler(queries *db.Queries) *StockMovementHandler {
	return &StockMovementHandler{queries: queries}
}

type CreateStockMovementRequest struct {
	ReferenceNumber string  `json:"reference_number"`
	ProductID       int64   `json:"product_id"`
	WarehouseID     int64   `json:"warehouse_id"`
	LocationID      *int64  `json:"location_id"`
	MovementType    string  `json:"movement_type"`
	QuantityBefore  int32   `json:"quantity_before"`
	QuantityChange  int32   `json:"quantity_change"`
	QuantityAfter   int32   `json:"quantity_after"`
	ReferenceID     *int64  `json:"reference_id"`
	ReferenceTable  *string `json:"reference_table"`
	Notes           *string `json:"notes"`
	CreatedBy       int64   `json:"created_by"`
}

func (h *StockMovementHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateStockMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	movement, err := h.queries.CreateStockMovement(ctx, db.CreateStockMovementParams{
		ReferenceNumber: req.ReferenceNumber,
		ProductID:       req.ProductID,
		WarehouseID:     req.WarehouseID,
		LocationID:      req.LocationID,
		MovementType:    db.MovementType(req.MovementType),
		QuantityBefore:  req.QuantityBefore,
		QuantityChange:  req.QuantityChange,
		QuantityAfter:   req.QuantityAfter,
		ReferenceID:     req.ReferenceID,
		ReferenceTable:  req.ReferenceTable,
		Notes:           req.Notes,
		CreatedBy:       req.CreatedBy,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create stock movement")
		return
	}

	respondJSON(w, http.StatusCreated, movement)
}

func (h *StockMovementHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid movement ID")
		return
	}

	movement, err := h.queries.GetStockMovement(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Stock movement not found")
		return
	}

	respondJSON(w, http.StatusOK, movement)
}

func (h *StockMovementHandler) ListByProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	productID, err := strconv.ParseInt(vars["productId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	limit := int32(50)
	offset := int32(0)

	movements, err := h.queries.ListStockMovementsByProduct(ctx, db.ListStockMovementsByProductParams{
		ProductID: productID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch stock movements")
		return
	}

	respondJSON(w, http.StatusOK, movements)
}

func (h *StockMovementHandler) ListByWarehouse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	warehouseID, err := strconv.ParseInt(vars["warehouseId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	limit := int32(50)
	offset := int32(0)

	movements, err := h.queries.ListStockMovementsByWarehouse(ctx, db.ListStockMovementsByWarehouseParams{
		WarehouseID: warehouseID,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch stock movements")
		return
	}

	respondJSON(w, http.StatusOK, movements)
}

func (h *StockMovementHandler) ListByType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	limit := int32(50)
	offset := int32(0)

	movements, err := h.queries.ListStockMovementsByType(ctx, db.ListStockMovementsByTypeParams{
		MovementType: db.MovementType(vars["type"]),
		Limit:        limit,
		Offset:       offset,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch stock movements")
		return
	}

	respondJSON(w, http.StatusOK, movements)
}

func (h *StockMovementHandler) GetProductHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	productID, err := strconv.ParseInt(vars["productId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	warehouseID, err := strconv.ParseInt(vars["warehouseId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	limit := int32(50)
	offset := int32(0)

	movements, err := h.queries.GetProductMovementHistory(ctx, db.GetProductMovementHistoryParams{
		ProductID:   productID,
		WarehouseID: warehouseID,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch movement history")
		return
	}

	respondJSON(w, http.StatusOK, movements)
}
