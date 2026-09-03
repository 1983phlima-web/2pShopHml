'use client';

import { useEffect, useState, useCallback } from 'react';
import { api } from '@/lib/api';

interface Order {
  id: string;
  status: string;
  total: number;
  created_at: string;
}

const FLOW: string[] = ['PENDING', 'CONFIRMED', 'SHIPPED', 'DELIVERED'];
const STATUS_LABELS: Record<string, string> = {
  PENDING: 'Pendente',
  CONFIRMED: 'Confirmado',
  PAID: 'Pago',
  SHIPPED: 'Enviado',
  DELIVERED: 'Entregue',
  CANCELLED: 'Cancelado',
  REFUNDED: 'Reembolsado',
};

export function SellerOrdersPanel() {
  const [orders, setOrders] = useState<Order[] | null>(null);
  const [updating, setUpdating] = useState<string | null>(null);

  const load = useCallback(async () => {
    const res = await api('/orders/all');
    if (res.ok) setOrders((await res.json()) || []);
    else setOrders([]);
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  async function advance(orderId: string, currentStatus: string) {
    const idx = FLOW.indexOf(currentStatus);
    if (idx === -1 || idx === FLOW.length - 1) return;
    const nextStatus = FLOW[idx + 1];
    setUpdating(orderId);
    try {
      await api(`/orders/${orderId}/status`, { method: 'PUT', body: JSON.stringify({ status: nextStatus }) });
      load();
    } finally {
      setUpdating(null);
    }
  }

  if (orders === null) return <p className="text-gray-500 text-sm">Carregando...</p>;
  if (orders.length === 0) return <p className="text-gray-500 text-sm">Nenhum pedido ainda.</p>;

  return (
    <div className="space-y-3">
      {orders.map((o) => {
        const idx = FLOW.indexOf(o.status);
        const canAdvance = idx !== -1 && idx < FLOW.length - 1;
        return (
          <div key={o.id} className="bg-white p-4 rounded-lg border border-gray-100">
            <div className="flex items-center justify-between mb-2">
              <span className="text-xs font-mono text-gray-400">{o.id.slice(0, 8)}...</span>
              <span className="text-xs font-medium bg-indigo-50 text-indigo-700 px-2 py-1 rounded">
                {STATUS_LABELS[o.status] || o.status}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm font-semibold">R$ {(o.total / 100).toFixed(2)}</span>
              {canAdvance && (
                <button
                  onClick={() => advance(o.id, o.status)}
                  disabled={updating === o.id}
                  className="text-xs font-bold brand-text hover:underline disabled:opacity-50"
                >
                  {updating === o.id ? 'Atualizando...' : `Mover para ${STATUS_LABELS[FLOW[idx + 1]]} →`}
                </button>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}
