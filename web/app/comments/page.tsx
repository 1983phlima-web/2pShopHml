'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/components/AuthContext';
import { api } from '@/lib/api';

interface Question {
  id: string;
  product_id: string;
  question: string;
  answer?: string;
  created_at: string;
}

interface Review {
  id: string;
  product_id: string;
  product_name: string;
  rating: number;
  comment: string;
  created_at: string;
}

function Stars({ value }: { value: number }) {
  return (
    <span className="text-amber-500 text-sm">
      {'★'.repeat(value)}
      <span className="text-gray-300">{'★'.repeat(5 - value)}</span>
    </span>
  );
}

export default function CommentsPage() {
  const { user } = useAuth();
  const [tab, setTab] = useState<'questions' | 'testimonials'>('testimonials');
  const [questions, setQuestions] = useState<Question[] | null>(null);
  const [reviews, setReviews] = useState<Review[] | null>(null);

  useEffect(() => {
    if (!user) return;
    api('/questions/mine').then(async (res) => {
      if (res.ok) setQuestions((await res.json()) || []);
      else setQuestions([]);
    });
    api('/reviews/mine').then(async (res) => {
      if (res.ok) setReviews((await res.json()) || []);
      else setReviews([]);
    });
  }, [user]);

  if (!user) {
    return (
      <div className="max-w-md mx-auto text-center py-20">
        <p className="text-gray-600 mb-4">Entre na sua conta para ver seus comentários.</p>
        <Link href="/login" className="text-indigo-600 hover:underline">Entrar</Link>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Meus Comentários</h1>

      <div className="flex gap-1 bg-gray-100 rounded-lg p-1 w-fit mb-6">
        <button
          onClick={() => setTab('testimonials')}
          className={`text-sm font-bold px-4 py-2 rounded-md transition ${tab === 'testimonials' ? 'bg-white shadow-sm brand-text' : 'text-gray-500'}`}
        >
          Testemunhos
        </button>
        <button
          onClick={() => setTab('questions')}
          className={`text-sm font-bold px-4 py-2 rounded-md transition ${tab === 'questions' ? 'bg-white shadow-sm brand-text' : 'text-gray-500'}`}
        >
          Perguntas
        </button>
      </div>

      {tab === 'testimonials' && (
        reviews === null ? (
          <p className="text-gray-500 text-sm">Carregando...</p>
        ) : reviews.length === 0 ? (
          <p className="text-gray-500 text-sm">
            Você ainda não deixou nenhum testemunho. Avalie um produto que já recebeu na página do produto.
          </p>
        ) : (
          <div className="space-y-3">
            {reviews.map((r) => (
              <div key={r.id} className="bg-white p-4 rounded-lg border border-gray-100">
                <div className="flex items-center justify-between mb-1">
                  <Link href={`/products/${r.product_id}`} className="text-sm font-medium hover:underline brand-text">
                    {r.product_name}
                  </Link>
                  <Stars value={r.rating} />
                </div>
                <p className="text-sm text-gray-600">{r.comment}</p>
                <p className="text-xs text-gray-400 mt-1">{new Date(r.created_at).toLocaleDateString('pt-BR')}</p>
              </div>
            ))}
          </div>
        )
      )}

      {tab === 'questions' && (
        questions === null ? (
          <p className="text-gray-500 text-sm">Carregando...</p>
        ) : questions.length === 0 ? (
          <p className="text-gray-500 text-sm">
            Você ainda não fez nenhuma pergunta. Pergunte algo na página de um produto anunciado.
          </p>
        ) : (
          <div className="space-y-3">
            {questions.map((q) => (
              <div key={q.id} className="bg-white p-4 rounded-lg border border-gray-100">
                <Link href={`/products/${q.product_id}`} className="text-xs brand-text hover:underline">
                  Ver produto →
                </Link>
                <p className="text-sm font-medium mt-1"><span className="text-gray-400">P:</span> {q.question}</p>
                {q.answer ? (
                  <p className="text-sm text-gray-600 mt-1"><span className="text-gray-400">R:</span> {q.answer}</p>
                ) : (
                  <p className="text-xs text-gray-400 italic mt-1">Aguardando resposta do vendedor</p>
                )}
              </div>
            ))}
          </div>
        )
      )}
    </div>
  );
}
