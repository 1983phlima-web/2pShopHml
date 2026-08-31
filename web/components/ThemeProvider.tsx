'use client';

import { useEffect } from 'react';
import { api } from '@/lib/api';
import { applyPalette, DEFAULT_PALETTE_KEY } from '@/lib/theme';

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    api('/settings/theme')
      .then(async (res) => {
        if (!res.ok) return applyPalette(DEFAULT_PALETTE_KEY);
        const data = await res.json();
        applyPalette(data.palette || DEFAULT_PALETTE_KEY);
      })
      .catch(() => applyPalette(DEFAULT_PALETTE_KEY));
  }, []);

  return <>{children}</>;
}
