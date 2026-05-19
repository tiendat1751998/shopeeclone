import { useState, useEffect, useCallback, useRef } from "react";
import { productsApi, categoriesApi } from "@/lib/api/products";
import { useDebounce } from "./useDebounce";
import type { Product, SearchFilters, SearchResult, Category } from "@/lib/types";

export function useProductSearch(initialFilters: SearchFilters = {}) {
  const [filters, setFilters] = useState<SearchFilters>(initialFilters);
  const [result, setResult] = useState<SearchResult | null>(null);
  const [isLoading, setIsLoading] = useState(false);
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
    return () => {
      if (abortRef.current) abortRef.current.abort();
    };
  }, [search]);

  return { result, isLoading, error, filters, setFilters, refetch: search };
}

export function useProduct(productId: string) {
  const [product, setProduct] = useState<Product | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const abortRef = useRef<AbortController | null>(null);

  useEffect(() => {
    if (abortRef.current) abortRef.current.abort();
    const controller = new AbortController();
    abortRef.current = controller;
    setIsLoading(true);

    productsApi.getById(productId, controller.signal)
      .then((p) => { if (!controller.signal.aborted) { setProduct(p); setIsLoading(false); } })
      .catch((e: Error) => { if (!controller.signal.aborted) { setError(e.message); setIsLoading(false); } });

    return () => controller.abort();
  }, [productId]);

  return { product, isLoading, error };
}

export function useCategories() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const abortRef = useRef<AbortController | null>(null);

  useEffect(() => {
    const controller = new AbortController();
    abortRef.current = controller;

    categoriesApi.getTree(undefined, controller.signal)
      .then((cats) => { if (!controller.signal.aborted) setCategories(cats); })
      .catch(() => undefined)
      .finally(() => { if (!controller.signal.aborted) setIsLoading(false); });

    return () => controller.abort();
  }, []);

  return { categories, isLoading };
}
