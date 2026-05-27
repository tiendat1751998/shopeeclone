"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Header } from "@/components/layout/header/Header";
import { Footer } from "@/components/layout/footer/Footer";
import { useCartStore } from "@/stores/cart";
import { useAuthStore } from "@/stores/auth";
import { useUIStore } from "@/stores/ui";

export default function CheckoutPage() {
  const router = useRouter();
  const items = useCartStore((s) => s.items);
  const getSubtotal = useCartStore((s) => s.getSubtotal);
  const clearCart = useCartStore((s) => s.clearCart);
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const addToast = useUIStore((s) => s.addToast);

  const [name, setName] = useState("");
  const [phone, setPhone] = useState("");
  const [address, setAddress] = useState("");
  const [city, setCity] = useState("Hà Nội");
  const [paymentMethod, setPaymentMethod] = useState("cod");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const selectedItems = items.filter((i) => i.is_selected);
  const subtotal = getSubtotal();
  const shippingFee = subtotal >= 45000 ? 0 : 30000;
  const total = subtotal + shippingFee;

  if (items.length === 0) {
    return (
      <>
        <Header />
        <main className="py-16 text-center">
          <div className="max-w-md mx-auto px-6">
            <div className="text-5xl mb-4">🛒</div>
            <h1 className="text-lg font-semibold text-tiki-text mb-2">Giỏ hàng trống</h1>
            <p className="text-sm text-tiki-text-secondary mb-6">Hãy thêm sản phẩm trước khi thanh toán nhé!</p>
            <Link href="/products" className="inline-block px-6 py-2.5 bg-tiki-blue text-white rounded-lg font-semibold text-sm hover:bg-tiki-blue-dark transition">
              Mua sắm ngay
            </Link>
          </div>
        </main>
        <Footer />
      </>
    );
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!name.trim() || !phone.trim() || !address.trim()) {
      addToast({ type: "error", title: "Thiếu thông tin", message: "Vui lòng điền đầy đủ thông tin giao hàng" });
      return;
    }
    setIsSubmitting(true);
    try {
      // Simulate order creation
      await new Promise((r) => setTimeout(r, 1500));
      clearCart();
      addToast({ type: "success", title: "Đặt hàng thành công!", message: "Đơn hàng của bạn đã được tạo" });
      router.push("/");
    } catch {
      addToast({ type: "error", title: "Có lỗi xảy ra", message: "Không thể tạo đơn hàng, vui lòng thử lại" });
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <>
      <Header />
      <main className="py-4" style={{ backgroundColor: "#F5F5FA" }}>
        <div className="max-w-[1270px] mx-auto px-[15px]">
          <div className="flex items-center gap-2 text-xs text-tiki-text-secondary mb-4">
            <Link href="/" className="hover:text-tiki-blue">Trang chủ</Link>
            <span>/</span>
            <Link href="/cart" className="hover:text-tiki-blue">Giỏ hàng</Link>
            <span>/</span>
            <span className="text-tiki-text">Thanh toán</span>
          </div>

          {!isAuthenticated && (
            <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-4 flex items-center justify-between">
              <p className="text-sm text-yellow-800">
                💡 <Link href="/login" className="text-tiki-blue hover:underline font-medium">Đăng nhập</Link> để theo dõi đơn hàng và nhận ưu đãi
              </p>
            </div>
          )}

          <form onSubmit={handleSubmit}>
            <div className="flex gap-4">
              {/* Left: Shipping info */}
              <div className="flex-1 min-w-0 space-y-4">
                <div className="bg-white rounded-lg border border-tiki-border p-4">
                  <h2 className="text-base font-semibold text-tiki-text mb-4">Thông tin giao hàng</h2>
                  <div className="space-y-3">
                    <div>
                      <label className="block text-sm font-medium text-tiki-text mb-1">Họ và tên *</label>
                      <input type="text" value={name} onChange={(e) => setName(e.target.value)} required
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:border-tiki-blue"
                        placeholder="Nguyễn Văn A" />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-tiki-text mb-1">Số điện thoại *</label>
                      <input type="tel" value={phone} onChange={(e) => setPhone(e.target.value)} required
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:border-tiki-blue"
                        placeholder="0912345678" />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-tiki-text mb-1">Địa chỉ *</label>
                      <input type="text" value={address} onChange={(e) => setAddress(e.target.value)} required
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:border-tiki-blue"
                        placeholder="Số nhà, đường, phường/xã" />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-tiki-text mb-1">Tỉnh/Thành phố</label>
                      <select value={city} onChange={(e) => setCity(e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:border-tiki-blue">
                        <option>Hà Nội</option>
                        <option>TP. Hồ Chí Minh</option>
                        <option>Đà Nẵng</option>
                        <option>Hải Phòng</option>
                        <option>Cần Thơ</option>
                        <option>Khác</option>
                      </select>
                    </div>
                  </div>
                </div>

                <div className="bg-white rounded-lg border border-tiki-border p-4">
                  <h2 className="text-base font-semibold text-tiki-text mb-4">Phương thức thanh toán</h2>
                  <div className="space-y-2">
                    {[
                      { value: "cod", label: "Thanh toán khi nhận hàng (COD)", icon: "💵" },
                      { value: "momo", label: "Ví MoMo", icon: "🟣" },
                      { value: "vnpay", label: "VNPay", icon: "🔵" },
                      { value: "bank", label: "Chuyển khoản ngân hàng", icon: "🏦" },
                    ].map((m) => (
                      <label key={m.value} className={`flex items-center gap-3 p-3 border rounded-lg cursor-pointer transition ${paymentMethod === m.value ? "border-tiki-blue bg-blue-50" : "border-tiki-border hover:border-tiki-blue"}`}>
                        <input type="radio" name="payment" value={m.value} checked={paymentMethod === m.value} onChange={() => setPaymentMethod(m.value)} className="w-4 h-4" />
                        <span className="text-lg">{m.icon}</span>
                        <span className="text-sm text-tiki-text">{m.label}</span>
                      </label>
                    ))}
                  </div>
                </div>
              </div>

              {/* Right: Order summary */}
              <div className="w-[380px] shrink-0">
                <div className="bg-white rounded-lg border border-tiki-border p-4 sticky top-4">
                  <h2 className="text-base font-semibold text-tiki-text mb-3">Đơn hàng ({selectedItems.length} sản phẩm)</h2>

                  <div className="max-h-[300px] overflow-y-auto space-y-3 mb-4">
                    {selectedItems.map((item) => (
                      <div key={item.id} className="flex items-center gap-3">
                        <img src={item.image_url || "/images/placeholder.svg"} alt={item.name} className="w-12 h-12 object-cover rounded border border-tiki-border shrink-0" />
                        <div className="flex-1 min-w-0">
                          <p className="text-sm text-tiki-text truncate">{item.name}</p>
                          <p className="text-xs text-tiki-text-secondary">x{item.quantity}</p>
                        </div>
                        <span className="text-sm font-medium text-tiki-text">{(item.price * item.quantity).toLocaleString("vi-VN")} ₫</span>
                      </div>
                    ))}
                  </div>

                  <div className="border-t border-tiki-border pt-3 space-y-2">
                    <div className="flex justify-between text-sm">
                      <span className="text-tiki-text-secondary">Tạm tính</span>
                      <span className="text-tiki-text">{subtotal.toLocaleString("vi-VN")} ₫</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-tiki-text-secondary">Phí vận chuyển</span>
                      <span className={shippingFee === 0 ? "text-tiki-green" : "text-tiki-text"}>
                        {shippingFee === 0 ? "Miễn phí" : `${shippingFee.toLocaleString("vi-VN")} ₫`}
                      </span>
                    </div>
                    {shippingFee > 0 && (
                      <p className="text-[11px] text-tiki-text-secondary">Mua thêm {(45000 - subtotal).toLocaleString("vi-VN")} ₫ để được miễn phí ship</p>
                    )}
                    <div className="border-t border-tiki-border pt-2">
                      <div className="flex justify-between">
                        <span className="text-sm font-medium text-tiki-text">Tổng cộng</span>
                        <span className="text-lg font-semibold text-tiki-red">{total.toLocaleString("vi-VN")} ₫</span>
                      </div>
                    </div>
                  </div>

                  <button
                    type="submit"
                    disabled={isSubmitting || selectedItems.length === 0}
                    className="w-full mt-4 py-3 bg-tiki-red text-white rounded-lg font-semibold text-sm hover:bg-tiki-red-dark transition disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {isSubmitting ? "Đang xử lý..." : `Đặt hàng — ${total.toLocaleString("vi-VN")} ₫`}
                  </button>

                  <Link href="/cart" className="block text-center mt-2 text-sm text-tiki-blue hover:underline">
                    ← Quay lại giỏ hàng
                  </Link>
                </div>
              </div>
            </div>
          </form>
        </div>
      </main>
      <Footer />
    </>
  );
}
