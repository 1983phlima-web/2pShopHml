'use client';

import { useEffect, useState, useCallback } from 'react';
import { api } from '@/lib/api';

interface ExchangeRequest {
  id: string;
  order_id: string;
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

export function SellerExchangesPanel() {
  const [requests, setRequests] = useState<ExchangeRequest[] | null>(null);
  const [updating, setUpdating] = useState<string | null>(null);

  const load = useCallback(async () => {
    const res = await api('/exchanges/all');
    if (res.ok) setRequests((await res.json()) || []);
    else setRequests([]);
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  async function updateStatus(id: string, status: string) {
    setUpdating(id);
    try {
      await api(`/exchanges/${id}/status`, { method: 'PUT', body: JSON.stringify({ status }) });
      load();
    } finally {
      setUpdating(null);
    }
  }

  if (requests === null) return <p className="text-gray-500 text-sm">Carregando...</p>;
  if (requests.length === 0) return <p className="text-gray-500 text-sm">Nenhuma solicitação de troca ainda.</p>;

  return (
    <div className="space-y-3">
      {requests.map((r) => (
        <div key={r.id} className="bg-white p-4 rounded-lg border border-gray-100">
          <div className="flex items-center justify-between mb-2">
            <span className="text-xs font-mono text-gray-400">Pedido {r.order_id.slice(0, 8)}...</span>
            <span className="text-xs font-medium bg-amber-50 text-amber-700 px-2 py-1 rounded">
              {STATUS_LABELS[r.status] || r.status}
            </span>
          </div>
          <p className="text-sm text-gray-700 mb-3">{r.reason}</p>
          {r.status === 'REQUESTED' && (
            <div className="flex gap-2">
              <button
                onClick={() => updateStatus(r.id, 'APPROVED')}
                disabled={updating === r.id}
                className="text-xs font-bold text-green-700 bg-green-50 px-3 py-1.5 rounded-md hover:bg-green-100 disabled:opacity-50"
              >
                Aprovar
              </button>
              <button
                onClick={() => updateStatus(r.id, 'REJECTED')}
                disabled={updating === r.id}
                className="text-xs font-bold text-red-700 bg-red-50 px-3 py-1.5 rounded-md hover:bg-red-100 disabled:opacity-50"
              >
                Recusar
              </button>
            </div>
          )}
          {r.status === 'APPROVED' && (
            <button
              onClick={() => updateStatus(r.id, 'COMPLETED')}
              disabled={updating === r.id}
              className="text-xs font-bold brand-text hover:underline disabled:opacity-50"
            >
              Marcar como concluído
            </button>
          )}
        </div>
      ))}
    </div>
  );
}
