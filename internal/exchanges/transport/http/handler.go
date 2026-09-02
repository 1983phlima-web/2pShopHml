package http

import (
	"encoding/json"
	"net/http"

	"github.com/2pshop/2pshop/internal/exchanges/application"
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

// Routes: requires an authenticated user (mounted behind RequireAuth).
func (h *Handler) Routes(r chi.Router) {
	r.Post("/exchanges", h.Create)
	r.Get("/exchanges/mine", h.ListMine)
}

type createRequest struct {
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Reason    string `json:"reason"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	claims := httpmw.ClaimsFromContext(r.Context())
	if claims == nil {
		respondError(w, errors.New(errors.ErrUnauthorized))
		return
	}
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("reason", "invalid json"))
		return
	}
	e, err := h.service.Request(r.Context(), tenantID, req.OrderID, req.ProductID, claims.UserID, req.Reason)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, e)
}

func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	claims := httpmw.ClaimsFromContext(r.Context())
	if claims == nil {
		respondError(w, errors.New(errors.ErrUnauthorized))
		return
	}
	list, err := h.service.ListMine(r.Context(), tenantID, claims.UserID, 50, 0)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, list)
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
