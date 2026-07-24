"use client";

import { useTranslations } from "next-intl";
import { IconCheck, IconX } from "@tabler/icons-react";
import { Button } from "@/components/ui/button";
import { formatPriceIn } from "@/lib/currency";
import { useGoToOffer } from "@/hooks/use-redirect";
import { useExchangeRate } from "@/hooks/use-exchange-rate";
import { usePreferencesStore } from "@/stores/preferences";
import type { Offer } from "@/types/api";

export function PriceComparison({ offers }: { offers: Offer[] }) {
  const t = useTranslations("Product");
  const goToOffer = useGoToOffer();
  const currency = usePreferencesStore((s) => s.currency);
  const { data: rate } = useExchangeRate();

  if (offers.length === 0) {
    return <p className="text-sm text-muted">{t("noOffers")}</p>;
  }

  const cheapestPrice = Math.min(...offers.filter((o) => o.in_stock).map((o) => o.price));

  return (
    <div className="overflow-hidden rounded-[var(--radius-card)] border border-border">
      {offers.map((offer) => {
        const isCheapest = offer.in_stock && offer.price === cheapestPrice;
        return (
          <div
            key={offer.id}
            className={`flex items-center justify-between gap-4 p-4 ${
              isCheapest ? "border-2 border-off-black dark:border-pure-white" : "border-b border-border last:border-b-0"
            }`}
          >
            <div>
              <p className="text-sm font-medium">{offer.store_name}</p>
              <p className="flex items-center gap-1 text-xs text-muted">
                {offer.in_stock ? (
                  <>
                    <IconCheck size={12} /> {t("inStock")}
                  </>
                ) : (
                  <>
                    <IconX size={12} /> {t("outOfStock")}
                  </>
                )}
                {isCheapest && ` · ${t("cheapest")}`}
              </p>
            </div>
            <div className="flex items-center gap-3">
              <span className={`font-mono tabular-nums ${isCheapest ? "text-lg font-bold" : "text-sm"}`}>
                {formatPriceIn(offer.price, currency, rate?.rate)}
              </span>
              <Button
                size="sm"
                disabled={!offer.in_stock || goToOffer.isPending}
                onClick={() => goToOffer.mutate(offer.id)}
              >
                {t("buy")}
              </Button>
            </div>
          </div>
        );
      })}
    </div>
  );
}
