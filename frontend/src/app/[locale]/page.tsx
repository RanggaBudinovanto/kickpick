import Image from "next/image";
import { getTranslations } from "next-intl/server";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/button";
import { HeroSearch } from "@/components/home/hero-search";
import { ProductCard } from "@/components/product/product-card";
import { getBrands, getPriceDrops, getProducts, getTrending } from "@/lib/server-api";
import type { Product } from "@/types/api";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL ?? "http://localhost:3000";

const organizationStructuredData = {
  "@context": "https://schema.org",
  "@type": "Organization",
  name: "KickPick",
  url: APP_URL,
  description: "Platform pencarian dan perbandingan harga sepatu lintas brand.",
};

// Homepage sections are independent and non-critical (Section 8 PRD: "Section
// disembunyikan sepenuhnya jika data kosong" — a failed section shouldn't
// take down the rest of the page), so each fetch is isolated and swallowed to
// an empty result on failure rather than one bad section throwing for all.
async function safeProducts(promise: Promise<{ data: Product[] }>): Promise<Product[]> {
  try {
    return (await promise).data;
  } catch {
    return [];
  }
}

async function safeBrands() {
  try {
    return (await getBrands()).data;
  } catch {
    return [];
  }
}

export default async function HomePage() {
  const t = await getTranslations("Home");
  const tFooter = await getTranslations("Footer");

  const [brands, bestSellers, trending, rare, priceDrops] = await Promise.all([
    safeBrands(),
    safeProducts(getProducts({ limit: 8 })),
    getTrending(8),
    safeProducts(getProducts({ filter: "rare", limit: 8 })),
    getPriceDrops(8),
  ]);

  return (
    <>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{
          __html: JSON.stringify(organizationStructuredData).replace(/</g, "\\u003c"),
        }}
      />

      <section className="mx-auto grid max-w-[1400px] gap-10 px-4 py-16 md:grid-cols-[60%_40%] md:px-6 md:py-24">
        <div className="flex flex-col justify-center gap-6">
          <h1 className="font-display text-4xl font-bold leading-[1.05] tracking-[-0.01em] md:text-[56px]">
            {t("heroHeadline")}
          </h1>
          <p className="text-sm text-muted">{t("heroSubtext", { count: 40 })}</p>
          <HeroSearch />
          <div className="flex gap-3">
            <Link href="/brand">
              <Button variant="secondary">{t("browseBrands")}</Button>
            </Link>
          </div>
        </div>
        <div className="hidden items-center justify-center rounded-[var(--radius-card)] border border-border bg-surface md:flex">
          <span className="text-sm text-muted">Visual pendukung (brand strip / grafik tren)</span>
        </div>
      </section>

      {brands.length > 0 && (
        <div className="mx-auto max-w-[1400px] px-4 py-6 md:px-6">
          <div className="flex items-center gap-8 overflow-x-auto pb-2">
            {brands.map((brand) => (
              <Link
                key={brand.id}
                href={`/brand/${brand.slug}`}
                className="shrink-0 text-sm font-medium text-muted hover:text-foreground"
              >
                {brand.logo_url ? (
                  <Image
                    src={brand.logo_url}
                    alt={brand.name}
                    width={96}
                    height={32}
                    className="h-8 w-auto object-contain"
                  />
                ) : (
                  brand.name
                )}
              </Link>
            ))}
            <Link href="/brand" className="shrink-0 text-sm font-medium underline">
              {tFooter("brands")}
            </Link>
          </div>
        </div>
      )}

      {bestSellers.length > 0 && (
        <section className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
          <div className="mb-6 flex items-end justify-between">
            <h2 className="font-display text-2xl font-semibold tracking-[-0.005em] md:text-4xl">
              Best Seller
            </h2>
            <Link href="/cari" className="text-sm font-medium underline">
              Lihat semua
            </Link>
          </div>
          <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
            {bestSellers.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        </section>
      )}

      {trending.length > 0 && (
        <section className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
          <h2 className="mb-6 font-display text-2xl font-semibold tracking-[-0.005em] md:text-4xl">
            Lagi Trending
          </h2>
          <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
            {trending.map((product) => (
              <ProductCard key={product.id} product={product} badge="Trending" />
            ))}
          </div>
        </section>
      )}

      {rare.length > 0 && (
        <section className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
          <div className="mb-6 flex items-end justify-between">
            <h2 className="font-display text-2xl font-semibold tracking-[-0.005em] md:text-4xl">
              Rare & Limited
            </h2>
            <Link href="/cari?filter=rare" className="text-sm font-medium underline">
              Lihat semua
            </Link>
          </div>
          <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
            {rare.map((product) => (
              <ProductCard key={product.id} product={product} badge="Limited" />
            ))}
          </div>
        </section>
      )}

      {priceDrops.length > 0 && (
        <section className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
          <h2 className="mb-6 font-display text-2xl font-semibold tracking-[-0.005em] md:text-4xl">
            Turun Harga
          </h2>
          <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
            {priceDrops.map((product) => (
              <ProductCard key={product.id} product={product} badge={`Turun ${product.drop_percent}%`} />
            ))}
          </div>
        </section>
      )}
    </>
  );
}
