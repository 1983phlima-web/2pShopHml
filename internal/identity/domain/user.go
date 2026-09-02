package domain

import (
	"math/rand"
	"strconv"
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
	Avatar    string    `json:"avatar"` // "preset:N" (1-12) or a base64 data URI
	Phone     string    `json:"phone,omitempty"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PresetAvatarCount is the number of built-in avatar options offered at
// signup and in the profile popup.
const PresetAvatarCount = 12

// randomPresetAvatar assigns a random built-in avatar to every new user,
// so nobody starts with a blank profile picture.
func randomPresetAvatar() string {
	return "preset:" + strconv.Itoa(rand.Intn(PresetAvatarCount)+1)
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
		Avatar:    randomPresetAvatar(),
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
