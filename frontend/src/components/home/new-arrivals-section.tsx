import { Link } from "@/i18n/navigation";
import { ProductCard } from "@/components/product/product-card";
import type { Product } from "@/types/api";

interface NewArrivalsSectionProps {
  products: Product[];
}

export function NewArrivalsSection({ products }: NewArrivalsSectionProps) {
  if (products.length === 0) return null;

  return (
    <section className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
      <div className="mb-6 flex items-end justify-between border-b border-border/60 pb-4">
        <div>
          <h2 className="font-display text-2xl font-bold tracking-tight text-foreground md:text-4xl">
            New Arrivals
          </h2>
          <p className="mt-1 text-xs text-muted md:text-sm">
            The latest sneakers just added to our price comparison engine.
          </p>
        </div>
        <Link href="/cari" className="text-sm font-medium underline text-foreground/80 hover:text-foreground">
          View All New →
        </Link>
      </div>
      <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
        {products.map((product) => (
          <ProductCard key={product.id} product={product} badge="New" />
        ))}
      </div>
    </section>
  );
}
