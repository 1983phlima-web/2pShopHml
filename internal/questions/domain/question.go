package domain

import (
	"time"

	"github.com/google/uuid"
)

// Question is a customer's question on an announced product — one half
// of "Meus Comentários" (the other half being reviews/testemunhos).
type Question struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	ProductID  string     `json:"product_id"`
	UserID     string     `json:"user_id"`
	UserName   string     `json:"user_name"`
	Question   string     `json:"question"`
	Answer     string     `json:"answer,omitempty"`
	AnsweredAt *time.Time `json:"answered_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

func NewQuestion(tenantID, productID, userID, userName, question string) *Question {
	return &Question{
		ID:        uuid.Must(uuid.NewV7()).String(),
		TenantID:  tenantID,
		ProductID: productID,
		UserID:    userID,
		UserName:  userName,
		Question:  question,
		CreatedAt: time.Now().UTC(),
	}
}
