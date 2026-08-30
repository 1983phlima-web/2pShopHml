'use client';

import Link from 'next/link';
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
      <Link href={`/products/${product.id}`}>
        <div className="h-40 bg-gray-100 flex items-center justify-center text-gray-400 text-4xl">
          📦
        </div>
      </Link>
      <div className="p-4">
        <Link href={`/products/${product.id}`}>
          <h3 className="font-semibold text-gray-900 truncate hover:text-indigo-600">{product.name}</h3>
        </Link>
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
