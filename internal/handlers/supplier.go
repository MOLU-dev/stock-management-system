package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	db "github.com/molu/stock-management-system/internal/db/sqlc"
	"github.com/shopspring/decimal"
)

type SupplierHandler struct {
	queries db.SingleDb
}

func NewSupplierHandler(queries db.SingleDb) *SupplierHandler {
	return &SupplierHandler{queries: queries}
}

// CreateSupplierRequest defines the request structure for creating a supplier
type CreateSupplierRequest struct {
	Code          string           `json:"code"`
	Name          string           `json:"name"`
	ContactPerson *string          `json:"contact_person"`
	Email         *string          `json:"email"`
	Phone         *string          `json:"phone"`
	Address       *string          `json:"address"`
	TaxID         *string          `json:"tax_id"`
	PaymentTerms  *string          `json:"payment_terms"`
	LeadTimeDays  *int32           `json:"lead_time_days"`
	Rating        *decimal.Decimal `json:"rating"`
	IsActive      bool             `json:"is_active"`
}

func (h *SupplierHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var rating decimal.Decimal
	if req.Rating != nil {
		rating = *req.Rating
	} else {
		rating = decimal.Zero
	}

	supplier, err := h.queries.CreateSupplier(ctx, db.CreateSupplierParams{
		Code:          req.Code,
		Name:          req.Name,
		ContactPerson: toNullString(req.ContactPerson),
		Email:         toNullString(req.Email),
		Phone:         toNullString(req.Phone),
		Address:       toNullString(req.Address),
		TaxID:         toNullString(req.TaxID),
		PaymentTerms:  toNullString(req.PaymentTerms),
		LeadTimeDays:  toNullInt32FromInt32(req.LeadTimeDays),
		Rating: rating,

	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create supplier")
		return
	}

	respondJSON(w, http.StatusCreated, supplier)
}

func (h *SupplierHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id64, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid supplier ID")
		return
	}
	id := int32(id64)

	supplier, err := h.queries.GetSupplier(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "Supplier not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch supplier")
		return
	}

	respondJSON(w, http.StatusOK, supplier)
}

func (h *SupplierHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	supplier, err := h.queries.GetSupplierByCode(ctx, vars["code"])
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "Supplier not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch supplier")
		return
	}

	respondJSON(w, http.StatusOK, supplier)
}

func (h *SupplierHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := int32(50)
	offset := int32(0)

	// Parse query parameters
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.ParseInt(l, 10, 32); err == nil {
			limit = int32(val)
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.ParseInt(o, 10, 32); err == nil {
			offset = int32(val)
		}
	}

	// Check if we want all suppliers (including inactive)
	showAll := r.URL.Query().Get("showAll") == "true"

	var suppliers []db.Supplier
	var err error

	if showAll {
		suppliers, err = h.queries.ListAllSuppliers(ctx, db.ListAllSuppliersParams{
			Limit:  limit,
			Offset: offset,
		})
	} else {
		suppliers, err = h.queries.ListSuppliers(ctx, db.ListSuppliersParams{
			Limit:  limit,
			Offset: offset,
		})
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch suppliers")
		return
	}

	respondJSON(w, http.StatusOK, suppliers)
}

func (h *SupplierHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id64, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid supplier ID")
		return
	}
	id := int32(id64)

	var req CreateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var rating decimal.Decimal
	if req.Rating != nil {
		rating = *req.Rating
	} else {
		rating = decimal.Zero
	}
	

	
	searchTerm := req.Name
	supplier, err := h.queries.UpdateSupplier(ctx, db.UpdateSupplierParams{
		SupplierID:    id,
		Name:          searchTerm,
		ContactPerson: toNullString(req.ContactPerson),
		Email:         toNullString(req.Email),
		Phone:         toNullString(req.Phone),
		Address:       toNullString(req.Address),
		TaxID:         toNullString(req.TaxID),
		PaymentTerms:  toNullString(req.PaymentTerms),
		LeadTimeDays:  toNullInt32FromInt32(req.LeadTimeDays),
		Rating:rating,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update supplier")
		return
	}

	respondJSON(w, http.StatusOK, supplier)
}

func (h *SupplierHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id64, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid supplier ID")
		return
	}
	id := int32(id64)

	if err := h.queries.DeactivateSupplier(ctx, id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to deactivate supplier")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SupplierHandler) Activate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id64, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid supplier ID")
		return
	}
	id := int32(id64)

	if err := h.queries.ActivateSupplier(ctx, id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to activate supplier")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Supplier activated successfully"})
}

func (h *SupplierHandler) Search(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    searchTerm := r.URL.Query().Get("q")
    if searchTerm == "" {
        respondError(w, http.StatusBadRequest, "Search term is required")
        return
    }

    limit := int32(50)
    offset := int32(0)

    if l := r.URL.Query().Get("limit"); l != "" {
        if val, err := strconv.ParseInt(l, 10, 32); err == nil {
            limit = int32(val)
        }
    }

    if o := r.URL.Query().Get("offset"); o != "" {
        if val, err := strconv.ParseInt(o, 10, 32); err == nil {
            offset = int32(val)
        }
    }

    suppliers, err := h.queries.SearchSuppliers(ctx, db.SearchSuppliersParams{
        Column1: sql.NullString{
            String: searchTerm,
            Valid:  true,
        },
        Limit:  limit,
        Offset: offset,
    })
    if err != nil {
        respondError(w, http.StatusInternalServerError, "Failed to search suppliers")
        return
    }

    respondJSON(w, http.StatusOK, suppliers)
}



func (h *SupplierHandler) ListActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	suppliers, err := h.queries.ListActiveSuppliers(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch active suppliers")
		return
	}

	respondJSON(w, http.StatusOK, suppliers)
}

func (h *SupplierHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id64, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid supplier ID")
		return
	}
	id := int32(id64)

	limit := int32(50)
	offset := int32(0)

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.ParseInt(l, 10, 32); err == nil {
			limit = int32(val)
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.ParseInt(o, 10, 32); err == nil {
			offset = int32(val)
		}
	}

	products, err := h.queries.GetSupplierProducts(ctx, db.GetSupplierProductsParams{
		SupplierID: id,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch supplier products")
		return
	}

	respondJSON(w, http.StatusOK, products)
}

func (h *SupplierHandler) GetPerformance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id64, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid supplier ID")
		return
	}
	id := int32(id64)

	performance, err := h.queries.GetSupplierPerformance(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch supplier performance")
		return
	}

	respondJSON(w, http.StatusOK, performance)
}
