import Image from "next/image";
import { Link } from "@/i18n/navigation";
import type { Product } from "@/types/api";

interface EditorialHeroProps {
  headline: string;
  products: Product[];
}

// Multi-panel editorial banner (DESIGN.md § 5 "Hero (revisi)") — replaces the
// old headline+search full-height hero. Panels use real trending product
// photos rather than static campaign photography, so the hero stays live
// content instead of a stock asset that goes stale.
export function EditorialHero({ headline, products }: EditorialHeroProps) {
  const panels = products.slice(0, 3);
  if (panels.length === 0) return null;

  return (
    <section className="grid grid-cols-1 gap-px overflow-hidden bg-border md:h-[480px] md:grid-cols-3">
      {panels.map((product, i) => (
        <Link
          key={product.id}
          href={`/produk/${product.slug}`}
          className="group relative h-[280px] overflow-hidden bg-surface md:h-full"
        >
          {product.image_url ? (
            <Image
              src={product.image_url}
              alt={product.name}
              fill
              priority={i === 0}
              sizes="(max-width: 768px) 100vw, 34vw"
              className="object-cover grayscale transition-transform duration-700 group-hover:scale-105"
            />
          ) : null}
          <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-transparent" />

          {i === 0 && (
            <h1 className="absolute left-4 top-4 max-w-[85%] font-display text-2xl font-bold leading-[1.05] tracking-[-0.01em] text-white md:left-6 md:top-6 md:text-4xl">
              {headline}
            </h1>
          )}

          <div className="absolute bottom-4 left-4 right-4 md:bottom-6 md:left-6 md:right-6">
            <p className="text-xs font-medium uppercase tracking-[0.01em] text-white/70">
              {product.brand_name}
            </p>
            <p className="font-display text-lg font-semibold text-white md:text-xl">{product.name}</p>
          </div>
        </Link>
      ))}
    </section>
  );
}
