'use client';

import Link from 'next/link';
import { useAuth } from '@/components/AuthContext';
import { AdminSummaryPanels } from '@/components/AdminSummaryPanels';

export default function SystemAdminPage() {
  const { user } = useAuth();

  if (!user) {
    return (
      <div className="max-w-md mx-auto text-center py-20">
        <p className="text-gray-600 mb-4">Entre com uma conta de administrador.</p>
        <Link href="/login" className="text-indigo-600 hover:underline">Entrar</Link>
      </div>
    );
  }

  if (!['SYSTEM_ADMIN', 'GLOBAL_ADMIN'].includes(user.role)) {
    return (
      <div className="max-w-md mx-auto text-center py-20">
        <p className="text-gray-600">Seu perfil ({user.role}) não tem acesso a esta área.</p>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-1">Administração do Sistema</h1>
      <p className="text-sm text-gray-500 mb-8">Parametrizações do sistema e métricas gerais da plataforma.</p>
      <AdminSummaryPanels canEditTheme />
    </div>
  );
}
