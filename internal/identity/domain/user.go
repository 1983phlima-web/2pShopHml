package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

// Role represents one of the 4 RBAC roles of the marketplace.
type Role string

const (
	// RoleSeller (Vendedor): manages own products, inventory and orders.
	RoleSeller Role = "SELLER"
	// RoleBuyer (Cliente): browses products, places orders, writes reviews.
	RoleBuyer Role = "BUYER"
	// RoleSystemAdmin (Administrador do Sistema): system parametrization,
	// manages catalog/categories/settings across the tenant.
	RoleSystemAdmin Role = "SYSTEM_ADMIN"
	// RoleGlobalAdmin (Administrador Global): full access, including
	// database-level operations and payment/Stripe integration config.
	RoleGlobalAdmin Role = "GLOBAL_ADMIN"
)

// IsAdmin returns true for either of the two administrative roles.
func (r Role) IsAdmin() bool {
	return r == RoleSystemAdmin || r == RoleGlobalAdmin
}

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

// SetPassword hashes and stores the given plaintext password.
func (u *User) SetPassword(plain string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// CheckPassword reports whether the given plaintext password matches
// the stored hash.
func (u *User) CheckPassword(plain string) bool {
	if u.Password == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain)) == nil
}
