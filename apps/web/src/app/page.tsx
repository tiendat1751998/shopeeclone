import Link from "next/link";
import { productsApi, categoriesApi } from "@/lib/api/products";
import { ProductCard } from "@/components/product/ProductCard";
import { ProductGridSkeleton } from "@/components/ui/Skeleton";
import { Suspense } from "react";
import type { Product, Category } from "@/lib/types";

async function FeaturedProducts() {
  try {
    const products: Product[] = await productsApi.getFeatured(12);
    return (
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3">
        {products.map((p) => <ProductCard key={p.id} product={p} priority />)}
      </div>
    );
  } catch {
    return <ProductGridSkeleton count={12} />;
  }
}

async function DealProducts() {
  try {
    const products: Product[] = await productsApi.getDeals(20);
    return (
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
        {products.map((p) => <ProductCard key={p.id} product={p} />)}
      </div>
    );
  } catch {
    return <ProductGridSkeleton count={10} />;
  }
}

async function TopCategories() {
  try {
    const categories: Category[] = await categoriesApi.getTree();
    return (
      <div className="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-6 lg:grid-cols-8 gap-4">
        {categories.slice(0, 8).map((cat) => (
          <Link key={cat.id} href={`/categories/${cat.slug}`} className="flex flex-col items-center gap-2 group">
            <div className="w-16 h-16 rounded-full bg-[#fff0ed] flex items-center justify-center group-hover:bg-[#ee4d2d] transition-colors overflow-hidden">
              {cat.image_url ? (
                <img src={cat.image_url} alt={cat.name} className="w-full h-full object-cover" />
              ) : (
                <span className="text-2xl font-bold text-[#ee4d2d] group-hover:text-white transition-colors">{cat.name.charAt(0)}</span>
              )}
            </div>
            <span className="text-xs text-center text-[#222] group-hover:text-[#ee4d2d] transition-colors line-clamp-2">{cat.name}</span>
          </Link>
        ))}
      </div>
    );
  } catch {
    return null;
  }
}

export default function HomePage() {
  return (
    <div className="container py-6 space-y-8">
      <div className="card overflow-hidden">
        <div className="bg-gradient-to-r from-[#ee4d2d] to-[#f5a623] p-8 md:p-12 text-white">
          <h2 className="text-2xl md:text-4xl font-bold mb-2">Welcome to Shopee</h2>
          <p className="text-white/80 mb-4 text-sm md:text-base">Discover millions of products from trusted sellers</p>
          <Link href="/products" className="inline-block bg-white text-[#ee4d2d] px-6 py-2.5 rounded font-medium text-sm hover:bg-gray-100 transition-colors">Shop Now</Link>
        </div>
      </div>
      <section>
        <h2 className="text-lg font-semibold mb-4">Categories</h2>
        <Suspense fallback={<div className="grid grid-cols-4 gap-4">{[...Array(8)].map((_, i) => <div key={i} className="skeleton h-24 rounded-lg" />)}</div>}>
          <TopCategories />
        </Suspense>
      </section>
      <section>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold flex items-center gap-2"><span className="text-[#ee4d2d]">⚡</span> Flash Deals</h2>
          <Link href="/products?sort=deals" className="text-sm text-[#ee4d2d] hover:underline">See All →</Link>
        </div>
        <Suspense fallback={<ProductGridSkeleton count={10} />}><DealProducts /></Suspense>
      </section>
      <section>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold">Featured Products</h2>
          <Link href="/products" className="text-sm text-[#ee4d2d] hover:underline">See All →</Link>
        </div>
        <Suspense fallback={<ProductGridSkeleton count={12} />}><FeaturedProducts /></Suspense>
      </section>
    </div>
  );
}
