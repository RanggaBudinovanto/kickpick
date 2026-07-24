import * as React from "react";
import { cn } from "@/lib/utils";

export const Input = React.forwardRef<HTMLInputElement, React.InputHTMLAttributes<HTMLInputElement>>(
  ({ className, ...props }, ref) => {
    return (
      <input
        ref={ref}
        className={cn(
          "h-11 w-full rounded-[var(--radius-control)] border border-zinc-500/40 bg-background px-4 text-sm text-foreground placeholder:text-muted outline-none transition-shadow focus:ring-2 focus:ring-off-black dark:focus:ring-pure-white",
          className,
        )}
        {...props}
      />
    );
  },
);
Input.displayName = "Input";
