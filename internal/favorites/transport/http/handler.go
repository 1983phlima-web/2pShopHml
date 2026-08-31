package http

import (
	"encoding/json"
	"net/http"

	"github.com/2pshop/2pshop/internal/favorites/application"
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

// Routes: every favorites route requires an authenticated user (mounted
// behind RequireAuth by the caller).
func (h *Handler) Routes(r chi.Router) {
	r.Get("/favorites", h.List)
	r.Get("/favorites/ids", h.ListIDs)
	r.Post("/favorites/{productId}", h.Add)
	r.Delete("/favorites/{productId}", h.Remove)
}

func (h *Handler) currentUser(r *http.Request) (tenantID, userID string, ok bool) {
	tenantID = httpmw.TenantFromContext(r.Context())
	claims := httpmw.ClaimsFromContext(r.Context())
	if claims == nil {
		return "", "", false
	}
	return tenantID, claims.UserID, true
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.currentUser(r)
	if !ok {
		respondError(w, errors.New(errors.ErrUnauthorized))
		return
	}
	products, err := h.service.ListWithProducts(r.Context(), tenantID, userID)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, products)
}

func (h *Handler) ListIDs(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.currentUser(r)
	if !ok {
		respondError(w, errors.New(errors.ErrUnauthorized))
		return
	}
	ids, err := h.service.ListProductIDs(r.Context(), tenantID, userID)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, ids)
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.currentUser(r)
	if !ok {
		respondError(w, errors.New(errors.ErrUnauthorized))
		return
	}
	productID := chi.URLParam(r, "productId")
	if err := h.service.Add(r.Context(), tenantID, userID, productID); err != nil {
		respondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.currentUser(r)
	if !ok {
		respondError(w, errors.New(errors.ErrUnauthorized))
		return
	}
	productID := chi.URLParam(r, "productId")
	if err := h.service.Remove(r.Context(), tenantID, userID, productID); err != nil {
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
