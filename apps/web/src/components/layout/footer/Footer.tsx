import Link from "next/link";

export function Footer() {
  return (
    <footer className="bg-white border-t border-tiki-border mt-6">
      <div className="max-w-tiki mx-auto px-3">
        {/* Main columns */}
        <div className="grid grid-cols-2 md:grid-cols-5 gap-4 py-4">
          {/* Support */}
          <div>
            <h4 className="text-[11px] font-semibold text-tiki-text mb-2">Hỗ trợ khách hàng</h4>
            <ul className="space-y-1.5 text-[10px] text-tiki-text-secondary">
              <li>Hotline: <a href="tel:19006035" className="text-tiki-blue font-medium">1900-6035</a></li>
              <li className="text-[9px]">(1000 đ/phút, 8-21h kể cả T7, CN)</li>
              <li><Link href="#" className="hover:text-tiki-blue">Các câu hỏi thường gặp</Link></li>
              <li><Link href="#" className="hover:text-tiki-blue">Gửi yêu cầu hỗ trợ</Link></li>
              <li><Link href="#" className="hover:text-tiki-blue">Hướng dẫn đặt hàng</Link></li>
              <li><Link href="#" className="hover:text-tiki-blue">Chính sách đổi trả</Link></li>
            </ul>
          </div>

          {/* About */}
          <div>
            <h4 className="text-[11px] font-semibold text-tiki-text mb-2">Về Tiki</h4>
            <ul className="space-y-1.5 text-[10px] text-tiki-text-secondary">
              <li><Link href="#" className="hover:text-tiki-blue">Giới thiệu Tiki</Link></li>
              <li><Link href="#" className="hover:text-tiki-blue">Tiki Blog</Link></li>
              <li><Link href="#" className="hover:text-tiki-blue">Tuyển dụng</Link></li>
              <li><Link href="#" className="hover:text-tiki-blue">Chính sách bảo mật</Link></li>
              <li><Link href="#" className="hover:text-tiki-blue">Điều khoản sử dụng</Link></li>
            </ul>
          </div>

          {/* Partner */}
          <div>
            <h4 className="text-[11px] font-semibold text-tiki-text mb-2">Hợp tác & Liên kết</h4>
            <ul className="space-y-1.5 text-[10px] text-tiki-text-secondary">
              <li><Link href="#" className="hover:text-tiki-blue">Quy chế hoạt động</Link></li>
              <li><Link href="/products" className="hover:text-tiki-blue">Bán hàng cùng Tiki</Link></li>
            </ul>
          </div>

          {/* Payment */}
          <div>
            <h4 className="text-[11px] font-semibold text-tiki-text mb-2">Phương thức thanh toán</h4>
            <div className="flex flex-wrap gap-1 mb-3">
              {["Visa", "MC", "JCB", "ZaloPay", "Momo", "VNPay", "COD"].map((m) => (
                <span key={m} className="px-1.5 py-0.5 bg-white border border-tiki-border rounded text-[9px] text-tiki-text-secondary">{m}</span>
              ))}
            </div>
            <h4 className="text-[11px] font-semibold text-tiki-text mb-1">Dịch vụ giao hàng</h4>
            <div className="flex items-center gap-1.5">
              <span className="px-1.5 py-0.5 bg-tiki-blue text-white rounded text-[8px] font-bold">TIKI NOW</span>
              <span className="text-[10px] text-tiki-text-secondary">Giao nhanh 2h</span>
            </div>
          </div>

          {/* Logo & Social */}
          <div>
            <div className="flex flex-col items-start gap-1.5 mb-3">
              <svg width="64" height="22" viewBox="0 0 64 22" fill="none">
                <rect width="64" height="22" rx="4" fill="#1A94FF"/>
                <text x="10" y="16" fill="white" fontSize="12" fontWeight="700" fontFamily="Inter, sans-serif">Tiki</text>
              </svg>
              <span className="text-[8px] text-[#003EA1] font-semibold tracking-tight">TỐT &amp; NHANH</span>
            </div>
            <div className="flex gap-1">
              {["F", "Y", "Z", "I"].map((s) => (
                <span key={s} className="w-6 h-6 flex items-center justify-center bg-gray-100 rounded text-[9px] text-tiki-text-secondary font-medium">{s}</span>
              ))}
            </div>
          </div>
        </div>

        {/* Bottom bar */}
        <div className="border-t border-tiki-border py-3 text-[9px] text-tiki-text-secondary leading-5">
          <div className="font-medium text-tiki-text text-[10px] mb-0.5">Công ty TNHH TI KI</div>
          <div>Tòa nhà số 52, đường Út Tịch, P.4, Q. Tân Bình, TP. Hồ Chí Minh</div>
          <div>ĐKKD số: 0309532909 — cấp lần đầu: 06/01/2010. ĐT: (028) 38 321 232</div>
          <div className="mt-1">© 2010 - 2025 - Bản quyền của Công ty TNHH TI KI</div>
        </div>
      </div>
    </footer>
  );
}
