import type { Metadata } from 'next';
import './globals.css';
import { CartProvider } from '@/components/CartContext';
import { AuthProvider } from '@/components/AuthContext';
import { FavoritesProvider } from '@/components/FavoritesContext';
import { LoyaltyProvider } from '@/components/LoyaltyContext';
import { ProductModalProvider } from '@/components/ProductModalContext';
import { ProductModal } from '@/components/ProductModal';
import { ThemeProvider } from '@/components/ThemeProvider';
import { LanguageProvider } from '@/components/LanguageContext';
import { Header } from '@/components/Header';
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
        <LanguageProvider>
          <ThemeProvider>
            <AuthProvider>
              <FavoritesProvider>
                <LoyaltyProvider>
                  <CartProvider>
                    <ProductModalProvider>
                      <Header />
                      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                        {children}
                      </main>
                      <ProductModal />
                    </ProductModalProvider>
                  </CartProvider>
                </LoyaltyProvider>
              </FavoritesProvider>
            </AuthProvider>
          </ThemeProvider>
        </LanguageProvider>
      </body>
    </html>
  );
}
