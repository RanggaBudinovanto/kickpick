import Image from "next/image";
import { IconStar, IconStarFilled } from "@tabler/icons-react";
import { Link } from "@/i18n/navigation";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { PriceDisplay } from "@/components/product/price-display";
import type { Product } from "@/types/api";

export function ProductCard({ product, badge }: { product: Product; badge?: string }) {
  return (
    <Link href={`/produk/${product.slug}`}>
      <Card className="group overflow-hidden transition-transform hover:-translate-y-0.5">
        <div className="relative aspect-square bg-surface">
          {product.image_url ? (
            <Image
              src={product.image_url}
              alt={product.name}
              fill
              className="object-cover"
              sizes="(max-width: 768px) 50vw, 25vw"
            />
          ) : (
            <div className="flex h-full items-center justify-center text-xs text-muted">
              Tidak ada foto
            </div>
          )}
          {badge && (
            <span className="absolute right-2 top-2">
              <Badge variant="strong">{badge}</Badge>
            </span>
          )}
        </div>
        <div className="p-3">
          <p className="text-xs font-medium uppercase tracking-[0.01em] text-muted">
            {product.brand_name}
          </p>
          <h3 className="mb-1 truncate text-sm font-medium">{product.name}</h3>
          <p className="font-mono text-base font-bold tabular-nums">
            <PriceDisplay min={product.min_price} max={product.max_price} />
          </p>
          <p className="mt-1 flex items-center gap-1 text-xs text-muted">
            {product.rating > 0 ? (
              <>
                <IconStarFilled size={12} />
                {product.rating.toFixed(1)}
              </>
            ) : (
              <>
                <IconStar size={12} />
                Belum ada rating
              </>
            )}
          </p>
        </div>
      </Card>
    </Link>
  );
}
