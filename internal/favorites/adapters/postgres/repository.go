package postgres

import (
	"context"
	"encoding/json"

	catalogDomain "github.com/2pshop/2pshop/internal/catalog/domain"
	"github.com/2pshop/2pshop/internal/favorites/ports"
	"github.com/2pshop/2pshop/internal/platform"
)

type Repository struct {
	db *platform.DB
}

func NewRepository(db *platform.DB) ports.Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(ctx context.Context, tenantID, userID, productID string) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO favorites (id, tenant_id, user_id, product_id, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, NOW())
		ON CONFLICT (tenant_id, user_id, product_id) DO NOTHING
	`, tenantID, userID, productID)
	return err
}

func (r *Repository) Remove(ctx context.Context, tenantID, userID, productID string) error {
	_, err := r.db.Pool.Exec(ctx, `
		DELETE FROM favorites WHERE tenant_id = $1 AND user_id = $2 AND product_id = $3
	`, tenantID, userID, productID)
	return err
}

func (r *Repository) ListProductIDs(ctx context.Context, tenantID, userID string) ([]string, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT product_id FROM favorites WHERE tenant_id = $1 AND user_id = $2 ORDER BY created_at DESC
	`, tenantID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *Repository) ListWithProducts(ctx context.Context, tenantID, userID string) ([]*catalogDomain.Product, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT p.id, p.tenant_id, p.name, p.slug, p.description, p.sku, p.price, p.state, p.category_id, p.attributes, p.seo, p.created_at, p.updated_at
		FROM favorites f
		JOIN products p ON p.id = f.product_id AND p.tenant_id = f.tenant_id
		WHERE f.tenant_id = $1 AND f.user_id = $2
		ORDER BY f.created_at DESC
	`, tenantID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*catalogDomain.Product
	for rows.Next() {
		var p catalogDomain.Product
		var attrs, seo []byte
		var categoryID *string
		if err := rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.Slug, &p.Description, &p.SKU, &p.Price, &p.State, &categoryID, &attrs, &seo, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if categoryID != nil {
			p.CategoryID = *categoryID
		}
		if len(attrs) > 0 {
			_ = json.Unmarshal(attrs, &p.Attributes)
		}
		if len(seo) > 0 {
			_ = json.Unmarshal(seo, &p.SEO)
		}
		products = append(products, &p)
	}
	return products, rows.Err()
}
