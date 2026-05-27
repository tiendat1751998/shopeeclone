import { NextResponse } from "next/server";
import categoriesData from "@/data/tiki-categories.json";

export async function GET() {
  return NextResponse.json({
    success: true,
    data: categoriesData,
  });
}
