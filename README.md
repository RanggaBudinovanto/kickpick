# KickPick

Platform pencarian dan perbandingan harga sepatu lintas brand. Lihat [KickPick_PRD.md](KickPick_PRD.md) untuk spesifikasi produk lengkap, [DESIGN.md](DESIGN.md) untuk kontrak visual, dan [PLAN.md](PLAN.md) untuk roadmap pembangunan.

## Struktur Monorepo

```
frontend/   Next.js 15 (App Router) + TypeScript + Tailwind CSS
backend/    Go + Fiber API, scraping engine, job scheduler
```

## Menjalankan Secara Lokal

### Frontend

```bash
cd frontend
npm install
cp .env.example .env.local
npm run dev
```

Berjalan di `http://localhost:3000`, redirect otomatis ke `/id` (locale default).

### Backend

```bash
cd backend
go mod tidy
cp .env.example .env
go run ./cmd/api
```

Berjalan di `http://localhost:8080`. Tanpa `DATABASE_URL` yang valid, server tetap jalan (untuk cek routing/health), tapi endpoint yang butuh data akan gagal.

### Database (PostgreSQL lokal)

Repo ini sudah ditest jalan dengan PostgreSQL 17 lokal (via `winget install PostgreSQL.PostgreSQL.17`, default password superuser `postgres`). Setup dari nol:

```bash
psql -U postgres -h 127.0.0.1 -c "CREATE DATABASE kickpick;"
psql -U postgres -h 127.0.0.1 -c "CREATE USER kickpick WITH PASSWORD 'kickpick_dev_pw';"
psql -U postgres -h 127.0.0.1 -c "GRANT ALL PRIVILEGES ON DATABASE kickpick TO kickpick;"
psql -U postgres -h 127.0.0.1 -d kickpick -c "GRANT ALL ON SCHEMA public TO kickpick;"

psql -U kickpick -h 127.0.0.1 -d kickpick -f backend/internal/db/migrations/0001_init.up.sql
psql -U kickpick -h 127.0.0.1 -d kickpick -f backend/internal/db/migrations/0002_product_views.up.sql
```

Set `DATABASE_URL=postgresql://kickpick:kickpick_dev_pw@127.0.0.1:5432/kickpick?sslmode=disable` di `backend/.env`.

### Seed Data

```bash
cd backend
go run ./cmd/seed
```

Mengisi 6 brand awal (Compass, Ventela, Aerostreet, Nike, Adidas, Vans), beberapa produk per brand, offers dari beberapa toko, 30 hari histori harga, dan 1 user demo (`demo@kickpick.id` / `Password123`).

### Generate Query Type-Safe (sqlc)

```bash
cd backend
CGO_ENABLED=0 go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate
```

> Catatan: `CGO_ENABLED=0` wajib di environment ini karena toolchain MinGW gcc lokal gagal compile parser cgo `pg_query_go`. Dengan CGO dimatikan, sqlc otomatis pakai parser SQL berbasis WASM (tidak perlu gcc).

### Scraping

```bash
cd backend
go run ./cmd/scraper   # jalankan sekali secara manual
go run ./cmd/worker    # jalan terus, scrape sekali saat start lalu terjadwal harian jam 03:00
```

Saat ini baru brand **Compass** yang punya adapter scraping nyata (`internal/scraper/compass/`) — lihat [PENDING.md](PENDING.md) untuk alasan teknis kenapa 5 brand kandidat lainnya belum bisa di-scrape. Jadwal cron bisa diubah lewat env var `SCRAPE_CRON_SPEC` (format cron standar, default `"0 3 * * *"`).

### Testing

```bash
# Backend: unit tests jalan tanpa database, integration tests butuh TEST_DATABASE_URL
cd backend
go test ./...
TEST_DATABASE_URL="postgresql://kickpick:kickpick_dev_pw@127.0.0.1:5432/kickpick_test?sslmode=disable" go test ./...

# Frontend: unit/component test (Vitest)
cd frontend
npm test

# Frontend: E2E (Playwright) — butuh backend (:8080) dan database jalan
cd frontend
npx playwright install chromium   # sekali saja
npm run test:e2e
```

CI (`.github/workflows/ci.yml`) menjalankan semuanya otomatis di tiap push/PR ke `main`/`staging`, termasuk `govulncheck` (Go) dan `npm audit` (Node).

## Deployment

Sesuai Section 22 PRD: frontend ke Vercel (native Next.js), backend ke Railway atau Fly.io.

- **Frontend**: hubungkan repo ke Vercel, set `NEXT_PUBLIC_API_URL` dan `NEXT_PUBLIC_APP_URL` ke domain production. Vercel otomatis build & deploy tiap push.
- **Backend API**: deploy `backend/Dockerfile` (expose port 8080) ke Railway/Fly.io. Set semua env var dari `backend/.env.example`.
- **Worker (scraper scheduler)**: deploy `backend/Dockerfile.worker` sebagai service terpisah dari API — jangan gabung jadi satu proses, supaya durasi scraping tidak mempengaruhi latency HTTP.
- **Database**: Neon atau Railway PostgreSQL, jalankan migrasi (`internal/db/migrations/*.up.sql`) sebelum deploy pertama.

`docker-compose.yml` di root tersedia sebagai referensi untuk menjalankan Postgres + API + worker via container secara lokal — **belum divalidasi** dengan `docker build` sungguhan di environment ini (Docker tidak terpasang saat development), lihat [PENDING.md](PENDING.md).

### Pre-deploy Checklist (Section 22 PRD)

- [x] Semua test pass di CI
- [ ] Migration sudah ditest di staging (baru ditest di local dev DB)
- [ ] Env vars production sudah di-set
- [x] Security checklist (Section 14) sudah diverifikasi — lihat [PENDING.md](PENDING.md) untuk item yang masih tersisa (Redis-backed rate limiting, dll)
- [ ] Backup database diambil sebelum migration destructive

## Status Pembangunan

Lihat [PLAN.md](PLAN.md) untuk fase pembangunan (Plan 1-4). Lihat [PENDING.md](PENDING.md) untuk daftar fitur yang masih ditunda/disederhanakan dan perlu dikerjakan di fase berikutnya.
