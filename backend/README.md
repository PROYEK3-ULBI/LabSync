# LabSync Backend

Backend REST API untuk sistem presensi lab dan manajemen aset berbasis Go + Fiber.

## Kebutuhan Sistem
- Go 1.22+ (teruji pada Go 1.27)
- PostgreSQL (Lokal atau Supabase)

## Menjalankan Layanan

### 1. Salin Environment
```bash
cp .env.example .env
```
Sesuaikan nilai `DATABASE_URL` dengan database lokal atau instance Supabase Anda.

### 2. Jalankan Server
```bash
go run cmd/server/main.go
```
Server akan aktif di `http://localhost:8000`.

### 3. Endpoint yang Tersedia
- **System & Docs:**
  - `GET /api/v1/health`: Health & readiness probe.
  - `GET /api/v1/openapi.json`: OpenAPI 3.0 specification.
  - `GET /api/v1/docs`: Swagger UI.
- **Autentikasi (`/api/v1/auth`):**
  - `POST /api/v1/auth/login`: Login akun (menerbitkan cookie `token` HttpOnly).
  - `POST /api/v1/auth/logout`: Logout dan pembatalan cookie (Memerlukan auth).
  - `GET /api/v1/auth/me`: Informasi profil user saat ini (Memerlukan auth).
- **Manajemen Pengguna (`/api/v1/users`):**
  - `GET /api/v1/users`: Daftar pengguna dengan paginasi & filter (Admin).
  - `POST /api/v1/users`: Pendaftaran pengguna baru (Admin).
  - `GET /api/v1/users/:id`: Detail profil pengguna (Admin atau diri sendiri).
  - `PUT /api/v1/users/:id`: Perbarui profil pengguna (Admin atau diri sendiri).
  - `PATCH /api/v1/users/:id/status`: Perbarui status ACTIVE/INACTIVE (Admin).
- **Sesi Praktikum (`/api/v1/sessions`):**
  - `GET /api/v1/sessions`: Daftar sesi praktikum dengan filter & paginasi (Admin & Mahasiswa).
  - `GET /api/v1/sessions/active`: Daftar sesi yang sedang aktif (Admin & Mahasiswa).
  - `GET /api/v1/sessions/:id`: Detail sesi praktikum (Admin & Mahasiswa).
  - `POST /api/v1/sessions`: Buat sesi praktikum baru (Admin).
  - `PUT /api/v1/sessions/:id`: Perbarui jadwal atau informasi sesi (Admin).
  - `PATCH /api/v1/sessions/:id/status`: Perbarui status siklus sesi SCHEDULED/ACTIVE/CLOSED/CANCELLED (Admin).
  - `DELETE /api/v1/sessions/:id`: Hapus sesi praktikum (Admin; ditolak jika memiliki data presensi).
- **Biometrik Wajah (`/api/v1/face`):**
  - `POST /api/v1/face/enroll`: Pendaftaran 3 pose wajah multi-angle (Mahasiswa).
  - `GET /api/v1/face/status`: Cek status pendaftaran wajah pengguna (Mahasiswa/Admin).
  - `DELETE /api/v1/face/template/:user_id`: Cabut status pendaftaran wajah pengguna (Admin).
- **Presensi Praktikum (`/api/v1/attendance`):**
  - `POST /api/v1/attendance/check-in`: Presensi masuk 1:1 server-side verification (Mahasiswa).
  - `POST /api/v1/attendance/check-out`: Presensi keluar stempel waktu server (Mahasiswa).
  - `GET /api/v1/attendance/history`: Riwayat presensi pribadi (Mahasiswa).
  - `GET /api/v1/attendance`: Monitoring daftar presensi semua mahasiswa (Admin).
  - `GET /api/v1/attendance/:id`: Detail presensi (Admin atau mahasiswa bersangkutan).
- **Master Aset Laboratorium (`/api/v1/assets`):**
  - `GET /api/v1/assets`: Daftar katalog master aset dengan filter status, kategori, kata kunci (Admin & Mahasiswa).
  - `GET /api/v1/assets/:id`: Detail master aset (Admin & Mahasiswa).
  - `GET /api/v1/assets/:id/history`: Histori dan timeline siklus peminjaman/kondisi aset (Admin & Mahasiswa ter-scope).
  - `GET /api/v1/assets/my-allocations`: Daftar aset yang sedang dipinjam mahasiswa (Mahasiswa).
  - `POST /api/v1/assets`: Tambah master aset baru (Admin).
  - `PUT /api/v1/assets/:id`: Perbarui data master aset (Admin).
  - `PATCH /api/v1/assets/:id/status`: Perbarui status siklus AVAILABLE/MAINTENANCE/RETIRED (Admin).
  - `DELETE /api/v1/assets/:id`: Hapus master aset (Admin; ditolak jika aset sedang dialokasikan/dalam perbaikan).
- **Alokasi & Peminjaman Aset (`/api/v1/allocations`):**
  - `POST /api/v1/allocations`: Alokasikan aset ke mahasiswa pada sesi aktif dengan validasi presensi masuk (Admin).
  - `POST /api/v1/allocations/:id/return`: Pengembalian aset dengan kondisi GOOD/DAMAGED (Admin).
  - `GET /api/v1/allocations`: Daftar seluruh peminjaman aset (Admin).
  - `GET /api/v1/allocations/my`: Daftar peminjaman aset pribadi (Mahasiswa).
  - `GET /api/v1/allocations/:id`: Detail data peminjaman aset (Admin atau mahasiswa bersangkutan).
- **Pemeliharaan Aset (`/api/v1/maintenance`):**
  - `POST /api/v1/maintenance`: Buat tiket pemeliharaan aset (Admin).
  - `GET /api/v1/maintenance`: Daftar tiket pemeliharaan aset (Admin).
  - `GET /api/v1/maintenance/:id`: Detail tiket pemeliharaan (Admin).
  - `PUT /api/v1/maintenance/:id`: Perbarui status, resolusi, atau biaya pemeliharaan (Admin).
- **Audit Trail (`/api/v1/audit-logs`):**
  - `GET /api/v1/audit-logs`: Inspeksi rekaman audit keamanan dan operasional append-only (Admin).

## Menjalankan Pengujian
```bash
go test -v ./...
```
