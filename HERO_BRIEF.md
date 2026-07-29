# Kriteria Visual Hero Homepage KickPick

> Merujuk balik ke `DESIGN.md` § 1, § 5, § 9 — jangan menyimpang dari sistem warna/komposisi yang sudah ditetapkan.

## 1. Konteks penempatan

Hero homepage **asimetris, bukan centered**:
- **Kiri (60%)**: headline besar (maks 2 baris, Barlow Condensed Bold) + search bar
- **Kanan (40%)**: **visual pendukung** — ini yang perlu dibuatkan asetnya

Saat ini area kanan masih placeholder teks kosong ("Visual pendukung"). Kita mau isi dengan gambar sungguhan.

Section hero pakai `min-height: 100dvh` (nyaris penuh layar), jadi visual ini harus cukup besar/kuat, bukan elemen kecil yang tenggelam.

## 2. Arah konten visual (pilih salah satu, atau kombinasi)

Sesuai DNA referensi di `DESIGN.md` § 1 ("Fotografi produk sebagai hero, ala Nike/Adidas"):

**Opsi A — Foto sepatu studio (paling sesuai DNA brand)**
- Satu sepasang sneakers, kualitas studio tinggi, **sedikit miring/dinamis** (bukan foto lurus datar dari depan)
- Background polos (putih/abu-abu terang atau gradient monokrom halus)
- Pencahayaan dramatis tapi bersih — bayangan tegas, bukan flat lighting e-commerce biasa

**Opsi B — Brand strip / grafik tren harga mini**
- Alternatif yang disebut eksplisit di `DESIGN.md` § 5 kalau tidak pakai foto produk
- Baris logo brand yang di-scroll horizontal, ATAU grafik garis tren harga minimalis (garis tipis hitam di atas putih, tanpa grid/axis berlebihan)

**Rekomendasi**: Opsi A untuk kesan premium/editorial yang lebih kuat sesuai § 1, kombinasi dengan sedikit elemen data (misal angka harga kecil mengambang di sudut, mono font) untuk mengisyaratkan "perbandingan harga" tanpa perlu ilustrasi rumit.

## 3. Aturan mutlak (non-negotiable)

- **Monokrom murni** — hitam, putih, abu-abu saja. Tidak ada warna sepatu yang mencolok (kalau foto sepatu, sepatu itu sendiri harus berwarna netral/hitam-putih, atau foto di-treatment jadi grayscale)
- **Tidak ada badge/pill/eyebrow dekoratif** menempel di gambar
- **Tidak ada gradient warna** (gradient grayscale/monokrom untuk background boleh, gradient warna tidak boleh)
- Foto harus terasa **dinamis, bukan statis** — sedikit sudut, bukan produk difoto lurus dari depan seperti katalog biasa

## 4. Kebutuhan teknis

- Rasio kira-kira **4:5 sampai 1:1** (kolom kanan 40% di desktop, cenderung agak tinggi bukan landscape lebar)
- Harus tetap kuat di **mobile** — di layar <768px, hero jadi 1 kolom stack dan visual ini bisa disembunyikan atau ditaruh di bawah search bar (jadi tidak wajib sangat detail di ukuran kecil)
- Kalau berupa foto (bukan generated art), butuh resolusi tinggi (minimal 1600px sisi terpanjang) supaya tidak pecah di retina display
- Background sebaiknya bisa **transparent atau seamless dengan warna section** (`bg-background`) supaya tidak ada kotak/frame terlihat pada gambar

## 5. Prompt siap pakai (untuk AI image generator — Midjourney/DALL-E/dsb, dalam Bahasa Inggris karena umumnya lebih akurat)

```
Premium studio product photography of a single pair of monochrome sneakers,
shot at a slight dynamic angle (not straight-on catalog shot), dramatic
directional lighting with crisp defined shadows, pure grayscale palette
(black, white, gray only, no color accents whatsoever), clean minimal
background with soft gradient from white to light gray, high-end editorial
sneaker showroom aesthetic, sharp focus on shoe texture and materials,
negative space on one side for text overlay, 4:5 aspect ratio, ultra
high resolution, commercial photography style similar to Nike/Adidas hero
campaigns but entirely monochrome
```

**Varian negative prompt (kalau generator-nya mendukung)**:
```
no color, no colored shoes, no logos of real brands, no watermark, no text,
no busy background, no multiple products, no flat lighting, no centered
straight-on angle
```

## 6. Yang harus dihindari

- Foto sepatu berwarna-warni (harus di-grayscale atau memang sepatu netral hitam/putih/abu)
- Komposisi datar/simetris/lurus dari depan (bertentangan dengan "Motion Intent 5" dan kesan dinamis di § 1)
- Ilustrasi kartun/flat design generik (bertentangan dengan kesan "premium showroom")
- Logo brand sepatu asli yang mencolok/dominan di foto (risiko hukum kalau dipakai sebagai aset komersial situs)
- Background ramai/bertekstur berat yang mengalihkan perhatian dari sepatu

---

*Referensi lengkap: `DESIGN.md` § 1 (DNA visual), § 5 (aturan layout hero), § 9 (anti-pattern yang dilarang).*
