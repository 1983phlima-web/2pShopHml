package ports

import (
	"context"

	"github.com/2pshop/2pshop/internal/settings/domain"
)

type Repository interface {
	Get(ctx context.Context, tenantID, key string) (*domain.Setting, error)
	Set(ctx context.Context, tenantID, key string, value map[string]any) error
}
