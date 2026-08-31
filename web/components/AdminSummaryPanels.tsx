'use client';

import { useEffect, useState } from 'react';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts';
import { api } from '@/lib/api';
import { PALETTES, applyPalette } from '@/lib/theme';

interface AdminSummary {
  users_by_role: Record<string, number>;
  products_by_state: Record<string, number>;
  orders_by_status: Record<string, number>;
  total_revenue: number;
  total_reviews: number;
}

const ROLE_LABELS: Record<string, string> = {
  SELLER: 'Vendedores',
  BUYER: 'Clientes',
  SYSTEM_ADMIN: 'Admins Sistema',
  GLOBAL_ADMIN: 'Admins Global',
};

const COLORS = ['#4f46e5', '#059669', '#d97706', '#e11d48', '#8350a9'];

function StatCard({ label, value }: { label: string; value: string | number }) {
  return (
    <div className="bg-white rounded-xl border border-gray-100 p-4">
      <p className="text-xs font-semibold text-gray-500">{label}</p>
      <p className="text-2xl font-extrabold mt-1 brand-text">{value}</p>
    </div>
  );
}

function DistributionChart({ title, data }: { title: string; data: Record<string, number> }) {
  const entries = Object.entries(data).map(([name, value]) => ({ name, value }));
  if (entries.length === 0) return null;
  return (
    <div className="bg-white rounded-xl border border-gray-100 p-4">
      <p className="text-sm font-semibold mb-3">{title}</p>
      <ResponsiveContainer width="100%" height={220}>
        <PieChart>
          <Pie data={entries} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={70} label>
            {entries.map((_, i) => (
              <Cell key={i} fill={COLORS[i % COLORS.length]} />
            ))}
          </Pie>
          <Tooltip />
          <Legend wrapperStyle={{ fontSize: 11 }} />
        </PieChart>
      </ResponsiveContainer>
    </div>
  );
}

export function AdminSummaryPanels({ canEditTheme }: { canEditTheme: boolean }) {
  const [summary, setSummary] = useState<AdminSummary | null>(null);
  const [palette, setPalette] = useState<string>('indigo');
  const [savingPalette, setSavingPalette] = useState(false);

  useEffect(() => {
    api('/analytics/admin-summary').then(async (res) => {
      if (res.ok) setSummary(await res.json());
    });
    api('/settings/theme').then(async (res) => {
      if (res.ok) {
        const data = await res.json();
        setPalette(data.palette || 'indigo');
      }
    });
  }, []);

  async function selectPalette(key: string) {
    setPalette(key);
    applyPalette(key);
    if (!canEditTheme) return;
    setSavingPalette(true);
    try {
      await api('/settings/theme', { method: 'PUT', body: JSON.stringify({ palette: key }) });
    } finally {
      setSavingPalette(false);
    }
  }

  const usersByRoleLabeled = summary
    ? Object.fromEntries(Object.entries(summary.users_by_role).map(([k, v]) => [ROLE_LABELS[k] || k, v]))
    : {};

  return (
    <div className="space-y-8">
      <section>
        <h2 className="text-lg font-bold mb-4">Métricas da plataforma</h2>
        {summary ? (
          <>
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
              <StatCard label="Receita total" value={`R$ ${(summary.total_revenue / 100).toFixed(2)}`} />
              <StatCard label="Usuários" value={Object.values(summary.users_by_role).reduce((a, b) => a + b, 0)} />
              <StatCard label="Produtos" value={Object.values(summary.products_by_state).reduce((a, b) => a + b, 0)} />
              <StatCard label="Avaliações" value={summary.total_reviews} />
            </div>
            <div className="grid md:grid-cols-3 gap-4">
              <DistributionChart title="Usuários por papel" data={usersByRoleLabeled} />
              <DistributionChart title="Produtos por status" data={summary.products_by_state} />
              <DistributionChart title="Pedidos por status" data={summary.orders_by_status} />
            </div>
          </>
        ) : (
          <p className="text-gray-500 text-sm">Carregando métricas...</p>
        )}
      </section>

      <section>
        <h2 className="text-lg font-bold mb-1">Paleta visual da loja</h2>
        <p className="text-sm text-gray-500 mb-4">
          {canEditTheme
            ? 'Escolha a paleta aplicada em toda a vitrine.'
            : 'Somente administradores podem salvar a paleta (você pode pré-visualizar).'}
          {savingPalette && ' Salvando...'}
        </p>
        <div className="grid grid-cols-2 sm:grid-cols-5 gap-3">
          {Object.values(PALETTES).map((p) => (
            <button
              key={p.key}
              onClick={() => selectPalette(p.key)}
              className={`rounded-xl border-2 p-3 text-left transition ${palette === p.key ? 'border-gray-900' : 'border-gray-100 hover:border-gray-300'}`}
            >
              <div className="flex gap-1 mb-2">
                <span className="h-6 w-6 rounded-full" style={{ backgroundColor: p.primary }} />
                <span className="h-6 w-6 rounded-full" style={{ backgroundColor: p.dark }} />
                <span className="h-6 w-6 rounded-full" style={{ backgroundColor: p.accent }} />
              </div>
              <span className="text-xs font-semibold">{p.label}</span>
            </button>
          ))}
        </div>
      </section>
    </div>
  );
}
