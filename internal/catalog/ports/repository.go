package ports

import (
	"context"

	"github.com/2pshop/2pshop/internal/catalog/domain"
)

type Repository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, tenantID, id string) (*domain.Product, error)
	GetBySlug(ctx context.Context, tenantID, slug string) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	List(ctx context.Context, tenantID string, state domain.PublicationState, limit, offset int) ([]*domain.Product, error)
	CountByTenant(ctx context.Context, tenantID string) (int, error)
}

type EventPublisher interface {
	PublishProductCreated(ctx context.Context, product *domain.Product) error
	PublishProductUpdated(ctx context.Context, product *domain.Product) error
}
