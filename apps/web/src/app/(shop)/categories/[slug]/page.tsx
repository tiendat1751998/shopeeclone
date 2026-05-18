"use client";
import { useState, useEffect } from "react";
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
  const [filters, setFilters] = useState<SearchFilters>({ category_id: "", page: 1, page_size: 24, sort_by: "relevance" });

  useEffect(() => {
    categoriesApi.getBySlug(params.slug).then((cat) => { setCategory(cat); setFilters((f) => ({ ...f, category_id: cat.id })); }).catch(() => {});
  }, [params.slug]);

  useEffect(() => {
    if (!filters.category_id) return;
    setIsLoading(true);
    productsApi.search(filters).then((res) => { setProducts(res.products); setTotal(res.total); }).catch(() => {}).finally(() => setIsLoading(false));
  }, [filters]);

  const totalPages = Math.ceil(total / (filters.page_size || 24));

  return (
    <div className="container py-6">
      <div className="mb-6"><h1 className="text-2xl font-bold">{category?.name || "Category"}</h1></div>
      <SearchFiltersBar filters={filters} onChange={(f) => setFilters({ ...f, page: 1 })} resultCount={total} />
      <ProductGrid products={products} isLoading={isLoading} />
      {totalPages > 1 && <Pagination currentPage={filters.page || 1} totalPages={totalPages} onPageChange={(page) => setFilters((f) => ({ ...f, page }))} />}
    </div>
  );
}
