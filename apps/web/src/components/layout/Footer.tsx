import Link from "next/link";

export function Footer() {
  return (
    <footer className="bg-white border-t border-[#e8e8e8] mt-8">
      <div className="container py-10">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
          <div>
            <h4 className="font-semibold text-sm mb-4">CUSTOMER SERVICE</h4>
            <ul className="space-y-2 text-sm text-[#757575]">
              <li><Link href="/help" className="hover:text-[#ee4d2d]">Help Centre</Link></li>
              <li><Link href="/how-to-buy" className="hover:text-[#ee4d2d]">How to Buy</Link></li>
              <li><Link href="/how-to-sell" className="hover:text-[#ee4d2d]">How to Sell</Link></li>
              <li><Link href="/returns" className="hover:text-[#ee4d2d]">Returns & Refunds</Link></li>
            </ul>
          </div>
          <div>
            <h4 className="font-semibold text-sm mb-4">ABOUT SHOPEE</h4>
            <ul className="space-y-2 text-sm text-[#757575]">
              <li><Link href="/about" className="hover:text-[#ee4d2d]">About Us</Link></li>
              <li><Link href="/careers" className="hover:text-[#ee4d2d]">Careers</Link></li>
              <li><Link href="/privacy" className="hover:text-[#ee4d2d]">Privacy Policy</Link></li>
              <li><Link href="/terms" className="hover:text-[#ee4d2d]">Terms of Service</Link></li>
            </ul>
          </div>
          <div>
            <h4 className="font-semibold text-sm mb-4">PAYMENT</h4>
            <div className="flex flex-wrap gap-2">
              {["Visa", "Mastercard", "PayNow", "GrabPay", "ShopeePay"].map((p) => (
                <span key={p} className="px-2 py-1 bg-gray-100 rounded text-xs text-[#757575]">{p}</span>
              ))}
            </div>
          </div>
          <div>
            <h4 className="font-semibold text-sm mb-4">FOLLOW US</h4>
            <div className="flex gap-3">
              {["Facebook", "Instagram", "Twitter"].map((s) => (
                <Link key={s} href="#" className="w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center text-[#757575] hover:bg-[#ee4d2d] hover:text-white transition-colors">
                  <span className="text-xs font-bold">{s[0]}</span>
                </Link>
              ))}
            </div>
          </div>
        </div>
        <div className="mt-8 pt-6 border-t border-[#e8e8e8] text-center text-xs text-[#757575]">
          <p>&copy; 2025 Shopee Clone. All rights reserved. This is a demo project.</p>
        </div>
      </div>
    </footer>
  );
}
