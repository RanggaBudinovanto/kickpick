import { IconTrendingDown } from "@tabler/icons-react";
import { ProductCard } from "@/components/product/product-card";
import type { Product } from "@/types/api";

export function PriceDropSpotlight({ products }: { products: Product[] }) {
  if (!products || products.length === 0) return null;

  return (
    <section className="mx-auto max-w-[1400px] px-4 py-12 md:px-6">
      {/* Deals Header Banner */}
      <div className="mb-8 flex flex-col justify-between gap-4 rounded-[var(--radius-card)] border border-emerald-500/30 bg-emerald-950/20 p-6 backdrop-blur-sm md:flex-row md:items-center">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 items-center justify-center rounded-[var(--radius-control)] border border-emerald-500/30 bg-emerald-500/10 text-emerald-400">
            <IconTrendingDown size={28} />
          </div>
          <div>
            <h2 className="font-display text-2xl font-bold tracking-tight text-foreground md:text-3xl">
              Latest Price Drops
            </h2>
            <p className="text-xs text-muted md:text-sm">
              Updated automatically. Compare verified price drops from official stores today.
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2 self-start rounded-[var(--radius-control)] border border-emerald-500/40 bg-emerald-500/10 px-3.5 py-1.5 text-xs font-semibold text-emerald-400 md:self-auto">
          <span>Best Savings This Week</span>
        </div>
      </div>

      {/* Product Cards Grid */}
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
        {products.slice(0, 5).map((product) => (
          <ProductCard
            key={product.id}
            product={product}
            badge={`-${product.drop_percent ?? 15}% OFF`}
          />
        ))}
      </div>
    </section>
  );
}
