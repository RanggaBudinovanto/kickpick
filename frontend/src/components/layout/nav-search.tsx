"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { IconSearch } from "@tabler/icons-react";
import { useRouter } from "@/i18n/navigation";
import { Input } from "@/components/ui/input";
import { useAutocomplete } from "@/hooks/use-products";

// Compact search used in the navbar (moved here from the old hero — DESIGN.md
// § 5 revision moves search out of the hero banner). Submits on Enter rather
// than a separate button, since navbar height doesn't have room for one.
export function NavSearch() {
  const [query, setQuery] = useState("");
  const [focused, setFocused] = useState(false);
  const router = useRouter();
  const t = useTranslations("Nav");
  const { data } = useAutocomplete(query);

  const suggestions = [...(data?.products ?? []), ...(data?.brands ?? [])].slice(0, 6);

  function submit(q: string) {
    if (!q.trim()) return;
    setFocused(false);
    router.push(`/cari?q=${encodeURIComponent(q)}`);
  }

  return (
    <div className="relative w-full">
      <form
        onSubmit={(e) => {
          e.preventDefault();
          submit(query);
        }}
      >
        <div className="relative">
          <IconSearch
            size={18}
            className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-muted"
          />
          <Input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onFocus={() => setFocused(true)}
            onBlur={() => setTimeout(() => setFocused(false), 150)}
            placeholder={t("searchPlaceholder")}
            className="pl-10"
          />
        </div>
      </form>

      {focused && suggestions.length > 0 && (
        <div className="absolute left-0 right-0 top-full z-10 mt-2 rounded-[var(--radius-card)] border border-border bg-background p-2 shadow-lg">
          {suggestions.map((s) => (
            <button
              key={s.id}
              type="button"
              onMouseDown={() => submit(s.name)}
              className="block w-full rounded-[var(--radius-control)] px-3 py-2 text-left text-sm hover:bg-surface"
            >
              {s.name}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
