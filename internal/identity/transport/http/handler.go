package http

import (
	"encoding/json"
	"net/http"

	"github.com/2pshop/2pshop/internal/identity/application"
	"github.com/2pshop/2pshop/internal/identity/domain"
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

// Routes registers the public auth routes (register/login).
func (h *Handler) Routes(r chi.Router) {
	r.Post("/auth/register", h.Register)
	r.Post("/auth/login", h.Login)
}

// ProtectedRoutes registers routes that require a valid Bearer token.
func (h *Handler) ProtectedRoutes(r chi.Router) {
	r.Get("/auth/me", h.Me)
	r.Put("/auth/me/avatar", h.UpdateAvatar)
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

type userResponse struct {
	ID     string      `json:"id"`
	Email  string      `json:"email"`
	Name   string      `json:"name"`
	Role   domain.Role `json:"role"`
	Avatar string      `json:"avatar"`
	Phone  string      `json:"phone,omitempty"`
	Active bool        `json:"active"`
}

// Register creates a new customer (BUYER) account. Privileged roles
// (Seller/SystemAdmin/GlobalAdmin) are provisioned separately, not through
// public self-signup.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("reason", "invalid json"))
		return
	}
	if req.Email == "" || req.Password == "" || req.Name == "" {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("reason", "name, email and password are required"))
		return
	}

	user, err := h.service.CreateUser(r.Context(), tenantID, req.Email, req.Name, req.Password, domain.RoleBuyer)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, toUserResponse(user))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	tenantID := httpmw.TenantFromContext(r.Context())
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("reason", "invalid json"))
		return
	}

	token, user, err := h.service.Authenticate(r.Context(), tenantID, req.Email, req.Password)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, authResponse{Token: token, User: toUserResponse(user)})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims := httpmw.ClaimsFromContext(r.Context())
	if claims == nil {
		respondError(w, errors.New(errors.ErrUnauthorized))
		return
	}
	user, err := h.service.GetUser(r.Context(), claims.TenantID, claims.UserID)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, toUserResponse(user))
}

type updateAvatarRequest struct {
	Avatar string `json:"avatar"`
}

func (h *Handler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	claims := httpmw.ClaimsFromContext(r.Context())
	if claims == nil {
		respondError(w, errors.New(errors.ErrUnauthorized))
		return
	}
	var req updateAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, errors.New(errors.ErrInvalidInput).WithDetail("reason", "invalid json"))
		return
	}
	if err := h.service.UpdateAvatar(r.Context(), claims.TenantID, claims.UserID, req.Avatar); err != nil {
		respondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toUserResponse(u *domain.User) userResponse {
	return userResponse{ID: u.ID, Email: u.Email, Name: u.Name, Role: u.Role, Avatar: u.Avatar, Phone: u.Phone, Active: u.Active}
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
