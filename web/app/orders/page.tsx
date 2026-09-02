'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/components/AuthContext';
import { api } from '@/lib/api';

interface OrderItem {
  product_id: string;
  name: string;
  quantity: number;
  unit_price: number;
  total: number;
}

interface Order {
  id: string;
  items: OrderItem[];
  total: number;
  status: string;
  created_at: string;
}

const STATUS_LABELS: Record<string, string> = {
  PENDING: 'Pendente',
  CONFIRMED: 'Confirmado',
  PAID: 'Pago',
  SHIPPED: 'Enviado',
  DELIVERED: 'Entregue',
  CANCELLED: 'Cancelado',
  REFUNDED: 'Reembolsado',
};

function ExchangeForm({ orderId, productId, onDone }: { orderId: string; productId: string; onDone: () => void }) {
  const [reason, setReason] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [done, setDone] = useState(false);

  async function submit() {
    if (!reason.trim()) return;
    setSubmitting(true);
    try {
      await api('/exchanges', {
        method: 'POST',
        body: JSON.stringify({ order_id: orderId, product_id: productId, reason }),
      });
      setDone(true);
      setTimeout(onDone, 1200);
    } finally {
      setSubmitting(false);
    }
  }

  if (done) return <p className="text-xs text-green-600 mt-1">Solicitação enviada! Acompanhe em "Minhas Trocas".</p>;

  return (
    <div className="mt-2 flex flex-col sm:flex-row gap-2">
      <input
        value={reason}
        onChange={(e) => setReason(e.target.value)}
        placeholder="Motivo da troca/devolução..."
        className="flex-1 border border-gray-300 rounded-md px-2 py-1.5 text-xs focus:outline-none focus:ring-2 focus:ring-indigo-500"
      />
      <button
        onClick={submit}
        disabled={submitting}
        className="text-xs font-bold text-white brand-bg px-3 py-1.5 rounded-md disabled:opacity-50"
      >
        {submitting ? 'Enviando...' : 'Confirmar'}
      </button>
    </div>
  );
}

export default function OrdersPage() {
  const { user } = useAuth();
  const [orders, setOrders] = useState<Order[] | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [exchangeItem, setExchangeItem] = useState<{ orderId: string; productId: string } | null>(null);

  useEffect(() => {
    if (!user) return;
    api('/orders')
      .then(async (res) => {
        if (!res.ok) {
          setError('Não foi possível carregar seus pedidos.');
          return;
        }
        const data = await res.json();
        setOrders(data || []);
      })
      .catch(() => setError('Falha de conexão'));
  }, [user]);

  if (!user) {
    return (
      <div className="max-w-md mx-auto text-center py-20">
        <p className="text-gray-600 mb-4">Entre na sua conta para ver suas compras.</p>
        <Link href="/login" className="text-indigo-600 hover:underline">Entrar</Link>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Minhas Compras</h1>
        <Link href="/exchanges" className="text-sm font-medium brand-text hover:underline">Minhas Trocas →</Link>
      </div>
      {error && <div className="text-red-600 text-sm bg-red-50 p-3 rounded mb-4">{error}</div>}
      {orders === null ? (
        <p className="text-gray-500">Carregando...</p>
      ) : orders.length === 0 ? (
        <p className="text-gray-500">Você ainda não fez nenhuma compra.</p>
      ) : (
        <div className="space-y-4">
          {orders.map((order) => (
            <div key={order.id} className="bg-white p-5 rounded-xl border border-gray-100">
              <div className="flex justify-between items-start mb-3">
                <div>
                  <p className="text-xs text-gray-400">Pedido</p>
                  <p className="font-mono text-sm">{order.id}</p>
                </div>
                <span className="text-xs font-medium bg-indigo-50 text-indigo-700 px-2 py-1 rounded">
                  {STATUS_LABELS[order.status] || order.status}
                </span>
              </div>
              <ul className="space-y-2 mb-3">
                {order.items.map((item, i) => {
                  const isEditing = exchangeItem?.orderId === order.id && exchangeItem?.productId === item.product_id;
                  return (
                    <li key={i} className="text-sm text-gray-600">
                      <div className="flex justify-between items-center">
                        <span>{item.quantity}x {item.name}</span>
                        <div className="flex items-center gap-2">
                          <span>R$ {(item.total / 100).toFixed(2)}</span>
                          {order.status === 'DELIVERED' && (
                            <button
                              onClick={() => setExchangeItem(isEditing ? null : { orderId: order.id, productId: item.product_id })}
                              className="text-xs text-gray-400 hover:text-rose-600 underline"
                            >
                              {isEditing ? 'cancelar' : 'trocar'}
                            </button>
                          )}
                        </div>
                      </div>
                      {isEditing && (
                        <ExchangeForm orderId={order.id} productId={item.product_id} onDone={() => setExchangeItem(null)} />
                      )}
                    </li>
                  );
                })}
              </ul>
              <div className="flex justify-between items-center pt-3 border-t">
                <span className="text-xs text-gray-400">
                  {new Date(order.created_at).toLocaleDateString('pt-BR')}
                </span>
                <span className="font-bold text-indigo-600">R$ {(order.total / 100).toFixed(2)}</span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
