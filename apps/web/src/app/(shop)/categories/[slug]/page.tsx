"use client";
import { useState, useEffect, useRef } from "react";
import { productsApi, categoriesApi } from "@/lib/api/products";
import { ProductGrid } from "@/components/product/ProductGrid";
import { Pagination } from "@/components/ui/Pagination";
import { SearchFiltersBar } from "@/components/search/SearchFilters";
import type { Product, Category, SearchFilters } from "@/lib/types";

export default function CategoryPage({ params }: { params: { slug: string } }) {
  const [category, setCategory] = useState<Category | null>(null);
  const [products, setProducts] = useState<Product[] | null>(null);
  const [total, setTotal] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<SearchFilters>({ category_id: "", page: 1, page_size: 24, sort_by: "relevance" });
  const abortRef = useRef<AbortController | null>(null);

  useEffect(() => {
    const controller = new AbortController();
    abortRef.current = controller;
    categoriesApi.getBySlug(params.slug, controller.signal)
      .then((cat) => {
        if (!controller.signal.aborted) {
          setCategory(cat);
          setFilters((f) => ({ ...f, category_id: cat.id || "" }));
        }
      })
      .catch((e: Error) => {
        if (!controller.signal.aborted) setError(e.message);
      });
    return () => controller.abort();
  }, [params.slug]);

  useEffect(() => {
    if (!filters.category_id) return;
    if (abortRef.current) abortRef.current.abort();
    const controller = new AbortController();
    abortRef.current = controller;
    setIsLoading(true);
    setError(null);

    productsApi.search(filters, controller.signal)
      .then((res) => {
        if (!controller.signal.aborted) {
          setProducts(res.products);
          setTotal(res.total);
          setIsLoading(false);
        }
      })
      .catch((e: Error) => {
        if (!controller.signal.aborted) {
          setError(e.message);
          setIsLoading(false);
        }
      });

    return () => controller.abort();
  }, [filters]);

  const totalPages = Math.max(1, Math.ceil(total / Math.max(1, filters.page_size || 24)));

  return (
    <div className="container py-6">
      <div className="mb-6">
        <h1 className="text-2xl font-bold">{category?.name || "Category"}</h1>
        {error && <p className="text-red-500 text-sm mt-2">{error}</p>}
      </div>
      <SearchFiltersBar filters={filters} onChange={(f) => setFilters({ ...f, page: 1 })} resultCount={total} />
      <ProductGrid products={products} isLoading={isLoading} />
      {totalPages > 1 && (
        <Pagination currentPage={filters.page || 1} totalPages={totalPages} onPageChange={(page) => setFilters((f) => ({ ...f, page }))} />
      )}
    </div>
  );
}
