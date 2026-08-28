'use client';

import { useCart } from './CartContext';

interface Product {
  id: string;
  name: string;
  slug: string;
  description: string;
  price: number;
  state: string;
}

export function ProductCard({ product }: { product: Product }) {
  const { add } = useCart();

  return (
    <div className="bg-white rounded-xl border border-gray-100 overflow-hidden hover:shadow-md transition">
      <div className="h-40 bg-gray-100 flex items-center justify-center text-gray-400">
        📦
      </div>
      <div className="p-4">
        <h3 className="font-semibold text-gray-900 truncate">{product.name}</h3>
        <p className="text-sm text-gray-500 line-clamp-2 mt-1">{product.description}</p>
        <div className="mt-4 flex items-center justify-between">
          <span className="text-lg font-bold text-indigo-600">
            R$ {(product.price / 100).toFixed(2)}
          </span>
          <button
            onClick={() => add({ id: product.id, name: product.name, price: product.price, quantity: 1 })}
            className="px-3 py-1.5 bg-indigo-50 text-indigo-700 text-sm font-medium rounded-md hover:bg-indigo-100 transition"
          >
            Adicionar
          </button>
        </div>
      </div>
    </div>
  );
}
