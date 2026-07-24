const idrFormatter = new Intl.NumberFormat("id-ID", {
  style: "currency",
  currency: "IDR",
  maximumFractionDigits: 0,
});

const usdFormatter = new Intl.NumberFormat("en-US", {
  style: "currency",
  currency: "USD",
  maximumFractionDigits: 2,
});

export function formatPrice(amount: number) {
  return idrFormatter.format(amount);
}

export function formatPriceRange(min: number, max: number) {
  if (min === max || max === 0) {
    return formatPrice(min);
  }
  return `${formatPrice(min)} - ${formatPrice(max)}`;
}

// amountIDR is always the source of truth (harga dasar disimpan IDR di
// database, Section 10/19 PRD); usdRate converts IDR -> USD for display only.
export function formatPriceIn(amountIDR: number, currency: "IDR" | "USD", usdRate?: number) {
  if (currency === "USD" && usdRate) {
    return usdFormatter.format(amountIDR * usdRate);
  }
  return formatPrice(amountIDR);
}

export function formatPriceRangeIn(
  min: number,
  max: number,
  currency: "IDR" | "USD",
  usdRate?: number,
) {
  if (min === max || max === 0) {
    return formatPriceIn(min, currency, usdRate);
  }
  return `${formatPriceIn(min, currency, usdRate)} - ${formatPriceIn(max, currency, usdRate)}`;
}
