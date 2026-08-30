'use client';

import Image from 'next/image';
import Link from 'next/link';
import { useAuth } from './AuthContext';
import { useCart } from './CartContext';
import { ROLE_LABELS } from '@/lib/auth';

export function Header() {
  const { user, logout } = useAuth();
  const { items } = useCart();
  const cartCount = items.reduce((sum, i) => sum + i.quantity, 0);

  return (
    <header className="bg-white border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <Link href="/" className="flex items-center gap-2">
          <Image src="/logo.png" alt="2pShop" width={36} height={36} className="rounded-md" priority />
          <span className="text-xl font-bold text-indigo-600">2pShop</span>
        </Link>
        <nav className="flex items-center gap-6">
          <Link href="/products" className="text-sm font-medium hover:text-indigo-600">Produtos</Link>
          {user && (
            <Link href="/orders" className="text-sm font-medium hover:text-indigo-600">Meus Pedidos</Link>
          )}
          <Link href="/checkout" className="text-sm font-medium hover:text-indigo-600 relative">
            Carrinho
            {cartCount > 0 && (
              <span className="absolute -top-2 -right-3 bg-indigo-600 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center">
                {cartCount}
              </span>
            )}
          </Link>
          {user ? (
            <div className="flex items-center gap-3">
              <span className="text-sm text-gray-500">
                {user.name} <span className="text-xs bg-indigo-50 text-indigo-700 px-1.5 py-0.5 rounded">{ROLE_LABELS[user.role]}</span>
              </span>
              <button onClick={logout} className="text-sm font-medium text-gray-500 hover:text-red-600">Sair</button>
            </div>
          ) : (
            <div className="flex items-center gap-3">
              <Link href="/login" className="text-sm font-medium hover:text-indigo-600">Entrar</Link>
              <Link href="/register" className="text-sm font-medium bg-indigo-600 text-white px-3 py-1.5 rounded-md hover:bg-indigo-700">Criar conta</Link>
            </div>
          )}
        </nav>
      </div>
    </header>
  );
}
