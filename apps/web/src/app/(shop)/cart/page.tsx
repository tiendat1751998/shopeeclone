"use client";
import { useEffect } from "react";
import Image from "next/image";
import Link from "next/link";
import { useCartStore } from "@/lib/store/cart";
import { Button } from "@/components/ui/Button";
import { Price } from "@/components/ui/Price";

export default function CartPage() {
  const { items, isLoading, fetchCart, updateQuantity, removeItem, toggleSelect, toggleSelectAll, selectedItems, subtotal } = useCartStore();
  useEffect(() => { fetchCart(); }, [fetchCart]);
  const selected = selectedItems();

  if (isLoading) return <div className="container py-6"><h1 className="text-xl font-semibold mb-6">Shopping Cart</h1><div className="space-y-4">{[...Array(3)].map((_, i) => <div key={i} className="skeleton h-24 rounded-lg" />)}</div></div>;
  if (items.length === 0) return <div className="container py-16 text-center"><h2 className="text-xl font-semibold mb-2">Your cart is empty</h2><Link href="/products"><Button variant="primary">Start Shopping</Button></Link></div>;

  return (
    <div className="container py-6">
      <h1 className="text-xl font-semibold mb-6">Shopping Cart ({items.length} items)</h1>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-4">
          <div className="card p-3 flex items-center gap-3">
            <input type="checkbox" checked={items.every((i) => i.is_selected)} onChange={toggleSelectAll} className="w-4 h-4 accent-[#ee4d2d]" />
            <span className="text-sm">Select All</span>
          </div>
          {items.map((item) => (
            <div key={item.id} className="card p-4 flex gap-4">
              <input type="checkbox" checked={item.is_selected} onChange={() => toggleSelect(item.id)} className="w-4 h-4 accent-[#ee4d2d] mt-2" />
              <div className="w-20 h-20 flex-shrink-0 rounded overflow-hidden bg-gray-100">
                <Image src={item.image_url || "/images/placeholder.png"} alt={item.name} width={80} height={80} className="w-full h-full object-cover" />
              </div>
              <div className="flex-1 min-w-0">
                <Link href={`/products/${item.product_id}`} className="text-sm text-[#222] hover:text-[#ee4d2d] line-clamp-2">{item.name}</Link>
                {item.sku_name && <p className="text-xs text-[#757575] mt-0.5">Variation: {item.sku_name}</p>}
                <Price amount={item.price} size="sm" className="mt-1" />
              </div>
              <div className="flex flex-col items-end justify-between">
                <button onClick={() => removeItem(item.id)} className="text-xs text-[#757575] hover:text-red-500">Delete</button>
                <div className="flex items-center border border-[#e8e8e8] rounded">
                  <button onClick={() => updateQuantity(item.id, item.quantity - 1)} className="px-2 py-1 text-sm">−</button>
                  <span className="px-3 py-1 text-sm border-x border-[#e8e8e8] min-w-[40px] text-center">{item.quantity}</span>
                  <button onClick={() => updateQuantity(item.id, item.quantity + 1)} className="px-2 py-1 text-sm" disabled={item.quantity >= item.stock}>+</button>
                </div>
              </div>
            </div>
          ))}
        </div>
        <div className="card p-4 h-fit sticky top-24">
          <h3 className="font-semibold mb-4">Order Summary</h3>
          <div className="space-y-2 text-sm">
            <div className="flex justify-between"><span className="text-[#757575]">Subtotal ({selected.length})</span><span>S${subtotal().toFixed(2)}</span></div>
            <div className="flex justify-between"><span className="text-[#757575]">Shipping</span><span className="text-[#00bfa5]">Free</span></div>
            <div className="border-t border-[#e8e8e8] pt-2 flex justify-between font-semibold"><span>Total</span><Price amount={subtotal()} size="md" /></div>
          </div>
          <Link href="/checkout"><Button variant="primary" fullWidth className="mt-4" disabled={selected.length === 0}>Checkout ({selected.length})</Button></Link>
        </div>
      </div>
    </div>
  );
}
