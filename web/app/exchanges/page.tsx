'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/components/AuthContext';
import { api } from '@/lib/api';

interface ExchangeRequest {
  id: string;
  order_id: string;
  product_id: string;
  reason: string;
  status: string;
  created_at: string;
}

const STATUS_LABELS: Record<string, string> = {
  REQUESTED: 'Solicitado',
  APPROVED: 'Aprovado',
  REJECTED: 'Recusado',
  COMPLETED: 'Concluído',
};

const STATUS_COLORS: Record<string, string> = {
  REQUESTED: 'bg-amber-50 text-amber-700',
  APPROVED: 'bg-blue-50 text-blue-700',
  REJECTED: 'bg-red-50 text-red-700',
  COMPLETED: 'bg-green-50 text-green-700',
};

export default function ExchangesPage() {
  const { user } = useAuth();
  const [requests, setRequests] = useState<ExchangeRequest[] | null>(null);

  useEffect(() => {
    if (!user) return;
    api('/exchanges/mine').then(async (res) => {
      if (res.ok) setRequests((await res.json()) || []);
      else setRequests([]);
    });
  }, [user]);

  if (!user) {
    return (
      <div className="max-w-md mx-auto text-center py-20">
        <p className="text-gray-600 mb-4">Entre na sua conta para ver suas trocas.</p>
        <Link href="/login" className="text-indigo-600 hover:underline">Entrar</Link>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Minhas Trocas</h1>
        <Link href="/orders" className="text-sm font-medium brand-text hover:underline">← Minhas Compras</Link>
      </div>
      {requests === null ? (
        <p className="text-gray-500">Carregando...</p>
      ) : requests.length === 0 ? (
        <div className="text-center py-16 text-gray-500 text-sm">
          Você ainda não solicitou nenhuma troca ou devolução.
          <br />
          Isso pode ser feito a partir de um item entregue em "Minhas Compras".
        </div>
      ) : (
        <div className="space-y-3">
          {requests.map((r) => (
            <div key={r.id} className="bg-white p-4 rounded-lg border border-gray-100">
              <div className="flex items-center justify-between mb-2">
                <span className="text-xs font-mono text-gray-400">Pedido {r.order_id.slice(0, 8)}...</span>
                <span className={`text-xs font-medium px-2 py-1 rounded ${STATUS_COLORS[r.status] || 'bg-gray-50 text-gray-600'}`}>
                  {STATUS_LABELS[r.status] || r.status}
                </span>
              </div>
              <p className="text-sm text-gray-700 mb-1">{r.reason}</p>
              <p className="text-xs text-gray-400">{new Date(r.created_at).toLocaleDateString('pt-BR')}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
