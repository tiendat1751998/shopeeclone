import Link from "next/link";
export function Header() {
  return (
    <header className="bg-white border-b border-gray-200 sticky top-0 z-50">
      <div className="container mx-auto px-4 h-14 flex items-center justify-between">
        <Link href="/" className="text-xl font-bold text-[#189eff]">Tiki</Link>
        <nav className="flex items-center gap-4 text-sm text-gray-600">
          <Link href="/products">Products</Link>
          <Link href="/categories">Categories</Link>
        </nav>
      </div>
    </header>
  );
}
