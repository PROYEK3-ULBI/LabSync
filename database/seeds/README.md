# Development Seed

Seed data development dijalankan dari Go agar password di-hash memakai bcrypt yang sama dengan backend.

Jalankan dari folder `backend/` setelah `backend/.env` berisi `DATABASE_URL` dev:

```bash
go run cmd/seed/main.go
```

Data yang dibuat:

- Admin demo: `ADM-DEMO-001` / `AdminDemo123!`
- Mahasiswa demo: `MHS-DEMO-001` / `StudentDemo123!`
- 2 sesi praktikum demo (`ACTIVE` dan `SCHEDULED`)
- 3 aset demo (`AVAILABLE`, `AVAILABLE`, `MAINTENANCE`)
- 1 tiket maintenance demo

Ganti password seed tanpa mengubah source:

```bash
DEV_SEED_ADMIN_PASSWORD='password-admin-kuat' DEV_SEED_STUDENT_PASSWORD='password-mahasiswa-kuat' go run cmd/seed/main.go
```

Catatan keamanan:

- Seed ditolak saat `APP_ENV=production`.
- Jangan gunakan akun/password demo di production.
- Jangan commit `.env` atau password asli.
- Seeder idempotent: aman dijalankan ulang untuk refresh data demo.

Membuat akun admin tambahan:

- Cara UI/API: login sebagai admin demo, lalu buat user baru role `ADMIN` lewat endpoint/user-management admin.
- Cara sementara via backend: gunakan `POST /api/v1/users` dengan cookie admin; backend akan hash password dan validasi role.
- Jangan insert langsung ke SQL kecuali menghasilkan bcrypt hash yang valid.

Membuat akun mahasiswa:

- Mahasiswa dibuat oleh admin lewat `POST /api/v1/users` dengan role `STUDENT`.
- Frontend admin Prompt 15 akan menghubungkan layar `Pengguna` ke endpoint ini.
- Tidak ada self-registration mahasiswa pada MVP.
