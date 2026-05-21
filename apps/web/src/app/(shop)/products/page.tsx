"use client";
import { Suspense, useState, useEffect, useCallback, useRef } from "react";
import { useSearchParams } from "next/navigation";
import { productsApi } from "@/lib/api/products";
import { ProductGrid } from "@/components/product/ProductGrid";
import { SearchFiltersBar, PriceFilter } from "@/components/search/SearchFilters";
import { Pagination } from "@/components/ui/Pagination";
import { CategorySidebar } from "@/components/layout/CategorySidebar";
import { useDebounce } from "@/lib/hooks/useDebounce";
import type { SearchFilters, SearchResult } from "@/lib/types";

export default function ProductsPage() {
  return (
    <Suspense fallback={<div className="p-8 text-center">Loading...</div>}>
      <ProductsContent />
    </Suspense>
  );
}

function ProductsContent() {
  const searchParams = useSearchParams();
  const query = searchParams.get("q") || "";
  const categorySlug = searchParams.get("category") || "";
  const [filters, setFilters] = useState<SearchFilters>({
    query, category_id: categorySlug || undefined, page: 1, page_size: 24, sort_by: "relevance",
  });
  const [result, setResult] = useState<SearchResult | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const debouncedQuery = useDebounce(filters.query, 300);
  const abortRef = useRef<AbortController | null>(null);

  const search = useCallback(async () => {
    if (abortRef.current) abortRef.current.abort();
    const controller = new AbortController();
    abortRef.current = controller;
    setIsLoading(true);
    setError(null);

    try {
      const res = await productsApi.search({ ...filters, query: debouncedQuery }, controller.signal);
      if (!controller.signal.aborted) {
        setResult(res);
        setIsLoading(false);
      }
    } catch (e: unknown) {
      if (!controller.signal.aborted) {
        setError(e instanceof Error ? e.message : "Search failed");
        setIsLoading(false);
      }
    }
  }, [filters, debouncedQuery]);

  useEffect(() => {
    search();
    return () => { if (abortRef.current) abortRef.current.abort(); };
  }, [search]);

  const totalPages = result ? Math.max(1, Math.ceil(result.total / (result.page_size || 24))) : 1;

  return (
    <div className="container py-6">
      <div className="flex gap-6">
        <div className="hidden lg:block w-56 flex-shrink-0">
          <CategorySidebar />
          <div className="card p-4 mt-4">
            <PriceFilter minPrice={filters.min_price} maxPrice={filters.max_price}
              onChange={(min, max) => setFilters((f) => ({ ...f, min_price: min, max_price: max, page: 1 }))} />
          </div>
        </div>
        <div className="flex-1 min-w-0">
          {error && <div className="bg-red-50 text-red-600 text-sm p-3 rounded mb-4">{error}</div>}
          <SearchFiltersBar filters={filters} onChange={(f) => setFilters({ ...f, page: 1 })} resultCount={result?.total} />
          <ProductGrid products={result?.products ?? null} isLoading={isLoading} />
          {result && totalPages > 1 && (
            <Pagination currentPage={result.page} totalPages={totalPages}
              onPageChange={(page) => setFilters((f) => ({ ...f, page }))} />
          )}
        </div>
      </div>
    </div>
  );
}
