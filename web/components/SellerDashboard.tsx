'use client';

import { useEffect, useState } from 'react';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts';
import { api } from '@/lib/api';

interface DailySales {
  date: string;
  revenue: number;
  orders: number;
}

interface ProductSales {
  product_id: string;
  name: string;
  units: number;
  revenue: number;
}

interface SellerSummary {
  total_revenue: number;
  total_orders: number;
  total_units_sold: number;
  pending_orders: number;
  sales_by_day: DailySales[];
  top_products: ProductSales[];
}

function KpiCard({ label, value, detail, tone }: { label: string; value: string; detail: string; tone: string }) {
  return (
    <div className="bg-white rounded-xl border border-gray-100 p-4">
      <p className="text-xs font-semibold text-gray-500">{label}</p>
      <p className="text-2xl font-extrabold mt-1" style={{ color: tone }}>{value}</p>
      <p className="text-xs text-gray-400 mt-1">{detail}</p>
    </div>
  );
}

export function SellerDashboard() {
  const [summary, setSummary] = useState<SellerSummary | null>(null);

  useEffect(() => {
    api('/analytics/seller-summary')
      .then(async (res) => {
        if (res.ok) setSummary(await res.json());
      })
      .catch(() => {});
  }, []);

  if (!summary) return null;

  const chartData = summary.sales_by_day.map((d) => ({
    day: d.date.slice(5), // MM-DD
    Vendas: d.revenue / 100,
  }));

  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-bold">Painel de Vendas</h2>
        <span className="text-xs text-gray-400">Últimos 14 dias</span>
      </div>

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <KpiCard
          label="Vendas totais"
          value={`R$ ${(summary.total_revenue / 100).toFixed(2)}`}
          detail={`${summary.total_orders} pedidos`}
          tone="var(--brand-primary)"
        />
        <KpiCard
          label="Unidades vendidas"
          value={String(summary.total_units_sold)}
          detail="itens em todos os pedidos"
          tone="#059669"
        />
        <KpiCard
          label="Pedidos pendentes"
          value={String(summary.pending_orders)}
          detail="aguardando confirmação"
          tone="#d97706"
        />
        <KpiCard
          label="Ticket médio"
          value={`R$ ${summary.total_orders ? (summary.total_revenue / summary.total_orders / 100).toFixed(2) : '0,00'}`}
          detail="por pedido"
          tone="#8350a9"
        />
      </div>

      {chartData.some((d) => d.Vendas > 0) && (
        <div className="bg-white rounded-xl border border-gray-100 p-4">
          <p className="text-sm font-semibold mb-3">Vendas por dia (R$)</p>
          <ResponsiveContainer width="100%" height={200}>
            <BarChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f0f0f0" />
              <XAxis dataKey="day" tick={{ fontSize: 11 }} />
              <YAxis tick={{ fontSize: 11 }} width={40} />
              <Tooltip formatter={(v: number) => `R$ ${v.toFixed(2)}`} />
              <Bar dataKey="Vendas" fill="var(--brand-primary)" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>
      )}

      {summary.top_products.length > 0 && (
        <div className="bg-white rounded-xl border border-gray-100 p-4">
          <p className="text-sm font-semibold mb-3">Produtos mais vendidos</p>
          <div className="space-y-2">
            {summary.top_products.map((p) => (
              <div key={p.product_id} className="flex items-center justify-between text-sm">
                <span className="text-gray-700 truncate">{p.name}</span>
                <span className="text-gray-400 shrink-0 ml-3">{p.units} un · R$ {(p.revenue / 100).toFixed(2)}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </section>
  );
}
