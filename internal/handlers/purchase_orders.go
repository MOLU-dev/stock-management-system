// ============================================================================
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	db "github.com/molu/stock-management-system/internal/db/sqlc"
	"github.com/shopspring/decimal"
)

type PurchaseOrderHandler struct {
	queries *db.Queries
}

func NewPurchaseOrderHandler(queries *db.Queries) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{queries: queries}
}

type CreatePurchaseOrderRequest struct {
	PONumber             string          `json:"po_number"`
	SupplierID           int64           `json:"supplier_id"`
	OrderDate            time.Time       `json:"order_date"`
	ExpectedDeliveryDate time.Time       `json:"expected_delivery_date"`
	Status               string          `json:"status"`
	TotalAmount          decimal.Decimal `json:"total_amount"`
	Notes                *string         `json:"notes"`
	CreatedBy            int64           `json:"created_by"`
}

func (h *PurchaseOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var req CreatePurchaseOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	po, err := h.queries.CreatePurchaseOrder(ctx, db.CreatePurchaseOrderParams{
		PoNumber:             req.PONumber,
		SupplierID:           req.SupplierID,
		OrderDate:            req.OrderDate,
		ExpectedDeliveryDate: req.ExpectedDeliveryDate,
		Status:               db.PurchaseOrderStatus(req.Status),
		TotalAmount:          req.TotalAmount,
		Notes:                req.Notes,
		CreatedBy:            req.CreatedBy,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create purchase order")
		return
	}

	respondJSON(w, http.StatusCreated, po)
}

func (h *PurchaseOrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid PO ID")
		return
	}

	po, err := h.queries.GetPurchaseOrder(ctx, id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Purchase order not found")
		return
	}

	respondJSON(w, http.StatusOK, po)
}

func (h *PurchaseOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := int32(50)
	offset := int32(0)

	orders, err := h.queries.ListPurchaseOrders(ctx, db.ListPurchaseOrdersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch purchase orders")
		return
	}

	respondJSON(w, http.StatusOK, orders)
}

func (h *PurchaseOrderHandler) ListByStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	limit := int32(50)
	offset := int32(0)

	orders, err := h.queries.ListPurchaseOrdersByStatus(ctx, db.ListPurchaseOrdersByStatusParams{
		Status: db.PurchaseOrderStatus(vars["status"]),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch purchase orders")
		return
	}

	respondJSON(w, http.StatusOK, orders)
}

type UpdatePurchaseOrderStatusRequest struct {
	Status      string          `json:"status"`
	TotalAmount decimal.Decimal `json:"total_amount"`
}

func (h *PurchaseOrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid PO ID")
		return
	}

	var req UpdatePurchaseOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	po, err := h.queries.UpdatePurchaseOrderStatus(ctx, db.UpdatePurchaseOrderStatusParams{
		PoID:        id,
		Status:      db.PurchaseOrderStatus(req.Status),
		TotalAmount: req.TotalAmount,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update purchase order")
		return
	}

	respondJSON(w, http.StatusOK, po)
}

func (h *PurchaseOrderHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid PO ID")
		return
	}

	items, err := h.queries.GetPurchaseOrderItems(ctx, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch items")
		return
	}

	respondJSON(w, http.StatusOK, items)
}

type CreatePurchaseOrderItemRequest struct {
	ProductID        int64           `json:"product_id"`
	QuantityOrdered  int32           `json:"quantity_ordered"`
	QuantityReceived int32           `json:"quantity_received"`
	UnitPrice        decimal.Decimal `json:"unit_price"`
	TotalPrice       decimal.Decimal `json:"total_price"`
}

func (h *PurchaseOrderHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	
	poID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid PO ID")
		return
	}

	var req CreatePurchaseOrderItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.queries.CreatePurchaseOrderItem(ctx, db.CreatePurchaseOrderItemParams{
		PoID:             poID,
		ProductID:        req.ProductID,
		QuantityOrdered:  req.QuantityOrdered,
		QuantityReceived: req.QuantityReceived,
		UnitPrice:        req.UnitPrice,
		TotalPrice:       req.TotalPrice,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	respondJSON(w, http.StatusCreated, item)
}

type ReceiveItemRequest struct {
	Quantity int32 `json:"quantity"`
}

func (h *PurchaseOrderHandler) ReceiveItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	
	itemID, err := strconv.ParseInt(vars["itemId"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var req ReceiveItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.queries.UpdatePurchaseOrderItemReceivedQty(ctx, db.UpdatePurchaseOrderItemReceivedQtyParams{
		PoItemID: itemID,
		Column2:  req.Quantity,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to receive item")
		return
	}

	respondJSON(w, http.StatusOK, item)
}