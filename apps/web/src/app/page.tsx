import { Suspense } from "react";
import Link from "next/link";
import Image from "next/image";
import { Header } from "@/components/layout/header/Header";
import { Footer } from "@/components/layout/footer/Footer";
import { ProductGrid } from "@/components/storefront/product/ProductCard";
import { ProductGridSkeleton } from "@/components/ui";
import { FlashSaleTimer } from "@/components/storefront/FlashSaleTimer";
import { mapProductArray } from "@/lib/api/mapper";
import { Product } from "@/types";
import categoriesData from "@/data/tiki-categories.json";
import promotionsData from "@/data/tiki-promotions.json";

const GATEWAY_URL = process.env.GATEWAY_URL || "http://gateway:8080";

/* ── Hero Banner ── */
function Hero() {
  const promo = (promotionsData[0] as Record<string, unknown>) || {};
  const quickLinks = (promo.quick_links as Array<Record<string, string>>) || [];
  const sideBanner1 = quickLinks[0];
  const sideBanner2 = quickLinks[1];
  const heroImage = (promo.image_url as string) || null;

  return (
    <section className="mb-3">
      <div className="max-w-tiki mx-auto px-3">
        <div className="flex gap-2">
          {/* Main hero — Tiki ratio: ~2/3 width, compact height */}
          <Link href="/promotions" className="hero-banner flex-[2] relative" style={{ minHeight: 240 }}>
            {heroImage ? (
              <Image
                src={heroImage}
                alt="promotion"
                fill
                sizes="(max-width: 768px) 66vw, 60vw"
                style={{ objectFit: "cover" }}
                className="rounded-lg"
                priority
              />
            ) : (
              <div className="hero-fallback">
                <span className="hero-fallback__icon">🎉</span>
                <span className="hero-fallback__text">Khuyến Mãi Đặc Biệt</span>
                <span className="hero-fallback__cta">Xem ngay →</span>
              </div>
            )}
          </Link>

          {/* Side banners — Tiki style: stacked, compact */}
          <div className="flex-[1] flex flex-col gap-2 min-w-0">
            {sideBanner1 ? (
              <Link href="/promotions" className="promo-banner flex-1 relative" style={{ minHeight: 116 }}>
                {sideBanner1.image_url ? (
                  <Image
                    src={sideBanner1.image_url}
                    alt={sideBanner1.name}
                    fill
                    sizes="(max-width: 768px) 33vw, 20vw"
                    style={{ objectFit: "cover" }}
                    className="rounded-lg"
                  />
                ) : (
                  <div className="promo-banner__fallback">
                    <span>{sideBanner1.name}</span>
                  </div>
                )}
                <div className="promo-banner__overlay" />
              </Link>
            ) : (
              <div className="promo-banner flex-1" style={{ minHeight: 116 }}>
                <div className="promo-banner__fallback"><span>Khuyến Mãi</span></div>
              </div>
            )}
            {sideBanner2 ? (
              <Link href="/promotions" className="promo-banner flex-1 relative" style={{ minHeight: 116 }}>
                {sideBanner2.image_url ? (
                  <Image
                    src={sideBanner2.image_url}
                    alt={sideBanner2.name}
                    fill
                    sizes="(max-width: 768px) 33vw, 20vw"
                    style={{ objectFit: "cover" }}
                    className="rounded-lg"
                  />
                ) : (
                  <div className="promo-banner__fallback promo-banner__fallback--blue">
                    <span>{sideBanner2.name}</span>
                  </div>
                )}
                <div className="promo-banner__overlay" />
              </Link>
            ) : (
              <div className="promo-banner flex-1" style={{ minHeight: 116 }}>
                <div className="promo-banner__fallback promo-banner__fallback--blue"><span>Ưu Đãi Đặc Biệt</span></div>
              </div>
            )}
          </div>
        </div>
      </div>
    </section>
  );
}

/* ── Categories Grid ── */
const CATEGORY_ICONS: Record<string, string> = {
  "nha-cua-doi-song": "🏠",
  "nha-bep": "🍳",
  "dien-tu": "📱",
  "thoi-trang": "👗",
  "lam-dep": "💄",
  "me-be": "👶",
  "sach": "📚",
  "the-thao": "⚽",
  "oto-xe-may": "🚗",
  "dien-lanh": "❄️",
  "bach-hoa": "🛒",
  "choi-choi": "🧸",
  "thu-cung": "🐾",
  "nha-sach": "📖",
  "dien-gia-dung": "🔌",
};

function getCategoryEmoji(slug: string): string {
  for (const [key, icon] of Object.entries(CATEGORY_ICONS)) {
    if (slug.includes(key)) return icon;
  }
  return "📦";
}

interface CategoryItem {
  id: string;
  slug: string;
  name: string;
  image_url?: string;
}

function CategoriesGrid({ categories }: { categories: CategoryItem[] }) {
  const cats = categories.slice(0, 20);

  return (
    <section className="mb-6">
      <div className="max-w-tiki mx-auto px-4 sm:px-6">
        <div className="bg-white rounded-xl border border-tiki-border overflow-hidden">
          <div className="grid grid-cols-5 sm:grid-cols-6 md:grid-cols-8 lg:grid-cols-10 gap-0.5 p-1.5">
            {cats.map((cat) => (
              <Link
                key={cat.id}
                href={`/categories/${cat.slug}`}
                className="category-card"
              >
                <div className="category-card__icon-wrap">
                  {cat.image_url ? (
                    <img
                      src={cat.image_url}
                      alt={cat.name}
                      className="category-card__img"
                      loading="lazy"
                    />
                  ) : (
                    <span className="category-card__icon">{getCategoryEmoji(cat.slug)}</span>
                  )}
                </div>
                <span className="category-card__name">{cat.name}</span>
              </Link>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}

/* ── Product Section ── */
function ProductSection({ title, products }: { title: string; products: Product[] }) {
  return (
    <section className="mb-4">
      <div className="max-w-tiki mx-auto px-3">
        <div className="bg-white rounded-lg border border-tiki-border overflow-hidden">
          <div className="px-3 py-2 flex items-center justify-between border-b border-tiki-border">
            <h2 className="text-sm font-semibold text-tiki-text">{title}</h2>
            <Link href="/products" className="text-xs font-medium text-tiki-blue hover:opacity-80 flex items-center gap-1">
              Xem tất cả <span>→</span>
            </Link>
          </div>
          <div className="p-2">
            <ProductGrid products={products} />
          </div>
        </div>
      </div>
    </section>
  );
}

/* ── Service Highlights ── */
function ServiceHighlights() {
  const services = [
    { icon: "🛡️", title: "100% hàng thật", desc: "Cam kết hàng chính hãng" },
    { icon: "🚚", title: "Freeship mọi đơn", desc: "Miễn phí giao hàng toàn quốc" },
    { icon: "💰", title: "Hoàn 200% nếu giả", desc: "Bảo vệ người tiêu dùng" },
    { icon: "🔄", title: "30 ngày đổi trả", desc: "Đổi trả miễn phí" },
    { icon: "⚡", title: "Giao nhanh 2h", desc: "Giao hàng nhanh chóng" },
  ];

  return (
    <section className="mb-4">
      <div className="max-w-tiki mx-auto px-3">
        <div className="grid grid-cols-3 sm:grid-cols-5 gap-2">
          {services.map((s) => (
            <div key={s.title} className="service-card">
              <div className="text-lg mb-1">{s.icon}</div>
              <div className="text-[10px] font-semibold text-tiki-text mb-0.5">{s.title}</div>
              <div className="text-[9px] text-tiki-text-secondary">{s.desc}</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

/* ── flatten categories ── */
function flattenCategories(data: typeof categoriesData): CategoryItem[] {
  const result: CategoryItem[] = [];
  for (const cat of data) {
    result.push({ id: cat.id, slug: cat.slug, name: cat.name, image_url: cat.image_url });
    if (cat.children) {
      for (const child of cat.children as Array<{ id: string; slug: string; name: string; image_url?: string }>) {
        result.push({ id: child.id, slug: child.slug, name: child.name, image_url: child.image_url });
      }
    }
  }
  return result;
}

/* ── Page ── */
export default async function HomePage() {
  let featured: Product[] = [];
  let deals: Product[] = [];
  const categories = flattenCategories(categoriesData);

  try {
    const baseUrl = process.env.GATEWAY_URL || "http://gateway:8080";
    const [featRes, dealsRes] = await Promise.all([
      fetch(`${baseUrl}/api/v1/products/featured?limit=10`, { cache: "no-store" }).then(r => r.json()).catch(() => null),
      fetch(`${baseUrl}/api/v1/products/deals?limit=10`, { cache: "no-store" }).then(r => r.json()).catch(() => null),
    ]);
    featured = featRes ? mapProductArray(Array.isArray(featRes) ? featRes : featRes.data || []) : [];
    deals = dealsRes ? mapProductArray(Array.isArray(dealsRes) ? dealsRes : dealsRes.data || []) : [];
  } catch {}

  return (
    <>
      <Header />
      <main className="py-2 sm:py-3">
        <Hero />
        <CategoriesGrid categories={categories} />
        <FlashSaleTimer />
        <Suspense fallback={<div className="max-w-tiki mx-auto px-6"><ProductGridSkeleton count={10} /></div>}>
          <ProductSection title="🔥 Giá Sốc Hôm Nay" products={deals} />
        </Suspense>
        <Suspense fallback={<div className="max-w-tiki mx-auto px-6"><ProductGridSkeleton count={10} /></div>}>
          <ProductSection title="⭐ Sản Phẩm Nổi Bật" products={featured} />
        </Suspense>
        <ServiceHighlights />
      </main>
      <Footer />
    </>
  );
}
