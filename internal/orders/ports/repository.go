package ports

import (
	"context"

	"github.com/2pshop/2pshop/internal/orders/domain"
)

type Repository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, tenantID, id string) (*domain.Order, error)
	UpdateStatus(ctx context.Context, tenantID, id string, status domain.Status) error
	ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*domain.Order, error)
	ListByCustomer(ctx context.Context, tenantID, customerID string, limit, offset int) ([]*domain.Order, error)
}

type EventPublisher interface {
	PublishOrderCreated(ctx context.Context, order *domain.Order) error
	PublishOrderUpdated(ctx context.Context, order *domain.Order) error
}
