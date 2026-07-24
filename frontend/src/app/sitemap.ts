import type { MetadataRoute } from "next";
import { API_URL } from "@/lib/api-client";
import { routing } from "@/i18n/routing";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL ?? "http://localhost:3000";

async function safeJson<T>(url: string): Promise<T | null> {
  try {
    const res = await fetch(url, { cache: "no-store" });
    if (!res.ok) return null;
    return res.json();
  } catch {
    return null;
  }
}

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const entries: MetadataRoute.Sitemap = [];

  const staticPaths = ["", "/cari", "/brand", "/tentang", "/privasi", "/disclosure"];
  for (const locale of routing.locales) {
    for (const path of staticPaths) {
      entries.push({ url: `${APP_URL}/${locale}${path}`, changeFrequency: "daily" });
    }
  }

  const products = await safeJson<{ data: { slug: string }[] }>(`${API_URL}/api/products?limit=100`);
  const brands = await safeJson<{ data: { slug: string }[] }>(`${API_URL}/api/brands`);

  for (const locale of routing.locales) {
    for (const p of products?.data ?? []) {
      entries.push({ url: `${APP_URL}/${locale}/produk/${p.slug}`, changeFrequency: "hourly" });
    }
    for (const b of brands?.data ?? []) {
      entries.push({ url: `${APP_URL}/${locale}/brand/${b.slug}`, changeFrequency: "daily" });
    }
  }

  return entries;
}
