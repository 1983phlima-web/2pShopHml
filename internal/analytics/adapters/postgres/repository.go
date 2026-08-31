package postgres

import (
	"context"
	"time"

	"github.com/2pshop/2pshop/internal/analytics/domain"
	"github.com/2pshop/2pshop/internal/analytics/ports"
	"github.com/2pshop/2pshop/internal/platform"
)

type Repository struct {
	db        *platform.DB
	startedAt time.Time
	version   string
}

func NewRepository(db *platform.DB, version string) ports.Repository {
	return &Repository{db: db, startedAt: time.Now().UTC(), version: version}
}

func (r *Repository) SellerSummary(ctx context.Context, tenantID string, days int) (*domain.SellerSummary, error) {
	s := &domain.SellerSummary{}

	if err := r.db.Pool.QueryRow(ctx, `
		SELECT COUNT(*), COALESCE(SUM(total), 0) FROM orders WHERE tenant_id = $1
	`, tenantID).Scan(&s.TotalOrders, &s.TotalRevenue); err != nil {
		return nil, err
	}

	if err := r.db.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM((item->>'quantity')::int), 0)
		FROM orders, jsonb_array_elements(items) AS item
		WHERE tenant_id = $1
	`, tenantID).Scan(&s.TotalUnitsSold); err != nil {
		return nil, err
	}

	if err := r.db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM orders WHERE tenant_id = $1 AND status IN ('PENDING', 'CONFIRMED')
	`, tenantID).Scan(&s.PendingOrders); err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, `
		WITH days AS (
			SELECT generate_series(CURRENT_DATE - ($2::int - 1), CURRENT_DATE, '1 day')::date AS day
		),
		daily AS (
			SELECT date_trunc('day', created_at)::date AS day, SUM(total) AS revenue, COUNT(*) AS orders
			FROM orders
			WHERE tenant_id = $1 AND created_at >= CURRENT_DATE - ($2::int - 1)
			GROUP BY 1
		)
		SELECT days.day, COALESCE(daily.revenue, 0), COALESCE(daily.orders, 0)
		FROM days LEFT JOIN daily ON daily.day = days.day
		ORDER BY days.day
	`, tenantID, days)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var d domain.DailySales
		var day time.Time
		if err := rows.Scan(&day, &d.Revenue, &d.Orders); err != nil {
			rows.Close()
			return nil, err
		}
		d.Date = day.Format("2006-01-02")
		s.SalesByDay = append(s.SalesByDay, d)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	prodRows, err := r.db.Pool.Query(ctx, `
		SELECT item->>'product_id', item->>'name', SUM((item->>'quantity')::int), SUM((item->>'total')::bigint)
		FROM orders, jsonb_array_elements(items) AS item
		WHERE tenant_id = $1
		GROUP BY item->>'product_id', item->>'name'
		ORDER BY 3 DESC
		LIMIT 5
	`, tenantID)
	if err != nil {
		return nil, err
	}
	for prodRows.Next() {
		var p domain.ProductSales
		if err := prodRows.Scan(&p.ProductID, &p.Name, &p.Units, &p.Revenue); err != nil {
			prodRows.Close()
			return nil, err
		}
		s.TopProducts = append(s.TopProducts, p)
	}
	prodRows.Close()
	if err := prodRows.Err(); err != nil {
		return nil, err
	}

	return s, nil
}

func (r *Repository) AdminSummary(ctx context.Context, tenantID string) (*domain.AdminSummary, error) {
	a := &domain.AdminSummary{
		UsersByRole:     make(map[string]int),
		ProductsByState: make(map[string]int),
		OrdersByStatus:  make(map[string]int),
	}

	rows, err := r.db.Pool.Query(ctx, `SELECT role, COUNT(*) FROM users WHERE tenant_id = $1 GROUP BY role`, tenantID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var role string
		var count int
		if err := rows.Scan(&role, &count); err != nil {
			rows.Close()
			return nil, err
		}
		a.UsersByRole[role] = count
	}
	rows.Close()

	rows, err = r.db.Pool.Query(ctx, `SELECT state, COUNT(*) FROM products WHERE tenant_id = $1 GROUP BY state`, tenantID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var state string
		var count int
		if err := rows.Scan(&state, &count); err != nil {
			rows.Close()
			return nil, err
		}
		a.ProductsByState[state] = count
	}
	rows.Close()

	rows, err = r.db.Pool.Query(ctx, `SELECT status, COUNT(*) FROM orders WHERE tenant_id = $1 GROUP BY status`, tenantID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			rows.Close()
			return nil, err
		}
		a.OrdersByStatus[status] = count
	}
	rows.Close()

	if err := r.db.Pool.QueryRow(ctx, `SELECT COALESCE(SUM(total), 0) FROM orders WHERE tenant_id = $1`, tenantID).Scan(&a.TotalRevenue); err != nil {
		return nil, err
	}
	if err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM product_reviews WHERE tenant_id = $1`, tenantID).Scan(&a.TotalReviews); err != nil {
		return nil, err
	}

	return a, nil
}

func (r *Repository) Health(ctx context.Context, tenantID string) (*domain.HealthReport, error) {
	h := &domain.HealthReport{Version: r.version, UptimeSeconds: int64(time.Since(r.startedAt).Seconds())}

	var one int
	if err := r.db.Pool.QueryRow(ctx, `SELECT 1`).Scan(&one); err != nil {
		h.DatabaseOK = false
	} else {
		h.DatabaseOK = true
	}

	_ = r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM schema_migrations`).Scan(&h.MigrationsApplied)
	_ = r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE tenant_id = $1`, tenantID).Scan(&h.TotalUsers)
	_ = r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE tenant_id = $1`, tenantID).Scan(&h.TotalProducts)
	_ = r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE tenant_id = $1`, tenantID).Scan(&h.TotalOrders)

	return h, nil
}

// RecordSnapshot measures current DB latency (via a lightweight ping) and
// platform-wide counts, then persists one row to health_snapshots. Called
// periodically by a background job in main.go — this is how the health
// history charts get real, structured data over time instead of mocks.
func (r *Repository) RecordSnapshot(ctx context.Context) error {
	start := time.Now()
	var one int
	pingErr := r.db.Pool.QueryRow(ctx, `SELECT 1`).Scan(&one)
	latencyMs := float64(time.Since(start).Microseconds()) / 1000.0

	var users, products, orders int
	var revenue int64
	_ = r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&users)
	_ = r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM products`).Scan(&products)
	_ = r.db.Pool.QueryRow(ctx, `SELECT COUNT(*), COALESCE(SUM(total), 0) FROM orders`).Scan(&orders, &revenue)

	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO health_snapshots (id, recorded_at, db_ok, db_latency_ms, total_users, total_products, total_orders, total_revenue)
		VALUES (gen_random_uuid(), NOW(), $1, $2, $3, $4, $5, $6)
	`, pingErr == nil, latencyMs, users, products, orders, revenue)
	return err
}

// HealthHistory returns sampled health snapshots recorded since the given
// time, ordered chronologically — powers the historical trend charts with
// period drill-down (24h/7d/30d) on the Global Admin panel.
func (r *Repository) HealthHistory(ctx context.Context, since time.Time) ([]domain.HealthPoint, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT recorded_at, db_ok, db_latency_ms, total_users, total_products, total_orders, total_revenue
		FROM health_snapshots
		WHERE recorded_at >= $1
		ORDER BY recorded_at ASC
	`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []domain.HealthPoint
	for rows.Next() {
		var p domain.HealthPoint
		var recordedAt time.Time
		if err := rows.Scan(&recordedAt, &p.DBOk, &p.DBLatencyMs, &p.TotalUsers, &p.TotalProducts, &p.TotalOrders, &p.TotalRevenue); err != nil {
			return nil, err
		}
		p.RecordedAt = recordedAt.Format(time.RFC3339)
		points = append(points, p)
	}
	return points, rows.Err()
}
