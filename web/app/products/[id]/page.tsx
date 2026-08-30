'use client';

import { useEffect, useState, useCallback } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/components/AuthContext';
import { useCart } from '@/components/CartContext';
import { api } from '@/lib/api';

interface Product {
  id: string;
  name: string;
  slug: string;
  description: string;
  price: number;
  state: string;
}

interface Review {
  id: string;
  user_name: string;
  rating: number;
  comment: string;
  created_at: string;
}

function Stars({ value }: { value: number }) {
  return (
    <span className="text-amber-500">
      {'★'.repeat(Math.round(value))}
      <span className="text-gray-300">{'★'.repeat(5 - Math.round(value))}</span>
    </span>
  );
}

export default function ProductDetailPage() {
  const params = useParams<{ id: string }>();
  const { user } = useAuth();
  const { add } = useCart();

  const [product, setProduct] = useState<Product | null>(null);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [average, setAverage] = useState(0);
  const [count, setCount] = useState(0);
  const [rating, setRating] = useState(5);
  const [comment, setComment] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [reviewError, setReviewError] = useState<string | null>(null);

  const loadReviews = useCallback(async () => {
    const res = await api(`/products/${params.id}/reviews`);
    if (res.ok) {
      const data = await res.json();
      setReviews(data.data || []);
      setAverage(data.average || 0);
      setCount(data.count || 0);
    }
  }, [params.id]);

  useEffect(() => {
    api(`/products/${params.id}`).then(async (res) => {
      if (res.ok) setProduct(await res.json());
    });
    loadReviews();
  }, [params.id, loadReviews]);

  async function submitReview(e: React.FormEvent) {
    e.preventDefault();
    setSubmitting(true);
    setReviewError(null);
    try {
      const res = await api(`/products/${params.id}/reviews`, {
        method: 'POST',
        body: JSON.stringify({ rating, comment }),
      });
      if (!res.ok) {
        const data = await res.json();
        setReviewError(data.message || 'Não foi possível enviar sua avaliação.');
        return;
      }
      setComment('');
      setRating(5);
      loadReviews();
    } catch {
      setReviewError('Falha de conexão');
    } finally {
      setSubmitting(false);
    }
  }

  if (!product) {
    return <p className="text-gray-500">Carregando...</p>;
  }

  return (
    <div className="max-w-3xl mx-auto space-y-10">
      <div>
        <Link href="/products" className="text-sm text-indigo-600 hover:underline">← Voltar aos produtos</Link>
        <div className="bg-white rounded-xl border border-gray-100 p-8 mt-4">
          <h1 className="text-2xl font-bold mb-2">{product.name}</h1>
          <div className="flex items-center gap-2 mb-4">
            <Stars value={average} />
            <span className="text-sm text-gray-500">
              {average.toFixed(1)} ({count} avaliaç{count === 1 ? 'ão' : 'ões'})
            </span>
          </div>
          <p className="text-gray-600 mb-6">{product.description}</p>
          <div className="flex items-center justify-between">
            <span className="text-2xl font-bold text-indigo-600">R$ {(product.price / 100).toFixed(2)}</span>
            <button
              onClick={() => add({ id: product.id, name: product.name, price: product.price, quantity: 1 })}
              className="px-4 py-2 bg-indigo-600 text-white rounded-md font-medium hover:bg-indigo-700 transition"
            >
              Adicionar ao carrinho
            </button>
          </div>
        </div>
      </div>

      <div>
        <h2 className="text-lg font-bold mb-4">Comentários</h2>
        {reviews.length === 0 ? (
          <p className="text-gray-500 text-sm">Ainda não há comentários para este produto.</p>
        ) : (
          <div className="space-y-4 mb-6">
            {reviews.map((r) => (
              <div key={r.id} className="bg-white p-4 rounded-lg border border-gray-100">
                <div className="flex items-center justify-between mb-1">
                  <span className="font-medium text-sm">{r.user_name}</span>
                  <Stars value={r.rating} />
                </div>
                <p className="text-sm text-gray-600">{r.comment}</p>
              </div>
            ))}
          </div>
        )}

        {user ? (
          <form onSubmit={submitReview} className="bg-white p-4 rounded-lg border border-gray-100 space-y-3">
            <p className="text-sm font-medium">Deixe seu comentário</p>
            <div className="flex gap-1">
              {[1, 2, 3, 4, 5].map((n) => (
                <button
                  type="button"
                  key={n}
                  onClick={() => setRating(n)}
                  className={n <= rating ? 'text-amber-500' : 'text-gray-300'}
                >
                  ★
                </button>
              ))}
            </div>
            <textarea
              required
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              placeholder="Conte como foi sua experiência com o produto..."
              className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              rows={3}
            />
            {reviewError && <div className="text-red-600 text-sm">{reviewError}</div>}
            <button
              type="submit"
              disabled={submitting}
              className="px-4 py-2 bg-indigo-600 text-white rounded-md text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 transition"
            >
              {submitting ? 'Enviando...' : 'Enviar comentário'}
            </button>
          </form>
        ) : (
          <p className="text-sm text-gray-500">
            <Link href="/login" className="text-indigo-600 hover:underline">Entre na sua conta</Link> para deixar um comentário.
          </p>
        )}
      </div>
    </div>
  );
}
