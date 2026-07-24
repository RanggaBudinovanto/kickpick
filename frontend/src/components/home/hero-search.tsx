"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { IconSearch } from "@tabler/icons-react";
import { useRouter } from "@/i18n/navigation";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useAutocomplete } from "@/hooks/use-products";

export function HeroSearch() {
  const [query, setQuery] = useState("");
  const [focused, setFocused] = useState(false);
  const router = useRouter();
  const t = useTranslations("Nav");
  const { data } = useAutocomplete(query);

  const suggestions = [...(data?.products ?? []), ...(data?.brands ?? [])].slice(0, 6);

  function submit(q: string) {
    if (!q.trim()) return;
    router.push(`/cari?q=${encodeURIComponent(q)}`);
  }

  return (
    <div className="relative">
      <form
        onSubmit={(e) => {
          e.preventDefault();
          submit(query);
        }}
        className="flex gap-2"
      >
        <div className="relative flex-1">
          <IconSearch size={18} className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-muted" />
          <Input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onFocus={() => setFocused(true)}
            onBlur={() => setTimeout(() => setFocused(false), 150)}
            placeholder={t("searchPlaceholder")}
            className="pl-10"
          />
        </div>
        <Button type="submit">{t("search")}</Button>
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
