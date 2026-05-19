"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { useAuthStore } from "@/lib/store/auth";
import { ApiError } from "@/lib/api/client";

function validateEmail(email: string): string | null {
  if (!email || typeof email !== "string") return "Email is required";
  const trimmed = email.trim().toLowerCase();
  if (trimmed.length > 254) return "Email is too long";
  const re = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  if (!re.test(trimmed)) return "Invalid email format";
  return null;
}

function validatePassword(password: string): string | null {
  if (!password) return "Password is required";
  if (password.length < 8) return "Password must be at least 8 characters";
  if (password.length > 128) return "Password is too long";
  return null;
}

function sanitizeInput(value: string): string {
  return value.trim().slice(0, 200).replace(/[<>]/g, "");
}

export default function LoginPage() {
  const router = useRouter();
  const login = useAuthStore((s) => s.login);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<{ email?: string; password?: string }>({});

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    const emailErr = validateEmail(email);
    const passwordErr = validatePassword(password);
    if (emailErr || passwordErr) {
      setFieldErrors({ email: emailErr || undefined, password: passwordErr || undefined });
      return;
    }
    setFieldErrors({});

    setIsLoading(true);
    try {
      await login(sanitizeInput(email), password);
      router.push("/");
    } catch (e: unknown) {
      if (e instanceof ApiError && e.status === 429) {
        setError("Too many login attempts. Please try again later.");
      } else if (e instanceof ApiError) {
        setError(e.message || "Login failed. Please check your credentials.");
      } else {
        setError("Login failed. Please try again.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-[60vh] flex items-center justify-center px-4">
      <div className="w-full max-w-md card p-8">
        <h1 className="text-2xl font-bold mb-6">Login</h1>
        {error && <div className="bg-red-50 text-red-600 text-sm p-3 rounded mb-4">{error}</div>}
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Email"
            type="email"
            value={email}
            onChange={(e) => { setEmail(e.target.value); setFieldErrors((p) => ({ ...p, email: undefined })); }}
            placeholder="your@email.com"
            required
            error={fieldErrors.email}
            autoComplete="email"
            maxLength={254}
          />
          <Input
            label="Password"
            type="password"
            value={password}
            onChange={(e) => { setPassword(e.target.value); setFieldErrors((p) => ({ ...p, password: undefined })); }}
            placeholder="Enter password"
            required
            error={fieldErrors.password}
            autoComplete="current-password"
            maxLength={128}
          />
          <Button variant="primary" fullWidth isLoading={isLoading}>Login</Button>
        </form>
        <p className="text-center text-sm text-[#757575] mt-4">
          Don&apos;t have an account? <Link href="/register" className="text-[#ee4d2d] hover:underline">Register</Link>
        </p>
      </div>
    </div>
  );
}
