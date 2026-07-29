import Image from "next/image";
import { getTranslations } from "next-intl/server";
import { Link } from "@/i18n/navigation";
import { HeroCarousel } from "@/components/home/hero-carousel";
import { CuratedBrandGrid, type CuratedTile } from "@/components/home/curated-brand-grid";
import { TrendingSpotlight } from "@/components/home/trending-spotlight";
import { RareVaultSection } from "@/components/home/rare-vault-section";
import { PriceDropSpotlight } from "@/components/home/price-drop-spotlight";
import { ShopByCategory } from "@/components/home/shop-by-category";
import { NewArrivalsSection } from "@/components/home/new-arrivals-section";
import { ProductCard } from "@/components/product/product-card";
import { getBrands, getPriceDrops, getProducts, getTrending } from "@/lib/server-api";
import type { Product, Brand } from "@/types/api";

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

/**
 * Diversify best-sellers so no single brand dominates.
 * - Sorts products: international brands (is_local=false) first
 * - Limits to max `perBrand` products per brand slug
 * - Returns up to `limit` products
 */
function diversifyBestSellers(
  products: Product[],
  brands: Brand[],
  limit = 8,
  perBrand = 2,
): Product[] {
  const localSlugs = new Set(brands.filter((b) => b.is_local).map((b) => b.slug));
  // International first, then local
  const sorted = [
    ...products.filter((p) => !localSlugs.has(p.brand_slug)),
    ...products.filter((p) => localSlugs.has(p.brand_slug)),
  ];
  const countPerBrand: Record<string, number> = {};
  const result: Product[] = [];
  for (const product of sorted) {
    const count = countPerBrand[product.brand_slug] ?? 0;
    if (count >= perBrand) continue;
    countPerBrand[product.brand_slug] = count + 1;
    result.push(product);
    if (result.length >= limit) break;
  }
  return result;
}

// placehold.co URLs come from the old demo/seed products (Nike Air Force 1
// 07, Adidas Samba OG, etc.) — they render as a flat gray square with the
// slug baked into the image as text, which looks broken in a hero-sized
// spot. Real scraped products all have real product photos, so this simply
// skips the placeholder ones rather than trying to special-case them.
function hasRealImage(product: Product): boolean {
  return Boolean(product.image_url) && !product.image_url.includes("placehold.co");
}

const FEATURED_BRAND_FALLBACKS: CuratedTile[] = [
  {
    brandName: "Nike",
    brandSlug: "nike",
    imageUrl: "/brands/nike.png",
    logoUrl: "/brands/nike.svg",
    accentColor: "#f97316",
  },
  {
    brandName: "Adidas",
    brandSlug: "adidas",
    imageUrl: "/hero/slide-1.png",
    logoUrl: "/brands/adidas.svg",
    accentColor: "#3b82f6",
  },
  {
    brandName: "Puma",
    brandSlug: "puma",
    imageUrl: "/hero/slide-2.png",
    logoUrl: "/brands/puma.svg",
    accentColor: "#ef4444",
  },
  {
    brandName: "Jordan",
    brandSlug: "jordan",
    imageUrl: "/hero/slide-3.png",
    logoUrl: "/brands/jordan.svg",
    accentColor: "#ef4444",
  },
  {
    brandName: "New Balance",
    brandSlug: "new-balance",
    imageUrl: "/hero/trending-banner.png",
    logoUrl: "/brands/new-balance.svg",
    accentColor: "#6366f1",
  },
  {
    brandName: "Vans",
    brandSlug: "vans",
    imageUrl: "/hero/category-lifestyle.png",
    logoUrl: "/brands/vans.svg",
    accentColor: "#a3a3a3",
  },
];

// Local mapping: brand slug → official logo SVG + accent colour.
// Used to enrich tiles built from product data (which only carry image_url)
// so the curated grid always shows official logos regardless of data source.
const BRAND_LOGO_MAP: Record<string, { logoUrl: string; accentColor: string }> = {
  "nike":        { logoUrl: "/brands/nike.svg",        accentColor: "#f97316" },
  "adidas":      { logoUrl: "/brands/adidas.svg",      accentColor: "#3b82f6" },
  "puma":        { logoUrl: "/brands/puma.svg",        accentColor: "#ef4444" },
  "jordan":      { logoUrl: "/brands/jordan.svg",      accentColor: "#ef4444" },
  "new-balance": { logoUrl: "/brands/new-balance.svg", accentColor: "#6366f1" },
  "vans":        { logoUrl: "/brands/vans.svg",        accentColor: "#a3a3a3" },
  "asics":       { logoUrl: "/brands/asics.svg",       accentColor: "#3b82f6" },
  "mizuno":      { logoUrl: "/brands/mizuno.svg",      accentColor: "#1e40af" },
  "crocs":       { logoUrl: "/brands/crocs.svg",       accentColor: "#22c55e" },
  "on":          { logoUrl: "/brands/on.svg",          accentColor: "#f59e0b" },
  "under-armour":{ logoUrl: "/brands/under-armour.svg",accentColor: "#2563eb" },
  "ventela":     { logoUrl: "/brands/ventela.svg",     accentColor: "#8b5cf6" },
  "aerostreet":  { logoUrl: "/brands/aerostreet.svg",  accentColor: "#10b981" },
  "compass":     { logoUrl: "/brands/compass.svg",     accentColor: "#64748b" },
  "brodo":       { logoUrl: "/brands/brodo.svg",       accentColor: "#92400e" },
  "geoff-max":   { logoUrl: "/brands/geoffmax.svg",    accentColor: "#7c3aed" },
};

// Builds curated tiles from already-fetched product lists, enriched with
// brand logo data. Falls back to FEATURED_BRAND_FALLBACKS if fewer than max.
function buildCuratedTiles(
  productLists: Product[][],
  brandLogoUrlMap: Record<string, string>,
  max: number = 6,
): CuratedTile[] {
  const seen = new Set<string>();
  const tiles: CuratedTile[] = [];

  for (const list of productLists) {
    for (const product of list) {
      if (seen.has(product.brand_slug) || !hasRealImage(product)) continue;
      seen.add(product.brand_slug);
      // Prefer API logo_url (from /brands), fallback to local SVG map
      const apiLogoUrl = brandLogoUrlMap[product.brand_slug];
      const localLogo = BRAND_LOGO_MAP[product.brand_slug];
      const logoUrl = apiLogoUrl || localLogo?.logoUrl;
      tiles.push({
        brandName: product.brand_name,
        brandSlug: product.brand_slug,
        imageUrl: product.image_url,
        logoUrl,
        accentColor: localLogo?.accentColor,
        wide: tiles.length === 0 || tiles.length === 5,
      });
      if (tiles.length >= max) return tiles;
    }
  }

  for (const fallback of FEATURED_BRAND_FALLBACKS) {
    if (tiles.length >= max) break;
    if (!seen.has(fallback.brandSlug)) {
      seen.add(fallback.brandSlug);
      tiles.push({
        ...fallback,
        wide: tiles.length === 0 || tiles.length === 5,
      });
    }
  }

  return tiles;
}

export default async function HomePage() {
  const t = await getTranslations("Home");
  const tFooter = await getTranslations("Footer");

  const [brands, rawBestSellers, trending, rare, priceDrops, newArrivals] = await Promise.all([
    safeBrands(),
    safeProducts(getProducts({ limit: 24 })),   // fetch more so diversify has enough variety
    getTrending(8),
    safeProducts(getProducts({ filter: "rare", limit: 8 })),
    getPriceDrops(8),
    safeProducts(getProducts({ limit: 8 })),
  ]);

  // Max 2 per brand, international brands shown first
  const bestSellers = diversifyBestSellers(rawBestSellers, brands, 8, 2);

  // Build slug → logo_url map from fetched brands (used to enrich curated tiles)
  const brandLogoUrlMap = Object.fromEntries(
    brands.filter((b) => b.logo_url).map((b) => [b.slug, b.logo_url]),
  );

  const curatedTiles = buildCuratedTiles([trending, bestSellers, rare, priceDrops], brandLogoUrlMap, 6);

  return (
    <>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{
          __html: JSON.stringify(organizationStructuredData).replace(/</g, "\\u003c"),
        }}
      />

      <HeroCarousel />

      {brands.length > 0 && (
        <div className="mx-auto max-w-[1400px] px-4 py-2 md:px-6">
          <div className="flex items-center gap-10 overflow-x-auto py-4 [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none] border-y border-border/40">
            {brands.map((brand) => (
              <Link
                key={brand.id}
                href={`/brand/${brand.slug}`}
                className="shrink-0 flex items-center justify-center h-8"
              >
                {brand.logo_url ? (
                  <Image
                    src={brand.logo_url}
                    alt={brand.name}
                    width={120}
                    height={40}
                    unoptimized
                    className="h-7 w-auto max-h-7 object-contain dark:invert opacity-75 hover:opacity-100 transition-opacity"
                  />
                ) : (
                  <span className="text-xs font-bold uppercase tracking-wider text-muted hover:text-foreground">
                    {brand.name}
                  </span>
                )}
              </Link>
            ))}
            <Link
              href="/brand"
              className="shrink-0 text-xs font-bold uppercase tracking-wider text-muted hover:text-foreground underline underline-offset-4"
            >
              {tFooter("brands")} →
            </Link>
          </div>
        </div>
      )}

      <CuratedBrandGrid title={t("curatedTitle")} tiles={curatedTiles} />

      {bestSellers.length > 0 && (
        <section className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
          <div className="mb-6 flex items-end justify-between border-b border-border/60 pb-4">
            <div>
              <h2 className="font-display text-2xl font-bold tracking-tight text-foreground md:text-4xl">
                Best Sellers
              </h2>
              <p className="mt-1 text-xs text-muted md:text-sm">
                Most wanted sneakers and top community picks this week.
              </p>
            </div>
            <Link href="/cari" className="text-sm font-medium underline text-foreground/80 hover:text-foreground">
              Explore All →
            </Link>
          </div>
          <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
            {bestSellers.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        </section>
      )}

      <TrendingSpotlight products={trending} />

      {/* ─── Divider ─── */}
      <div className="mx-auto max-w-[1400px] px-4 md:px-6">
        <div className="border-t border-border/40" />
      </div>

      <ShopByCategory />

      {/* ─── Divider ─── */}
      <div className="mx-auto max-w-[1400px] px-4 md:px-6">
        <div className="border-t border-border/40" />
      </div>

      <NewArrivalsSection products={newArrivals} />

      {/* ─── Divider ─── */}
      <div className="mx-auto max-w-[1400px] px-4 md:px-6">
        <div className="border-t border-border/40" />
      </div>

      <RareVaultSection products={rare} />

      <PriceDropSpotlight products={priceDrops} />
    </>
  );
}
