package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	db "github.com/molu/stock-management-system/internal/db/sqlc"
)

func NullInt32(i int64) sql.NullInt32 {
	if i == 0 {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{
		Int32: int32(i),
		Valid: true,
	}
}



type WarehouseHandler struct {
	queries db.SingleDb
}

func NewWarehouseHandler(queries db.SingleDb) *WarehouseHandler {
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
		Address:       toNullString(req.Address),
		ContactPerson: toNullString(req.ContactPerson),
		ContactPhone:  toNullString(req.ContactPhone),
		ContactEmail:  toNullString(req.ContactEmail),
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

	warehouse, err := h.queries.GetWarehouse(ctx, int32(id))
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
		WarehouseID:   int32(id),
		Name:          req.Name,
		Address:       toNullString(req.Address),
		ContactPerson: toNullString(req.ContactPerson),
		ContactPhone:  toNullString(req.ContactPhone),
		ContactEmail:  toNullString(req.ContactEmail),
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

	if err := h.queries.DeactivateWarehouse(ctx, int32(id)); err != nil {
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
		WarehouseID:  int32(warehouseID),
		LocationCode: req.LocationCode,
		Aisle:        toNullString(req.Aisle),
		Shelf:        toNullString(req.Shelf),
		Bin:          toNullString(req.Bin),
		MaxCapacity:  toNullInt32FromInt32(req.MaxCapacity),
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

	location, err := h.queries.GetLocation(ctx, int32(id))
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

	locations, err := h.queries.ListLocationsByWarehouse(ctx, int32(warehouseID))
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
		LocationID:  int32(id),
		Aisle:       toNullString(req.Aisle),
		Shelf:       toNullString(req.Shelf),
		Bin:         toNullString(req.Bin),
		MaxCapacity: toNullInt32FromInt32(req.MaxCapacity),
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

	if err := h.queries.DeactivateLocation(ctx, int32(id)); err != nil {
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
