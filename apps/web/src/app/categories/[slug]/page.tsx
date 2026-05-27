import Link from "next/link";
import { Header } from "@/components/layout/header/Header";
import { Footer } from "@/components/layout/footer/Footer";
import { ProductGrid } from "@/components/storefront/product/ProductCard";
import { extractProducts } from "@/lib/api/mapper";
import { Product, Category } from "@/types";
import categoriesData from "@/data/tiki-categories.json";

const GATEWAY_URL = process.env.GATEWAY_URL || "http://gateway:8080";

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

async function getCategoryProducts(slug: string) {
  try {
    const res = await fetch(`${GATEWAY_URL}/api/v1/products?category_slug=${slug}`, { cache: "no-store" });
    const data = await res.json();
    return extractProducts(data);
  } catch { return { products: [], total: 0, page: 1, page_size: 20, total_pages: 0 }; }
}

export default async function CategoryPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const category = findCategory(categoriesData, slug);
  const categoryName = category?.name || slug.replace(/-/g, " ");
  const { products } = await getCategoryProducts(slug);

  const sortOptions = [
    { label: "Phù hợp", value: "" },
    { label: "Mới nhất", value: "created_at" },
    { label: "Bán chạy", value: "sales" },
    { label: "Giá thấp → cao", value: "price_asc" },
    { label: "Giá cao → thấp", value: "price_desc" },
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
            <span className="text-tiki-text">{categoryName}</span>
          </div>

          <div className="flex gap-3">
            {/* Sidebar filters */}
            <aside className="w-[200px] shrink-0 hidden md:block">
              <div className="bg-white rounded-lg border border-tiki-border overflow-hidden sticky top-[60px]">
                {/* Sub-categories */}
                {category?.children && category.children.length > 0 && (
                  <div>
                    <div className="px-3 py-2 border-b border-tiki-border">
                      <span className="text-[11px] font-semibold text-tiki-text">DANH MỤC</span>
                    </div>
                    <div className="px-3 py-2 space-y-1.5">
                      {(category.children as Category[]).map((child) => (
                        <Link
                          key={child.id}
                          href={`/categories/${child.slug}`}
                          className="block text-[11px] text-tiki-text-secondary hover:text-tiki-blue transition"
                        >
                          {child.name}
                        </Link>
                      ))}
                    </div>
                  </div>
                )}

                {/* Price filter */}
                <div className="border-t border-tiki-border">
                  <div className="px-3 py-2 border-b border-tiki-border">
                    <span className="text-[11px] font-semibold text-tiki-text">GIÁ BÁN</span>
                  </div>
                  <div className="px-3 py-2 space-y-1.5">
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

            {/* Main content */}
            <div className="flex-1 min-w-0">
              {/* Header + sort */}
              <div className="bg-white rounded-lg border border-tiki-border mb-2">
                <div className="px-3 py-2 flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <h1 className="text-sm font-semibold text-tiki-text">{categoryName}</h1>
                    <span className="text-[11px] text-tiki-text-secondary">({products.length} sản phẩm)</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <span className="text-[10px] text-tiki-text-secondary mr-1">Sắp xếp:</span>
                    {sortOptions.map((opt) => (
                      <Link
                        key={opt.value}
                        href={`/categories/${slug}${opt.value ? `?sort_by=${opt.value === "price_asc" ? "price" : opt.value === "price_desc" ? "price" : opt.value}&sort_order=${opt.value === "price_asc" ? "ASC" : opt.value === "price_desc" ? "DESC" : "DESC"}` : ""}`}
                        className={`px-2 py-1 text-[10px] rounded transition ${
                          !opt.value ? "bg-tiki-blue text-white" : "text-tiki-text-secondary hover:bg-gray-50"
                        }`}
                      >
                        {opt.label}
                      </Link>
                    ))}
                  </div>
                </div>
              </div>

              {/* Sub-category chips */}
              {category?.children && category.children.length > 0 && (
                <div className="flex flex-wrap gap-1.5 mb-2">
                  {(category.children as Category[]).map((child) => (
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

              {/* Product grid */}
              <ProductGrid products={products as Product[]} />
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </>
  );
}
