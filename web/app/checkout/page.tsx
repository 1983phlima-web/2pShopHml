'use client';

import { useCart } from '@/components/CartContext';
import { useState } from 'react';
import { api } from '@/lib/api';

export default function CheckoutPage() {
  const { items, clear } = useCart();
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<{ order_id: string; status: string } | null>(null);
  const [error, setError] = useState<string | null>(null);

  const total = items.reduce((sum, item) => sum + item.price * item.quantity, 0);

  async function handleCheckout() {
    setLoading(true);
    setError(null);
    try {
      const res = await api('/checkout', {
        method: 'POST',
        body: JSON.stringify({
          tenant_id: process.env.NEXT_PUBLIC_TENANT_ID || 'tenant_01',
          customer_id: 'user_01',
          items: items.map((i) => ({ product_id: i.id, quantity: i.quantity })),
          payment_method: { type: 'card', token: 'tok_visa' },
          idempotency_key: `checkout-${Date.now()}`,
        }),
      });
      const data = await res.json();
      if (!res.ok) {
        setError(data.message || 'Erro no checkout');
      } else {
        setResult(data);
        clear();
      }
    } catch (e) {
      setError('Falha de conexão');
    } finally {
      setLoading(false);
    }
  }

  if (result) {
    return (
      <div className="max-w-md mx-auto text-center py-20">
        <div className="text-green-600 text-5xl mb-4">✓</div>
        <h2 className="text-2xl font-bold mb-2">Pedido realizado!</h2>
        <p className="text-gray-600 mb-4">Order ID: {result.order_id}</p>
        <a href="/products" className="text-indigo-600 hover:underline">Continuar comprando</a>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Checkout</h1>
      {items.length === 0 ? (
        <p className="text-gray-500">Seu carrinho está vazio.</p>
      ) : (
        <div className="space-y-4">
          {items.map((item) => (
            <div key={item.id} className="flex justify-between items-center bg-white p-4 rounded-lg border">
              <div>
                <p className="font-medium">{item.name}</p>
                <p className="text-sm text-gray-500">Qtd: {item.quantity}</p>
              </div>
              <p className="font-semibold">R$ {((item.price * item.quantity) / 100).toFixed(2)}</p>
            </div>
          ))}
          <div className="flex justify-between items-center pt-4 border-t">
            <span className="text-lg font-bold">Total</span>
            <span className="text-xl font-bold text-indigo-600">R$ {(total / 100).toFixed(2)}</span>
          </div>
          {error && <div className="text-red-600 text-sm bg-red-50 p-3 rounded">{error}</div>}
          <button
            onClick={handleCheckout}
            disabled={loading}
            className="w-full py-3 bg-indigo-600 text-white rounded-lg font-medium hover:bg-indigo-700 disabled:opacity-50 transition"
          >
            {loading ? 'Processando...' : 'Finalizar compra'}
          </button>
        </div>
      )}
    </div>
  );
}
