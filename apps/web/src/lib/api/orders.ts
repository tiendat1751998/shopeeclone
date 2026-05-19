import { api } from "./client";
import type { Order, PaginatedResponse, ShippingAddress } from "@/lib/types";

export interface CheckoutRequest {
  items: { product_id: string; sku_id: string; quantity: number }[];
  shipping_address: Omit<ShippingAddress, "id" | "is_default">;
  payment_method: string;
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
