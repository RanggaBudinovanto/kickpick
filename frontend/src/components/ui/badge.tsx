import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

const badgeVariants = cva(
  "inline-flex items-center gap-1 rounded-[var(--radius-control)] px-2.5 py-1 text-xs font-medium tracking-[0.01em] uppercase",
  {
    variants: {
      variant: {
        neutral: "bg-zinc-100 text-off-black dark:bg-zinc-800 dark:text-pure-white",
        strong: "bg-zinc-950 text-pure-white",
      },
    },
    defaultVariants: {
      variant: "neutral",
    },
  },
);

export interface BadgeProps
  extends React.HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof badgeVariants> {}

export function Badge({ className, variant, ...props }: BadgeProps) {
  return <span className={cn(badgeVariants({ variant, className }))} {...props} />;
}
