package domain

import (
	"time"

	"github.com/google/uuid"
)

type PublicationState string

const (
	StateDraft    PublicationState = "DRAFT"
	StateActive   PublicationState = "ACTIVE"
	StateInactive PublicationState = "INACTIVE"
	StateArchived PublicationState = "ARCHIVED"
)

type Product struct {
	ID          string           `json:"id"`
	TenantID    string           `json:"tenant_id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description string           `json:"description"`
	SKU         string           `json:"sku"`
	Price       int64            `json:"price"` // cents
	State       PublicationState `json:"state"`
	CategoryID  string           `json:"category_id"`
	BrandID     string           `json:"brand_id,omitempty"`
	Attributes  map[string]any   `json:"attributes,omitempty"`
	SEO         SEO              `json:"seo,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type SEO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
}

func NewProduct(tenantID, name, slug, sku string, price int64) *Product {
	now := time.Now().UTC()
	return &Product{
		ID:         uuid.Must(uuid.NewV7()).String(),
		TenantID:   tenantID,
		Name:       name,
		Slug:       slug,
		SKU:        sku,
		Price:      price,
		State:      StateDraft,
		Attributes: make(map[string]any),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (p *Product) Publish() {
	p.State = StateActive
	p.UpdatedAt = time.Now().UTC()
}

func (p *Product) Archive() {
	p.State = StateArchived
	p.UpdatedAt = time.Now().UTC()
}
