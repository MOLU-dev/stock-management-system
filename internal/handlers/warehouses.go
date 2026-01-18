package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	db "github.com/molu/stock-management-system/internal/db/sqlc"
)

type WarehouseHandler struct {
	queries *db.Queries
}

func NewWarehouseHandler(queries *db.Queries) *WarehouseHandler {
	return &WarehouseHandler{queries: queries}
}

type CreateWarehouseRequest struct {
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	Address       *string `json:"address"`
	ContactPerson *string `json:"contact_person"`
	ContactPhone  *string `json:"contact_phone"`
	ContactEmail  *string `json:"contact_email"`
}

func (h *WarehouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateWarehouseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	warehouse, err := h.queries.CreateWarehouse(ctx, db.CreateWarehouseParams{
		Code:          req.Code,
		Name:          req.Name,
		Address:       req.Address,
		ContactPerson: req.ContactPerson,
		ContactPhone:  req.ContactPhone,
		ContactEmail:  req.ContactEmail,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create warehouse")
		return
	}

	respondJSON(w, http.StatusCreated, warehouse)
}

func (h *WarehouseHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	warehouse, err := h.queries.GetWarehouse(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Warehouse not found")
		return
	}

	respondJSON(w, http.StatusOK, warehouse)
}

func (h *WarehouseHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	warehouse, err := h.queries.GetWarehouseByCode(ctx, vars["code"])
	if err != nil {
		respondError(w, http.StatusNotFound, "Warehouse not found")
		return
	}

	respondJSON(w, http.StatusOK, warehouse)
}

func (h *WarehouseHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	warehouses, err := h.queries.ListWarehouses(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch warehouses")
		return
	}

	respondJSON(w, http.StatusOK, warehouses)
}

func (h *WarehouseHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	var req CreateWarehouseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	warehouse, err := h.queries.UpdateWarehouse(ctx, db.UpdateWarehouseParams{
		WarehouseID:   id,
		Name:          req.Name,
		Address:       req.Address,
		ContactPerson: req.ContactPerson,
		ContactPhone:  req.ContactPhone,
		ContactEmail:  req.ContactEmail,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update warehouse")
		return
	}

	respondJSON(w, http.StatusOK, warehouse)
}

func (h *WarehouseHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	if err := h.queries.DeactivateWarehouse(ctx, id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to deactivate warehouse")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type CreateLocationRequest struct {
	LocationCode string  `json:"location_code"`
	Aisle        *string `json:"aisle"`
	Shelf        *string `json:"shelf"`
	Bin          *string `json:"bin"`
	MaxCapacity  *int32  `json:"max_capacity"`
}

func (h *WarehouseHandler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	warehouseID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	var req CreateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	location, err := h.queries.CreateLocation(ctx, db.CreateLocationParams{
		WarehouseID:  warehouseID,
		LocationCode: req.LocationCode,
		Aisle:        req.Aisle,
		Shelf:        req.Shelf,
		Bin:          req.Bin,
		MaxCapacity:  req.MaxCapacity,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create location")
		return
	}

	respondJSON(w, http.StatusCreated, location)
}

func (h *WarehouseHandler) GetLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid location ID")
		return
	}

	location, err := h.queries.GetLocation(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Location not found")
		return
	}

	respondJSON(w, http.StatusOK, location)
}

func (h *WarehouseHandler) ListLocations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	warehouseID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	locations, err := h.queries.ListLocationsByWarehouse(ctx, warehouseID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch locations")
		return
	}

	respondJSON(w, http.StatusOK, locations)
}

func (h *WarehouseHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid location ID")
		return
	}

	var req CreateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	location, err := h.queries.UpdateLocation(ctx, db.UpdateLocationParams{
		LocationID:  id,
		Aisle:       req.Aisle,
		Shelf:       req.Shelf,
		Bin:         req.Bin,
		MaxCapacity: req.MaxCapacity,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update location")
		return
	}

	respondJSON(w, http.StatusOK, location)
}

func (h *WarehouseHandler) DeactivateLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid location ID")
		return
	}

	if err := h.queries.DeactivateLocation(ctx, id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to deactivate location")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WarehouseHandler) GetInventorySummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	summary, err := h.queries.GetWarehouseInventorySummary(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch summary")
		return
	}

	respondJSON(w, http.StatusOK, summary)
}
