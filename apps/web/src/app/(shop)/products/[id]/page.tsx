"use client";
import { useState, useEffect, useRef } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { Price } from "@/components/ui/Price";
import { Badge } from "@/components/ui/Badge";
import { RatingStars } from "@/components/ui/RatingStars";
import { Button } from "@/components/ui/Button";
import { Skeleton } from "@/components/ui/Skeleton";
import { useCartStore } from "@/lib/store/cart";
import { productsApi } from "@/lib/api/products";
import type { Product, SKU } from "@/lib/types";

function sanitizeHTML(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#x27;");
}

function safeQuantity(value: number, min: number, max: number): number {
  if (!Number.isFinite(value)) return min;
  return Math.max(min, Math.min(max, Math.floor(value)));
}

export default function ProductDetailPage({ params: paramsPromise }: { params: Promise<{ id: string }> }) {
  const router = useRouter();
  const [product, setProduct] = useState<Product | null>(null);
  const [skus, setSkus] = useState<SKU[]>([]);
  const [selectedSKU, setSelectedSKU] = useState<SKU | null>(null);
  const [selectedImage, setSelectedImage] = useState(0);
  const [quantity, setQuantity] = useState(1);
  const [isLoading, setIsLoading] = useState(true);
  const [isAdding, setIsAdding] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [id, setId] = useState<string>("");
  const addItem = useCartStore((s) => s.addItem);
  const abortRef = useRef<AbortController | null>(null);

  useEffect(() => {
    paramsPromise.then((p) => setId(p.id));
  }, [paramsPromise]);

  useEffect(() => {
    if (!id) return;
    if (abortRef.current) abortRef.current.abort();
    const controller = new AbortController();
    abortRef.current = controller;
    setIsLoading(true);
    setError(null);

    productsApi.getById(id, controller.signal)
      .then((p) => {
        if (!controller.signal.aborted) {
          setProduct(p);
          const productSkus = p.skus || [];
          setSkus(productSkus);
          if (productSkus.length > 0) setSelectedSKU(productSkus[0]);
          setIsLoading(false);
        }
      })
      .catch((e: Error) => {
        if (!controller.signal.aborted) {
          setError(e.message);
          setIsLoading(false);
        }
      });

    return () => controller.abort();
  }, [id]);

  const handleAddToCart = async () => {
    if (!product || !selectedSKU) return;
    setIsAdding(true);
    try {
      await addItem(product.id, selectedSKU.id, quantity, selectedSKU.name, selectedSKU.price, product.shop?.id || "", product.shop?.name || "", images[0]?.url || "");
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to add to cart");
    } finally {
      setIsAdding(false);
    }
  };

  const handleBuyNow = async () => {
    if (!product || !selectedSKU) return;
    setIsAdding(true);
    try {
      await addItem(product.id, selectedSKU.id, 1, selectedSKU.name, selectedSKU.price, product.shop?.id || "", product.shop?.name || "", images[0]?.url || "");
      router.push("/checkout");
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to process");
    } finally {
      setIsAdding(false);
    }
  };

  if (isLoading) {
    return (
      <div className="container py-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          <Skeleton variant="rectangular" className="aspect-square rounded-lg" />
          <div className="space-y-4">
            <Skeleton className="h-8 w-full" /><Skeleton className="h-6 w-1/3" />
            <Skeleton className="h-10 w-1/2" /><Skeleton className="h-32 w-full" />
          </div>
        </div>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="container py-16 text-center">
        <h2 className="text-xl font-semibold mb-2">{error || "Product not found"}</h2>
      </div>
    );
  }

  const images = product.media?.filter((m) => m.type === "image" && m.status === "active") || [];
  const currentImage = images[selectedImage]?.url || "/images/placeholder.svg";
  const displayPrice = selectedSKU?.price || skus[0]?.price || 0;
  const displayComparePrice = selectedSKU?.compare_price || skus[0]?.compare_price;
  const availableStock = selectedSKU ? Math.max(0, selectedSKU.stock - selectedSKU.reserved_stock) : 0;

  return (
    <div className="container py-6">
      {error && <div className="bg-red-50 text-red-600 text-sm p-3 rounded mb-4">{error}</div>}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div className="space-y-3">
          <div className="card overflow-hidden"><div className="relative aspect-square">
            <Image src={currentImage} alt={sanitizeHTML(product.name)} fill className="object-cover" priority />
          </div></div>
          {images.length > 1 && (
            <div className="flex gap-2 overflow-x-auto">
              {images.map((img, i) => (
                <button key={img.id} onClick={() => setSelectedImage(i)}
                  className={`flex-shrink-0 w-16 h-16 rounded border-2 overflow-hidden transition-colors ${i === selectedImage ? "border-[#ee4d2d]" : "border-transparent hover:border-[#e8e8e8]"}`}>
                  <img src={img.thumbnail_url || img.url} alt="" className="w-full h-full object-cover" />
                </button>
              ))}
            </div>
          )}
        </div>
        <div className="space-y-4">
          <div>
            {product.shop?.is_official && <Badge variant="official" className="mb-2">Official Shop</Badge>}
            <h1 className="text-xl font-semibold text-[#222] leading-tight">{sanitizeHTML(product.name)}</h1>
          </div>
          {product.rating && (
            <div className="flex items-center gap-4 text-sm">
              <div className="flex items-center gap-1"><span className="font-semibold text-[#ee4d2d]">{product.rating.average.toFixed(1)}</span><RatingStars rating={product.rating.average} showCount={false} /></div>
              <span className="text-[#757575]">|</span>
              <span className="text-[#757575]">{product.rating.count.toLocaleString()} Ratings</span>
            </div>
          )}
          <div className="bg-[#fafafa] p-4 rounded"><Price amount={displayPrice} originalAmount={displayComparePrice} size="lg" /></div>
          {skus.length > 0 && (
            <div>
              <label className="text-sm font-medium text-[#222] mb-2 block">Variation</label>
              <div className="flex flex-wrap gap-2">
                {skus.map((sku) => (
                  <button key={sku.id} onClick={() => { setSelectedSKU(sku); setQuantity(1); }}
                    className={`px-3 py-1.5 text-sm border rounded transition-colors ${selectedSKU?.id === sku.id ? "border-[#ee4d2d] bg-[#fff0ed] text-[#ee4d2d]" : "border-[#e8e8e8] hover:border-[#ee4d2d]"}`}>
                    {sanitizeHTML(sku.name)}
                  </button>
                ))}
              </div>
            </div>
          )}
          <div>
            <label className="text-sm font-medium text-[#222] mb-2 block">Quantity</label>
            <div className="flex items-center gap-3">
              <div className="flex items-center border border-[#e8e8e8] rounded">
                <button onClick={() => setQuantity((q) => safeQuantity(q - 1, 1, availableStock))} className="px-3 py-1.5 text-[#757575]">−</button>
                <input type="number" value={quantity}
                  onChange={(e) => setQuantity(safeQuantity(Number(e.target.value), 1, availableStock))}
                  className="w-12 text-center text-sm border-x border-[#e8e8e8] py-1.5" min={1} max={availableStock} />
                <button onClick={() => setQuantity((q) => safeQuantity(q + 1, 1, availableStock))} className="px-3 py-1.5 text-[#757575]" disabled={quantity >= availableStock}>+</button>
              </div>
              <span className="text-sm text-[#757575]">{availableStock} available</span>
            </div>
          </div>
          <div className="flex gap-3 pt-2">
            <Button variant="outline" fullWidth isLoading={isAdding} onClick={handleAddToCart} disabled={availableStock === 0}>Add to Cart</Button>
            <Button variant="primary" fullWidth disabled={availableStock === 0} onClick={handleBuyNow}>Buy Now</Button>
          </div>
          {product.description && (
            <div className="border-t border-[#e8e8e8] pt-4">
              <h3 className="font-semibold text-sm mb-2">Description</h3>
              <p className="text-sm text-[#757575] whitespace-pre-line">{sanitizeHTML(product.description)}</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
