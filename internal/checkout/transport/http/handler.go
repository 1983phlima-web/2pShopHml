package http

import (
	"encoding/json"
	"net/http"

	"github.com/2pshop/2pshop/internal/checkout/application"
	paymentDomain "github.com/2pshop/2pshop/internal/payments/domain"
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

func (h *Handler) Routes(r chi.Router) {
	r.Post("/checkout", h.Checkout)
}

type checkoutRequest struct {
	Items         []application.CheckoutItem  `json:"items"`
	PaymentMethod paymentDomain.PaymentMethod `json:"payment_method"`
}

func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	claims := httpmw.ClaimsFromContext(r.Context())
	if claims == nil {
		respondError(w, errors.New(errors.ErrUnauthorized))
		return
	}

	var body checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("reason", "invalid json"))
		return
	}
	if len(body.Items) == 0 {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("field", "items"))
		return
	}

	req := application.CheckoutRequest{
		TenantID:       tenantID,
		CustomerID:     claims.UserID,
		Items:          body.Items,
		PaymentMethod:  body.PaymentMethod,
		IdempotencyKey: r.Header.Get("Idempotency-Key"),
	}

	result, err := h.service.Checkout(r.Context(), req)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, result)
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
