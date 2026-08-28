package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending    Status = "PENDING"
	StatusConfirmed  Status = "CONFIRMED"
	StatusPaid       Status = "PAID"
	StatusShipped    Status = "SHIPPED"
	StatusDelivered  Status = "DELIVERED"
	StatusCancelled  Status = "CANCELLED"
	StatusRefunded   Status = "REFUNDED"
)

type Order struct {
	ID         string      `json:"id"`
	TenantID   string      `json:"tenant_id"`
	CustomerID string      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
	Total      int64       `json:"total"` // cents
	Status     Status      `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	UnitPrice int64  `json:"unit_price"`
	Total     int64  `json:"total"`
}

func NewOrder(tenantID, customerID string, items []OrderItem) *Order {
	now := time.Now().UTC()
	var total int64
	for _, item := range items {
		item.Total = int64(item.Quantity) * item.UnitPrice
		total += item.Total
	}
	return &Order{
		ID:         uuid.Must(uuid.NewV7()).String(),
		TenantID:   tenantID,
		CustomerID: customerID,
		Items:      items,
		Total:      total,
		Status:     StatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (o *Order) Confirm() {
	o.Status = StatusConfirmed
	o.UpdatedAt = time.Now().UTC()
}

func (o *Order) Cancel() {
	o.Status = StatusCancelled
	o.UpdatedAt = time.Now().UTC()
}
