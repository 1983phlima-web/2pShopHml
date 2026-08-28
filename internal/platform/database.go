package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/2pshop/2pshop/pkg/telemetry"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DB struct {
	Pool    *pgxpool.Pool
	metrics *telemetry.DBMetrics
}

func NewDB(ctx context.Context, databaseURL string, metrics *telemetry.DBMetrics) (*DB, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	config.HealthCheckPeriod = time.Minute * 5

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &DB{Pool: pool, metrics: metrics}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) WithTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantContextKey{}, tenantID)
}

func (db *DB) TenantFromContext(ctx context.Context) string {
	if v := ctx.Value(tenantContextKey{}); v != nil {
		return v.(string)
	}
	return ""
}

type tenantContextKey struct{}

func (db *DB) QueryRow(ctx context.Context, operation string, sql string, args ...any) error {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.operation", operation),
		)
	}

	start := time.Now()
	_, err := db.Pool.Exec(ctx, sql, args...)
	telemetry.RecordQuery(ctx, db.metrics, operation, time.Since(start), err)
	return err
}
