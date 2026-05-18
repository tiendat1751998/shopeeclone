"use client";
import Link from "next/link";
import { useState, useEffect } from "react";
import { categoriesApi } from "@/lib/api/products";
import type { Category } from "@/lib/types";
import { Skeleton } from "@/components/ui/Skeleton";

export function CategorySidebar() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [expanded, setExpanded] = useState<Set<string>>(new Set());

  useEffect(() => {
    categoriesApi.getTree()
      .then(setCategories)
      .catch(() => {})
      .finally(() => setIsLoading(false));
  }, []);

  const toggle = (id: string) => {
    setExpanded((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id); else next.add(id);
      return next;
    });
  };

  const renderCategory = (cat: Category, depth = 0) => {
    const hasChildren = cat.children && cat.children.length > 0;
    const isExpanded = expanded.has(cat.id);

    return (
      <li key={cat.id} className="group">
        <div className="flex items-center justify-between">
          <Link href={`/categories/${cat.slug}`}
            className={`flex-1 py-1.5 text-sm hover:text-[#ee4d2d] transition-colors ${depth > 0 ? "pl-4" : "font-medium"}`}>
            {cat.name}
            {cat.product_count !== undefined && (
              <span className="ml-1 text-xs text-[#bdbdbd]">({cat.product_count})</span>
            )}
          </Link>
          {hasChildren && (
            <button onClick={() => toggle(cat.id)} className="p-1 text-[#757575] hover:text-[#ee4d2d]">
              <svg className={`w-3 h-3 transition-transform ${isExpanded ? "rotate-90" : ""}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </button>
          )}
        </div>
        {hasChildren && isExpanded && (
          <ul className="ml-2 border-l border-[#e8e8e8]">
            {cat.children!.map((child) => renderCategory(child, depth + 1))}
          </ul>
        )}
      </li>
    );
  };

  return (
    <aside className="card p-4">
      <h3 className="font-semibold text-sm mb-3 flex items-center gap-2">
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 10h16M4 14h16M4 18h16" /></svg>
        All Categories
      </h3>
      {isLoading ? (
        <div className="space-y-2">
          {Array.from({ length: 8 }).map((_, i) => <Skeleton key={i} className="h-6 w-full" />)}
        </div>
      ) : (
        <ul className="space-y-0.5">{categories.map((c) => renderCategory(c))}</ul>
      )}
    </aside>
  );
}
