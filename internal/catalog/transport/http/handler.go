package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/2pshop/2pshop/internal/catalog/application"
	"github.com/2pshop/2pshop/internal/catalog/domain"
	"github.com/2pshop/2pshop/internal/platform/httpmw"
	"github.com/2pshop/2pshop/pkg/errors"
	"github.com/2pshop/2pshop/pkg/pagination"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

// PublicRoutes: browsing the catalog requires no authentication.
func (h *Handler) PublicRoutes(r chi.Router) {
	r.Get("/products", h.List)
	r.Get("/products/facets", h.Facets)
	r.Get("/products/{id}", h.Get)
}

// ManageRoutes: creating/publishing products is restricted to
// Seller/SystemAdmin/GlobalAdmin (enforced by the caller via RequireRole).
func (h *Handler) ManageRoutes(r chi.Router) {
	r.Post("/products", h.Create)
	r.Post("/products/{id}/publish", h.Publish)
}

type createProductRequest struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	SKU   string `json:"sku"`
	Price int64  `json:"price"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	var req createProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("reason", "invalid json"))
		return
	}
	product, err := h.service.CreateProduct(r.Context(), tenantID, req.Name, req.Slug, req.SKU, req.Price)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, product)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	id := chi.URLParam(r, "id")
	product, err := h.service.GetProduct(r.Context(), tenantID, id)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, product)
}

func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	id := chi.URLParam(r, "id")
	if err := h.service.PublishProduct(r.Context(), tenantID, id); err != nil {
		respondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type listResponse struct {
	Data  []*domain.Product `json:"data"`
	Total int               `json:"total"`
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	q := r.URL.Query()
	params := pagination.FromQuery(q.Get("cursor"), q.Get("limit"), q.Get("sort"), q.Get("order"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	minPrice, _ := strconv.ParseInt(q.Get("min_price"), 10, 64)
	maxPrice, _ := strconv.ParseInt(q.Get("max_price"), 10, 64)

	filter := domain.ListFilter{
		Query:        q.Get("q"),
		CategorySlug: q.Get("category"),
		Brand:        q.Get("brand"),
		Gender:       q.Get("gender"),
		MinPrice:     minPrice,
		MaxPrice:     maxPrice,
	}

	products, total, err := h.service.ListFiltered(r.Context(), tenantID, filter, params.Limit, offset)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, listResponse{Data: products, Total: total})
}

func (h *Handler) Facets(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	facets, err := h.service.GetFacets(r.Context(), tenantID)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, facets)
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
