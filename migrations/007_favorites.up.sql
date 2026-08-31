-- Migration 007: Favorites (lista pessoal de itens favoritados)
-- Created: 2026-08-31

CREATE TABLE favorites (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    user_id     UUID NOT NULL REFERENCES users(id),
    product_id  UUID NOT NULL REFERENCES products(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, user_id, product_id)
);

CREATE INDEX idx_favorites_user ON favorites(tenant_id, user_id, created_at DESC);
