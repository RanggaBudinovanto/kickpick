import { Link } from "@/i18n/navigation";
import { ProductCard } from "@/components/product/product-card";
import type { Product } from "@/types/api";

export function TrendingSpotlight({ products }: { products: Product[] }) {
  if (!products || products.length === 0) return null;

  return (
    <section className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
      <div className="mb-6 flex items-end justify-between border-b border-border/60 pb-4">
        <div>
          <h2 className="font-display text-2xl font-bold tracking-tight text-foreground md:text-4xl">
            Trending Spotlight
          </h2>
          <p className="mt-1 text-xs text-muted md:text-sm">
            Discover the most searched and talked-about sneakers in the community this week.
          </p>
        </div>
        <Link href="/cari" className="text-sm font-medium underline text-foreground/80 hover:text-foreground">
          Explore All Trending →
        </Link>
      </div>

      {/* Product Cards Grid */}
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
        {products.slice(0, 5).map((product) => (
          <ProductCard key={product.id} product={product} badge="Trending" />
        ))}
      </div>
    </section>
  );
}
