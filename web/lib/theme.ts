export interface Palette {
  key: string;
  label: string;
  primary: string;
  primaryDark: string;
  accent: string;
  dark: string;
}

// Curated palette catalog. "navyCoral" mirrors the 2pShop reference
// app's signature colors (#142033 / #ef6c4d); the others keep our
// existing indigo identity and offer harmonious alternatives.
export const PALETTES: Record<string, Palette> = {
  indigo: {
    key: 'indigo',
    label: 'Índigo (padrão)',
    primary: '#4f46e5',
    primaryDark: '#4338ca',
    accent: '#818cf8',
    dark: '#111827',
  },
  navyCoral: {
    key: 'navyCoral',
    label: 'Navy & Coral (2pShop)',
    primary: '#ef6c4d',
    primaryDark: '#d95c40',
    accent: '#ffb47a',
    dark: '#142033',
  },
  emerald: {
    key: 'emerald',
    label: 'Esmeralda',
    primary: '#059669',
    primaryDark: '#047857',
    accent: '#6ee7b7',
    dark: '#064e3b',
  },
  rose: {
    key: 'rose',
    label: 'Rosé',
    primary: '#e11d48',
    primaryDark: '#be123c',
    accent: '#fda4af',
    dark: '#1f2937',
  },
  amber: {
    key: 'amber',
    label: 'Âmbar',
    primary: '#d97706',
    primaryDark: '#b45309',
    accent: '#fcd34d',
    dark: '#292524',
  },
};

export const DEFAULT_PALETTE_KEY = 'indigo';

export function applyPalette(key: string) {
  const palette = PALETTES[key] || PALETTES[DEFAULT_PALETTE_KEY];
  const root = document.documentElement;
  root.style.setProperty('--brand-primary', palette.primary);
  root.style.setProperty('--brand-primary-dark', palette.primaryDark);
  root.style.setProperty('--brand-accent', palette.accent);
  root.style.setProperty('--brand-dark', palette.dark);
}
