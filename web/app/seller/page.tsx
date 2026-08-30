'use client';

import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useAuth } from '@/components/AuthContext';
import { api } from '@/lib/api';

interface Product {
  id: string;
  name: string;
  price: number;
  state: string;
}

const MANAGE_ROLES = ['SELLER', 'SYSTEM_ADMIN', 'GLOBAL_ADMIN'];

export default function SellerPage() {
  const { user } = useAuth();
  const [products, setProducts] = useState<Product[] | null>(null);
  const [name, setName] = useState('');
  const [sku, setSku] = useState('');
  const [price, setPrice] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [message, setMessage] = useState<{ type: 'ok' | 'error'; text: string } | null>(null);

  const slug = name
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)/g, '');

  const loadProducts = useCallback(async () => {
    const res = await api('/products?limit=20');
    if (res.ok) {
      const data = await res.json();
      setProducts(data.data || []);
    }
  }, []);

  useEffect(() => {
    if (user && MANAGE_ROLES.includes(user.role)) loadProducts();
  }, [user, loadProducts]);

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    setSubmitting(true);
    setMessage(null);
    try {
      const res = await api('/products', {
        method: 'POST',
        body: JSON.stringify({
          name,
          slug: `${slug}-${Date.now()}`,
          sku: sku || `SKU-${Date.now()}`,
          price: Math.round(parseFloat(price.replace(',', '.')) * 100),
        }),
      });
      const data = await res.json();
      if (!res.ok) {
        setMessage({ type: 'error', text: data.message || 'Não foi possível criar o produto.' });
        return;
      }
      setMessage({ type: 'ok', text: `Produto "${data.name}" criado como rascunho. Publique-o para aparecer na vitrine.` });
      setName('');
      setSku('');
      setPrice('');
      loadProducts();
    } catch {
      setMessage({ type: 'error', text: 'Falha de conexão' });
    } finally {
      setSubmitting(false);
    }
  }

  async function handlePublish(id: string) {
    await api(`/products/${id}/publish`, { method: 'POST' });
    loadProducts();
  }

  if (!user) {
    return (
      <div className="max-w-md mx-auto text-center py-20">
        <p className="text-gray-600 mb-4">Entre com uma conta de vendedor ou administrador para acessar este painel.</p>
        <Link href="/login" className="text-indigo-600 hover:underline">Entrar</Link>
      </div>
    );
  }

  if (!MANAGE_ROLES.includes(user.role)) {
    return (
      <div className="max-w-md mx-auto text-center py-20">
        <p className="text-gray-600">Seu perfil ({user.role}) não tem acesso ao painel do vendedor.</p>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto space-y-10">
      <div>
        <h1 className="text-2xl font-bold mb-1">Painel do Vendedor</h1>
        <p className="text-sm text-gray-500 mb-6">Cadastre novos produtos e publique-os na vitrine.</p>

        <form onSubmit={handleCreate} className="bg-white p-6 rounded-xl border border-gray-100 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Nome do produto</label>
            <input
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">SKU (opcional)</label>
              <input
                value={sku}
                onChange={(e) => setSku(e.target.value)}
                className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Preço (R$)</label>
              <input
                required
                inputMode="decimal"
                placeholder="99,90"
                value={price}
                onChange={(e) => setPrice(e.target.value)}
                className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
            </div>
          </div>
          {message && (
            <div className={`text-sm p-3 rounded ${message.type === 'ok' ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-600'}`}>
              {message.text}
            </div>
          )}
          <button
            type="submit"
            disabled={submitting}
            className="px-5 py-2.5 bg-indigo-600 text-white rounded-lg font-medium hover:bg-indigo-700 disabled:opacity-50 transition"
          >
            {submitting ? 'Criando...' : 'Criar produto (rascunho)'}
          </button>
        </form>
      </div>

      <div>
        <h2 className="text-lg font-bold mb-4">Produtos recentes</h2>
        {products === null ? (
          <p className="text-gray-500 text-sm">Carregando...</p>
        ) : products.length === 0 ? (
          <p className="text-gray-500 text-sm">Nenhum produto cadastrado ainda.</p>
        ) : (
          <div className="space-y-2">
            {products.map((p) => (
              <div key={p.id} className="flex items-center justify-between bg-white p-4 rounded-lg border border-gray-100">
                <div>
                  <p className="font-medium text-sm">{p.name}</p>
                  <p className="text-xs text-gray-400">R$ {(p.price / 100).toFixed(2)} · {p.state}</p>
                </div>
                {p.state !== 'ACTIVE' && (
                  <button
                    onClick={() => handlePublish(p.id)}
                    className="text-xs font-bold text-indigo-600 hover:underline"
                  >
                    Publicar
                  </button>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
