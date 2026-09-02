package ports

import (
	"context"

	"github.com/2pshop/2pshop/internal/questions/domain"
)

type Repository interface {
	Create(ctx context.Context, q *domain.Question) error
	ListByProduct(ctx context.Context, tenantID, productID string, limit, offset int) ([]*domain.Question, error)
	ListByUser(ctx context.Context, tenantID, userID string, limit, offset int) ([]*domain.Question, error)
	Answer(ctx context.Context, tenantID, id, answer string) error
}
