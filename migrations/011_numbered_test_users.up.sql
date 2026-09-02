-- Migration 011: Usuários de teste adicionais (numeração sequencial)
-- Created: 2026-09-02
-- Mesma senha por papel dos usuários originais (já com hash validado).

INSERT INTO users (id, tenant_id, email, name, role, password_hash, avatar, active) VALUES
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'cliente1@2pshop.com.br', 'Cliente 1', 'BUYER', '$2b$10$vpluvd8c92uRiJtK6Kvf1.ErWy0ouC.Q7pVwJp0nylKaALBuTXEg6', 'preset:3', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'cliente2@2pshop.com.br', 'Cliente 2', 'BUYER', '$2b$10$vpluvd8c92uRiJtK6Kvf1.ErWy0ouC.Q7pVwJp0nylKaALBuTXEg6', 'preset:6', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'cliente3@2pshop.com.br', 'Cliente 3', 'BUYER', '$2b$10$vpluvd8c92uRiJtK6Kvf1.ErWy0ouC.Q7pVwJp0nylKaALBuTXEg6', 'preset:9', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'cliente4@2pshop.com.br', 'Cliente 4', 'BUYER', '$2b$10$vpluvd8c92uRiJtK6Kvf1.ErWy0ouC.Q7pVwJp0nylKaALBuTXEg6', 'preset:12', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'cliente5@2pshop.com.br', 'Cliente 5', 'BUYER', '$2b$10$vpluvd8c92uRiJtK6Kvf1.ErWy0ouC.Q7pVwJp0nylKaALBuTXEg6', 'preset:1', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'cliente6@2pshop.com.br', 'Cliente 6', 'BUYER', '$2b$10$vpluvd8c92uRiJtK6Kvf1.ErWy0ouC.Q7pVwJp0nylKaALBuTXEg6', 'preset:5', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'cliente7@2pshop.com.br', 'Cliente 7', 'BUYER', '$2b$10$vpluvd8c92uRiJtK6Kvf1.ErWy0ouC.Q7pVwJp0nylKaALBuTXEg6', 'preset:8', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'cliente8@2pshop.com.br', 'Cliente 8', 'BUYER', '$2b$10$vpluvd8c92uRiJtK6Kvf1.ErWy0ouC.Q7pVwJp0nylKaALBuTXEg6', 'preset:2', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'cliente9@2pshop.com.br', 'Cliente 9', 'BUYER', '$2b$10$vpluvd8c92uRiJtK6Kvf1.ErWy0ouC.Q7pVwJp0nylKaALBuTXEg6', 'preset:10', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'cliente10@2pshop.com.br', 'Cliente 10', 'BUYER', '$2b$10$vpluvd8c92uRiJtK6Kvf1.ErWy0ouC.Q7pVwJp0nylKaALBuTXEg6', 'preset:7', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'vendedor1@2pshop.com.br', 'Vendedor 1', 'SELLER', '$2b$10$81jmdRgEQdSt8JNftYbkLOL0eFMANa/ATBW8HFvsuT/u4C8WwyMDm', 'preset:4', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'vendedor2@2pshop.com.br', 'Vendedor 2', 'SELLER', '$2b$10$81jmdRgEQdSt8JNftYbkLOL0eFMANa/ATBW8HFvsuT/u4C8WwyMDm', 'preset:11', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'vendedor3@2pshop.com.br', 'Vendedor 3', 'SELLER', '$2b$10$81jmdRgEQdSt8JNftYbkLOL0eFMANa/ATBW8HFvsuT/u4C8WwyMDm', 'preset:3', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'vendedor4@2pshop.com.br', 'Vendedor 4', 'SELLER', '$2b$10$81jmdRgEQdSt8JNftYbkLOL0eFMANa/ATBW8HFvsuT/u4C8WwyMDm', 'preset:6', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'vendedor5@2pshop.com.br', 'Vendedor 5', 'SELLER', '$2b$10$81jmdRgEQdSt8JNftYbkLOL0eFMANa/ATBW8HFvsuT/u4C8WwyMDm', 'preset:9', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'vendedor6@2pshop.com.br', 'Vendedor 6', 'SELLER', '$2b$10$81jmdRgEQdSt8JNftYbkLOL0eFMANa/ATBW8HFvsuT/u4C8WwyMDm', 'preset:12', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'vendedor7@2pshop.com.br', 'Vendedor 7', 'SELLER', '$2b$10$81jmdRgEQdSt8JNftYbkLOL0eFMANa/ATBW8HFvsuT/u4C8WwyMDm', 'preset:1', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'admin.sistema1@2pshop.com.br', 'Administrador do Sistema 1', 'SYSTEM_ADMIN', '$2b$10$28rmZcDYo8BbUm/S/aFXXObxD6plKZKwFuyeoVlIkHYfyeuLU5GSy', 'preset:5', true)
ON CONFLICT (tenant_id, email) DO NOTHING;
