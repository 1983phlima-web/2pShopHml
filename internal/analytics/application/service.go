package application

import (
	"context"
	"time"

	"github.com/2pshop/2pshop/internal/analytics/domain"
	"github.com/2pshop/2pshop/internal/analytics/ports"
	"github.com/2pshop/2pshop/pkg/errors"
)

type Service struct {
	repo ports.Repository
}

func NewService(repo ports.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SellerSummary(ctx context.Context, tenantID string) (*domain.SellerSummary, error) {
	summary, err := s.repo.SellerSummary(ctx, tenantID, 14)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to compute seller summary", err)
	}
	return summary, nil
}

func (s *Service) AdminSummary(ctx context.Context, tenantID string) (*domain.AdminSummary, error) {
	summary, err := s.repo.AdminSummary(ctx, tenantID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to compute admin summary", err)
	}
	return summary, nil
}

func (s *Service) Health(ctx context.Context, tenantID string) (*domain.HealthReport, error) {
	report, err := s.repo.Health(ctx, tenantID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to compute health report", err)
	}
	return report, nil
}

// HealthHistory returns the sampled health curve for the given period
// ("24h", "7d", or "30d" — defaults to 24h for anything else).
func (s *Service) HealthHistory(ctx context.Context, period string) ([]domain.HealthPoint, error) {
	var since time.Time
	now := time.Now().UTC()
	switch period {
	case "7d":
		since = now.AddDate(0, 0, -7)
	case "30d":
		since = now.AddDate(0, 0, -30)
	default:
		since = now.Add(-24 * time.Hour)
	}
	points, err := s.repo.HealthHistory(ctx, since)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to load health history", err)
	}
	return points, nil
}
