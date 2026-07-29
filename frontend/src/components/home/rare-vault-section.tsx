import { Link } from "@/i18n/navigation";
import { ProductCard } from "@/components/product/product-card";
import type { Product } from "@/types/api";

export function RareVaultSection({ products }: { products: Product[] }) {
  if (!products || products.length === 0) return null;

  return (
    <section className="mx-auto max-w-[1400px] px-4 py-12 md:px-6">
      <div className="rounded-[var(--radius-card)] border border-amber-500/30 bg-neutral-950 p-6 md:p-10 text-white shadow-2xl relative overflow-hidden">
        {/* Subtle Vault Background Accent Glow */}
        <div className="pointer-events-none absolute -top-24 -right-24 h-96 w-96 rounded-full bg-amber-500/10 blur-3xl" />
        <div className="pointer-events-none absolute -bottom-24 -left-24 h-96 w-96 rounded-full bg-amber-600/10 blur-3xl" />

        <div className="relative z-10 mb-8 flex flex-col items-start justify-between gap-4 md:flex-row md:items-end">
          <div>
            <h2 className="font-display text-3xl font-bold tracking-tight text-white md:text-4xl">
              The Rare Vault — Rare & Limited Drops
            </h2>
            <p className="mt-2 text-sm text-neutral-400 max-w-lg">
              Exclusive releases and limited edition grails with minimal market availability.
            </p>
          </div>

          <Link
            href="/cari?filter=rare"
            className="shrink-0 text-sm font-semibold text-amber-400 hover:text-amber-300 underline underline-offset-4"
          >
            Explore All Rare Drops →
          </Link>
        </div>

        {/* Product Cards Grid */}
        <div className="relative z-10 grid grid-cols-2 gap-4 md:grid-cols-4">
          {products.slice(0, 4).map((product) => (
            <ProductCard key={product.id} product={product} badge="Limited" />
          ))}
        </div>
      </div>
    </section>
  );
}
