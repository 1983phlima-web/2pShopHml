package postgres

import (
	"context"
	"time"

	"github.com/2pshop/2pshop/internal/loyalty/ports"
	"github.com/2pshop/2pshop/internal/platform"
)

type Repository struct {
	db *platform.DB
}

func NewRepository(db *platform.DB) ports.Repository {
	return &Repository{db: db}
}

func (r *Repository) LifetimeSpend(ctx context.Context, tenantID, userID string) (int64, error) {
	var total int64
	err := r.db.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(total), 0) FROM orders WHERE tenant_id = $1 AND customer_id = $2
	`, tenantID, userID).Scan(&total)
	return total, err
}

func (r *Repository) PeriodSpend(ctx context.Context, tenantID, userID string, since time.Time) (int64, error) {
	var total int64
	err := r.db.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(total), 0) FROM orders WHERE tenant_id = $1 AND customer_id = $2 AND created_at >= $3
	`, tenantID, userID, since).Scan(&total)
	return total, err
}

func (r *Repository) GrantAward(ctx context.Context, tenantID, userID, periodType, periodKey, badge string, coins int) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO loyalty_awards (id, tenant_id, user_id, period_type, period_key, badge, coins, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (tenant_id, user_id, period_type, period_key, badge) DO NOTHING
	`, tenantID, userID, periodType, periodKey, badge, coins)
	return err
}

func (r *Repository) ListAwards(ctx context.Context, tenantID, userID string) ([]ports.Award, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT badge, coins, created_at FROM loyalty_awards
		WHERE tenant_id = $1 AND user_id = $2
		ORDER BY created_at ASC
	`, tenantID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var awards []ports.Award
	for rows.Next() {
		var a ports.Award
		if err := rows.Scan(&a.Badge, &a.Coins, &a.EarnedAt); err != nil {
			return nil, err
		}
		awards = append(awards, a)
	}
	return awards, rows.Err()
}

func (r *Repository) TotalCoins(ctx context.Context, tenantID, userID string) (int64, error) {
	var total int64
	err := r.db.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(coins), 0) FROM loyalty_awards WHERE tenant_id = $1 AND user_id = $2
	`, tenantID, userID).Scan(&total)
	return total, err
}
