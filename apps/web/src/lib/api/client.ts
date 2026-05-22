import type { ApiResponse } from "@/lib/types";

const isServer = typeof window === "undefined";
const API_BASE = isServer
  ? (process.env.API_GATEWAY_URL || "http://gateway:8080") + "/api/v1"
  : (process.env.NEXT_PUBLIC_API_URL || "/api/gateway");

class ApiError extends Error {
  public code: string;
  public status: number;
  constructor(status: number, code: string, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
  }
}

function getCSRFToken(): string | null {
  if (typeof document === "undefined") return null;
  const match = document.cookie.match(/csrf_token=([^;]+)/);
  return match ? decodeURIComponent(match[1]) : null;
}

async function request<T>(path: string, options: RequestInit = {}, signal?: AbortSignal): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };

  const csrfToken = getCSRFToken();
  if (csrfToken) {
    headers["X-CSRF-Token"] = csrfToken;
  }

  const fetchOptions: RequestInit = { ...options, headers, credentials: "include" };
  if (signal) fetchOptions.signal = signal;

  const res = await fetch(`${API_BASE}${path}`, fetchOptions);

  if (res.status === 401) {
    throw new ApiError(401, "UNAUTHORIZED", "Authentication required");
  }
  if (res.status === 403) {
    throw new ApiError(403, "FORBIDDEN", "Access denied");
  }
  if (res.status === 429) {
    throw new ApiError(429, "RATE_LIMITED", "Too many requests. Please try again later.");
  }

  let data: Record<string, unknown>;
  try {
    data = await res.json();
  } catch {
    throw new ApiError(res.status, "PARSE_ERROR", `Server returned ${res.status}`);
  }

  if (!res.ok) {
    throw new ApiError(
      res.status,
      (data as unknown as { code?: string }).code || "UNKNOWN",
      (data as unknown as { error?: string; message?: string }).error ||
        (data as unknown as { message?: string }).message ||
        `Request failed with status ${res.status}`
    );
  }

  if (data && typeof data === "object" && "success" in data) {
    return (data as unknown as ApiResponse<T>).data as T;
  }

  return data as unknown as T;
}

export const api = {
  get: <T>(path: string, signal?: AbortSignal) => request<T>(path, {}, signal),
  post: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "POST", body: JSON.stringify(body) }),
  put: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "PUT", body: JSON.stringify(body) }),
  patch: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "PATCH", body: JSON.stringify(body) }),
  delete: <T>(path: string) => request<T>(path, { method: "DELETE" }),
};

export { ApiError };
