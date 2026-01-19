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

type ProductHandler struct {
	queries db.SingleDb
}

func NewProductHandler(queries db.SingleDb) *ProductHandler {
	return &ProductHandler{queries: queries}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	products, err := h.queries.ListProducts(ctx, db.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	respondJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id64, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	id := int32(id64)

	product, err := h.queries.GetProduct(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Product not found")
		return
	}

	respondJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) GetBySKU(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	product, err := h.queries.GetProductBySKU(ctx, vars["sku"])
	if err != nil {
		respondError(w, http.StatusNotFound, "Product not found")
		return
	}

	respondJSON(w, http.StatusOK, product)
}

type CreateProductRequest struct {
	SKU           string           `json:"sku"`
	Name          string           `json:"name"`
	Description   *string          `json:"description"`
	CategoryID    int64            `json:"category_id"`
	UnitPrice     decimal.Decimal  `json:"unit_price"`
	CostPrice     decimal.Decimal  `json:"cost_price"`
	Barcode       *string          `json:"barcode"`
	Weight        *decimal.Decimal `json:"weight"`
	Dimensions    *string          `json:"dimensions"`
	SupplierID    *int64           `json:"supplier_id"`
	MinStockLevel int32            `json:"min_stock_level"`
	MaxStockLevel int32            `json:"max_stock_level"`
	ReorderPoint  int32            `json:"reorder_point"`
	SafetyStock   int32            `json:"safety_stock"`
	LeadTimeDays  *int32           `json:"lead_time_days"`
	AutoReorder   bool             `json:"auto_reorder"`
	IsActive      bool             `json:"is_active"`
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Convert nullable string fields
	var description sql.NullString
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}

	var barcode sql.NullString
	if req.Barcode != nil {
		barcode = sql.NullString{String: *req.Barcode, Valid: true}
	}

	var dimensions sql.NullString
	if req.Dimensions != nil {
		dimensions = sql.NullString{String: *req.Dimensions, Valid: true}
	}

	// Convert weight - handle nil pointer
	var weight decimal.Decimal
	if req.Weight != nil {
		weight = *req.Weight
	} else {
		weight = decimal.Zero
	}

	// Convert category ID
	var categoryID sql.NullInt32
	categoryID = sql.NullInt32{Int32: int32(req.CategoryID), Valid: true}

	// Convert supplier ID
	var supplierID sql.NullInt32
	if req.SupplierID != nil {
		supplierID = sql.NullInt32{Int32: int32(*req.SupplierID), Valid: true}
	}

	// Convert lead time days
	var leadTimeDays sql.NullInt32
	if req.LeadTimeDays != nil {
		leadTimeDays = sql.NullInt32{Int32: *req.LeadTimeDays, Valid: true}
	}

	product, err := h.queries.CreateProduct(ctx, db.CreateProductParams{
		Sku:           req.SKU,
		Name:          req.Name,
		Description:   description,
		CategoryID:    categoryID,
		UnitPrice:     req.UnitPrice,
		CostPrice:     req.CostPrice,
		Barcode:       barcode,
		Weight:        weight,
		Dimensions:    dimensions,
		SupplierID:    supplierID,
		MinStockLevel: req.MinStockLevel,
		MaxStockLevel: sql.NullInt32{Int32: req.MaxStockLevel, Valid: true},
		ReorderPoint:  sql.NullInt32{Int32: req.ReorderPoint, Valid: true},
		SafetyStock:   sql.NullInt32{Int32: req.SafetyStock, Valid: true},
		LeadTimeDays:  leadTimeDays,
		AutoReorder:   req.AutoReorder,
		IsActive:      req.IsActive,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	respondJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id64, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	id := int32(id64)

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Convert nullable string fields
	var description sql.NullString
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}

	var barcode sql.NullString
	if req.Barcode != nil {
		barcode = sql.NullString{String: *req.Barcode, Valid: true}
	}

	var dimensions sql.NullString
	if req.Dimensions != nil {
		dimensions = sql.NullString{String: *req.Dimensions, Valid: true}
	}

	// Convert weight - handle nil pointer
	var weight decimal.Decimal
	if req.Weight != nil {
		weight = *req.Weight
	} else {
		weight = decimal.Zero
	}

	// Convert category ID
	var categoryID sql.NullInt32
	categoryID = sql.NullInt32{Int32: int32(req.CategoryID), Valid: true}

	// Convert supplier ID
	var supplierID sql.NullInt32
	if req.SupplierID != nil {
		supplierID = sql.NullInt32{Int32: int32(*req.SupplierID), Valid: true}
	}

	// Convert lead time days
	var leadTimeDays sql.NullInt32
	if req.LeadTimeDays != nil {
		leadTimeDays = sql.NullInt32{Int32: *req.LeadTimeDays, Valid: true}
	}

	product, err := h.queries.UpdateProduct(ctx, db.UpdateProductParams{
		ProductID:     id,
		Name:          req.Name,
		Description:   description,
		CategoryID:    categoryID,
		UnitPrice:     req.UnitPrice,
		CostPrice:     req.CostPrice,
		Barcode:       barcode,
		Weight:        weight,
		Dimensions:    dimensions,
		SupplierID:    supplierID,
		MinStockLevel: req.MinStockLevel,
		MaxStockLevel: sql.NullInt32{Int32: req.MaxStockLevel, Valid: true},
		ReorderPoint:  sql.NullInt32{Int32: req.ReorderPoint, Valid: true},
		SafetyStock:   sql.NullInt32{Int32: req.SafetyStock, Valid: true},
		LeadTimeDays:  leadTimeDays,
		AutoReorder:   req.AutoReorder,
		IsActive:      req.IsActive,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	respondJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id64, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	id := int32(id64)

	if err := h.queries.SoftDeleteProduct(ctx, id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) ListByCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	categoryID64, err := strconv.ParseInt(vars["categoryId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	categoryID := int32(categoryID64)

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

	products, err := h.queries.ListProductsByCategory(ctx, db.ListProductsByCategoryParams{
		CategoryID: sql.NullInt32{Int32: categoryID, Valid: true},
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	respondJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) ListBelowReorderPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, err := h.queries.ListProductsBelowReorderPoint(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	respondJSON(w, http.StatusOK, products)
}