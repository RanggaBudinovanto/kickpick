import { test, expect } from "@playwright/test";

// Critical path per Section 20 PRD: Registrasi -> verifikasi email -> login -> logout.
// Email verification itself can't be driven end-to-end here (no test mailbox
// wired up — Resend isn't configured with a real key in this environment), so
// this covers registration through to the "check your email" confirmation, and
// separately exercises login/logout against the seeded demo account.
test.describe("Auth flow", () => {
  test("register shows the check-your-email confirmation", async ({ page }) => {
    const email = `e2e-${Date.now()}@example.com`;

    await page.goto("/id/registrasi");
    await page.getByLabel("Nama").fill("E2E Test User");
    await page.getByLabel("Email").fill(email);
    await page.getByLabel("Password", { exact: true }).fill("Password123");
    await page.getByLabel("Konfirmasi Password").fill("Password123");
    await page.getByRole("button", { name: "Daftar" }).click();

    await expect(page.getByText(/registrasi berhasil/i)).toBeVisible();
  });

  test("login with the seeded demo account, then logout", async ({ page }) => {
    await page.goto("/id/login");
    await page.getByLabel("Email").fill("demo@kickpick.id");
    await page.getByLabel("Password").fill("Password123");
    await page.getByRole("button", { name: "Masuk" }).click();

    // Successful login redirects to the homepage and the navbar shows a
    // profile link instead of "Masuk".
    await expect(page).toHaveURL(/\/id\/?$/);
    await page.goto("/id/profil");
    await expect(page.locator("#profile-email")).toHaveValue("demo@kickpick.id");

    await page.getByRole("button", { name: "Keluar" }).click();
    await page.goto("/id/wishlist");
    await expect(page).toHaveURL(/\/id\/login/);
  });

  test("rejects login with wrong password", async ({ page }) => {
    await page.goto("/id/login");
    await page.getByLabel("Email").fill("demo@kickpick.id");
    await page.getByLabel("Password").fill("wrong-password");
    await page.getByRole("button", { name: "Masuk" }).click();

    await expect(page.getByText(/email atau password salah/i)).toBeVisible();
    await expect(page).toHaveURL(/\/id\/login/);
  });
});
