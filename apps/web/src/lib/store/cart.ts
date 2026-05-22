import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { CartItem } from "@/lib/types";
import { cartApi } from "@/lib/api/cart";

const MAX_QUANTITY = 99;
const MIN_QUANTITY = 1;

function clampQuantity(qty: number, stock: number): number {
  if (!Number.isFinite(qty) || qty < MIN_QUANTITY) return MIN_QUANTITY;
  if (qty > stock) return Math.max(MIN_QUANTITY, stock);
  if (qty > MAX_QUANTITY) return MAX_QUANTITY;
  return qty;
}

interface CartState {
  items: CartItem[];
  isLoading: boolean;
  error: string | null;
  fetchCart: () => Promise<void>;
  addItem: (productId: string, skuId: string, quantity: number, name?: string, price?: number, shopId?: string, shopName?: string, imageUrl?: string) => Promise<void>;
  updateQuantity: (itemId: string, quantity: number) => Promise<void>;
  removeItem: (itemId: string) => Promise<void>;
  toggleSelect: (itemId: string) => void;
  toggleSelectAll: () => void;
  clearCart: () => Promise<void>;
  selectedItems: () => CartItem[];
  subtotal: () => number;
  totalItems: () => number;
}

export const useCartStore = create<CartState>()(
  persist(
    (set, get) => ({
      items: [],
      isLoading: false,
      error: null,

      fetchCart: async () => {
        set({ isLoading: true, error: null });
        try {
          const cart = await cartApi.get();
          set({ items: cart.items || [], isLoading: false });
        } catch (e: unknown) {
          set({ isLoading: false, error: e instanceof Error ? e.message : "Failed to fetch cart" });
        }
      },

      addItem: async (productId, skuId, quantity, name?, price?, shopId?, shopName?, imageUrl?) => {
        const safeQty = clampQuantity(quantity, MAX_QUANTITY);
        set({ isLoading: true, error: null });
        try {
          const cart = await cartApi.addItem(productId, skuId, safeQty, name, price, shopId, shopName, imageUrl);
          set({ items: cart.items || [], isLoading: false });
        } catch (e: unknown) {
          set({ isLoading: false, error: e instanceof Error ? e.message : "Failed to add item" });
        }
      },

      updateQuantity: async (itemId, quantity) => {
        if (quantity < MIN_QUANTITY) {
          await get().removeItem(itemId);
          return;
        }
        const item = get().items.find((i) => i.id === itemId);
        const safeQty = clampQuantity(quantity, item?.stock ?? MAX_QUANTITY);
        set({ isLoading: true, error: null });
        try {
          const cart = await cartApi.updateItem(itemId, safeQty);
          set({ items: cart.items || [], isLoading: false });
        } catch (e: unknown) {
          set({ isLoading: false, error: e instanceof Error ? e.message : "Failed to update quantity" });
        }
      },

      removeItem: async (itemId) => {
        set({ isLoading: true, error: null });
        try {
          const cart = await cartApi.removeItem(itemId);
          set({ items: cart.items || [], isLoading: false });
        } catch (e: unknown) {
          set({ isLoading: false, error: e instanceof Error ? e.message : "Failed to remove item" });
        }
      },

      toggleSelect: (itemId) => {
        set({ items: get().items.map((i) => i.id === itemId ? { ...i, is_selected: !i.is_selected } : i) });
      },

      toggleSelectAll: () => {
        const allSelected = get().items.length > 0 && get().items.every((i) => i.is_selected);
        set({ items: get().items.map((i) => ({ ...i, is_selected: !allSelected })) });
      },

      clearCart: async () => {
        try { await cartApi.clear(); } catch { /* ignore */ }
        set({ items: [], error: null });
      },

      selectedItems: () => get().items.filter((i) => i.is_selected),

      subtotal: () =>
        get().items
          .filter((i) => i.is_selected)
          .reduce((sum, i) => {
            const price = Number.isFinite(i.price) && i.price >= 0 ? i.price : 0;
            const qty = Number.isFinite(i.quantity) && i.quantity > 0 ? i.quantity : 0;
            return sum + price * qty;
          }, 0),

      totalItems: () =>
        get().items.reduce((sum, i) => {
          const qty = Number.isFinite(i.quantity) && i.quantity > 0 ? i.quantity : 0;
          return sum + qty;
        }, 0),
    }),
    { name: "shopee-cart", partialize: (state) => ({ items: state.items }) }
  )
);
