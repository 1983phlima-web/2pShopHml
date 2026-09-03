package http

import (
	"encoding/json"
	"net/http"

	"github.com/2pshop/2pshop/internal/orders/application"
	"github.com/2pshop/2pshop/internal/orders/domain"
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

// Routes: every order route requires authentication (mounted behind
// RequireAuth by the caller).
func (h *Handler) Routes(r chi.Router) {
	r.Get("/orders", h.MyOrders)
	r.Get("/orders/{id}", h.Get)
}

// SellerRoutes: the fulfillment queue and status transitions are
// restricted to Seller/Admin roles.
func (h *Handler) SellerRoutes(r chi.Router) {
	r.Get("/orders/all", h.AllOrders)
	r.Put("/orders/{id}/status", h.UpdateStatus)
}

func (h *Handler) AllOrders(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	orders, err := h.service.ListAllOrders(r.Context(), tenantID, 100, 0)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, orders)
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	id := chi.URLParam(r, "id")
	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("reason", "invalid json"))
		return
	}
	if err := h.service.UpdateStatus(r.Context(), tenantID, id, domain.Status(req.Status)); err != nil {
		respondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) MyOrders(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	claims := httpmw.ClaimsFromContext(r.Context())
	if claims == nil {
		respondError(w, errors.New(errors.ErrUnauthorized))
		return
	}
	orders, err := h.service.ListMyOrders(r.Context(), tenantID, claims.UserID, 50, 0)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, orders)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	id := chi.URLParam(r, "id")
	order, err := h.service.GetOrder(r.Context(), tenantID, id)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, order)
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
