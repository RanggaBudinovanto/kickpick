import type { Metadata } from "next";
import { getBrands } from "@/lib/server-api";
import { BrandDirectoryClient } from "@/components/brand/brand-directory-client";

export const metadata: Metadata = {
  title: "Brand Sepatu - Direktori Lengkap | KickPick",
  description: "Lihat direktori lengkap brand sepatu lokal dan internasional di KickPick.",
};

export default async function BrandDirectoryPage() {
  const brandsRes = await getBrands();

  return (
    <div className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
      <h1 className="mb-6 font-display text-3xl font-bold tracking-[-0.01em]">Brand</h1>
      <BrandDirectoryClient brands={brandsRes.data} />
    </div>
  );
}
