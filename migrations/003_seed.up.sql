-- Migration 003: Seed data for HML
-- Created: 2026-08-29
-- Provides one demo tenant, one test user per RBAC role, a small sample
-- catalog with stock, one completed order (compra) and one product
-- review (comentário) — usage examples requested for the HML delivery.

-- Demo tenant (matches the frontend's default X-Tenant-ID: tenant_01)
INSERT INTO tenants (id, name, slug, status, plan_id, limits)
VALUES ('00000000-0000-0000-0000-000000000001', '2pShop Demo', 'tenant_01', 'ACTIVE', 'demo', '{}')
ON CONFLICT (slug) DO NOTHING;

-- Test users — one per RBAC role. Passwords (bcrypt-hashed below) are
-- shared with the requester out of band; see delivery notes.
INSERT INTO users (id, tenant_id, email, name, role, password_hash, active) VALUES
  ('00000000-0000-0000-0000-000000000101', '00000000-0000-0000-0000-000000000001', 'vendedor@2pshop.com.br', 'Vendedor Demo', 'SELLER', '$2b$10$VzjhlFyyB1TMstR2BT/AG.ZtW02IUFAxfBTP2rB6wMf9EuR2hdBHi', true),
  ('00000000-0000-0000-0000-000000000102', '00000000-0000-0000-0000-000000000001', 'cliente@2pshop.com.br', 'Cliente Demo', 'BUYER', '$2b$10$0.ISbR1BwMhGa2yFBc0hWOCHyVtZf/.4nBpJN3BMyxjGPPmg.5MkC', true),
  ('00000000-0000-0000-0000-000000000103', '00000000-0000-0000-0000-000000000001', 'admin.sistema@2pshop.com.br', 'Administrador do Sistema', 'SYSTEM_ADMIN', '$2b$10$e2RF37E7yOJssSibVF01d.QKVfxZMfXVyMSGkGMj3ylUdWtsnaClC', true),
  ('00000000-0000-0000-0000-000000000104', '00000000-0000-0000-0000-000000000001', 'admin.global@2pshop.com.br', 'Administrador Global', 'GLOBAL_ADMIN', '$2b$10$wpI9CFFG9BCbIfALIslqX.UTezkGusn9efLkJcZJ8TiuPDcsnO1zi', true)
ON CONFLICT (tenant_id, email) DO NOTHING;

-- Sample catalog
INSERT INTO products (id, tenant_id, name, slug, description, sku, price, state) VALUES
  ('00000000-0000-0000-0000-000000000201', '00000000-0000-0000-0000-000000000001', 'Fone de Ouvido Bluetooth 2P Sound', 'fone-bluetooth-2p-sound', 'Fone de ouvido sem fio com cancelamento de ruído e 30h de bateria.', 'FONE-2P-001', 24990, 'ACTIVE'),
  ('00000000-0000-0000-0000-000000000202', '00000000-0000-0000-0000-000000000001', 'Mochila Executiva 2P Urban', 'mochila-executiva-2p-urban', 'Mochila resistente à água com compartimento acolchoado para notebook.', 'MOCH-2P-002', 18990, 'ACTIVE'),
  ('00000000-0000-0000-0000-000000000203', '00000000-0000-0000-0000-000000000001', 'Garrafa Térmica Inox 2P Hydro', 'garrafa-termica-2p-hydro', 'Garrafa térmica de aço inox, mantém a temperatura por até 12 horas.', 'GARR-2P-003', 8990, 'ACTIVE'),
  ('00000000-0000-0000-0000-000000000204', '00000000-0000-0000-0000-000000000001', 'Teclado Mecânico 2P TypeMaster', 'teclado-mecanico-2p-typemaster', 'Teclado mecânico com switches táteis e iluminação RGB.', 'TECL-2P-004', 34990, 'ACTIVE')
ON CONFLICT (tenant_id, slug) DO NOTHING;

INSERT INTO inventory (id, tenant_id, product_id, quantity, reserved) VALUES
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000201', 50, 0),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000202', 50, 0),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000203', 50, 0),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000204', 50, 0)
ON CONFLICT (tenant_id, product_id) WHERE variant_id IS NULL DO NOTHING;

-- Sample order (compra) placed by the demo customer, already delivered.
INSERT INTO orders (id, tenant_id, customer_id, items, total, status) VALUES (
  '00000000-0000-0000-0000-000000000301',
  '00000000-0000-0000-0000-000000000001',
  '00000000-0000-0000-0000-000000000102',
  '[
    {"product_id":"00000000-0000-0000-0000-000000000201","name":"Fone de Ouvido Bluetooth 2P Sound","quantity":2,"unit_price":24990,"total":49980},
    {"product_id":"00000000-0000-0000-0000-000000000203","name":"Garrafa Térmica Inox 2P Hydro","quantity":1,"unit_price":8990,"total":8990}
  ]'::jsonb,
  58970,
  'DELIVERED'
) ON CONFLICT (id) DO NOTHING;

-- Sample review (comentário) on the headphones product.
INSERT INTO product_reviews (id, tenant_id, product_id, user_id, user_name, rating, comment) VALUES (
  '00000000-0000-0000-0000-000000000401',
  '00000000-0000-0000-0000-000000000001',
  '00000000-0000-0000-0000-000000000201',
  '00000000-0000-0000-0000-000000000102',
  'Cliente Demo',
  5,
  'Excelente qualidade de som e ótimo custo-benefício! Chegou rápido e bem embalado.'
) ON CONFLICT (id) DO NOTHING;
