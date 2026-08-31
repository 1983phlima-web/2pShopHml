'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/components/AuthContext';
import { AdminSummaryPanels } from '@/components/AdminSummaryPanels';
import { api } from '@/lib/api';

interface HealthReport {
  database_ok: boolean;
  migrations_applied: number;
  total_users: number;
  total_products: number;
  total_orders: number;
  uptime_seconds: number;
  version: string;
}

function formatUptime(seconds: number) {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  return `${h}h ${m}min`;
}

function HealthPanel() {
  const [health, setHealth] = useState<HealthReport | null>(null);

  useEffect(() => {
    api('/analytics/health').then(async (res) => {
      if (res.ok) setHealth(await res.json());
    });
  }, []);

  if (!health) return <p className="text-gray-500 text-sm">Carregando status do sistema...</p>;

  return (
    <div className="bg-white rounded-xl border border-gray-100 p-5">
      <div className="flex items-center gap-2 mb-4">
        <span className={`h-2.5 w-2.5 rounded-full ${health.database_ok ? 'bg-green-500' : 'bg-red-500'}`} />
        <span className="text-sm font-semibold">
          {health.database_ok ? 'Banco de dados operacional' : 'Falha no banco de dados'}
        </span>
      </div>
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 text-sm">
        <div>
          <p className="text-gray-400 text-xs">Migrations aplicadas</p>
          <p className="font-bold">{health.migrations_applied}</p>
        </div>
        <div>
          <p className="text-gray-400 text-xs">Uptime</p>
          <p className="font-bold">{formatUptime(health.uptime_seconds)}</p>
        </div>
        <div>
          <p className="text-gray-400 text-xs">Versão</p>
          <p className="font-bold">{health.version}</p>
        </div>
        <div>
          <p className="text-gray-400 text-xs">Registros</p>
          <p className="font-bold">{health.total_users}u / {health.total_products}p / {health.total_orders}o</p>
        </div>
      </div>
    </div>
  );
}

export default function GlobalAdminPage() {
  const { user } = useAuth();

  if (!user) {
    return (
      <div className="max-w-md mx-auto text-center py-20">
        <p className="text-gray-600 mb-4">Entre com uma conta de administrador global.</p>
        <Link href="/login" className="text-indigo-600 hover:underline">Entrar</Link>
      </div>
    );
  }

  if (user.role !== 'GLOBAL_ADMIN') {
    return (
      <div className="max-w-md mx-auto text-center py-20">
        <p className="text-gray-600">Seu perfil ({user.role}) não tem acesso a esta área.</p>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto space-y-10">
      <div>
        <h1 className="text-2xl font-bold mb-1">Administração Global</h1>
        <p className="text-sm text-gray-500 mb-6">
          Saúde da infraestrutura, banco de dados e integrações de pagamento (Stripe — sandbox em HML).
        </p>
        <HealthPanel />
      </div>
      <AdminSummaryPanels canEditTheme />
    </div>
  );
}
