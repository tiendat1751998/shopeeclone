// Next.js API Routes - Proxy to Go Gateway
// apps/web/src/app/api/v1/[...path]/route.ts
import { NextRequest, NextResponse } from "next/server";

const GATEWAY_URL = process.env.GATEWAY_URL || "http://gateway:8080";

function getClientIP(request: NextRequest): string {
  // Check X-Forwarded-For first (from nginx/other proxies)
  const xff = request.headers.get("x-forwarded-for");
  if (xff) return xff.split(",")[0].trim();
  // Check X-Real-IP
  const xri = request.headers.get("x-real-ip");
  if (xri) return xri.trim();
  // Fallback to Next.js connection info
  return request.headers.get("x-client-ip") || "unknown";
}

async function proxy(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> },
  method: string,
  hasBody: boolean
) {
  const { path } = await params;
  const pathStr = path.join("/");
  const searchParams = request.nextUrl.searchParams.toString();
  const url = `${GATEWAY_URL}/api/v1/${pathStr}${searchParams ? `?${searchParams}` : ""}`;
  const clientIP = getClientIP(request);

  try {
    const body = hasBody ? await request.json().catch(() => undefined) : undefined;
    const res = await fetch(url, {
      method,
      headers: {
        "Content-Type": "application/json",
        "X-Forwarded-For": clientIP,
        "X-Real-IP": clientIP,
        ...(request.headers.get("authorization")
          ? { Authorization: request.headers.get("authorization")! }
          : {}),
      },
      ...(body !== undefined ? { body: JSON.stringify(body) } : {}),
      next: { revalidate: 0 },
    });

    const data = await res.json().catch(() => null);
    return NextResponse.json(data, { status: res.status });
  } catch (error) {
    console.error("API proxy error:", error);
    return NextResponse.json(
      { success: false, error: "Gateway unavailable" },
      { status: 502 }
    );
  }
}

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  return proxy(request, { params }, "GET", false);
}

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  return proxy(request, { params }, "POST", true);
}

export async function PATCH(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  return proxy(request, { params }, "PATCH", true);
}

export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  return proxy(request, { params }, "PUT", true);
}

export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  return proxy(request, { params }, "DELETE", false);
}
