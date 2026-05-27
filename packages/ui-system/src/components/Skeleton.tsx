import { clsx } from "clsx";

interface SkeletonProps { className?: string; variant?: "text" | "circular" | "rectangular"; width?: string; height?: string; }

export function Skeleton({ className, variant = "text", width, height }: SkeletonProps) {
  const base = "animate-pulse bg-gray-200";
  const variants = { text: "rounded", circular: "rounded-full", rectangular: "rounded-md" };
  return <div className={clsx(base, variants[variant], className)} style={{ width, height }} />;
}
