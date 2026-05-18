import { clsx } from "clsx";
import type { ReactNode } from "react";

interface BadgeProps { variant?: "sale" | "new" | "official" | "warning" | "info"; children: ReactNode; className?: string; }

export function Badge({ variant = "info", children, className }: BadgeProps) {
  const variants = {
    sale: "bg-[#ee4d2d] text-white",
    new: "bg-[#00bfa5] text-white",
    official: "bg-[#f5a623] text-white",
    warning: "bg-yellow-500 text-white",
    info: "bg-gray-100 text-[#757575]",
  };
  return <span className={clsx("inline-flex items-center px-2 py-0.5 rounded text-xs font-medium", variants[variant], className)}>{children}</span>;
}
