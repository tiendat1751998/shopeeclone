# TIKI CLONE - BUILD STATUS

## Completed

### Database
- tiki_categories: 20 categories with Vietnamese names
- tiki_products: 100 products across 16 categories in VND prices
- Crawler scripts ready at scripts/

### Frontend Components (NEW)
- src/components/tiki/TikiHeader.tsx - Blue header with logo, search, account, cart, location
- src/components/tiki/TikiProductCard.tsx - Product card with discount badge, stars, sold count, seller
- src/components/tiki/TikiHeroCarousel.tsx - Hero carousel + quick links + categories grid
- src/components/tiki/TikiFooter.tsx - Full footer with support, about, payments, shipping
- src/styles/globals.css - Tiki-themed global styles

### Frontend Pages (UPDATED)
- src/app/layout.tsx - Uses TikiHeader + TikiFooter
- src/app/page.tsx - Home page with carousel, categories grid, product sections, service highlights

### API Routes (NEW)
- src/app/api/tiki/products/route.ts - GET /api/tiki/products with pagination, sorting, filtering

## Pending Manual Steps

### 1. Fix MySQL connection in API
The Next.js API route at src/app/api/tiki/products/route.ts needs to connect to MySQL.
Since the web container can't resolve "mysql-primary" hostname, either:
- a) Change host from "mysql-primary" to "172.18.0.10" in route.ts
- b) Or rebuild web container with proper network access

### 2. Rebuild web container
docker build -t shopeeclone-web -f apps/web/Dockerfile apps/web
docker restart shopeeclone-web-1

### 3. Test the frontend
Visit http://localhost (through nginx proxy) and verify:
- Header shows blue Tiki-style with search bar, account, cart
- Home page shows carousel, categories grid, product cards
- Product cards show discount %, rating stars, sold count, seller name
- Footer shows proper links

## API Endpoints

### GET /api/tiki/products
Query params:
- category_id: filter by category
- page: page number (default 1)
- limit: page size (default 24, max 100)
- sort: popular|price_asc|price_desc|discount|rating|newest
- q: search query
- min_price/max_price: price range filter

Response:
{
  success: true,
  data: {
    products: [...],
    total: 100,
    page: 1,
    page_size: 24,
    total_pages: 5
  }
}
