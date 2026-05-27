import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactStrictMode: true,
  poweredByHeader: false,
  compress: true,
  output: "standalone",
  images: {
    remotePatterns: [
      { protocol: "https", hostname: "salt.tikicdn.com" },
      { protocol: "https", hostname: "tiki.vn" },
    ],
    formats: ["image/avif", "image/webp"],
    deviceSizes: [640, 750, 828, 1080, 1200],
    minimumCacheTTL: 86400,
  },
  experimental: {
    optimizePackageImports: [
      "lucide-react",
      "recharts",
    ],
  },
};

export default nextConfig;
