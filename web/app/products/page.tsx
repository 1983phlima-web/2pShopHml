import { ProductCard } from '@/components/ProductCard';
import { api } from '@/lib/api';

export const dynamic = 'force-dynamic';

interface Product {
  id: string;
  name: string;
  slug: string;
  description: string;
  price: number;
  state: string;
}

async function getProducts(): Promise<Product[]> {
  try {
    const res = await api('/products?limit=20');
    if (!res.ok) return [];
    const data = await res.json();
    return data.data || [];
  } catch {
    return [];
  }
}

export default async function ProductsPage() {
  const products = await getProducts();

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Produtos</h1>
      {products.length === 0 ? (
        <div className="text-center py-20 text-gray-500">
          Nenhum produto disponível no momento.
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {products.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      )}
    </div>
  );
}
