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
	// ListByTenant returns stock for every product that has an inventory
	// row, joined with the product name — powers the seller's stock panel.
	ListByTenant(ctx context.Context, tenantID string) ([]domain.StockItem, error)
}
