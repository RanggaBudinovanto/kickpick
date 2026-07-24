"use client";

import { useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { useRouter } from "@/i18n/navigation";
import type { Brand } from "@/types/api";

const CATEGORIES = ["running", "lifestyle", "training", "formal"];

export function SearchFilters({ brands }: { brands: Brand[] }) {
  const router = useRouter();
  const searchParams = useSearchParams();

  const kategori = searchParams.get("kategori") ?? "";
  const brandIdsParam = searchParams.get("brand_ids") ?? "";
  const selectedBrandIds = brandIdsParam ? brandIdsParam.split(",") : [];
  const filter = searchParams.get("filter") ?? "";

  function updateFilter(key: string, value: string) {
    const next = new URLSearchParams(searchParams.toString());
    if (value) {
      next.set(key, value);
    } else {
      next.delete(key);
    }
    router.push(`/cari?${next.toString()}`);
  }

  function toggleBrand(brandId: string, checked: boolean) {
    const nextIds = checked
      ? [...selectedBrandIds, brandId]
      : selectedBrandIds.filter((id) => id !== brandId);
    updateFilter("brand_ids", nextIds.join(","));
  }

  function resetFilters() {
    router.push("/cari");
  }

  return (
    <div className="rounded-[var(--radius-card)] border border-border p-4">
      <div className="mb-5">
        <label
          htmlFor="filter-kategori"
          className="mb-2 block text-xs font-medium uppercase tracking-[0.01em] text-muted"
        >
          Kategori
        </label>
        <select
          id="filter-kategori"
          value={kategori}
          onChange={(e) => updateFilter("kategori", e.target.value)}
          className="h-10 w-full rounded-[var(--radius-control)] border border-zinc-500/40 bg-background px-3 text-sm"
        >
          <option value="">Semua kategori</option>
          {CATEGORIES.map((c) => (
            <option key={c} value={c}>
              {c.charAt(0).toUpperCase() + c.slice(1)}
            </option>
          ))}
        </select>
      </div>

      <div className="mb-5">
        <label className="mb-2 block text-xs font-medium uppercase tracking-[0.01em] text-muted">
          Brand
        </label>
        <div className="flex max-h-48 flex-col gap-1.5 overflow-y-auto">
          {brands.map((b) => (
            <label key={b.id} className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                checked={selectedBrandIds.includes(b.id)}
                onChange={(e) => toggleBrand(b.id, e.target.checked)}
              />
              {b.name}
            </label>
          ))}
        </div>
      </div>

      <label className="flex items-center gap-2 text-sm">
        <input
          type="checkbox"
          checked={filter === "rare"}
          onChange={(e) => updateFilter("filter", e.target.checked ? "rare" : "")}
        />
        Rare / Limited saja
      </label>

      <Button variant="ghost" size="sm" className="mt-5 w-full" onClick={resetFilters}>
        Reset filter
      </Button>
    </div>
  );
}
