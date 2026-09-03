package http

import (
	"encoding/json"
	"net/http"

	"github.com/2pshop/2pshop/internal/inventory/application"
	"github.com/2pshop/2pshop/internal/platform/httpmw"
	"github.com/2pshop/2pshop/pkg/errors"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

// SellerRoutes: stock management is restricted to Seller/Admin roles.
func (h *Handler) SellerRoutes(r chi.Router) {
	r.Get("/inventory", h.List)
	r.Put("/inventory/{productId}", h.SetStock)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	items, err := h.service.ListStock(r.Context(), tenantID)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, items)
}

type setStockRequest struct {
	Quantity int `json:"quantity"`
}

func (h *Handler) SetStock(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	productID := chi.URLParam(r, "productId")
	var req setStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("reason", "invalid json"))
		return
	}
	if err := h.service.SetStock(r.Context(), tenantID, productID, req.Quantity); err != nil {
		respondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, err error) {
	status := errors.HTTPStatus(err)
	var appErr *errors.AppError
	if ae, ok := err.(*errors.AppError); ok {
		appErr = ae
	} else {
		appErr = errors.New(errors.ErrInternal)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(appErr)
}
