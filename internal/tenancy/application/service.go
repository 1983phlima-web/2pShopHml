package application

import (
	"context"
	"fmt"

	"github.com/2pshop/2pshop/internal/tenancy/domain"
	"github.com/2pshop/2pshop/internal/tenancy/ports"
	"github.com/2pshop/2pshop/pkg/errors"
)

type Service struct {
	repo ports.Repository
}

func NewService(repo ports.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTenant(ctx context.Context, name, slug, planID string) (*domain.Tenant, error) {
	existing, err := s.repo.GetBySlug(ctx, slug)
	if err != nil && !errors.IsNotFound(err) {
		return nil, errors.Wrap(errors.ErrInternal, "failed to check existing tenant", err)
	}
	if existing != nil {
		return nil, errors.New(errors.ErrConflict).WithDetail("field", "slug")
	}

	tenant := domain.NewTenant(name, slug, planID)
	if err := s.repo.Create(ctx, tenant); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to create tenant", err)
	}
	return tenant, nil
}

func (s *Service) GetTenant(ctx context.Context, id string) (*domain.Tenant, error) {
	tenant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.New(errors.ErrNotFound).WithDetail("resource", "tenant")
		}
		return nil, errors.Wrap(errors.ErrInternal, "failed to get tenant", err)
	}
	return tenant, nil
}

func (s *Service) ValidateTenant(ctx context.Context, tenantID string) error {
	tenant, err := s.repo.GetByID(ctx, tenantID)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "tenant validation failed", err)
	}
	if tenant == nil {
		return errors.New(errors.ErrNotFound).WithDetail("resource", "tenant")
	}
	if !tenant.IsActive() {
		return errors.New(errors.ErrTenantSuspended).WithDetail("status", string(tenant.Status))
	}
	return nil
}
