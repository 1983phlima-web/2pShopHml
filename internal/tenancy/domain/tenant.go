package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusTrial      Status = "TRIAL"
	StatusActive     Status = "ACTIVE"
	StatusPastDue    Status = "PAST_DUE"
	StatusSuspended  Status = "SUSPENDED"
	StatusCancelled  Status = "CANCELLED"
	StatusExpired    Status = "EXPIRED"
)

type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Status    Status    `json:"status"`
	PlanID    string    `json:"plan_id"`
	Limits    Limits    `json:"limits"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Limits struct {
	MaxProducts   int `json:"max_products"`
	MaxOrders     int `json:"max_orders_monthly"`
	MaxUsers      int `json:"max_users"`
	MaxStorageGB  int `json:"max_storage_gb"`
	APIRateLimit  int `json:"api_rate_limit"`
}

func NewTenant(name, slug, planID string) *Tenant {
	now := time.Now().UTC()
	return &Tenant{
		ID:        uuid.Must(uuid.NewV7()).String(),
		Name:      name,
		Slug:      slug,
		Status:    StatusTrial,
		PlanID:    planID,
		Limits:    Limits{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (t *Tenant) IsActive() bool {
	return t.Status == StatusActive || t.Status == StatusTrial
}

func (t *Tenant) CanCreateProduct(currentCount int) bool {
	if t.Limits.MaxProducts <= 0 {
		return true
	}
	return currentCount < t.Limits.MaxProducts
}
