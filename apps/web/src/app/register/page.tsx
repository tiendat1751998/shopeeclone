"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Header } from "@/components/layout/header/Header";
import { Footer } from "@/components/layout/footer/Footer";
import { useAuthStore } from "@/stores/auth";

export default function RegisterPage() {
  const router = useRouter();
  const register = useAuthStore((s) => s.register);
  const isLoading = useAuthStore((s) => s.isLoading);
  const [form, setForm] = useState({
    display_name: "",
    email: "",
    phone: "",
    password: "",
    confirm_password: "",
  });
  const [error, setError] = useState("");

  function update(field: string, value: string) {
    setForm((prev) => ({ ...prev, [field]: value }));
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    if (form.password !== form.confirm_password) {
      setError("Mật khẩu xác nhận không khớp");
      return;
    }

    try {
      await register({
        display_name: form.display_name,
        email: form.email,
        username: form.email.split("@")[0],
        password: form.password,
        confirm_password: form.confirm_password,
      });
      router.push("/");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Đăng ký thất bại";
      setError(msg);
    }
  }

  return (
    <>
      <Header />
      <main className="py-12">
        <div className="max-w-[400px] mx-auto px-6">
          <div className="bg-white rounded-lg shadow-sm border border-tiki-border p-8">
            <h1 className="text-lg font-semibold text-tiki-text text-center mb-6">Đăng ký</h1>
            {error && (
              <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-600">
                {error}
              </div>
            )}
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-tiki-text mb-1">Họ tên</label>
                <input
                  type="text"
                  value={form.display_name}
                  onChange={(e) => update("display_name", e.target.value)}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:border-tiki-blue"
                  placeholder="Nguyễn Văn A"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-tiki-text mb-1">Email</label>
                <input
                  type="email"
                  value={form.email}
                  onChange={(e) => update("email", e.target.value)}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:border-tiki-blue"
                  placeholder="Email của bạn"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-tiki-text mb-1">Số điện thoại</label>
                <input
                  type="tel"
                  value={form.phone}
                  onChange={(e) => update("phone", e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:border-tiki-blue"
                  placeholder="0912345678"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-tiki-text mb-1">Mật khẩu</label>
                <input
                  type="password"
                  value={form.password}
                  onChange={(e) => update("password", e.target.value)}
                  required
                  minLength={8}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:border-tiki-blue"
                  placeholder="Mật khẩu (tối thiểu 8 ký tự)"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-tiki-text mb-1">Xác nhận mật khẩu</label>
                <input
                  type="password"
                  value={form.confirm_password}
                  onChange={(e) => update("confirm_password", e.target.value)}
                  required
                  minLength={8}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:border-tiki-blue"
                  placeholder="Nhập lại mật khẩu"
                />
              </div>
              <button
                type="submit"
                disabled={isLoading}
                className="w-full py-2.5 bg-tiki-blue text-white rounded-lg font-semibold text-sm hover:bg-tiki-blue-dark transition disabled:opacity-50"
              >
                {isLoading ? "Đang xử lý..." : "Đăng ký"}
              </button>
            </form>
            <div className="mt-6 text-center text-xs text-tiki-text-secondary">
              Đã có tài khoản?{" "}
              <Link href="/login" className="text-tiki-blue font-medium hover:underline">Đăng nhập</Link>
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </>
  );
}
