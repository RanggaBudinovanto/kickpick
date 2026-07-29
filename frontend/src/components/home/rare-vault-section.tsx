import { Link } from "@/i18n/navigation";
import { ProductCard } from "@/components/product/product-card";
import type { Product } from "@/types/api";

export function RareVaultSection({ products }: { products: Product[] }) {
  if (!products || products.length === 0) return null;

  return (
    <section className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
      <div className="rounded-[var(--radius-card)] border border-border/80 bg-card p-6 md:p-8 shadow-sm">
        <div className="mb-6 flex flex-col items-start justify-between gap-4 border-b border-border/60 pb-4 md:flex-row md:items-end">
          <div>
            <h2 className="font-display text-2xl font-bold tracking-tight text-foreground md:text-4xl">
              Rare & Limited Drops
            </h2>
            <p className="mt-1 text-xs text-muted md:text-sm">
              Exclusive releases and limited edition grails with minimal market availability.
            </p>
          </div>

          <Link
            href="/cari?filter=rare"
            className="shrink-0 text-sm font-medium underline text-foreground/80 hover:text-foreground"
          >
            Explore All Rare Drops →
          </Link>
        </div>

        {/* Product Cards Grid */}
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
          {products.slice(0, 5).map((product) => (
            <ProductCard key={product.id} product={product} badge="Limited" />
          ))}
        </div>
      </div>
    </section>
  );
}
