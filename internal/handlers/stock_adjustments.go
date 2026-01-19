package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	db "github.com/molu/stock-management-system/internal/db/sqlc"
	"github.com/shopspring/decimal"
)

func NullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}


type StockAdjustmentHandler struct {
	queries db.SingleDb
}

func NewStockAdjustmentHandler(queries db.SingleDb) *StockAdjustmentHandler {
	return &StockAdjustmentHandler{queries: queries}
}

type CreateStockAdjustmentRequest struct {
	AdjustmentNumber string          `json:"adjustment_number"`
	WarehouseID      int64           `json:"warehouse_id"`
	AdjustmentDate   time.Time       `json:"adjustment_date"`
	Reason           string          `json:"reason"`
	Status           string          `json:"status"`
	TotalValue       decimal.Decimal `json:"total_value"`
	Notes            *string         `json:"notes"`
	CreatedBy        int64           `json:"created_by"`
}

func (h *StockAdjustmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateStockAdjustmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	adjustment, err := h.queries.CreateStockAdjustment(ctx, db.CreateStockAdjustmentParams{
		AdjustmentNumber: req.AdjustmentNumber,
		WarehouseID:      int32(req.WarehouseID),
		AdjustmentDate:   req.AdjustmentDate,
		Reason:           db.AdjustmentReason(req.Reason),
		Status:           db.AdjustmentStatus(req.Status),
		TotalValue:       req.TotalValue,
		Notes:            NullString(req.Notes),
		CreatedBy:        NullInt32(req.CreatedBy),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create adjustment")
		return
	}

	respondJSON(w, http.StatusCreated, adjustment)
}

type ApproveAdjustmentRequest struct {
	ApprovedBy int64 `json:"approved_by"`
}

func (h *StockAdjustmentHandler) Approve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid adjustment ID")
		return
	}

	var req ApproveAdjustmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	adjustment, err := h.queries.ApproveStockAdjustment(ctx, db.ApproveStockAdjustmentParams{
		AdjustmentID: int32(id),
		ApprovedBy:   NullInt32(req.ApprovedBy),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to approve adjustment")
		return
	}

	respondJSON(w, http.StatusOK, adjustment)
}

type CreateAdjustmentItemRequest struct {
	ProductID        int64           `json:"product_id"`
	QuantityBefore   int32           `json:"quantity_before"`
	QuantityAdjusted int32           `json:"quantity_adjusted"`
	CostPrice        decimal.Decimal `json:"cost_price"`
	Reason           string          `json:"reason"`
}

func (h *StockAdjustmentHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	adjustmentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid adjustment ID")
		return
	}

	var req CreateAdjustmentItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.queries.CreateStockAdjustmentItem(ctx, db.CreateStockAdjustmentItemParams{
		AdjustmentID:     int32(adjustmentID),
		ProductID:        int32(req.ProductID),
		QuantityBefore:   req.QuantityBefore,
		QuantityAdjusted: req.QuantityAdjusted,
		CostPrice:        req.CostPrice,
		Reason:           sql.NullString{String: req.Reason, Valid: true},
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	respondJSON(w, http.StatusCreated, item)
}
