package postgres

import (
	"context"
	"encoding/json"

	"github.com/2pshop/2pshop/internal/catalog/domain"
	"github.com/2pshop/2pshop/internal/catalog/ports"
	"github.com/2pshop/2pshop/internal/platform"
	"github.com/2pshop/2pshop/pkg/errors"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *platform.DB
}

func NewRepository(db *platform.DB) ports.Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, p *domain.Product) error {
	attrs, err := json.Marshal(p.Attributes)
	if err != nil {
		return err
	}
	seo, err := json.Marshal(p.SEO)
	if err != nil {
		return err
	}
	_, err = r.db.Pool.Exec(ctx, `
		INSERT INTO products (id, tenant_id, name, slug, description, sku, price, state, category_id, attributes, seo, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NULLIF($9, '')::uuid, $10, $11, $12, $13)
	`, p.ID, p.TenantID, p.Name, p.Slug, p.Description, p.SKU, p.Price, p.State, p.CategoryID, attrs, seo, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *Repository) scan(row pgx.Row) (*domain.Product, error) {
	var p domain.Product
	var attrs, seo []byte
	var categoryID *string
	err := row.Scan(&p.ID, &p.TenantID, &p.Name, &p.Slug, &p.Description, &p.SKU, &p.Price, &p.State, &categoryID, &attrs, &seo, &p.CreatedAt, &p.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound)
	}
	if err != nil {
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
	return &p, nil
}

const selectCols = `id, tenant_id, name, slug, description, sku, price, state, category_id, attributes, seo, created_at, updated_at`

func (r *Repository) GetByID(ctx context.Context, tenantID, id string) (*domain.Product, error) {
	row := r.db.Pool.QueryRow(ctx, `SELECT `+selectCols+` FROM products WHERE tenant_id = $1 AND id = $2`, tenantID, id)
	return r.scan(row)
}

func (r *Repository) GetBySlug(ctx context.Context, tenantID, slug string) (*domain.Product, error) {
	row := r.db.Pool.QueryRow(ctx, `SELECT `+selectCols+` FROM products WHERE tenant_id = $1 AND slug = $2`, tenantID, slug)
	p, err := r.scan(row)
	if errors.IsNotFound(err) {
		return nil, nil
	}
	return p, err
}

func (r *Repository) Update(ctx context.Context, p *domain.Product) error {
	attrs, _ := json.Marshal(p.Attributes)
	seo, _ := json.Marshal(p.SEO)
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE products SET name=$3, description=$4, price=$5, state=$6, attributes=$7, seo=$8, updated_at=$9
		WHERE tenant_id=$1 AND id=$2
	`, p.TenantID, p.ID, p.Name, p.Description, p.Price, p.State, attrs, seo, p.UpdatedAt)
	return err
}

func (r *Repository) List(ctx context.Context, tenantID string, state domain.PublicationState, limit, offset int) ([]*domain.Product, error) {
	var rows pgx.Rows
	var err error
	if state == "" {
		rows, err = r.db.Pool.Query(ctx, `SELECT `+selectCols+` FROM products WHERE tenant_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, tenantID, limit, offset)
	} else {
		rows, err = r.db.Pool.Query(ctx, `SELECT `+selectCols+` FROM products WHERE tenant_id=$1 AND state=$2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`, tenantID, state, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *Repository) CountByTenant(ctx context.Context, tenantID string) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE tenant_id=$1`, tenantID).Scan(&count)
	return count, err
}
