'use client';

import { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { api } from '@/lib/api';
import { useAuth } from './AuthContext';

interface FavoritesContextType {
  ids: Set<string>;
  isFavorite: (productId: string) => boolean;
  toggle: (productId: string) => Promise<void>;
  loading: boolean;
}

const FavoritesContext = createContext<FavoritesContextType | undefined>(undefined);

export function FavoritesProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const [ids, setIds] = useState<Set<string>>(new Set());
  const [loading, setLoading] = useState(false);

  const load = useCallback(async () => {
    if (!user) {
      setIds(new Set());
      return;
    }
    try {
      const res = await api('/favorites/ids');
      if (res.ok) {
        const data = await res.json();
        setIds(new Set(data || []));
      }
    } catch {
      // silently ignore — favorites are a non-critical enhancement
    }
  }, [user]);

  useEffect(() => {
    load();
  }, [load]);

  const toggle = useCallback(
    async (productId: string) => {
      if (!user) return;
      const wasFavorite = ids.has(productId);

      // Optimistic update
      setIds((prev) => {
        const next = new Set(prev);
        if (wasFavorite) next.delete(productId);
        else next.add(productId);
        return next;
      });

      setLoading(true);
      try {
        await api(`/favorites/${productId}`, { method: wasFavorite ? 'DELETE' : 'POST' });
      } catch {
        // revert on failure
        setIds((prev) => {
          const next = new Set(prev);
          if (wasFavorite) next.add(productId);
          else next.delete(productId);
          return next;
        });
      } finally {
        setLoading(false);
      }
    },
    [ids, user]
  );

  const isFavorite = useCallback((productId: string) => ids.has(productId), [ids]);

  return (
    <FavoritesContext.Provider value={{ ids, isFavorite, toggle, loading }}>
      {children}
    </FavoritesContext.Provider>
  );
}

export function useFavorites() {
  const ctx = useContext(FavoritesContext);
  if (!ctx) throw new Error('useFavorites must be used within FavoritesProvider');
  return ctx;
}
