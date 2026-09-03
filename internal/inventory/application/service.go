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

// ListStock returns stock for every product in the tenant's catalog —
// powers the seller's stock management panel.
func (s *Service) ListStock(ctx context.Context, tenantID string) ([]domain.StockItem, error) {
	items, err := s.repo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to list stock", err)
	}
	return items, nil
}

// SetStock sets the absolute quantity for a product (not a delta).
func (s *Service) SetStock(ctx context.Context, tenantID, productID string, quantity int) error {
	if quantity < 0 {
		return errors.New(errors.ErrInvalidInput).WithDetail("field", "quantity")
	}
	if err := s.repo.UpdateStock(ctx, tenantID, productID, quantity); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to update stock", err)
	}
	return nil
}
