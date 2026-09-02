package http

import (
	"encoding/json"
	"net/http"

	"github.com/2pshop/2pshop/internal/platform/httpmw"
	"github.com/2pshop/2pshop/internal/reviews/application"
	"github.com/2pshop/2pshop/pkg/errors"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

// PublicRoutes: reading reviews requires no authentication.
func (h *Handler) PublicRoutes(r chi.Router) {
	r.Get("/products/{id}/reviews", h.List)
}

// ProtectedRoutes: posting a review requires an authenticated buyer.
func (h *Handler) ProtectedRoutes(r chi.Router) {
	r.Post("/products/{id}/reviews", h.Create)
	r.Get("/reviews/mine", h.ListMine)
}

type createReviewRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type listReviewsResponse struct {
	Data    []reviewResponse `json:"data"`
	Average float64          `json:"average"`
	Count   int              `json:"count"`
}

type reviewResponse struct {
	ID          string `json:"id"`
	ProductID   string `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	UserName    string `json:"user_name"`
	Rating      int    `json:"rating"`
	Comment     string `json:"comment"`
	CreatedAt   string `json:"created_at"`
}

func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	claims := httpmw.ClaimsFromContext(r.Context())
	if claims == nil {
		respondError(w, errors.New(errors.ErrUnauthorized))
		return
	}
	reviews, err := h.service.ListMyReviews(r.Context(), tenantID, claims.UserID, 50, 0)
	if err != nil {
		respondError(w, err)
		return
	}
	data := make([]reviewResponse, 0, len(reviews))
	for _, rv := range reviews {
		data = append(data, reviewResponse{
			ID:          rv.ID,
			ProductID:   rv.ProductID,
			ProductName: rv.ProductName,
			UserName:    rv.UserName,
			Rating:      rv.Rating,
			Comment:     rv.Comment,
			CreatedAt:   rv.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	respondJSON(w, http.StatusOK, data)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	productID := chi.URLParam(r, "id")
	claims := httpmw.ClaimsFromContext(r.Context())
	if claims == nil {
		respondError(w, errors.New(errors.ErrUnauthorized))
		return
	}

	var req createReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("reason", "invalid json"))
		return
	}

	review, err := h.service.AddReview(r.Context(), tenantID, productID, claims.UserID, claims.Email, req.Comment, req.Rating)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, review)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	productID := chi.URLParam(r, "id")

	reviews, summary, err := h.service.ListReviews(r.Context(), tenantID, productID, 50, 0)
	if err != nil {
		respondError(w, err)
		return
	}

	data := make([]reviewResponse, 0, len(reviews))
	for _, rv := range reviews {
		data = append(data, reviewResponse{
			ID:        rv.ID,
			UserName:  rv.UserName,
			Rating:    rv.Rating,
			Comment:   rv.Comment,
			CreatedAt: rv.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	respondJSON(w, http.StatusOK, listReviewsResponse{Data: data, Average: summary.Average, Count: summary.Count})
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
