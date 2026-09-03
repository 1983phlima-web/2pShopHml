package domain

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	ProductID string    `json:"product_id"`
	VariantID string    `json:"variant_id,omitempty"`
	Quantity  int       `json:"quantity"`
	OrderID   string    `json:"order_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Stock struct {
	TenantID  string    `json:"tenant_id"`
	ProductID string    `json:"product_id"`
	VariantID string    `json:"variant_id,omitempty"`
	Quantity  int       `json:"quantity"`
	Reserved  int       `json:"reserved"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StockItem is a Stock row enriched with the product name — used by the
// seller's stock management panel so it doesn't need a second fetch.
type StockItem struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Reserved    int    `json:"reserved"`
}

func (s *Stock) Available() int {
	avail := s.Quantity - s.Reserved
	if avail < 0 {
		return 0
	}
	return avail
}

func NewReservation(tenantID, productID, orderID string, quantity int, duration time.Duration) *Reservation {
	now := time.Now().UTC()
	return &Reservation{
		ID:        uuid.Must(uuid.NewV7()).String(),
		TenantID:  tenantID,
		ProductID: productID,
		OrderID:   orderID,
		Quantity:  quantity,
		ExpiresAt: now.Add(duration),
		CreatedAt: now,
	}
}
