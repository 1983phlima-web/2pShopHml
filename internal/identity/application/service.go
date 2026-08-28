package application

import (
	"context"
	"fmt"

	"github.com/2pshop/2pshop/internal/identity/domain"
	"github.com/2pshop/2pshop/internal/identity/ports"
	"github.com/2pshop/2pshop/pkg/errors"
)

type Service struct {
	repo  ports.Repository
	token ports.TokenService
}

func NewService(repo ports.Repository, token ports.TokenService) *Service {
	return &Service{repo: repo, token: token}
}

func (s *Service) CreateUser(ctx context.Context, tenantID, email, name string, role domain.Role) (*domain.User, error) {
	existing, err := s.repo.GetByEmail(ctx, tenantID, email)
	if err != nil && !errors.IsNotFound(err) {
		return nil, errors.Wrap(errors.ErrInternal, "failed to check existing user", err)
	}
	if existing != nil {
		return nil, errors.New(errors.ErrConflict).WithDetail("field", "email")
	}

	user := domain.NewUser(tenantID, email, name, role)
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to create user", err)
	}
	return user, nil
}

func (s *Service) Authenticate(ctx context.Context, tenantID, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, tenantID, email)
	if err != nil {
		if errors.IsNotFound(err) {
			return "", errors.New(errors.ErrUnauthorized).WithDetail("reason", "invalid credentials")
		}
		return "", errors.Wrap(errors.ErrInternal, "authentication failed", err)
	}
	if !user.Active {
		return "", errors.New(errors.ErrUnauthorized).WithDetail("reason", "user inactive")
	}

	// TODO: password hash comparison
	_ = password

	token, err := s.token.Generate(ctx, user)
	if err != nil {
		return "", errors.Wrap(errors.ErrInternal, "token generation failed", err)
	}
	return token, nil
}
