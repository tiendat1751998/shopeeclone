import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { CartItem } from "@/lib/types";
import { cartApi } from "@/lib/api/cart";

interface CartState {
  items: CartItem[];
  isLoading: boolean;
  error: string | null;
  fetchCart: () => Promise<void>;
  addItem: (productId: string, skuId: string, quantity: number) => Promise<void>;
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
        set({ isLoading: true });
        try {
          const cart = await cartApi.get();
          set({ items: cart.items, isLoading: false });
        } catch {
          set({ isLoading: false });
        }
      },

      addItem: async (productId, skuId, quantity) => {
        set({ isLoading: true, error: null });
        try {
          const cart = await cartApi.addItem(productId, skuId, quantity);
          set({ items: cart.items, isLoading: false });
        } catch (e: unknown) {
          const msg = e instanceof Error ? e.message : "Failed to add item";
          set({ isLoading: false, error: msg });
        }
      },

      updateQuantity: async (itemId, quantity) => {
        if (quantity < 1) { get().removeItem(itemId); return; }
        set({ isLoading: true });
        try {
          const cart = await cartApi.updateItem(itemId, quantity);
          set({ items: cart.items, isLoading: false });
        } catch { set({ isLoading: false }); }
      },

      removeItem: async (itemId) => {
        set({ isLoading: true });
        try {
          const cart = await cartApi.removeItem(itemId);
          set({ items: cart.items, isLoading: false });
        } catch { set({ isLoading: false }); }
      },

      toggleSelect: (itemId) => {
        set({ items: get().items.map((i) => i.id === itemId ? { ...i, is_selected: !i.is_selected } : i) });
      },

      toggleSelectAll: () => {
        const allSelected = get().items.every((i) => i.is_selected);
        set({ items: get().items.map((i) => ({ ...i, is_selected: !allSelected })) });
      },

      clearCart: async () => {
        try { await cartApi.clear(); } catch { /* ignore */ }
        set({ items: [] });
      },

      selectedItems: () => get().items.filter((i) => i.is_selected),

      subtotal: () => get().items.filter((i) => i.is_selected).reduce((sum, i) => sum + i.price * i.quantity, 0),

      totalItems: () => get().items.reduce((sum, i) => sum + i.quantity, 0),
    }),
    { name: "shopee-cart", partialize: (state) => ({ items: state.items }) }
  )
);
