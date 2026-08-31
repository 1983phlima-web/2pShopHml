package application

import (
	"context"

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
