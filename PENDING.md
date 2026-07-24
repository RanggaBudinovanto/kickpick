# Yang Masih Ditunda / Belum Rampung

> Dicatat per fase. Ini bukan bug — ini adalah simplifikasi sadar supaya tiap fase bisa selesai dan dites end-to-end, sebelum masuk ke scraping data asli (Plan 3) dan hardening penuh (Plan 4). Item yang sudah diselesaikan dicoret/dipindah ke bagian "Selesai" di tiap kategori supaya riwayatnya tetap jelas.

---

## Dari Plan 2 (Core Product Experience) — update setelah sesi perbaikan

### Selesai di sesi ini
- ✅ **Trending section** — sekarang berdasarkan agregasi view produk 7 hari terakhir (tabel `product_views`, endpoint `/api/products/trending`), bukan lagi placeholder.
- ✅ **Price Drop Deals section** — query `ListPriceDropProducts` membandingkan harga termurah saat ini vs rata-rata 30 hari (`price_history`), endpoint `/api/products/price-drops`.
- ✅ **Filter brand di `/cari` sekarang multi-select** (checkbox, param `brand_ids` sebagai `uuid[]` pakai `= ANY()`).
- ✅ **Email verifikasi & reset password sekarang benar-benar terkirim** lewat Resend API (`internal/email/`). Kalau `RESEND_API_KEY` kosong (mis. dev lokal), sistem log ke console alih-alih gagal — sudah dites, tidak mem-block request registrasi/forgot-password.
- ✅ **CSRF protection** — custom header `X-Requested-With: kickpick` wajib di semua endpoint yang mengubah state (refresh, logout, review, wishlist, notifikasi, profil), sudah dites 403 tanpa header / 200 dengan header.
- ✅ **Audit log diperluas** — sekarang mencatat `login`, `logout`, `password_reset`, `account_deleted`, selain `affiliate_click` yang sudah ada sebelumnya.
- ✅ **`/produk/[slug]` sekarang Server Component** dengan `generateMetadata` dinamis (title per produk, deskripsi, Open Graph, hreflang `id`/`en`), structured data `Product` (JSON-LD). HTML awal dari server sudah berisi data produk asli, bukan skeleton kosong — dicek langsung lewat curl.
- ✅ **`app/sitemap.ts` dan `app/robots.ts`** — sitemap mencakup semua produk & brand di kedua locale, robots blokir `/wishlist`, `/notifikasi`, `/profil`, `/api`.
- ✅ **Structured data `Organization`** di homepage.
- ✅ **Halaman `/tentang` (FAQ), `/privasi`, `/disclosure`** dibangun (sebelumnya cuma link di footer yang menuju halaman kosong/404).

### Masih tertunda
- **`/cari` dan `/brand/[slug]` masih CSR**, cuma `/produk/[slug]` dan homepage yang sudah dikonversi ke Server Component. Idealnya semua halaman public di-SSR untuk SEO penuh sesuai Section 21 PRD — ini bisa dilanjutkan dengan pola yang sama (server fetch + `generateMetadata`) kapan saja, tidak butuh perubahan arsitektur baru.
- **Filter warna, ukuran, dan status stok di `/cari` belum ada.** Sidebar filter baru mendukung kategori, brand (multi-select), dan rare/limited. Data warna/ukuran granular belum konsisten di `products.attributes` (jsonb) karena belum ada data scraping asli yang bentuknya jelas — masuk akal digarap bareng Plan 3.
- **Rating produk selalu 0 kalau belum ada review** — belum ada fallback UI khusus (cuma disembunyikan di card), bukan bug, tapi belum polish.
- **`preferred_currency` user tersimpan di DB tapi belum dipakai untuk konversi tampilan.** Semua harga masih tampil IDR saja. `EXCHANGE_RATES` tabel sudah ada tapi belum ada cron yang mengisi otomatis (Plan 3/4 — cron kurs harian).
- **ISR/caching untuk `/produk/[slug]` sengaja pakai `cache: "no-store"`**, bukan `revalidate: 3600` sesuai Section 19 PRD, supaya review yang baru disubmit langsung kelihatan tanpa nunggu cache expire. Perlu direvisit di Plan 4: pindah ke ISR + `revalidateTag` yang dipanggil saat ada perubahan harga (scraping) atau review baru, biar dapat manfaat cache tanpa data basi.

### Push & Notifikasi
- **Push notification (OneSignal) belum diintegrasikan** — sesuai PRD ini memang kanal opsional untuk rilis pertama, bukan blocker.
- **Belum ada job scheduler yang membuat notifikasi restock/price-drop otomatis.** Template email `E04`/`E05` (`RestockAlertEmail`, `PriceDropAlertEmail`) sudah ditulis di `internal/email/templates.go` supaya tinggal dipanggil begitu ada job scheduler di Plan 3 yang mendeteksi perubahan harga/stok asli.

### Auth & Keamanan
- **Rate limiting masih in-memory (default Fiber limiter), bukan Redis-backed.** Section 19 PRD minta rate limit counter di Redis supaya konsisten across multiple server instance. Cukup untuk single-instance dev, tidak akan scale ke multi-instance production tanpa Redis. (Masih berlaku setelah Plan 4 — lihat bawah.)
- **Redirect endpoint `/api/redirect/:offer_id` sengaja return JSON (bukan HTTP redirect 302 langsung)** karena method-nya POST sesuai PRD — keputusan desain sadar (dicatat di komentar `internal/handler/redirect_handler.go`), bukan gap, supaya tidak dikira bug oleh siapa pun yang expect redirect HTTP native.

---

## Dari Plan 3 (Scraping Engine & Data Pipeline)

### Temuan lapangan penting: dari 6 brand kandidat, cuma 1 yang benar-benar scrapable sekarang
Sebelum coding, tiap brand dicek robots.txt-nya dan struktur situsnya secara langsung (bukan asumsi). Hasilnya jauh dari yang diharapkan Section 24 PRD:

| Brand | Status | Alasan |
|---|---|---|
| **Compass** | ✅ Scraping nyata jalan | `sepatucompass.com` (bukan `compass.co.id` — domain itu ternyata bisnis furniture tidak terkait, sempat salah tebak di awal). robots.txt izinkan crawl, situs Shopify Hydrogen yang SSR seluruh katalog sebagai JSON di `window.__remixContext`. 107 produk asli berhasil di-scrape dan disimpan. |
| **Ventela** | ❌ Tidak bisa sekarang | Domain resmi (`ventela.id`) DAN domain toko (`ventelashoes.com`) dua-duanya dalam mode "Coming Soon" — tidak ada katalog live sama sekali saat dicek. |
| **Aerostreet** | ❌ Tidak bisa langsung | Domain sendiri (`aerostreet.co.id`) tidak punya katalog/checkout — cuma halaman profil brand yang mengarahkan ke marketplace (Shopee, Tokopedia, Lazada, Blibli). Butuh integrasi affiliate API, bukan scraping situs sendiri. |
| **Nike** | ❌ Diblokir | `nike.com/id` robots.txt eksplisit `Disallow: */p/` — pola URL itu mencakup semua halaman produk. Dihormati sesuai Section 14 PRD, tidak di-scrape. |
| **Adidas** | ❌ Diblokir | `adidas.co.id` pakai proteksi bot Akamai yang me-reject request dasar (403) bahkan untuk baca robots.txt. |
| **Vans** | ❌ Belum terverifikasi | Domain resmi yang benar belum ditemukan/dicek (tebakan awal tidak resolve). |

**Kesimpulan:** strategi "scraping situs brand langsung" cuma realistis untuk sebagian kecil brand lokal yang memang punya toko online sendiri dan mengizinkan crawling. Brand besar internasional (Nike, Adidas) dan sebagian brand lokal (Aerostreet) secara struktural butuh jalur **affiliate marketplace API** (Shopee/Tokopedia/Involve Asia), bukan scraping — persis seperti yang diprediksi Section 24 PRD, hanya saja polanya berbeda dari asumsi awal (Compass ternyata scrapable, Aerostreet ternyata tidak).

### Selesai
- ✅ **Scraper interface modular** (`internal/scraper/types.go`) — `Adapter` interface generik, tidak terikat ke Colly/chromedp tertentu, jadi adapter brand berikutnya bisa pilih tool sesuai kebutuhan situsnya.
- ✅ **Adapter Compass nyata** (`internal/scraper/compass/`) — pakai Colly dengan rate limit 1 request/2 detik, User-Agent jujur (`KickPickBot/1.0`, bukan menyamar jadi browser), parser JSON rekursif yang tahan terhadap perubahan struktur halaman (tidak bergantung pada path JSON yang persis).
- ✅ **Pipeline persistensi** (`internal/scraper/pipeline.go`) — upsert produk by slug, replace gambar, upsert offer per toko, catat histori harga harian. Dites idempoten: dijalankan 2x berturut-turut menghasilkan jumlah produk/offer/histori harga yang identik (tidak ada duplikat).
- ✅ **Job scheduler** (`cmd/worker`, `internal/scheduler/`) — cron in-process (robfig/cron, bukan Asynq+Redis karena Redis tidak terpasang di environment ini) yang jalan sekali saat start lalu terjadwal harian jam 03:00.
- ✅ **CLI manual** (`cmd/scraper`) — untuk trigger scrape sekali tanpa nunggu jadwal, berguna untuk testing/debugging.
- ✅ **Bug ditemukan & diperbaiki**: `Float64ToNumeric` di `internal/db/pgtypes.go` diam-diam gagal (error di-swallow) sehingga semua harga masuk sebagai NULL — kena constraint NOT NULL sebelum sempat masuk data salah. Diperbaiki dengan scan lewat string, bukan raw float64.

### Masih tertunda
- **5 dari 6 brand kandidat belum punya adapter** — lihat tabel temuan di atas untuk alasan spesifik per brand. Aerostreet, Nike, Adidas butuh integrasi affiliate API (Shopee/Tokopedia/Involve Asia) yang belum ada kredensialnya. Ventela perlu dicek ulang berkala sampai tokonya keluar dari mode "Coming Soon". Vans perlu riset domain yang benar.
- **Belum ada integrasi affiliate network sungguhan** (Shopee Affiliate, Tokopedia, Involve Asia) — masih asumsi PRD Section 2, belum ada API key/partnership nyata.
- **Deduplikasi lintas-sumber belum dibangun** — PRD Section 24 mencatat brand lokal bisa punya 2 sumber data (situs resmi + marketplace) untuk produk yang sama, butuh logic dedup. Belum relevan sampai ada adapter kedua untuk brand yang sama.
- **Fake-discount detector sudah ada versi dasarnya** (endpoint `/api/products/price-drops` dari Plan 2, bandingkan harga vs rata-rata 30 hari) tapi belum divalidasi dengan data harga yang benar-benar bergerak dari waktu ke waktu (baru 1 hari data scraping asli sejauh ini).
- **Currency conversion cron belum jalan** — `EXCHANGE_RATES` tabel ada, tapi belum ada job yang mengisi kurs harian dari API pihak ketiga.
- **Scheduler cron in-process, bukan Redis-backed Asynq** seperti tech stack PRD Section 10 — konsekuensinya: kalau nanti `cmd/worker` dijalankan lebih dari satu instance, tiap instance akan trigger scrape sendiri-sendiri (duplikasi kerja, walau tidak duplikasi data karena pipeline idempoten). Aman untuk single-instance deployment, perlu Redis+Asynq kalau sudah butuh scale.
- **robots.txt dicek manual sekali per brand saat development**, belum ada pengecekan otomatis berkala (Section 14 PRD minta brand yang di-scraping "ditinjau ToS-nya secara berkala").

---

## Dari Plan 4 (Growth Features, Keamanan, & Launch Readiness)

### Selesai
- ✅ **Security headers eksplisit** — CSP (`default-src 'none'`, cocok untuk API JSON-only), `X-Frame-Options: DENY`, HSTS (muncul otomatis begitu di-serve lewat HTTPS di production), `Referrer-Policy`. Sebelumnya cuma pakai default helmet yang tidak set CSP.
- ✅ **Ownership check diaudit** — dikonfirmasi semua endpoint yang akses data user (wishlist, notifikasi, profil) sudah benar filter `WHERE user_id = $token_user_id`, bukan cuma cek auth tapi tidak cek kepemilikan.
- ✅ **Automated testing lengkap ditambahkan** (Section 20 PRD):
  - Backend: unit test (`internal/auth`, `internal/db`, `internal/middleware`, `internal/scraper/compass`) + integration test nyata (`internal/router`) yang jalan lewat HTTP stack sungguhan ke database test terpisah (`kickpick_test`), bukan mock.
  - Frontend: component test (Vitest + React Testing Library) untuk `FormField` dan halaman login (validasi kosong, email salah format, submit sukses).
  - E2E (Playwright): 9 test mencakup keempat critical path Section 20 — register→login→logout, cari→detail→beli (redirect affiliate ke tab baru), login→wishlist→alert, submit review→duplikat ditolak. Semua dites terhadap app & database nyata, bukan mock.
- ✅ **3 bug nyata ditemukan & diperbaiki lewat testing** (bukan cuma "test hijau", testing ini benar-benar menemukan masalah):
  1. Label form tidak terhubung ke input (`htmlFor`/`id` hilang) di `FormField`, filter `/cari`, size converter, review form, halaman profil — masalah aksesibilitas nyata (screen reader tidak bisa umumkan label yang benar), ditemukan karena test butuh `getByLabelText`.
  2. Native browser validation (`type="email"`) memblokir event submit sebelum validasi Zod sempat jalan, jadi pesan error custom (sesuai DESIGN.md) tidak pernah muncul — pengguna malah lihat popup native browser yang tidak konsisten dengan desain. Diperbaiki dengan `noValidate` di ketiga form auth.
  3. Checkbox alert wishlist tidak ada optimistic update, jadi UI sempat "snap back" ke status lama selama menunggu round-trip API — diperbaiki dengan optimistic update + rollback di `useSetWishlistAlert`.
- ✅ **`govulncheck` dijalankan nyata, menemukan 10 kerentanan reachable** dari kode scraper (`golang.org/x/net`, `github.com/antchfx/xpath` — parsing HTML/XPath yang dipakai Colly). Diperbaiki dengan upgrade dependency (`x/net` v0.47→v0.55, `xpath` v1.3.5→v1.3.6). Sisa 3 temuan ada di Go standard library sendiri (butuh upgrade toolchain ke go1.26.4+, dicatat di bawah, bukan diabaikan).
- ✅ **GitHub Actions CI** (`.github/workflows/ci.yml`) — job backend (vet, migrate ke Postgres service container, test, build, govulncheck) dan frontend (lint, test, build, npm audit), jalan di tiap push/PR ke `main`/`staging`.
- ✅ **Dockerfile backend** (`backend/Dockerfile` untuk API, `backend/Dockerfile.worker` untuk scheduler terpisah) + `docker-compose.yml` referensi untuk orkestrasi lokal.
- ✅ **Bug seed data**: gambar seed pakai `placehold.co` tanpa ekstensi (default SVG), Next.js Image menolak render karena `dangerouslyAllowSVG` sengaja dimatikan (praktik aman, SVG bisa bawa script). Diperbaiki dengan minta format PNG eksplisit dari placehold.co.
- ✅ **Bug konfigurasi Vitest**: `npm test` awalnya ikut menjalankan file Playwright di `e2e/` (glob default Vitest cocok dengan `*.spec.ts` di mana pun), bikin 4 test file gagal dengan error tidak relevan. Diperbaiki dengan exclude eksplisit folder `e2e/` di `vitest.config.ts` — ini juga akan mem-block CI kalau tidak ditemukan sebelum push pertama.

### Masih tertunda
- **Dockerfile belum divalidasi dengan `docker build` sungguhan** — Docker tidak terpasang di environment development ini. Ditulis mengikuti pola multi-stage standar Go, kemungkinan besar benar, tapi belum ada bukti build sukses.
- **Toolchain Go masih 1.26.3, ada 3 kerentanan di standard library yang baru fix di go1.26.4** (`net/textproto`, `crypto/x509`) — govulncheck traces-nya agak tidak langsung (lewat `rand.Read`/`x509.HostnameError.Error`), risiko rendah untuk pola pakai kita saat ini, tapi tetap perlu upgrade toolchain saat memungkinkan.
- **`npm audit` melaporkan 3 high severity vuln** (`postcss`, `sharp`) — keduanya dependency transitif **di dalam Next.js sendiri** (dipakai untuk tooling build/image processing internal Next, bukan dipanggil langsung oleh kode kita). Fix yang disarankan npm (`--force`) akan downgrade Next.js ke v9 (breaking change besar, bukan solusi nyata) — ditandai `|| true` di CI supaya tidak false-block deploy, tapi perlu dipantau sampai Next.js sendiri bump dependency-nya.
- **Rate limiting tetap in-memory**, belum Redis-backed (item ini juga sudah dicatat sejak Plan 2, masih relevan — Redis belum pernah terpasang di sesi manapun sejauh ini karena kendala disk berulang).
- **CI belum pernah benar-benar dijalankan di GitHub** — proyek ini belum jadi git repository (belum ada `git init`/remote), jadi `.github/workflows/ci.yml` baru diverifikasi secara manual (menjalankan langkah yang sama satu-satu secara lokal), belum lewat GitHub Actions runner sungguhan.
- **Backup database terjadwal belum diatur** — ini murni tanggung jawab platform hosting (Railway/Neon biasanya punya opsi backup otomatis), belum ada konfigurasi eksplisit karena belum deploy ke platform manapun.
- **Lighthouse score belum diukur** — Section 23 PRD minta skor ≥90, belum pernah dijalankan audit Lighthouse sungguhan terhadap build production.

---

## Dari Plan 1 (Fondasi)

- **sqlc codegen butuh `CGO_ENABLED=0`** di environment Windows ini karena toolchain MinGW gcc lokal tidak bisa compile parser cgo `pg_query_go`. Sudah didokumentasikan di README, bukan blocker, tapi perlu diketahui siapa pun yang lanjut development di mesin ini.
- **Next.js yang terpasang versi 16, bukan 15** seperti disebut di PRD Section 10 — `create-next-app@latest` menarik versi terbaru yang tersedia saat scaffolding. Secara API App Router kompatibel, tidak ada masalah fungsional, tapi versi persisnya beda dari yang tertulis di dokumen.

---

## Blocker Bisnis

- Daftar brand final untuk cakupan scraping asli — 6 brand awal (Compass, Ventela, Aerostreet, Nike, Adidas, Vans) masih kandidat prioritas dari sisi bisnis, tapi kelayakan teknisnya sudah diverifikasi nyata di Plan 3 (lihat tabel temuan di atas) — cuma Compass yang scrapable langsung hari ini.
- Legal/ToS review per brand sebelum scraper di-deploy ke production — untuk Compass sudah dicek robots.txt-nya (mengizinkan), tapi review legal menyeluruh (di luar robots.txt teknis) tetap seperti yang disepakati sebelumnya, ditunda ke kamu.
- Affiliate network mana yang benar-benar disetujui (Shopee/Tokopedia/Involve Asia masih asumsi PRD Section 2) — ini sekarang jadi blocker konkret, bukan cuma asumsi dokumen, karena 3 dari 6 brand kandidat (Aerostreet, Nike, Adidas) terbukti butuh jalur affiliate marketplace, bukan scraping situs sendiri.
