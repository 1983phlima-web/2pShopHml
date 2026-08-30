-- Migration 004: Simplify seed user passwords
-- Created: 2026-08-30
-- Migration 003 already ran, so its INSERT ... ON CONFLICT DO NOTHING would
-- silently skip on a password change. Update the existing rows explicitly.

UPDATE users SET password_hash = '$2b$10$81jmdRgEQdSt8JNftYbkLOL0eFMANa/ATBW8HFvsuT/u4C8WwyMDm'
WHERE tenant_id = '00000000-0000-0000-0000-000000000001' AND email = 'vendedor@2pshop.com.br';

UPDATE users SET password_hash = '$2b$10$vpluvd8c92uRiJtK6Kvf1.ErWy0ouC.Q7pVwJp0nylKaALBuTXEg6'
WHERE tenant_id = '00000000-0000-0000-0000-000000000001' AND email = 'cliente@2pshop.com.br';

UPDATE users SET password_hash = '$2b$10$28rmZcDYo8BbUm/S/aFXXObxD6plKZKwFuyeoVlIkHYfyeuLU5GSy'
WHERE tenant_id = '00000000-0000-0000-0000-000000000001' AND email = 'admin.sistema@2pshop.com.br';

UPDATE users SET password_hash = '$2b$10$nLVOQ8Y0mt7qU9jTsECMhOuWI9fvGcdqx4SIWvpPOue7kQUQBjpya'
WHERE tenant_id = '00000000-0000-0000-0000-000000000001' AND email = 'admin.global@2pshop.com.br';
