"use client";

import { useState, useMemo } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useCartStore } from "@/stores/cart";
import { useAuthStore } from "@/stores/auth";
import { useMounted } from "@/hooks/useMounted";
import categoriesData from "@/data/tiki-categories.json";

// Get first 7 categories from the data
const categoryItems = (categoriesData as any[])[0]?.children?.slice(0, 7) || [];

function SearchIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M15.7832 14.1911L19.2708 17.6787C19.5725 17.9804 19.5725 18.4679 19.2708 18.7696C18.9692 19.0712 18.4816 19.0712 18.18 18.7696L14.6923 15.2819C13.3012 17.0504 11.2383 18.0362 8.98366 17.9881C5.01229 17.9008 1.83319 14.6286 1.83319 10.6558C1.83319 6.5924 5.12652 3.29907 9.18994 3.29907C13.1628 3.29907 16.4335 6.47818 16.5223 10.4495C16.5704 12.7042 15.5846 14.7671 13.8161 16.1582L15.7832 14.1911ZM9.18994 4.82073C6.00087 4.82073 3.35485 7.46525 3.35485 10.6558C3.35485 13.8449 6.00087 16.4909 9.18994 16.4909C12.3805 16.4909 15.0265 13.8449 15.0265 10.6558C15.0265 7.46525 12.3805 4.82073 9.18994 4.82073Z" fill="#0A68FF"/>
    </svg>
  );
}

function MenuIcon() {
  return (
    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <line x1="3" y1="6" x2="21" y2="6" />
      <line x1="3" y1="12" x2="21" y2="12" />
      <line x1="3" y1="18" x2="21" y2="18" />
    </svg>
  );
}

function CloseIcon() {
  return (
    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <line x1="18" y1="6" x2="6" y2="18" />
      <line x1="6" y1="6" x2="18" y2="18" />
    </svg>
  );
}

export function Header() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState("");
  const [userDropdownOpen, setUserDropdownOpen] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const mounted = useMounted();

  // Cart store
  const cartItems = useCartStore((s) => s.items);
  const cartCount = useMemo(
    () => (cartItems ?? []).reduce((sum, item) => sum + item.quantity, 0),
    [cartItems]
  );

  // Auth store — only read after hydration to avoid mismatch
  const user = useAuthStore((s) => s.user);
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const logout = useAuthStore((s) => s.logout);

  // SSR-safe defaults: use guest values until client hydrates
  const showAuth = mounted && isAuthenticated && user;
  const showCartBadge = mounted && cartCount > 0;
  const showGuestCartNotice = mounted && !isAuthenticated && cartCount > 0;

  function handleSearch(e: React.FormEvent) {
    e.preventDefault();
    const q = searchQuery.trim();
    if (q) {
      router.push(`/search?q=${encodeURIComponent(q)}`);
    }
  }

  async function handleLogout() {
    await logout();
    setUserDropdownOpen(false);
    router.push("/");
  }

  return (
    <header>
      {/* Top banner */}
      <div className="bg-[#EFFFF4] text-[#00AB56] text-xs py-1.5 text-center font-medium">
        <div className="max-w-tiki mx-auto px-6">
          Freeship đơn từ 45k, giảm nhiều hơn cùng <strong>TikiCARD</strong>
        </div>
      </div>

      {/* Main header */}
      <div className="bg-white border-b border-tiki-border">
        <div className="max-w-tiki mx-auto px-6 py-2">
          <div className="flex items-center gap-6">
            {/* Logo */}
            <Link href="/" className="flex-shrink-0">
              <div className="flex flex-col items-center">
                <svg width="64" height="22" viewBox="0 0 64 22" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <rect width="64" height="22" rx="4" fill="#1A94FF"/>
                  <text x="10" y="16" fill="white" fontSize="12" fontWeight="700" fontFamily="Inter, sans-serif">Tiki</text>
                </svg>
                <span className="text-[8px] text-[#003EA1] font-semibold mt-0.5 tracking-tight">TỐT &amp; NHANH</span>
              </div>
            </Link>

            {/* Search */}
            <div className="flex-1 max-w-[560px]">
              <form onSubmit={handleSearch} className="flex items-center border border-[#DDDDE3] rounded-lg overflow-hidden focus-within:border-tiki-blue">
                <div className="pl-4 pr-2">
                  <SearchIcon />
                </div>
                <input
                  type="text"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder="Tìm sản phẩm, danh mục hay thương hiệu mong muốn ..."
                  className="flex-1 py-2 text-sm outline-none text-tiki-text placeholder:text-[#808089]"
                />
                <button
                  type="submit"
                  className="w-[92px] h-[38px] text-tiki-blue text-sm font-medium border-l border-[#DDDDE3] hover:bg-blue-50 transition"
                >
                  Tìm kiếm
                </button>
              </form>
            </div>

            {/* Delivery location */}
            <div className="hidden lg:flex items-center gap-1 text-xs text-tiki-text-secondary shrink-0">
              <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                <path d="M8 1C5.24 1 3 3.24 3 6c0 3.5 5 9 5 9s5-5.5 5-9c0-2.76-2.24-5-5-5zm0 7.5C6.62 8.5 5.5 7.38 5.5 6S6.62 3.5 8 3.5 10.5 4.62 10.5 6 9.38 8.5 8 8.5z" fill="#808089"/>
              </svg>
              <span className="text-tiki-text font-medium">Giao đến:</span>
              <span className="underline text-tiki-text truncate max-w-[140px]">Q. Hoàn Kiếm, Hà Nội</span>
            </div>

            {/* Mobile menu button */}
            <button
              type="button"
              className="lg:hidden p-2 text-tiki-text-secondary"
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              aria-label="Toggle menu"
            >
              {mobileMenuOpen ? <CloseIcon /> : <MenuIcon />}
            </button>

            {/* Action icons - shown on desktop, hidden on mobile unless menu open */}
            <div className={`flex items-center gap-1 shrink-0 ${mobileMenuOpen ? 'absolute top-full left-0 right-0 bg-white border-b border-tiki-border px-6 py-4 z-50 flex-row justify-end' : 'hidden lg:flex'}`}>
              {/* User */}
              {showAuth ? (
                <div className="relative">
                  <button
                    type="button"
                    onClick={() => setUserDropdownOpen(!userDropdownOpen)}
                    className="flex items-center gap-1 px-3 py-2 rounded-lg hover:bg-gray-100 cursor-pointer"
                  >
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                      <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z" fill="#808089"/>
                    </svg>
                    <span className="text-sm text-tiki-text-secondary max-w-[100px] truncate">{user.display_name}</span>
                  </button>
                  {/* Dropdown */}
                  {userDropdownOpen && (
                    <div className="absolute right-0 top-full mt-1 w-52 bg-white border border-tiki-border rounded-lg shadow-lg z-50">
                      <div className="px-4 py-3 border-b border-tiki-border">
                        <p className="text-sm font-medium text-tiki-text truncate">{user.display_name}</p>
                        <p className="text-xs text-tiki-text-secondary truncate">{user.email}</p>
                      </div>
                      <Link
                        href="/account"
                        onClick={() => setUserDropdownOpen(false)}
                        className="block px-4 py-2.5 text-sm text-tiki-text-secondary hover:bg-gray-50 transition"
                      >
                        Tài khoản của tôi
                      </Link>
                      <Link
                        href="/account/orders"
                        onClick={() => setUserDropdownOpen(false)}
                        className="block px-4 py-2.5 text-sm text-tiki-text-secondary hover:bg-gray-50 transition"
                      >
                        Đơn hàng của tôi
                      </Link>
                      <Link
                        href="/account/addresses"
                        onClick={() => setUserDropdownOpen(false)}
                        className="block px-4 py-2.5 text-sm text-tiki-text-secondary hover:bg-gray-50 transition"
                      >
                        Sổ địa chỉ
                      </Link>
                      <Link
                        href="/account/wishlist"
                        onClick={() => setUserDropdownOpen(false)}
                        className="block px-4 py-2.5 text-sm text-tiki-text-secondary hover:bg-gray-50 transition"
                      >
                        Sản phẩm yêu thích
                      </Link>
                      {user.role === "admin" && (
                        <Link
                          href="/admin"
                          onClick={() => setUserDropdownOpen(false)}
                          className="block px-4 py-2.5 text-sm text-tiki-blue hover:bg-gray-50 transition border-t border-tiki-border"
                        >
                          📊 Trang quản trị
                        </Link>
                      )}
                      <div className="border-t border-tiki-border">
                        <button
                          onClick={handleLogout}
                          className="w-full text-left px-4 py-2.5 text-sm text-red-600 hover:bg-gray-50 transition"
                        >
                          Đăng xuất
                        </button>
                      </div>
                    </div>
                  )}
                </div>
              ) : (
                <Link href="/login" className="flex items-center gap-1 px-3 py-2 rounded-lg hover:bg-gray-100">
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                    <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z" fill="#808089"/>
                  </svg>
                  <span className="text-sm text-tiki-text-secondary">Đăng nhập</span>
                </Link>
              )}

              {/* Cart */}
              <Link href="/cart" className="relative flex items-center gap-1 px-3 py-2 rounded-lg hover:bg-gray-100">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                  <path d="M7 18c-1.1 0-1.99.9-1.99 2S5.9 22 7 22s2-.9 2-2-.9-2-2-2zM1 2v2h2l3.6 7.59-1.35 2.45c-.16.28-.25.61-.25.96 0 1.1.9 2 2 2h12v-2H7.42c-.14 0-.25-.11-.25-.25l.03-.12.9-1.63h7.45c.75 0 1.41-.41 1.75-1.03l3.58-6.49c.08-.14.12-.31.12-.48 0-.55-.45-1-1-1H5.21l-.94-2H1zm16 16c-1.1 0-1.99.9-1.99 2s.89 2 1.99 2 2-.9 2-2-.9-2-2-2z" fill="#808089"/>
                </svg>
                {showCartBadge && (
                  <span className="absolute -top-0.5 right-1 bg-tiki-red text-white text-[10px] font-bold rounded-full min-w-[16px] h-4 flex items-center justify-center px-1">
                    {cartCount}
                  </span>
                )}
              </Link>
            </div>
          </div>
        </div>
      </div>

      {/* Category bar - shown on mobile via hamburger, always on desktop */}
      <div className={`bg-white border-b border-tiki-border ${mobileMenuOpen ? 'block' : 'hidden md:block'}`}>
        <div className="max-w-tiki mx-auto px-6 flex items-center gap-6 h-10 text-sm overflow-x-auto">
          <Link href="/" className="flex items-center gap-1.5 text-tiki-text font-medium hover:text-tiki-blue shrink-0">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
              <rect x="1" y="1" width="7.5" height="7.5" rx="1.5" stroke="#27272A" strokeWidth="1.5"/>
              <rect x="11.5" y="1" width="7.5" height="7.5" rx="1.5" stroke="#27272A" strokeWidth="1.5"/>
              <rect x="1" y="11.5" width="7.5" height="7.5" rx="1.5" stroke="#27272A" strokeWidth="1.5"/>
              <rect x="11.5" y="11.5" width="7.5" height="7.5" rx="1.5" stroke="#27272A" strokeWidth="1.5"/>
            </svg>
            Danh mục
          </Link>
          <div className="flex gap-5">
            {categoryItems.map((cat: any) => (
              <Link
                key={cat.id}
                href={`/categories/${cat.slug}`}
                className="text-tiki-text-secondary hover:text-tiki-blue whitespace-nowrap text-xs"
              >
                {cat.name}
              </Link>
            ))}
          </div>
        </div>
      </div>

      {/* Guest cart notice */}
      {showGuestCartNotice && (
        <div className="bg-yellow-50 border-b border-yellow-200 py-2">
          <div className="max-w-tiki mx-auto px-6 flex items-center justify-between text-xs">
            <span className="text-yellow-800">
              🛒 Bạn có {cartCount} sản phẩm trong giỏ hàng — <Link href="/login" className="text-tiki-blue hover:underline font-medium">Đăng nhập</Link> để đồng bộ
            </span>
            <Link href="/cart" className="text-tiki-blue hover:underline font-medium">
              Xem giỏ hàng →
            </Link>
          </div>
        </div>
      )}
    </header>
  );
}
