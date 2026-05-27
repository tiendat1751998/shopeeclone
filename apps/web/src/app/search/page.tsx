import Link from "next/link";
import { Header } from "@/components/layout/header/Header";
import { Footer } from "@/components/layout/footer/Footer";
import { ProductGrid } from "@/components/storefront/product/ProductCard";
import { extractProducts } from "@/lib/api/mapper";
import { Product } from "@/types";

const GATEWAY_URL = process.env.GATEWAY_URL || "http://gateway:8080";

async function getSearchResults(q: string, sort_by?: string, sort_order?: string) {
  try {
    const params = new URLSearchParams({ q });
    if (sort_by) params.set("sort_by", sort_by);
    if (sort_order) params.set("sort_order", sort_order);
    const res = await fetch(`${GATEWAY_URL}/api/v1/products/search?${params}`, { cache: "no-store" });
    const data = await res.json();
    return extractProducts(data);
  } catch { return { products: [], total: 0, page: 1, page_size: 20, total_pages: 0 }; }
}

export default async function SearchPage({ searchParams }: { searchParams: Promise<{ q?: string; sort_by?: string; sort_order?: string }> }) {
  const sp = await searchParams;
  const q = sp.q || "";
  const { products } = q ? await getSearchResults(q, sp.sort_by, sp.sort_order) : { products: [] };

  const sortOptions = [
    { label: "Phù hợp", value: "" },
    { label: "Mới nhất", value: "created_at" },
    { label: "Bán chạy", value: "sales" },
    { label: "Giá thấp → cao", value: "price_asc" },
    { label: "Giá cao → thấp", value: "price_desc" },
  ];

  const currentSort = sortOptions.find(o => o.value === (sp.sort_by || "")) || sortOptions[0];

  function buildSortUrl(val: string) {
    const p = new URLSearchParams();
    if (q) p.set("q", q);
    if (val === "price_asc") { p.set("sort_by", "price"); p.set("sort_order", "ASC"); }
    else if (val === "price_desc") { p.set("sort_by", "price"); p.set("sort_order", "DESC"); }
    else if (val) p.set("sort_by", val);
    return `/search?${p.toString()}`;
  }

  return (
    <>
      <Header />
      <main style={{ backgroundColor: "#F5F5FA" }}>
        <div className="max-w-tiki mx-auto">
          {/* Breadcrumb */}
          <div className="flex items-center h-[36px] text-xs">
            <Link href="/" className="text-tiki-text-secondary hover:text-tiki-blue hover:underline">Trang chủ</Link>
            <svg className="mx-[5px]" width="5" height="8" viewBox="0 0 5 8" fill="none"><path d="M1 1L4 4L1 7" stroke="#808089" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/></svg>
            <span className="text-tiki-text">{q ? `Tìm kiếm: ${q}` : "Tìm kiếm"}</span>
          </div>

          {q ? (
            <div className="flex gap-3">
              {/* Filter sidebar */}
              <aside className="w-[200px] shrink-0 hidden md:block">
                <div className="bg-white rounded-lg border border-tiki-border overflow-hidden sticky top-[60px]">
                  <div className="px-3 py-2 border-b border-tiki-border">
                    <span className="text-xs font-semibold text-tiki-text">BỘ LỌC</span>
                  </div>
                  <div className="px-3 py-2 border-b border-tiki-border">
                    <div className="text-[11px] font-semibold text-tiki-text mb-1.5">CHUYÊN TRANG</div>
                    <div className="flex flex-col gap-1">
                      <span className="text-[11px] text-tiki-text-secondary">Tiki Trading</span>
                      <span className="text-[11px] text-tiki-text-secondary">Tiki Global</span>
                      <span className="text-[11px] text-tiki-text-secondary">Nhà bán hàng</span>
                    </div>
                  </div>
                  <div className="px-3 py-2">
                    <div className="text-[11px] font-semibold text-tiki-text mb-1.5">GIÁ BÁN</div>
                    <div className="flex flex-col gap-1">
                      <span className="text-[11px] text-tiki-text-secondary">Dưới 500.000</span>
                      <span className="text-[11px] text-tiki-text-secondary">500.000 - 1.000.000</span>
                      <span className="text-[11px] text-tiki-text-secondary">1.000.000 - 3.000.000</span>
                      <span className="text-[11px] text-tiki-text-secondary">Trên 3.000.000</span>
                    </div>
                  </div>
                </div>
              </aside>

              {/* Main content */}
              <div className="flex-1 min-w-0">
                {/* Result count + sort bar */}
                <div className="bg-white rounded-lg border border-tiki-border mb-2">
                  <div className="px-3 py-2 flex items-center justify-between">
                    <span className="text-xs text-tiki-text-secondary">{products.length} kết quả cho "<strong className="text-tiki-text">{q}</strong>"</span>
                    <div className="flex items-center gap-1">
                      <span className="text-[11px] text-tiki-text-secondary mr-1">Sắp xếp:</span>
                      {sortOptions.map((opt) => (
                        <Link
                          key={opt.value}
                          href={buildSortUrl(opt.value)}
                          className={`px-2 py-1 text-[11px] rounded transition ${
                            currentSort.value === opt.value
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

                {/* Product grid */}
                <ProductGrid products={products as Product[]} />
              </div>
            </div>
          ) : (
            <div className="bg-white rounded-lg p-12 text-center border border-tiki-border">
              <p className="text-tiki-text-secondary text-sm">Nhập từ khóa để tìm kiếm</p>
            </div>
          )}
        </div>
      </main>
      <Footer />
    </>
  );
}
