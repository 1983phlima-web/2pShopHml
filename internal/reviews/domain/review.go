package domain

import (
	"time"

	"github.com/google/uuid"
)

// Review represents a customer's comment and rating on a product —
// the "comentários de compra" shown on the product page.
type Review struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	ProductID string    `json:"product_id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	Rating    int       `json:"rating"` // 1-5
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

func NewReview(tenantID, productID, userID, userName, comment string, rating int) *Review {
	if rating < 1 {
		rating = 1
	}
	if rating > 5 {
		rating = 5
	}
	return &Review{
		ID:        uuid.Must(uuid.NewV7()).String(),
		TenantID:  tenantID,
		ProductID: productID,
		UserID:    userID,
		UserName:  userName,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: time.Now().UTC(),
	}
}

// Summary aggregates a product's reviews for display (average + count).
type Summary struct {
	ProductID string  `json:"product_id"`
	Average   float64 `json:"average"`
	Count     int     `json:"count"`
}
