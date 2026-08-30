'use client';

import { useEffect, useState, useCallback } from 'react';
import { ProductCard, Product } from '@/components/ProductCard';
import { api } from '@/lib/api';

const PAGE_SIZE = 24;

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[] | null>(null);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);

  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE));

  const load = useCallback(async (p: number) => {
    setProducts(null);
    try {
      const offset = (p - 1) * PAGE_SIZE;
      const res = await api(`/products?limit=${PAGE_SIZE}&offset=${offset}`);
      if (!res.ok) {
        setProducts([]);
        return;
      }
      const data = await res.json();
      setProducts(data.data || []);
      setTotal(data.total || 0);
    } catch {
      setProducts([]);
    }
  }, []);

  useEffect(() => {
    load(page);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }, [page, load]);

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Produtos</h1>
      {products === null ? (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-6">
          {Array.from({ length: 8 }).map((_, i) => (
            <div key={i} className="bg-white rounded-xl border border-gray-100 h-72 animate-pulse" />
          ))}
        </div>
      ) : products.length === 0 ? (
        <div className="text-center py-20 text-gray-500">Nenhum produto disponível no momento.</div>
      ) : (
        <>
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-6">
            {products.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
          <nav className="mt-10 flex flex-wrap items-center justify-between gap-3" aria-label="Paginação de produtos">
            <span className="text-xs font-semibold text-gray-500">
              Página {page} de {totalPages} · {total} produtos
            </span>
            <div className="flex items-center gap-2">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page === 1}
                className="rounded-lg border border-gray-300 bg-white px-4 py-2 text-xs font-bold text-gray-900 disabled:cursor-not-allowed disabled:opacity-40 hover:border-indigo-300 transition"
              >
                Anterior
              </button>
              <button
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page === totalPages}
                className="rounded-lg bg-indigo-600 px-4 py-2 text-xs font-bold text-white disabled:cursor-not-allowed disabled:opacity-40 hover:bg-indigo-700 transition"
              >
                Próxima
              </button>
            </div>
          </nav>
        </>
      )}
    </div>
  );
}
