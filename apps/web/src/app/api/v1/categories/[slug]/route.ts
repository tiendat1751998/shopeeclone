import { NextResponse } from "next/server";
import categoriesData from "@/data/tiki-categories.json";

function findCategory(cats: any[], slug: string): any | null {
  for (const cat of cats) {
    if (cat.slug === slug) return cat;
    if (cat.children) {
      const found = findCategory(cat.children, slug);
      if (found) return found;
    }
  }
  return null;
}

export async function GET(
  _request: Request,
  { params }: { params: Promise<{ slug: string }> }
) {
  const { slug } = await params;
  const category = findCategory(categoriesData, slug);
  if (!category) {
    return NextResponse.json(
      { success: false, error: "Category not found" },
      { status: 404 }
    );
  }
  return NextResponse.json({ success: true, data: category });
}
