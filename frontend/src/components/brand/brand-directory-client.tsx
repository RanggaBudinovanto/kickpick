"use client";

import { useState } from "react";
import Image from "next/image";
import { Link } from "@/i18n/navigation";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import type { Brand } from "@/types/api";

// Local official logo assets already bundled in /public/brands/.
// These take priority over the database logo_url so the grid always
// shows the real brand mark regardless of what the API stores.
const LOCAL_LOGOS: Record<string, string> = {
  "nike":          "/brands/nike.svg",
  "adidas":        "/brands/adidas.svg",
  "puma":          "/brands/puma.svg",
  "jordan":        "/brands/jordan.svg",
  "new-balance":   "/brands/new-balance.svg",
  "vans":          "/brands/vans.svg",
  "asics":         "/brands/asics.svg",
  "mizuno":        "/brands/mizuno.svg",
  "crocs":         "/brands/crocs.svg",
  "on":            "/brands/on.svg",
  "under-armour":  "/brands/under-armour.svg",
  "ventela":       "/brands/ventela.svg",
  "aerostreet":    "/brands/aerostreet.svg",
  "compass":       "/brands/compass.svg",
  "brodo":         "/brands/brodo.svg",
  "geoff-max":     "/brands/geoffmax.svg",
};

function getLogoUrl(brand: Brand): string | null {
  return LOCAL_LOGOS[brand.slug] ?? brand.logo_url ?? null;
}

export function BrandDirectoryClient({ brands }: { brands: Brand[] }) {
  const [query, setQuery] = useState("");

  const filtered = brands.filter((b) =>
    b.name.toLowerCase().includes(query.toLowerCase()),
  );

  return (
    <div>
      <div className="mb-8 max-w-sm">
        <Input
          placeholder="Search brands"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
      </div>

      <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
        {filtered.map((brand) => {
          const logoUrl = getLogoUrl(brand);
          return (
            <Link key={brand.id} href={`/brand/${brand.slug}`}>
              <Card className="flex h-28 items-center justify-center p-6 text-center transition-all duration-300 hover:-translate-y-1 hover:shadow-lg bg-card/40 backdrop-blur-sm border border-border/80 hover:border-border-hover">
                {logoUrl ? (
                  <Image
                    src={logoUrl}
                    alt={brand.name}
                    width={140}
                    height={70}
                    unoptimized
                    className="h-12 w-auto object-contain dark:invert"
                  />
                ) : (
                  <span className="font-display text-lg font-bold tracking-tight text-foreground/80">
                    {brand.name}
                  </span>
                )}
              </Card>
            </Link>
          );
        })}
      </div>

      {filtered.length === 0 && (
        <div className="py-12 text-center text-sm text-muted">
          Tidak ada brand yang cocok dengan pencarian &quot;{query}&quot;
        </div>
      )}
    </div>
  );
}
