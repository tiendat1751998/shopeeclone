"use client";

import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { CartStore, CartItem } from "@/types";
import { cartApi } from "@/lib/api/client";
import { useAuthStore } from "./auth";

export const useCartStore = create<CartStore>()(
  persist(
    (set, get) => ({
      items: [],
      isLoading: false,

      addItem: async (item) => {
        const prev = get().items;
        set({ isLoading: true });

        // Optimistic update
        const existing = prev.find(
          (i) => i.product_id === item.product_id && i.sku_id === item.sku_id
        );
        if (existing) {
          set({
            items: prev.map((i) =>
              i.id === existing.id
                ? { ...i, quantity: i.quantity + item.quantity }
                : i
            ),
          });
        } else {
          set({
            items: [
              ...prev,
              { ...item, id: crypto.randomUUID(), is_selected: true },
            ],
          });
        }

        // Only call API if user is authenticated
        const isAuthenticated = useAuthStore.getState().isAuthenticated;
        if (!isAuthenticated) {
          set({ isLoading: false });
          return;
        }

        try {
          const cart = await cartApi.addItem(
            item.product_id,
            item.sku_id,
            item.quantity
          );
          set({ items: cart.items ?? [], isLoading: false });
        } catch {
          set({ items: prev, isLoading: false });
        }
      },

      removeItem: async (id) => {
        const prev = get().items;
        set({ items: prev.filter((i) => i.id !== id) });

        // Only call API if user is authenticated
        const isAuthenticated = useAuthStore.getState().isAuthenticated;
        if (!isAuthenticated) {
          return;
        }

        try {
          await cartApi.removeItem(id);
        } catch {
          set({ items: prev });
        }
      },

      updateQuantity: async (id, quantity) => {
        const prev = get().items;
        set({
          items: prev.map((i) =>
            i.id === id ? { ...i, quantity } : i
          ),
        });

        // Only call API if user is authenticated
        const isAuthenticated = useAuthStore.getState().isAuthenticated;
        if (!isAuthenticated) {
          return;
        }

        try {
          await cartApi.updateItem(id, quantity);
        } catch {
          set({ items: prev });
        }
      },

      toggleSelect: (id) => {
        set({
          items: get().items.map((i) =>
            i.id === id ? { ...i, is_selected: !i.is_selected } : i
          ),
        });
      },

      selectAll: () => {
        set({ items: get().items.map((i) => ({ ...i, is_selected: true })) });
      },

      clearCart: () => {
        const prev = get().items;
        set({ items: [] });

        // Only call API if user is authenticated
        const isAuthenticated = useAuthStore.getState().isAuthenticated;
        if (!isAuthenticated) {
          return;
        }

        cartApi.clear().catch(() => {
          set({ items: prev });
        });
      },

      getSelectedItems: () => get().items.filter((i) => i.is_selected),

      getSubtotal: () =>
        (get().items ?? [])
          .filter((i) => i.is_selected)
          .reduce((sum, i) => sum + i.price * i.quantity, 0),

      getTotal: () => get().getSubtotal(),
    }),
    {
      name: "shopee-cart",
      partialize: (state) => ({ items: state.items }),
    }
  )
);
