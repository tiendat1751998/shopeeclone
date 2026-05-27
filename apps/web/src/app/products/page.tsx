import Link from "next/link";
import { Header } from "@/components/layout/header/Header";
import { Footer } from "@/components/layout/footer/Footer";
import { ProductGrid } from "@/components/storefront/product/ProductCard";
import { extractProducts } from "@/lib/api/mapper";
import { Product } from "@/types";

const GATEWAY_URL = process.env.GATEWAY_URL || "http://gateway:8080";

async function getProducts(searchParams: Record<string, string>) {
  const sp = new URLSearchParams(searchParams);
  try {
    const res = await fetch(`${GATEWAY_URL}/api/v1/products?${sp}`, { cache: "no-store" });
    const data = await res.json();
    return extractProducts(data);
  } catch { return { products: [], total: 0, page: 1, page_size: 20, total_pages: 0 }; }
}

async function getCategories() {
  try {
    const res = await fetch(`${GATEWAY_URL}/api/v1/categories`, { cache: "no-store" });
    const data = await res.json();
    if (!Array.isArray(data)) return [];
    const children = data.flatMap((c: Record<string, unknown>) =>
      Array.isArray(c.children) ? (c.children as Array<{ id: string; name: string; slug: string }>) : []
    );
    return children;
  } catch { return []; }
}

export default async function ProductsPage({ searchParams }: { searchParams: Promise<Record<string, string>> }) {
  const sp = await searchParams;
  const [{ products }, categories] = await Promise.all([getProducts(sp), getCategories()]);

  const sortOptions = [
    { label: "Phổ biến", value: "popular", sort_by: "popularity", sort_order: "DESC" },
    { label: "Mới nhất", value: "newest", sort_by: "created_at", sort_order: "DESC" },
    { label: "Giá thấp → cao", value: "price_asc", sort_by: "price", sort_order: "ASC" },
    { label: "Giá cao → thấp", value: "price_desc", sort_by: "price", sort_order: "DESC" },
  ];

  return (
    <>
      <Header />
      <main style={{ backgroundColor: "#F5F5FA" }}>
        <div className="max-w-tiki mx-auto">
          {/* Breadcrumb */}
          <div className="flex items-center h-[36px] text-xs">
            <Link href="/" className="text-tiki-text-secondary hover:text-tiki-blue hover:underline">Trang chủ</Link>
            <svg className="mx-[5px]" width="5" height="8" viewBox="0 0 5 8" fill="none"><path d="M1 1L4 4L1 7" stroke="#808089" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/></svg>
            <span className="text-tiki-text">Tất cả sản phẩm</span>
          </div>

          <div className="flex gap-3">
            {/* Sidebar */}
            <aside className="w-[200px] shrink-0 hidden md:block">
              <div className="bg-white rounded-lg border border-tiki-border overflow-hidden sticky top-[60px]">
                <div className="px-3 py-2 border-b border-tiki-border">
                  <span className="text-[11px] font-semibold text-tiki-text">BỘ LỌC</span>
                </div>
                <div className="px-3 py-2 border-b border-tiki-border">
                  <div className="text-[10px] font-semibold text-tiki-text mb-1.5">DANH MỤC</div>
                  <div className="flex flex-col gap-1">
                    {(categories as Array<{ id: string; name: string; slug: string }>).map((cat) => (
                      <Link key={cat.id} href={`/categories/${cat.slug}`} className="text-[11px] text-tiki-text-secondary hover:text-tiki-blue transition">{cat.name}</Link>
                    ))}
                  </div>
                </div>
                <div className="px-3 py-2">
                  <div className="text-[10px] font-semibold text-tiki-text mb-1.5">GIÁ BÁN</div>
                  <div className="flex flex-col gap-1">
                    {["Dưới 500.000", "500.000 - 1.000.000", "1.000.000 - 3.000.000", "Trên 3.000.000"].map((r) => (
                      <label key={r} className="flex items-center gap-2 text-[11px] text-tiki-text-secondary cursor-pointer">
                        <input type="checkbox" className="w-3 h-3 rounded border-gray-300 accent-tiki-blue" />
                        <span>{r}</span>
                      </label>
                    ))}
                  </div>
                </div>
              </div>
            </aside>

            {/* Main */}
            <div className="flex-1 min-w-0">
              {/* Sort bar */}
              <div className="bg-white rounded-lg border border-tiki-border mb-2">
                <div className="px-3 py-2 flex items-center justify-between">
                  <span className="text-xs text-tiki-text-secondary">{(products as unknown[]).length} sản phẩm</span>
                  <div className="flex items-center gap-1">
                    <span className="text-[10px] text-tiki-text-secondary mr-1">Sắp xếp:</span>
                    {sortOptions.map((opt) => (
                      <Link
                        key={opt.value}
                        href={`/products?sort_by=${opt.sort_by}&sort_order=${opt.sort_order}`}
                        className={`px-2 py-1 text-[10px] rounded transition ${
                          sp.sort === opt.value
                            ? "bg-tiki-blue text-white"
                            : "text-tiki-text-secondary hover:bg-gray-50"
                        }`}
                      >
                        {opt.label}
                      </Link>
                    ))}
                  </div>
                </div>
              </div>

              <ProductGrid products={products as Product[]} />
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </>
  );
}
