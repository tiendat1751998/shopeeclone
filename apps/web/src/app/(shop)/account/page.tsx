"use client";
import { useEffect, useState } from "react";
import Link from "next/link";
import { useAuthStore } from "@/lib/store/auth";
import { ordersApi } from "@/lib/api/orders";
import { Button } from "@/components/ui/Button";
import type { Order } from "@/lib/types";

export default function AccountPage() {
  const { user, isAuthenticated, logout, fetchProfile } = useAuthStore();
  const [orders, setOrders] = useState<Order[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated) return;
    fetchProfile();
    ordersApi.list(1, 10).then((res) => setOrders(res.data)).catch(() => {}).finally(() => setIsLoading(false));
  }, [isAuthenticated, fetchProfile]);

  if (!isAuthenticated) return <div className="container py-16 text-center"><h2 className="text-xl font-semibold mb-4">Please login</h2><Link href="/login"><Button variant="primary">Login</Button></Link></div>;

  return (
    <div className="container py-6">
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="card p-4 h-fit">
          <div className="flex items-center gap-3 mb-4 pb-4 border-b border-[#e8e8e8]">
            <div className="w-12 h-12 rounded-full bg-[#ee4d2d] flex items-center justify-center text-white text-lg font-bold">{user?.display_name?.charAt(0)?.toUpperCase() || "U"}</div>
            <div><p className="font-medium text-sm">{user?.display_name}</p><p className="text-xs text-[#757575]">{user?.email}</p></div>
          </div>
          <nav className="space-y-1">
            <Link href="/account" className="block px-3 py-2 text-sm rounded bg-[#fff0ed] text-[#ee4d2d] font-medium">My Orders</Link>
            <button onClick={logout} className="block w-full text-left px-3 py-2 text-sm rounded text-red-500 hover:bg-red-50">Logout</button>
          </nav>
        </div>
        <div className="md:col-span-3">
          <h2 className="text-lg font-semibold mb-4">My Orders</h2>
          {isLoading ? <div className="space-y-3">{[...Array(3)].map((_, i) => <div key={i} className="skeleton h-20 rounded-lg" />)}</div>
          : orders.length === 0 ? <div className="card p-8 text-center"><p className="text-[#757575] mb-4">No orders yet.</p><Link href="/products"><Button variant="primary">Start Shopping</Button></Link></div>
          : <div className="space-y-4">{orders.map((o) => (
            <div key={o.id} className="card p-4 flex items-center justify-between">
              <div><p className="text-sm font-medium">{o.order_number}</p><p className="text-xs text-[#757575]">{new Date(o.created_at).toLocaleDateString()}</p></div>
              <span className="text-sm font-semibold text-[#ee4d2d]">S${o.total.toFixed(2)}</span>
            </div>
          ))}</div>}
        </div>
      </div>
    </div>
  );
}
