import { api } from "./client";
import type { Product, PaginatedResponse, SearchResult, SearchFilters, Category, SKU } from "@/lib/types";

export const productsApi = {
  search: (filters: SearchFilters, signal?: AbortSignal) => {
    const params = new URLSearchParams();
    if (filters.query) params.set("q", filters.query);
    if (filters.category_id) params.set("category_id", filters.category_id);
    if (filters.min_price != null && Number.isFinite(filters.min_price) && filters.min_price > 0) {
      params.set("min_price", String(filters.min_price));
    }
    if (filters.max_price != null && Number.isFinite(filters.max_price) && filters.max_price > 0) {
      params.set("max_price", String(filters.max_price));
    }
    if (filters.sort_by) params.set("sort_by", filters.sort_by);
    params.set("page", String(Math.max(1, filters.page || 1)));
    params.set("page_size", String(Math.min(100, Math.max(1, filters.page_size || 24))));
    return api.get<SearchResult>(`/products/search?${params}`, signal);
  },
  getById: (id: string, signal?: AbortSignal) => {
    if (!id || typeof id !== "string") return Promise.reject(new Error("Invalid product ID"));
    const safeId = encodeURIComponent(id);
    return api.get<Product>(`/products/${safeId}`, signal);
  },
  getByShop: (shopId: string, page = 1, pageSize = 20) => {
    const safeId = encodeURIComponent(shopId);
    return api.get<PaginatedResponse<Product>>(`/products?shop_id=${safeId}&page=${page}&page_size=${pageSize}`);
  },
  getFeatured: (limit = 10) => api.get<Product[]>(`/products/featured?limit=${Math.min(100, limit)}`),
  getDeals: (limit = 20) => api.get<Product[]>(`/products/deals?limit=${Math.min(100, limit)}`),
};

export const categoriesApi = {
  getTree: (rootId?: string, signal?: AbortSignal) => {
    const params = rootId ? `?root_id=${encodeURIComponent(rootId)}` : "";
    return api.get<Category[]>(`/categories${params}`, signal);
  },
  getBySlug: (slug: string, signal?: AbortSignal) => {
    const safeSlug = encodeURIComponent(slug);
    return api.get<Category>(`/categories/${safeSlug}`, signal);
  },
};

export const skusApi = {
  getByProduct: (productId: string, signal?: AbortSignal) => {
    const safeId = encodeURIComponent(productId);
    return api.get<SKU[]>(`/products/${safeId}/skus`, signal);
  },
};
