"use client";
import { useState, useEffect, useCallback } from "react";
import { useSearchParams } from "next/navigation";
import { productsApi } from "@/lib/api/products";
import { ProductGrid } from "@/components/product/ProductGrid";
import { SearchFiltersBar, PriceFilter } from "@/components/search/SearchFilters";
import { Pagination } from "@/components/ui/Pagination";
import { CategorySidebar } from "@/components/layout/CategorySidebar";
import { useDebounce } from "@/lib/hooks/useDebounce";
import type { SearchFilters, SearchResult } from "@/lib/types";

export default function ProductsPage() {
  const searchParams = useSearchParams();
  const query = searchParams.get("q") || "";
  const categorySlug = searchParams.get("category") || "";
  const [filters, setFilters] = useState<SearchFilters>({ query, category_id: categorySlug, page: 1, page_size: 24, sort_by: "relevance" });
  const [result, setResult] = useState<SearchResult | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const debouncedQuery = useDebounce(filters.query, 300);

  const search = useCallback(async () => {
    setIsLoading(true);
    try { const res = await productsApi.search({ ...filters, query: debouncedQuery }); setResult(res); }
    catch { /* handle */ }
    finally { setIsLoading(false); }
  }, [filters, debouncedQuery]);

  useEffect(() => { search(); }, [search]);

  return (
    <div className="container py-6">
      <div className="flex gap-6">
        <div className="hidden lg:block w-56 flex-shrink-0">
          <CategorySidebar />
          <div className="card p-4 mt-4">
            <PriceFilter minPrice={filters.min_price} maxPrice={filters.max_price} onChange={(min, max) => setFilters((f) => ({ ...f, min_price: min, max_price: max, page: 1 }))} />
          </div>
        </div>
        <div className="flex-1 min-w-0">
          <SearchFiltersBar filters={filters} onChange={(f) => setFilters({ ...f, page: 1 })} resultCount={result?.total} />
          <ProductGrid products={result?.products ?? null} isLoading={isLoading} />
          {result && result.total_pages > 1 && <Pagination currentPage={result.page} totalPages={result.total_pages} onPageChange={(page) => setFilters((f) => ({ ...f, page }))} />}
        </div>
      </div>
    </div>
  );
}
