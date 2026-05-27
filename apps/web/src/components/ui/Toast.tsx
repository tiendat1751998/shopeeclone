"use client";

import { useUIStore } from "@/stores/ui";

export function ToastContainer() {
  const toasts = useUIStore((s) => s.toastNotifications);
  const dismissToast = useUIStore((s) => s.dismissToast);

  if (toasts.length === 0) return null;

  return (
    <div className="toast-container">
      {toasts.map((toast) => (
        <div key={toast.id} className={`toast toast-${toast.type}`}>
          <div className="toast-title">{toast.title}</div>
          <div>{toast.message}</div>
          <button
            onClick={() => dismissToast(toast.id)}
            style={{ position: "absolute", top: "8px", right: "8px", background: "none", border: "none", cursor: "pointer", fontSize: "14px", opacity: 0.5 }}
          >
            ✕
          </button>
        </div>
      ))}
    </div>
  );
}
