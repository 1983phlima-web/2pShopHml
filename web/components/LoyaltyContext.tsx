'use client';

import { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api';
import { useAuth } from './AuthContext';

export interface LoyaltyProfile {
  xp: number;
  coins: number;
  badges: { key: string; label: string; earned_at: string }[];
  period_15day_spend: number;
  period_15day_target: number;
  period_month_spend: number;
  period_month_target: number;
}

interface LoyaltyContextType {
  loyalty: LoyaltyProfile | null;
  refresh: () => void;
}

const LoyaltyContext = createContext<LoyaltyContextType | undefined>(undefined);

export function LoyaltyProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const [loyalty, setLoyalty] = useState<LoyaltyProfile | null>(null);

  const refresh = useCallback(() => {
    if (!user || user.role !== 'BUYER') {
      setLoyalty(null);
      return;
    }
    api('/loyalty/profile').then(async (res) => {
      if (res.ok) setLoyalty(await res.json());
    });
  }, [user]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  return <LoyaltyContext.Provider value={{ loyalty, refresh }}>{children}</LoyaltyContext.Provider>;
}

export function useLoyalty() {
  const ctx = useContext(LoyaltyContext);
  if (!ctx) throw new Error('useLoyalty must be used within LoyaltyProvider');
  return ctx;
}
