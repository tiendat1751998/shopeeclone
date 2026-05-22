import type { NextConfig } from 'next';

const API_GATEWAY_URL = process.env.API_GATEWAY_URL || 'http://localhost:8080';

const nextConfig: NextConfig = {
  output: 'standalone',
  images: {
    remotePatterns: [
      { protocol: 'https', hostname: '**.shopee.vn' },
      { protocol: 'https', hostname: '**.shopee.sg' },
      { protocol: 'https', hostname: '**.shopee.com' },
      { protocol: 'https', hostname: 'cdn.example.com' },
      { protocol: 'https', hostname: 'store.storeimages.cdn-apple.com' },
      { protocol: 'https', hostname: 'images.samsung.com' },
      { protocol: 'http', hostname: 'localhost', port: '9000' },
      { protocol: 'https', hostname: 'hoanghamobile.com' },
      { protocol: 'https', hostname: '**.hoanghamobile.com' },
    ],
    formats: ['image/avif', 'image/webp'],
  },
  async rewrites() {
    return [
      {
        source: '/api/gateway/:path*',
        destination: `${API_GATEWAY_URL}/api/v1/:path*`,
      },
    ];
  },
  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          { key: 'X-Content-Type-Options', value: 'nosniff' },
          { key: 'X-Frame-Options', value: 'DENY' },
          { key: 'X-XSS-Protection', value: '1; mode=block' },
          { key: 'Referrer-Policy', value: 'strict-origin-when-cross-origin' },
        ],
      },
    ];
  },
};

export default nextConfig;
