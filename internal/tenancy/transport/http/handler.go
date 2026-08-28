package http

import (
	"encoding/json"
	"net/http"

	"github.com/2pshop/2pshop/internal/tenancy/application"
	"github.com/2pshop/2pshop/internal/tenancy/domain"
	"github.com/2pshop/2pshop/pkg/errors"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes(r chi.Router) {
	r.Post("/tenants", h.Create)
	r.Get("/tenants/{id}", h.Get)
}

type createTenantRequest struct {
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	PlanID string `json:"plan_id"`
}

type tenantResponse struct {
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	Slug   string        `json:"slug"`
	Status domain.Status `json:"status"`
	PlanID string        `json:"plan_id"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("reason", "invalid json"))
		return
	}

	tenant, err := h.service.CreateTenant(r.Context(), req.Name, req.Slug, req.PlanID)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, toResponse(tenant))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant, err := h.service.GetTenant(r.Context(), id)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, toResponse(tenant))
}

func toResponse(t *domain.Tenant) tenantResponse {
	return tenantResponse{
		ID:     t.ID,
		Name:   t.Name,
		Slug:   t.Slug,
		Status: t.Status,
		PlanID: t.PlanID,
	}
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
