# Migrasi Database

Direktori ini berisi file migrasi SQL untuk database PostgreSQL (Supabase).

## Cara Menjalankan Migrasi

Migrasi ditulis menggunakan standar SQL murni dan dapat dijalankan pada database PostgreSQL mana pun, termasuk instance lokal, Docker container, maupun Supabase.

### 1. Menggunakan `psql` Standar (Lokal / Uji Coba)
Jika Anda memiliki PostgreSQL yang berjalan secara lokal, Anda dapat menerapkan migrasi `UP` dengan perintah berikut:

```bash
psql -U postgres -d labsync -f migrations/000001_init_schema.up.sql
```

Untuk membatalkan (rollback):
```bash
psql -U postgres -d labsync -f migrations/000001_init_schema.down.sql
```

### 2. Menggunakan Supabase CLI (Direkomendasikan)
Jika Anda menginstal Supabase CLI, Anda dapat menggunakan file ini. Supabase mengelola migrasi di dalam folder `supabase/migrations`, Anda cukup memindahkan atau menyalin file dari sini ke folder tersebut.

```bash
# Jalankan supabase lokal
supabase start

# Terapkan migrasi (asumsi file sudah berada di supabase/migrations)
supabase db reset
```

## Constraint yang Diterapkan
- **Roles & Users**: Relasi role dasar. `users.status` dibatasi hanya ACTIVE/INACTIVE.
- **Face Templates**: Menggunakan `Partial Unique Index` untuk memastikan satu user hanya dapat memiliki **satu** set pendaftaran wajah yang aktif pada satu waktu.
- **Face Template Samples**: `Unique Index` memastikan di dalam satu set pendaftaran, hanya boleh ada tepat satu pose FRONT, satu LEFT, dan satu RIGHT.
- **Sessions & Attendance**: `CHECK` constraint memastikan `end_at > start_at` dan `checked_out_at >= checked_in_at`. Terdapat juga `Unique Index` untuk mencegah presensi ganda bagi user pada sesi yang sama.
- **Assets & Allocations**: `Partial Unique Index` pada tabel alokasi dengan kondisi `returned_at IS NULL` menjamin bahwa sebuah aset tidak dapat dialokasikan kepada lebih dari satu pengguna secara bersamaan.
- **Maintenance**: Validasi status standar (OPEN, IN_PROGRESS, COMPLETED, CANCELLED) dan validasi tingkat prioritas.

## Seed Development

Seeder development ada di `backend/cmd/seed/main.go` dan didokumentasikan di `database/seeds/README.md`.

Jalankan dari folder `backend/`:

```bash
go run cmd/seed/main.go
```

Akun demo default:
- Admin: `ADM-DEMO-001` / `AdminDemo123!`
- Mahasiswa: `MHS-DEMO-001` / `StudentDemo123!`

Seeder hanya untuk development, ditolak saat `APP_ENV=production`, dan password dapat dioverride lewat `DEV_SEED_ADMIN_PASSWORD` / `DEV_SEED_STUDENT_PASSWORD`.
