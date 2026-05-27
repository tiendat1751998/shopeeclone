import { NextResponse } from "next/server";
import categoriesData from "@/data/tiki-categories.json";

function flattenCategories(cats: any[]): any[] {
  const result: any[] = [];
  for (const cat of cats) {
    result.push(cat);
    if (cat.children) {
      result.push(...flattenCategories(cat.children));
    }
  }
  return result;
}

export async function GET() {
  const flat = flattenCategories(categoriesData);
  return NextResponse.json({
    success: true,
    data: flat,
  });
}
