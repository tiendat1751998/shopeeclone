import { Suspense } from "react";
import Link from "next/link";
import { Header } from "@/components/layout/header/Header";
import { Footer } from "@/components/layout/footer/Footer";
import { ProductGrid } from "@/components/storefront/product/ProductCard";
import { PriceFilter } from "@/components/storefront/PriceFilter";
import { CategoryListingClient } from "./CategoryListingClient";
import { Product, Category } from "@/types";
import categoriesData from "@/data/tiki-categories.json";

const SORT_OPTIONS = [
  { label: "Phổ biến", value: "popular", by: "popularity", order: "DESC" },
  { label: "Mới nhất", value: "newest", by: "created_at", order: "DESC" },
  { label: "Bán chạy", value: "best_selling", by: "sales_count", order: "DESC" },
  { label: "Giá thấp → cao", value: "price_asc", by: "price", order: "ASC" },
  { label: "Giá cao → thấp", value: "price_desc", by: "price", order: "DESC" },
];

function findCategory(cats: Category[], slug: string): Category | null {
  for (const cat of cats) {
    if (cat.slug === slug) return cat;
    if (cat.children) {
      const found = findCategory(cat.children, slug);
      if (found) return found;
    }
  }
  return null;
}

export default async function CategoryPage({
  params,
  searchParams,
}: {
  params: Promise<{ slug: string }>;
  searchParams: Promise<Record<string, string>>;
}) {
  const { slug } = await params;
  const sp = await searchParams;
  const API_BASE = process.env.GATEWAY_URL || "http://gateway:8080";
  const apiBase = `${API_BASE}/api/v1`;

  const category = findCategory(categoriesData, slug);
  const categoryName = category?.name || slug.replace(/-/g, " ");

  const queryParams = new URLSearchParams();
  queryParams.set("category_slug", slug);
  queryParams.set("page", "1");
  queryParams.set("size", "20");
  if (sp.sort_by) queryParams.set("sort_by", sp.sort_by);
  if (sp.sort_order) queryParams.set("sort_order", sp.sort_order);

  const price = sp.price;
  if (price) {
    const ranges = price.split(",");
    const minMax = ranges.reduce(
      (acc, r) => {
        const [min, max] = r.split("-");
        if (min && (!acc.min || Number(min) < acc.min)) acc.min = Number(min);
        if (max && (!acc.max || Number(max) > acc.max)) acc.max = Number(max);
        return acc;
      },
      { min: undefined as number | undefined, max: undefined as number | undefined }
    );
    if (minMax.min !== undefined) queryParams.set("min_price", String(minMax.min));
    if (minMax.max !== undefined) queryParams.set("max_price", String(minMax.max));
  }

  let initialProducts: Product[] = [];
  let total = 0;
  let totalPages = 0;

  try {
    const res = await fetch(`${apiBase}/products?${queryParams}`, { cache: "no-store" });
    const data = await res.json();
    if (Array.isArray(data)) {
      initialProducts = data as Product[];
      total = data.length;
      totalPages = 1;
    } else {
      const obj = data as Record<string, unknown>;
      initialProducts = (obj.products || obj.data || []) as Product[];
      total = (obj.total as number) || initialProducts.length;
      const ps = (obj.page_size as number) || (obj.size as number) || 20;
      totalPages = Math.ceil(total / ps) || 1;
    }
  } catch {}

  const currentSortBy = sp.sort_by || "";
  const currentSortOrder = sp.sort_order || "DESC";

  function sortActive(opt: { by: string; order: string }) {
    return currentSortBy === opt.by && currentSortOrder === opt.order;
  }

  return (
    <>
      <Header />
      <main style={{ backgroundColor: "#F5F5FA" }}>
        <div className="max-w-tiki mx-auto">
          <div className="flex items-center h-[36px] text-xs">
            <Link href="/" className="text-tiki-text-secondary hover:text-tiki-blue hover:underline">Trang chủ</Link>
            <svg className="mx-[5px]" width="5" height="8" viewBox="0 0 5 8" fill="none"><path d="M1 1L4 4L1 7" stroke="#808089" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/></svg>
            <span className="text-tiki-text">{categoryName}</span>
          </div>

          <div className="flex gap-3">
            <aside className="w-[200px] shrink-0 hidden md:block">
              <div className="bg-white rounded-lg border border-tiki-border overflow-hidden sticky top-[60px]">
                {category?.children && category.children.length > 0 && (
                  <div>
                    <div className="px-3 py-2 border-b border-tiki-border">
                      <span className="text-[11px] font-semibold text-tiki-text">DANH MỤC</span>
                    </div>
                    <div className="px-3 py-2 space-y-1.5">
                      {category.children.map((child) => (
                        <Link key={child.id} href={`/categories/${child.slug}`} className="block text-[11px] text-tiki-text-secondary hover:text-tiki-blue transition">
                          {child.name}
                        </Link>
                      ))}
                    </div>
                  </div>
                )}
                <div className="border-t border-tiki-border">
                  <div className="px-3 py-2 border-b border-tiki-border">
                    <span className="text-[11px] font-semibold text-tiki-text">GIÁ BÁN</span>
                  </div>
                  <PriceFilter basePath={`/categories/${slug}`} />
                </div>
              </div>
            </aside>

            <div className="flex-1 min-w-0">
              <div className="bg-white rounded-lg border border-tiki-border mb-2">
                <div className="px-3 py-2 flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <h1 className="text-sm font-semibold text-tiki-text">{categoryName}</h1>
                    <span className="text-[11px] text-tiki-text-secondary">({total} sản phẩm)</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <span className="text-[10px] text-tiki-text-secondary mr-1">Sắp xếp:</span>
                    {SORT_OPTIONS.map((opt) => (
                      <Link
                        key={opt.value}
                        href={`/categories/${slug}?sort_by=${opt.by}&sort_order=${opt.order}`}
                        className={`px-2 py-1 text-[10px] rounded transition ${
                          sortActive(opt) ? "bg-tiki-blue text-white" : "text-tiki-text-secondary hover:bg-gray-50"
                        }`}
                      >
                        {opt.label}
                      </Link>
                    ))}
                  </div>
                </div>
              </div>

              {category?.children && category.children.length > 0 && (
                <div className="flex flex-wrap gap-1.5 mb-2">
                  {category.children.map((child) => (
                    <Link
                      key={child.id}
                      href={`/categories/${child.slug}`}
                      className="px-2.5 py-1 text-[11px] bg-white border border-tiki-border rounded-full text-tiki-text-secondary hover:border-tiki-blue hover:text-tiki-blue transition"
                    >
                      {child.name}
                    </Link>
                  ))}
                </div>
              )}

              {/* SSR product grid */}
              <ProductGrid products={initialProducts} />

              {/* Client-side load-more */}
              <CategoryListingClient slug={slug} initialTotalPages={totalPages} />
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </>
  );
}
