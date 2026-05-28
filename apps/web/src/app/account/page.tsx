"use client";

import { useState } from "react";
import { useAuthStore } from "@/stores/auth";
import { authApi } from "@/lib/api/client";

export default function AccountPage() {
  const user = useAuthStore((s) => s.user);
  const refreshUser = useAuthStore((s) => s.refreshUser);
  const [displayName, setDisplayName] = useState(user?.display_name || "");
  const [phone, setPhone] = useState(user?.phone || "");
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState("");

  async function handleSave(e: React.FormEvent) {
    e.preventDefault();
    setSaving(true);
    setMessage("");
    try {
      await authApi.put("/auth/profile", { display_name: displayName, phone });
      await refreshUser();
      setMessage("Cập nhật thành công!");
    } catch {
      setMessage("Cập nhật thất bại. Vui lòng thử lại.");
    } finally {
      setSaving(false);
    }
  }

  if (!user) return null;

  return (
    <main className="bg-[#F5F5FA] py-6 min-h-[60vh]">
      <div className="max-w-2xl mx-auto px-3">
        <h1 className="text-lg font-semibold text-tiki-text mb-4">Thông tin tài khoản</h1>

        <div className="bg-white rounded-lg border border-tiki-border p-6">
          {/* Email (read-only) */}
          <div className="mb-4">
            <label className="block text-xs font-medium text-tiki-text-secondary mb-1">Email</label>
            <div className="text-sm text-tiki-text bg-gray-50 rounded-lg px-3 py-2">{user.email}</div>
          </div>

          <form onSubmit={handleSave} className="space-y-4">
            <div>
              <label className="block text-xs font-medium text-tiki-text-secondary mb-1">Họ và tên</label>
              <input
                type="text"
                value={displayName}
                onChange={(e) => setDisplayName(e.target.value)}
                className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-tiki-blue focus:ring-1 focus:ring-tiki-blue outline-none"
                placeholder="Nhập họ và tên"
              />
            </div>

            <div>
              <label className="block text-xs font-medium text-tiki-text-secondary mb-1">Số điện thoại</label>
              <input
                type="tel"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-tiki-blue focus:ring-1 focus:ring-tiki-blue outline-none"
                placeholder="Nhập số điện thoại"
              />
            </div>

            <div>
              <label className="block text-xs font-medium text-tiki-text-secondary mb-1">Tên đăng nhập</label>
              <div className="text-sm text-tiki-text-secondary bg-gray-50 rounded-lg px-3 py-2">{user.username}</div>
            </div>

            <div>
              <label className="block text-xs font-medium text-tiki-text-secondary mb-1">Vai trò</label>
              <div className="text-sm text-tiki-text-secondary bg-gray-50 rounded-lg px-3 py-2 capitalize">{user.role}</div>
            </div>

            <div>
              <label className="block text-xs font-medium text-tiki-text-secondary mb-1">Ngày tham gia</label>
              <div className="text-sm text-tiki-text-secondary bg-gray-50 rounded-lg px-3 py-2">
                {new Date(user.created_at).toLocaleDateString("vi-VN")}
              </div>
            </div>

            {message && (
              <div className={`text-sm px-3 py-2 rounded-lg ${message.includes("thành công") ? "bg-green-50 text-green-700" : "bg-red-50 text-red-700"}`}>
                {message}
              </div>
            )}

            <button
              type="submit"
              disabled={saving}
              className="px-6 py-2 bg-tiki-blue text-white rounded-lg text-sm font-medium hover:bg-tiki-blue-dark transition disabled:opacity-50"
            >
              {saving ? "Đang lưu..." : "Lưu thay đổi"}
            </button>
          </form>
        </div>
      </div>
    </main>
  );
}
