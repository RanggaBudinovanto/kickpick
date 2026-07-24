"use client";

import { usePreferencesStore } from "@/stores/preferences";
import { useExchangeRate } from "@/hooks/use-exchange-rate";
import { formatPriceRangeIn } from "@/lib/currency";

export function PriceDisplay({ min, max }: { min: number; max: number }) {
  const currency = usePreferencesStore((s) => s.currency);
  const { data: rate } = useExchangeRate();

  return <>{formatPriceRangeIn(min, max, currency, rate?.rate)}</>;
}
