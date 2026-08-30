package postgres

import (
	"context"
	"encoding/json"

	"github.com/2pshop/2pshop/internal/orders/domain"
	"github.com/2pshop/2pshop/internal/orders/ports"
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

func (r *Repository) Create(ctx context.Context, o *domain.Order) error {
	items, err := json.Marshal(o.Items)
	if err != nil {
		return err
	}
	_, err = r.db.Pool.Exec(ctx, `
		INSERT INTO orders (id, tenant_id, customer_id, items, total, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, o.ID, o.TenantID, o.CustomerID, items, o.Total, o.Status, o.CreatedAt, o.UpdatedAt)
	return err
}

func (r *Repository) scan(row pgx.Row) (*domain.Order, error) {
	var o domain.Order
	var items []byte
	err := row.Scan(&o.ID, &o.TenantID, &o.CustomerID, &items, &o.Total, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if len(items) > 0 {
		_ = json.Unmarshal(items, &o.Items)
	}
	return &o, nil
}

const selectCols = `id, tenant_id, customer_id, items, total, status, created_at, updated_at`

func (r *Repository) GetByID(ctx context.Context, tenantID, id string) (*domain.Order, error) {
	row := r.db.Pool.QueryRow(ctx, `SELECT `+selectCols+` FROM orders WHERE tenant_id=$1 AND id=$2`, tenantID, id)
	return r.scan(row)
}

func (r *Repository) UpdateStatus(ctx context.Context, tenantID, id string, status domain.Status) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE orders SET status=$3, updated_at=NOW() WHERE tenant_id=$1 AND id=$2`, tenantID, id, status)
	return err
}

func (r *Repository) ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*domain.Order, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT `+selectCols+` FROM orders WHERE tenant_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		o, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

// ListByCustomer returns orders for a single customer — used for the
// "Meus Pedidos" (my purchases) screen.
func (r *Repository) ListByCustomer(ctx context.Context, tenantID, customerID string, limit, offset int) ([]*domain.Order, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT `+selectCols+` FROM orders WHERE tenant_id=$1 AND customer_id=$2 ORDER BY created_at DESC LIMIT $3 OFFSET $4
	`, tenantID, customerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		o, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}
