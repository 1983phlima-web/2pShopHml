package application

import (
	"context"

	catalogDomain "github.com/2pshop/2pshop/internal/catalog/domain"
	"github.com/2pshop/2pshop/internal/favorites/ports"
	"github.com/2pshop/2pshop/pkg/errors"
)

type Service struct {
	repo ports.Repository
}

func NewService(repo ports.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Add(ctx context.Context, tenantID, userID, productID string) error {
	if err := s.repo.Add(ctx, tenantID, userID, productID); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to favorite product", err)
	}
	return nil
}

func (s *Service) Remove(ctx context.Context, tenantID, userID, productID string) error {
	if err := s.repo.Remove(ctx, tenantID, userID, productID); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to unfavorite product", err)
	}
	return nil
}

func (s *Service) ListProductIDs(ctx context.Context, tenantID, userID string) ([]string, error) {
	ids, err := s.repo.ListProductIDs(ctx, tenantID, userID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to list favorite ids", err)
	}
	return ids, nil
}

func (s *Service) ListWithProducts(ctx context.Context, tenantID, userID string) ([]*catalogDomain.Product, error) {
	products, err := s.repo.ListWithProducts(ctx, tenantID, userID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to list favorites", err)
	}
	return products, nil
}
