'use client';

import Image from 'next/image';
import Link from 'next/link';
import { useState } from 'react';
import { useAuth } from './AuthContext';
import { useCart } from './CartContext';
import { useFavorites } from './FavoritesContext';
import { Avatar } from './Avatar';
import { ProfilePopup } from './ProfilePopup';
import { useLanguage } from './LanguageContext';
import { ROLE_LABELS } from '@/lib/auth';

export function Header() {
  const { user, logout } = useAuth();
  const { items } = useCart();
  const { ids: favoriteIds } = useFavorites();
  const { t } = useLanguage();
  const [menuOpen, setMenuOpen] = useState(false);
  const [profileOpen, setProfileOpen] = useState(false);
  const cartCount = items.reduce((sum, i) => sum + i.quantity, 0);

  const canManage = user && ['SELLER', 'SYSTEM_ADMIN', 'GLOBAL_ADMIN'].includes(user.role);
  const isSystemAdmin = user && ['SYSTEM_ADMIN', 'GLOBAL_ADMIN'].includes(user.role);
  const isGlobalAdmin = user?.role === 'GLOBAL_ADMIN';
  const isBuyer = user?.role === 'BUYER';

  const navLinks = (
    <>
      <Link href="/products" className="text-sm font-medium hover:text-indigo-600" onClick={() => setMenuOpen(false)}>
        {t('nav.products')}
      </Link>
      {isBuyer && (
        <>
          <Link href="/orders" className="text-sm font-medium hover:text-indigo-600" onClick={() => setMenuOpen(false)}>
            {t('nav.myOrders')}
          </Link>
          <Link href="/exchanges" className="text-sm font-medium hover:text-indigo-600" onClick={() => setMenuOpen(false)}>
            {t('nav.myExchanges')}
          </Link>
          <Link href="/comments" className="text-sm font-medium hover:text-indigo-600" onClick={() => setMenuOpen(false)}>
            {t('nav.myComments')}
          </Link>
        </>
      )}
      {user && (
        <Link href="/favorites" className="text-sm font-medium hover:text-indigo-600 relative" onClick={() => setMenuOpen(false)}>
          {t('nav.favorites')}
          {favoriteIds.size > 0 && (
            <span className="absolute -top-2 -right-3 bg-rose-600 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center">
              {favoriteIds.size}
            </span>
          )}
        </Link>
      )}
      {canManage && (
        <Link href="/seller" className="text-sm font-medium hover:text-indigo-600" onClick={() => setMenuOpen(false)}>
          {t('nav.sellerPanel')}
        </Link>
      )}
      {isSystemAdmin && (
        <Link href="/admin/system" className="text-sm font-medium hover:text-indigo-600" onClick={() => setMenuOpen(false)}>
          {t('nav.systemAdmin')}
        </Link>
      )}
      {isGlobalAdmin && (
        <Link href="/admin/global" className="text-sm font-medium hover:text-indigo-600" onClick={() => setMenuOpen(false)}>
          {t('nav.globalAdmin')}
        </Link>
      )}
      <Link href="/checkout" className="text-sm font-medium hover:text-indigo-600 relative" onClick={() => setMenuOpen(false)}>
        {t('nav.cart')}
        {cartCount > 0 && (
          <span className="absolute -top-2 -right-3 brand-bg text-white text-xs rounded-full h-4 w-4 flex items-center justify-center">
            {cartCount}
          </span>
        )}
      </Link>
    </>
  );

  const authLinks = user ? (
    <div className="flex items-center gap-3">
      <button
        onClick={() => { setProfileOpen(true); setMenuOpen(false); }}
        className="flex items-center gap-2 hover:opacity-80 transition"
      >
        <Avatar avatar={user.avatar} size={32} />
        <span className="text-sm text-gray-500 hidden sm:inline">
          {user.name}{' '}
          <span className="text-xs bg-indigo-50 text-indigo-700 px-1.5 py-0.5 rounded">{ROLE_LABELS[user.role]}</span>
        </span>
      </button>
      <button onClick={() => { logout(); setMenuOpen(false); }} className="text-sm font-medium text-gray-500 hover:text-red-600">
        {t('auth.logout')}
      </button>
    </div>
  ) : (
    <div className="flex items-center gap-3">
      <Link href="/login" className="text-sm font-medium hover:text-indigo-600" onClick={() => setMenuOpen(false)}>{t('auth.login')}</Link>
      <Link href="/register" className="text-sm font-medium bg-indigo-600 text-white px-3 py-1.5 rounded-md hover:bg-indigo-700" onClick={() => setMenuOpen(false)}>
        {t('auth.register')}
      </Link>
    </div>
  );

  return (
    <header className="bg-white border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-20 flex items-center justify-between">
        <Link href="/" className="flex items-center gap-3 shrink-0">
          <Image src="/logo.png" alt="2pShop" width={546} height={633} className="h-14 w-auto" priority />
          <span className="text-2xl font-extrabold brand-text tracking-tight">2pShop</span>
        </Link>

        {/* Desktop nav */}
        <nav className="hidden md:flex items-center gap-6">
          {navLinks}
          {authLinks}
        </nav>

        {/* Mobile: avatar + hamburger */}
        <div className="md:hidden flex items-center gap-2">
          {user && (
            <button onClick={() => setProfileOpen(true)} aria-label="Meu perfil">
              <Avatar avatar={user.avatar} size={32} />
            </button>
          )}
          <button
            className="p-2 -mr-2 text-gray-700"
            aria-label="Abrir menu"
            aria-expanded={menuOpen}
            onClick={() => setMenuOpen((v) => !v)}
          >
            {menuOpen ? (
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            ) : (
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            )}
          </button>
        </div>
      </div>

      {/* Mobile menu panel */}
      {menuOpen && (
        <div className="md:hidden border-t border-gray-100 bg-white px-4 pb-4 pt-3 space-y-3">
          <div className="flex flex-col gap-3">{navLinks}</div>
          <div className="pt-3 border-t border-gray-100">{authLinks}</div>
        </div>
      )}

      {profileOpen && <ProfilePopup onClose={() => setProfileOpen(false)} />}
    </header>
  );
}
