-- Migration 002: Product reviews (comments) + inventory upsert fix
-- Created: 2026-08-29

-- Product reviews / comments — shown on the product detail page.
CREATE TABLE product_reviews (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    product_id  UUID NOT NULL REFERENCES products(id),
    user_id     UUID NOT NULL REFERENCES users(id),
    user_name   VARCHAR(255) NOT NULL,
    rating      SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment     TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_product_reviews_product ON product_reviews(tenant_id, product_id, created_at DESC);

-- The original UNIQUE(tenant_id, product_id, variant_id) constraint on
-- inventory does not behave as an upsert target for the common
-- no-variant case, because Postgres treats each NULL variant_id as
-- distinct for uniqueness purposes. Add a partial unique index that
-- covers that case explicitly, used by the inventory adapter's
-- ON CONFLICT upsert.
CREATE UNIQUE INDEX idx_inventory_unique_no_variant ON inventory(tenant_id, product_id) WHERE variant_id IS NULL;
