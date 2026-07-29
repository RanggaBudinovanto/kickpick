import Image from "next/image";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/button";
import { ProductCard } from "@/components/product/product-card";
import type { Product } from "@/types/api";

export function TrendingSpotlight({ products }: { products: Product[] }) {
  if (!products || products.length === 0) return null;

  return (
    <section className="mx-auto max-w-[1400px] px-4 py-12 md:px-6">
      {/* Editorial Feature Banner */}
      <div className="relative mb-8 overflow-hidden rounded-[var(--radius-card)] border border-border/80 bg-neutral-950 text-white shadow-xl min-h-[320px] flex items-center">
        <div className="absolute inset-0 z-0">
          <Image
            src="/hero/trending-banner.png"
            alt="Trending Spotlight"
            fill
            className="object-cover object-right md:object-center opacity-85 transition-transform duration-700 hover:scale-105"
          />
          <div className="absolute inset-0 bg-gradient-to-r from-neutral-950 via-neutral-950/70 to-transparent" />
        </div>

        <div className="relative z-10 flex flex-col items-start justify-center p-8 md:p-12 max-w-xl">
          <h2 className="mb-4 font-display text-3xl font-bold tracking-tight text-white md:text-5xl">
            Street Culture & Performance
          </h2>
          <p className="mb-6 text-sm text-neutral-300 md:text-base leading-relaxed">
            Discover the most searched and talked-about sneakers in the community this week.
          </p>
          <Link href="/cari">
            <Button className="bg-white text-neutral-950 hover:bg-neutral-200 font-bold px-6 py-3 shadow-lg border-0">
              Explore All Trending →
            </Button>
          </Link>
        </div>
      </div>

      {/* Product Cards Grid */}
      <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
        {products.slice(0, 4).map((product) => (
          <ProductCard key={product.id} product={product} badge="Trending" />
        ))}
      </div>
    </section>
  );
}
