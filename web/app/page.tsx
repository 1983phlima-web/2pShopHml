import Link from 'next/link';

export default function HomePage() {
  return (
    <div className="space-y-12">
      <section className="text-center py-20 bg-gradient-to-br from-indigo-50 to-white rounded-2xl">
        <h1 className="text-4xl md:text-6xl font-extrabold text-gray-900 mb-6">
          Seu marketplace, <span className="text-indigo-600">sem limites</span>
        </h1>
        <p className="text-lg text-gray-600 max-w-2xl mx-auto mb-8">
          Plataforma multi-tenant orientada a eventos, com observabilidade nativa e escala horizontal.
        </p>
        <Link
          href="/products"
          className="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 transition"
        >
          Explorar produtos
        </Link>
      </section>

      <section className="grid md:grid-cols-3 gap-8">
        {[
          { title: 'Catálogo', desc: 'Gerencie produtos, variantes, categorias e SEO.' },
          { title: 'Checkout', desc: 'Saga distribuída com compensação automática.' },
          { title: 'Observabilidade', desc: 'Traces, métricas e logs com OpenTelemetry.' },
        ].map((f) => (
          <div key={f.title} className="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
            <h3 className="text-lg font-semibold mb-2">{f.title}</h3>
            <p className="text-gray-600 text-sm">{f.desc}</p>
          </div>
        ))}
      </section>
    </div>
  );
}
