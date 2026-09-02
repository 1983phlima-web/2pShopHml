package http

import (
	"encoding/json"
	"net/http"

	"github.com/2pshop/2pshop/internal/platform/httpmw"
	"github.com/2pshop/2pshop/internal/questions/application"
	"github.com/2pshop/2pshop/pkg/errors"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

// PublicRoutes: reading a product's Q&A requires no authentication.
func (h *Handler) PublicRoutes(r chi.Router) {
	r.Get("/products/{id}/questions", h.ListByProduct)
}

// ProtectedRoutes: asking a question, and listing "my questions", require
// an authenticated user.
func (h *Handler) ProtectedRoutes(r chi.Router) {
	r.Post("/products/{id}/questions", h.Ask)
	r.Get("/questions/mine", h.ListMine)
}

// SellerRoutes: answering a question is restricted to Seller/Admin roles.
func (h *Handler) SellerRoutes(r chi.Router) {
	r.Post("/questions/{id}/answer", h.Answer)
}

type askRequest struct {
	Question string `json:"question"`
}

func (h *Handler) Ask(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	productID := chi.URLParam(r, "id")
	claims := httpmw.ClaimsFromContext(r.Context())
	if claims == nil {
		respondError(w, errors.New(errors.ErrUnauthorized))
		return
	}
	var req askRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("reason", "invalid json"))
		return
	}
	q, err := h.service.Ask(r.Context(), tenantID, productID, claims.UserID, claims.Email, req.Question)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, q)
}

func (h *Handler) ListByProduct(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	productID := chi.URLParam(r, "id")
	questions, err := h.service.ListByProduct(r.Context(), tenantID, productID, 50, 0)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, questions)
}

func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	claims := httpmw.ClaimsFromContext(r.Context())
	if claims == nil {
		respondError(w, errors.New(errors.ErrUnauthorized))
		return
	}
	questions, err := h.service.ListByUser(r.Context(), tenantID, claims.UserID, 50, 0)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, questions)
}

type answerRequest struct {
	Answer string `json:"answer"`
}

func (h *Handler) Answer(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	id := chi.URLParam(r, "id")
	var req answerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("reason", "invalid json"))
		return
	}
	if err := h.service.Answer(r.Context(), tenantID, id, req.Answer); err != nil {
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
