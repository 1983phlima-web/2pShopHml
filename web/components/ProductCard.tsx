'use client';

import Image from 'next/image';
import { useRouter } from 'next/navigation';
import { useCart } from './CartContext';
import { useFavorites } from './FavoritesContext';
import { useAuth } from './AuthContext';
import { useProductModal } from './ProductModalContext';

export interface Product {
  id: string;
  name: string;
  slug: string;
  description: string;
  price: number;
  state: string;
  attributes?: {
    image?: string;
    badge?: string;
    brand?: string;
    compare_at?: number;
    sold?: number;
  };
}

function HeartIcon({ filled }: { filled: boolean }) {
  return (
    <svg
      viewBox="0 0 24 24"
      className="h-6 w-6"
      fill={filled ? '#e11d48' : 'none'}
      stroke={filled ? '#e11d48' : '#ffffff'}
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M12 21s-6.716-4.35-9.428-8.06C.79 10.31 1.2 6.6 4.2 4.9c2.3-1.3 5-0.7 6.5 1.3l1.3 1.7 1.3-1.7c1.5-2 4.2-2.6 6.5-1.3 3 1.7 3.41 5.41 1.63 8.04C18.716 16.65 12 21 12 21z"
      />
    </svg>
  );
}

export function ProductCard({ product }: { product: Product }) {
  const { add } = useCart();
  const { isFavorite, toggle } = useFavorites();
  const { user } = useAuth();
  const router = useRouter();
  const { open } = useProductModal();
  const image = product.attributes?.image;
  const badge = product.attributes?.badge;
  const compareAt = product.attributes?.compare_at;
  const favorited = isFavorite(product.id);

  function handleFavoriteClick(e: React.MouseEvent) {
    e.preventDefault();
    e.stopPropagation();
    if (!user) {
      router.push('/login');
      return;
    }
    toggle(product.id);
  }

  return (
    <div className="bg-white rounded-xl border border-gray-100 overflow-hidden hover:shadow-md transition group">
      <div className="block relative cursor-pointer" onClick={() => open(product.id)}>
        <div className="h-40 bg-gray-100 overflow-hidden relative">
          {image ? (
            <Image
              src={image}
              alt={product.name}
              fill
              sizes="(max-width: 640px) 50vw, (max-width: 1024px) 25vw, 20vw"
              className="object-cover group-hover:scale-105 transition duration-300"
            />
          ) : (
            <div className="h-full w-full flex items-center justify-center text-gray-400 text-4xl">📦</div>
          )}
        </div>
        {badge && (
          <span className="absolute top-2 left-2 bg-indigo-600 text-white text-[10px] font-bold px-2 py-1 rounded-full uppercase tracking-wide z-10">
            {badge}
          </span>
        )}
        <button
          onClick={handleFavoriteClick}
          aria-label={favorited ? 'Remover dos favoritos' : 'Adicionar aos favoritos'}
          aria-pressed={favorited}
          className="absolute top-2 right-2 h-8 w-8 flex items-center justify-center hover:scale-110 transition z-10 drop-shadow"
        >
          <HeartIcon filled={favorited} />
        </button>
      </div>
      <div className="p-4">
        <h3
          onClick={() => open(product.id)}
          className="font-semibold text-gray-900 truncate hover:text-indigo-600 cursor-pointer"
        >
          {product.name}
        </h3>
        <p className="text-sm text-gray-500 line-clamp-2 mt-1">{product.description}</p>
        <div className="mt-4 flex items-center justify-between gap-2">
          <div className="min-w-0">
            {compareAt && compareAt > product.price && (
              <span className="text-xs text-gray-400 line-through block whitespace-nowrap">R$ {(compareAt / 100).toFixed(2)}</span>
            )}
            <span className="text-lg font-bold text-indigo-600 whitespace-nowrap">
              R$ {(product.price / 100).toFixed(2)}
            </span>
          </div>
          <button
            onClick={() => add({ id: product.id, name: product.name, price: product.price, quantity: 1 })}
            aria-label="Adicionar ao carrinho"
            className="shrink-0 h-9 w-9 flex items-center justify-center bg-indigo-50 text-indigo-700 rounded-md hover:bg-indigo-100 transition"
          >
            <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
}
