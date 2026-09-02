'use client';

import { useEffect, useState, useCallback } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/components/AuthContext';
import { useCart } from '@/components/CartContext';
import { useFavorites } from '@/components/FavoritesContext';
import { ImageCarousel } from '@/components/ImageCarousel';
import { api } from '@/lib/api';

interface Product {
  id: string;
  name: string;
  slug: string;
  description: string;
  price: number;
  state: string;
  attributes?: {
    image?: string;
    images?: string[];
    badge?: string;
    brand?: string;
    compare_at?: number;
    sold?: number;
  };
}

interface Question {
  id: string;
  user_name: string;
  question: string;
  answer?: string;
  created_at: string;
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
  const router = useRouter();
  const { user } = useAuth();
  const { add } = useCart();
  const { isFavorite, toggle } = useFavorites();

  const [product, setProduct] = useState<Product | null>(null);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [average, setAverage] = useState(0);
  const [count, setCount] = useState(0);
  const [rating, setRating] = useState(5);
  const [comment, setComment] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [reviewError, setReviewError] = useState<string | null>(null);

  const [questions, setQuestions] = useState<Question[]>([]);
  const [newQuestion, setNewQuestion] = useState('');
  const [submittingQuestion, setSubmittingQuestion] = useState(false);
  const [questionError, setQuestionError] = useState<string | null>(null);

  const loadQuestions = useCallback(async () => {
    const res = await api(`/products/${params.id}/questions`);
    if (res.ok) setQuestions((await res.json()) || []);
  }, [params.id]);

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
    loadQuestions();
  }, [params.id, loadReviews, loadQuestions]);

  async function submitQuestion(e: React.FormEvent) {
    e.preventDefault();
    setSubmittingQuestion(true);
    setQuestionError(null);
    try {
      const res = await api(`/products/${params.id}/questions`, {
        method: 'POST',
        body: JSON.stringify({ question: newQuestion }),
      });
      if (!res.ok) {
        const data = await res.json();
        setQuestionError(data.message || 'Não foi possível enviar sua pergunta.');
        return;
      }
      setNewQuestion('');
      loadQuestions();
    } catch {
      setQuestionError('Falha de conexão');
    } finally {
      setSubmittingQuestion(false);
    }
  }

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
        <div className="bg-white rounded-xl border border-gray-100 overflow-hidden mt-4 grid md:grid-cols-2">
          {(product.attributes?.images?.length || product.attributes?.image) && (
            <div className="relative h-64 md:h-full bg-gray-100">
              <ImageCarousel
                images={product.attributes?.images?.length ? product.attributes.images : [product.attributes?.image || '']}
                alt={product.name}
              />
              {product.attributes?.badge && (
                <span className="absolute top-3 left-3 bg-indigo-600 text-white text-[10px] font-bold px-2 py-1 rounded-full uppercase tracking-wide z-10">
                  {product.attributes.badge}
                </span>
              )}
            </div>
          )}
          <div className="p-8">
            {product.attributes?.brand && (
              <p className="text-xs font-bold uppercase tracking-wide text-indigo-500 mb-1">{product.attributes.brand}</p>
            )}
            <h1 className="text-2xl font-bold mb-2">{product.name}</h1>
            <div className="flex items-center gap-2 mb-4">
              <Stars value={average} />
              <span className="text-sm text-gray-500">
                {average.toFixed(1)} ({count} avaliaç{count === 1 ? 'ão' : 'ões'})
              </span>
            </div>
            <p className="text-gray-600 mb-6">{product.description}</p>
            <div className="flex items-center gap-3 mb-6">
              {product.attributes?.compare_at && product.attributes.compare_at > product.price && (
                <span className="text-sm text-gray-400 line-through">
                  R$ {(product.attributes.compare_at / 100).toFixed(2)}
                </span>
              )}
              <span className="text-2xl font-bold text-indigo-600">R$ {(product.price / 100).toFixed(2)}</span>
            </div>
            <div className="flex items-center gap-3">
              <button
                onClick={() => add({ id: product.id, name: product.name, price: product.price, quantity: 1 })}
                className="flex-1 sm:flex-none px-6 py-2.5 bg-indigo-600 text-white rounded-md font-medium hover:bg-indigo-700 transition"
              >
                Adicionar ao carrinho
              </button>
              <button
                onClick={() => (user ? toggle(product.id) : router.push('/login'))}
                aria-label={isFavorite(product.id) ? 'Remover dos favoritos' : 'Adicionar aos favoritos'}
                aria-pressed={isFavorite(product.id)}
                className="h-11 w-11 shrink-0 rounded-md border border-gray-200 flex items-center justify-center hover:border-rose-300 transition"
              >
                <svg
                  viewBox="0 0 24 24"
                  className="h-5 w-5"
                  fill={isFavorite(product.id) ? '#e11d48' : 'none'}
                  stroke={isFavorite(product.id) ? '#e11d48' : 'currentColor'}
                  strokeWidth={2}
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M12 21s-6.716-4.35-9.428-8.06C.79 10.31 1.2 6.6 4.2 4.9c2.3-1.3 5-0.7 6.5 1.3l1.3 1.7 1.3-1.7c1.5-2 4.2-2.6 6.5-1.3 3 1.7 3.41 5.41 1.63 8.04C18.716 16.65 12 21 12 21z"
                  />
                </svg>
              </button>
            </div>
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

      <div>
        <h2 className="text-lg font-bold mb-4">Perguntas</h2>
        {questions.length === 0 ? (
          <p className="text-gray-500 text-sm mb-6">Ainda não há perguntas sobre este produto. Seja o primeiro a perguntar!</p>
        ) : (
          <div className="space-y-4 mb-6">
            {questions.map((q) => (
              <div key={q.id} className="bg-white p-4 rounded-lg border border-gray-100">
                <p className="text-sm font-medium mb-1">
                  <span className="text-gray-400">P:</span> {q.question}
                </p>
                {q.answer ? (
                  <p className="text-sm text-gray-600 pl-4 border-l-2 border-indigo-100">
                    <span className="text-gray-400">R:</span> {q.answer}
                  </p>
                ) : (
                  <p className="text-xs text-gray-400 italic">Aguardando resposta do vendedor</p>
                )}
              </div>
            ))}
          </div>
        )}

        {user ? (
          <form onSubmit={submitQuestion} className="bg-white p-4 rounded-lg border border-gray-100 space-y-3">
            <p className="text-sm font-medium">Faça uma pergunta sobre este produto</p>
            <textarea
              required
              value={newQuestion}
              onChange={(e) => setNewQuestion(e.target.value)}
              placeholder="Ex.: Este produto tem garantia? Qual o material?"
              className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              rows={2}
            />
            {questionError && <div className="text-red-600 text-sm">{questionError}</div>}
            <button
              type="submit"
              disabled={submittingQuestion}
              className="px-4 py-2 bg-indigo-600 text-white rounded-md text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 transition"
            >
              {submittingQuestion ? 'Enviando...' : 'Perguntar'}
            </button>
          </form>
        ) : (
          <p className="text-sm text-gray-500">
            <Link href="/login" className="text-indigo-600 hover:underline">Entre na sua conta</Link> para fazer uma pergunta.
          </p>
        )}
      </div>
    </div>
  );
}
