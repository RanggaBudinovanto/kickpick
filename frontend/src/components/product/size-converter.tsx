"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useBrands, useSizeConversion } from "@/hooks/use-products";

export function SizeConverter({ productSlug }: { productSlug: string }) {
  const t = useTranslations("Product");
  const { data: brandsData } = useBrands();
  const [referenceBrandSlug, setReferenceBrandSlug] = useState("");
  const [size, setSize] = useState("");

  const { data, isFetching } = useSizeConversion(productSlug, referenceBrandSlug, size);

  return (
    <Card className="p-5">
      <h2 className="mb-4 font-display text-xl font-semibold">{t("sizeConverter")}</h2>

      <div className="grid grid-cols-2 gap-3">
        <div>
          <label htmlFor="size-reference-brand" className="mb-1.5 block text-sm font-medium">
            {t("referenceBrand")}
          </label>
          <select
            id="size-reference-brand"
            value={referenceBrandSlug}
            onChange={(e) => setReferenceBrandSlug(e.target.value)}
            className="h-11 w-full rounded-[var(--radius-control)] border border-zinc-500/40 bg-background px-3 text-sm"
          >
            <option value="">-</option>
            {brandsData?.data.map((b) => (
              <option key={b.id} value={b.slug}>
                {b.name}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label htmlFor="size-your-size" className="mb-1.5 block text-sm font-medium">
            {t("yourSize")}
          </label>
          <Input id="size-your-size" value={size} onChange={(e) => setSize(e.target.value)} placeholder="42" />
        </div>
      </div>

      {referenceBrandSlug && size && (
        <div className="mt-4 rounded-[var(--radius-control)] bg-surface p-3">
          {isFetching ? (
            <p className="text-sm text-muted">...</p>
          ) : data?.data ? (
            <p className="font-mono text-lg font-bold tabular-nums">
              {t("convertedSize")}: {data.data}
            </p>
          ) : (
            <p className="text-sm text-muted">{data?.message}</p>
          )}
        </div>
      )}

      <p className="mt-3 text-xs text-muted">{t("sizeConverterDisclaimer")}</p>
    </Card>
  );
}
