'use client';

import { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api';
import { LangCode, translations } from '@/lib/i18n';

interface LanguageContextType {
  lang: LangCode;
  setLang: (code: LangCode, persist?: boolean) => void;
  t: (key: string) => string;
}

const LanguageContext = createContext<LanguageContextType | undefined>(undefined);

export function LanguageProvider({ children }: { children: React.ReactNode }) {
  const [lang, setLangState] = useState<LangCode>('pt');

  useEffect(() => {
    api('/settings/language').then(async (res) => {
      if (res.ok) {
        const data = await res.json();
        if (data.code) setLangState(data.code as LangCode);
      }
    });
  }, []);

  const setLang = useCallback(async (code: LangCode, persist = false) => {
    setLangState(code);
    if (persist) {
      await api('/settings/language', { method: 'PUT', body: JSON.stringify({ code }) });
    }
  }, []);

  const t = useCallback((key: string) => translations[lang][key] || translations.pt[key] || key, [lang]);

  return <LanguageContext.Provider value={{ lang, setLang, t }}>{children}</LanguageContext.Provider>;
}

export function useLanguage() {
  const ctx = useContext(LanguageContext);
  if (!ctx) throw new Error('useLanguage must be used within LanguageProvider');
  return ctx;
}
