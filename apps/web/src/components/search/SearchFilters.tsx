"use client";
import { clsx } from "clsx";
import type { SearchFilters } from "@/lib/types";

interface SearchFiltersProps { filters: SearchFilters; onChange: (filters: SearchFilters) => void; resultCount?: number; }

const sortOptions = [
  { value: "relevance", label: "Relevance" },
  { value: "best_selling", label: "Best Selling" },
  { value: "newest", label: "Newest" },
  { value: "price_asc", label: "Price: Low to High" },
  { value: "price_desc", label: "Price: High to Low" },
];

export function SearchFiltersBar({ filters, onChange, resultCount }: SearchFiltersProps) {
  return (
    <div className="card p-3 mb-4 flex flex-wrap items-center justify-between gap-3">
      <div className="text-sm text-[#757575]">
        {resultCount !== undefined && <span>{resultCount.toLocaleString()} results found</span>}
      </div>
      <div className="flex items-center gap-3">
        <label className="text-sm text-[#757575]">Sort by:</label>
        <select
          value={filters.sort_by || "relevance"}
          onChange={(e) => onChange({ ...filters, sort_by: e.target.value })}
          className="text-sm border border-[#e8e8e8] rounded px-3 py-1.5 focus:outline-none focus:ring-2 focus:ring-[#ee4d2d]"
        >
          {sortOptions.map((opt) => <option key={opt.value} value={opt.value}>{opt.label}</option>)}
        </select>
      </div>
    </div>
  );
}

interface PriceFilterProps { minPrice?: number; maxPrice?: number; onChange: (min?: number, max?: number) => void; }

export function PriceFilter({ minPrice, maxPrice, onChange }: PriceFilterProps) {
  return (
    <div className="space-y-3">
      <h4 className="text-sm font-medium">Price Range</h4>
      <div className="flex items-center gap-2">
        <input type="number" placeholder="Min" value={minPrice || ""} onChange={(e) => onChange(e.target.value ? Number(e.target.value) : undefined, maxPrice)}
          className="w-full px-2 py-1.5 text-sm border border-[#e8e8e8] rounded focus:outline-none focus:ring-1 focus:ring-[#ee4d2d]" />
        <span className="text-[#bdbdbd]">—</span>
        <input type="number" placeholder="Max" value={maxPrice || ""} onChange={(e) => onChange(minPrice, e.target.value ? Number(e.target.value) : undefined)}
          className="w-full px-2 py-1.5 text-sm border border-[#e8e8e8] rounded focus:outline-none focus:ring-1 focus:ring-[#ee4d2d]" />
      </div>
      {[[0, 20], [20, 50], [50, 100], [100, 500], [500, undefined]].map(([min, max]) => (
        <button key={`${min}-${max}`} onClick={() => onChange(min, max)}
          className={clsx("block w-full text-left text-sm py-1 px-2 rounded hover:bg-[#fff0ed] hover:text-[#ee4d2d] transition-colors",
            minPrice === min && maxPrice === max ? "bg-[#fff0ed] text-[#ee4d2d] font-medium" : "text-[#757575]"
          )}>
          {max ? `S$${min} - S$${max}` : `S$${min}+`}
        </button>
      ))}
    </div>
  );
}
