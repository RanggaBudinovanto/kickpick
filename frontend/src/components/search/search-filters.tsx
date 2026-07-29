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
  const minPrice = searchParams.get("min_price") ?? "";
  const maxPrice = searchParams.get("max_price") ?? "";

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
    <div className="rounded-[var(--radius-card)] border border-border/80 bg-card p-5 shadow-sm space-y-6">
      <div className="flex items-center justify-between border-b border-border/60 pb-3">
        <h3 className="font-display text-base font-bold tracking-tight text-foreground">Filter</h3>
        <Button variant="ghost" size="sm" className="h-7 text-xs text-muted hover:text-foreground p-0" onClick={resetFilters}>
          Reset
        </Button>
      </div>

      {/* Category Filter */}
      <div>
        <label
          htmlFor="filter-kategori"
          className="mb-2 block text-xs font-mono font-bold uppercase tracking-wider text-muted"
        >
          Kategori
        </label>
        <select
          id="filter-kategori"
          value={kategori}
          onChange={(e) => updateFilter("kategori", e.target.value)}
          className="h-10 w-full rounded-[var(--radius-control)] border border-border bg-background px-3 text-sm focus:outline-none focus:ring-1 focus:ring-foreground"
        >
          <option value="">Semua Kategori</option>
          {CATEGORIES.map((c) => (
            <option key={c} value={c}>
              {c.charAt(0).toUpperCase() + c.slice(1)}
            </option>
          ))}
        </select>
      </div>

      {/* Price Range Filter */}
      <div>
        <label className="mb-2 block text-xs font-mono font-bold uppercase tracking-wider text-muted">
          Rentang Harga (IDR)
        </label>
        <div className="flex items-center gap-2">
          <input
            type="number"
            placeholder="Min"
            value={minPrice}
            onChange={(e) => updateFilter("min_price", e.target.value)}
            className="h-9 w-full rounded-[var(--radius-control)] border border-border bg-background px-3 text-xs focus:outline-none focus:ring-1 focus:ring-foreground"
          />
          <span className="text-muted text-xs">-</span>
          <input
            type="number"
            placeholder="Max"
            value={maxPrice}
            onChange={(e) => updateFilter("max_price", e.target.value)}
            className="h-9 w-full rounded-[var(--radius-control)] border border-border bg-background px-3 text-xs focus:outline-none focus:ring-1 focus:ring-foreground"
          />
        </div>
      </div>

      {/* Brand Filter */}
      <div>
        <label className="mb-2 block text-xs font-mono font-bold uppercase tracking-wider text-muted">
          Brand
        </label>
        <div className="flex max-h-52 flex-col gap-2 overflow-y-auto pr-1">
          {brands.map((b) => (
            <label key={b.id} className="flex cursor-pointer items-center gap-2 text-sm text-foreground/90 hover:text-foreground">
              <input
                type="checkbox"
                className="h-4 w-4 rounded border-border text-foreground accent-foreground focus:ring-0"
                checked={selectedBrandIds.includes(b.id)}
                onChange={(e) => toggleBrand(b.id, e.target.checked)}
              />
              {b.name}
            </label>
          ))}
        </div>
      </div>

      {/* Special Collection */}
      <div className="border-t border-border/60 pt-4">
        <label className="flex cursor-pointer items-center gap-2 text-sm font-medium text-foreground">
          <input
            type="checkbox"
            className="h-4 w-4 rounded border-border accent-foreground"
            checked={filter === "rare"}
            onChange={(e) => updateFilter("filter", e.target.checked ? "rare" : "")}
          />
          Rare & Limited Only
        </label>
      </div>
    </div>
  );
}
