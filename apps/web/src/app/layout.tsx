import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Providers } from "./providers";
import { ToastContainer } from "@/components/ui/Toast";

const inter = Inter({
  subsets: ["latin", "vietnamese"],
  display: "swap",
  variable: "--font-inter",
});

export const metadata: Metadata = {
  title: {
    default: "Tiki - Mua hàng online giá tốt, hàng chuẩn, ship nhanh",
    template: "%s | Tiki",
  },
  description: "Mua sắm trực tuyến hàng triệu sản phẩm. Giá tốt, giao nhanh, miễn phí ship.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="vi" className={inter.variable}>
      <body className="min-h-screen bg-tiki-bg text-tiki-text antialiased font-sans">
        <Providers>
          {children}
          <ToastContainer />
        </Providers>
      </body>
    </html>
  );
}
