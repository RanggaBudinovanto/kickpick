import { API_URL } from "@/lib/api-client";
import type { Brand, Product, ProductDetail } from "@/types/api";

// Server-side fetch (used in generateMetadata + Server Components) so search engine
// crawlers get real HTML on first load, instead of an empty CSR shell.
// Returns null only for a genuine 404 (triggers notFound()); throws for any
// other failure so it doesn't get misreported as "product doesn't exist".
// ISR: cached up to 1 hour (Section 19 PRD), invalidated on demand via
// /api/revalidate + revalidateTag("product:<slug>") — called by the frontend
// after a review submit and by the scraper worker after a price update, so
// the cache window doesn't mean stale prices/reviews in practice.
export async function getProductDetail(slug: string): Promise<ProductDetail | null> {
  const res = await fetch(`${API_URL}/api/products/${slug}`, {
    next: { revalidate: 3600, tags: [`product:${slug}`] },
  });
  if (res.status === 404) return null;
  if (!res.ok) throw new Error(`Failed to load product ${slug} (${res.status})`);
  return res.json();
}

export interface GetProductsParams {
  q?: string;
  kategori?: string;
  brand_ids?: string;
  filter?: string;
  min_price?: string | number;
  max_price?: string | number;
  limit?: number;
}

// Throws on failure (rather than swallowing to an empty list) so the /cari
// route's error.tsx boundary can show a distinct "failed to load" + retry
// state instead of it looking identical to a genuine zero-result search.
export async function getProducts(
  params: GetProductsParams = {},
): Promise<{ data: Product[]; page: number; limit: number }> {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== "") search.set(key, String(value));
  }
  const qs = search.toString();

  const res = await fetch(`${API_URL}/api/products${qs ? `?${qs}` : ""}`, {
    cache: "no-store",
  });
  if (!res.ok) throw new Error(`Failed to load products (${res.status})`);
  return res.json();
}

export async function getBrands(): Promise<{ data: Brand[] }> {
  const res = await fetch(`${API_URL}/api/brands`, { cache: "no-store" });
  if (!res.ok) throw new Error(`Failed to load brands (${res.status})`);
  return res.json();
}

export async function getBrandDetail(
  slug: string,
): Promise<{ brand: Brand; products: Product[] } | null> {
  const res = await fetch(`${API_URL}/api/brands/${slug}`, { cache: "no-store" });
  if (res.status === 404) return null;
  if (!res.ok) throw new Error(`Failed to load brand ${slug} (${res.status})`);
  return res.json();
}

export async function getTrending(limit = 8): Promise<Product[]> {
  try {
    const res = await fetch(`${API_URL}/api/products/trending?limit=${limit}`, { cache: "no-store" });
    if (!res.ok) return [];
    const data: { data: Product[] } = await res.json();
    return data.data;
  } catch {
    return [];
  }
}

export async function getPriceDrops(
  limit = 8,
): Promise<(Product & { drop_percent: number })[]> {
  try {
    const res = await fetch(`${API_URL}/api/products/price-drops?limit=${limit}`, { cache: "no-store" });
    if (!res.ok) return [];
    const data: { data: (Product & { drop_percent: number })[] } = await res.json();
    return data.data;
  } catch {
    return [];
  }
}
