"use client";

import { useState, useRef, useEffect } from "react";
import { useLocale, useTranslations } from "next-intl";
import { IconChevronDown } from "@tabler/icons-react";
import { usePathname, useRouter } from "@/i18n/navigation";
import { routing } from "@/i18n/routing";
import { usePreferencesStore, type Currency } from "@/stores/preferences";
import { cn } from "@/lib/utils";

export function PreferenceSwitcher() {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  const locale = useLocale();
  const pathname = usePathname();
  const router = useRouter();
  const currency = usePreferencesStore((s) => s.currency);
  const setCurrency = usePreferencesStore((s) => s.setCurrency);
  const tCurrency = useTranslations("Currency");

  useEffect(() => {
    function onClickOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", onClickOutside);
    return () => document.removeEventListener("mousedown", onClickOutside);
  }, []);

  function switchLocale(nextLocale: string) {
    router.replace(pathname, { locale: nextLocale });
    setOpen(false);
  }

  function switchCurrency(next: Currency) {
    setCurrency(next);
    setOpen(false);
  }

  return (
    <div className="relative" ref={ref}>
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        className="flex h-11 items-center gap-1 rounded-[var(--radius-control)] px-3 text-sm font-medium text-foreground hover:bg-surface"
      >
        {locale.toUpperCase()} · {currency}
        <IconChevronDown size={16} />
      </button>

      {open && (
        <div className="absolute right-0 top-full mt-2 w-48 rounded-[var(--radius-card)] border border-border bg-background p-3 shadow-lg">
          <p className="mb-2 text-xs font-medium uppercase tracking-[0.01em] text-muted">
            Bahasa / Language
          </p>
          <div className="mb-3 flex gap-2">
            {routing.locales.map((l) => (
              <button
                key={l}
                onClick={() => switchLocale(l)}
                className={cn(
                  "flex-1 rounded-[var(--radius-control)] border border-zinc-500/40 py-1.5 text-sm",
                  l === locale && "bg-off-black text-pure-white dark:bg-pure-white dark:text-off-black",
                )}
              >
                {l.toUpperCase()}
              </button>
            ))}
          </div>

          <p className="mb-2 text-xs font-medium uppercase tracking-[0.01em] text-muted">
            Mata Uang / Currency
          </p>
          <div className="flex flex-col gap-1">
            {(["IDR", "USD"] as Currency[]).map((c) => (
              <button
                key={c}
                onClick={() => switchCurrency(c)}
                className={cn(
                  "rounded-[var(--radius-control)] px-2 py-1.5 text-left text-sm hover:bg-surface",
                  c === currency && "font-medium",
                )}
              >
                {c === "IDR" ? tCurrency("idr") : tCurrency("usd")}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
