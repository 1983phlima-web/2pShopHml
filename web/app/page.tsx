'use client';

import { useEffect, useState, useCallback, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import { ProductCard, Product } from '@/components/ProductCard';
import { useAuth } from '@/components/AuthContext';
import { roleHomePath } from '@/lib/auth';
import { api } from '@/lib/api';

const CATEGORIES = [
  { name: 'Casa & decoração', icon: '🏠' },
  { name: 'Tecnologia', icon: '💻' },
  { name: 'Moda', icon: '👗' },
  { name: 'Beleza', icon: '💄' },
  { name: 'Esportes', icon: '⚽' },
  { name: 'Jardinagem', icon: '🌿' },
];

export default function HomePage() {
  const { user } = useAuth();
  const router = useRouter();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [featuredIndex, setFeaturedIndex] = useState(0);

  // Vendedor e Administração têm o Painel deles como página principal —
  // a home "/" (vitrine) é exclusiva do Cliente.
  useEffect(() => {
    if (user && user.role !== 'BUYER') {
      router.replace(roleHomePath(user.role));
    }
  }, [user, router]);

  useEffect(() => {
    api('/products?limit=48')
      .then(async (res) => {
        if (!res.ok) return;
        const data = await res.json();
        setProducts(data.data || []);
      })
      .finally(() => setLoading(false));
  }, []);

  const offerProduct = useMemo(
    () => products.find((p) => p.attributes?.badge === 'Oferta do dia') || products[0],
    [products]
  );

  const featured = useMemo(() => {
    const withBadge = products.filter((p) => p.attributes?.badge);
    return (withBadge.length >= 3 ? withBadge : products).slice(0, 6);
  }, [products]);

  const featuredProduct = featured[featuredIndex];

  const nextFeatured = useCallback(() => {
    setFeaturedIndex((i) => (i + 1) % Math.max(featured.length, 1));
  }, [featured.length]);
  const prevFeatured = useCallback(() => {
    setFeaturedIndex((i) => (i - 1 + featured.length) % Math.max(featured.length, 1));
  }, [featured.length]);

  if (loading || (user && user.role !== 'BUYER')) {
    return (
      <div className="space-y-8">
        <div className="h-64 bg-gray-100 rounded-2xl animate-pulse" />
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-6">
          {Array.from({ length: 8 }).map((_, i) => (
            <div key={i} className="bg-white rounded-xl border border-gray-100 h-72 animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-14">
      {/* Hero: Oferta do dia */}
      {offerProduct && (
        <section className="relative overflow-hidden rounded-2xl bg-gray-900 text-white grid md:grid-cols-[1.1fr_.9fr] min-h-[260px]">
          <div className="flex flex-col justify-center p-8 sm:p-10">
            <p className="text-xs font-extrabold uppercase tracking-[.2em] text-indigo-300">Oferta do dia</p>
            <h1 className="mt-3 text-2xl sm:text-3xl font-extrabold leading-tight">{offerProduct.name}</h1>
            <p className="mt-2 text-sm text-white/60 max-w-md line-clamp-2">{offerProduct.description}</p>
            <div className="mt-5 flex items-center gap-3">
              <span className="text-2xl font-extrabold">R$ {(offerProduct.price / 100).toFixed(2)}</span>
              {offerProduct.attributes?.compare_at && (
                <span className="text-sm text-white/40 line-through">
                  R$ {(offerProduct.attributes.compare_at / 100).toFixed(2)}
                </span>
              )}
            </div>
            <Link
              href={`/products/${offerProduct.id}`}
              className="mt-6 inline-flex w-fit items-center gap-2 rounded-lg brand-bg px-5 py-2.5 text-sm font-bold text-white transition"
            >
              Ver oferta →
            </Link>
          </div>
          {offerProduct.attributes?.image && (
            <div className="relative min-h-[220px]">
              <Image
                src={offerProduct.attributes.image}
                alt={offerProduct.name}
                fill
                sizes="(max-width: 768px) 100vw, 45vw"
                className="object-cover"
                priority
              />
            </div>
          )}
        </section>
      )}

      {/* Categories */}
      <section>
        <h2 className="text-lg font-bold mb-4">Categorias</h2>
        <div className="grid grid-cols-3 sm:grid-cols-6 gap-3">
          {CATEGORIES.map((c) => (
            <Link
              key={c.name}
              href="/products"
              className="flex flex-col items-center justify-center gap-1.5 rounded-xl border border-gray-100 bg-white py-4 hover:border-indigo-300 hover:shadow-sm transition"
            >
              <span className="text-2xl">{c.icon}</span>
              <span className="text-xs font-medium text-gray-700 text-center px-1">{c.name}</span>
            </Link>
          ))}
        </div>
      </section>

      {/* Destaque da semana - featured carousel */}
      {featuredProduct && (
        <section>
          <div className="flex items-end justify-between mb-4">
            <div>
              <p className="text-xs font-extrabold uppercase tracking-[.17em] brand-text">Destaque da semana</p>
              <h2 className="text-xl font-extrabold mt-1">{featuredProduct.name}</h2>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-xs text-gray-400">{featuredIndex + 1}/{featured.length}</span>
              <button
                onClick={prevFeatured}
                aria-label="Produto destacado anterior"
                className="h-8 w-8 grid place-items-center rounded-full border border-gray-200 hover:bg-gray-50"
              >
                ←
              </button>
              <button
                onClick={nextFeatured}
                aria-label="Próximo produto destacado"
                className="h-8 w-8 grid place-items-center rounded-full border border-gray-200 hover:bg-gray-50"
              >
                →
              </button>
            </div>
          </div>
          <div className="grid md:grid-cols-[1.1fr_.9fr] gap-0 rounded-2xl overflow-hidden bg-gray-900 text-white min-h-[220px]">
            <div className="p-8 flex flex-col justify-center">
              <p className="text-sm text-white/60 line-clamp-3 max-w-md">{featuredProduct.description}</p>
              <div className="mt-4 flex items-center gap-3">
                <span className="text-xl font-extrabold">R$ {(featuredProduct.price / 100).toFixed(2)}</span>
              </div>
              <Link
                href={`/products/${featuredProduct.id}`}
                className="mt-5 inline-flex w-fit items-center rounded-lg brand-bg px-4 py-2 text-xs font-bold text-white transition"
              >
                Ver produto →
              </Link>
            </div>
            {featuredProduct.attributes?.image && (
              <div className="relative min-h-[220px]">
                <Image
                  src={featuredProduct.attributes.image}
                  alt={featuredProduct.name}
                  fill
                  sizes="(max-width: 768px) 100vw, 45vw"
                  className="object-cover"
                />
              </div>
            )}
          </div>
        </section>
      )}

      {/* Vitrine preview */}
      <section id="vitrine">
        <div className="flex items-end justify-between mb-4">
          <h2 className="text-lg font-bold">Vitrine</h2>
          <Link href="/products" className="text-sm font-medium brand-text hover:underline">
            Ver todos os produtos →
          </Link>
        </div>
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-6">
          {products.slice(0, 8).map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      </section>
    </div>
  );
}
