package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	db "github.com/molu/stock-management-system/internal/db/sqlc"
)

type CategoryHandler struct {
	queries db.SingleDb
}

func NewCategoryHandler(queries db.SingleDb) *CategoryHandler {
	return &CategoryHandler{queries: queries}
}

// Request/Response types
type CreateCategoryRequest struct {
	CategoryCode     string  `json:"category_code"`
	Name             string  `json:"name"`
	ParentCategoryID *int64  `json:"parent_category_id"`
	Description      *string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name             string  `json:"name"`
	ParentCategoryID *int64  `json:"parent_category_id"`
	Description      *string `json:"description"`
}

// List retrieves all categories with pagination
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := int32(50)
	offset := int32(0)

	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	if offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 32); err == nil {
			offset = int32(o)
		}
	}

	categories, err := h.queries.ListCategories(ctx, db.ListCategoriesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		log.Printf("Error listing categories: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}

	respondJSON(w, http.StatusOK, categories)
}

// ListRoot retrieves all root categories (no parent)
func (h *CategoryHandler) ListRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.queries.ListRootCategories(ctx)
	if err != nil {
		log.Printf("Error listing root categories: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to fetch root categories")
		return
	}

	respondJSON(w, http.StatusOK, categories)
}

// ListSubCategories retrieves all subcategories of a parent category
func (h *CategoryHandler) ListSubCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	parentID, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid parent category ID")
		return
	}

	categories, err := h.queries.ListSubCategories(ctx, sql.NullInt32{
		Int32: int32(parentID),
		Valid: true,
	})
	if err != nil {
		log.Printf("Error listing subcategories: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to fetch subcategories")
		return
	}

	respondJSON(w, http.StatusOK, categories)
}

// Get retrieves a single category by ID
func (h *CategoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.queries.GetCategory(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "Category not found")
		} else {
			log.Printf("Error getting category: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to fetch category")
		}
		return
	}

	respondJSON(w, http.StatusOK, category)
}

// GetByCode retrieves a category by its code
func (h *CategoryHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	code := vars["code"]
	if code == "" {
		respondError(w, http.StatusBadRequest, "Category code is required")
		return
	}

	category, err := h.queries.GetCategoryByCode(ctx, code)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "Category not found")
		} else {
			log.Printf("Error getting category by code: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to fetch category")
		}
		return
	}

	respondJSON(w, http.StatusOK, category)
}

// Create creates a new category
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.CategoryCode == "" || req.Name == "" {
		respondError(w, http.StatusBadRequest, "Category code and name are required")
		return
	}

	// Prepare parameters
	params := db.CreateCategoryParams{
		CategoryCode: req.CategoryCode,
		Name:         req.Name,
		ParentCategoryID: sql.NullInt32{
			Int32: 0,
			Valid: false,
		},
		Description: sql.NullString{
			String: "",
			Valid:  false,
		},
	}

	if req.ParentCategoryID != nil {
		params.ParentCategoryID = sql.NullInt32{
			Int32: int32(*req.ParentCategoryID),
			Valid: true,
		}
	}

	if req.Description != nil && *req.Description != "" {
		params.Description = sql.NullString{
			String: *req.Description,
			Valid:  true,
		}
	}

	category, err := h.queries.CreateCategory(ctx, params)
	if err != nil {
		log.Printf("Error creating category: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	respondJSON(w, http.StatusCreated, category)
}

// Update updates an existing category
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var req UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "Category name is required")
		return
	}

	// Prepare parameters
	params := db.UpdateCategoryParams{
		CategoryID: int32(id),
		Name:       req.Name,
		ParentCategoryID: sql.NullInt32{
			Int32: 0,
			Valid: false,
		},
		Description: sql.NullString{
			String: "",
			Valid:  false,
		},
	}

	if req.ParentCategoryID != nil {
		params.ParentCategoryID = sql.NullInt32{
			Int32: int32(*req.ParentCategoryID),
			Valid: true,
		}
	}

	if req.Description != nil && *req.Description != "" {
		params.Description = sql.NullString{
			String: *req.Description,
			Valid:  true,
		}
	}

	category, err := h.queries.UpdateCategory(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "Category not found")
		} else {
			log.Printf("Error updating category: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to update category")
		}
		return
	}

	respondJSON(w, http.StatusOK, category)
}

// Delete deletes a category
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	err = h.queries.DeleteCategory(ctx, int32(id))
	if err != nil {
		log.Printf("Error deleting category: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}
