import { create } from "zustand";
import type { User, AuthTokens } from "@/lib/types";
import { authApi } from "@/lib/api/auth";

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (data: {
    email: string;
    password: string;
    username: string;
    display_name: string;
  }) => Promise<void>;
  logout: () => void;
  fetchProfile: () => Promise<void>;
}

function setCookie(name: string, value: string, maxAge: number) {
  document.cookie = `${name}=${encodeURIComponent(value)}; path=/; SameSite=Strict; Max-Age=${maxAge}`;
}

function clearCookie(name: string) {
  document.cookie = `${name}=; path=/; Max-Age=0; SameSite=Strict`;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  isLoading: false,

  login: async (email, password) => {
    set({ isLoading: true });
    try {
      const tokens: AuthTokens = await authApi.login({ email, password });
      setCookie("access_token", tokens.access_token, tokens.expires_in);
      setCookie("refresh_token", tokens.refresh_token, 7 * 24 * 3600);
      const user = await authApi.getProfile();
      set({ user, isAuthenticated: true, isLoading: false });
    } catch (e: unknown) {
      set({ isLoading: false, isAuthenticated: false, user: null });
      throw e;
    }
  },

  register: async (data) => {
    set({ isLoading: true });
    try {
      const tokens: AuthTokens = await authApi.register(data);
      setCookie("access_token", tokens.access_token, tokens.expires_in);
      setCookie("refresh_token", tokens.refresh_token, 7 * 24 * 3600);
      const user = await authApi.getProfile();
      set({ user, isAuthenticated: true, isLoading: false });
    } catch (e: unknown) {
      set({ isLoading: false, isAuthenticated: false, user: null });
      throw e;
    }
  },

  logout: () => {
    authApi.logout().catch(() => undefined);
    clearCookie("access_token");
    clearCookie("refresh_token");
    set({ user: null, isAuthenticated: false });
  },

  fetchProfile: async () => {
    try {
      const user = await authApi.getProfile();
      set({ user, isAuthenticated: true });
    } catch {
      set({ isAuthenticated: false, user: null });
    }
  },
}));
