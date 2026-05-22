"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { productsApi, categoriesApi } from "@/lib/api/products";
import { ProductCard } from "@/components/product/ProductCard";
import { ProductGridSkeleton } from "@/components/ui/Skeleton";
import type { Product, Category } from "@/lib/types";

function FeaturedProducts() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    productsApi.getFeatured(12).then(setProducts).catch(() => {}).finally(() => setLoading(false));
  }, []);

  if (loading) return <ProductGridSkeleton count={12} />;

  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3">
      {products.map((p) => <ProductCard key={p.id} product={p} priority />)}
    </div>
  );
}

function DealProducts() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    productsApi.getDeals(20).then(setProducts).catch(() => {}).finally(() => setLoading(false));
  }, []);

  if (loading) return <ProductGridSkeleton count={10} />;

  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
      {products.map((p) => <ProductCard key={p.id} product={p} />)}
    </div>
  );
}

function TopCategories() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    categoriesApi.getTree().then(setCategories).catch(() => {}).finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div className="grid grid-cols-4 gap-4">{[...Array(8)].map((_, i) => <div key={i} className="skeleton h-24 rounded-lg" />)}</div>;
  }

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
        <TopCategories />
      </section>
      <section>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold flex items-center gap-2"><span className="text-[#ee4d2d]">⚡</span> Flash Deals</h2>
          <Link href="/products?sort=deals" className="text-sm text-[#ee4d2d] hover:underline">See All →</Link>
        </div>
        <DealProducts />
      </section>
      <section>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold">Featured Products</h2>
          <Link href="/products" className="text-sm text-[#ee4d2d] hover:underline">See All →</Link>
        </div>
        <FeaturedProducts />
      </section>
    </div>
  );
}
