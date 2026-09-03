package postgres

import (
	"context"

	"github.com/2pshop/2pshop/internal/exchanges/domain"
	"github.com/2pshop/2pshop/internal/exchanges/ports"
	"github.com/2pshop/2pshop/internal/platform"
)

type Repository struct {
	db *platform.DB
}

func NewRepository(db *platform.DB) ports.Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, e *domain.ExchangeRequest) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO exchange_requests (id, tenant_id, order_id, product_id, user_id, reason, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, e.ID, e.TenantID, e.OrderID, e.ProductID, e.UserID, e.Reason, e.Status, e.CreatedAt, e.UpdatedAt)
	return err
}

func (r *Repository) ListByUser(ctx context.Context, tenantID, userID string, limit, offset int) ([]*domain.ExchangeRequest, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, tenant_id, order_id, product_id, user_id, reason, status, created_at, updated_at
		FROM exchange_requests WHERE tenant_id=$1 AND user_id=$2 ORDER BY created_at DESC LIMIT $3 OFFSET $4
	`, tenantID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.ExchangeRequest
	for rows.Next() {
		var e domain.ExchangeRequest
		if err := rows.Scan(&e.ID, &e.TenantID, &e.OrderID, &e.ProductID, &e.UserID, &e.Reason, &e.Status, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	return list, rows.Err()
}

// ListByTenant returns every exchange request in the tenant — the
// seller/admin approval queue.
func (r *Repository) ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*domain.ExchangeRequest, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, tenant_id, order_id, product_id, user_id, reason, status, created_at, updated_at
		FROM exchange_requests WHERE tenant_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.ExchangeRequest
	for rows.Next() {
		var e domain.ExchangeRequest
		if err := rows.Scan(&e.ID, &e.TenantID, &e.OrderID, &e.ProductID, &e.UserID, &e.Reason, &e.Status, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	return list, rows.Err()
}

func (r *Repository) UpdateStatus(ctx context.Context, tenantID, id string, status domain.Status) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE exchange_requests SET status = $3, updated_at = NOW() WHERE tenant_id = $1 AND id = $2
	`, tenantID, id, status)
	return err
}
