package ports

import (
	"context"

	"github.com/2pshop/2pshop/internal/identity/domain"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, tenantID, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, tenantID, email string) (*domain.User, error)
	ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*domain.User, error)
	UpdateAvatar(ctx context.Context, tenantID, userID, avatar string) error
}

type TokenService interface {
	Generate(ctx context.Context, user *domain.User) (string, error)
	Validate(ctx context.Context, token string) (*TokenClaims, error)
}

type TokenClaims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
