package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	db "github.com/molu/stock-management-system/internal/db/sqlc"
)

type StocktakeHandler struct {
	queries *db.Queries
}

func NewStocktakeHandler(queries *db.Queries) *StocktakeHandler {
	return &StocktakeHandler{queries: queries}
}

type CreateStocktakeRequest struct {
	StocktakeNumber string     `json:"stocktake_number"`
	WarehouseID     int64      `json:"warehouse_id"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	Status          string     `json:"status"`
	Notes           *string    `json:"notes"`
	CreatedBy       int64      `json:"created_by"`
}

func (h *StocktakeHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateStocktakeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	stocktake, err := h.queries.CreateStocktake(ctx, db.CreateStocktakeParams{
		StocktakeNumber: req.StocktakeNumber,
		WarehouseID:     req.WarehouseID,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		Status:          db.StocktakeStatus(req.Status),
		Notes:           req.Notes,
		CreatedBy:       req.CreatedBy,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create stocktake")
		return
	}

	respondJSON(w, http.StatusCreated, stocktake)
}

func (h *StocktakeHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid stocktake ID")
		return
	}

	stocktake, err := h.queries.GetStocktake(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Stocktake not found")
		return
	}

	respondJSON(w, http.StatusOK, stocktake)
}

func (h *StocktakeHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := int32(50)
	offset := int32(0)

	stocktakes, err := h.queries.ListStocktakes(ctx, db.ListStocktakesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch stocktakes")
		return
	}

	respondJSON(w, http.StatusOK, stocktakes)
}

func (h *StocktakeHandler) ListByWarehouse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	warehouseID, err := strconv.ParseInt(vars["warehouseId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	stocktakes, err := h.queries.ListStocktakesByWarehouse(ctx, warehouseID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch stocktakes")
		return
	}

	respondJSON(w, http.StatusOK, stocktakes)
}

type UpdateStocktakeStatusRequest struct {
	Status string `json:"status"`
}

func (h *StocktakeHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid stocktake ID")
		return
	}

	var req UpdateStocktakeStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	stocktake, err := h.queries.UpdateStocktakeStatus(ctx, db.UpdateStocktakeStatusParams{
		StocktakeID: id,
		Status:      db.StocktakeStatus(req.Status),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update stocktake")
		return
	}

	respondJSON(w, http.StatusOK, stocktake)
}

type CreateStocktakeItemRequest struct {
	ProductID       int64      `json:"product_id"`
	LocationID      *int64     `json:"location_id"`
	SystemQuantity  int32      `json:"system_quantity"`
	CountedQuantity *int32     `json:"counted_quantity"`
	Variance        *int32     `json:"variance"`
	CountedBy       *int64     `json:"counted_by"`
	CountedAt       *time.Time `json:"counted_at"`
	Notes           *string    `json:"notes"`
}

func (h *StocktakeHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	stocktakeID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid stocktake ID")
		return
	}

	var req CreateStocktakeItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.queries.CreateStocktakeItem(ctx, db.CreateStocktakeItemParams{
		StocktakeID:     stocktakeID,
		ProductID:       req.ProductID,
		LocationID:      req.LocationID,
		SystemQuantity:  req.SystemQuantity,
		CountedQuantity: req.CountedQuantity,
		Variance:        req.Variance,
		CountedBy:       req.CountedBy,
		CountedAt:       req.CountedAt,
		Notes:           req.Notes,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	respondJSON(w, http.StatusCreated, item)
}

func (h *StocktakeHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid stocktake ID")
		return
	}

	items, err := h.queries.GetStocktakeItems(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch items")
		return
	}

	respondJSON(w, http.StatusOK, items)
}

type UpdateStocktakeItemCountRequest struct {
	CountedQuantity int32 `json:"counted_quantity"`
	CountedBy       int64 `json:"counted_by"`
}

func (h *StocktakeHandler) UpdateItemCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	itemID, err := strconv.ParseInt(vars["itemId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var req UpdateStocktakeItemCountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.queries.UpdateStocktakeItemCount(ctx, db.UpdateStocktakeItemCountParams{
		StocktakeItemID: itemID,
		CountedQuantity: req.CountedQuantity,
		CountedBy:       &req.CountedBy,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update count")
		return
	}

	respondJSON(w, http.StatusOK, item)
}

func (h *StocktakeHandler) GetVariances(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid stocktake ID")
		return
	}

	variances, err := h.queries.GetStocktakeVariances(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch variances")
		return
	}

	respondJSON(w, http.StatusOK, variances)
}

func (h *StocktakeHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stocktakes, err := h.queries.GetActiveStocktakes(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch active stocktakes")
		return
	}

	respondJSON(w, http.StatusOK, stocktakes)
}
