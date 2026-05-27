"use client";

import { create } from "zustand";
import type { UIStore, NotificationItem } from "@/types";

export const useUIStore = create<UIStore>((set, get) => ({
  sidebarOpen: true,
  mobileMenuOpen: false,
  theme: "light",
  toastNotifications: [],

  toggleSidebar: () => set({ sidebarOpen: !get().sidebarOpen }),
  toggleMobileMenu: () => set({ mobileMenuOpen: !get().mobileMenuOpen }),
  setTheme: (theme) => set({ theme }),

  addToast: (toast) => {
    const id = crypto.randomUUID();
    const notification: NotificationItem = {
      ...toast,
      id,
      is_read: false,
      created_at: new Date().toISOString(),
    };
    set({ toastNotifications: [...get().toastNotifications, notification] });
    // Auto dismiss after 5s
    setTimeout(() => {
      set({
        toastNotifications: get().toastNotifications.filter((n) => n.id !== id),
      });
    }, 5000);
  },

  dismissToast: (id) => {
    set({
      toastNotifications: get().toastNotifications.filter((n) => n.id !== id),
    });
  },
}));
