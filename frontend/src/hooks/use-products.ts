import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api-client";
import type { Brand, PricePoint, Product, ProductDetail } from "@/types/api";

export interface ProductListParams {
  page?: number;
  limit?: number;
  kategori?: string;
  brand_id?: string;
  brand_ids?: string;
  filter?: "rare" | "trending";
  q?: string;
}

function toQueryString(params: ProductListParams) {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== "") search.set(key, String(value));
  }
  const qs = search.toString();
  return qs ? `?${qs}` : "";
}

export function useProducts(params: ProductListParams = {}) {
  return useQuery({
    queryKey: ["products", params],
    queryFn: () =>
      apiFetch<{ data: Product[]; page: number; limit: number }>(
        `/api/products${toQueryString(params)}`,
      ),
  });
}

export function useProductDetail(slug: string) {
  return useQuery({
    queryKey: ["product", slug],
    queryFn: () => apiFetch<ProductDetail>(`/api/products/${slug}`),
    enabled: !!slug,
  });
}

export function usePriceHistory(slug: string, days: 30 | 90 = 30) {
  return useQuery({
    queryKey: ["price-history", slug, days],
    queryFn: () => apiFetch<{ data: PricePoint[] }>(`/api/products/${slug}/price-history?days=${days}`),
    enabled: !!slug,
  });
}

export function useBrands() {
  return useQuery({
    queryKey: ["brands"],
    queryFn: () => apiFetch<{ data: Brand[] }>("/api/brands"),
  });
}

export function useBrandDetail(slug: string) {
  return useQuery({
    queryKey: ["brand", slug],
    queryFn: () => apiFetch<{ brand: Brand; products: Product[] }>(`/api/brands/${slug}`),
    enabled: !!slug,
  });
}

export function useSizeConversion(slug: string, referenceBrandSlug: string, size: string) {
  return useQuery({
    queryKey: ["size-conversion", slug, referenceBrandSlug, size],
    queryFn: () =>
      apiFetch<{ data: string | null; message?: string }>(
        `/api/products/${slug}/size-conversion?reference_brand=${referenceBrandSlug}&size=${encodeURIComponent(size)}`,
      ),
    enabled: !!slug && !!referenceBrandSlug && !!size,
  });
}

export function useTrendingProducts(limit = 8) {
  return useQuery({
    queryKey: ["products", "trending", limit],
    queryFn: () => apiFetch<{ data: Product[] }>(`/api/products/trending?limit=${limit}`),
  });
}

export function usePriceDropProducts(limit = 8) {
  return useQuery({
    queryKey: ["products", "price-drops", limit],
    queryFn: () =>
      apiFetch<{ data: (Product & { drop_percent: number })[] }>(`/api/products/price-drops?limit=${limit}`),
  });
}

export function useAutocomplete(query: string) {
  return useQuery({
    queryKey: ["autocomplete", query],
    queryFn: () =>
      apiFetch<{ products: { id: string; name: string; slug: string }[]; brands: { id: string; name: string; slug: string }[] }>(
        `/api/search/autocomplete?q=${encodeURIComponent(query)}`,
      ),
    enabled: query.length > 1,
  });
}
