"use client";

import { useEffect } from "react";
import { Header } from "@/components/layout/header/Header";
import { Footer } from "@/components/layout/footer/Footer";

export default function ErrorPage({ error, reset }: { error: Error & { digest?: string }; reset: () => void }) {
  useEffect(() => { console.error(error); }, [error]);

  return (
    <>
      <Header />
      <main className="py-16 text-center">
        <p className="text-5xl mb-4">😕</p>
        <h1 className="text-lg font-semibold text-tiki-text mb-2">Có lỗi xảy ra</h1>
        <p className="text-sm text-tiki-text-secondary mb-6">Vui lòng thử lại sau</p>
        <button onClick={reset} className="px-6 py-2.5 bg-tiki-blue text-white rounded-lg font-semibold text-sm hover:bg-tiki-blue-dark transition">
          Thử lại
        </button>
      </main>
      <Footer />
    </>
  );
}
