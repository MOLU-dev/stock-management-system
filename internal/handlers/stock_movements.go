package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	db "github.com/molu/stock-management-system/internal/db/sqlc"
)

func toNullStringFromValue(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func toNullInt32FromInt32(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{
		Int32: *i,
		Valid: true,
	}
}

func toNullInt32FromValue(i int64) sql.NullInt32 {
	if i == 0 {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{
		Int32: int32(i),
		Valid: true,
	}
}

type StockMovementHandler struct {
	queries db.SingleDb
}

func NewStockMovementHandler(queries db.SingleDb) *StockMovementHandler {
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
		ReferenceNumber: toNullStringFromValue(req.ReferenceNumber), // string → sql.NullString
		ProductID:       int32(req.ProductID),
		WarehouseID:     int32(req.WarehouseID),
		LocationID:      toNullInt32FromInt64(req.LocationID), // *int64 → sql.NullInt32
		MovementType:    db.MovementType(req.MovementType),

		// these are int32 in your request, so convert properly
		QuantityBefore: func() sql.NullInt32 {
			q := req.QuantityBefore
			return toNullInt32FromInt32(&q)
		}(),

		QuantityChange: req.QuantityChange,

		QuantityAfter: func() sql.NullInt32 {
			q := req.QuantityAfter
			return toNullInt32FromInt32(&q)
		}(),

		ReferenceID:    toNullInt32FromInt64(req.ReferenceID), // *int64 → sql.NullInt32
		ReferenceTable: toNullString(req.ReferenceTable),      // *string → sql.NullString
		Notes:          toNullString(req.Notes),               // *string → sql.NullString
		CreatedBy:      toNullInt32FromValue(req.CreatedBy),   // int64 → sql.NullInt32
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

	movement, err := h.queries.GetStockMovement(ctx, int32(id))
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
		ProductID: int32(productID),
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
		WarehouseID: int32(warehouseID),
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
		ProductID:   int32(productID),
		WarehouseID: int32(warehouseID),
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch movement history")
		return
	}

	respondJSON(w, http.StatusOK, movements)
}
