import Link from "next/link";
import Image from "next/image";
import { Price } from "@/components/ui/Price";
import { Badge } from "@/components/ui/Badge";
import { RatingStars } from "@/components/ui/RatingStars";
import type { Product } from "@/lib/types";

interface ProductCardProps { product: Product; priority?: boolean; }

export function ProductCard({ product, priority = false }: ProductCardProps) {
  const mainImage = product.media?.find((m) => m.type === "image" && m.status === "active");
  const imageUrl = mainImage?.url || mainImage?.thumbnail_url || "/images/placeholder.svg";
  const lowestPrice = product.skus && product.skus.length > 0
    ? Math.min(...product.skus.map((s) => s.price))
    : 0;
  const hasDiscount = product.skus?.some((s) => s.compare_price > s.price);
  const discountPercent = hasDiscount && product.skus
    ? Math.round((1 - lowestPrice / (product.skus[0].compare_price || lowestPrice)) * 100)
    : 0;

  return (
    <Link href={`/products/${product.id}`} className="card group hover:shadow-md transition-shadow duration-200 block">
      <div className="relative aspect-square overflow-hidden">
        <Image
          src={imageUrl}
          alt={product.name}
          fill
          sizes="(max-width: 640px) 50vw, (max-width: 1024px) 33vw, 20vw"
          className="object-cover group-hover:scale-105 transition-transform duration-300"
          priority={priority}
        />
        {hasDiscount && discountPercent > 0 && <Badge variant="sale" className="absolute top-2 left-2">-{discountPercent}%</Badge>}
        {product.shop?.is_official && <Badge variant="official" className="absolute top-2 right-2">Official</Badge>}
      </div>
      <div className="p-3">
        <h3 className="text-sm text-[#222] truncate-2 leading-tight mb-1 group-hover:text-[#ee4d2d] transition-colors">
          {product.name}
        </h3>
        <div className="mb-1.5">
          {lowestPrice > 0 ? (
            <Price amount={lowestPrice} size="md" />
          ) : (
            <span className="text-sm text-[#757575]">Price unavailable</span>
          )}
        </div>
        <div className="flex items-center justify-between text-xs text-[#757575]">
          {product.rating ? (
            <RatingStars rating={product.rating.average} count={product.rating.count} size="sm" />
          ) : (
            <span className="text-[#bdbdbd]">No ratings yet</span>
          )}
          {product.sold_count !== undefined && product.sold_count > 0 && (
            <span>{product.sold_count > 1000 ? `${(product.sold_count / 1000).toFixed(1)}k` : product.sold_count} sold</span>
          )}
        </div>
        {product.shop && (
          <div className="mt-2 flex items-center gap-1.5 text-xs text-[#757575]">
            <div className="w-4 h-4 rounded-full bg-gray-200 flex items-center justify-center text-[10px] font-bold text-[#757575]">
              {product.shop.name.charAt(0)}
            </div>
            <span className="truncate">{product.shop.name}</span>
          </div>
        )}
      </div>
    </Link>
  );
}
