import type { Metadata } from 'next';
import './globals.css';
import { CartProvider } from '@/components/CartContext';
import { OtelInit } from '@/lib/otel';
import { WebVitals } from '@/components/WebVitals';

export const metadata: Metadata = {
  title: '2pShop - Marketplace',
  description: 'Plataforma de marketplace multi-tenant',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="pt-BR">
      <body className="bg-gray-50 text-gray-900 min-h-screen">
        <OtelInit />
        <WebVitals />
        <CartProvider>
          <header className="bg-white border-b border-gray-200 sticky top-0 z-50">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
              <a href="/" className="text-xl font-bold text-indigo-600">2pShop</a>
              <nav className="flex gap-6">
                <a href="/products" className="text-sm font-medium hover:text-indigo-600">Produtos</a>
                <a href="/checkout" className="text-sm font-medium hover:text-indigo-600">Carrinho</a>
              </nav>
            </div>
          </header>
          <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            {children}
          </main>
        </CartProvider>
      </body>
    </html>
  );
}
