package domain

import "time"

// Favorite links a user to a product they've saved to their personal
// "Favoritos" list.
type Favorite struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
}
