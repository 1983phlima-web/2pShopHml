'use client';

import { useState } from 'react';
import Image from 'next/image';

const WHY_CARDS = [
  {
    bg: 'bg-orange-50',
    iconBg: 'bg-orange-100 text-orange-600',
    icon: '🛡️',
    title: 'Compra protegida',
    text: 'Pagamento seguro, suporte próximo e acompanhamento até a entrega.',
  },
  {
    bg: 'bg-blue-50',
    iconBg: 'bg-blue-100 text-blue-600',
    icon: '🏪',
    title: 'Lojas com história',
    text: 'Descubra pequenos negócios e marcas que fazem diferente.',
  },
  {
    bg: 'bg-emerald-50',
    iconBg: 'bg-emerald-100 text-emerald-600',
    icon: '💬',
    title: 'Atendimento humano',
    text: 'Quando precisar, fale com alguém que conhece o seu pedido.',
  },
];

const FAQ = [
  {
    q: 'Cliente: como comprar?',
    a: 'Busque ou filtre um produto na vitrine, adicione ao carrinho, revise o pedido e confirme o checkout.',
  },
  {
    q: 'Cliente: como acompanhar?',
    a: 'Abra "Minhas Compras" para acompanhar as etapas do pedido (confirmado, enviado, entregue).',
  },
  {
    q: 'Cliente: como solicitar troca?',
    a: 'Em "Minhas Compras", itens já entregues têm a opção "trocar" — acompanhe o status em "Minhas Trocas".',
  },
  {
    q: 'Vendedor: como publicar um produto?',
    a: 'Acesse o Painel do Vendedor, cadastre o produto na aba Produtos e publique-o para aparecer na vitrine.',
  },
  {
    q: 'Vendedor: como gerenciar pedidos e estoque?',
    a: 'No Painel do Vendedor, as abas Pedidos e Estoque permitem avançar o status de entrega e ajustar quantidades.',
  },
  {
    q: 'Vendedor: como responder perguntas?',
    a: 'Perguntas feitas nos seus produtos aparecem na página do produto — responda diretamente por lá.',
  },
];

function FaqItem({ q, a }: { q: string; a: string }) {
  const [open, setOpen] = useState(false);
  return (
    <div className="rounded-2xl bg-white/10 p-4">
      <button onClick={() => setOpen((v) => !v)} className="w-full text-left flex items-center justify-between gap-2">
        <span className="text-sm font-extrabold">{q}</span>
        <span className="text-white/50 text-xs shrink-0">{open ? '−' : '+'}</span>
      </button>
      {open && <p className="mt-3 text-xs leading-relaxed text-white/70">{a}</p>}
    </div>
  );
}

export function Footer() {
  return (
    <div className="mt-16 -mx-4 sm:-mx-6 lg:-mx-8">
      {/* Por que 2pShop */}
      <section className="border-y border-gray-100 bg-orange-50/30 py-14 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto">
          <p className="text-[11px] font-bold uppercase tracking-[.17em] brand-text">Por que 2pShop</p>
          <h2 className="mt-1 text-2xl font-extrabold mb-6">Mais que um marketplace. Um jeito melhor de comprar.</h2>
          <div className="grid gap-4 md:grid-cols-3">
            {WHY_CARDS.map((card) => (
              <div key={card.title} className={`rounded-2xl ${card.bg} p-6`}>
                <span className={`grid h-11 w-11 place-items-center rounded-xl ${card.iconBg} text-xl`}>
                  {card.icon}
                </span>
                <h3 className="mt-5 font-extrabold text-gray-900">{card.title}</h3>
                <p className="mt-2 text-sm leading-relaxed text-gray-500">{card.text}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Central de Ajuda */}
      <section className="py-14 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto">
          <div className="rounded-3xl bg-gray-900 p-6 sm:p-10 text-white">
            <p className="text-[11px] font-bold uppercase tracking-[.17em] text-amber-300">Central de ajuda</p>
            <h2 className="mt-2 text-2xl sm:text-3xl font-extrabold">Como usar a 2pShop</h2>
            <div className="mt-6 grid gap-3 md:grid-cols-2">
              {FAQ.map((item) => (
                <FaqItem key={item.q} {...item} />
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Footer bar */}
      <footer className="border-t border-gray-100 bg-white py-8 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto flex flex-col sm:flex-row justify-between items-start sm:items-center gap-5">
          <div className="flex items-center gap-3">
            <Image src="/logo.png" alt="2pShop" width={546} height={633} className="h-9 w-auto" />
            <div>
              <p className="font-extrabold brand-text">2pShop</p>
              <p className="text-xs text-gray-400">Marketplace que aproxima quem compra de quem faz.</p>
            </div>
          </div>
          <div className="flex flex-wrap items-center gap-4 text-xs font-semibold text-gray-500">
            <span>Como funciona</span>
            <span>Ajuda</span>
            <span>Privacidade</span>
          </div>
        </div>
      </footer>
    </div>
  );
}
