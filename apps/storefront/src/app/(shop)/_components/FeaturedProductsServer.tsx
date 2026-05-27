// SERVER COMPONENT - Fetches featured products on server
import { productsApi } from "@/lib/api/products";
import { ProductCard } from "@/components/product/ProductCard";

export async function FeaturedProductsServer() {
  const products = await productsApi.getFeatured(12);
  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-2 md:gap-3">
      {products.map((p, i) => <ProductCard key={p.id} product={p} priority={i < 6} />)}
    </div>
  );
}
