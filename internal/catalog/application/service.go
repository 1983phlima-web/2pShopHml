package application

import (
	"context"

	"github.com/2pshop/2pshop/internal/catalog/domain"
	"github.com/2pshop/2pshop/internal/catalog/ports"
	"github.com/2pshop/2pshop/pkg/errors"
)

type Service struct {
	repo      ports.Repository
	publisher ports.EventPublisher
}

func NewService(repo ports.Repository, publisher ports.EventPublisher) *Service {
	return &Service{repo: repo, publisher: publisher}
}

func (s *Service) CreateProduct(ctx context.Context, tenantID, name, slug, sku string, price int64) (*domain.Product, error) {
	existing, err := s.repo.GetBySlug(ctx, tenantID, slug)
	if err != nil && !errors.IsNotFound(err) {
		return nil, errors.Wrap(errors.ErrInternal, "failed to check product", err)
	}
	if existing != nil {
		return nil, errors.New(errors.ErrConflict).WithDetail("field", "slug")
	}

	product := domain.NewProduct(tenantID, name, slug, sku, price)
	if err := s.repo.Create(ctx, product); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to create product", err)
	}

	if s.publisher != nil {
		_ = s.publisher.PublishProductCreated(ctx, product)
	}
	return product, nil
}

func (s *Service) GetProduct(ctx context.Context, tenantID, id string) (*domain.Product, error) {
	product, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.New(errors.ErrNotFound).WithDetail("resource", "product")
		}
		return nil, errors.Wrap(errors.ErrInternal, "failed to get product", err)
	}
	return product, nil
}

func (s *Service) PublishProduct(ctx context.Context, tenantID, id string) error {
	product, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	product.Publish()
	if err := s.repo.Update(ctx, product); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to publish product", err)
	}
	if s.publisher != nil {
		_ = s.publisher.PublishProductUpdated(ctx, product)
	}
	return nil
}
