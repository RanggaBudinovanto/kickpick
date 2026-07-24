# PRD — KickPick

> Dokumen ini dirancang untuk vibe coding. AI coder (Claude Code atau sejenisnya) harus bisa membaca PRD ini dan langsung membangun tanpa perlu bertanya lebih lanjut.

---

## 1. Project Overview

**Nama proyek:** KickPick
**Status:** Siap Dibangun
**Tanggal:** 22 Juli 2026

### Apa yang dibangun
KickPick adalah platform pencarian dan perbandingan harga sepatu lintas brand (lokal maupun internasional) dalam satu tempat. Pengguna memasukkan kriteria (budget, warna, jenis, ukuran), lalu KickPick menampilkan produk yang sesuai dari berbagai toko/marketplace, lengkap dengan histori harga, konversi ukuran antar brand, dan alert restock. Pengguna menyelesaikan pembelian di toko resmi/marketplace tujuan (KickPick tidak menjual produk sendiri).

### Siapa penggunanya
1. **Pembeli umum** — mencari sepatu sesuai kebutuhan sehari-hari dengan budget tertentu, malas membuka satu-satu website brand untuk membandingkan harga.
2. **Sneakerhead/kolektor** — memburu rilisan terbatas (limited edition, kolaborasi), butuh info restock secepat mungkin dan validasi keaslian harga/diskon.

### Apa yang harus dicapai di versi pertama
- Pengguna bisa mencari dan membandingkan harga sepatu dari minimal beberapa brand lokal dan internasional dalam satu grid.
- Pengguna bisa menyimpan produk (wishlist) dan mengaktifkan alert restock/price-drop setelah login.
- Setiap klik "Beli" menghasilkan redirect ber-tracking affiliate ke toko resmi/marketplace.
- Sistem scraping + data affiliate berjalan otomatis dan terjadwal, menyimpan histori harga harian.

---

## 2. Asumsi dan Pertanyaan Terbuka

| Item | Status | Keterangan | Dampak jika Salah |
|---|---|---|---|
| Affiliate network yang dipakai di awal | ⚠️ Asumsi | Diasumsikan Shopee Affiliate + Tokopedia Affiliate + Involve Asia sebagai fondasi | Sedang — integrasi API bisa berubah tergantung network yang benar-benar disetujui |
| Daftar brand awal yang di-cover | ❓ Belum Jelas | Belum ditentukan brand spesifik mana yang jadi prioritas scraping pertama | Tinggi — menentukan adapter scraping mana yang dibangun duluan |
| Bahasa website | ✅ Konfirmasi | Dwibahasa penuh: Bahasa Indonesia dan Inggris, user bisa switch kapan saja lewat navbar | Rendah |
| Mata uang | ✅ Konfirmasi | Menampilkan IDR dan USD, user bisa switch lewat navbar; harga dasar tetap disimpan dalam IDR di database, USD adalah hasil konversi tampilan | Sedang |
| Kebijakan legal scraping per brand | ❓ Belum Jelas | Tiap brand perlu dicek ToS masing-masing sebelum scraper aktif diproduksi | Tinggi — risiko hukum/pemblokiran IP |
| Verifikasi review komunitas | ⚠️ Asumsi | Review hanya bisa dari user yang login, dengan rate limit dan sistem report | Sedang |
| Provider notifikasi push | ⚠️ Asumsi | Email (Resend) di rilis ini, push notification (OneSignal) disediakan sebagai kanal tambahan | Rendah |
| Sumber kurs mata uang | ⚠️ Asumsi | Kurs IDR-USD diambil dari API kurs pihak ketiga (mis. exchangerate-api atau Bank Indonesia), diperbarui harian | Sedang — kurs yang stale bisa bikin harga USD tidak akurat |
| Bahasa default | ⚠️ Asumsi | Bahasa Indonesia sebagai default, Inggris sebagai pilihan kedua (bukan auto-detect lokasi) | Rendah |

---

## 3. Tujuan Bisnis dan Konversi

| Field | Isi |
|---|---|
| Model bisnis | Affiliate commission (komisi dari klik beli ke toko resmi/marketplace) + fitur premium (restock alert instan) |
| Konversi utama | Klik tombol "Beli" pada halaman detail produk (redirect ber-tracking affiliate) |
| Konversi sekunder | Registrasi akun untuk mengaktifkan wishlist dan alert restock/price-drop |
| Indikator keberhasilan | Jumlah klik-beli (klik-through ke toko), jumlah akun aktif dengan alert terpasang, akurasi data harga (tidak stale) |
| Bukti kepercayaan yang dibutuhkan | Disclosure transparan soal komisi affiliate, indikator "data harga diperbarui otomatis", jumlah brand yang di-cover, review komunitas asli |
| Risiko drop-off utama | Data harga/stok yang tidak akurat (user kecewa kalau harga di KickPick beda dengan toko asli), proses login yang menghalangi user yang cuma mau cek harga cepat |

---

## 4. Target Pengguna

| Tipe Pengguna | Tujuan | Masalah Utama | Aksi yang Diharapkan di Website |
|---|---|---|---|
| Pembeli umum (budget-conscious) | Menemukan sepatu terbaik sesuai budget tanpa buka banyak tab website brand | Harus cek satu-satu website brand, sulit membandingkan harga real-time | Cari & filter produk, lihat detail harga dari berbagai toko, klik beli |
| Sneakerhead/kolektor | Mendapatkan rilisan terbatas secepat mungkin saat restock | Kehabisan stok dalam hitungan menit, sulit memantau banyak brand sekaligus | Login, pasang wishlist & alert restock, terima notifikasi instan, klik beli saat alert masuk |

---

## 5. Scope — Apa yang Dibangun dan Apa yang Tidak

### Dibangun di Rilis Ini (lengkap, siap launch — full build)

- Homepage dengan kategori, brand strip, best seller, trending, rare/limited section
- Halaman search/listing produk dengan filter lengkap (harga, brand, warna, jenis, ukuran, status stok, trending/rare)
- Halaman detail produk: multi-source price comparison, grafik histori harga, size converter antar brand, review komunitas, tombol beli (redirect affiliate)
- Halaman brand directory
- Sistem akun: registrasi, login, verifikasi email, lupa password
- Wishlist dan sistem alert (restock + price-drop), termasuk tier alert instan (premium)
- Sistem review komunitas dengan moderasi dasar (rate limit + report)
- Sistem notifikasi (in-app + email; push disiapkan sebagai kanal opsional)
- Backend scraping modular per-brand + integrasi affiliate network (Shopee, Tokopedia, Involve Asia)
- Job scheduler untuk update harga harian dan cek restock berkala
- Dukungan dwibahasa penuh (Bahasa Indonesia dan Inggris), termasuk switcher bahasa di navbar dan SEO per-bahasa
- Dukungan multi-currency tampilan (Rupiah/IDR dan Dollar/USD), termasuk switcher mata uang di navbar dan konversi kurs otomatis
- Dark mode dan light mode
- Halaman profil, halaman FAQ/disclosure affiliate, halaman kebijakan privasi
- Keamanan penuh: auth aman, rate limiting, validasi input, enkripsi data sensitif, audit log dasar

### Tidak Dibangun (tidak relevan untuk proyek ini)

- Sistem pembayaran/checkout di dalam KickPick — karena model bisnisnya redirect ke toko resmi, bukan menjual produk sendiri
- Iklan banner berbayar dari brand — karena berpotensi merusak positioning "netral & terpercaya" yang sudah disepakati
- Fitur jual-beli sesama user (marketplace C2C) — di luar konsep aggregator/komparasi ini

---

## 6. Sitemap

```
Home (/)
├── Search & Listing (/cari)
├── Detail Produk (/produk/[slug])
├── Brand Directory (/brand)
│   └── Brand Detail (/brand/[slug])
├── Wishlist (/wishlist) — auth
├── Notifikasi (/notifikasi) — auth
├── Profil (/profil) — auth
├── Login (/login)
├── Registrasi (/registrasi)
├── Lupa Password (/lupa-password)
├── Tentang & FAQ (/tentang)
├── Kebijakan Privasi (/privasi)
└── Disclosure Affiliate (/disclosure)
```

---

## 7. Alur Pengguna Utama

```mermaid
flowchart TD
    A([Buka KickPick]) --> B{Fokus onboarding}
    B -- Cari harga terbaik --> C1[Homepage netral]
    B -- Berburu limited edition --> C2[Homepage + drop section diutamakan]
    C1 --> D[Cari & filter produk]
    C2 --> D
    D --> E[Lihat detail produk]
    E --> F{Punya akun?}
    F -- Tidak, cuma cek harga --> G[Klik beli]
    F -- Ya, mau simpan/alert --> H[Login]
    H --> I[Wishlist & alert aktif]
    I --> J[Menunggu notifikasi restock/harga turun]
    J -- Notifikasi masuk --> E
    G --> K([Redirect ke toko resmi])
```

---

## 8. Halaman-Halaman Website

### Home

**Route:** `/`
**Tipe:** Public

**Halaman ini untuk apa:**
Titik masuk utama untuk semua pengguna. Menampilkan gambaran brand yang di-cover, produk populer, dan penawaran terbaik, tanpa condong ke satu segmen (umum vs sneakerhead). Urutan section dapat dipersonalisasi berdasarkan preferensi onboarding atau perilaku pengguna.

**CTA Utama:** Cari sepatu (search bar fungsional di hero)
**CTA Sekunder:** Jelajahi brand

---

#### Hero — Pencarian Langsung

> Titik masuk utama, tanpa badge/eyebrow dekoratif, langsung fungsional.

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Headline singkat (maks 2 baris) menyampaikan value proposition, search bar fungsional (input harga/brand/nama produk), teks kecil non-dekoratif "Membandingkan harga dari [n]+ brand lokal & internasional" |
| Komponen UI | Search input dengan autocomplete, tombol submit |
| Interaksi / Behavior | Ketik lalu submit atau pilih dari autocomplete → redirect ke `/cari?q=...` |
| Sumber data | API — autocomplete dari index produk/brand |
| Animasi / Scroll Behavior | Fade-in halus saat halaman dimuat (durasi 0.4s), tanpa parallax/scroll-scrub |
| Catatan | Tidak ada badge/pill/eyebrow di atas headline. Tidak ada dot dekoratif. Tidak ada teks "BETA"/"V1.0" |

---

#### Brand Strip

> Menampilkan seluruh brand yang di-cover sebagai social proof dan navigasi cepat.

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Logo brand (lokal + internasional) dalam grid/scroll horizontal |
| Komponen UI | Logo grid, scrollable di mobile |
| Interaksi / Behavior | Klik logo → ke `/brand/[slug]` |
| Sumber data | Database — tabel `brands` |
| Animasi / Scroll Behavior | Tidak ada animasi khusus (statis, load langsung) |
| Catatan | Maks 1 baris di desktop dengan tombol "Lihat semua brand" jika lebih dari kapasitas |

---

#### Kategori Cepat

> Filter cepat berbasis tag, terinspirasi struktur Kick Avenue tapi disesuaikan brand voice KickPick.

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Tag kategori: Running, Lifestyle, Training, Formal, dll |
| Komponen UI | Grid ikon + label (ikon dari Tabler Icons) |
| Interaksi / Behavior | Klik → ke `/cari?kategori=...` |
| Sumber data | Statis (daftar kategori tetap) |
| Animasi / Scroll Behavior | Hover: scale 1.03, transisi 0.15s |
| Catatan | — |

---

#### Best Sellers

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Grid 4 kolom (2 kolom di mobile): foto, nama, harga termurah, rating |
| Komponen UI | `ProductCard` |
| Interaksi / Behavior | Klik kartu → `/produk/[slug]` |
| Sumber data | API — diurutkan dari volume klik-beli 30 hari terakhir |
| Animasi / Scroll Behavior | Fade-in staggered saat section masuk viewport (stagger 0.05s) |
| Catatan | Skeleton loading saat data belum siap |

---

#### Trending & Rare/Limited

> Ini BUKAN halaman terpisah — section kurasi di atas grid produk yang sama, sesuai keputusan "1 grid, banyak filter".

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Dua sub-section: "Lagi Trending" (berdasarkan lonjakan klik/pencarian) dan "Rare & Limited" (metadata rilisan terbatas) |
| Komponen UI | `ProductCard` dengan badge kecil "Trending" / "Limited" (bukan dot dekoratif, teks label fungsional) |
| Interaksi / Behavior | Klik kartu → detail produk. Klik "Lihat semua" → `/cari?filter=trending` atau `/cari?filter=rare` |
| Sumber data | API — trending dihitung dari agregasi klik 7 hari terakhir, rare dari flag `is_limited` di database |
| Animasi / Scroll Behavior | Sama seperti Best Sellers |
| Catatan | Section ini tampil untuk semua user (bukan cuma yang fokus sneakerhead), sesuai keputusan homepage netral |

---

#### Price Drop Deals

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Produk dengan penurunan harga terverifikasi (bukan diskon palsu) dalam 7 hari terakhir |
| Komponen UI | `ProductCard` dengan badge persentase turun harga (hijau) |
| Interaksi / Behavior | Klik → detail produk, otomatis scroll ke grafik histori harga |
| Sumber data | API — dihitung dari tabel `price_history`, hanya tampil jika harga turun dari rata-rata 30 hari, bukan sekadar harga naik-turun sesaat |
| Animasi / Scroll Behavior | Sama seperti Best Sellers |
| Catatan | Ini bagian dari fitur fake-discount detector |

---

**Kondisi yang perlu ditangani di halaman ini:**

| Kondisi | Tampilan yang Diharapkan |
|---|---|
| Data kosong / belum ada konten | Section disembunyikan sepenuhnya jika data kosong (bukan tampil kotak kosong) |
| Loading data | Skeleton card sesuai bentuk grid produk |
| Error | Toast "Gagal memuat data. Coba lagi." + tombol retry per section |
| Sukses (untuk pencarian di hero) | Redirect ke `/cari` dengan query terisi |

**Halaman ini selesai ketika:**
- [ ] Semua section di atas tampil dengan data nyata dari API
- [ ] Search bar hero berfungsi dengan autocomplete
- [ ] Homepage personalisasi ringan berjalan (urutan section menyesuaikan onboarding_focus user)
- [ ] Tampil benar di mobile (375px) dan desktop (1280px)
- [ ] Dark mode dan light mode berfungsi penuh
- [ ] Tidak ada badge/eyebrow/dot dekoratif di hero

---

### Search & Listing Produk

**Route:** `/cari`
**Tipe:** Public

**Halaman ini untuk apa:**
Halaman inti KickPick tempat semua produk (reguler maupun limited) ditampilkan dalam satu grid yang bisa difilter. Pengguna datang dari search hero, klik kategori, klik brand, atau klik section trending/rare di homepage.

**CTA Utama:** Klik kartu produk untuk lihat detail
**CTA Sekunder:** Ubah/reset filter

---

#### Filter Sidebar / Bar

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Filter: rentang harga (slider), brand (checkbox multi-select), warna (swatch), jenis sepatu, ukuran, status (tersedia/habis), tag (trending/rare/limited/restock terbaru) |
| Komponen UI | Sidebar di desktop, bottom sheet di mobile |
| Interaksi / Behavior | Filter diterapkan real-time tanpa reload halaman (fade in/out produk) |
| Sumber data | API — hasil query dengan parameter filter |
| Animasi / Scroll Behavior | Produk fade-out lalu fade-in (0.2s) saat filter berubah |
| Catatan | Filter tersimpan di URL query agar bisa di-share/bookmark |

---

#### Grid Produk

| Field | Detail |
|---|---|
| Konten yang ditampilkan | `ProductCard`: foto, nama, brand, rentang harga (termurah–termahal antar toko), rating, badge status |
| Komponen UI | Grid 4 kolom desktop, 2 kolom mobile, infinite scroll atau pagination |
| Interaksi / Behavior | Klik kartu → `/produk/[slug]` |
| Sumber data | API paginated |
| Animasi / Scroll Behavior | Fade-in staggered saat load |
| Catatan | Skeleton loading, empty state jika filter tidak menghasilkan produk |

---

**Kondisi yang perlu ditangani di halaman ini:**

| Kondisi | Tampilan yang Diharapkan |
|---|---|
| Data kosong / belum ada konten | Ilustrasi sederhana + teks "Tidak ada sepatu yang cocok dengan filter ini" + tombol reset filter |
| Loading data | Skeleton grid |
| Error | Toast error + tombol retry |
| Sukses (filter diterapkan) | Grid update, jumlah hasil ditampilkan di atas grid ("124 sepatu ditemukan") |

**Halaman ini selesai ketika:**
- [ ] Semua filter berfungsi dan bisa dikombinasikan
- [ ] Filter tersimpan di URL
- [ ] Pagination/infinite scroll berjalan benar (halaman pertama, halaman terakhir, hasil kosong)
- [ ] Tampil benar di mobile dan desktop

---

### Detail Produk

**Route:** `/produk/[slug]`
**Tipe:** Public

**Halaman ini untuk apa:**
Halaman paling krusial — tempat keputusan beli terjadi. Menampilkan semua diferensiasi KickPick: perbandingan harga multi-toko, histori harga, size converter, dan review komunitas.

**CTA Utama:** Klik Beli (redirect affiliate ke toko dengan harga/ketersediaan pilihan user)
**CTA Sekunder:** Simpan ke wishlist / pasang alert

---

#### Galeri & Info Dasar Produk

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Foto produk (studio, dari sumber data), nama, brand, kategori, rating rata-rata |
| Komponen UI | Image gallery dengan thumbnail |
| Interaksi / Behavior | Klik thumbnail ganti gambar utama |
| Sumber data | Database — tabel `products` + `product_images` |
| Animasi / Scroll Behavior | Tidak ada, statis |
| Catatan | Foto di-crop/resize otomatis agar rasio konsisten meski sumber beda-beda |

---

#### Multi-Source Price Comparison Card

> Komponen unik KickPick — inti dari value proposition.

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Tabel/list toko yang menjual produk ini: nama toko, harga, status stok, badge "Termurah" pada baris terendah |
| Komponen UI | Tabel responsif, badge harga termurah di-highlight |
| Interaksi / Behavior | Klik "Beli" pada baris toko tertentu → redirect ke link affiliate toko tersebut |
| Sumber data | API — agregasi dari `product_offers` (hasil scraping + affiliate feed) |
| Animasi / Scroll Behavior | Tidak ada |
| Catatan | Jika stok habis di semua toko, tampilkan tombol "Aktifkan Alert Restock" menggantikan tombol beli |

---

#### Grafik Histori Harga

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Line chart 30/90 hari, titik harga terendah di-highlight, indikator "diskon asli" jika harga saat ini memang di bawah rata-rata historis |
| Komponen UI | Recharts line chart, toggle 30 hari / 90 hari |
| Interaksi / Behavior | Hover titik menampilkan tooltip harga & tanggal |
| Sumber data | API — dari tabel `price_history` |
| Animasi / Scroll Behavior | Garis chart digambar progresif saat pertama masuk viewport (0.6s) |
| Catatan | Jika data histori kurang dari 30 hari (produk baru di-index), tampilkan pesan "Data histori masih terkumpul" alih-alih grafik kosong |

---

#### Size Converter Widget

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Dropdown pilih brand referensi + ukuran user, otomatis menampilkan ukuran setara di brand produk ini |
| Komponen UI | Dua dropdown (brand referensi, ukuran) + hasil konversi |
| Interaksi / Behavior | Pilih brand & ukuran → hasil konversi muncul instan |
| Sumber data | Database — tabel `size_conversion_matrix`, diperkaya dari agregasi field `fit_feedback` di review komunitas |
| Animasi / Scroll Behavior | Fade-in hasil konversi (0.2s) |
| Catatan | Tampilkan disclaimer kecil: "Estimasi berdasarkan data komunitas, bisa bervariasi per model" |

---

#### Review Komunitas

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Daftar review: rating, teks, tag fit ("kekecilan"/"pas"/"kebesaran"), nama user (username, bukan nama asli) |
| Komponen UI | List review + form tambah review (hanya untuk user login) |
| Interaksi / Behavior | Submit review → validasi → masuk moderasi ringan (rate limit + flag report) |
| Sumber data | Database — tabel `reviews` |
| Animasi / Scroll Behavior | Tidak ada |
| Catatan | Tombol "Laporkan" di tiap review untuk moderasi komunitas. Rate limit: maksimal 1 review per produk per user |

---

**Kondisi yang perlu ditangani di halaman ini:**

| Kondisi | Tampilan yang Diharapkan |
|---|---|
| Data kosong (belum ada offer/toko) | Pesan "Belum ada data harga untuk produk ini" |
| Loading data | Skeleton untuk tiap section |
| Error | Toast error, section tetap tampil dengan data terakhir yang ter-cache jika ada |
| Sukses submit review | Toast "Review berhasil dikirim, menunggu moderasi" |

**Halaman ini selesai ketika:**
- [ ] Semua section (galeri, price comparison, histori harga, size converter, review) berfungsi dengan data nyata
- [ ] Tombol beli menghasilkan redirect affiliate yang benar dan ter-track
- [ ] Alert restock bisa diaktifkan (redirect ke login jika belum login)
- [ ] Tampil benar di mobile dan desktop

---

### Brand Directory

**Route:** `/brand`, `/brand/[slug]`
**Tipe:** Public

**Halaman ini untuk apa:**
Menampilkan seluruh brand yang di-cover KickPick, bisa diklik untuk melihat semua produk brand tersebut.

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Grid logo brand + nama, di halaman detail brand: semua produk brand tersebut (reuse komponen grid dari `/cari`) |
| Komponen UI | Logo grid, product grid |
| Interaksi / Behavior | Klik brand → filter grid produk otomatis by brand |
| Sumber data | Database `brands`, `products` |
| Animasi / Scroll Behavior | Fade-in staggered |
| Catatan | Search bar kecil untuk cari brand by nama |

---

### Wishlist

**Route:** `/wishlist`
**Tipe:** Auth (wajib login)

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Grid produk yang disimpan user, status alert aktif per produk |
| Komponen UI | `ProductCard` + toggle alert per kartu |
| Interaksi / Behavior | Klik toggle untuk aktif/nonaktifkan alert, klik hapus untuk keluarkan dari wishlist |
| Sumber data | API — tabel `wishlists` join `products` |
| Animasi / Scroll Behavior | Fade-out saat item dihapus |
| Catatan | Empty state: "Belum ada sepatu yang disimpan" + tombol "Mulai cari sepatu" |

**Kondisi yang perlu ditangani:**

| Kondisi | Tampilan |
|---|---|
| Wishlist kosong | Empty state seperti di atas |
| Loading | Skeleton grid |
| Error | Toast retry |

---

### Notifikasi

**Route:** `/notifikasi`
**Tipe:** Auth

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Riwayat notifikasi (restock, price-drop), pengaturan kanal notifikasi (email/push) |
| Komponen UI | List notifikasi dengan unread indicator, toggle setting |
| Interaksi / Behavior | Klik notifikasi → ke detail produk terkait, tandai terbaca |
| Sumber data | API — tabel `notifications` |
| Animasi / Scroll Behavior | Tidak ada |
| Catatan | Badge unread count di ikon navbar |

---

### Profil

**Route:** `/profil`
**Tipe:** Auth

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Data akun, ukuran favorit tersimpan per brand, riwayat pencarian, pengaturan onboarding_focus |
| Komponen UI | Form edit profil, list ukuran favorit |
| Interaksi / Behavior | Edit dan simpan data |
| Sumber data | API — tabel `users`, `user_size_preferences` |
| Animasi / Scroll Behavior | Tidak ada |
| Catatan | Termasuk tombol hapus akun (dengan konfirmasi ganda) |

---

### Login / Registrasi / Lupa Password

**Route:** `/login`, `/registrasi`, `/lupa-password`
**Tipe:** Public (redirect ke dashboard/profil jika sudah login)

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Form login (email+password), form registrasi (email, password, konfirmasi password), form lupa password |
| Komponen UI | Form dengan validasi inline |
| Interaksi / Behavior | Submit → validasi → API call → redirect sesuai konteks |
| Sumber data | API `/api/auth/*` |
| Animasi / Scroll Behavior | Shake halus pada field error (0.3s) |
| Catatan | Lihat Section 16 untuk detail auth flow dan keamanan |

**Kondisi:**

| Kondisi | Tampilan |
|---|---|
| Validasi gagal | Pesan error inline per field |
| Email sudah terdaftar | "Email sudah digunakan, coba login" |
| Sukses registrasi | Redirect ke halaman verifikasi email dengan pesan "Cek email kamu" |
| Sukses login | Redirect ke halaman asal (atau homepage) |

---

### Tentang & FAQ / Kebijakan Privasi / Disclosure Affiliate

**Route:** `/tentang`, `/privasi`, `/disclosure`
**Tipe:** Public

| Field | Detail |
|---|---|
| Konten yang ditampilkan | Penjelasan cara kerja KickPick, FAQ, kebijakan privasi lengkap, disclosure transparan soal komisi affiliate |
| Komponen UI | Accordion FAQ, teks statis |
| Interaksi / Behavior | Klik accordion untuk expand/collapse |
| Sumber data | Statis (Markdown/CMS ringan) |
| Animasi / Scroll Behavior | Expand/collapse 0.2s |
| Catatan | Wajib mencantumkan disclosure affiliate secara eksplisit sesuai etika bisnis yang sudah disepakati |

---

### Komponen Global

#### Navbar

**Muncul di:** Semua halaman
**Isi:** Logo KickPick, search bar ringkas, menu (Cari, Brand, Drops & Alert), switcher bahasa (ID/EN), switcher mata uang (IDR/USD), ikon notifikasi (badge unread), ikon wishlist, toggle dark/light mode, avatar/menu akun
**Behavior:** Sticky di scroll, collapse jadi hamburger menu di mobile. Switcher bahasa dan mata uang disatukan dalam satu dropdown kecil untuk hemat ruang di mobile
**Catatan:** Menu "Drops & Alert" selalu ada permanen di navbar (bukan cuma section homepage) agar sneakerhead tetap punya akses cepat, sesuai kesepakatan personalisasi sebelumnya. Pilihan bahasa dan mata uang disimpan di `preferred_language`/`preferred_currency` untuk user login, atau di cookie untuk guest

#### Footer

**Isi:** Logo, link kategori, link brand, link about/FAQ/privasi/disclosure, ikon sosial media, form subscribe update harga (opsional)
**Behavior:** Statis di semua halaman

---

## 9. Sistem Desain

| Aspek | Keputusan |
|---|---|
| Palet warna | Monokrom murni: Off-Black `#0A0A0A`, Pure White `#FFFFFF`, Zinc 100/500/800/950 sebagai variasi abu-abu. **Tidak ada warna aksen apapun** (tidak hijau, biru, kuning, merah). Status (harga termurah, tersedia/habis, alert, turun harga) disampaikan lewat kombinasi bold, ukuran font, ikon Tabler, dan border — bukan warna |
| Mode tampilan | Light mode dan dark mode, toggle di navbar, tersimpan di local state/cookie preferensi |
| Tipografi | `Barlow Condensed` (Bold/SemiBold) untuk headline — grotesque tebal dan sedikit rapat, karakter dekat dengan Trade Gothic Condensed yang dipakai Nike dan AdiHaus DIN yang dipakai Adidas, memberi kesan atletik dan tegas. `Barlow` (Regular) untuk body, satu keluarga font yang sama untuk konsistensi. **Tidak memakai `Inter`, `Poppins`, atau `Montserrat`** |
| Ikon | `@tabler/icons-react` — satu keluarga ikon konsisten di seluruh proyek |
| Animasi | Motion (`motion/react`) — intensitas Balanced: fade-in, hover scale ringan, transisi filter halus. Tidak ada scroll-scrub berat atau parallax dramatis |
| Fotografi produk | Studio bersih, di-crop/resize otomatis dari sumber data untuk konsistensi rasio |
| Radius | Satu skala radius konsisten: card 12px, button 8px, input 8px |

### Larangan Desain (Anti-AI-Slop) — Wajib Dipatuhi

- ❌ Tidak ada badge/pill/eyebrow dekoratif di atas H1 mana pun, termasuk hero homepage
- ❌ Tidak ada dot bulat dekoratif (●) di depan label/nav item/badge kecuali menandakan status nyata (contoh: indikator "data live")
- ❌ Tidak ada label versi ala startup ("BETA", "V1.0", "EARLY ACCESS")
- ❌ Tidak ada em dash (—/–) di semua copy yang tampil ke user
- ❌ Tidak ada filler verbs ("Revolusioner", "Solusi Seamless", "Elevate")
- ❌ Tidak ada 3 kartu fitur simetris identik berturut-turut
- ❌ Tidak ada scroll cue ("Scroll untuk explore")
- ❌ Font default `Inter` + `slate-900` tanpa alasan eksplisit
- ✅ Maksimal 1 eyebrow per 3 section jika benar-benar diperlukan (bukan default)
- ✅ Data contoh (nama reviewer, angka rating) harus realistis, bukan angka bulat sempurna (`4.7` bukan `5.0`, nama user asli bukan "John Doe")

---

## 10. Tech Stack Recommendation

| Layer | Pilihan | Alternatif | Alasan Dipilih | Tradeoff | Status |
|---|---|---|---|---|---|
| Frontend framework | Next.js 15 (App Router) | Remix, SvelteKit | SSR/ISR penting untuk SEO halaman produk, ekosistem terbesar | Bundle lebih besar dari SvelteKit | ✅ Konfirmasi User |
| Bahasa frontend | TypeScript | JavaScript | Type-safety krusial karena data produk dari banyak sumber rawan beda struktur | Sedikit overhead setup | ✅ |
| Styling | Tailwind CSS | CSS Modules | Cepat, konsisten dengan desain monokrom yang disepakati | — | ✅ |
| Komponen UI | shadcn/ui (dikustom total) | Chakra UI | Kode dimiliki sendiri, mudah dikustom agar tidak generik | Perlu disiplin kustomisasi agar tidak terlihat default | ✅ |
| Backend framework | Go + Fiber | Node.js + NestJS | Performa dan concurrency terbaik untuk scraping banyak brand sekaligus, resource server lebih efisien untuk beban scraping berat | Learning curve lebih tinggi, tapi AI coding tool sudah kompeten menulis Go | ✅ Konfirmasi User |
| Scraping engine | Colly (situs statis) + chromedp (situs dinamis/anti-bot) | Playwright (Node) | Resource jauh lebih hemat untuk scraping besar dan berkelanjutan | Ekosistem sedikit lebih kecil dari Playwright | ✅ |
| Job queue/scheduler | Asynq (Redis-based) | BullMQ | Native Go, minim overhead, terintegrasi mulus dengan Fiber | — | ✅ |
| ORM/Query | sqlc | GORM, Prisma | Query SQL native type-safe, performa mendekati raw SQL — penting untuk query price history yang sering diakses | Perlu menulis SQL manual (bukan query builder otomatis) | ✅ |
| Database utama | PostgreSQL | MySQL | JSONB untuk data produk tidak seragam antar brand, query analitik lebih kuat untuk price history | — | ✅ Konfirmasi User |
| Cache/Queue | Redis | — | Cache hasil pencarian, antrian job scraping | — | ✅ |
| Auth | JWT custom (access + refresh token) di Go, terhubung ke NextAuth.js di Next.js sebagai session layer | Supabase Auth | Kontrol penuh atas security policy karena data user + wishlist sensitif | Perlu implementasi manual lebih hati-hati | ✅ |
| Email | Resend | SendGrid | API modern, deliverability baik, harga kompetitif untuk skala awal | — | ✅ |
| Push notification | OneSignal | Firebase Cloud Messaging | Setup lebih cepat, cross-platform | Biaya di skala besar | ✅ |
| Image storage/CDN | Cloudflare R2 | AWS S3 | Biaya egress lebih murah, terintegrasi baik dengan CDN | — | ✅ |
| Search | PostgreSQL full-text search (awal), Meilisearch (jika data membesar) | Algolia | Mulai sederhana dulu, upgrade saat traffic besar | Perlu migrasi index saat upgrade | ⚠️ Asumsi |
| Deployment frontend | Vercel | Netlify | Native untuk Next.js, preview deployment otomatis | — | ✅ |
| Deployment backend | Railway atau Fly.io | VPS manual | Mudah deploy Go binary, scaling sederhana | — | ⚠️ Asumsi |
| CI/CD | GitHub Actions | — | Standar industri, gratis untuk repo publik/kecil | — | ✅ |
| Error tracking | Sentry | — | Standar untuk Go dan Next.js sekaligus | — | ✅ |
| Ikon | @tabler/icons-react | Phosphor Icons | Gaya outline bersih cocok tema monokrom Nike/Adidas-like | — | ✅ |
| Animasi | Motion (`motion/react`) | GSAP | Cukup untuk intensitas Balanced, lebih ringan dari GSAP | — | ✅ |
| Chart | Recharts | Chart.js | Ringan, gampang dikustom untuk line chart histori harga | — | ✅ |
| Form & validasi | React Hook Form + Zod | Formik | Type-safe validasi, terintegrasi baik dengan TypeScript | — | ✅ |
| State management | Zustand | Redux | Ringan, cukup untuk state wishlist/filter | — | ✅ |
| Data fetching | TanStack Query | SWR | Cache otomatis, auto-refresh data harga/stok | — | ✅ |
| Internasionalisasi (i18n) | `next-intl` | `react-i18next` | Terintegrasi native dengan Next.js App Router, mendukung routing per-locale (`/id/...`, `/en/...`) untuk SEO yang benar per bahasa | — | ✅ Konfirmasi User |
| Currency conversion | API kurs pihak ketiga (mis. exchangerate-api.com) + cron harian simpan ke Redis | Hardcode kurs manual | Kurs selalu update otomatis tanpa intervensi manual | Biaya API kecil di skala tinggi | ✅ Konfirmasi User |

---

## 11. Skema Database

```mermaid
erDiagram
  USERS ||--o{ WISHLISTS : has
  USERS ||--o{ REVIEWS : writes
  USERS ||--o{ NOTIFICATIONS : receives
  USERS ||--o{ USER_SIZE_PREFERENCES : sets
  BRANDS ||--o{ PRODUCTS : offers
  PRODUCTS ||--o{ PRODUCT_IMAGES : has
  PRODUCTS ||--o{ PRODUCT_OFFERS : has
  PRODUCTS ||--o{ PRICE_HISTORY : tracks
  PRODUCTS ||--o{ REVIEWS : receives
  PRODUCTS ||--o{ WISHLISTS : saved_in
  PRODUCTS ||--o{ PRODUCT_TRANSLATIONS : has
  STORES ||--o{ PRODUCT_OFFERS : lists

  USERS {
    uuid id PK
    string email
    string password_hash
    string name
    string onboarding_focus
    string preferred_language
    string preferred_currency
    boolean email_verified
    timestamp created_at
  }
  BRANDS {
    uuid id PK
    string name
    string slug
    string logo_url
    boolean is_local
  }
  PRODUCTS {
    uuid id PK
    uuid brand_id FK
    string name
    string slug
    string category
    jsonb attributes
    boolean is_limited
    timestamp created_at
  }
  PRODUCT_IMAGES {
    uuid id PK
    uuid product_id FK
    string url
    int sort_order
  }
  STORES {
    uuid id PK
    string name
    string type
    string affiliate_network
  }
  PRODUCT_OFFERS {
    uuid id PK
    uuid product_id FK
    uuid store_id FK
    decimal price
    string currency
    boolean in_stock
    string size
    string affiliate_url
    timestamp scraped_at
  }
  PRICE_HISTORY {
    uuid id PK
    uuid product_id FK
    uuid store_id FK
    decimal price
    date recorded_date
  }
  REVIEWS {
    uuid id PK
    uuid product_id FK
    uuid user_id FK
    int rating
    text comment
    string fit_feedback
    boolean is_flagged
    timestamp created_at
  }
  WISHLISTS {
    uuid id PK
    uuid user_id FK
    uuid product_id FK
    boolean alert_active
    string alert_type
    timestamp created_at
  }
  NOTIFICATIONS {
    uuid id PK
    uuid user_id FK
    string type
    string title
    string body
    string action_url
    boolean is_read
    timestamp created_at
  }
  USER_SIZE_PREFERENCES {
    uuid id PK
    uuid user_id FK
    uuid brand_id FK
    string size
  }
  PRODUCT_TRANSLATIONS {
    uuid id PK
    uuid product_id FK
    string locale
    string name
    text description
  }
  EXCHANGE_RATES {
    uuid id PK
    string base_currency
    string target_currency
    decimal rate
    date recorded_date
  }
```

---

## 12. API Endpoints

| Method | Endpoint | Deskripsi | Auth |
|---|---|---|---|
| POST | `/api/auth/register` | Registrasi akun baru | Tidak |
| POST | `/api/auth/login` | Login, mengembalikan access & refresh token | Tidak |
| POST | `/api/auth/refresh` | Refresh access token | Refresh token (httpOnly cookie) |
| POST | `/api/auth/logout` | Invalidasi refresh token | Ya |
| POST | `/api/auth/verify-email` | Verifikasi email dari link | Tidak |
| POST | `/api/auth/forgot-password` | Kirim email reset password | Tidak |
| POST | `/api/auth/reset-password` | Set password baru dari token reset | Tidak |
| GET | `/api/products` | List produk dengan filter & pagination | Tidak |
| GET | `/api/products/:slug` | Detail produk lengkap (offers, price history, reviews) | Tidak |
| GET | `/api/products/:slug/price-history` | Histori harga produk | Tidak |
| GET | `/api/products/:slug/size-conversion` | Konversi ukuran antar brand | Tidak |
| GET | `/api/brands` | List semua brand | Tidak |
| GET | `/api/brands/:slug` | Detail brand + produk terkait | Tidak |
| GET | `/api/search/autocomplete` | Autocomplete pencarian | Tidak |
| POST | `/api/reviews` | Submit review baru | Ya |
| POST | `/api/reviews/:id/report` | Laporkan review | Ya |
| GET | `/api/wishlist` | List wishlist user | Ya |
| POST | `/api/wishlist` | Tambah produk ke wishlist | Ya |
| DELETE | `/api/wishlist/:id` | Hapus dari wishlist | Ya (ownership check) |
| PATCH | `/api/wishlist/:id/alert` | Aktif/nonaktifkan alert | Ya (ownership check) |
| GET | `/api/notifications` | List notifikasi user | Ya |
| GET | `/api/notifications/unread-count` | Jumlah notifikasi belum dibaca | Ya |
| PATCH | `/api/notifications/:id/read` | Tandai notifikasi terbaca | Ya (ownership check) |
| GET | `/api/profile` | Data profil user | Ya |
| PATCH | `/api/profile` | Update profil | Ya |
| DELETE | `/api/profile` | Hapus akun (soft delete) | Ya, dengan konfirmasi password |
| POST | `/api/redirect/:offer_id` | Log klik-beli lalu redirect ke affiliate URL | Tidak (tapi di-rate-limit per IP) |

---

## 13. Auth Flow & Role Matrix

```
REGISTRASI:
1. User isi form registrasi (email, password, konfirmasi password)
2. FE validasi client-side (Zod) — email format, password minimal 8 karakter + kombinasi huruf/angka
3. FE POST /api/auth/register
4. BE validasi ulang di server, cek email belum terdaftar
5. BE hash password dengan bcrypt (cost 12), simpan user dengan email_verified=false
6. BE kirim email verifikasi (link berisi token, expired 24 jam)
7. FE redirect ke halaman "Cek email kamu"

VERIFIKASI EMAIL:
8. User klik link di email → /verifikasi?token=xxx
9. FE POST /api/auth/verify-email dengan token
10. BE validasi token, set email_verified=true

LOGIN:
11. User isi form login
12. FE POST /api/auth/login
13. BE validasi kredensial, bandingkan password hash
14. BE generate access_token (JWT, expired 15 menit) dan refresh_token (expired 7 hari)
15. BE simpan refresh_token di DB (hashed), kirim sebagai httpOnly + Secure + SameSite=Strict cookie
16. FE simpan access_token di memory (bukan localStorage)
17. FE redirect ke halaman asal atau homepage

REQUEST DENGAN AUTH:
18. FE kirim access_token di header Authorization: Bearer
19. Jika access_token expired (401), FE panggil /api/auth/refresh
20. BE validasi refresh_token dari cookie, cek masih ada di DB dan belum expired
21. Jika valid: BE kembalikan access_token baru
22. FE retry request original
23. Jika refresh juga invalid: FE logout, redirect ke /login

LOGOUT:
24. FE POST /api/auth/logout
25. BE invalidasi refresh_token di DB
26. FE hapus access_token dari memory, cookie refresh_token dihapus BE
27. Redirect ke /
```

### Role dan Permission Matrix

| Aksi | Guest | User (login) |
|---|---|---|
| Lihat homepage, search, detail produk, brand directory | ✅ | ✅ |
| Klik beli (redirect affiliate) | ✅ | ✅ |
| Tambah wishlist | ❌ | ✅ |
| Aktifkan alert restock/price-drop | ❌ | ✅ |
| Submit review | ❌ | ✅ |
| Laporkan review | ❌ | ✅ |
| Lihat/edit profil sendiri | ❌ | ✅ |
| Lihat data user lain | ❌ | ❌ |
| Hapus akun sendiri | ❌ | ✅ (dengan konfirmasi) |

---

## 14. Keamanan (Security Spec)

> 📌 Section ini wajib dipatuhi penuh mengingat KickPick menyimpan data akun, preferensi pribadi, dan berinteraksi dengan sistem scraping pihak ketiga.

### Autentikasi & Sesi

- Password di-hash dengan **bcrypt cost ≥ 12**, tidak pernah disimpan/dilog dalam bentuk plain text.
- Access token (JWT) berumur pendek (15 menit), refresh token (7 hari) disimpan sebagai **httpOnly, Secure, SameSite=Strict cookie** — tidak pernah diakses lewat JavaScript (mencegah XSS mencuri token).
- Refresh token disimpan di database dalam bentuk **hashed**, bukan plain text, dan bisa di-invalidate manual (logout paksa/ganti perangkat).
- Rate limiting pada endpoint auth: maksimal 10 percobaan login per IP per menit, lockout sementara setelah 5 kali gagal berturut-turut untuk akun yang sama.
- Reset password token berumur pendek (1 jam) dan **single-use** (invalid setelah dipakai sekali).

### Validasi & Sanitasi Input

- Semua input divalidasi di **backend**, bukan hanya frontend (validasi frontend hanya untuk UX, bukan security boundary).
- Gunakan parameterized query / prepared statement (sqlc otomatis menangani ini) untuk mencegah **SQL injection**.
- Sanitasi semua output yang berasal dari user-generated content (review, nama profil) sebelum dirender untuk mencegah **XSS** — escape HTML secara default di React (jangan pernah pakai `dangerouslySetInnerHTML` untuk data user).
- Validasi tipe dan ukuran file untuk upload (jika ada foto profil ke depan): maksimal 2MB, hanya `.jpg/.png/.webp`.

### Proteksi API & Rate Limiting

- Rate limiting per-IP dan per-user di semua endpoint publik (terutama `/api/redirect/:offer_id` untuk mencegah abuse klik-affiliate palsu).
- **CORS** dikonfigurasi ketat — hanya origin domain resmi KickPick yang diizinkan mengakses API.
- **CSRF protection** untuk aksi yang mengubah state (submit review, ubah wishlist) menggunakan token CSRF atau memanfaatkan SameSite cookie + custom header check.
- Setiap endpoint yang mengambil/mengubah data milik user wajib melakukan **ownership check** (memastikan `user_id` di token sama dengan pemilik data yang diakses) — mencegah IDOR (Insecure Direct Object Reference).
- Semua endpoint API berjalan di atas **HTTPS**, tidak ada endpoint yang bisa diakses via HTTP plain.

### Keamanan Data

- Data sensitif (password hash, refresh token) tidak pernah dikirim ke frontend dalam respons API apapun.
- Environment variables (kunci JWT, kredensial database, API key affiliate) disimpan di secret manager platform hosting (Vercel/Railway env vars), tidak pernah di-commit ke repository.
- Backup database terjadwal (harian), retensi minimal 30 hari, backup terenkripsi.
- Kebijakan privasi menjelaskan data apa yang dikumpulkan (dari onboarding, review, wishlist) dan bagaimana digunakan, sesuai UU PDP (Perlindungan Data Pribadi) Indonesia.

### Keamanan Sistem Scraping

- Setiap adapter scraping per-brand menghormati `robots.txt` dan menerapkan rate limiting/delay antar request agar tidak dianggap serangan DDoS oleh server target.
- User-Agent scraper jujur (tidak menyamar sebagai browser biasa jika target melarang crawling non-browser).
- Kredensial/API key affiliate network disimpan terenkripsi, akses dibatasi hanya untuk service scraping (least privilege).
- Setiap brand yang di-scraping ditinjau ToS-nya secara berkala; brand yang secara eksplisit melarang scraping dan tidak tersedia lewat affiliate network dikeluarkan dari daftar cakupan.

### Dependency & Infrastruktur

- Dependency (npm packages, Go modules) di-scan otomatis untuk kerentanan (Dependabot atau `govulncheck` untuk Go, `npm audit` untuk Node/frontend) di setiap CI run.
- Security headers diaktifkan: `Content-Security-Policy`, `X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`, `Strict-Transport-Security`.
- Logging error tidak pernah mencatat data sensitif (password, token) — gunakan redaction di Sentry/logger.
- Audit log dasar untuk aksi sensitif (login, ganti password, hapus akun) disimpan untuk investigasi jika terjadi insiden.

### Security QA Checklist

- [ ] Semua password di-hash bcrypt cost ≥ 12
- [ ] Refresh token httpOnly + Secure + SameSite=Strict
- [ ] Rate limiting aktif di endpoint auth dan redirect affiliate
- [ ] CORS hanya izinkan origin resmi
- [ ] Ownership check di setiap endpoint yang akses data user
- [ ] Tidak ada SQL injection (parameterized query di semua query)
- [ ] Tidak ada XSS (output user-generated content di-escape)
- [ ] HTTPS enforced di semua environment production
- [ ] Environment variables tidak ter-commit ke repository
- [ ] Dependency scanning aktif di CI
- [ ] Security headers terpasang (CSP, X-Frame-Options, dll)
- [ ] Backup database terjadwal dan terenkripsi
- [ ] Kebijakan privasi sesuai UU PDP tersedia dan mudah diakses

---

## 15. Environment Variables

```env
# App
APP_URL=https://kickpick.id
NODE_ENV=production

# Database
DATABASE_URL=postgresql://user:password@host:5432/kickpick

# Redis
REDIS_URL=redis://user:password@host:6379

# Auth
JWT_ACCESS_SECRET=xxxxxxxx
JWT_REFRESH_SECRET=xxxxxxxx
JWT_ACCESS_EXPIRES=15m
JWT_REFRESH_EXPIRES=7d

# Email (Resend)
RESEND_API_KEY=xxxxxxxx
EMAIL_FROM="KickPick <noreply@kickpick.id>"

# Push Notification
ONESIGNAL_APP_ID=xxxxxxxx
ONESIGNAL_API_KEY=xxxxxxxx

# Storage
CLOUDFLARE_R2_ACCESS_KEY=xxxxxxxx
CLOUDFLARE_R2_SECRET_KEY=xxxxxxxx
CLOUDFLARE_R2_BUCKET=kickpick-assets

# Affiliate Networks
SHOPEE_AFFILIATE_API_KEY=xxxxxxxx
TOKOPEDIA_AFFILIATE_API_KEY=xxxxxxxx
INVOLVE_ASIA_API_KEY=xxxxxxxx

# Error Tracking
SENTRY_DSN=xxxxxxxx
```

---

## 16. Error Handling Spec

| Error Scenario | User Experience | Tech Implementation |
|---|---|---|
| Form validasi gagal | Inline error di bawah field, tombol submit disabled | Zod schema + React Hook Form |
| API call gagal (network) | Toast "Gagal memuat. Coba lagi." + tombol retry | try/catch + retry wrapper |
| Session expired | Redirect ke `/login` dengan pesan "Sesi berakhir, silakan login kembali" | Auto-refresh token, fallback redirect |
| Halaman tidak ditemukan | Custom 404 + link kembali ke home | Next.js `not-found.tsx` |
| Server error (5xx) | Toast "Terjadi kesalahan di server" + log ke Sentry | Error boundary + Sentry |
| Data scraping gagal untuk 1 brand | Produk brand tersebut tetap tampil dengan data terakhir yang tersimpan + label "Data mungkin belum terbaru" | Fallback ke cache/data lama, alert internal ke tim |
| Rate limit tercapai | Toast "Terlalu banyak permintaan, coba lagi dalam X detik" | Handle status 429 |
| Koneksi offline | Banner persisten "Tidak ada koneksi internet" | `useOnlineStatus` hook |

**Tools:** Toast — Sonner. Error boundary — komponen `ErrorBoundary`. Error tracking — Sentry.

---

## 17. Notification System Spec

**Channels:** In-app notification center, Email transaksional, Push notification (opsional)

### Email Templates

| ID | Template | Trigger | Subject |
|---|---|---|---|
| E01 | welcome | Registrasi berhasil | "Selamat datang di KickPick!" |
| E02 | verify-email | Registrasi | "Verifikasi email kamu" |
| E03 | reset-password | Lupa password | "Reset password KickPick kamu" |
| E04 | restock-alert | Produk wishlist restock | "[Nama Produk] baru saja restock!" |
| E05 | price-drop-alert | Harga wishlist turun | "Harga [Nama Produk] turun sekarang" |

### In-App Notification Center

- Real-time method: Polling 30 detik (upgrade ke WebSocket jika traffic besar)
- Unread badge di navbar: Ya
- Tabel DB: `notifications(id, user_id, type, title, body, action_url, is_read, created_at)`

---

## 18. Onboarding Flow Spec

**Pattern:** Pertanyaan onboarding singkat 1 langkah (bukan wizard multi-step)

**Setelah registrasi:** Tampilkan pertanyaan "Apa yang paling kamu cari di KickPick?" dengan opsi "Cari harga terbaik" atau "Berburu rilisan limited" → disimpan sebagai `onboarding_focus`, dipakai untuk personalisasi urutan section homepage.

**Empty States:**

| Halaman | Pesan | CTA |
|---|---|---|
| Wishlist | "Belum ada sepatu yang disimpan" | "Mulai cari sepatu" |
| Notifikasi | "Belum ada notifikasi" | — |
| Hasil pencarian kosong | "Tidak ada sepatu yang cocok dengan filter ini" | "Reset filter" |

---

## 19. Caching Strategy

| Data | Cache Layer | Tool | TTL | Invalidasi |
|---|---|---|---|---|
| Halaman produk | Next.js ISR | `revalidate` | 1 jam | `revalidateTag` saat harga baru di-scrape |
| List produk/search | Client | TanStack Query | staleTime 5 menit | Auto saat filter berubah |
| Autocomplete brand/produk | Redis | Redis | 10 menit | TTL otomatis |
| Rate limit counter | Redis | Redis | 60 detik | TTL otomatis |

---

## 20. Testing Strategy

| Layer | Tool | Target Coverage | Dijalankan |
|---|---|---|---|
| Unit | Vitest (FE), Go testing (BE) | 70% logic/utils | Setiap commit (CI) |
| Component | React Testing Library | Semua form + komponen interaktif | Setiap PR |
| Integration | Vitest/Go testing + HTTP client | Semua endpoint kritis | Setiap PR |
| E2E | Playwright | Critical paths di bawah | Sebelum deploy produksi |

**Critical paths wajib E2E:**
- [ ] Registrasi → verifikasi email → login → logout
- [ ] Pencarian & filter produk → detail produk → klik beli (redirect affiliate)
- [ ] Login → tambah wishlist → aktifkan alert
- [ ] Submit review → moderasi rate limit teruji

---

## 21. SEO & Meta Spec

| Halaman | Route | Title Tag | Meta Description | Index? |
|---|---|---|---|---|
| Homepage | `/id` , `/en` | "KickPick — Bandingkan Harga Sepatu Semua Brand" / "KickPick — Compare Shoe Prices Across Brands" | Sesuai locale | Ya |
| Detail Produk | `/id/produk/[slug]`, `/en/produk/[slug]` | "[Nama Produk] — Bandingkan Harga \| KickPick" | Deskripsi dinamis dari data produk sesuai locale | Ya |
| Wishlist, Notifikasi, Profil | — | — | — | Tidak |

- [ ] `app/sitemap.ts` mencakup semua halaman produk & brand, untuk KEDUA locale
- [ ] `app/robots.ts` blokir `/wishlist`, `/notifikasi`, `/profil`, `/api`
- [ ] Structured data `Product` dan `Organization` di halaman relevan
- [ ] Tag `hreflang` (`id`, `en`, `x-default`) dipasang di semua halaman public untuk mencegah duplicate content antar bahasa
- [ ] URL menyertakan locale prefix (`/id/...`, `/en/...`), redirect otomatis ke `/id/` sebagai default kalau user buka root domain tanpa preferensi tersimpan

---

## 22. Deployment Spec

| Pertanyaan | Jawaban |
|---|---|
| Platform hosting frontend | Vercel |
| Platform hosting backend | Railway atau Fly.io |
| Platform database | Neon atau Railway PostgreSQL |
| CI/CD | GitHub Actions |
| Branch strategy | `main` (production) ← `staging` ← `feature/...` |
| Environments | development → preview (per PR) → staging → production |
| Error tracking | Sentry |

### Pre-deploy Checklist

- [ ] Semua test pass di CI
- [ ] Migration sudah ditest di staging
- [ ] Env vars production sudah di-set
- [ ] Security checklist (Section 14) sudah diverifikasi penuh
- [ ] Backup database diambil sebelum migration destructive

---

## 23. QA Checklist Final

**Auth & Security:** lihat Section 14 secara lengkap.

**Data & API:**
- [ ] Pagination benar di semua kondisi (halaman pertama, terakhir, kosong)
- [ ] Filter & search mengembalikan hasil benar
- [ ] Soft delete berjalan untuk akun yang dihapus

**Performance:**
- [ ] Lighthouse score ≥ 90 di halaman public
- [ ] API response time < 500ms read, < 1000ms write
- [ ] Gambar pakai `next/image` dengan lazy load

**Accessibility:**
- [ ] Semua gambar punya alt text
- [ ] Kontras warna ≥ 4.5:1
- [ ] Keyboard navigation berjalan penuh

**Anti-AI-Slop:**
- [ ] Tidak ada badge/eyebrow dekoratif di hero
- [ ] Tidak ada em dash di semua copy
- [ ] Tidak ada 3 kartu fitur simetris identik

---

## 24. Appendix — Daftar Brand Hasil Riset (Cakupan Kandidat)

> 📌 Daftar ini hasil riset pasar brand sepatu lokal dan internasional yang relevan untuk pasar Indonesia. Ini adalah KANDIDAT cakupan, bukan daftar final — prioritas brand mana yang di-scraping/di-integrasikan duluan tetap perlu disepakati terpisah (lihat Section 2, item "Daftar brand awal yang di-cover").

### Brand Lokal Indonesia (~40 brand)

| Kategori | Brand |
|---|---|
| Casual/Lifestyle Sneakers | Compass, Ventela, Aerostreet, Patrobas, Nah Project, Geoff Max, Ando, Kodachi, Wakai, 1-999 (One Triple Nine) |
| Sport/Olahraga & Running | League, 910 (Nineten), Ortuseight, Eagle, Spotec, Kanky, Xpreme, Logic, Mills, Ardiles, Precise, Specs, Warrior |
| Formal/Kulit (termasuk premium handmade) | Brodo, Sagara Bootmaker, Portee Goods, Buccheri, Nappa Milano, Paulmay, Saba, Seis |
| Wanita/Fashion | Donatello, Berrybenka, Pix Footwear |
| Outdoor/Hiking | Eiger, Consina, Forester, Arei |
| Safety Shoes | Krisbow, Safety Ranger |
| Anak-anak/Sepatu Roda | AZKO |

### Brand Internasional (~26 brand)

| Kategori | Brand |
|---|---|
| Sportswear/Performance Global | Nike (termasuk Air Jordan), Adidas, Puma, New Balance, Under Armour, Asics, Reebok, Skechers, Fila |
| Running Specialist | Salomon, On Running, Hoka, Brooks, Saucony, Mizuno |
| Streetwear/Skate/Klasik | Vans, Converse |
| Casual/Formal Internasional | Crocs, Birkenstock, Timberland, Dr. Martens, Clarks |
| Outdoor/Hiking Global | The North Face, Arc'teryx, Mammut, Osprey, Karrimor |

### Catatan Penting dari Riset

- Brand global besar (Nike, Adidas, Puma) kemungkinan besar **tidak punya toko online resmi langsung untuk pasar Indonesia** — cakupan mereka realistisnya lewat **affiliate marketplace** (official store di Shopee/Tokopedia, atau reseller resmi seperti Sport Station/Planet Sports/Foot Locker Indonesia), bukan scraping situs global brand tersebut.
- Brand lokal (Compass, Ventela, Aerostreet, dll) umumnya **punya toko online sendiri SEKALIGUS jualan di marketplace** — artinya satu brand yang sama bisa punya 2 sumber data (scraping situs resmi + feed affiliate marketplace), perlu deduplikasi produk yang sama dari 2 sumber saat digabung ke database.
- Ditemukan kompetitor langsung bernama **TokoSepatu.id** yang sudah punya fitur size chart converter dan quiz rekomendasi sepatu — perlu dipantau sebagai referensi diferensiasi kompetitif.

---

*Dokumen ini adalah rilis satu build lengkap — semua fitur di Section 5 wajib selesai sebelum diklaim siap launch.*
