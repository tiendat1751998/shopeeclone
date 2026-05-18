import { api } from "./client";
import type { Cart, CartItem } from "@/lib/types";

export const cartApi = {
  get: () => api.get<Cart>("/cart"),
  addItem: (productId: string, skuId: string, quantity: number) =>
    api.post<Cart>("/cart/items", { product_id: productId, sku_id: skuId, quantity }),
  updateItem: (itemId: string, quantity: number) =>
    api.patch<Cart>(`/cart/items/${itemId}`, { quantity }),
  removeItem: (itemId: string) => api.delete<Cart>(`/cart/items/${itemId}`),
  selectItem: (itemId: string, selected: boolean) =>
    api.patch<Cart>(`/cart/items/${itemId}/select`, { selected }),
  clear: () => api.delete("/cart"),
};
