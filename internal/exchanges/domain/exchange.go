package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusRequested Status = "REQUESTED"
	StatusApproved  Status = "APPROVED"
	StatusRejected  Status = "REJECTED"
	StatusCompleted Status = "COMPLETED"
)

// ExchangeRequest is a customer's post-delivery return/exchange request
// for an item from a past order — "Minhas Trocas".
type ExchangeRequest struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	OrderID   string    `json:"order_id"`
	ProductID string    `json:"product_id"`
	UserID    string    `json:"user_id"`
	Reason    string    `json:"reason"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewExchangeRequest(tenantID, orderID, productID, userID, reason string) *ExchangeRequest {
	now := time.Now().UTC()
	return &ExchangeRequest{
		ID:        uuid.Must(uuid.NewV7()).String(),
		TenantID:  tenantID,
		OrderID:   orderID,
		ProductID: productID,
		UserID:    userID,
		Reason:    reason,
		Status:    StatusRequested,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
