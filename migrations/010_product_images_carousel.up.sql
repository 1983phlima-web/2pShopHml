-- Migration 010: Carrossel de imagens do produto (até 3)
-- Created: 2026-08-31
-- Deriva 3 recortes reais (mesma foto original, enquadramentos
-- diferentes via parâmetros da própria Unsplash) para popular
-- attributes.images — usado pelo carrossel no modal de produto.

UPDATE products
SET attributes = attributes || jsonb_build_object(
  'images', jsonb_build_array(
    attributes->>'image',
    split_part(attributes->>'image', '?', 1) || '?auto=format&fit=crop&w=900&h=1100&crop=top&q=80',
    split_part(attributes->>'image', '?', 1) || '?auto=format&fit=crop&w=900&h=600&crop=bottom&q=80'
  )
)
WHERE attributes ? 'image' AND COALESCE(attributes->>'image', '') <> '';
