package postgres

import (
	"context"
	"time"

	"github.com/2pshop/2pshop/internal/inventory/domain"
	"github.com/2pshop/2pshop/internal/inventory/ports"
	"github.com/2pshop/2pshop/internal/platform"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *platform.DB
}

func NewRepository(db *platform.DB) ports.Repository {
	return &Repository{db: db}
}

func (r *Repository) GetStock(ctx context.Context, tenantID, productID string) (*domain.Stock, error) {
	var s domain.Stock
	err := r.db.Pool.QueryRow(ctx, `
		SELECT tenant_id, product_id, quantity, reserved, updated_at
		FROM inventory WHERE tenant_id=$1 AND product_id=$2 AND variant_id IS NULL
	`, tenantID, productID).Scan(&s.TenantID, &s.ProductID, &s.Quantity, &s.Reserved, &s.UpdatedAt)
	if err == pgx.ErrNoRows {
		// No inventory row yet means the product simply has no stock tracked.
		return &domain.Stock{TenantID: tenantID, ProductID: productID, Quantity: 0, Reserved: 0}, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) Reserve(ctx context.Context, reservation *domain.Reservation) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		INSERT INTO reservations (id, tenant_id, product_id, quantity, order_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, reservation.ID, reservation.TenantID, reservation.ProductID, reservation.Quantity, reservation.OrderID, reservation.ExpiresAt, reservation.CreatedAt); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO inventory (id, tenant_id, product_id, quantity, reserved, updated_at)
		VALUES (gen_random_uuid(), $1, $2, 0, $3, NOW())
		ON CONFLICT (tenant_id, product_id) WHERE variant_id IS NULL
		DO UPDATE SET reserved = inventory.reserved + $3, updated_at = NOW()
	`, reservation.TenantID, reservation.ProductID, reservation.Quantity); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) Release(ctx context.Context, tenantID, reservationID string) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var productID string
	var quantity int
	err = tx.QueryRow(ctx, `
		SELECT product_id, quantity FROM reservations WHERE tenant_id=$1 AND id=$2
	`, tenantID, reservationID).Scan(&productID, &quantity)
	if err == pgx.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM reservations WHERE tenant_id=$1 AND id=$2`, tenantID, reservationID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE inventory SET reserved = GREATEST(reserved - $3, 0), updated_at = NOW()
		WHERE tenant_id=$1 AND product_id=$2 AND variant_id IS NULL
	`, tenantID, productID, quantity); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repository) UpdateStock(ctx context.Context, tenantID, productID string, quantity int) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO inventory (id, tenant_id, product_id, quantity, reserved, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, 0, NOW())
		ON CONFLICT (tenant_id, product_id) WHERE variant_id IS NULL
		DO UPDATE SET quantity = $3, updated_at = NOW()
	`, tenantID, productID, quantity)
	return err
}

func (r *Repository) ExpireReservations(ctx context.Context, before time.Time) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM reservations WHERE expires_at < $1`, before)
	return err
}
