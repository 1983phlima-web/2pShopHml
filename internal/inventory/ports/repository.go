package ports

import (
	"context"
	"time"

	"github.com/2pshop/2pshop/internal/inventory/domain"
)

type Repository interface {
	GetStock(ctx context.Context, tenantID, productID string) (*domain.Stock, error)
	Reserve(ctx context.Context, reservation *domain.Reservation) error
	Release(ctx context.Context, tenantID, reservationID string) error
	UpdateStock(ctx context.Context, tenantID, productID string, quantity int) error
	ExpireReservations(ctx context.Context, before time.Time) error
}
