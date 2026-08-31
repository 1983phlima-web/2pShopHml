package http

import (
	"encoding/json"
	"net/http"

	"github.com/2pshop/2pshop/internal/analytics/application"
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

// SellerRoutes: Seller/SystemAdmin/GlobalAdmin can see sales analytics.
func (h *Handler) SellerRoutes(r chi.Router) {
	r.Get("/analytics/seller-summary", h.SellerSummary)
}

// AdminRoutes: SystemAdmin/GlobalAdmin only.
func (h *Handler) AdminRoutes(r chi.Router) {
	r.Get("/analytics/admin-summary", h.AdminSummary)
}

// GlobalRoutes: GlobalAdmin only — infra-level health signals.
func (h *Handler) GlobalRoutes(r chi.Router) {
	r.Get("/analytics/health", h.Health)
}

func (h *Handler) SellerSummary(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	summary, err := h.service.SellerSummary(r.Context(), tenantID)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, summary)
}

func (h *Handler) AdminSummary(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	summary, err := h.service.AdminSummary(r.Context(), tenantID)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, summary)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	report, err := h.service.Health(r.Context(), tenantID)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, report)
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
