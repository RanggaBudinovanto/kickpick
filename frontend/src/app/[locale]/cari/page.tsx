import type { Metadata } from "next";
import { IconSearchOff } from "@tabler/icons-react";
import { ProductCard } from "@/components/product/product-card";
import { SearchFilters } from "@/components/search/search-filters";
import { ResetFilterLink } from "@/components/search/reset-filter-link";
import { getBrands, getProducts } from "@/lib/server-api";

interface PageProps {
  searchParams: Promise<{
    q?: string;
    kategori?: string;
    brand_ids?: string;
    filter?: string;
  }>;
}

export async function generateMetadata({ searchParams }: PageProps): Promise<Metadata> {
  const params = await searchParams;
  const title = params.q
    ? `Cari "${params.q}" | KickPick`
    : "Cari Sepatu - Bandingkan Harga Semua Brand | KickPick";

  return {
    title,
    description: "Cari dan bandingkan harga sepatu dari brand lokal dan internasional di KickPick.",
    alternates: {
      languages: { id: "/id/cari", en: "/en/cari", "x-default": "/id/cari" },
    },
  };
}

export default async function SearchPage({ searchParams }: PageProps) {
  const params = await searchParams;

  const [productsRes, brandsRes] = await Promise.all([
    getProducts({
      q: params.q,
      kategori: params.kategori,
      brand_ids: params.brand_ids,
      filter: params.filter,
      limit: 24,
    }),
    getBrands(),
  ]);

  const products = productsRes.data;

  return (
    <div className="mx-auto flex max-w-[1400px] flex-col gap-6 px-4 py-8 md:flex-row md:px-6">
      <aside className="w-full shrink-0 md:w-64">
        <SearchFilters brands={brandsRes.data} />
      </aside>

      <div className="flex-1">
        <p className="mb-4 text-sm text-muted">{products.length} sepatu ditemukan</p>

        {products.length === 0 && (
          <div className="flex flex-col items-center gap-3 py-16 text-center">
            <IconSearchOff size={32} />
            <p className="text-sm font-medium">Tidak ada sepatu yang cocok dengan filter ini</p>
            <ResetFilterLink />
          </div>
        )}

        <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
          {products.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      </div>
    </div>
  );
}
