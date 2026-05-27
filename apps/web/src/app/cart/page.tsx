"use client";

import { useState } from "react";
import Link from "next/link";
import { Header } from "@/components/layout/header/Header";
import { Footer } from "@/components/layout/footer/Footer";
import { useCartStore } from "@/stores/cart";
import { useAddToCart } from "@/hooks/useApi";
import { useAuthStore } from "@/stores/auth";

export default function CartPage() {
  const items = useCartStore((s) => s.items);
  const removeItem = useCartStore((s) => s.removeItem);
  const updateQuantity = useCartStore((s) => s.updateQuantity);
  const toggleSelect = useCartStore((s) => s.toggleSelect);
  const selectAll = useCartStore((s) => s.selectAll);
  const clearCart = useCartStore((s) => s.clearCart);
  const getSubtotal = useCartStore((s) => s.getSubtotal);
  const selectedItems = items.filter((i) => i.is_selected);
  const allSelected = items.length > 0 && items.every((i) => i.is_selected);
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  if (items.length === 0) {
    return (
      <>
        <Header />
        <main className="py-4">
          <div className="max-w-tiki mx-auto px-6">
            <div className="flex items-center gap-2 text-xs text-tiki-text-secondary mb-4">
              <Link href="/" className="hover:text-tiki-blue">Trang chủ</Link>
              <span>/</span>
              <span className="text-tiki-text">Giỏ hàng</span>
            </div>
            <div className="bg-white rounded-lg p-12 text-center border border-tiki-border">
              <div className="text-5xl mb-4">🛒</div>
              <h1 className="text-base font-semibold text-tiki-text mb-2">Giỏ hàng trống</h1>
              <p className="text-sm text-tiki-text-secondary mb-6">Hãy thêm sản phẩm để mua sắm nhé!</p>
              <Link
                href="/products"
                className="inline-block px-6 py-2.5 bg-tiki-blue text-white rounded-lg font-semibold text-sm hover:bg-tiki-blue-dark transition"
              >
                Mua sắm ngay
              </Link>
            </div>
          </div>
        </main>
        <Footer />
      </>
    );
  }

  return (
    <>
      <Header />
      <main className="py-2" style={{ backgroundColor: "#F5F5FA" }}>
        <div className="max-w-[1270px] mx-auto px-[12px]">
          {/* Breadcrumb */}
          <div className="flex items-center h-9 text-xs">
            <Link href="/" className="text-tiki-text-secondary hover:text-tiki-blue hover:underline">Trang chủ</Link>
            <svg className="mx-[5px]" width="5" height="8" viewBox="0 0 5 8" fill="none"><path d="M1 1L4 4L1 7" stroke="#808089" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/></svg>
            <span className="text-tiki-text">Giỏ hàng</span>
          </div>

          {/* Guest user notice */}
          {!isAuthenticated && items.length > 0 && (
            <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-4 flex items-center justify-between">
              <div>
                <p className="text-sm text-yellow-800 font-medium">
                  Đăng nhập để đồng bộ giỏ hàng của bạn
                </p>
                <p className="text-xs text-yellow-600 mt-0.5">
                  Giỏ hàng hiện tại đang được lưu cục bộ. Đăng nhập để lưu trữ đám mây và truy cập trên mọi thiết bị.
                </p>
              </div>
              <Link
                href="/login?redirect=/cart"
                className="shrink-0 ml-4 px-4 py-2 bg-tiki-blue text-white rounded-lg text-sm font-medium hover:bg-tiki-blue-dark transition"
              >
                Đăng nhập
              </Link>
            </div>
          )}

          <div className="flex gap-4">
            {/* Cart items */}
            <div className="flex-1 min-w-0">
              <div className="bg-white rounded-lg border border-tiki-border overflow-hidden">
                {/* Header row */}
                <div className="flex items-center px-4 py-3 border-b border-tiki-border bg-gray-50">
                  <label className="flex items-center gap-2 text-sm text-tiki-text">
                    <input
                      type="checkbox"
                      checked={allSelected}
                      onChange={selectAll}
                      className="w-4 h-4 rounded border-gray-300"
                    />
                    Chọn tất cả ({items.length} sản phẩm)
                  </label>
                </div>

                {/* Item list */}
                <div className="divide-y divide-tiki-border">
                  {items.map((item) => (
                    <div key={item.id} className="flex items-center gap-4 px-4 py-3">
                      <input
                        type="checkbox"
                        checked={item.is_selected}
                        onChange={() => toggleSelect(item.id)}
                        className="w-4 h-4 rounded border-gray-300 shrink-0"
                      />
                      <img
                        src={item.image_url || "/images/placeholder.svg"}
                        alt={item.name}
                        className="w-16 h-16 object-cover rounded border border-tiki-border shrink-0"
                      />
                      <div className="flex-1 min-w-0">
                        <h3 className="text-sm text-tiki-text truncate">{item.name}</h3>
                        {item.sku_name && item.sku_name !== "default" && (
                          <p className="text-xs text-tiki-text-secondary mt-0.5">{item.sku_name}</p>
                        )}
                        <div className="flex items-center gap-2 mt-1">
                          <span className="text-sm font-semibold text-tiki-red">
                            {item.price?.toLocaleString("vi-VN")} ₫
                          </span>
                          {item.original_price && item.original_price > item.price && (
                            <span className="text-xs text-tiki-text-secondary line-through">
                              {item.original_price.toLocaleString("vi-VN")} ₫
                            </span>
                          )}
                        </div>
                      </div>
                      {/* Quantity controls */}
                      <div className="flex items-center gap-2 shrink-0">
                        <button
                          onClick={() => updateQuantity(item.id, Math.max(1, item.quantity - 1))}
                          className="w-8 h-8 flex items-center justify-center border border-tiki-border rounded text-sm hover:bg-gray-50"
                        >
                          −
                        </button>
                        <input
                          type="text"
                          value={item.quantity}
                          readOnly
                          className="w-10 h-8 text-center border border-tiki-border rounded text-sm"
                        />
                        <button
                          onClick={() => updateQuantity(item.id, Math.min(item.stock, item.quantity + 1))}
                          className="w-8 h-8 flex items-center justify-center border border-tiki-border rounded text-sm hover:bg-gray-50"
                        >
                          +
                        </button>
                      </div>
                      <button
                        onClick={() => removeItem(item.id)}
                        className="text-xs text-tiki-text-secondary hover:text-tiki-red shrink-0"
                      >
                        Xóa
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Sidebar */}
            <div className="w-[320px] shrink-0">
              <div className="bg-white rounded-lg border border-tiki-border p-4 sticky top-4">
                <h2 className="text-sm font-semibold text-tiki-text mb-3">Đơn hàng</h2>
                <div className="flex justify-between text-sm mb-2">
                  <span className="text-tiki-text-secondary">Tạm tính ({selectedItems.length} sản phẩm)</span>
                  <span className="text-tiki-text font-medium">{getSubtotal().toLocaleString("vi-VN")} ₫</span>
                </div>
                <div className="border-t border-tiki-border pt-3 mt-3">
                  <div className="flex justify-between text-sm">
                    <span className="text-tiki-text-secondary">Tổng tiền</span>
                    <span className="text-lg font-semibold text-tiki-red">{getSubtotal().toLocaleString("vi-VN")} ₫</span>
                  </div>
                  <p className="text-[11px] text-tiki-text-secondary mt-1">(Đã bao gồm VAT nếu có)</p>
                </div>
                <button
                  disabled={selectedItems.length === 0}
                  onClick={() => window.location.href = "/checkout"}
                  className="w-full mt-4 py-3 bg-tiki-red text-white rounded-lg font-semibold text-sm hover:bg-tiki-red-dark transition disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Mua hàng ({selectedItems.length})
                </button>
                <button
                  onClick={clearCart}
                  className="w-full mt-2 py-2 border border-tiki-border text-tiki-text-secondary rounded-lg text-sm hover:bg-gray-50 transition"
                >
                  Xóa giỏ hàng
                </button>
              </div>
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </>
  );
}
