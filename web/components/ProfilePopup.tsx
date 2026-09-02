'use client';

import { useEffect, useState, useRef } from 'react';
import { useAuth } from './AuthContext';
import { Avatar } from './Avatar';
import { AVATAR_PRESETS } from '@/lib/avatars';
import { ROLE_LABELS } from '@/lib/auth';
import { api } from '@/lib/api';

interface LoyaltyProfile {
  xp: number;
  coins: number;
  badges: { key: string; label: string; earned_at: string }[];
  period_15day_spend: number;
  period_15day_target: number;
  period_month_spend: number;
  period_month_target: number;
}

const MAX_AVATAR_BYTES = 500_000; // ~500KB before base64 overhead

export function ProfilePopup({ onClose }: { onClose: () => void }) {
  const { user, updateUser } = useAuth();
  const [loyalty, setLoyalty] = useState<LoyaltyProfile | null>(null);
  const [pickerOpen, setPickerOpen] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    api('/loyalty/profile').then(async (res) => {
      if (res.ok) setLoyalty(await res.json());
    });
  }, []);

  async function selectPreset(id: number) {
    const avatar = `preset:${id}`;
    updateUser({ avatar });
    setPickerOpen(false);
    await api('/auth/me/avatar', { method: 'PUT', body: JSON.stringify({ avatar }) });
  }

  function handleFile(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    setUploadError(null);
    if (file.size > MAX_AVATAR_BYTES) {
      setUploadError('Imagem muito grande — escolha uma foto de até 500KB.');
      return;
    }
    const reader = new FileReader();
    reader.onload = async () => {
      const dataUri = reader.result as string;
      updateUser({ avatar: dataUri });
      setPickerOpen(false);
      await api('/auth/me/avatar', { method: 'PUT', body: JSON.stringify({ avatar: dataUri }) });
    };
    reader.readAsDataURL(file);
  }

  if (!user) return null;

  const period15Pct = loyalty ? Math.min(100, (loyalty.period_15day_spend / loyalty.period_15day_target) * 100) : 0;
  const periodMonthPct = loyalty ? Math.min(100, (loyalty.period_month_spend / loyalty.period_month_target) * 100) : 0;

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center sm:items-center px-4 pt-20 sm:pt-4">
      <div className="absolute inset-0 bg-black/40" onClick={onClose} />
      <div className="relative bg-white w-full max-w-sm rounded-2xl shadow-xl max-h-[85vh] overflow-y-auto">
        <button onClick={onClose} className="absolute top-3 right-3 h-8 w-8 flex items-center justify-center text-gray-400 hover:text-gray-700">
          ✕
        </button>

        <div className="p-6 text-center border-b border-gray-100">
          <div className="relative inline-block">
            <Avatar avatar={user.avatar} size={80} />
            <button
              onClick={() => setPickerOpen((v) => !v)}
              className="absolute -bottom-1 -right-1 h-7 w-7 rounded-full brand-bg text-white flex items-center justify-center text-xs shadow"
              aria-label="Trocar avatar"
            >
              ✎
            </button>
          </div>
          <h2 className="text-lg font-bold mt-3">{user.name}</h2>
          <p className="text-xs text-gray-500">{ROLE_LABELS[user.role]}</p>

          {pickerOpen && (
            <div className="mt-4 text-left">
              <p className="text-xs font-semibold text-gray-500 mb-2">Escolha um avatar</p>
              <div className="grid grid-cols-6 gap-2 mb-3">
                {AVATAR_PRESETS.map((p) => (
                  <button
                    key={p.id}
                    onClick={() => selectPreset(p.id)}
                    className={`rounded-full bg-gradient-to-br ${p.gradient} h-9 w-9 flex items-center justify-center text-base hover:scale-110 transition`}
                  >
                    {p.emoji}
                  </button>
                ))}
              </div>
              <input ref={fileInputRef} type="file" accept="image/*" onChange={handleFile} className="hidden" />
              <button
                onClick={() => fileInputRef.current?.click()}
                className="text-xs font-medium brand-text hover:underline"
              >
                Ou enviar uma foto...
              </button>
              {uploadError && <p className="text-xs text-red-600 mt-1">{uploadError}</p>}
            </div>
          )}
        </div>

        <div className="p-6 space-y-5">
          <div>
            <p className="text-xs font-semibold text-gray-500 mb-2">Contato</p>
            <div className="text-sm space-y-1">
              <p className="text-gray-700">{user.email}</p>
              <p className="text-gray-400">{user.phone || 'Telefone não informado'}</p>
            </div>
          </div>

          {user.role === 'BUYER' && (
            <>
              <div>
                <div className="flex items-center justify-between text-xs mb-1">
                  <span className="font-semibold text-gray-500">XP acumulado</span>
                  <span className="font-bold brand-text">{loyalty?.xp ?? 0} XP</span>
                </div>
                <div className="flex items-center justify-between text-xs mb-3">
                  <span className="font-semibold text-gray-500">Coins</span>
                  <span className="font-bold text-amber-600">🪙 {loyalty?.coins ?? 0}</span>
                </div>
              </div>

              <div>
                <div className="flex items-center justify-between text-xs mb-1">
                  <span className="text-gray-500">Gasto no período de 15 dias</span>
                  <span className="text-gray-400">
                    R$ {((loyalty?.period_15day_spend ?? 0) / 100).toFixed(0)} / R$ {((loyalty?.period_15day_target ?? 100000) / 100).toFixed(0)}
                  </span>
                </div>
                <div className="h-2 bg-gray-100 rounded-full overflow-hidden">
                  <div className="h-full brand-bg rounded-full transition-all" style={{ width: `${period15Pct}%` }} />
                </div>
              </div>

              <div>
                <div className="flex items-center justify-between text-xs mb-1">
                  <span className="text-gray-500">Gasto no mês (rumo a VIP)</span>
                  <span className="text-gray-400">
                    R$ {((loyalty?.period_month_spend ?? 0) / 100).toFixed(0)} / R$ {((loyalty?.period_month_target ?? 300000) / 100).toFixed(0)}
                  </span>
                </div>
                <div className="h-2 bg-gray-100 rounded-full overflow-hidden">
                  <div className="h-full bg-amber-500 rounded-full transition-all" style={{ width: `${periodMonthPct}%` }} />
                </div>
              </div>

              <div>
                <p className="text-xs font-semibold text-gray-500 mb-2">Badges conquistadas</p>
                {loyalty?.badges.length ? (
                  <div className="flex flex-wrap gap-2">
                    {loyalty.badges.map((b) => (
                      <span key={b.key} className="text-xs font-medium bg-amber-50 text-amber-700 px-2.5 py-1 rounded-full">
                        🏅 {b.label}
                      </span>
                    ))}
                  </div>
                ) : (
                  <p className="text-xs text-gray-400">
                    Nenhuma badge ainda — gaste R$500+ em 15 dias para conquistar a primeira!
                  </p>
                )}
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
