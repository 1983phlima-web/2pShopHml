package http

import (
	"encoding/json"
	"net/http"

	"github.com/2pshop/2pshop/internal/platform/httpmw"
	"github.com/2pshop/2pshop/internal/settings/application"
	"github.com/2pshop/2pshop/pkg/errors"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

// PublicRoutes: reading the active theme requires no authentication — the
// storefront applies it for every visitor, logged in or not.
func (h *Handler) PublicRoutes(r chi.Router) {
	r.Get("/settings/theme", h.GetTheme)
}

// AdminRoutes: only System/Global Admin can change the palette.
func (h *Handler) AdminRoutes(r chi.Router) {
	r.Put("/settings/theme", h.SetTheme)
}

func (h *Handler) GetTheme(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	theme, err := h.service.GetTheme(r.Context(), tenantID)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, theme)
}

type setThemeRequest struct {
	Palette string `json:"palette"`
}

func (h *Handler) SetTheme(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	var req setThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Palette == "" {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("field", "palette"))
		return
	}
	if err := h.service.SetTheme(r.Context(), tenantID, req.Palette); err != nil {
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
