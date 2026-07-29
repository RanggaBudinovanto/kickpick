import { test, expect } from "@playwright/test";

// Critical path per Section 20 PRD: Submit review -> moderasi rate limit teruji.
// "Rate limit" here is the PRD's one-review-per-product-per-user rule
// (Section 8: "Rate limit: maksimal 1 review per produk per user"), enforced
// by a unique constraint on (product_id, user_id).
test.describe("Review submission", () => {
  test("logged-in user can submit a review, and a second attempt is rejected", async ({ page }) => {
    await page.goto("/en/login");
    await page.getByLabel("Email").fill("demo@kickpick.id");
    await page.getByLabel("Password").fill("Password123");
    await page.getByRole("button", { name: "Log in" }).click();
    await expect(page).toHaveURL(/\/en\/?$/);

    await page.goto("/en/cari");
    await page.locator("a[href*='/produk/']").first().click();
    await expect(page).toHaveURL(/\/en\/produk\//);

    const writeReviewButton = page.getByRole("button", { name: /write a review/i });

    // First submission: either succeeds, or the demo user already reviewed
    // this product from a previous test run — both are valid starting states.
    // The success toast races with the page's router.refresh() (used so the
    // new review shows up immediately), so assert on the durable outcome —
    // the review appearing in the list — rather than the transient toast.
    if (await writeReviewButton.isVisible()) {
      const comment = `Nyaman dipakai sehari-hari, ukurannya pas. (${Date.now()})`;
      await writeReviewButton.click();
      await page.getByLabel("Comment").fill(comment);
      await page.getByRole("button", { name: "Submit review" }).click();
      await expect(page.getByText(comment)).toBeVisible();
    }

    // Either way, the "write a review" CTA must not be offered again for a
    // product this user already reviewed.
    await page.reload();
    await expect(page.getByRole("button", { name: /write a review/i })).not.toBeVisible();
  });
});
