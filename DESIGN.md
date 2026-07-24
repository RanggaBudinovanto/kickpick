# Design System: KickPick

## Configuration
| Dial | Level |
|---|---|
| Density | 5 |
| Variance | 7 |
| Motion Intent | 5 |

> Archetype: Marketplace Search-First. Density 5 karena halaman search/detail produk cukup padat informasi (filter, tabel perbandingan harga, grafik), tapi tetap diimbangi whitespace di homepage. Variance 7 karena hero dan section homepage sengaja asimetris (bukan template marketplace generik). Motion Intent 5 — Balanced, cukup terasa hidup tapi tidak norak untuk tool yang dipakai orang buru-buru bandingkan harga.

---

## 1. Visual Theme & Atmosphere

Interface monokrom murni yang percaya diri, terinspirasi lantai pamer sepatu premium: hanya hitam, putih, dan abu-abu, tanpa satu pun warna aksen. Atmosfernya klinis dan objektif, seperti etalase yang membiarkan produk itu sendiri jadi satu-satunya warna, sementara UI-nya sendiri diam dan netral. Hierarki dan makna dibawa lewat kontras, berat tipografi, ikon, dan garis, bukan warna. Setiap section punya tujuan fungsional yang jelas, tidak ada dekorasi tanpa makna.

### Referensi Visual Langsung (DNA yang wajib terasa di setiap layar)

- **Fotografi produk sebagai hero** (ala Nike/Adidas): foto sepatu besar, kualitas studio tinggi, sedikit miring/dinamis (bukan foto lurus datar) — foto adalah elemen paling menonjol di tiap card, bukan teks atau chrome UI.
- **Tipografi besar dan tebal, kata sedikit** (ala Nike/Adidas): headline hero tidak lebih dari 2 baris, ukuran besar (minimal 48px desktop), tanpa kalimat panjang menjelaskan.
- **Harga dengan hierarki eksplisit** (ala Kick Avenue): kalau ada harga asli vs harga diskon, harga asli dicoret (strikethrough) dan lebih kecil/abu-abu, harga final ditebalkan dan lebih besar — persis pola "coret-tebal" yang dipakai Kick Avenue, hanya saja tanpa warna merah untuk badge diskon (ganti dengan teks tebal + border, sesuai aturan monokrom di § 2).
- **Grid produk padat tapi rapi** (ala Kick Avenue): jarak antar card rapat terkontrol, banyak produk terlihat sekaligus dalam satu layar tanpa terasa sesak — bukan grid renggang ala galeri seni.
- **Filter berbasis kurasi, bukan cuma kategori** (ala Kick Avenue "Top 50"/"Under Retail"/"For Her"): tag filter cepat di homepage dan listing terasa "dipilihkan", bukan sekadar daftar kategori generik.
- **Chrome UI minim** (ala Nike/Adidas): navbar sederhana, tidak ada border/divider berlebihan antar section, whitespace yang disengaja memisahkan section alih-alih garis pembatas.

---

## 2. Color Palette & Roles

> Tidak ada warna aksen sama sekali. Palet ini murni grayscale — semua makna dan hierarki dibawa oleh kontras, berat font, ikon, dan border, bukan warna.

- **Off-Black** (#0A0A0A) — teks utama di light mode, background utama di dark mode, CTA primary (fill hitam + teks putih)
- **Pure White** (#FFFFFF) — background utama di light mode, teks utama di dark mode
- **Zinc 100** (#F4F4F5) — background section sekunder, latar badge status netral (di light mode)
- **Zinc 800** (#27272A) — card dan border di dark mode
- **Zinc 500** (#71717A) — teks sekunder/metadata di kedua mode
- **Zinc 950** (#09090B) — latar badge status "penting/habis" (kontras tertinggi selain hitam murni)

Aturan: TIDAK ADA warna aksen brand dalam bentuk apapun (tidak hijau, tidak biru, tidak kuning). Status yang biasanya diwakili warna (harga termurah, tersedia/habis, alert) diwakili dengan kombinasi: bold + ukuran font, ikon Tabler yang relevan, border/underline, dan variasi shade abu-abu — lihat § 4 untuk detail per komponen.

---

## 3. Typography Rules

> Font sebelumnya (`Geist`) diganti karena terlalu "SaaS netral" — kurang punya karakter atletik. Nike memakai font custom berbasis **Trade Gothic Bold Condensed** ("Nike TG"), Adidas memakai font custom berbasis **DIN** ("AdiHaus DIN"). Keduanya proprietary dan tidak bisa dipakai langsung, jadi dipilih font gratis dengan DNA visual paling dekat: grotesque tebal, sedikit rapat/condensed, tegas.

### Font Family

- **Display/Headline**: `Barlow Condensed` — SemiBold (600) untuk H2-H6, Bold (700) untuk H1/hero. Karakternya rapat dan tegas, dekat dengan Trade Gothic Condensed (Nike) dan DIN (Adidas).
- **Body**: `Barlow` (bukan versi condensed) — Regular (400) untuk paragraf, Medium (500) untuk label/emphasis kecil. Satu keluarga font yang sama dengan headline supaya tetap terasa satu sistem, tapi lebar normal untuk keterbacaan di ukuran kecil.
- **Mono**: `JetBrains Mono` atau `Geist Mono` — WAJIB untuk semua angka harga, persentase diskon, dan data grafik histori harga. Gunakan `font-variant-numeric: tabular-nums` supaya digit selalu rata lebar yang sama (krusial untuk alignment tabel perbandingan harga multi-toko).
- Tidak ada serif di manapun. Tidak ada `Inter`/`Poppins`/`Montserrat`.

### Type Scale (Desktop → Mobile)

| Level | Font | Weight | Desktop | Mobile | Line-height | Tracking |
|---|---|---|---|---|---|---|
| H1 (Hero) | Barlow Condensed | Bold 700 | 56px | 36px | 1.05 | -0.01em |
| H2 (Section title) | Barlow Condensed | SemiBold 600 | 36px | 26px | 1.1 | -0.005em |
| H3 (Card/subsection title) | Barlow Condensed | SemiBold 600 | 22px | 18px | 1.2 | normal |
| Body | Barlow | Regular 400 | 16px | 15px | 1.6 | normal |
| Small/Caption (metadata, label) | Barlow | Medium 500 | 13px | 13px | 1.4 | 0.01em (sedikit renggang untuk uppercase label) |
| Harga Utama (Mono) | JetBrains Mono | Bold 700 | 28px | 22px | 1.2 | normal, tabular-nums |
| Harga Sekunder/Dicoret (Mono) | JetBrains Mono | Regular 400 | 15px | 14px | 1.2 | normal, tabular-nums |

### Aturan Hierarki

- Headline hero **maksimal 2 baris**, tidak ada kalimat panjang — biarkan `Barlow Condensed` Bold besar yang membawa dampak visual, bukan banyak kata.
- Body text maksimal **65 karakter per baris** untuk deskripsi produk/review, supaya tetap nyaman dibaca meski leading relaks (1.6).
- Label/caption uppercase (misal "TERSEDIA", "TREND") memakai tracking sedikit lebih renggang (0.01em) dan Medium weight — bukan Bold, supaya tidak bersaing dengan H1/H2 sebagai fokus utama halaman.
- Angka harga selalu memakai font mono dengan `tabular-nums` — ini non-negotiable karena tabel perbandingan harga multi-toko HARUS rata secara visual tanpa digit "meloncat" lebar antar baris.

---

## 4. Component Stylings

- **Buttons**: flat, tanpa outer glow. Primary = fill Off-Black + teks putih (light mode), fill putih + teks off-black (dark mode). Secondary = ghost/outline zinc. Active state: `scale(0.98)` saat ditekan.
- **Product Card**: foto studio di atas, nama + brand + rentang harga (mono font) + rating di bawah, shadow tipis di-tint zinc (bukan hitam polos), rounded 12px. Badge kecil pojok kanan atas HANYA kalau ada makna (trending/limited/turun harga) — teks label + ikon Tabler kecil, latar Zinc 100 (light) / Zinc 800 (dark), tidak ada warna.
- **Multi-Source Price Table**: baris toko dengan harga (mono). Baris termurah dibedakan lewat **border 2px solid hitam di sekeliling baris** + harga di baris itu memakai font size lebih besar dan bold — bukan warna latar.
- **Status Badge "Tersedia"**: ikon centang (Tabler `ti-check`) + teks "Tersedia", latar Zinc 100.
- **Status Badge "Habis"**: ikon silang (Tabler `ti-x`) + teks "Habis", latar Zinc 950 + teks putih (kontras tertinggi, menandakan "tidak bisa diklik").
- **Badge "Turun Harga"**: ikon panah bawah (Tabler `ti-arrow-down`) + teks persentase ("Turun 23%"), latar Zinc 100, teks bold.
- **Inputs/Forms**: label di atas input, border zinc, focus ring hitam (ring 2px), error text tetap pakai teks bold + ikon peringatan (bukan warna merah) di bawah field.
- **Loading states**: skeletal shimmer abu-abu yang match dimensi card/grid asli — bukan spinner lingkaran.
- **Empty states**: komposisi ikon sederhana (Tabler Icons, garis outline hitam) + headline singkat + CTA jelas, bukan sekadar teks "Tidak ada data".
- **Error states**: inline, kontekstual, ikon peringatan outline + teks bold, dengan aksi pemulihan (tombol retry).
- **Live/Verified Indicator**: satu-satunya "dot" yang diizinkan — abu-abu gelap (Zinc 500), berdenyut halus, menandakan "data harga diperbarui otomatis". Tetap monokrom, bukan warna.

---

## 5. Layout Principles

- **Grid-first**: CSS Grid untuk grid produk dan filter, bukan flexbox percentage math.
- **No overlapping**: setiap elemen (card, badge, chart) punya zona spasialnya sendiri.
- **Hero**: TIDAK centered (Variance 7 > 4) — layout asimetris: headline + search bar di kiri (60%), elemen visual pendukung (misal ilustrasi ringan brand strip atau grafik tren harga mini) di kanan (40%).
- **Feature/kategori section**: tidak pernah 3 kartu setara sejajar — pakai bento grid asimetris atau grid 4-6 kolom untuk kategori dengan ukuran bervariasi.
- **Containment**: max-width 1400px, centered.
- **Full-height**: `min-height: 100dvh` untuk section hero.

---

## 6. Responsive Rules

- Mobile (<768px): filter sidebar jadi bottom sheet, grid produk 2 kolom, hero jadi 1 kolom stack (search bar di atas, visual di bawah atau disembunyikan).
- Tidak ada horizontal scroll kecuali brand strip (scroll horizontal sengaja untuk logo brand).
- Touch target minimal 44px.
- Testing viewport: 375px, 768px, 1024px, 1440px.

---

## 7. Motion & Interaction (Code-Phase Intent)

- Physics engine: spring-based (`stiffness: 100, damping: 20`) untuk hover dan transisi filter.
- Perpetual micro-loop: HANYA pada dot indikator "live/verified" (pulse halus, 2s loop) — tidak ada micro-loop dekoratif lain.
- Staggered orchestration: grid produk muncul dengan cascade delay 0.05s per card saat scroll masuk viewport.
- Grafik histori harga: garis digambar progresif (path draw-in) saat pertama masuk viewport, durasi 0.6s.
- Hardware rule: animasikan hanya `transform` dan `opacity`.
- Tidak ada interaksi lanjutan immersive (no 3D tilt, no magnetic cursor, no particle background) — Motion Intent 5 tidak memerlukan itu, produk ini fungsional bukan campaign/portfolio.

---

## 8. Anti-Generic Signals (minimal 3 dari 5 diterapkan)

1. **Editorial Scale** — angka harga termurah di halaman detail produk ditampilkan jauh lebih besar dari harga toko lain, menciptakan hierarki visual yang jelas tanpa perlu warna.
2. **Intentional Restraint** — tidak ada warna aksen sama sekali adalah bentuk restraint paling besar di proyek ini; tidak ada border/background berlebih pada hero, biarkan whitespace dan tipografi besar yang bicara.
3. **Texture & Depth** — shadow di-tint zinc (bukan hitam polos) pada product card, memberi kedalaman halus tanpa terasa berat, meski seluruhnya grayscale.

---

## 9. Anti-Patterns (Banned)

**Visual & Warna**
- **Tidak ada warna aksen dalam bentuk apapun** — tidak hijau, biru, kuning, merah, ungu, atau warna lain di luar hitam/putih/abu-abu. Ini larangan mutlak untuk proyek ini, bukan "1 aksen saja" seperti default biasa.
- Status (harga termurah, tersedia/habis, alert) HARUS dibawa lewat kontras, bold, ikon, dan border — tidak pernah lewat warna
- Tidak ada emoji di mana pun
- Tidak ada font `Inter`, `Poppins`, atau `Montserrat`
- Tidak ada serif generik maupun `Fraunces`/`Instrument Serif`
- Tidak ada pure black (`#000000`) — pakai Off-Black `#0A0A0A`
- Tidak ada neon glow/outer glow shadow
- Tidak ada gradient warna apapun

**Layout & Komposisi**
- Tidak ada badge/pill/eyebrow dekoratif di atas heading manapun, terutama hero homepage
- Tidak ada dot dekoratif kecuali live/verified indicator yang disebut di § 4
- Tidak ada 3-kolom kartu fitur setara
- Tidak ada hero centered
- Maksimal 1 eyebrow per 3 section (dan disarankan nol untuk homepage)
- Tidak ada 3 section berturut-turut dengan pola split image-teks yang sama

**Copy & Konten**
- Tidak ada em dash (`—`/`–`) di teks manapun
- Tidak ada filler text ("Scroll untuk explore")
- Tidak ada nama/data generik ("John Doe", angka bulat sempurna seperti "50% off" — pakai angka realistis seperti "34% off")
- Tidak ada klise copywriting ("Revolusioner", "Solusi Seamless")

**Production-Test Tells**
- Tidak ada label versi ("BETA", "V1.0")
- Tidak ada section-number eyebrow ("001 · Kategori")
- Tidak ada fake product UI dari div (fake dashboard/terminal)
- Tidak ada scroll cue

---

*File ini adalah kontrak visual KickPick. Setiap prompt layar di bawah merujuk balik ke sini — jangan mengubah warna/font/dial per layar.*
