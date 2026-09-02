'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useProductModal } from './ProductModalContext';
import { ImageCarousel } from './ImageCarousel';
import { useCart } from './CartContext';
import { useFavorites } from './FavoritesContext';
import { useAuth } from './AuthContext';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';

interface ProductDetail {
  id: string;
  name: string;
  description: string;
  price: number;
  attributes?: {
    image?: string;
    images?: string[];
    badge?: string;
    brand?: string;
    compare_at?: number;
  };
}

export function ProductModal() {
  const { openProductId, close } = useProductModal();
  const { add } = useCart();
  const { isFavorite, toggle } = useFavorites();
  const { user } = useAuth();
  const router = useRouter();
  const [product, setProduct] = useState<ProductDetail | null>(null);

  useEffect(() => {
    if (!openProductId) {
      setProduct(null);
      return;
    }
    api(`/products/${openProductId}`).then(async (res) => {
      if (res.ok) setProduct(await res.json());
    });
  }, [openProductId]);

  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') close();
    }
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [close]);

  if (!openProductId) return null;

  const images = product?.attributes?.images?.length
    ? product.attributes.images
    : product?.attributes?.image
    ? [product.attributes.image]
    : [];

  return (
    <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center">
      <div className="absolute inset-0 bg-black/50" onClick={close} />
      <div className="relative bg-white w-full sm:max-w-2xl sm:rounded-2xl rounded-t-2xl max-h-[92vh] overflow-y-auto grid sm:grid-cols-2">
        <button
          onClick={close}
          aria-label="Fechar"
          className="absolute top-3 right-3 z-10 h-8 w-8 rounded-full bg-white/90 flex items-center justify-center text-gray-700 hover:bg-white shadow"
        >
          ✕
        </button>

        <div className="relative h-64 sm:h-full">
          {product ? <ImageCarousel images={images} alt={product.name} /> : <div className="h-full bg-gray-100 animate-pulse" />}
        </div>

        <div className="p-6">
          {!product ? (
            <div className="space-y-3">
              <div className="h-5 bg-gray-100 rounded animate-pulse" />
              <div className="h-4 bg-gray-100 rounded animate-pulse w-2/3" />
            </div>
          ) : (
            <>
              {product.attributes?.brand && (
                <p className="text-xs font-bold uppercase tracking-wide brand-text mb-1">{product.attributes.brand}</p>
              )}
              <h2 className="text-xl font-bold mb-2">{product.name}</h2>
              <p className="text-sm text-gray-500 line-clamp-4 mb-4">{product.description}</p>
              <div className="flex items-center gap-3 mb-5">
                {product.attributes?.compare_at && product.attributes.compare_at > product.price && (
                  <span className="text-sm text-gray-400 line-through">
                    R$ {(product.attributes.compare_at / 100).toFixed(2)}
                  </span>
                )}
                <span className="text-2xl font-bold brand-text">R$ {(product.price / 100).toFixed(2)}</span>
              </div>
              <div className="flex items-center gap-3">
                <button
                  onClick={() => add({ id: product.id, name: product.name, price: product.price, quantity: 1 })}
                  className="flex-1 py-2.5 brand-bg text-white rounded-md font-medium transition"
                >
                  Adicionar ao carrinho
                </button>
                <button
                  onClick={() => (user ? toggle(product.id) : router.push('/login'))}
                  aria-label="Favoritar"
                  className="h-11 w-11 shrink-0 rounded-md border border-gray-200 flex items-center justify-center hover:border-rose-300 transition"
                >
                  <svg
                    viewBox="0 0 24 24"
                    className="h-5 w-5"
                    fill={isFavorite(product.id) ? '#e11d48' : 'none'}
                    stroke={isFavorite(product.id) ? '#e11d48' : 'currentColor'}
                    strokeWidth={2}
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M12 21s-6.716-4.35-9.428-8.06C.79 10.31 1.2 6.6 4.2 4.9c2.3-1.3 5-0.7 6.5 1.3l1.3 1.7 1.3-1.7c1.5-2 4.2-2.6 6.5-1.3 3 1.7 3.41 5.41 1.63 8.04C18.716 16.65 12 21 12 21z"
                    />
                  </svg>
                </button>
              </div>
              <Link
                href={`/products/${product.id}`}
                onClick={close}
                className="block text-center text-sm brand-text hover:underline mt-4"
              >
                Ver comentários, perguntas e detalhes completos →
              </Link>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
