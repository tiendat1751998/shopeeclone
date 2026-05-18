import { create } from "zustand";
import type { User, AuthTokens } from "@/lib/types";
import { authApi } from "@/lib/api/auth";

interface AuthState {
  user: User | null;
  tokens: AuthTokens | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (data: { email: string; password: string; username: string; display_name: string }) => Promise<void>;
  logout: () => void;
  refreshToken: () => Promise<boolean>;
  fetchProfile: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  tokens: null,
  isAuthenticated: false,
  isLoading: false,

  login: async (email, password) => {
    set({ isLoading: true });
    try {
      const tokens = await authApi.login({ email, password });
      localStorage.setItem("access_token", tokens.access_token);
      localStorage.setItem("refresh_token", tokens.refresh_token);
      const user = await authApi.getProfile();
      set({ user, tokens, isAuthenticated: true, isLoading: false });
    } catch (e: unknown) {
      set({ isLoading: false });
      throw e;
    }
  },

  register: async (data) => {
    set({ isLoading: true });
    try {
      const tokens = await authApi.register(data);
      localStorage.setItem("access_token", tokens.access_token);
      localStorage.setItem("refresh_token", tokens.refresh_token);
      const user = await authApi.getProfile();
      set({ user, tokens, isAuthenticated: true, isLoading: false });
    } catch (e: unknown) {
      set({ isLoading: false });
      throw e;
    }
  },

  logout: () => {
    authApi.logout().catch(() => {});
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    set({ user: null, tokens: null, isAuthenticated: false });
  },

  refreshToken: async () => {
    const refreshToken = localStorage.getItem("refresh_token");
    if (!refreshToken) { get().logout(); return false; }
    try {
      const tokens = await authApi.refresh(refreshToken);
      localStorage.setItem("access_token", tokens.access_token);
      set({ tokens, isAuthenticated: true });
      return true;
    } catch { get().logout(); return false; }
  },

  fetchProfile: async () => {
    try {
      const user = await authApi.getProfile();
      set({ user, isAuthenticated: true });
    } catch { set({ isAuthenticated: false }); }
  },
}));
