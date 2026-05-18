interface PriceProps { amount: number; originalAmount?: number; currency?: string; size?: "sm" | "md" | "lg"; className?: string; }

export function Price({ amount, originalAmount, currency = "S$", size = "md", className }: PriceProps) {
  const sizes = { sm: "text-sm", md: "text-lg", lg: "text-2xl" };
  const formatted = `${currency}${amount.toLocaleString("en-SG", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
  return (
    <div className={className}>
      <span className={clsx("font-semibold text-[#ee4d2d]", sizes[size])}>{formatted}</span>
      {originalAmount && originalAmount > amount && (
        <span className="ml-2 text-sm text-[#757575] line-through">{currency}{originalAmount.toLocaleString("en-SG", { minimumFractionDigits: 2 })}</span>
      )}
    </div>
  );
}

import { clsx } from "clsx";
