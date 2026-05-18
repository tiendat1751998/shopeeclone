import { api } from "./client";
import type { Product, PaginatedResponse, SearchResult, SearchFilters, Category, SKU } from "@/lib/types";

export const productsApi = {
  search: (filters: SearchFilters) => {
    const params = new URLSearchParams();
    if (filters.query) params.set("q", filters.query);
    if (filters.category_id) params.set("category_id", filters.category_id);
    if (filters.min_price) params.set("min_price", String(filters.min_price));
    if (filters.max_price) params.set("max_price", String(filters.max_price));
    if (filters.sort_by) params.set("sort_by", filters.sort_by);
    params.set("page", String(filters.page || 1));
    params.set("page_size", String(filters.page_size || 24));
    return api.get<SearchResult>(`/products/search?${params}`);
  },
  getById: (id: string) => api.get<Product>(`/products/${id}`),
  getByShop: (shopId: string, page = 1, pageSize = 20) =>
    api.get<PaginatedResponse<Product>>(`/products?shop_id=${shopId}&page=${page}&page_size=${pageSize}`),
  getFeatured: (limit = 10) => api.get<Product[]>(`/products/featured?limit=${limit}`),
  getDeals: (limit = 20) => api.get<Product[]>(`/products/deals?limit=${limit}`),
};

export const categoriesApi = {
  getTree: (rootId?: string) => api.get<Category[]>(`/categories${rootId ? `?root_id=${rootId}` : ""}`),
  getBySlug: (slug: string) => api.get<Category>(`/categories/${slug}`),
};

export const skusApi = {
  getByProduct: (productId: string) => api.get<SKU[]>(`/products/${productId}/skus`),
};
