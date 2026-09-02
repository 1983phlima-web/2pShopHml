-- Migration 009: Perfil do cliente (avatar/contato), Perguntas, Trocas, Fidelidade (XP/Coins/Badges)
-- Created: 2026-08-31

ALTER TABLE users ADD COLUMN avatar VARCHAR(500);
ALTER TABLE users ADD COLUMN phone VARCHAR(30);

-- Perguntas feitas nos produtos anunciados (parte de "Meus Comentários" -> Perguntas)
CREATE TABLE product_questions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    product_id  UUID NOT NULL REFERENCES products(id),
    user_id     UUID NOT NULL REFERENCES users(id),
    user_name   VARCHAR(255) NOT NULL,
    question    TEXT NOT NULL,
    answer      TEXT,
    answered_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_product_questions_product ON product_questions(tenant_id, product_id, created_at DESC);
CREATE INDEX idx_product_questions_user ON product_questions(tenant_id, user_id, created_at DESC);

-- Solicitações de troca/devolução pós-entrega ("Minhas Trocas")
CREATE TABLE exchange_requests (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    order_id    UUID NOT NULL REFERENCES orders(id),
    product_id  UUID NOT NULL REFERENCES products(id),
    user_id     UUID NOT NULL REFERENCES users(id),
    reason      TEXT NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'REQUESTED',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_exchange_requests_user ON exchange_requests(tenant_id, user_id, created_at DESC);

-- Fidelidade: ledger de badges/coins concedidos por período (15 dias / mês),
-- idempotente via UNIQUE — cada período só concede o mesmo badge uma vez.
CREATE TABLE loyalty_awards (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    user_id     UUID NOT NULL REFERENCES users(id),
    period_type VARCHAR(10) NOT NULL,   -- '15day' | 'month'
    period_key  VARCHAR(20) NOT NULL,   -- ex: '2026-08-P1', '2026-08'
    badge       VARCHAR(20) NOT NULL,   -- '500' | '1000' | 'vip'
    coins       INT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, user_id, period_type, period_key, badge)
);
CREATE INDEX idx_loyalty_awards_user ON loyalty_awards(tenant_id, user_id);

-- Avatar padrão aleatório para os 4 usuários de teste já existentes.
UPDATE users SET avatar = 'preset:2' WHERE email = 'vendedor@2pshop.com.br';
UPDATE users SET avatar = 'preset:7' WHERE email = 'cliente@2pshop.com.br';
UPDATE users SET avatar = 'preset:11' WHERE email = 'admin.sistema@2pshop.com.br';
UPDATE users SET avatar = 'preset:4' WHERE email = 'admin.global@2pshop.com.br';

-- Idioma padrão do site (parametrizável pelos Admins).
INSERT INTO system_settings (tenant_id, key, value)
VALUES ('00000000-0000-0000-0000-000000000001', 'language', '{"code": "pt"}'::jsonb)
ON CONFLICT (tenant_id, key) DO NOTHING;
