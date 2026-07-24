import * as React from "react";
import { cn } from "@/lib/utils";

export function Card({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        "rounded-[var(--radius-card)] border border-border bg-background shadow-[0_1px_3px_rgba(39,39,42,0.08)]",
        className,
      )}
      {...props}
    />
  );
}
