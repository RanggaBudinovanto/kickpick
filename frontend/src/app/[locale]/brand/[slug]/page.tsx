import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { ProductCard } from "@/components/product/product-card";
import { getBrandDetail } from "@/lib/server-api";

interface PageProps {
  params: Promise<{ slug: string; locale: string }>;
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { slug } = await params;
  const data = await getBrandDetail(slug);

  if (!data) {
    return { title: "Brand not found | KickPick" };
  }

  return {
    title: `${data.brand.name} - Compare Prices | KickPick`,
    description: `Compare prices for ${data.brand.name} products from various stores on KickPick.`,
    alternates: {
      languages: {
        id: `/id/brand/${slug}`,
        en: `/en/brand/${slug}`,
        "x-default": `/en/brand/${slug}`,
      },
    },
  };
}

export default async function BrandDetailPage({ params }: PageProps) {
  const { slug } = await params;
  const data = await getBrandDetail(slug);

  if (!data) {
    notFound();
  }

  return (
    <div className="mx-auto max-w-[1400px] px-4 py-10 md:px-6">
      <h1 className="mb-6 font-display text-3xl font-bold tracking-[-0.01em]">{data.brand.name}</h1>

      {data.products.length === 0 ? (
        <p className="text-sm text-muted">No products found for this brand yet.</p>
      ) : (
        <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
          {data.products.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      )}
    </div>
  );
}
