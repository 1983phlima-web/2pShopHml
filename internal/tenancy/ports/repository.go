package ports

import (
	"context"

	"github.com/2pshop/2pshop/internal/tenancy/domain"
)

type Repository interface {
	Create(ctx context.Context, tenant *domain.Tenant) error
	GetByID(ctx context.Context, id string) (*domain.Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error)
	Update(ctx context.Context, tenant *domain.Tenant) error
	List(ctx context.Context, limit, offset int) ([]*domain.Tenant, error)
}
