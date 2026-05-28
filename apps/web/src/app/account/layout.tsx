"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Header } from "@/components/layout/header/Header";
import { Footer } from "@/components/layout/footer/Footer";
import { useAuthStore } from "@/stores/auth";

const SIDEBAR_ITEMS = [
  { href: "/account", label: "Thông tin tài khoản", icon: "👤" },
  { href: "/account/orders", label: "Đơn hàng của tôi", icon: "📦" },
  { href: "/account/addresses", label: "Sổ địa chỉ", icon: "📍" },
  { href: "/account/wishlist", label: "Sản phẩm yêu thích", icon: "❤️" },
];

export default function AccountLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, router]);

  return (
    <>
      <Header />
      <main className="bg-tiki-bg py-4 min-h-[60vh]">
        <div className="max-w-tiki mx-auto px-3">
          <div className="flex gap-4">
            <aside className="w-56 shrink-0 hidden md:block">
              <div className="account-sidebar">
                <div className="text-xs font-semibold text-tiki-text mb-3">TÀI KHOẢN CỦA TÔI</div>
                <nav className="space-y-0.5">
                  {SIDEBAR_ITEMS.map((item) => (
                    <Link
                      key={item.href}
                      href={item.href}
                      className="account-sidebar__link"
                    >
                      <span className="text-sm">{item.icon}</span>
                      {item.label}
                    </Link>
                  ))}
                </nav>
              </div>
            </aside>
            <div className="flex-1 min-w-0">{children}</div>
          </div>
        </div>
      </main>
      <Footer />
    </>
  );
}
