"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { useAuthStore } from "@/lib/store/auth";
import { ApiError } from "@/lib/api/client";

function validateEmail(email: string): string | null {
  if (!email?.trim()) return "Email is required";
  const trimmed = email.trim().toLowerCase();
  if (trimmed.length > 254) return "Email is too long";
  if (!/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(trimmed)) return "Invalid email format";
  return null;
}

function validatePassword(password: string): string | null {
  if (!password) return "Password is required";
  if (password.length < 8) return "Password must be at least 8 characters";
  if (password.length > 128) return "Password is too long";
  return null;
}

function validateUsername(username: string): string | null {
  if (!username?.trim()) return "Username is required";
  const trimmed = username.trim();
  if (trimmed.length < 3) return "Username must be at least 3 characters";
  if (trimmed.length > 32) return "Username is too long";
  if (!/^[a-zA-Z0-9_-]+$/.test(trimmed)) return "Username can only contain letters, numbers, underscores, and hyphens";
  return null;
}

function sanitize(value: string): string {
  return value.trim().slice(0, 200).replace(/[<>]/g, "");
}

export default function RegisterPage() {
  const router = useRouter();
  const register = useAuthStore((s) => s.register);
  const [form, setForm] = useState({ email: "", password: "", username: "", display_name: "" });
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    const errors: Record<string, string> = {};
    const emailErr = validateEmail(form.email);
    const passwordErr = validatePassword(form.password);
    const usernameErr = validateUsername(form.username);
    if (emailErr) errors.email = emailErr;
    if (passwordErr) errors.password = passwordErr;
    if (usernameErr) errors.username = usernameErr;
    if (!form.display_name.trim()) errors.display_name = "Display name is required";

    if (Object.keys(errors).length > 0) {
      setFieldErrors(errors);
      return;
    }
    setFieldErrors({});

    setIsLoading(true);
    try {
      await register({
        email: sanitize(form.email.toLowerCase()),
        password: form.password,
        username: sanitize(form.username),
        display_name: sanitize(form.display_name),
      });
      router.push("/");
    } catch (e: unknown) {
      if (e instanceof ApiError && e.status === 429) {
        setError("Too many registration attempts. Please try again later.");
      } else if (e instanceof ApiError) {
        setError(e.message || "Registration failed.");
      } else {
        setError("Registration failed. Please try again.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const updateField = (field: string, value: string) => {
    setForm((prev) => ({ ...prev, [field]: value }));
    setFieldErrors((prev) => ({ ...prev, [field]: "" }));
  };

  return (
    <div className="min-h-[60vh] flex items-center justify-center px-4">
      <div className="w-full max-w-md card p-8">
        <h1 className="text-2xl font-bold mb-6">Create Account</h1>
        {error && <div className="bg-red-50 text-red-600 text-sm p-3 rounded mb-4">{error}</div>}
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input label="Email" type="email" value={form.email} onChange={(e) => updateField("email", e.target.value)} required error={fieldErrors.email} maxLength={254} />
          <Input label="Username" value={form.username} onChange={(e) => updateField("username", e.target.value)} required error={fieldErrors.username} maxLength={32} />
          <Input label="Display Name" value={form.display_name} onChange={(e) => updateField("display_name", e.target.value)} required error={fieldErrors.display_name} maxLength={100} />
          <Input label="Password" type="password" value={form.password} onChange={(e) => updateField("password", e.target.value)} required error={fieldErrors.password} maxLength={128} />
          <Button variant="primary" fullWidth isLoading={isLoading}>Register</Button>
        </form>
        <p className="text-center text-sm text-[#757575] mt-4">Already have an account? <Link href="/login" className="text-[#ee4d2d] hover:underline">Login</Link></p>
      </div>
    </div>
  );
}
