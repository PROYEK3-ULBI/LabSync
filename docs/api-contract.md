# REST API Contract (`/api/v1`)

## 1. Auth & Users
- **`POST /api/v1/auth/login` [PRD]**
  - **Role:** Public
  - **Req:** `{ "identifier": "TI240001", "password": "..." }`
  - **Res:** `200 OK`, `Set-Cookie: token=...; HttpOnly` | `{ "message": "Login berhasil" }`
  - **Error:** `401 Unauthorized` (Salah password), `403 Forbidden` (Nonaktif).
- **`POST /api/v1/auth/logout` [PRD]**
  - **Role:** All
  - **Res:** `200 OK`, `Set-Cookie: token=; Max-Age=0`
- **`GET /api/v1/auth/me` [PROPOSED]**
  - **Role:** All
  - **Res:** `200 OK` | `{ "id": "...", "name": "...", "role": "STUDENT", "face_enrolled": true }`
  - *Alasan:* Frontend (React) butuh state auth dan role dari cookie yang tidak bisa diakses JS.

## 2. Face Enrollment & Verification
- **`POST /api/v1/face/enroll` [PRD]**
  - **Role:** STUDENT
  - **Req (Multipart/Form-Data):** `front_image` (file), `left_image` (file), `right_image` (file), `consent_agreed` (boolean), `re_enroll` (boolean, opsional).
  - **Proses:** Go memvalidasi ownership user login & consent -> memvalidasi tipe file & batas ukuran -> memanggil internal Face Service untuk ekstraksi pose FRONT, LEFT, RIGHT -> mengenkripsi embedding dengan AES-256-GCM -> menyimpan `face_templates` dan 3 `face_template_samples` dalam satu transaksi atomik -> membuang raw image & embedding plaintext dari RAM.
  - **Res:** `201 Created` | `{ "message": "Pendaftaran wajah berhasil", "template_id": "...", "model_name": "ArcFace", "model_version": "1.0.0", "sample_count": 3, "enrolled_at": "..." }`
  - **Error:** `400 Bad Request` (Kurang pose, consent false), `403 Forbidden` (Bukan mahasiswa / akun nonaktif), `409 Conflict` (Sudah ada enrollment aktif dan re_enroll tidak aktif), `422 Unprocessable Entity` (NO_FACE_DETECTED, MULTIPLE_FACES_DETECTED), `503 Service Unavailable` (Face service down), `504 Gateway Timeout`.
- **`GET /api/v1/face/status` [PRD/PROPOSED]**
  - **Role:** STUDENT, ADMIN
  - **Req Query:** `?user_id=...` (opsional, hanya untuk Admin)
  - **Res:** `200 OK` | `{ "enrolled": true, "status": "ACTIVE", "model_name": "ArcFace", "model_version": "1.0.0", "consent_version": "v1.0", "enrolled_at": "...", "sample_count": 3, "poses": ["FRONT", "LEFT", "RIGHT"] }`
- **`DELETE /api/v1/face/template/:user_id` [PRD]**
  - **Role:** ADMIN
  - **Res:** `200 OK` | `{ "message": "Template biometrik wajah berhasil dicabut" }`
  - **Error:** `403 Forbidden`, `404 Not Found`.

## 3. Session & Attendance
### Definisi Session ACTIVE
Sesi berstatus `ACTIVE` ditentukan oleh status lifecycle database (`status = 'ACTIVE'`). Admin/instruktur mengontrol pembukaan sesi (`SCHEDULED` -> `ACTIVE`) dan penutupan sesi (`ACTIVE` -> `CLOSED` atau `CANCELLED`). Sesi yang telah memiliki referensi presensi tidak boleh dihapus secara fisik (`DELETE` -> `409 Conflict`), melainkan diubah statusnya menjadi `CLOSED` atau `CANCELLED`.
Waktu disimpan dan divalidasi dengan semantik UTC (`RFC3339`), dan tanggal disimpan dengan format ISO `YYYY-MM-DD`.

- **`GET /api/v1/sessions` [PRD]**
  - **Role:** STUDENT, ADMIN
  - **Req Query:** `?page=1&limit=10&status=ACTIVE&session_date=2026-09-18&search=PWB`
  - **Res:** `200 OK` | Daftar sesi praktikum + pagination metadata.
- **`GET /api/v1/sessions/active` [PROPOSED]**
  - **Role:** STUDENT, ADMIN
  - **Res:** `200 OK` | Daftar sesi yang sedang `ACTIVE`.
- **`GET /api/v1/sessions/:id` [PRD]**
  - **Role:** STUDENT, ADMIN
  - **Res:** `200 OK` | Detail sesi praktikum.
  - **Error:** `404 Not Found`.
- **`POST /api/v1/sessions` [PRD]**
  - **Role:** ADMIN
  - **Req:** `{ "code": "PWB-07", "title": "Praktikum Pemrograman Web", "course_code": "IF2101", "location": "Lab Komputer 1", "session_date": "2026-09-18", "start_at": "2026-09-18T13:00:00Z", "end_at": "2026-09-18T15:30:00Z", "status": "SCHEDULED" }`
  - **Res:** `201 Created`
  - **Error:** `400 Bad Request` (Waktu end_at <= start_at atau format salah), `409 Conflict` (Kode sesi duplikat).
- **`PUT /api/v1/sessions/:id` [PRD]**
  - **Role:** ADMIN
  - **Req:** `{ "title": "...", "location": "...", "start_at": "...", "end_at": "..." }`
  - **Res:** `200 OK`
  - **Error:** `400 Bad Request`, `404 Not Found`, `409 Conflict`.
- **`PATCH /api/v1/sessions/:id/status` [PRD]**
  - **Role:** ADMIN
  - **Req:** `{ "status": "ACTIVE" | "CLOSED" | "CANCELLED" | "SCHEDULED" }`
  - **Res:** `200 OK`
  - **Error:** `400 Bad Request`, `404 Not Found`.
- **`DELETE /api/v1/sessions/:id` [PRD/Design Rule]**
  - **Role:** ADMIN
  - **Res:** `200 OK`
  - **Error:** `404 Not Found`, `409 Conflict` (Jika sesi sudah memiliki data attendance).

- **`POST /api/v1/attendance/check-in` [PRD]**
  - **Role:** STUDENT
  - **Req (Multipart/Form-Data):** `session_id` (string), `image` (live frame binary).
  - **Proses:** Backend memvalidasi sesi aktif -> mengambil 3 embedding terenkripsi milik user login -> memanggil Python service untuk evaluasi 1:1 -> jika verifikasi lolos, insert presensi dengan stempel waktu server -> buang citra & plaintext embedding dari memori.
  - **Res:** `200 OK` | `{ "message": "Presensi masuk (check-in) berhasil", "attendance": { "id": "...", "session_id": "...", "user_id": "...", "checked_in_at": "...", "verification_status": "VERIFIED_FACE" } }`
  - **Error:** `400 Bad Request` (Sesi belum aktif / `SESSION_NOT_ACTIVE`, verifikasi biometrik gagal / `VERIFICATION_FAILED`, belum enrollment / `FACE_NOT_ENROLLED`), `403 Forbidden` (Bukan mahasiswa / akun nonaktif), `404 Not Found` (Sesi tidak ditemukan), `409 Conflict` (Sudah presensi pada sesi ini / `ALREADY_CHECKED_IN`), `422 Unprocessable Entity` (Wajah tidak ada / multi-wajah), `503 Service Unavailable`.
- **`POST /api/v1/attendance/check-out` [PRD]**
  - **Role:** STUDENT
  - **Req JSON:** `{ "session_id": "..." }`
  - **Res:** `200 OK` | `{ "message": "Presensi keluar (check-out) berhasil", "attendance": { "id": "...", "checked_out_at": "..." } }`
  - **Error:** `404 Not Found` (Belum check-in / `ATTENDANCE_NOT_FOUND`), `409 Conflict` (Sudah check-out sebelumnya / `ALREADY_CHECKED_OUT`).
- **`GET /api/v1/attendance/history` [PRD]**
  - **Role:** STUDENT
  - **Req Query:** `?page=1&limit=10`
  - **Res:** `200 OK` | Riwayat presensi milik mahasiswa yang sedang login.
- **`GET /api/v1/attendance` [PRD]**
  - **Role:** ADMIN
  - **Req Query:** `?page=1&limit=10&session_id=...&user_id=...&date=...`
  - **Res:** `200 OK` | Daftar presensi laboratorium seluruh mahasiswa.
- **`GET /api/v1/attendance/:id` [PRD]**
  - **Role:** ADMIN, STUDENT (Hanya miliknya)
  - **Res:** `200 OK` | Detail presensi.
  - **Error:** `403 Forbidden`, `404 Not Found`.

## 4. Master Assets & Allocations
- **`GET /api/v1/assets` [PRD]**
  - **Role:** ADMIN, STUDENT (Katalog aset dapat dibaca mahasiswa)
  - **Req Query:** `?page=1&limit=10&status=AVAILABLE&category=Laptop&search=Dell`
  - **Res:** `200 OK` | Daftar master aset + metadata paginasi.
- **`GET /api/v1/assets/:id` [PRD]**
  - **Role:** ADMIN, STUDENT
  - **Res:** `200 OK` | Detail data master aset.
  - **Error:** `404 Not Found`.
- **`POST /api/v1/assets` [PRD]**
  - **Role:** ADMIN
  - **Req:** `{ "asset_code": "LAB-LT-007", "asset_name": "Dell Latitude 5490", "category": "Laptop", "brand": "Dell", "model": "Latitude 5490", "serial_number": "SN-DELL-007", "location": "Lab Komputer 1", "status": "AVAILABLE" }`
  - **Res:** `201 Created`
  - **Error:** `400 Bad Request`, `403 Forbidden`, `409 Conflict` (Kode aset atau nomor seri duplikat).
- **`PUT /api/v1/assets/:id` [PRD]**
  - **Role:** ADMIN
  - **Req:** `{ "asset_name": "...", "location": "...", "notes": "..." }`
  - **Res:** `200 OK`
  - **Error:** `400 Bad Request`, `403 Forbidden`, `404 Not Found`, `409 Conflict` (Aset sedang dipinjam atau dalam perbaikan).
- **`PATCH /api/v1/assets/:id/status` [PRD]**
  - **Role:** ADMIN
  - **Req:** `{ "status": "AVAILABLE" | "MAINTENANCE" | "RETIRED" }`
  - **Res:** `200 OK`
  - **Error:** `400 Bad Request` (Transisi status tidak diizinkan), `403 Forbidden`, `404 Not Found`, `409 Conflict` (Aset sedang dipinjam atau dalam perbaikan aktif).
- **`DELETE /api/v1/assets/:id` [PRD]**
  - **Role:** ADMIN
  - **Res:** `200 OK`
  - **Error:** `403 Forbidden`, `404 Not Found`, `409 Conflict` (Aset sedang dipinjam atau dalam perbaikan).
- **`POST /api/v1/allocations` [PRD]**
  - **Role:** ADMIN
  - **Req:** `{ "asset_id": "...", "user_id": "...", "session_id": "..." }`
  - **Res:** `201 Created`
  - **Error:** `409 Conflict` (Aset tidak AVAILABLE).
- **`POST /api/v1/allocations/:id/return` [PRD]**
  - **Role:** ADMIN
  - **Req:** `{ "condition": "GOOD" | "DAMAGED" }`
  - **Res:** `200 OK`
- **`GET /api/v1/assets/my-allocations` [PRD]**
  - **Role:** STUDENT
  - **Res:** `200 OK` | Aset yang saat ini sedang dipinjam mahasiswa.

## 5. Maintenance & Audit Trail (Admin Only)
- **`POST /api/v1/maintenance` [PRD]**
  - **Role:** ADMIN
  - **Req:** `{ "asset_id": "...", "issue_type": "HARDWARE", "description": "...", "priority": "HIGH" | "MEDIUM" | "LOW" }`
  - **Res:** `201 Created`
  - **Error:** `400 Bad Request`, `403 Forbidden`, `404 Not Found`, `409 Conflict` (Aset sedang dipinjam).
- **`GET /api/v1/maintenance` [PRD]**
  - **Role:** ADMIN
  - **Req Query:** `?page=1&limit=10&asset_id=...&status=OPEN|IN_PROGRESS|COMPLETED|CANCELLED&priority=HIGH`
  - **Res:** `200 OK` | Daftar tiket pemeliharaan aset + pagination.
- **`GET /api/v1/maintenance/:id` [PRD]**
  - **Role:** ADMIN
  - **Res:** `200 OK` | Detail tiket pemeliharaan.
- **`PUT /api/v1/maintenance/:id` [PRD]**
  - **Role:** ADMIN
  - **Req:** `{ "status": "IN_PROGRESS" | "COMPLETED" | "CANCELLED", "resolution": "...", "cost": 150000, "result_status": "AVAILABLE" | "RETIRED" }`
  - **Res:** `200 OK`
  - **Error:** `400 Bad Request` (Transisi status ilegal atau resolusi belum diisi), `403 Forbidden`, `404 Not Found`.
- **`GET /api/v1/audit-logs` [PRD/PROPOSED]**
  - **Role:** ADMIN
  - **Req Query:** `?page=1&limit=10&action=...&entity_type=...&user_id=...`
  - **Res:** `200 OK` | Log audit operasional dan keamanan append-only (tidak ada endpoint edit/delete).

## 6. Reports & System
- **`GET /api/v1/reports/attendance` [PRD]**
  - **Role:** ADMIN
  - **Req Query:** `?session_id=...&start_date=...`
  - **Res:** `200 OK` | Agregasi data laporan.
- **`GET /api/v1/health` [PROPOSED]**
  - **Role:** Public
  - **Res:** `200 OK` | Liveness probe Go dan koneksi DB.
- **`GET /api/v1/docs` [PRD]**
  - **Role:** Public
  - **Res:** Swagger/OpenAPI UI.