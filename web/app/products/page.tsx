'use client';

import { useEffect, useState, useCallback } from 'react';
import { ProductCard, Product } from '@/components/ProductCard';
import { api } from '@/lib/api';

const PAGE_SIZE = 24;

interface Facets {
  categories: { slug: string; name: string }[];
  brands: string[];
  genders: string[];
}

interface Filters {
  q: string;
  category: string;
  brand: string;
  gender: string;
  minPrice: string;
  maxPrice: string;
}

const EMPTY_FILTERS: Filters = { q: '', category: '', brand: '', gender: '', minPrice: '', maxPrice: '' };

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[] | null>(null);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [facets, setFacets] = useState<Facets | null>(null);
  const [filters, setFilters] = useState<Filters>(EMPTY_FILTERS);
  const [draft, setDraft] = useState<Filters>(EMPTY_FILTERS);
  const [filtersOpen, setFiltersOpen] = useState(false);

  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE));
  const activeCount = Object.values(filters).filter(Boolean).length;

  useEffect(() => {
    api('/products/facets').then(async (res) => {
      if (res.ok) setFacets(await res.json());
    });
  }, []);

  const load = useCallback(async (p: number, f: Filters) => {
    setProducts(null);
    try {
      const params = new URLSearchParams({ limit: String(PAGE_SIZE), offset: String((p - 1) * PAGE_SIZE) });
      if (f.q) params.set('q', f.q);
      if (f.category) params.set('category', f.category);
      if (f.brand) params.set('brand', f.brand);
      if (f.gender) params.set('gender', f.gender);
      if (f.minPrice) params.set('min_price', String(Math.round(parseFloat(f.minPrice) * 100)));
      if (f.maxPrice) params.set('max_price', String(Math.round(parseFloat(f.maxPrice) * 100)));

      const res = await api(`/products?${params.toString()}`);
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
    load(page, filters);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }, [page, filters, load]);

  function applyFilters() {
    setFilters(draft);
    setPage(1);
    setFiltersOpen(false);
  }

  function clearFilters() {
    setDraft(EMPTY_FILTERS);
    setFilters(EMPTY_FILTERS);
    setPage(1);
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-2xl font-bold">Produtos</h1>
        <button
          onClick={() => setFiltersOpen((v) => !v)}
          className="flex items-center gap-2 text-sm font-bold border border-gray-200 rounded-lg px-3 py-2 hover:border-gray-300 transition"
        >
          <svg viewBox="0 0 24 24" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M3 4h18M6 8h12M9 12h6M11 16h2" />
          </svg>
          Filtros
          {activeCount > 0 && (
            <span className="h-4 w-4 rounded-full brand-bg text-white text-[10px] flex items-center justify-center">
              {activeCount}
            </span>
          )}
        </button>
      </div>

      {filtersOpen && (
        <div className="bg-white border border-gray-100 rounded-xl p-4 mb-6 space-y-4">
          <input
            value={draft.q}
            onChange={(e) => setDraft({ ...draft, q: e.target.value })}
            placeholder="Buscar por nome, descrição ou marca..."
            className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <label className="text-xs font-bold text-gray-600">
              Categoria
              <select
                value={draft.category}
                onChange={(e) => setDraft({ ...draft, category: e.target.value })}
                className="mt-1 w-full h-9 rounded-md border border-gray-300 px-2 text-sm"
              >
                <option value="">Todas</option>
                {facets?.categories.map((c) => (
                  <option key={c.slug} value={c.slug}>{c.name}</option>
                ))}
              </select>
            </label>
            <label className="text-xs font-bold text-gray-600">
              Marca
              <select
                value={draft.brand}
                onChange={(e) => setDraft({ ...draft, brand: e.target.value })}
                className="mt-1 w-full h-9 rounded-md border border-gray-300 px-2 text-sm"
              >
                <option value="">Todas</option>
                {facets?.brands.map((b) => (
                  <option key={b} value={b}>{b}</option>
                ))}
              </select>
            </label>
            <label className="text-xs font-bold text-gray-600">
              Gênero
              <select
                value={draft.gender}
                onChange={(e) => setDraft({ ...draft, gender: e.target.value })}
                className="mt-1 w-full h-9 rounded-md border border-gray-300 px-2 text-sm"
              >
                <option value="">Todos</option>
                {facets?.genders.map((g) => (
                  <option key={g} value={g}>{g}</option>
                ))}
              </select>
            </label>
            <div className="text-xs font-bold text-gray-600">
              Preço (R$)
              <div className="mt-1 flex gap-1.5">
                <input
                  value={draft.minPrice}
                  onChange={(e) => setDraft({ ...draft, minPrice: e.target.value })}
                  placeholder="Mín"
                  inputMode="decimal"
                  className="w-full h-9 rounded-md border border-gray-300 px-2 text-sm"
                />
                <input
                  value={draft.maxPrice}
                  onChange={(e) => setDraft({ ...draft, maxPrice: e.target.value })}
                  placeholder="Máx"
                  inputMode="decimal"
                  className="w-full h-9 rounded-md border border-gray-300 px-2 text-sm"
                />
              </div>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <button onClick={applyFilters} className="text-sm font-bold text-white brand-bg px-4 py-2 rounded-md transition">
              Aplicar filtros
            </button>
            <button onClick={clearFilters} className="text-sm font-medium text-gray-500 hover:text-gray-700">
              Limpar
            </button>
          </div>
        </div>
      )}

      {products === null ? (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-6">
          {Array.from({ length: 8 }).map((_, i) => (
            <div key={i} className="bg-white rounded-xl border border-gray-100 h-72 animate-pulse" />
          ))}
        </div>
      ) : products.length === 0 ? (
        <div className="text-center py-20 text-gray-500">Nenhum produto encontrado com esses filtros.</div>
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
                className="rounded-lg brand-bg px-4 py-2 text-xs font-bold text-white disabled:cursor-not-allowed disabled:opacity-40 transition"
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
