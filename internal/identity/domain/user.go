package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      Role      `json:"role"`
	Password  string    `json:"-"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Role string

const (
	RoleAdmin   Role = "ADMIN"
	RoleSeller  Role = "SELLER"
	RoleBuyer   Role = "BUYER"
	RoleSupport Role = "SUPPORT"
)

func NewUser(tenantID, email, name string, role Role) *User {
	now := time.Now().UTC()
	return &User{
		ID:        uuid.Must(uuid.NewV7()).String(),
		TenantID:  tenantID,
		Email:     email,
		Name:      name,
		Role:      role,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
