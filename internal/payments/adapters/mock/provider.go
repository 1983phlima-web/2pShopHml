// Package mock provides a sandbox payment provider for environments (like
// HML) where no real payment gateway credentials are configured. It
// implements the same domain.Provider interface a real Stripe adapter
// would, so swapping it out later requires no changes to checkout.
package mock

import (
	"context"
	"time"

	"github.com/2pshop/2pshop/internal/payments/domain"
	"github.com/google/uuid"
)

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Authorize(ctx context.Context, req domain.AuthorizeRequest) (domain.Authorization, error) {
	return domain.Authorization{
		ID:        "auth_" + uuid.NewString(),
		Status:    "AUTHORIZED",
		Amount:    req.Amount,
		Currency:  req.Currency,
		Provider:  "mock",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (p *Provider) Capture(ctx context.Context, req domain.CaptureRequest) (domain.Capture, error) {
	return domain.Capture{
		ID:        "cap_" + uuid.NewString(),
		Status:    "CAPTURED",
		Amount:    req.Amount,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (p *Provider) Refund(ctx context.Context, req domain.RefundRequest) (domain.Refund, error) {
	return domain.Refund{
		ID:        "refund_" + uuid.NewString(),
		Status:    "REFUNDED",
		Amount:    req.Amount,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}
