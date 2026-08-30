package application

import (
	"context"

	"github.com/2pshop/2pshop/internal/reviews/domain"
	"github.com/2pshop/2pshop/internal/reviews/ports"
	"github.com/2pshop/2pshop/pkg/errors"
)

type Service struct {
	repo ports.Repository
}

func NewService(repo ports.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddReview(ctx context.Context, tenantID, productID, userID, userName, comment string, rating int) (*domain.Review, error) {
	if comment == "" {
		return nil, errors.New(errors.ErrInvalidInput).WithDetail("field", "comment")
	}
	review := domain.NewReview(tenantID, productID, userID, userName, comment, rating)
	if err := s.repo.Create(ctx, review); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to create review", err)
	}
	return review, nil
}

func (s *Service) ListReviews(ctx context.Context, tenantID, productID string, limit, offset int) ([]*domain.Review, domain.Summary, error) {
	reviews, err := s.repo.ListByProduct(ctx, tenantID, productID, limit, offset)
	if err != nil {
		return nil, domain.Summary{}, errors.Wrap(errors.ErrInternal, "failed to list reviews", err)
	}
	summary, err := s.repo.Summary(ctx, tenantID, productID)
	if err != nil {
		return nil, domain.Summary{}, errors.Wrap(errors.ErrInternal, "failed to summarize reviews", err)
	}
	return reviews, summary, nil
}
