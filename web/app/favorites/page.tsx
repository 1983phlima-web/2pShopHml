'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/components/AuthContext';
import { ProductCard, Product } from '@/components/ProductCard';
import { api } from '@/lib/api';

export default function FavoritesPage() {
  const { user } = useAuth();
  const [products, setProducts] = useState<Product[] | null>(null);

  useEffect(() => {
    if (!user) return;
    api('/favorites')
      .then(async (res) => {
        if (res.ok) setProducts((await res.json()) || []);
        else setProducts([]);
      })
      .catch(() => setProducts([]));
  }, [user]);

  if (!user) {
    return (
      <div className="max-w-md mx-auto text-center py-20">
        <p className="text-gray-600 mb-4">Entre na sua conta para ver seus favoritos.</p>
        <Link href="/login" className="text-indigo-600 hover:underline">Entrar</Link>
      </div>
    );
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Favoritos</h1>
      {products === null ? (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-6">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="bg-white rounded-xl border border-gray-100 h-72 animate-pulse" />
          ))}
        </div>
      ) : products.length === 0 ? (
        <div className="text-center py-20 text-gray-500">
          Você ainda não favoritou nenhum produto. Toque no coraçãozinho de um produto para adicioná-lo aqui.
        </div>
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-6">
          {products.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      )}
    </div>
  );
}
