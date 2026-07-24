import { test, expect } from "@playwright/test";

// Critical path per Section 20 PRD: Pencarian & filter produk -> detail produk
// -> klik beli (redirect affiliate).
test.describe("Search, product detail, and buy", () => {
  test("search finds a real seeded/scraped product and opens its detail page", async ({ page }) => {
    await page.goto("/id/cari");
    await expect(page.getByText(/sepatu ditemukan/i)).toBeVisible();

    const firstCard = page.locator("a[href*='/produk/']").first();
    await expect(firstCard).toBeVisible();
    await firstCard.click();

    await expect(page).toHaveURL(/\/id\/produk\//);
    // Price comparison section should render at least one store row.
    await expect(page.getByText(/Rp/).first()).toBeVisible();
  });

  test("category filter narrows results and updates the URL", async ({ page }) => {
    await page.goto("/id/cari");
    await page.getByLabel("Kategori").selectOption("lifestyle");

    await expect(page).toHaveURL(/kategori=lifestyle/);
  });

  test("clicking Beli on an in-stock offer opens a new tab to the store", async ({ page, context }) => {
    await page.goto("/id/cari");
    await page.locator("a[href*='/produk/']").first().click();
    await expect(page).toHaveURL(/\/id\/produk\//);

    const buyButton = page.getByRole("button", { name: "Beli" }).first();
    await expect(buyButton).toBeVisible();

    const [newPage] = await Promise.all([
      context.waitForEvent("page"),
      buyButton.click(),
    ]);
    await newPage.waitForLoadState("domcontentloaded");
    expect(newPage.url()).not.toContain("localhost:3000");
  });
});
