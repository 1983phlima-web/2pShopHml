-- Migration 001: Foundation
-- Created: 2026-08-28

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tenants
CREATE TABLE tenants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL UNIQUE,
    status      VARCHAR(20) NOT NULL DEFAULT 'TRIAL',
    plan_id     VARCHAR(50) NOT NULL,
    limits      JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tenants_status ON tenants(status);
CREATE INDEX idx_tenants_slug ON tenants(slug);

-- Users
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    email       VARCHAR(255) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    role        VARCHAR(20) NOT NULL DEFAULT 'BUYER',
    password_hash VARCHAR(255) NOT NULL,
    active      BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, email)
);

CREATE INDEX idx_users_tenant_email ON users(tenant_id, email);

-- Products
CREATE TABLE products (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    name        VARCHAR(500) NOT NULL,
    slug        VARCHAR(500) NOT NULL,
    description TEXT,
    sku         VARCHAR(100) NOT NULL,
    price       BIGINT NOT NULL DEFAULT 0,
    state       VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    category_id UUID,
    brand_id    UUID,
    attributes  JSONB DEFAULT '{}',
    seo         JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, slug),
    UNIQUE(tenant_id, sku)
);

CREATE INDEX idx_products_tenant_status ON products(tenant_id, state);
CREATE INDEX idx_products_tenant_slug ON products(tenant_id, slug);

-- Categories
CREATE TABLE categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL,
    parent_id   UUID REFERENCES categories(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, slug)
);

-- Inventory
CREATE TABLE inventory (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    product_id  UUID NOT NULL REFERENCES products(id),
    variant_id  UUID,
    quantity    INT NOT NULL DEFAULT 0,
    reserved    INT NOT NULL DEFAULT 0,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, product_id, variant_id)
);

CREATE INDEX idx_inventory_tenant_product ON inventory(tenant_id, product_id);

-- Reservations
CREATE TABLE reservations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    product_id  UUID NOT NULL REFERENCES products(id),
    variant_id  UUID,
    quantity    INT NOT NULL,
    order_id    VARCHAR(100) NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reservations_expires ON reservations(expires_at);
CREATE INDEX idx_reservations_tenant_order ON reservations(tenant_id, order_id);

-- Orders
CREATE TABLE orders (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    customer_id UUID NOT NULL REFERENCES users(id),
    items       JSONB NOT NULL DEFAULT '[]',
    total       BIGINT NOT NULL DEFAULT 0,
    status      VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_tenant_created ON orders(tenant_id, created_at DESC);
CREATE INDEX idx_orders_tenant_status ON orders(tenant_id, status);

-- Outbox (event sourcing pattern)
CREATE TABLE outbox_events (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id   VARCHAR(100) NOT NULL,
    event_type  VARCHAR(100) NOT NULL,
    payload     JSONB NOT NULL,
    headers     JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    retry_count INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_outbox_unpublished ON outbox_events(published_at) WHERE published_at IS NULL;
CREATE INDEX idx_outbox_aggregate ON outbox_events(aggregate_type, aggregate_id);

-- Idempotency keys
CREATE TABLE idempotency_keys (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    key         VARCHAR(255) NOT NULL,
    operation   VARCHAR(100) NOT NULL,
    request_hash VARCHAR(64) NOT NULL,
    status      VARCHAR(20) NOT NULL,
    response_code INT,
    response_body JSONB,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, key)
);

CREATE INDEX idx_idempotency_expires ON idempotency_keys(expires_at);

-- Payments
CREATE TABLE payments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    order_id        UUID NOT NULL REFERENCES orders(id),
    provider        VARCHAR(50) NOT NULL,
    provider_tx_id  VARCHAR(255),
    amount          BIGINT NOT NULL,
    currency        VARCHAR(3) NOT NULL DEFAULT 'BRL',
    status          VARCHAR(20) NOT NULL,
    method_type     VARCHAR(20),
    method_last4    VARCHAR(4),
    idempotency_key VARCHAR(255),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_tenant_order ON payments(tenant_id, order_id);

-- Audit log
CREATE TABLE audit_log (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL,
    table_name  VARCHAR(100) NOT NULL,
    record_id   VARCHAR(100) NOT NULL,
    action      VARCHAR(20) NOT NULL,
    old_values  JSONB,
    new_values  JSONB,
    performed_by UUID,
    performed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_tenant ON audit_log(tenant_id, performed_at DESC);
