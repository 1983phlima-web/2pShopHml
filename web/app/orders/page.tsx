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

export default function OrdersPage() {
  const { user } = useAuth();
  const [orders, setOrders] = useState<Order[] | null>(null);
  const [error, setError] = useState<string | null>(null);

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
        <p className="text-gray-600 mb-4">Entre na sua conta para ver seus pedidos.</p>
        <Link href="/login" className="text-indigo-600 hover:underline">Entrar</Link>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Meus Pedidos</h1>
      {error && <div className="text-red-600 text-sm bg-red-50 p-3 rounded mb-4">{error}</div>}
      {orders === null ? (
        <p className="text-gray-500">Carregando...</p>
      ) : orders.length === 0 ? (
        <p className="text-gray-500">Você ainda não fez nenhum pedido.</p>
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
              <ul className="space-y-1 mb-3">
                {order.items.map((item, i) => (
                  <li key={i} className="text-sm text-gray-600 flex justify-between">
                    <span>{item.quantity}x {item.name}</span>
                    <span>R$ {(item.total / 100).toFixed(2)}</span>
                  </li>
                ))}
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
