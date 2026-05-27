export const env = {
  API_GATEWAY_URL: process.env.NEXT_PUBLIC_API_GATEWAY_URL || "http://localhost:8080",
  GRAPHQL_BFF_URL: process.env.NEXT_PUBLIC_GRAPHQL_BFF_URL || "http://localhost:4000/graphql",
  WS_URL: process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/ws",
  CDN_URL: process.env.NEXT_PUBLIC_CDN_URL || "",
  SENTRY_DSN: process.env.NEXT_PUBLIC_SENTRY_DSN || "",
  NODE_ENV: process.env.NODE_ENV || "development",
  isDev: process.env.NODE_ENV === "development",
  isProd: process.env.NODE_ENV === "production",
} as const;
