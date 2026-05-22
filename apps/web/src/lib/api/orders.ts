import { api } from "./client";
import type { Order, PaginatedResponse } from "@/lib/types";

export interface CheckoutItem {
  product_id: string;
  sku_id: string;
  shop_id: string;
  name: string;
  quantity: number;
  unit_price: number;
  image_url?: string;
}

export interface CheckoutRequest {
  items: CheckoutItem[];
  seller_id: string;
  shipping_address: {
    street1: string;
    city: string;
    state: string;
    postal_code: string;
    country: string;
    phone: string;
  };
  billing_address: {
    street1: string;
    city: string;
    state: string;
    postal_code: string;
    country: string;
    phone: string;
  };
  payment_method?: string;
  coupon_code?: string;
  idempotency_key: string;
}

function generateIdempotencyKey(): string {
  return `${Date.now()}-${Math.random().toString(36).substring(2, 15)}`;
}

export const ordersApi = {
  checkout: (data: Omit<CheckoutRequest, "idempotency_key">) => {
    const payload: CheckoutRequest = { ...data, idempotency_key: generateIdempotencyKey() };
    return api.post<Order>("/orders", payload);
  },
  getById: (id: string) => {
    const safeId = encodeURIComponent(id);
    return api.get<Order>(`/orders/${safeId}`);
  },
  list: (page = 1, pageSize = 20, status?: string) => {
    const params = new URLSearchParams({
      page: String(Math.max(1, page)),
      page_size: String(Math.min(100, Math.max(1, pageSize))),
    });
    if (status) params.set("status", encodeURIComponent(status));
    return api.get<PaginatedResponse<Order>>(`/orders?${params}`);
  },
  cancel: (id: string, reason: string) => {
    const safeId = encodeURIComponent(id);
    const safeReason = reason?.trim().slice(0, 500) || "No reason provided";
    return api.post<Order>(`/orders/${safeId}/cancel`, { reason: safeReason });
  },
};
