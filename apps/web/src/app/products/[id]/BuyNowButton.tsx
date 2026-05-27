"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useCartStore } from "@/stores/cart";
import { useUIStore } from "@/stores/ui";

export function BuyNowButton({ product }: { product: { id: string; name: string; image_url: string; price: number; stock?: number } }) {
  const [isBuying, setIsBuying] = useState(false);
  const addItem = useCartStore((s) => s.addItem);
  const clearCart = useCartStore((s) => s.clearCart);
  const addToast = useUIStore((s) => s.addToast);
  const router = useRouter();

  const handleBuyNow = async () => {
    setIsBuying(true);
    try {
      // Clear cart first, then add this item as the only item
      clearCart();
      await addItem({
        product_id: product.id,
        sku_id: "default",
        name: product.name,
        image_url: product.image_url,
        price: product.price,
        quantity: 1,
        stock: product.stock || 0,
        is_selected: true,
        sku_name: "default",
      });
      addToast({ type: "success", title: "Đã thêm vào giỏ hàng", message: product.name });
      router.push("/cart");
    } catch {
      addToast({ type: "error", title: "Có lỗi xảy ra", message: "Không thể thêm sản phẩm" });
    } finally {
      setIsBuying(false);
    }
  };

  return (
    <button
      onClick={handleBuyNow}
      disabled={isBuying}
      style={{
        flex: 1,
        height: "40px",
        background: "#FF424E",
        color: "white",
        fontSize: "16px",
        fontWeight: 400,
        borderRadius: "4px",
        border: "none",
        cursor: "pointer",
      }}
      className="disabled:opacity-50 hover:opacity-90 transition-opacity"
    >
      {isBuying ? "Đang xử lý..." : "Mua ngay"}
    </button>
  );
}
