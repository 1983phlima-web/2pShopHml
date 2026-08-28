package application

import (
	"context"
	"time"

	"github.com/2pshop/2pshop/internal/inventory/domain"
	"github.com/2pshop/2pshop/internal/inventory/ports"
	"github.com/2pshop/2pshop/pkg/errors"
)

type Service struct {
	repo ports.Repository
}

func NewService(repo ports.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Reserve(ctx context.Context, tenantID, productID, orderID string, quantity int) (*domain.Reservation, error) {
	stock, err := s.repo.GetStock(ctx, tenantID, productID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to get stock", err)
	}
	if stock == nil {
		return nil, errors.New(errors.ErrNotFound).WithDetail("resource", "stock")
	}
	if stock.Available() < quantity {
		return nil, errors.New(errors.ErrInventoryUnavailable).
			WithDetail("available", stock.Available()).
			WithDetail("requested", quantity)
	}

	reservation := domain.NewReservation(tenantID, productID, orderID, quantity, 15*time.Minute)
	if err := s.repo.Reserve(ctx, reservation); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to reserve inventory", err)
	}
	return reservation, nil
}

func (s *Service) Release(ctx context.Context, tenantID, reservationID string) error {
	if err := s.repo.Release(ctx, tenantID, reservationID); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to release inventory", err)
	}
	return nil
}
