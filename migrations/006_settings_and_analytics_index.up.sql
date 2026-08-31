-- Migration 006: System settings (layout/palette parametrization) + analytics index
-- Created: 2026-08-30

CREATE TABLE system_settings (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    key         VARCHAR(100) NOT NULL,
    value       JSONB NOT NULL DEFAULT '{}',
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, key)
);

INSERT INTO system_settings (tenant_id, key, value)
VALUES ('00000000-0000-0000-0000-000000000001', 'theme', '{"palette": "indigo"}'::jsonb)
ON CONFLICT (tenant_id, key) DO NOTHING;

-- Speeds up per-product aggregation (units sold, top products) over the
-- orders.items JSONB array used by the seller/admin analytics dashboards.
CREATE INDEX idx_orders_items_gin ON orders USING GIN (items);
