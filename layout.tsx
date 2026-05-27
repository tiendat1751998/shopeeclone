"use client";
import { useEffect, useState } from "react";
import { usePathname } from "next/navigation";

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const [page, setPage] = useState("overview");
  useEffect(() => { setPage(pathname.split("/").pop() || "overview"); }, [pathname]);
  
  const nav = [
    { id: "overview", label: "📊 Overview" },
    { id: "services", label: "⚙️ Services" },
    { id: "topology", label: "🌐 Topology" },
    { id: "deployments", label: "🚀 Deployments" },
    { id: "alerts", label: "🔔 Alerts" },
    { id: "incidents", label: "🔥 Incidents" },
    { id: "capacity", label: "📈 Capacity" },
    { id: "audit", label: "📋 Audit" },
  ];
  
  return (
    <div style={{ display: "flex", height: "100vh", background: "#f5f5f5" }}>
      <aside style={{ width: 220, background: "#0B1E33", color: "#fff", flexShrink: 0, display: "flex", flexDirection: "column" }}>
        <div style={{ padding: 18, fontSize: 20, fontWeight: 800, background: "rgba(0,0,0,.2)" }}>shopee<span style={{ color: "#189eff" }}>admin</span></div>
        <nav style={{ flex: 1, overflowY: "auto", padding: "8px 0" }}>
          {nav.map(n => (
            <a key={n.id} href={`/dashboard/${n.id}`} onClick={e => { e.preventDefault(); setPage(n.id); }}
              style={{ display: "block", padding: "10px 18px", color: page === n.id ? "#fff" : "rgba(255,255,255,.6)", textDecoration: "none", fontSize: 13, background: page === n.id ? "rgba(255,255,255,.1)" : "transparent", borderLeft: page === n.id ? "3px solid #189eff" : "3px solid transparent" }}>
              {n.label}
            </a>
          ))}
        </nav>
      </aside>
      <div style={{ flex: 1, display: "flex", flexDirection: "column", overflow: "hidden" }}>
        <header style={{ height: 52, background: "#fff", borderBottom: "1px solid #e8e8e8", display: "flex", alignItems: "center", padding: "0 24px" }}>
          <h1 style={{ fontSize: 16, fontWeight: 600 }}>{nav.find(n => n.id === page)?.label || "Dashboard"}</h1>
        </header>
        <main style={{ flex: 1, overflowY: "auto", padding: 24 }}>
          {page === "overview" ? children : (
            <div style={{ background: "#fff", border: "1px solid #e8e8e8", borderRadius: 8, padding: 40, textAlign: "center", color: "#757575" }}>
              <h3>{nav.find(n => n.id === page)?.label}</h3>
              <p>Section under development</p>
            </div>
          )}
        </main>
      </div>
    </div>
  );
}
