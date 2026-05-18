import { useState, useEffect, useCallback } from "react";
import { productsApi } from "@/lib/api/products";
import { useDebounce } from "./useDebounce";
import type { Product, SearchFilters, SearchResult, Category } from "@/lib/types";

export function useProductSearch(initialFilters: SearchFilters = {}) {
  const [filters, setFilters] = useState<SearchFilters>(initialFilters);
  const [result, setResult] = useState<SearchResult | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const debouncedQuery = useDebounce(filters.query, 300);

  const search = useCallback(async () => {
    setIsLoading(true); setError(null);
    try {
      const res = await productsApi.search({ ...filters, query: debouncedQuery });
      setResult(res);
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Search failed");
    } finally { setIsLoading(false); }
  }, [filters, debouncedQuery]);

  useEffect(() => { search(); }, [search]);

  return { result, isLoading, error, filters, setFilters, refetch: search };
}

export function useProduct(productId: string) {
  const [product, setProduct] = useState<Product | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setIsLoading(true);
    productsApi.getById(productId)
      .then(setProduct)
      .catch((e: Error) => setError(e.message))
      .finally(() => setIsLoading(false));
  }, [productId]);

  return { product, isLoading, error };
}

export function useCategories() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    categoriesApi.getTree()
      .then(setCategories)
      .catch(() => {})
      .finally(() => setIsLoading(false));
  }, []);

  return { categories, isLoading };
}
