'use client';

import Link from 'next/link';
import { useAuth } from '@/components/AuthContext';
import { AdminSummaryPanels } from '@/components/AdminSummaryPanels';
import { HealthHistoryPanel } from '@/components/HealthHistoryPanel';

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
        <HealthHistoryPanel />
      </div>
      <AdminSummaryPanels canEditTheme />
    </div>
  );
}
