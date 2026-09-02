package application

import (
	"context"

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

func (s *Service) CreateUser(ctx context.Context, tenantID, email, name, password string, role domain.Role) (*domain.User, error) {
	existing, err := s.repo.GetByEmail(ctx, tenantID, email)
	if err != nil && !errors.IsNotFound(err) {
		return nil, errors.Wrap(errors.ErrInternal, "failed to check existing user", err)
	}
	if existing != nil {
		return nil, errors.New(errors.ErrConflict).WithDetail("field", "email")
	}

	user := domain.NewUser(tenantID, email, name, role)
	if err := user.SetPassword(password); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to hash password", err)
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to create user", err)
	}
	return user, nil
}

func (s *Service) Authenticate(ctx context.Context, tenantID, email, password string) (string, *domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, tenantID, email)
	if err != nil {
		if errors.IsNotFound(err) {
			return "", nil, errors.New(errors.ErrUnauthorized).WithDetail("reason", "invalid credentials")
		}
		return "", nil, errors.Wrap(errors.ErrInternal, "authentication failed", err)
	}
	if !user.Active {
		return "", nil, errors.New(errors.ErrUnauthorized).WithDetail("reason", "user inactive")
	}
	if !user.CheckPassword(password) {
		return "", nil, errors.New(errors.ErrUnauthorized).WithDetail("reason", "invalid credentials")
	}

	token, err := s.token.Generate(ctx, user)
	if err != nil {
		return "", nil, errors.Wrap(errors.ErrInternal, "token generation failed", err)
	}
	return token, user, nil
}

func (s *Service) GetUser(ctx context.Context, tenantID, id string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.New(errors.ErrNotFound).WithDetail("resource", "user")
		}
		return nil, errors.Wrap(errors.ErrInternal, "failed to get user", err)
	}
	return user, nil
}

// maxAvatarLength guards against abusive uploads — ~700KB of base64 is
// generous for a profile picture while keeping row size sane.
const maxAvatarLength = 700_000

func (s *Service) UpdateAvatar(ctx context.Context, tenantID, userID, avatar string) error {
	if avatar == "" {
		return errors.New(errors.ErrInvalidInput).WithDetail("field", "avatar")
	}
	if len(avatar) > maxAvatarLength {
		return errors.New(errors.ErrInvalidInput).WithDetail("reason", "avatar too large")
	}
	if err := s.repo.UpdateAvatar(ctx, tenantID, userID, avatar); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to update avatar", err)
	}
	return nil
}
