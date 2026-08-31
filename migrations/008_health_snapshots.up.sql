-- Migration 008: Health snapshots (histórico real para os gráficos de saúde)
-- Created: 2026-08-31
-- Populated periodically by a background goroutine in the API service —
-- this is genuine sampled telemetry, not seed/mock data.

CREATE TABLE health_snapshots (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recorded_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    db_ok          BOOLEAN NOT NULL,
    db_latency_ms  DOUBLE PRECISION NOT NULL,
    total_users    INT NOT NULL,
    total_products INT NOT NULL,
    total_orders   INT NOT NULL,
    total_revenue  BIGINT NOT NULL
);

CREATE INDEX idx_health_snapshots_recorded_at ON health_snapshots(recorded_at DESC);
