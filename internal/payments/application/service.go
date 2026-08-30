package application

import (
	"context"

	"github.com/2pshop/2pshop/internal/payments/domain"
	"github.com/2pshop/2pshop/pkg/errors"
)

type Service struct {
	providers map[string]domain.Provider
}

func NewService(providers map[string]domain.Provider) *Service {
	return &Service{providers: providers}
}

func (s *Service) Authorize(ctx context.Context, providerName string, req domain.AuthorizeRequest) (domain.Authorization, error) {
	provider, ok := s.providers[providerName]
	if !ok {
		return domain.Authorization{}, errors.New(errors.ErrInvalidInput).WithDetail("field", "provider")
	}
	auth, err := provider.Authorize(ctx, req)
	if err != nil {
		return domain.Authorization{}, errors.Wrap(errors.ErrPaymentDeclined, "authorization failed", err)
	}
	return auth, nil
}

func (s *Service) Capture(ctx context.Context, providerName string, req domain.CaptureRequest) (domain.Capture, error) {
	provider, ok := s.providers[providerName]
	if !ok {
		return domain.Capture{}, errors.New(errors.ErrInvalidInput).WithDetail("field", "provider")
	}
	return provider.Capture(ctx, req)
}

func (s *Service) Refund(ctx context.Context, providerName string, req domain.RefundRequest) (domain.Refund, error) {
	provider, ok := s.providers[providerName]
	if !ok {
		return domain.Refund{}, errors.New(errors.ErrInvalidInput).WithDetail("field", "provider")
	}
	return provider.Refund(ctx, req)
}
