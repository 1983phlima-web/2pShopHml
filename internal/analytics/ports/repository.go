package ports

import (
	"context"

	"github.com/2pshop/2pshop/internal/analytics/domain"
)

type Repository interface {
	SellerSummary(ctx context.Context, tenantID string, days int) (*domain.SellerSummary, error)
	AdminSummary(ctx context.Context, tenantID string) (*domain.AdminSummary, error)
	Health(ctx context.Context, tenantID string) (*domain.HealthReport, error)
}
