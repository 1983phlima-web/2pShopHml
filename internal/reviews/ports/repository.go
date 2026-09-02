package ports

import (
	"context"

	"github.com/2pshop/2pshop/internal/reviews/domain"
)

type Repository interface {
	Create(ctx context.Context, review *domain.Review) error
	ListByProduct(ctx context.Context, tenantID, productID string, limit, offset int) ([]*domain.Review, error)
	ListByUser(ctx context.Context, tenantID, userID string, limit, offset int) ([]*domain.Review, error)
	Summary(ctx context.Context, tenantID, productID string) (domain.Summary, error)
}
