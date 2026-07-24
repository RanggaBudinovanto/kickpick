"use client";

import { useState } from "react";
import Image from "next/image";
import { Link } from "@/i18n/navigation";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import { useBrands } from "@/hooks/use-products";

export default function BrandDirectoryPage() {
  const { data, isLoading } = useBrands();
  const [query, setQuery] = useState("");

  const filtered = data?.data.filter((b) => b.name.toLowerCase().includes(query.toLowerCase())) ?? [];

  return (
    <div className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
      <h1 className="mb-6 font-display text-3xl font-bold tracking-[-0.01em]">Brand</h1>

      <div className="mb-8 max-w-sm">
        <Input placeholder="Cari brand" value={query} onChange={(e) => setQuery(e.target.value)} />
      </div>

      <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
        {isLoading
          ? Array.from({ length: 8 }).map((_, i) => (
              <div key={i} className="h-24 animate-pulse rounded-[var(--radius-card)] bg-surface" />
            ))
          : filtered.map((brand) => (
              <Link key={brand.id} href={`/brand/${brand.slug}`}>
                <Card className="flex h-24 flex-col items-center justify-center gap-2 p-4 text-center transition-transform hover:-translate-y-0.5">
                  {brand.logo_url ? (
                    <Image src={brand.logo_url} alt={brand.name} width={64} height={32} className="h-8 w-auto object-contain" />
                  ) : (
                    <span className="font-display text-lg font-semibold">{brand.name}</span>
                  )}
                </Card>
              </Link>
            ))}
      </div>
    </div>
  );
}
