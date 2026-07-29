import type { Metadata } from "next";
import { notFound } from "next/navigation";
import Image from "next/image";
import { getProductDetail } from "@/lib/server-api";
import { PriceComparison } from "@/components/product/price-comparison";
import { PriceHistoryChart } from "@/components/product/price-history-chart";
import { SizeConverter } from "@/components/product/size-converter";
import { ReviewSection } from "@/components/product/review-section";
import { WishlistButton } from "@/components/product/wishlist-button";
import { Badge } from "@/components/ui/badge";
import { formatPriceRange } from "@/lib/currency";

interface PageProps {
  params: Promise<{ slug: string; locale: string }>;
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { slug, locale } = await params;
  const data = await getProductDetail(slug);

  if (!data) {
    return { title: "Product not found | KickPick" };
  }

  const { product, offers } = data;
  const prices = offers.map((o) => o.price);
  const description =
    prices.length > 0
      ? `Compare prices for ${product.name} from ${offers.length} stores, starting at ${formatPriceRange(Math.min(...prices), Math.max(...prices))}.`
      : `Compare prices for ${product.name} from various stores on KickPick.`;

  return {
    title: `${product.name} - Compare Prices | KickPick`,
    description,
    alternates: {
      languages: {
        id: `/id/produk/${slug}`,
        en: `/en/produk/${slug}`,
        "x-default": `/en/produk/${slug}`,
      },
    },
    openGraph: {
      title: product.name,
      description,
      images: product.image_url ? [product.image_url] : undefined,
    },
  };
}

export default async function ProductDetailPage({ params }: PageProps) {
  const { slug } = await params;
  const data = await getProductDetail(slug);

  if (!data) {
    notFound();
  }

  const { product, offers, reviews } = data;
  const prices = offers.map((o) => o.price);

  const structuredData = {
    "@context": "https://schema.org",
    "@type": "Product",
    name: product.name,
    brand: { "@type": "Brand", name: product.brand_name },
    image: product.image_url || undefined,
    offers:
      prices.length > 0
        ? {
            "@type": "AggregateOffer",
            priceCurrency: "IDR",
            lowPrice: Math.min(...prices),
            highPrice: Math.max(...prices),
            offerCount: offers.length,
          }
        : undefined,
  };

  return (
    <div className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(structuredData).replace(/</g, "\\u003c") }}
      />

      <div className="mb-10 grid gap-8 md:grid-cols-2">
        <div className="relative aspect-square overflow-hidden rounded-[var(--radius-card)] bg-surface">
          {product.image_url ? (
            <Image
              src={product.image_url}
              alt={product.name}
              fill
              preload
              fetchPriority="high"
              className="object-cover"
              sizes="50vw"
            />
          ) : null}
          {product.is_limited && (
            <span className="absolute right-3 top-3">
              <Badge variant="strong">Limited</Badge>
            </span>
          )}
        </div>

        <div>
          <p className="mb-1 text-sm font-medium uppercase tracking-[0.01em] text-muted">
            {product.brand_name}
          </p>
          <h1 className="mb-4 font-display text-3xl font-bold tracking-[-0.01em] md:text-4xl">
            {product.name}
          </h1>

          <div className="mb-6">
            <WishlistButton productId={product.id} />
          </div>

          <PriceComparison offers={offers} />
        </div>
      </div>

      <div className="grid gap-10 md:grid-cols-2">
        <PriceHistoryChart slug={slug} />
        <SizeConverter productSlug={slug} />
      </div>

      <div className="mt-10 max-w-2xl">
        <ReviewSection productId={product.id} productSlug={slug} reviews={reviews} />
      </div>
    </div>
  );
}
