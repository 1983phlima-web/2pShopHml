package ports

import (
	"context"

	"github.com/2pshop/2pshop/internal/exchanges/domain"
)

type Repository interface {
	Create(ctx context.Context, e *domain.ExchangeRequest) error
	ListByUser(ctx context.Context, tenantID, userID string, limit, offset int) ([]*domain.ExchangeRequest, error)
}
