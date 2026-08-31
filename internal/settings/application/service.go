package application

import (
	"context"

	"github.com/2pshop/2pshop/internal/settings/ports"
	"github.com/2pshop/2pshop/pkg/errors"
)

type Service struct {
	repo ports.Repository
}

func NewService(repo ports.Repository) *Service {
	return &Service{repo: repo}
}

var defaultTheme = map[string]any{"palette": "indigo"}

// GetTheme never fails the caller with a 404 — a tenant without a saved
// preference simply gets the default palette.
func (s *Service) GetTheme(ctx context.Context, tenantID string) (map[string]any, error) {
	setting, err := s.repo.Get(ctx, tenantID, "theme")
	if err != nil {
		if errors.IsNotFound(err) {
			return defaultTheme, nil
		}
		return nil, errors.Wrap(errors.ErrInternal, "failed to load theme", err)
	}
	if setting.Value == nil || setting.Value["palette"] == nil {
		return defaultTheme, nil
	}
	return setting.Value, nil
}

func (s *Service) SetTheme(ctx context.Context, tenantID, palette string) error {
	if err := s.repo.Set(ctx, tenantID, "theme", map[string]any{"palette": palette}); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save theme", err)
	}
	return nil
}
