import { describe, expect, it } from "vitest";
import { formatPrice, formatPriceRange } from "./currency";

// formatPrice must never show decimals (no sen/cents for IDR) and must use
// thousands separators — asserted structurally rather than against one exact
// string, since ICU inserts a non-breaking space between "Rp" and the amount
// that varies subtly across Node/browser Intl implementations.
describe("formatPrice", () => {
  it("formats a whole number with thousands separators and no decimals", () => {
    const result = formatPrice(1698000);
    expect(result).toContain("1.698.000");
    expect(result).not.toContain(",");
    expect(result.startsWith("Rp")).toBe(true);
  });

  it("formats zero", () => {
    expect(formatPrice(0)).toContain("0");
  });
});

describe("formatPriceRange", () => {
  it("shows a single price when min equals max", () => {
    const result = formatPriceRange(500000, 500000);
    expect(result).toContain("500.000");
    expect(result).not.toContain("-");
  });

  it("shows a single price when max is 0 (no offers yet)", () => {
    const result = formatPriceRange(500000, 0);
    expect(result).toContain("500.000");
    expect(result).not.toContain("-");
  });

  it("shows a range when min and max differ", () => {
    const result = formatPriceRange(400000, 600000);
    expect(result).toContain("400.000");
    expect(result).toContain("600.000");
    expect(result).toContain("-");
  });
});
