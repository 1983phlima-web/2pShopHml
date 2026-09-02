'use client';

import { createContext, useContext, useState, useCallback } from 'react';

interface ProductModalContextType {
  openProductId: string | null;
  open: (id: string) => void;
  close: () => void;
}

const ProductModalContext = createContext<ProductModalContextType | undefined>(undefined);

export function ProductModalProvider({ children }: { children: React.ReactNode }) {
  const [openProductId, setOpenProductId] = useState<string | null>(null);
  const open = useCallback((id: string) => setOpenProductId(id), []);
  const close = useCallback(() => setOpenProductId(null), []);
  return (
    <ProductModalContext.Provider value={{ openProductId, open, close }}>
      {children}
    </ProductModalContext.Provider>
  );
}

export function useProductModal() {
  const ctx = useContext(ProductModalContext);
  if (!ctx) throw new Error('useProductModal must be used within ProductModalProvider');
  return ctx;
}
