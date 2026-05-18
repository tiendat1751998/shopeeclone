import { api } from "./client";
import type { Order, PaginatedResponse, ShippingAddress } from "@/lib/types";

export interface CheckoutRequest {
  items: { product_id: string; sku_id: string; quantity: number }[];
  shipping_address: Omit<ShippingAddress, "id" | "is_default">;
  payment_method: string;
  coupon_code?: string;
}

export const ordersApi = {
  checkout: (data: CheckoutRequest) => api.post<Order>("/orders", data),
  getById: (id: string) => api.get<Order>(`/orders/${id}`),
  list: (page = 1, pageSize = 20, status?: string) => {
    const params = new URLSearchParams({ page: String(page), page_size: String(pageSize) });
    if (status) params.set("status", status);
    return api.get<PaginatedResponse<Order>>(`/orders?${params}`);
  },
  cancel: (id: string, reason: string) => api.post<Order>(`/orders/${id}/cancel`, { reason }),
};
