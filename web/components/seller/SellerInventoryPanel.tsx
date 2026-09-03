'use client';

import { useEffect, useState, useCallback } from 'react';
import { api } from '@/lib/api';

interface StockItem {
  product_id: string;
  product_name: string;
  quantity: number;
  reserved: number;
}

export function SellerInventoryPanel() {
  const [items, setItems] = useState<StockItem[] | null>(null);
  const [edits, setEdits] = useState<Record<string, string>>({});
  const [saving, setSaving] = useState<string | null>(null);

  const load = useCallback(async () => {
    const res = await api('/inventory');
    if (res.ok) setItems((await res.json()) || []);
    else setItems([]);
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  async function save(productId: string) {
    const raw = edits[productId];
    if (raw === undefined) return;
    const quantity = parseInt(raw, 10);
    if (isNaN(quantity) || quantity < 0) return;
    setSaving(productId);
    try {
      await api(`/inventory/${productId}`, { method: 'PUT', body: JSON.stringify({ quantity }) });
      setEdits((e) => {
        const next = { ...e };
        delete next[productId];
        return next;
      });
      load();
    } finally {
      setSaving(null);
    }
  }

  if (items === null) return <p className="text-gray-500 text-sm">Carregando...</p>;
  if (items.length === 0) return <p className="text-gray-500 text-sm">Nenhum produto cadastrado ainda.</p>;

  return (
    <div className="space-y-2">
      {items.map((item) => {
        const editing = edits[item.product_id] !== undefined;
        const low = item.quantity - item.reserved <= 5;
        return (
          <div key={item.product_id} className="flex items-center justify-between bg-white p-3 rounded-lg border border-gray-100 gap-3">
            <div className="min-w-0">
              <p className="text-sm font-medium truncate">{item.product_name}</p>
              <p className={`text-xs ${low ? 'text-amber-600' : 'text-gray-400'}`}>
                {item.quantity} em estoque · {item.reserved} reservado{item.reserved !== 1 ? 's' : ''}
                {low && ' · estoque baixo'}
              </p>
            </div>
            <div className="flex items-center gap-2 shrink-0">
              <input
                type="number"
                min={0}
                value={editing ? edits[item.product_id] : item.quantity}
                onChange={(e) => setEdits((prev) => ({ ...prev, [item.product_id]: e.target.value }))}
                className="w-16 border border-gray-300 rounded-md px-2 py-1 text-sm text-right"
              />
              {editing && (
                <button
                  onClick={() => save(item.product_id)}
                  disabled={saving === item.product_id}
                  className="text-xs font-bold brand-text hover:underline disabled:opacity-50"
                >
                  Salvar
                </button>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}
