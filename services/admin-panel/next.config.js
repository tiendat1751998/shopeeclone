const nextConfig = {
  output: 'standalone',
  reactStrictMode: true,
  swcMinify: true,
  images: {
    domains: ['localhost'],
  },
  rewrites: async () => {
    return [
      {
        source: '/api/auth/:path*',
        destination: 'http://auth:8080/api/v1/auth/:path*',
      },
      {
        source: '/api/users/:path*',
        destination: 'http://auth:8080/api/v1/admin/users/:path*',
      },
      {
        source: '/api/products/:path*',
        destination: 'http://gateway:8080/api/v1/products/:path*',
      },
      {
        source: '/api/orders/:path*',
        destination: 'http://gateway:8080/api/v1/orders/:path*',
      },
      {
        source: '/api/inventory/:path*',
        destination: 'http://gateway:8080/api/v1/inventory/:path*',
      },
      {
        source: '/api/promotions/:path*',
        destination: 'http://gateway:8080/api/v1/promotions/:path*',
      },
      {
        source: '/api/payments/:path*',
        destination: 'http://gateway:8080/api/v1/payments/:path*',
      },
      {
        source: '/api/shipments/:path*',
        destination: 'http://gateway:8080/api/v1/shipments/:path*',
      },
      {
        source: '/api/categories/:path*',
        destination: 'http://gateway:8080/api/v1/categories/:path*',
      },
      {
        source: '/api/search/:path*',
        destination: 'http://gateway:8080/api/v1/search/:path*',
      },
    ];
  },
};

module.exports = nextConfig;
