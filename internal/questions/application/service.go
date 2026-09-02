package application

import (
	"context"

	"github.com/2pshop/2pshop/internal/questions/domain"
	"github.com/2pshop/2pshop/internal/questions/ports"
	"github.com/2pshop/2pshop/pkg/errors"
)

type Service struct {
	repo ports.Repository
}

func NewService(repo ports.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Ask(ctx context.Context, tenantID, productID, userID, userName, question string) (*domain.Question, error) {
	if question == "" {
		return nil, errors.New(errors.ErrInvalidInput).WithDetail("field", "question")
	}
	q := domain.NewQuestion(tenantID, productID, userID, userName, question)
	if err := s.repo.Create(ctx, q); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to create question", err)
	}
	return q, nil
}

func (s *Service) ListByProduct(ctx context.Context, tenantID, productID string, limit, offset int) ([]*domain.Question, error) {
	questions, err := s.repo.ListByProduct(ctx, tenantID, productID, limit, offset)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to list questions", err)
	}
	return questions, nil
}

func (s *Service) ListByUser(ctx context.Context, tenantID, userID string, limit, offset int) ([]*domain.Question, error) {
	questions, err := s.repo.ListByUser(ctx, tenantID, userID, limit, offset)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to list my questions", err)
	}
	return questions, nil
}

func (s *Service) Answer(ctx context.Context, tenantID, id, answer string) error {
	if answer == "" {
		return errors.New(errors.ErrInvalidInput).WithDetail("field", "answer")
	}
	if err := s.repo.Answer(ctx, tenantID, id, answer); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to answer question", err)
	}
	return nil
}
