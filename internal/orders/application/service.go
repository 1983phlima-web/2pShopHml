package application

import (
	"context"

	"github.com/2pshop/2pshop/internal/orders/domain"
	"github.com/2pshop/2pshop/internal/orders/ports"
	"github.com/2pshop/2pshop/pkg/errors"
)

type Service struct {
	repo      ports.Repository
	publisher ports.EventPublisher
}

func NewService(repo ports.Repository, publisher ports.EventPublisher) *Service {
	return &Service{repo: repo, publisher: publisher}
}

func (s *Service) CreateOrder(ctx context.Context, tenantID, customerID string, items []domain.OrderItem) (*domain.Order, error) {
	order := domain.NewOrder(tenantID, customerID, items)
	if err := s.repo.Create(ctx, order); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to create order", err)
	}
	if s.publisher != nil {
		_ = s.publisher.PublishOrderCreated(ctx, order)
	}
	return order, nil
}

func (s *Service) GetOrder(ctx context.Context, tenantID, id string) (*domain.Order, error) {
	order, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.New(errors.ErrNotFound).WithDetail("resource", "order")
		}
		return nil, errors.Wrap(errors.ErrInternal, "failed to get order", err)
	}
	return order, nil
}

func (s *Service) UpdateStatus(ctx context.Context, tenantID, id string, status domain.Status) error {
	if err := s.repo.UpdateStatus(ctx, tenantID, id, status); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to update order status", err)
	}
	order, _ := s.repo.GetByID(ctx, tenantID, id)
	if order != nil && s.publisher != nil {
		_ = s.publisher.PublishOrderUpdated(ctx, order)
	}
	return nil
}

func (s *Service) ListMyOrders(ctx context.Context, tenantID, customerID string, limit, offset int) ([]*domain.Order, error) {
	orders, err := s.repo.ListByCustomer(ctx, tenantID, customerID, limit, offset)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to list orders", err)
	}
	return orders, nil
}

// ListAllOrders returns every order in the tenant — the seller/admin
// fulfillment queue (as opposed to ListMyOrders, which is customer-scoped).
func (s *Service) ListAllOrders(ctx context.Context, tenantID string, limit, offset int) ([]*domain.Order, error) {
	orders, err := s.repo.ListByTenant(ctx, tenantID, limit, offset)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to list all orders", err)
	}
	return orders, nil
}
