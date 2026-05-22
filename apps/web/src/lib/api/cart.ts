import { api } from "./client";
import type { Cart, CartItem } from "@/lib/types";

const MAX_QTY = 99;
const MIN_QTY = 1;

function validateQuantity(qty: number): number {
  if (!Number.isFinite(qty)) return MIN_QTY;
  return Math.max(MIN_QTY, Math.min(MAX_QTY, Math.floor(qty)));
}

function validateId(id: string): string {
  if (!id || typeof id !== "string") throw new Error("Invalid ID");
  return encodeURIComponent(id);
}

export const cartApi = {
  get: () => api.get<Cart>("/cart"),
  addItem: (productId: string, skuId: string, quantity: number, name?: string, price?: number, shopId?: string, shopName?: string, imageUrl?: string) => {
    const safePid = validateId(productId);
    const safeSid = validateId(skuId);
    const safeQty = validateQuantity(quantity);
    return api.post<Cart>(`/cart/items`, { product_id: safePid, sku_id: safeSid, quantity: safeQty, name, price, shop_id: shopId, shop_name: shopName, image_url: imageUrl });
  },
  updateItem: (itemId: string, quantity: number) => {
    const safeId = validateId(itemId);
    const safeQty = validateQuantity(quantity);
    return api.patch<Cart>(`/cart/items/${safeId}`, { quantity: safeQty });
  },
  removeItem: (itemId: string) => {
    const safeId = validateId(itemId);
    return api.delete<Cart>(`/cart/items/${safeId}`);
  },
  selectItem: (itemId: string, selected: boolean) => {
    const safeId = validateId(itemId);
    return api.patch<Cart>(`/cart/items/${safeId}/select`, { selected: !!selected });
  },
  clear: () => api.delete("/cart"),
};
