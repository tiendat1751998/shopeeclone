import { api } from "./client";
import type { AuthTokens, LoginRequest, RegisterRequest, User } from "@/lib/types";

export const authApi = {
  login: (data: LoginRequest) => api.post<AuthTokens>("/auth/login", data),
  register: (data: RegisterRequest) => api.post<AuthTokens>("/auth/register", data),
  refresh: (refreshToken: string) => api.post<AuthTokens>("/auth/refresh", { refresh_token: refreshToken }),
  logout: () => api.post("/auth/logout", {}),
  getProfile: () => api.get<User>("/auth/me"),
  updateProfile: (data: Partial<User>) => api.put<User>("/auth/me", data),
};
