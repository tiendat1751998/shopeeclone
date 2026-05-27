import Link from "next/link";
import { Header } from "@/components/layout/header/Header";
import { Footer } from "@/components/layout/footer/Footer";
import { AddToCartButton } from "./AddToCartButton";
import { BuyNowButton } from "./BuyNowButton";
import { mapSingleProduct } from "@/lib/api/mapper";
import type { Product } from "@/types";

const GATEWAY_URL = process.env.GATEWAY_URL || "http://gateway:8080";

type ProductDetail = Product & { attributes?: { name: string; value: string }[] };

async function getProduct(id: string): Promise<ProductDetail | null> {
  try {
    const res = await fetch(`${GATEWAY_URL}/api/v1/products/${id}`, { cache: "no-store" });
    const data = await res.json();
    return mapSingleProduct(data) as ProductDetail | null;
  } catch { return null; }
}

function StarIcons({ rating, size = "text-sm" }: { rating: number; size?: string }) {
  return (
    <div className="flex items-center">
      {[1, 2, 3, 4, 5].map((i) => (
        <svg key={i} width="12" height="12" viewBox="0 0 12 12" fill={i <= Math.round(rating) ? "#FDD835" : "#EBEBF0"}>
          <path d="M6 0L7.96 4.3L12 4.56L9 7.5L9.84 12L6 9.7L2.16 12L3 7.5L0 4.56L4.04 4.3L6 0Z" />
        </svg>
      ))}
    </div>
  );
}

export default async function ProductPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const p = await getProduct(id);

  if (!p) {
    return (
      <>
        <Header />
        <main className="py-16 text-center">
          <p className="text-5xl mb-4">😕</p>
          <h1 className="text-lg font-semibold text-tiki-text mb-2">Sản phẩm không tồn tại</h1>
          <Link href="/products" className="text-tiki-blue hover:underline text-sm">← Quay lại</Link>
        </main>
        <Footer />
      </>
    );
  }

  return (
    <>
      <Header />
      <main style={{ backgroundColor: "#F5F5FA" }}>
        <div className="max-w-[1270px] mx-auto px-[12px]">
          {/* Breadcrumb */}
          <div className="flex items-center h-9 text-xs">
            <Link href="/" className="text-tiki-text-secondary hover:text-tiki-blue hover:underline whitespace-nowrap">Trang chủ</Link>
            <svg className="mx-[5px]" width="5" height="8" viewBox="0 0 5 8" fill="none"><path d="M1 1L4 4L1 7" stroke="#808089" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/></svg>
            <Link href="/products" className="text-tiki-text-secondary hover:text-tiki-blue hover:underline whitespace-nowrap">{p.category_name || "Sản phẩm"}</Link>
            <svg className="mx-[5px]" width="5" height="8" viewBox="0 0 5 8" fill="none"><path d="M1 1L4 4L1 7" stroke="#808089" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/></svg>
            <span className="text-tiki-text truncate">{p.name}</span>
          </div>

          <div style={{ display: "grid", gridTemplateColumns: "1fr 340px", gap: "16px" }}>
            {/* Left column */}
            <div style={{ display: "grid", gridTemplateColumns: "100%", gap: "12px" }}>
              {/* Product main section */}
              <div className="bg-white rounded-lg p-3">
                <div style={{ display: "grid", gridTemplateColumns: "360px 1fr", gap: "16px", alignItems: "start" }}>
                  {/* Image gallery */}
                  <div className="flex gap-2">
                    <div className="flex flex-col gap-1.5">
                      <div className="w-[48px] h-[48px] p-[3px] rounded border-2 border-tiki-blue overflow-hidden cursor-pointer">
                        <img src={p.image_url || "/images/placeholder.svg"} alt="" className="w-full h-full object-cover rounded-sm" />
                      </div>
                    </div>
                    <div className="flex-1">
                      <div className="rounded-lg overflow-hidden border border-[#ebebf0]">
                        <img src={p.image_url || "/images/placeholder.svg"} alt={p.name} className="w-full aspect-square object-contain" />
                      </div>
                    </div>
                  </div>

                  {/* Product info */}
                  <div className="flex flex-col gap-3">
                    {/* Official badge + name */}
                    <div>
                      {p.is_official && (
                        <span className="inline-block bg-tiki-blue text-white text-[9px] font-semibold px-1.5 py-0.5 rounded mr-2 align-middle">
                          CHÍNH HÃNG
                        </span>
                      )}
                      <h1 className="text-base font-medium text-[#27272A] leading-relaxed">{p.name}</h1>
                    </div>

                    {/* Brand & rating row */}
                    <div className="flex items-center gap-3 flex-wrap">
                      {p.brand && (
                        <div className="flex items-center">
                          <span className="text-xs text-[#242424]">
                            Thương hiệu: <Link href="#" className="text-[#0d5cb6] text-xs">{p.brand}</Link>
                          </span>
                        </div>
                      )}
                      {p.rating_average && p.rating_average > 0 && (
                        <div className="flex items-center gap-1.5">
                          <span className="text-xs text-[#787878]">|</span>
                          <StarIcons rating={p.rating_average} />
                          <span className="text-xs text-[#787878]">{p.rating_average?.toFixed(1)}</span>
                          <span className="text-xs text-[#787878]">({p.rating_count || 0})</span>
                          <span className="text-xs text-[#787878]">|</span>
                          <span className="text-xs text-[#787878]">{p.quantity_sold_text || `Đã bán ${p.sold_count || 0}`}</span>
                        </div>
                      )}
                    </div>

                    {/* Price */}
                    <div className="bg-[#fafafa] rounded-lg p-3">
                      <div className="flex items-center gap-2">
                        <span className="text-xl font-semibold text-tiki-red">{p.price?.toLocaleString("vi-VN")} ₫</span>
                        {p.original_price && p.original_price > p.price && (
                          <>
                            <span className="text-xs text-tiki-text-secondary line-through">{p.original_price?.toLocaleString("vi-VN")} ₫</span>
                            <span className="text-[10px] font-medium bg-[#F5F5FA] px-1.5 py-0.5 rounded text-tiki-text-secondary">-{p.discount_percent}%</span>
                          </>
                        )}
                      </div>
                    </div>

                    {/* Short description */}
                    {p.short_description && (
                      <p className="text-xs text-tiki-text-secondary leading-relaxed">{p.short_description}</p>
                    )}

                    {/* Attributes */}
                    {p.attributes && p.attributes.length > 0 && (
                      <div className="bg-white rounded-lg border border-[#ebebf0]">
                        {p.attributes.slice(0, 5).map((attr, i) => (
                          <div key={i} className="flex items-center py-1.5 px-3 border-b border-[#ebebf0] last:border-b-0 text-xs">
                            <span className="text-tiki-text-secondary w-1/3">{attr.name}</span>
                            <span className="text-tiki-text font-medium">{attr.value}</span>
                          </div>
                        ))}
                      </div>
                    )}

                    {/* Action buttons */}
                    <div className="flex gap-2">
                      <AddToCartButton product={{ id: p.id, name: p.name, image_url: p.image_url, price: p.price, stock: p.stock }} />
                      <BuyNowButton product={{ id: p.id, name: p.name, image_url: p.image_url, price: p.price, stock: p.stock }} />
                    </div>
                  </div>
                </div>
              </div>

              {/* Description */}
              {p.description && (
                <div className="bg-white rounded-lg p-3">
                  <h2 className="text-sm font-semibold text-[#27272A] mb-2">Mô tả sản phẩm</h2>
                  <p className="text-xs text-tiki-text-secondary leading-relaxed whitespace-pre-line">{p.description}</p>
                </div>
              )}
            </div>

            {/* Right sidebar */}
            <div className="flex flex-col gap-3">
              {/* Delivery info */}
              <div className="bg-white rounded-lg p-3 border border-[#ebebf0]">
                <div className="flex items-center gap-2 mb-2">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="#00AB56"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg>
                  <span className="text-xs font-medium text-tiki-text">Miễn phí vận chuyển</span>
                </div>
                <div className="text-[10px] text-tiki-text-secondary mb-1">Giao đến <strong className="text-tiki-text underline">Q. Hoàn Kiếm, Hà Nội</strong></div>
                <div className="flex items-center gap-1 text-[10px] text-tiki-text-secondary mb-2">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="#808089"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
                  <span>Phí vận chuyển: <strong className="text-tiki-text">30.000 ₫ - 50.000 ₫</strong></span>
                </div>
                <div className="flex items-center gap-1">
                  <span className="px-1.5 py-0.5 bg-tiki-blue text-white rounded text-[8px] font-bold">TIKI NOW</span>
                  <span className="text-[10px] text-tiki-text-secondary">Giao nhanh 2h</span>
                </div>
              </div>

              {/* Seller info */}
              {p.seller_name && (
                <div className="bg-white rounded-lg p-3 border border-[#ebebf0]">
                  <h3 className="text-xs font-semibold text-tiki-text mb-2">Thông tin người bán</h3>
                  <div className="flex items-center gap-2">
                    <div className="w-8 h-8 rounded-full bg-tiki-blue flex items-center justify-center text-white text-xs font-bold shrink-0">
                      {p.seller_name.charAt(0)}
                    </div>
                    <div>
                      <div className="text-xs font-medium text-tiki-text">{p.seller_name}</div>
                      <div className="text-[10px] text-tiki-text-secondary">Online 15 phút trước</div>
                    </div>
                  </div>
                  <button className="w-full mt-2 py-1.5 border border-tiki-blue text-tiki-blue rounded text-xs font-medium hover:bg-blue-50 transition">
                    Chat với người bán
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </>
  );
}
