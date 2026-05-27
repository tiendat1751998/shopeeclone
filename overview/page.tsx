"use client";
import { useEffect, useState } from "react";

export default function OverviewPage() {
  const [data, setData] = useState<any>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("/api/gateway/admin/dashboard/summary")
      .then(r => r.json())
      .then(d => { setData(d); setLoading(false); })
      .catch(e => { setError(e.message); setLoading(false); });
  }, []);

  if (loading) return <div style={{ padding: 20, color: "#757575" }}>Loading...</div>;
  if (error) return <div style={{ background: "#fff3cd", color: "#856404", padding: 16, borderRadius: 8 }}>⚠️ {error}</div>;
  const s = data || {};
  return (
    <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(200px, 1fr))", gap: 16 }}>
      <div style={{ background: "#fff", border: "1px solid #e8e8e8", borderRadius: 8, padding: 20 }}>
        <div style={{ fontSize: 13, color: "#757575" }}>Total Services</div>
        <div style={{ fontSize: 28, fontWeight: 700 }}>{s.total_services || 0}</div>
      </div>
      <div style={{ background: "#fff", border: "1px solid #e8e8e8", borderRadius: 8, padding: 20 }}>
        <div style={{ fontSize: 13, color: "#757575" }}>Healthy</div>
        <div style={{ fontSize: 28, fontWeight: 700, color: "#00bfa5" }}>{s.healthy_services || 0}</div>
      </div>
      <div style={{ background: "#fff", border: "1px solid #e8e8e8", borderRadius: 8, padding: 20 }}>
        <div style={{ fontSize: 13, color: "#757575" }}>Degraded</div>
        <div style={{ fontSize: 28, fontWeight: 700, color: "#f5a623" }}>{s.degraded_services || 0}</div>
      </div>
      <div style={{ background: "#fff", border: "1px solid #e8e8e8", borderRadius: 8, padding: 20 }}>
        <div style={{ fontSize: 13, color: "#757575" }}>Unhealthy</div>
        <div style={{ fontSize: 28, fontWeight: 700, color: "#ff424f" }}>{s.unhealthy_services || 0}</div>
      </div>
    </div>
  );
}
