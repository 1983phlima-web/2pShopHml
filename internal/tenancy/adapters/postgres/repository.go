package postgres

import (
	"context"
	"fmt"

	"github.com/2pshop/2pshop/internal/platform"
	"github.com/2pshop/2pshop/internal/tenancy/domain"
	"github.com/2pshop/2pshop/internal/tenancy/ports"
	"github.com/2pshop/2pshop/pkg/errors"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *platform.DB
}

func NewRepository(db *platform.DB) ports.Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, tenant *domain.Tenant) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO tenants (id, name, slug, status, plan_id, limits, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, tenant.ID, tenant.Name, tenant.Slug, tenant.Status, tenant.PlanID, tenant.Limits, tenant.CreatedAt, tenant.UpdatedAt)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	var limits domain.Limits
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, name, slug, status, plan_id, limits, created_at, updated_at
		FROM tenants WHERE id = $1
	`, id).Scan(&tenant.ID, &tenant.Name, &tenant.Slug, &tenant.Status, &tenant.PlanID, &limits, &tenant.CreatedAt, &tenant.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	tenant.Limits = limits
	return &tenant, nil
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	var limits domain.Limits
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, name, slug, status, plan_id, limits, created_at, updated_at
		FROM tenants WHERE slug = $1
	`, slug).Scan(&tenant.ID, &tenant.Name, &tenant.Slug, &tenant.Status, &tenant.PlanID, &limits, &tenant.CreatedAt, &tenant.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	tenant.Limits = limits
	return &tenant, nil
}

func (r *Repository) Update(ctx context.Context, tenant *domain.Tenant) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE tenants SET name = $2, status = $3, plan_id = $4, limits = $5, updated_at = $6
		WHERE id = $1
	`, tenant.ID, tenant.Name, tenant.Status, tenant.PlanID, tenant.Limits, tenant.UpdatedAt)
	return err
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]*domain.Tenant, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, name, slug, status, plan_id, limits, created_at, updated_at
		FROM tenants ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []*domain.Tenant
	for rows.Next() {
		var t domain.Tenant
		var limits domain.Limits
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.Status, &t.PlanID, &limits, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		t.Limits = limits
		tenants = append(tenants, &t)
	}
	return tenants, rows.Err()
}
