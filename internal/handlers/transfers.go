package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	db "github.com/molu/stock-management-system/internal/db/sqlc"
)

type TransferHandler struct {
	queries *db.Queries
}

func NewTransferHandler(queries *db.Queries) *TransferHandler {
	return &TransferHandler{queries: queries}
}

type CreateStockTransferRequest struct {
	TransferNumber         string     `json:"transfer_number"`
	FromWarehouseID        int64      `json:"from_warehouse_id"`
	ToWarehouseID          int64      `json:"to_warehouse_id"`
	Status                 string     `json:"status"`
	TransferDate           *time.Time `json:"transfer_date"`
	ExpectedCompletionDate time.Time  `json:"expected_completion_date"`
	Notes                  *string    `json:"notes"`
	CreatedBy              int64      `json:"created_by"`
}

func (h *TransferHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateStockTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	transfer, err := h.queries.CreateStockTransfer(ctx, db.CreateStockTransferParams{
		TransferNumber:         req.TransferNumber,
		FromWarehouseID:        req.FromWarehouseID,
		ToWarehouseID:          req.ToWarehouseID,
		Status:                 db.TransferStatus(req.Status),
		TransferDate:           req.TransferDate,
		ExpectedCompletionDate: req.ExpectedCompletionDate,
		Notes:                  req.Notes,
		CreatedBy:              req.CreatedBy,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create transfer")
		return
	}

	respondJSON(w, http.StatusCreated, transfer)
}

type UpdateTransferStatusRequest struct {
	Status string `json:"status"`
}

func (h *TransferHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid transfer ID")
		return
	}

	var req UpdateTransferStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	transfer, err := h.queries.UpdateStockTransferStatus(ctx, db.UpdateStockTransferStatusParams{
		TransferID: id,
		Status:     db.TransferStatus(req.Status),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update transfer")
		return
	}

	respondJSON(w, http.StatusOK, transfer)
}

type CreateTransferItemRequest struct {
	ProductID      int64  `json:"product_id"`
	Quantity       int32  `json:"quantity"`
	FromLocationID *int64 `json:"from_location_id"`
	ToLocationID   *int64 `json:"to_location_id"`
}

func (h *TransferHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	transferID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid transfer ID")
		return
	}

	var req CreateTransferItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.queries.CreateStockTransferItem(ctx, db.CreateStockTransferItemParams{
		TransferID:     transferID,
		ProductID:      req.ProductID,
		Quantity:       req.Quantity,
		FromLocationID: req.FromLocationID,
		ToLocationID:   req.ToLocationID,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	respondJSON(w, http.StatusCreated, item)
}

type UpdateTransferQuantitiesRequest struct {
	QuantitySent     int32 `json:"quantity_sent"`
	QuantityReceived int32 `json:"quantity_received"`
}

func (h *TransferHandler) UpdateItemQuantities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	itemID, err := strconv.ParseInt(vars["itemId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var req UpdateTransferQuantitiesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.queries.UpdateStockTransferItemQuantities(ctx, db.UpdateStockTransferItemQuantitiesParams{
		TransferItemID:   itemID,
		QuantitySent:     req.QuantitySent,
		QuantityReceived: req.QuantityReceived,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update quantities")
		return
	}

	respondJSON(w, http.StatusOK, item)
}
