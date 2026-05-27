import { NextResponse } from "next/server";
import promotionsData from "@/data/tiki-promotions.json";

export async function GET() {
  return NextResponse.json({ success: true, data: promotionsData });
}
