import { clsx } from "clsx";

interface SkeletonProps { className?: string; variant?: "text" | "circular" | "rectangular"; width?: string; height?: string; }

export function Skeleton({ className, variant = "text", width, height }: SkeletonProps) {
  const base = "skeleton";
  const variants = { text: "rounded", circular: "rounded-full", rectangular: "rounded-md" };
  return <div className={clsx(base, variants[variant], className)} style={{ width, height }} />;
}

export function ProductCardSkeleton() {
  return (
    <div className="card p-0">
      <Skeleton variant="rectangular" className="w-full aspect-square" />
      <div className="p-3 space-y-2">
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-3/4" />
        <Skeleton className="h-5 w-1/3" />
        <div className="flex gap-2">
          <Skeleton className="h-3 w-12" />
          <Skeleton className="h-3 w-16" />
        </div>
      </div>
    </div>
  );
}

export function ProductGridSkeleton({ count = 12 }: { count?: number }) {
  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-3">
      {Array.from({ length: count }).map((_, i) => <ProductCardSkeleton key={i} />)}
    </div>
  );
}
