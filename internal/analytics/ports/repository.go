package ports

import (
	"context"
	"time"

	"github.com/2pshop/2pshop/internal/analytics/domain"
)

type Repository interface {
	SellerSummary(ctx context.Context, tenantID string, days int) (*domain.SellerSummary, error)
	AdminSummary(ctx context.Context, tenantID string) (*domain.AdminSummary, error)
	Health(ctx context.Context, tenantID string) (*domain.HealthReport, error)
	// RecordSnapshot samples current DB latency + platform counts and
	// persists a health_snapshots row — called periodically by a
	// background job so the health history has real data to chart.
	RecordSnapshot(ctx context.Context) error
	HealthHistory(ctx context.Context, since time.Time) ([]domain.HealthPoint, error)
}
