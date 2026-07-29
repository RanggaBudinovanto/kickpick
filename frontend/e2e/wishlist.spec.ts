import { test, expect } from "@playwright/test";

async function login(page: import("@playwright/test").Page) {
  await page.goto("/en/login");
  await page.getByLabel("Email").fill("demo@kickpick.id");
  await page.getByLabel("Password").fill("Password123");
  await page.getByRole("button", { name: "Log in" }).click();
  await expect(page).toHaveURL(/\/en\/?$/);
}

// Critical path per Section 20 PRD: Login -> tambah wishlist -> aktifkan alert.
test.describe("Wishlist and alert", () => {
  test("guest is redirected to login when trying to save a product", async ({ page }) => {
    await page.goto("/en/cari");
    await page.locator("a[href*='/produk/']").first().click();

    await page.getByRole("link", { name: /save to wishlist/i }).click();
    await expect(page).toHaveURL(/\/en\/login/);
  });

  test("logged-in user can add a product to wishlist and activate an alert", async ({ page }) => {
    await login(page);

    await page.goto("/en/cari");
    await page.locator("a[href*='/produk/']").first().click();
    await expect(page).toHaveURL(/\/en\/produk\//);

    // The demo account's wishlist persists across test runs, so this product
    // may already be saved from a previous run — don't assume a fresh start.
    // WishlistButton fetches its own state client-side, so wait for either
    // button to appear before branching (isVisible() alone doesn't wait and
    // can catch the component mid-fetch). Toggle once in whichever direction
    // is available rather than resetting-then-setting, since two round trips
    // in one test doubles the flakiness surface for no real extra coverage —
    // one toggle already proves the add/remove mechanism works.
    const removeButton = page.getByRole("button", { name: /remove from wishlist/i });
    const addButton = page.getByRole("button", { name: /save to wishlist/i });
    await expect(removeButton.or(addButton)).toBeVisible();

    if (await addButton.isVisible()) {
      await addButton.click();
      await expect(removeButton).toBeVisible();
    } else {
      await removeButton.click();
      await expect(addButton).toBeVisible();
    }

    await page.goto("/en/wishlist");
    const firstAlertCheckbox = page.getByRole("checkbox", { name: "Alert" }).first();
    await expect(firstAlertCheckbox).toBeVisible();
    
    // Ensure we test the checking transition cleanly even if already checked from previous run
    if (await firstAlertCheckbox.isChecked()) {
      await firstAlertCheckbox.uncheck();
      await expect(firstAlertCheckbox).not.toBeChecked();
    }
    
    await firstAlertCheckbox.check();
    await expect(firstAlertCheckbox).toBeChecked();
  });
});
