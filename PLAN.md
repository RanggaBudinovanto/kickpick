# Rencana Pembangunan KickPick

> Dokumen ini memecah pembangunan KickPick (lihat `KickPick_PRD.md` dan `DESIGN.md`) menjadi 4 fase berurutan. Tiap fase punya definition of done sendiri sehingga bisa didemo/dites sebelum lanjut ke fase berikutnya — bukan big-bang di akhir.

---

## Plan 1 — Fondasi & Design System

**Tujuan:** kerangka teknis siap, belum ada data nyata.

- Setup monorepo (Next.js 15 App Router + TypeScript + Tailwind, Go + Fiber backend terpisah)
- Skema database PostgreSQL penuh (Section 11 PRD) + migrasi via sqlc
- Design system diimplementasikan sebagai kode: font Barlow/Barlow Condensed, palet monokrom, komponen dasar shadcn dikustom (Button, Card, Badge, Input) sesuai DESIGN.md
- Layout global: Navbar (dengan switcher bahasa/currency), Footer
- i18n routing (`next-intl`, `/id`, `/en`) dan dark/light mode toggle

**Selesai kalau:** halaman kosong dengan navbar/footer tampil benar di kedua bahasa, kedua tema, dan mobile/desktop, tanpa data asli.

---

## Plan 2 — Core Product Experience (dengan data seed/dummy)

**Tujuan:** semua halaman utama berfungsi penuh secara UI/UX, sebelum scraping nyata ada.

- Auth penuh (register, login, verifikasi email, JWT+refresh sesuai Section 13)
- Homepage (hero, brand strip, kategori, best seller, trending/rare, price drop)
- Search & Listing `/cari` dengan semua filter
- Detail Produk `/produk/[slug]`: price comparison table, grafik histori harga, size converter, review komunitas
- Brand Directory, Wishlist, Notifikasi, Profil
- Data diisi lewat seed script (bukan scraping asli) supaya UI bisa langsung dites end-to-end

**Selesai kalau:** semua critical path di Section 20 PRD bisa dijalankan manual dengan data seed.

---

## Plan 3 — Scraping Engine & Data Pipeline (bagian paling berisiko)

**Tujuan:** data asli mengalir masuk, menggantikan seed data.

- Adapter scraping modular per-brand (Colly untuk situs statis, chromedp untuk situs dinamis) — mulai dari 3-5 brand prioritas dulu (perlu keputusan brand mana, lihat Section 2 PRD)
- Integrasi affiliate network (Shopee, Tokopedia, Involve Asia) untuk brand yang tidak scraping langsung
- Job scheduler (Asynq + Redis): update harga harian, cek restock berkala, simpan `price_history`
- Deduplikasi produk antar sumber (situs resmi vs marketplace)
- Fake-discount detector, currency conversion cron

**Selesai kalau:** minimal 5 brand punya data harga live yang ter-update otomatis dan histori harga mulai terkumpul.

---

## Plan 4 — Growth Features, Keamanan, & Launch Readiness

**Tujuan:** dari "berfungsi" ke "siap launch" sesuai klaim full-build di PRD.

- Notifikasi (in-app + email Resend, restock/price-drop alert)
- Review komunitas + moderasi (rate limit + report)
- Security hardening penuh (Section 14 checklist: CORS, CSRF, rate limiting, security headers, ownership check)
- SEO (sitemap, hreflang, structured data) + Lighthouse ≥90
- Testing (unit, integration, E2E Playwright untuk 4 critical path)
- Deployment (Vercel + Railway/Fly.io + CI/CD GitHub Actions) + backup terjadwal

**Selesai kalau:** semua checklist Section 22 & 23 PRD lolos.

---

## Catatan Penting / Blocker

Dua item di PRD masih **belum jelas** dan akan memblokir Plan 3 kalau tidak diputuskan dulu:

1. Daftar brand prioritas awal untuk scraping (lihat Section 2 PRD)
2. Status legal ToS scraping per brand

Disarankan kedua hal ini diputuskan sebelum masuk Plan 3.
