package application

import (
	"context"

	"github.com/2pshop/2pshop/internal/exchanges/domain"
	"github.com/2pshop/2pshop/internal/exchanges/ports"
	"github.com/2pshop/2pshop/pkg/errors"
)

type Service struct {
	repo ports.Repository
}

func NewService(repo ports.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Request(ctx context.Context, tenantID, orderID, productID, userID, reason string) (*domain.ExchangeRequest, error) {
	if reason == "" {
		return nil, errors.New(errors.ErrInvalidInput).WithDetail("field", "reason")
	}
	e := domain.NewExchangeRequest(tenantID, orderID, productID, userID, reason)
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to create exchange request", err)
	}
	return e, nil
}

func (s *Service) ListMine(ctx context.Context, tenantID, userID string, limit, offset int) ([]*domain.ExchangeRequest, error) {
	list, err := s.repo.ListByUser(ctx, tenantID, userID, limit, offset)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to list exchange requests", err)
	}
	return list, nil
}
