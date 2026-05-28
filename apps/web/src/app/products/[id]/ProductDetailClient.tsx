"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { StarRating, Price } from "@/components/ui";
import { AddToCartButton } from "./AddToCartButton";
import { useCartStore } from "@/stores/cart";
import { useUIStore } from "@/stores/ui";
import ProductDetailWithReviews from "@/components/storefront/product/ProductReview";
import RelatedProducts from "@/components/storefront/product/RelatedProducts";

interface ProductImage {
  id: string;
  url: string;
  is_primary: boolean;
}

interface ProductDetail {
  id: string;
  name: string;
  description?: string;
  short_description?: string;
  image_url: string;
  images?: ProductImage[];
  price: number;
  original_price?: number | null;
  discount_percent?: number | null;
  stock: number;
  sold_count: number;
  quantity_sold_text?: string;
  rating_average?: number | null;
  rating_count?: number;
  review_count?: number;
  brand?: string;
  seller_name?: string;
  seller_avatar_url?: string;
  is_official?: boolean;
  attributes?: { name: string; value: string }[];
  category_name?: string;
  category_id: string;
  weight?: number;
  dimensions?: string;
  status: string;
  shop_id?: string;
  shop_name?: string;
}

export default function ProductDetailClient({
  product,
  allImages,
}: {
  product: ProductDetail;
  allImages: ProductImage[];
}) {
  const router = useRouter();
  const [selectedImage, setSelectedImage] = useState(0);
  const [quantity, setQuantity] = useState(1);
  const addItem = useCartStore((s) => s.addItem);
  const addToast = useUIStore((s) => s.addToast);
  const [isBuying, setIsBuying] = useState(false);

  const currentImage = allImages[selectedImage]?.url || product.image_url;

  const handleBuyNow = async () => {
    setIsBuying(true);
    try {
      await addItem({
        product_id: product.id, sku_id: product.id, name: product.name,
        image_url: product.image_url, price: product.price, quantity,
        stock: product.stock || 0, shop_id: product.shop_id, shop_name: product.shop_name,
      });
      router.push("/checkout");
    } catch {
      addToast({ type: "error", title: "Có lỗi xảy ra", message: "Không thể thêm sản phẩm" });
    } finally {
      setIsBuying(false);
    }
  };

  return (
    <>
      <div className="bg-white rounded-lg border border-tiki-border p-4 md:p-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-8">
          <div className="flex flex-col gap-3">
            <div className="aspect-square bg-gray-50 rounded-lg overflow-hidden flex items-center justify-center">
              <img
                src={currentImage}
                alt={product.name}
                className="w-full h-full object-contain"
              />
            </div>
            {allImages.length > 1 && (
              <div className="flex gap-2 overflow-x-auto">
                {allImages.map((img, idx) => (
                  <button
                    key={img.id}
                    onClick={() => setSelectedImage(idx)}
                    className={`w-14 h-14 shrink-0 rounded border-2 overflow-hidden ${
                      idx === selectedImage ? "border-tiki-blue" : "border-gray-200"
                    }`}
                  >
                    <img src={img.url} alt="" className="w-full h-full object-cover" />
                  </button>
                ))}
              </div>
            )}
          </div>

          <div>
            <h1 className="text-lg font-semibold text-tiki-text mb-2">{product.name}</h1>

            <div className="flex items-center gap-3 text-xs text-tiki-text-secondary mb-3">
              {product.rating_average != null && product.rating_average > 0 && (
                <span className="flex items-center gap-1">
                  <StarRating rating={product.rating_average} size="sm" />
                  <span>{product.rating_average.toFixed(1)}</span>
                </span>
              )}
              {product.sold_count > 0 && (
                <span>Đã bán {product.quantity_sold_text || product.sold_count.toLocaleString("vi-VN")}</span>
              )}
              {product.review_count != null && product.review_count > 0 && (
                <span>{product.review_count} đánh giá</span>
              )}
            </div>

            <div className="bg-red-50 -mx-6 px-6 py-3 rounded">
              <Price
                amount={product.price}
                originalAmount={product.original_price}
                discountPercent={product.discount_percent}
                size="lg"
              />
            </div>

            <div className="mt-4 space-y-3">
              {product.brand && (
                <div className="flex text-sm">
                  <span className="text-tiki-text-secondary w-20 shrink-0">Thương hiệu</span>
                  <span className="text-tiki-text">{product.brand}</span>
                </div>
              )}

              <div className="flex items-center gap-2 text-sm">
                <span className="text-tiki-text-secondary">Số lượng</span>
                <div className="flex items-center border border-gray-300 rounded">
                  <button className="w-8 h-8 flex items-center justify-center text-sm text-tiki-text border-r border-gray-300 hover:bg-gray-50" onClick={() => setQuantity(Math.max(1, quantity - 1))}>−</button>
                  <input className="w-10 h-8 text-center text-xs border-none outline-none" type="text" value={quantity} readOnly />
                  <button className="w-8 h-8 flex items-center justify-center text-sm text-tiki-text border-l border-gray-300 hover:bg-gray-50" onClick={() => setQuantity(Math.min(product.stock || 999, quantity + 1))}>+</button>
                </div>
                <span className="text-xs text-tiki-text-secondary">{product.stock} sản phẩm có sẵn</span>
              </div>
            </div>

            <div className="flex gap-3 mt-4">
              <button
                onClick={handleBuyNow}
                disabled={isBuying}
                className="flex-1 py-3 bg-tiki-red text-white rounded font-semibold text-sm hover:bg-red-600 transition disabled:opacity-50"
              >
                {isBuying ? "Đang xử lý..." : "Mua ngay"}
              </button>
              <AddToCartButton product={product} quantity={quantity} />
            </div>

            {product.seller_name && (
              <div className="mt-4 p-3 bg-gray-50 rounded-lg flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center text-sm font-bold text-tiki-blue shrink-0">
                  {product.seller_name.charAt(0)}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-1">
                    <span className="text-sm font-medium text-tiki-text truncate">{product.seller_name}</span>
                    {product.is_official && (
                      <span className="text-[10px] bg-tiki-blue text-white px-1 rounded shrink-0">Official</span>
                    )}
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      <ProductDetailWithReviews product={product}>
        {product.attributes && product.attributes.length > 0 && (
          <div className="mt-4 bg-white rounded-lg border border-tiki-border p-4">
            <h3 className="text-sm font-semibold text-tiki-text mb-3">Thông số sản phẩm</h3>
            <table className="w-full text-sm">
              <tbody>
                {product.attributes.map((attr, idx) => (
                  <tr key={idx} className={idx % 2 === 0 ? "bg-gray-50" : ""}>
                    <td className="py-2 px-3 text-tiki-text-secondary w-40">{attr.name}</td>
                    <td className="py-2 px-3 text-tiki-text">{attr.value}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </ProductDetailWithReviews>

      <RelatedProducts productId={product.id} />
    </>
  );
}
