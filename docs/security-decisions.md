# Security & Architecture Decisions

## 1. Batas Tanggung Jawab (Boundaries)
- **Go + Fiber (Main API):** Entry point HTTP utama. Menangani routing, autentikasi, otorisasi (RBAC), validasi input, transaksi database, dan orkestrasi ke face service.
- **Python + FastAPI (Face Service):** Microservice internal eksklusif untuk ekstraksi dan verifikasi biometrik. Menerima gambar di memori, mengekstrak embedding via DeepFace, menghitung distance, dan mengembalikan hasil (boolean `verified` atau array of floats untuk embedding).
- **Supabase PostgreSQL:** Sumber data utama (Relational DB). Tidak diakses langsung dari Frontend.
- **Supabase Storage:** **TIDAK DIPAKAI** untuk raw image pada MVP sesuai PRD.
- **Frontend (React/Vite):** Layer presentasi murni. Tidak menyimpan kredensial database. Berkomunikasi hanya via HTTP/HTTPS ke Go API menggunakan cookie credentials.

## 2. Struktur Folder Backend (Usulan)
Mengikuti standar idiomatic Go dan simpel:
```text
/backend
├── cmd/server/main.go
├── internal/
│   ├── config/      # Load env, DB config, Python service URL
│   ├── handler/     # HTTP controllers (Fiber)
│   ├── middleware/  # Auth, Error Handling, Logger
│   ├── model/       # DTOs, DB structs (GORM/sqlc)
│   ├── repository/  # Query & DB Transactions
│   ├── router/      # API Routes mapping
│   └── service/     # Business logic, panggil internal repo & Python API
└── pkg/             # Utility (Encryption, dll)

/face_service
├── main.py
├── requirements.txt
├── routes/
└── services/
```

## 3. Autentikasi & Cookie
- **Mekanisme:** JWT yang disimpan di dalam **HttpOnly, Secure, SameSite=Lax Cookie**.
- **Aturan Log:** Token tidak pernah dicatat di server log.
- **Exposure:** JavaScript frontend tidak dapat mengakses JWT (melindungi dari XSS). Frontend bergantung pada endpoint `GET /api/v1/auth/me` untuk mengetahui status autentikasi dan role (Admin/Student).
- **Password:** Hashed menggunakan bcrypt dengan salt otomatis sebelum masuk DB.
- **Resolusi Ambiguitas PRD (Token Response):** Tabel PRD 11.1 menyebutkan kolom response `token/session + user`. Ditetapkan bahwa token JWT **TIDAK PERNAH** disertakan di dalam JSON response body `POST /api/v1/auth/login` demi mencegah eksfiltrasi token dari memori/storage JS. Token ditransmisikan eksklusif via `Set-Cookie` HttpOnly. Frontend tidak menyimpan token di `localStorage`/`sessionStorage`. Status auth diverifikasi via `GET /api/v1/auth/me`. Invalidation logout dilakukan via `Set-Cookie: token=; Max-Age=0`.

## 4. Komunikasi Internal Go -> Python (Aman)
- Browser/Mahasiswa **TIDAK PERNAH** memanggil Python API secara langsung.
- Go API meneruskan request file (di memori) ke Python melalui internal network (misal: `http://localhost:8001` atau `http://face-service:8001` dalam Docker).
- **Keamanan Internal:** Endpoint Python dilindungi oleh header `X-Internal-API-Key` yang dicocokkan dengan `INTERNAL_API_KEY` dari environment variable.
- **Zero-Retention Citra:** Python service memproses data citra murni di memori dan membuang byte citra maupun array numpy segera setelah ekstraksi/inferensi selesai. Tidak ada penyimpanan ke disk, folder lokal, maupun Supabase Storage.
- **Logging Bersih:** Logging terstruktur mencatat method, path, status, latency, dan request ID tanpa pernah mencatat token rahasia, binary citra, atau nilai array embedding biometrik.
- **Aturan Multi-Angle & Kalibrasi:** Verifikasi 1:1 terhadap 3 sampel pendaftaran (`FRONT`, `LEFT`, `RIGHT`) menggunakan strategi `MIN_DISTANCE` berbasis Cosine distance. Evaluasi empiris (ROC / False Acceptance Rate vs False Rejection Rate) wajib diselesaikan sebelum `IS_CALIBRATED=true` diaktifkan di tingkat produksi.

## 5. Standar Error & Pagination
- **Error Format:** Semua HTTP error mengembalikan format JSON konsisten tanpa membocorkan internal database (e.g. stack trace).
  ```json
  {
    "error": {
      "code": "UNAUTHORIZED",
      "message": "Kredensial tidak valid",
      "request_id": "req-12345"
    }
  }
  ```
- **Pagination [PROPOSED]:** Endpoint yang mengembalikan daftar (seperti `/users`, `/sessions`, `/assets`, `/attendance/history`) menggunakan parameter query `?page=1&limit=10`. Response mencakup metadata: `{"data": [...], "meta": {"total": 50, "page": 1, "limit": 10}}`.

## 6. Keputusan Desain Modul Lab Sessions
- **Definisi Sesi ACTIVE:** Ditentukan secara definitif oleh lifecycle status `status = 'ACTIVE'` pada database. Status ini dikontrol oleh Admin/Dosen untuk membuka sesi presensi secara operasional. Hal ini mengakomodasi keterlambatan atau pergeseran menit jadwal riil di laboratorium dibandingkan asumsi kaku jadwal statis.
- **Semantik Timezone:** Semua stempel waktu `start_at` dan `end_at` divalidasi dan dinormalisasi menjadi UTC (`RFC3339`) sebelum disimpan ke kolom `TIMESTAMP WITH TIME ZONE`. Tanggal kalender sesi menggunakan tipe ISO 8601 `DATE` (`YYYY-MM-DD`) murni untuk mencegah pergeseran hari akibat perbedaan offset lokal klien.
- **Larangan Hapus Fisik Sesi Berpresensi:** Sesi yang telah memiliki keterkaitan relasional dengan data presensi mahasiswa (`attendance`) tidak boleh dihapus secara fisik (`DELETE` ditolak dengan `409 Conflict`). Admin diarahkan untuk mengubah status sesi menjadi `CLOSED` atau `CANCELLED`.

## 7. Enkripsi Biometrik & Key Management (Face Enrollment)
- **Skema Enkripsi:** Menggunakan algoritma terotentikasi **AES-256-GCM** (`crypto/aes`, `crypto/cipher`). Setiap embedding (array 512 float) diserialisasi, dienkripsi menggunakan nonce acak 12-byte, dan disimpan pada kolom `embedding_ciphertext` tipe `BYTEA`.
- **Manajemen Kunci (Key Management):** Kunci rahasia 256-bit (32 bytes) dimuat eksklusif dari environment variable `BIOMETRIC_ENCRYPTION_KEY`. Kunci **TIDAK PERNAH** disimpan di dalam tabel database, script migrasi, atau source code repositori.
- **Kebijakan Zero-Retention:** Citra mentah mahasiswa tidak pernah disimpan ke disk lokal, Supabase Storage, atau database. Data citra hanya berada di buffer RAM selama proses validasi dan transfer internal ke microservice Python, lalu langsung dibersihkan.
- **Transaksional Atomik 3-Pose:** Pendaftaran wajah wajib memuat tepat 3 pose (`FRONT`, `LEFT`, `RIGHT`). Template induk (`face_templates`) dan ketiga sampel pose (`face_template_samples`) disimpan dalam **satu transaksi database atomik**. Jika salah satu pose gagal (tidak ada wajah, buram, multi-wajah), seluruh transaksi dibatalkan (`ROLLBACK`) dan tidak ada template parsial yang tersimpan.
- **Aturan Pendaftaran Ulang (Re-enrollment):** Pengguna dengan template `ACTIVE` tidak dapat mengganti template secara diam-diam (`POST /api/v1/face/enroll` ditolak dengan `409 Conflict: FACE_ALREADY_ENROLLED`). Penggantian hanya diizinkan dengan parameter `re_enroll=true`, di mana template lama baru diubah statusnya menjadi `REVOKED` setelah ketiga sampel baru terverifikasi valid dan transaksi berhasil di-commit.
- **Audit Tanpa Kebocoran Data:** Aktivitas `FACE_ENROLLED`, `FACE_RE_ENROLLED`, dan `FACE_TEMPLATE_REVOKED` dicatat di `audit_logs` tanpa mencatat nilai array embedding biometrik maupun binary citra.

## 8. Desain Alur Presensi Biometrik & Server-Side Verification
- **Keputusan Otoritatif Server:** Frontend **TIDAK PERNAH** mengirim status boolean `verified: true`. Backend Go mendekripsi sampel template milik user login dan memverifikasi live frame secara server-side via microservice Python. Keputusan keabsahan presensi murni ditentukan oleh server.
- **Stempel Waktu Server (Server-Authoritative Timestamps):** Nilai `checked_in_at` dan `checked_out_at` selalu ditentukan oleh server menggunakan UTC timestamp (`time.Now().UTC()`), mengabaikan manipulasi jam/zona waktu di perangkat klien.
- **Pencegahan Double Check-in & Race Condition:** Constraint database `idx_attendance_user_session` (UNIQUE `user_id, session_id`) menjamin satu mahasiswa hanya dapat memiliki satu baris presensi per sesi. Jika terjadi double-submit konkuren, error PostgreSQL 23505 ditangkap dan ditranslasikan menjadi respons `409 Conflict: ALREADY_CHECKED_IN`.
- **Diferensiasi Error Presensi:** Kegagalan verifikasi biometrik mengembalikan error `400 Bad Request` dengan kode `VERIFICATION_FAILED`. Error ini dibedakan secara tegas dari kesalahan input citra (`NO_FACE_DETECTED`, `MULTIPLE_FACES_DETECTED`, `INVALID_IMAGE`), sesi belum aktif (`SESSION_NOT_ACTIVE`), atau kegagalan koneksi internal (`FACE_SERVICE_UNAVAILABLE`).
- **Prosedur Kalibrasi Threshold Biometrik:**
  1. Pengujian empiris dilakukan dengan mengumpulkan citra wajah pasangan orang sama (*genuine pairs*) dan orang berbeda (*impostor pairs*) di bawah variasi pencahayaan laboratorium komputer universitas.
  2. Dihitung kurva ROC (*Receiver Operating Characteristic*) untuk metrik Cosine Distance ArcFace.
  3. Ditetapkan nilai threshold di mana *False Acceptance Rate* (FAR) mendekati 0.001% dengan *False Rejection Rate* (FRR) di bawah rentang yang dapat diterima pengguna (2-5%).
  4. Nilai threshold hasil kalibrasi diset pada `DEFAULT_SIMILARITY_THRESHOLD` dan `IS_CALIBRATED=true` di environment Python service.
  5. Selama kalibrasi belum selesai (`IS_CALIBRATED=false`), microservice Python secara otomatis menandai `verified=false` untuk mencegah lolosnya presensi yang tidak terkalibrasi di lingkungan produksi.

## 9. Aturan Bisnis & Transaksional Alokasi/Pengembalian Aset
- **Transaksi Atomik Alokasi:** Pembuatan alokasi baru dan perubahan status aset dari `AVAILABLE` menjadi `ALLOCATED` diproses dalam **satu transaksi database**. Baris aset dikunci dengan `SELECT ... FOR UPDATE` untuk mencegah race condition. Jika terjadi kegagalan atau rollback, status aset tetap `AVAILABLE`.
- **Pencegahan Alokasi Ganda (Double Allocation):** Dilindungi secara berlapis melalui partial unique index PostgreSQL `idx_asset_allocations_active` (`UNIQUE (asset_id) WHERE returned_at IS NULL`) serta pengecekan aplikasi level database. Pelanggaran menghasilkan `409 Conflict: ASSET_ALREADY_ALLOCATED`.
- **Prasyarat Kehadiran (Attendance Prerequisite):** Mahasiswa wajib telah menyelesaikan presensi masuk (`checked_in_at IS NOT NULL` dan `checked_out_at IS NULL`) pada sesi praktikum yang bersangkutan sebelum unit aset dapat dialokasikan kepadanya (`400 Bad Request: ATTENDANCE_REQUIRED`).
- **Logika Pengembalian (Return Flow):** Stempel waktu `returned_at` dicatat secara otoritatif dari waktu server UTC (`time.Now().UTC()`). Kondisi fisik setelah pengembalian menentukan status siklus aset:
  - `GOOD` -> Status aset kembali menjadi `AVAILABLE`.
  - `DAMAGED` -> Status aset otomatis beralih menjadi `MAINTENANCE`.
  Pengembalian berulang atas alokasi yang telah selesai ditolak dengan `409 Conflict: ALREADY_RETURNED`.
- **Pembatasan Scoping & Histori Permanen:** Mahasiswa hanya diizinkan membaca daftar alokasi dan riwayat timeline miliknya (`GET /api/v1/allocations/my` dan `GET /api/v1/assets/my-allocations`). Seluruh rekaman peminjaman dan pengembalian disimpan permanen (tidak pernah dihapus fisik) guna keperluan audit operasional laboratorium.

## 10. Desain Pemeliharaan (Maintenance) & Audit Trail Append-Only
- **Mesin Keadaan (State Machine) Pemeliharaan:**
  - Status tiket: `OPEN` -> `IN_PROGRESS` -> `COMPLETED` atau `CANCELLED`.
  - `COMPLETED` dan `CANCELLED` merupakan *terminal state* yang tidak dapat diubah lagi.
  - Transisi ke `IN_PROGRESS` menetapkan stempel waktu `started_at = NOW()`.
  - Transisi ke `COMPLETED` menetapkan `completed_at = NOW()`, mewajibkan teks `resolution`, dan memperbarui status aset menjadi `AVAILABLE` (atau `RETIRED` jika perbaikan tidak memungkinkan).
  - Transisi ke `CANCELLED` mengembalikan status aset ke `AVAILABLE`.
- **Konsistensi Invariant terhadap Alokasi:**
  - Aset berstatus `MAINTENANCE` tidak dapat dialokasikan (`409 Conflict: ASSET_NOT_AVAILABLE`).
  - Aset yang sedang dipinjam (`ALLOCATED` / alokasi aktif) tidak dapat dimasukkan ke pemeliharaan sebelum pengembalian aset diselesaikan (`409 Conflict: ASSET_ALLOCATED`).
- **Penyimpanan Histori & Keputusan Skema Tabel:**
  - Sesuai DDL 10 tabel resmi, histori tiket perbaikan dikelola langsung pada tabel `maintenance` yang berelasi dengan `assets`.
  - Seluruh peristiwa perubahan status (`MAINTENANCE_CREATED`, `MAINTENANCE_STATUS_UPDATED`, `MAINTENANCE_COMPLETED`, `MAINTENANCE_CANCELLED`, `ALLOCATION_CREATED`, `ALLOCATION_RETURNED`, `ASSET_CREATED`, dll.) dicatat ke tabel `audit_logs`.
  - Histori aset (`GET /api/v1/assets/:id/history`) mengagregasi peristiwa alokasi, pengembalian, dan pemeliharaan secara dinamis tanpa memerlukan tabel baru yang tidak disetujui.
  - Audit log bersifat **append-only murni** dari perspektif API: tidak ada endpoint untuk mengedit atau menghapus rekaman log audit.
  - Metadata audit disaring ketat: dilarang keras memuat password, cookie/token JWT, binary citra wajah, atau vektor embedding biometrik.

## 11. Kontrak Frontend
- **API Client:** Menggunakan wrapper fetch/Axios (`src/lib/api.ts`).
- **Config:** Base URL dari `import.meta.env.VITE_API_URL`.
- **Credentials:** Request harus menyertakan `withCredentials: true` (Axios) atau `credentials: 'include'` (Fetch) agar HttpOnly cookie terkirim.
- **State Handling:** Menangani status `loading` (spinner) dan `error` (toast UI) tanpa merusak struktur halaman yang sudah ada. Menggantikan pemanggilan fungsi mock dengan async call.

## 12. Prioritas Implementasi & Acceptance Criteria
1. **Fase 1: Skeleton & Database** (AC: Go & FastAPI jalan, DB migrate sukses).
2. **Fase 2: Auth & RBAC** (AC: Login sukses dapat HttpOnly Cookie, non-admin tertolak dari API admin).
3. **Fase 3: Face Microservice** (AC: POST ke Python return embedding, raw image hilang).
4. **Fase 4: Enrollment API** (AC: Mahasiswa POST 3 gambar -> DB simpan 3 encrypted embedding).
5. **Fase 5: Session & Verification/Check-in** (AC: Verifikasi wajah sukses -> DB attendance tersimpan. Wajah salah -> Ditolak).
6. **Fase 6: Asset & Maintenance** (AC: Admin bisa alokasi, double allocation ditolak DB constraint).
7. **Fase 7: Integrasi Frontend** (AC: UI berfungsi penuh dengan backend riil, data.ts dihapus).