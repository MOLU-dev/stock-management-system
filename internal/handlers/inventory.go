package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	db "github.com/molu/stock-management-system/internal/db/sqlc"
)

type InventoryHandler struct {
	queries *db.Queries
}

func NewInventoryHandler(queries *db.Queries) *InventoryHandler {
	return &InventoryHandler{queries: queries}
}

func (h *InventoryHandler) List(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"message": "list all inventory"})
}

func (h *InventoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	inventory, err := h.queries.GetInventory(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Inventory not found")
		return
	}

	respondJSON(w, http.StatusOK, inventory)
}

func (h *InventoryHandler) GetByProductWarehouse(w http.ResponseWriter, r *http.Request) {
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

	inventory, err := h.queries.GetInventoryByProductWarehouse(ctx, db.GetInventoryByProductWarehouseParams{
		ProductID:   productID,
		WarehouseID: warehouseID,
	})
	if err != nil {
		respondError(w, http.StatusNotFound, "Inventory not found")
		return
	}

	respondJSON(w, http.StatusOK, inventory)
}

func (h *InventoryHandler) ListByWarehouse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	warehouseID, err := strconv.ParseInt(vars["warehouseId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	limit := int32(50)
	offset := int32(0)

	inventory, err := h.queries.ListInventoryByWarehouse(ctx, db.ListInventoryByWarehouseParams{
		WarehouseID: warehouseID,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch inventory")
		return
	}

	respondJSON(w, http.StatusOK, inventory)
}

func (h *InventoryHandler) ListByProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	productID, err := strconv.ParseInt(vars["productId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	inventory, err := h.queries.ListInventoryByProduct(ctx, productID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch inventory")
		return
	}

	respondJSON(w, http.StatusOK, inventory)
}

func (h *InventoryHandler) ListExpiring(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	inventory, err := h.queries.ListExpiringInventory(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch expiring inventory")
		return
	}

	respondJSON(w, http.StatusOK, inventory)
}

type UpdateQuantityRequest struct {
	Quantity         int32 `json:"quantity"`
	ReservedQuantity int32 `json:"reserved_quantity"`
}

func (h *InventoryHandler) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	var req UpdateQuantityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	inventory, err := h.queries.UpdateInventoryQuantity(ctx, db.UpdateInventoryQuantityParams{
		InventoryID:      id,
		Quantity:         req.Quantity,
		ReservedQuantity: req.ReservedQuantity,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update inventory")
		return
	}

	respondJSON(w, http.StatusOK, inventory)
}

type ReserveRequest struct {
	Quantity int32 `json:"quantity"`
}

func (h *InventoryHandler) Reserve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	var req ReserveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	inventory, err := h.queries.ReserveInventory(ctx, db.ReserveInventoryParams{
		InventoryID: id,
		Column2:     req.Quantity,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to reserve inventory")
		return
	}

	respondJSON(w, http.StatusOK, inventory)
}

func (h *InventoryHandler) Release(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	var req ReserveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	inventory, err := h.queries.ReleaseInventoryReservation(ctx, db.ReleaseInventoryReservationParams{
		InventoryID: id,
		Column2:     req.Quantity,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to release inventory")
		return
	}

	respondJSON(w, http.StatusOK, inventory)
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

func (h *InventoryHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	inventory, err := h.queries.UpdateInventoryStatus(ctx, db.UpdateInventoryStatusParams{
		InventoryID: id,
		Status:      db.InventoryStatus(req.Status),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update inventory status")
		return
	}

	respondJSON(w, http.StatusOK, inventory)
}
