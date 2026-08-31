'use client';

import { useEffect, useState, useCallback } from 'react';
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid, Legend } from 'recharts';
import { api } from '@/lib/api';

interface HealthPoint {
  recorded_at: string;
  db_ok: boolean;
  db_latency_ms: number;
  total_users: number;
  total_products: number;
  total_orders: number;
  total_revenue: number;
}

const PERIODS = [
  { key: '24h', label: '24 horas' },
  { key: '7d', label: '7 dias' },
  { key: '30d', label: '30 dias' },
];

function formatTick(iso: string, period: string) {
  const d = new Date(iso);
  if (period === '24h') return d.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' });
  return d.toLocaleDateString('pt-BR', { day: '2-digit', month: '2-digit' });
}

export function HealthHistoryPanel() {
  const [period, setPeriod] = useState('24h');
  const [points, setPoints] = useState<HealthPoint[] | null>(null);

  const load = useCallback((p: string) => {
    setPoints(null);
    api(`/analytics/health-history?period=${p}`).then(async (res) => {
      if (res.ok) setPoints((await res.json()) || []);
      else setPoints([]);
    });
  }, []);

  useEffect(() => {
    load(period);
  }, [period, load]);

  const chartData = (points || []).map((p) => ({
    t: formatTick(p.recorded_at, period),
    'Latência (ms)': Math.round(p.db_latency_ms * 100) / 100,
    Usuários: p.total_users,
    Produtos: p.total_products,
    Pedidos: p.total_orders,
  }));

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm font-semibold">Histórico de saúde</p>
        <div className="flex gap-1 bg-gray-100 rounded-lg p-1">
          {PERIODS.map((p) => (
            <button
              key={p.key}
              onClick={() => setPeriod(p.key)}
              className={`text-xs font-bold px-3 py-1.5 rounded-md transition ${
                period === p.key ? 'bg-white shadow-sm brand-text' : 'text-gray-500'
              }`}
            >
              {p.label}
            </button>
          ))}
        </div>
      </div>

      {points === null ? (
        <p className="text-gray-500 text-sm">Carregando histórico...</p>
      ) : points.length < 2 ? (
        <div className="bg-white rounded-xl border border-gray-100 p-6 text-center">
          <p className="text-sm text-gray-500">
            Ainda não há histórico suficiente para este período — o sistema registra uma amostra real a cada 5
            minutos. Volte em alguns minutos para ver a curva se formando.
          </p>
        </div>
      ) : (
        <>
          <div className="bg-white rounded-xl border border-gray-100 p-4">
            <p className="text-xs font-semibold text-gray-500 mb-3">Latência do banco de dados (ms)</p>
            <ResponsiveContainer width="100%" height={200}>
              <LineChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f0f0f0" />
                <XAxis dataKey="t" tick={{ fontSize: 11 }} />
                <YAxis tick={{ fontSize: 11 }} width={40} />
                <Tooltip />
                <Line type="monotone" dataKey="Latência (ms)" stroke="var(--brand-primary)" strokeWidth={2} dot={false} />
              </LineChart>
            </ResponsiveContainer>
          </div>

          <div className="bg-white rounded-xl border border-gray-100 p-4">
            <p className="text-xs font-semibold text-gray-500 mb-3">Crescimento da plataforma</p>
            <ResponsiveContainer width="100%" height={220}>
              <LineChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f0f0f0" />
                <XAxis dataKey="t" tick={{ fontSize: 11 }} />
                <YAxis tick={{ fontSize: 11 }} width={40} />
                <Tooltip />
                <Legend wrapperStyle={{ fontSize: 11 }} />
                <Line type="monotone" dataKey="Usuários" stroke="#059669" strokeWidth={2} dot={false} />
                <Line type="monotone" dataKey="Produtos" stroke="#d97706" strokeWidth={2} dot={false} />
                <Line type="monotone" dataKey="Pedidos" stroke="#e11d48" strokeWidth={2} dot={false} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </>
      )}
    </div>
  );
}
