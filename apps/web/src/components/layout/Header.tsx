"use client";
import Link from "next/link";
import { useState } from "react";
import { useCartStore } from "@/lib/store/cart";
import { useAuthStore } from "@/lib/store/auth";

export function Header() {
  const totalItems = useCartStore((s) => s.totalItems());
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const user = useAuthStore((s) => s.user);
  const [searchQuery, setSearchQuery] = useState("");
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <header className="sticky top-0 z-40 bg-gradient-to-r from-[#ee4d2d] to-[#f53d2d] shadow-md">
      {/* Top bar */}
      <div className="container">
        <div className="flex items-center justify-between h-8 text-xs text-white/80">
          <div className="flex items-center gap-4">
            <Link href="/seller" className="hover:text-white">Seller Centre</Link>
            <Link href="/download" className="hover:text-white">Download</Link>
            <span className="hidden sm:inline">Follow us on</span>
          </div>
          <div className="flex items-center gap-4">
            <Link href="/notifications" className="hover:text-white flex items-center gap-1">
              <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" /></svg>
              Notifications
            </Link>
            <Link href="/help" className="hover:text-white">Help</Link>
            {isAuthenticated ? (
              <Link href="/account" className="hover:text-white flex items-center gap-1">
                <div className="w-5 h-5 rounded-full bg-white/20 flex items-center justify-center text-xs font-bold">
                  {user?.display_name?.charAt(0)?.toUpperCase() || "U"}
                </div>
                <span className="hidden sm:inline">{user?.display_name || "Account"}</span>
              </Link>
            ) : (
              <div className="flex items-center gap-2">
                <Link href="/register" className="hover:text-white font-medium">Register</Link>
                <span>|</span>
                <Link href="/login" className="hover:text-white font-medium">Login</Link>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Main header */}
      <div className="container pb-3">
        <div className="flex items-center gap-4">
          {/* Logo */}
          <Link href="/" className="flex-shrink-0">
            <h1 className="text-2xl font-bold text-white tracking-tight">
              Shopee<span className="text-yellow-300">.</span>
            </h1>
          </Link>

          {/* Search bar */}
          <form className="flex-1 max-w-2xl hidden md:flex" action="/products" method="GET">
            <div className="flex w-full">
              <input
                type="text"
                name="q"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search for products, brands and shops"
                className="flex-1 px-4 py-2.5 text-sm rounded-l-sm bg-white focus:outline-none"
              />
              <button type="submit" className="px-6 py-2.5 bg-white border-l border-[#e8e8e8] rounded-r-sm hover:bg-gray-50 transition-colors">
                <svg className="w-5 h-5 text-[#ee4d2d]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
              </button>
            </div>
          </form>

          {/* Cart */}
          <Link href="/cart" className="relative flex-shrink-0 text-white hover:text-white/80 transition-colors">
            <svg className="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z" /></svg>
            {totalItems > 0 && (
              <span className="absolute -top-1.5 -right-1.5 bg-white text-[#ee4d2d] text-xs font-bold rounded-full min-w-[18px] h-[18px] flex items-center justify-center px-1">
                {totalItems > 99 ? "99+" : totalItems}
              </span>
            )}
          </Link>

          {/* Mobile menu toggle */}
          <button className="md:hidden text-white p-2" onClick={() => setMobileMenuOpen(!mobileMenuOpen)}>
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={mobileMenuOpen ? "M6 18L18 6M6 6l12 12" : "M4 6h16M4 12h16M4 18h16"} /></svg>
          </button>
        </div>

        {/* Mobile search */}
        {mobileMenuOpen && (
          <div className="md:hidden mt-3">
            <form className="flex" action="/products" method="GET">
              <input type="text" name="q" placeholder="Search..." className="flex-1 px-4 py-2 text-sm rounded-l-sm bg-white focus:outline-none" />
              <button type="submit" className="px-4 py-2 bg-white rounded-r-sm">
                <svg className="w-5 h-5 text-[#ee4d2d]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
              </button>
            </form>
          </div>
        )}
      </div>
    </header>
  );
}
