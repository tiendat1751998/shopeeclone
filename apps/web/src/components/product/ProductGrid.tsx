import { ProductCard } from "./ProductCard";
import { ProductGridSkeleton } from "@/components/ui/Skeleton";
import type { Product } from "@/lib/types";

interface ProductGridProps { products: Product[] | null; isLoading?: boolean; emptyMessage?: string; columns?: number; }

export function ProductGrid({ products, isLoading, emptyMessage = "No products found", columns = 6 }: ProductGridProps) {
  if (isLoading) return <ProductGridSkeleton count={columns * 2} />;
  if (!products || products.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-16 text-[#757575]">
        <svg className="w-16 h-16 mb-4 text-[#e8e8e8]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" /></svg>
        <p className="text-lg font-medium">{emptyMessage}</p>
      </div>
    );
  }

  const colClasses = {
    2: "grid-cols-2", 3: "grid-cols-2 sm:grid-cols-3", 4: "grid-cols-2 sm:grid-cols-3 md:grid-cols-4",
    5: "grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5",
    6: "grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6",
  }[columns] || "grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6";

  return (
    <div className={`grid ${colClasses} gap-3`}>
      {products.map((product, i) => (
        <ProductCard key={product.id} product={product} priority={i < 6} />
      ))}
    </div>
  );
}
