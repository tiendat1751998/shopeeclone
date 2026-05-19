"use client";
import { clsx } from "clsx";

interface PaginationProps { currentPage: number; totalPages: number; onPageChange: (page: number) => void; }

export function Pagination({ currentPage, totalPages, onPageChange }: PaginationProps) {
  if (totalPages <= 1) return null;

  const pages: (number | "...")[] = [];
  for (let i = 1; i <= totalPages; i++) {
    if (i === 1 || i === totalPages || (i >= currentPage - 2 && i <= currentPage + 2)) {
      pages.push(i);
    } else if (pages[pages.length - 1] !== "...") {
      pages.push("...");
    }
  }

  return (
    <div className="flex items-center justify-center gap-1 py-6">
      <button onClick={() => onPageChange(currentPage - 1)} disabled={currentPage === 1}
        className="px-3 py-2 text-sm rounded border border-[#e8e8e8] disabled:opacity-40 hover:border-[#ee4d2d] transition-colors">
        &laquo;
      </button>
      {pages.map((p, i) =>
        p === "..." ? <span key={`dots-${i}`} className="px-2 text-[#757575]">...</span> : (
          <button key={`page-${p}`} onClick={() => onPageChange(p)}
            className={clsx("px-3 py-2 text-sm rounded border transition-colors",
              p === currentPage ? "bg-[#ee4d2d] text-white border-[#ee4d2d]" : "border-[#e8e8e8] hover:border-[#ee4d2d]"
            )}>{p}</button>
        )
      )}
      <button onClick={() => onPageChange(currentPage + 1)} disabled={currentPage === totalPages}
        className="px-3 py-2 text-sm rounded border border-[#e8e8e8] disabled:opacity-40 hover:border-[#ee4d2d] transition-colors">
        &raquo;
      </button>
    </div>
  );
}
