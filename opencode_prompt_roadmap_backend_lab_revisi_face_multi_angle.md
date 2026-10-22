# Roadmap Prompt OpenCode — Backend Presensi & Manajemen Aset Laboratorium

Dokumen ini berisi prompt bertahap untuk mengembangkan backend dengan OpenCode berdasarkan:
- `prd_app_proyek3_revisi_face_enrollment_multi_angle.md`
- `diagram-proyek3_plantuml_revisi_face_multi_angle.md`
- kode frontend yang sudah ada di repository

## Cara memakai

1. Pastikan dua file dokumen ada di repository dan OpenCode sudah melakukan `/init`.
2. Jalankan prompt **satu per satu sesuai urutan**. Jangan tempel semua prompt sekaligus.
3. Setelah setiap prompt, review ringkasan, file yang berubah, hasil test, dan `git diff`. Jika fase belum lulus kriteria selesai, perbaiki fase tersebut dulu.
4. Untuk Prompt 00 dan 01, minta OpenCode bekerja sebagai analis/perencana dan jangan mengubah kode aplikasi.
5. Simpan secret hanya dalam file `.env` lokal yang tidak di-commit. File `.env.example` hanya berisi placeholder.
6. Jangan lanjut ke fase berikutnya bila prompt sebelumnya melaporkan blocker arsitektur atau acceptance criteria yang belum tercapai.

## Keputusan arsitektur dari PRD yang wajib dijaga

- Frontend yang ada dipertahankan: jangan redesign, rewrite, atau mengganti komponen secara massal.
- Frontend berkomunikasi hanya dengan REST API; tidak boleh mengakses PostgreSQL/Supabase secara langsung.
- Backend API utama: **Go + Fiber**.
- Face verification: **Python + FastAPI** sebagai service internal yang dipanggil oleh Go backend, tidak diakses langsung oleh browser.
- Database: **Supabase PostgreSQL** dengan migration SQL yang tersimpan di repository.
- Face stack: **DeepFace** dengan pretrained model; ArcFace adalah kandidat yang harus diuji, bukan asumsi threshold/akurasi final.
- Base path API: `/api/v1`.
- Role MVP: `STUDENT` dan `ADMIN`.
- Modul inti: auth, user, face enrollment/verification 1:1, lab sessions, attendance check-in/out, assets, allocation/return/history, maintenance, audit log, dashboard, reporting, OpenAPI/Swagger, testing, deployment.
- Biometric: enrollment wajib punya tiga pose `FRONT`, `LEFT`, `RIGHT`. Simpan 3 embedding terenkripsi pada `face_template_samples` di bawah satu enrollment set `face_templates`; raw image hanya diproses di memory lalu dibuang. Jangan log password, token, raw image, atau embedding plaintext.
- Jangan menggunakan Supabase Auth kecuali perubahan arsitektur disetujui secara eksplisit; PRD menempatkan autentikasi pada Go API dengan password hash dan cookie aman.
- Jangan menambahkan QR/RFID/IoT, face identification 1:N, training model dari nol, atau fitur di luar scope.

## Catatan desain biometrik yang wajib dijaga

PRD revisi sudah menetapkan bahwa raw face image tidak disimpan permanen dan Supabase Storage tidak digunakan untuk wajah pada MVP. Enrollment mengirim tiga pose sekaligus (`front_image`, `left_image`, `right_image`) ke Go API; Go memanggil Python face service, lalu menyimpan satu metadata set dan tiga embedding terenkripsi secara atomik. Saat verification, backend membandingkan live frame/burst hanya dengan sampel milik user yang login. Threshold dan strategi agregasi wajib dikalibrasi melalui pengujian.

Jika ada kontradiksi lain antara PRD revisi, diagram revisi, dan implementasi frontend, jangan diam-diam memilih salah satunya. Dokumentasikan lokasinya, dampaknya, keputusan sementara, dan apakah perlu keputusan manusia. Jangan berhenti untuk hal non-blocking: tulis asumsi yang jelas lalu lanjut sesuai scope.

---

# PROMPT 00 — Audit Repository dan Peta Sistem (tanpa implementasi)

```text
Bertindak sebagai senior software architect dan tech lead. Baca seluruh konteks repository, file AGENTS.md atau instruksi OpenCode yang ada, PRD `prd_app_proyek3_revisi_face_enrollment_multi_angle.md`, diagram `diagram-proyek3_plantuml_revisi_face_multi_angle.md`, dan implementasi frontend yang sudah ada. `/init` sudah pernah dijalankan; jangan menganggap hasilnya selalu lengkap, verifikasi struktur repo sendiri.

Untuk fase ini ANALISIS SAJA. Jangan edit source code aplikasi, jangan install dependency, jangan mengubah UI, jangan membuat koneksi ke layanan remote.

Inventarisasi:
1. Struktur folder, package manager, versi runtime, framework frontend, router, design system, form/state management, komponen halaman, mock data, dan semua titik yang saat ini membaca/menulis data.
2. Halaman/fitur frontend yang ada untuk login, dashboard, users, sessions, attendance, face enrollment/verification, assets, allocations, maintenance, reports, dan profile. Bedakan yang sungguh sudah ada dari yang belum ada.
3. Konfirmasi tidak ada backend/database yang sudah dibuat; bila ada, jangan menimpa—catat kondisinya.
4. Ekstrak aturan bisnis, tabel, endpoint, roles, acceptance criteria, security requirements, dan diagram sequence/state dari PRD/PlantUML.
5. Catat konflik atau ambiguitas dokumen, khususnya cookie/JWT, threshold/strategi multi-sample face verification, timestamp/timezone, persetujuan biometrik, dan detail field tabel/API yang belum konsisten.
6. Tandai fitur Must Have, Should Have, dan Nice to Have.

Buat hanya dokumen perencanaan:
- `docs/repo-audit.md`: kondisi repo apa adanya dengan path/file nyata.
- `docs/requirements-questions.md`: tabel isu, sumber konflik, dampak, rekomendasi/default sementara, dan apakah perlu keputusan manusia.
- `docs/implementation-plan.md`: urutan implementasi, dependency, definition of done per fase, serta file/folder target yang disarankan berdasarkan repo nyata.

Jangan mengarang file atau komponen yang tidak ditemukan. Jika lokasi dokumen berbeda, cari dulu berdasarkan nama. Akhiri dengan ringkasan kondisi repo, blocker, asumsi, dan daftar file yang dibuat. Jangan mulai Prompt 01 atau implementasi apa pun.
```

**Selesai jika:** audit punya path yang dapat diverifikasi; tidak ada source code aplikasi berubah; semua konflik penting dicatat.

# PROMPT 01 — Finalisasi Desain Backend dan Kontrak API (belum coding)

```text
Gunakan `docs/repo-audit.md`, `docs/requirements-questions.md`, `docs/implementation-plan.md`, PRD, diagram, dan frontend yang benar-benar ada.

Buat desain implementasi yang cukup spesifik untuk dikerjakan, tanpa mengimplementasikan source backend/frontend dan tanpa menjalankan perubahan remote di Supabase.

Tentukan dan dokumentasikan:
1. Batas tanggung jawab Go + Fiber, Python FastAPI face service, Supabase PostgreSQL, dan frontend. Supabase Storage tidak dipakai untuk menyimpan raw face image pada MVP.
2. Struktur folder backend mengikuti kebiasaan repo yang ada; bila belum ada, usulkan struktur sederhana yang tidak over-engineered.
3. Kontrak endpoint REST `/api/v1` yang mencakup auth, users, sessions, face enrollment/status/verification, attendance check-in/out/me/history, assets, allocations/returns/history, maintenance, dashboard, reporting, health, dan OpenAPI. Tandai endpoint berdasarkan bukti PRD atau sebagai usulan. Kontrak enrollment menerima tiga file pose (`front_image`, `left_image`, `right_image`) sekaligus; kontrak verification menerima live frame/burst dan hanya mencocokkan dengan sampel user login.
4. Untuk setiap endpoint: method, path, role, request/response DTO, status sukses, error codes, validasi, dan aturan akses data milik sendiri vs admin.
5. Model relasional sepuluh tabel: `roles`, `users`, `face_templates`, `face_template_samples`, `lab_sessions`, `attendance`, `assets`, `asset_allocations`, `maintenance`, `audit_logs`. `face_templates` hanya menyimpan metadata enrollment set; `face_template_samples` menyimpan embedding terenkripsi berlabel `FRONT`/`LEFT`/`RIGHT`. Sertakan constraint, partial unique index untuk satu set aktif per user, unique `(face_template_id, pose)`, status enum, serta operasi transaksi.
6. Format error konsisten yang tidak membocorkan internal detail, serta pagination/filter conventions bila daftar membutuhkan.
7. Alur login berbasis password hash dan cookie HttpOnly/Secure/SameSite yang tepat untuk environment; token tidak pernah dicatat di log dan tidak boleh diekspos ke frontend JavaScript.
8. Cara Go memanggil Python face service secara internal dan aman. Browser tidak pernah memanggil Python service secara langsung.
9. Kontrak frontend: API client, base URL dari environment, credentials/cookie behavior, tipe response, loading/error state; pertahankan tampilan dan struktur halaman yang ada.
10. Prioritas implementasi per fase dengan acceptance criteria yang bisa diuji.

Buat/update:
- `docs/api-contract.md`
- `docs/data-model.md`
- `docs/security-decisions.md`
- draft `docs/openapi.yaml` jika seluruh kontrak sudah cukup jelas; jika belum, tandai bagian belum diputuskan.

Jangan mulai coding. Jangan menciptakan endpoint hanya karena lazim; semua usulan harus diberi label usulan dan ditautkan ke kebutuhan frontend/PRD. Laporkan pertanyaan yang benar-benar blocking sebelum implementasi database dimulai.
```

**Selesai jika:** desain API dan model data cocok dengan kebutuhan frontend, aturan akses tertulis, serta blocker utama dipisahkan dari asumsi sementara.

# PROMPT 02 — Supabase PostgreSQL dan Migration SQL

```text
Implementasikan fase database saja berdasarkan `docs/data-model.md`, PRD, dan diagram. Ikuti keputusan tech stack dan instruksi repository yang sudah ada.

Kebutuhan:
1. Siapkan struktur Supabase CLI/migration bila cocok dengan repo, termasuk `supabase/migrations/` dan dokumentasi cara menjalankannya. Jangan mengasumsikan project Supabase remote sudah tersedia.
2. Buat migration SQL untuk sepuluh tabel: `roles`, `users`, `face_templates`, `face_template_samples`, `lab_sessions`, `attendance`, `assets`, `asset_allocations`, `maintenance`, `audit_logs`. `face_templates` menyimpan metadata set; jangan tempatkan satu `embedding_ciphertext` langsung pada tabel parent.
3. Tambahkan primary key UUID, foreign key, NOT NULL, unique/check constraints, index yang relevan, timestamp, dan aturan status sesuai PRD dan `docs/data-model.md`.
4. Terapkan constraint untuk satu attendance per user per session, unique identifier/email, unique asset code, dan satu allocation aktif per asset. Jika constraint parsial diperlukan di PostgreSQL, implementasikan dengan benar dan jelaskan.
5. Seed role `STUDENT` dan `ADMIN` secara idempotent. Jangan seed akun admin dengan password default atau secret.
6. Tambahkan migration rollback hanya bila pola migrasi yang dipilih mendukungnya dengan aman; jelaskan jika rollback membutuhkan migration baru.
7. Jangan expose credential PostgreSQL, Supabase service role, atau secret ke frontend. Jangan menambahkan Supabase client di aplikasi React.
8. Jangan membuat bucket foto wajah dan jangan simpan raw image. Simpan `consent_version`/`consented_at` pada metadata enrollment set, dan tepat tiga sampel embedding terenkripsi per set lengkap dengan pose FRONT/LEFT/RIGHT. Encryption key hanya di environment/secret manager; tidak boleh menyimpan plaintext embedding.
9. Jangan mengubah schema remote. Berikan command yang aman untuk menjalankan lokal, memeriksa migration, dan nanti menerapkan ke project yang benar.
10. Sertakan seed data development yang fiktif dan aman hanya jika diperlukan; jangan menggunakan data mahasiswa nyata.

Tambahkan/update dokumentasi schema dan `.env.example` hanya dengan nama variable serta placeholder, tanpa nilai rahasia. Jalankan pemeriksaan migration yang tersedia. Jangan lanjut membangun endpoint API pada fase ini. Akhiri dengan file yang berubah, command test yang dijalankan, dan langkah manual yang masih diperlukan.
```

**Selesai jika:** migration tersimpan di repo, constraint kritis tercakup, dapat dijalankan secara lokal/dry-run, dan tidak ada perubahan remote atau secret.

# PROMPT 03 — Bootstrap Go + Fiber API

```text
Bangun kerangka backend utama Go + Fiber sesuai `docs/implementation-plan.md`, `docs/api-contract.md`, dan pola repo yang sudah diaudit. Implementasi hanya fondasi, belum fitur domain lengkap.

Buat/implementasikan:
1. Entry point, konfigurasi environment tervalidasi, graceful shutdown, routing `/api/v1`, dan pemisahan config, handler/transport, service/business logic, repository/data access, middleware, DTO, dan error handling tanpa abstraksi berlebihan.
2. Koneksi PostgreSQL Supabase melalui connection string dari environment. Gunakan driver/query approach yang konsisten dengan keputusan desain. Jangan hardcode URL/password.
3. Endpoint health/readiness yang tidak membocorkan credential atau detail infrastruktur.
4. Middleware request ID, panic recovery, logging terstruktur dengan redaksi secret, CORS berdasarkan allowlist environment, batas ukuran request, dan rate limit dasar untuk route sensitif.
5. Validasi request dan error response konsisten sesuai kontrak API.
6. Dokumentasi cara run lokal dan file `.env.example` jika belum tersedia.
7. OpenAPI/Swagger awal yang selaras dengan route yang benar-benar ada.

Jangan implementasikan login, CRUD, face, atau fitur domain pada fase ini. Jangan mengubah frontend. Gunakan versi dependency yang cocok dengan repository; verifikasi kompatibilitas, jangan memasang versi yang tidak teruji. Jalankan `gofmt`, `go test ./...`, dan lint yang memang tersedia. Perbaiki error yang berasal dari perubahan ini sebelum selesai.
```

**Selesai jika:** server dapat start dengan konfigurasi development, health endpoint berfungsi, database connection failure ditangani aman, dan test/build fondasi lulus.

# PROMPT 04 — Authentication, Authorization, dan User Management

```text
Implementasikan authentication dan user management di Go + Fiber sesuai PRD dan `docs/api-contract.md`. Fokus hanya auth/user pada fase ini.

Scope:
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/me`
- `GET /api/v1/users`
- `POST /api/v1/users`
- `GET /api/v1/users/{id}`
- `PUT /api/v1/users/{id}`
- `PATCH /api/v1/users/{id}/status`

Persyaratan:
1. Password disimpan dengan password hashing modern yang sesuai untuk Go; tidak pernah plaintext dan tidak pernah dikembalikan oleh API.
2. Gunakan cookie HttpOnly, Secure pada HTTPS/production, SameSite yang sesuai, expiry dan invalidasi saat logout. Jangan menaruh bearer token dalam localStorage/sessionStorage atau response yang dapat dibaca JS. Jika kontrak PRD ambigu tentang response token, jelaskan keputusan di `docs/security-decisions.md`.
3. Terapkan auth middleware dan RBAC `STUDENT`/`ADMIN`; endpoint admin harus ditolak untuk student. User tidak dapat membaca atau mengubah data pengguna lain kecuali rule admin mengizinkan.
4. Akun nonaktif tidak boleh login. Credential salah menghasilkan error yang tidak membocorkan apakah identifier terdaftar. Tambahkan rate limiting untuk login.
5. Validasi identifier/email unik dan payload; tangani conflict database menjadi error API yang tepat.
6. Lindungi cookie-based auth dari CSRF berdasarkan arsitektur frontend/origin yang ditemukan (SameSite dan CSRF token bila memang diperlukan). CORS tidak boleh wildcard dengan credentials.
7. Audit event auth tanpa mencatat password, cookie, token, atau secret.
8. Sediakan unit/API tests untuk login valid/invalid, akun nonaktif, logout, auth tanpa cookie, role student/admin, duplicate email/identifier, dan password tidak bocor.

Jangan membangun modul lain atau mengubah UI. Periksa hasil dengan test, `go test ./...`, dan lint yang tersedia. Perbarui OpenAPI, `.env.example`, dan dokumentasi.
```

**Selesai jika:** auth end-to-end di backend dan RBAC teruji; secret tidak keluar ke log/response; semua endpoint user terdokumentasi.

# PROMPT 05 — Lab Sessions API

```text
Implementasikan modul lab sessions di Go + Fiber dengan kontrak `docs/api-contract.md` dan business rule PRD. Jangan membangun attendance pada fase ini.

Kebutuhan:
1. Admin dapat membuat, membaca daftar/detail, memperbarui, dan mengubah status session sesuai endpoint yang dirancang.
2. Student hanya dapat membaca session yang memang tersedia baginya; ia tidak dapat mengelola session.
3. Validasi code, title, lokasi, tanggal, start/end, status, dan rentang waktu. Waktu disimpan dengan semantik timezone yang konsisten dan didokumentasikan.
4. Session yang sudah direferensikan oleh attendance tidak dihapus fisik; gunakan status `CANCELLED`/`CLOSED` sesuai PRD.
5. Definisikan secara eksplisit cara menentukan session `ACTIVE`—berdasarkan waktu atau status—dan pastikan sama di backend serta API contract. Bila PRD menyisakan pilihan, catat asumsi.
6. Tambahkan pagination/filter bila kontrak membutuhkannya, test RBAC dan validasi rentang waktu, serta dokumentasi OpenAPI.

Jangan mengubah frontend atau modul lain. Jalankan test dan laporkan acceptance criteria satu per satu.
```

**Selesai jika:** CRUD session dan akses berdasarkan role teruji; attendance belum disentuh.

# PROMPT 06 — Python FastAPI Face Service (mandiri/internal)

```text
Siapkan microservice Python FastAPI untuk face verification sebagai service internal saja. Gunakan PRD Part face verification dan kontrak internal yang disepakati. Jangan hubungkan ke frontend dan jangan menambahkan attendance logic ke Python service.

Scope:
1. Buat struktur service, environment config, dependency pins, health endpoint, validasi MIME/format/ukuran file, batas request, structured errors, request ID dan logging aman.
2. Gunakan DeepFace untuk pipeline pretrained embedding/verification. Uji kompatibilitas ArcFace sebagai kandidat model; jangan mengklaim model/threshold final sebelum test lokal. Catat model name/version.
3. Validasi file sebagai image, deteksi tepat satu wajah, validasi ukuran/quality minimal yang tersedia, dan gagal secara terstruktur jika tidak ada wajah/lebih dari satu wajah/model error.
4. Berikan operasi internal untuk mengekstrak embedding dari satu pose enrollment (FRONT/LEFT/RIGHT) dan untuk membandingkan satu live embedding atau burst singkat terhadap tiga embedding enrollment milik user yang dipilih Go backend. Backend tetap menentukan user identity; service tidak melakukan pencarian 1:N.
5. Python service tidak menyimpan user account, password, sesi, attendance, atau raw image. Image diproses di memory lalu dibuang. Jangan gunakan Supabase Storage untuk foto wajah tanpa keputusan eksplisit.
6. Jangan hardcode threshold atau strategi sampel. Dokumentasikan aturan keputusan yang membandingkan live sample terhadap FRONT/LEFT/RIGHT, lakukan kalibrasi false acceptance/false rejection, dan jangan aktifkan verification untuk presensi saat calibration belum selesai.
7. Lindungi endpoint dari akses publik: dokumentasikan jaringan internal/secret service-to-service, dan jangan memasukkan embedding/raw image ke log.
8. Tambahkan unit tests untuk image invalid, file terlalu besar, no-face, multi-face, model error, dan response shape. Mock model pada test unit; pisahkan integration tests yang memerlukan model download/resource.
9. Dokumenkan kompatibilitas runtime, ukuran model, kebutuhan CPU/memory, lisensi dependensi/model yang perlu diverifikasi, dan cara menjalankan service lokal.

Jangan mengubah Go atau frontend kecuali untuk menambahkan kontrak terdokumentasi. Jalankan test yang tidak memerlukan download model dan laporkan test yang belum bisa jalan berikut penyebabnya.
```

**Selesai jika:** service dapat start, health check lulus, input/error terkontrol, tidak menyimpan gambar, dan test aman tersedia. Belum boleh membuat presensi dengan threshold yang belum dikalibrasi.

# PROMPT 07 — Integrasi Go ↔ Face Service dan Face Enrollment

```text
Integrasikan Go backend dengan Python face service untuk enrollment dan status template. Jangan implementasikan attendance/check-in pada fase ini.

Kebutuhan:
1. Browser hanya mengirim image melalui Go API; Go memvalidasi autentikasi, role, ownership, batas file, lalu memanggil Python service secara internal dengan timeout dan autentikasi service-to-service.
2. Enrollment hanya untuk pengguna yang sedang login (`user_id` dari session/cookie), bukan `user_id` bebas yang dikirim client. Re-enrollment mengikuti rule PRD dan tidak mengganti template diam-diam.
3. Ekstrak embedding untuk ketiga pose dan simpan pada `face_template_samples.embedding_ciphertext`; `face_templates` adalah metadata set. Simpan `consent_version`/`consented_at`, pose, model version, dan quality status. Enrollment parent dan seluruh sampel harus tersimpan dalam satu transaction setelah semua pose lolos. Enkripsi/key management dicatat di `docs/security-decisions.md`; jangan menyimpan secret key di database/source code.
4. Catat model name/version, status, timestamps, dan persetujuan. Simpan sampel berlabel `FRONT`, `LEFT`, `RIGHT`; embedding plaintext hanya ada di memory selama diperlukan. Jangan kirim embedding atau raw image kembali ke frontend atau log.
5. Tangani error service tidak tersedia, timeout, image invalid, no-face, multi-face, dan template already exists melalui error code yang konsisten.
6. Jangan upload image ke Supabase Storage; ini sudah menjadi keputusan MVP pada PRD revisi. Citra mentah hanya diproses di memory lalu dibuang.
7. Tes enrollment success (mock service), no face, multi-face, no session, inactive user, replace/re-enroll rules, no template leakage, service timeout.
8. Update OpenAPI dan dokumentasi alur enrollment.

Jangan mengimplementasikan attendance. Jalankan semua test Go terkait dan test Python yang relevan.
```

**Selesai jika:** enrollment mematuhi ownership dan privacy; template tersimpan terlindungi; backend tidak bergantung pada endpoint Python publik.

# PROMPT 08 — Attendance Check-in, Check-out, dan Verification Gate

```text
Implementasikan attendance di Go + Fiber sesuai PRD dan sequence/activity diagram. Modul hanya boleh membuat attendance jika seluruh prasyarat lulus.

Kebutuhan:
1. Check-in menerima `session_id` dan live image/burst via multipart, lalu backend menjalankan face verification server-side terhadap enrollment set lengkap (FRONT/LEFT/RIGHT) milik `user_id` yang login. Jangan percaya boolean `verified` dari client.
2. Server menentukan `checked_in_at` dan `checked_out_at`; jangan menerima timestamp client sebagai sumber kebenaran.
3. Satu user hanya memiliki satu attendance per session. Terapkan unique constraint database dan tangani race condition/double-submit; pengecekan aplikasi saja tidak cukup.
4. Jika verification gagal/error, tidak ada attendance. Bedakan `VERIFICATION_FAILED` dari camera/upload/service error. Check-in harus menjalankan verification dan insert attendance sebagai satu alur server-side; frontend tidak boleh mengirim keputusan verified sebagai otoritas.
5. Check-out memerlukan check-in yang ada dan belum check-out; operasi duplikat ditolak.
6. Admin dapat melihat/monitor attendance; student hanya melihat miliknya. Jangan sediakan update attendance oleh student di luar flow resmi.
7. Attendance response tidak berisi embedding, image, internal score yang dibatasi, atau detail model sensitif.
8. Tambahkan tests untuk session inactive, missing template, verification fail, duplicate, race/unique violation, check-out before check-in, repeated check-out, role access, timestamps server, rollback bila DB gagal.
9. Jangan menetapkan angka threshold arbitrer. Jika calibration belum disetujui, buat verification adapter mudah dites dan blokir mode demo/production dari menerima threshold yang tidak terkalibrasi; dokumentasikan cara mengkalibrasinya.
10. Update OpenAPI dan dokumentasi end-to-end flow.

Jangan mengubah UI atau asset modules. Jalankan test berulang untuk memastikan duplicate check-in tidak terjadi.
```

**Selesai jika:** business rules presensi dan constraint database teruji; verification gagal tidak membuat record.

# PROMPT 09 — Asset Management

```text
Implementasikan CRUD asset di Go + Fiber berdasarkan PRD dan `docs/api-contract.md`. Scope hanya master asset.

Kebutuhan:
1. Admin dapat list/detail/create/update dan mengubah status asset sesuai kontrak. Student tidak dapat mengelola asset; student hanya dapat melihat aset/alokasi miliknya melalui endpoint self-scoped yang telah dirancang.
2. `asset_code` wajib unik. Serial number unik jika nilai diberikan dan keputusan data model menetapkannya.
3. Terapkan status yang disetujui (`AVAILABLE`, `ALLOCATED`, `MAINTENANCE`, `RETIRED`) dan cegah transisi status ilegal. Jangan mengizinkan edit manual yang merusak allocation/maintenance aktif.
4. Gunakan DTO agar field sensitif/internal tidak bocor. Implementasikan filter/search/pagination hanya bila kontrak mendefinisikannya.
5. Catat event penting pada audit log tanpa payload rahasia.
6. Buat tests CRUD, validation, duplicate asset code, role, status transition, not found, dan konflik saat asset sedang dialokasikan/maintenance.
7. Update OpenAPI dan dokumentasi.

Jangan implementasikan allocation atau maintenance dalam prompt ini, selain guard agar data tidak menjadi inconsistent.
```

**Selesai jika:** CRUD asset teruji dan semua mutasi admin-only.

# PROMPT 10 — Allocation, Return, dan Asset History

```text
Implementasikan allocation/return dan asset history. Ikuti diagram status asset serta PRD.

Kebutuhan:
1. Admin mengalokasikan satu unit asset kepada user dan lab session yang valid. Validasi availability, user/session aktif, dan aturan attendance prerequisite seperti yang telah diputuskan dalam `docs/requirements-questions.md`.
2. Allocation insert dan update status asset harus satu transaksi database. Jika salah satu gagal, semuanya rollback.
3. Cegah dua allocation aktif atas asset yang sama menggunakan constraint database dan transaksi/locking yang sesuai untuk PostgreSQL.
4. Return mencatat `returned_at` dari server. Kondisi setelah return harus menentukan status asset: `AVAILABLE` bila layak, `MAINTENANCE` bila rusak; ikuti API contract untuk request kondisi return.
5. Asset berstatus `MAINTENANCE` atau `RETIRED` tidak dapat dialokasikan.
6. Riwayat menampilkan urutan allocation, return, maintenance dan event terkait. Student hanya boleh melihat allocation/history miliknya; admin bisa menelusuri berdasarkan asset/user/session.
7. Catat audit event create allocation dan return. Jangan menghapus histori operasional secara fisik.
8. Buat test untuk double allocation, race condition, transaction rollback, return twice, asset maintenance, ownership scoping, dan timeline.
9. Update OpenAPI serta dokumentasi business rules.

Jangan mengubah frontend dalam fase ini.
```

**Selesai jika:** status asset dan allocation selalu konsisten, termasuk pada request bersamaan.

# PROMPT 11 — Maintenance dan Audit Trail

```text
Implementasikan modul maintenance dan audit trail sesuai PRD/diagram state.

Kebutuhan:
1. Admin dapat membuat maintenance record untuk asset, menetapkan issue/description/priority, memperbarui status (`OPEN`, `IN_PROGRESS`, `COMPLETED`, `CANCELLED` sesuai rule), dan mencatat resolution.
2. Ketika maintenance dimulai, status asset berubah secara konsisten ke `MAINTENANCE`; ketika selesai, asset kembali ke `AVAILABLE` atau `RETIRED` sesuai hasil dan keputusan bisnis. Gunakan transaction dan cegah perubahan ilegal.
3. Asset maintenance tidak bisa dialokasikan. Asset yang masih memiliki allocation aktif tidak bisa masuk maintenance sebelum aturan return/penanganan khusus dipenuhi.
4. Audit event untuk maintenance create/update/complete/cancel dan status asset. Audit log append-only dari perspektif API umum: tidak ada endpoint delete/edit audit log biasa.
5. Audit metadata tidak boleh mengandung password, cookie/token, raw biometric image, atau embedding plaintext.
6. Terapkan admin-only permissions, validation, structured errors dan tests untuk semua transisi status, race/conflict, dan akses.
7. Dokumentasikan apakah histori perubahan maintenance memakai audit logs atau tabel khusus; jangan tambahkan tabel baru tanpa alasan dan persetujuan desain.
8. Update OpenAPI.

Jangan membangun laporan kompleks dulu. Jalankan tests dan cek kembali invariant allocation/maintenance.
```

**Selesai jika:** transisi maintenance/asset konsisten dan event kritis tercatat tanpa secret.

# PROMPT 12 — Dashboard dan Reporting API

```text
Implementasikan dashboard dan reporting backend berdasarkan data yang telah tersedia dan PRD. Jangan menambahkan analitik di luar scope.

Kebutuhan:
1. Dashboard student: status/riwayat presensi sendiri, session yang relevan, dan allocation sendiri.
2. Dashboard admin: session aktif, ringkasan attendance, count asset per status, allocation aktif, maintenance yang perlu ditangani, dan ringkasan aktivitas yang didukung data.
3. Reports minimum: attendance, asset, allocation, maintenance. Filter tanggal/entity harus divalidasi dan dibatasi berdasarkan role.
4. Gunakan aggregation/query efisien; hindari N+1 queries. Angka dashboard harus cocok dengan query/list sumbernya.
5. CSV export dan filter lanjutan adalah Should Have: implementasikan hanya setelah seluruh Must Have stabil dan test dasar lulus. Jika belum, dokumentasikan sebagai backlog, jangan menunda core API.
6. Tests authorization, filters, date range, empty states, aggregation consistency dan CSV injection mitigation jika export dibuat.
7. Update OpenAPI dengan schema dan query parameters.

Jangan menyentuh UI kecuali memperbarui kontrak.
```

**Selesai jika:** summary/report sesuai database, filter tidak melewati authorization, dan query aman.

# PROMPT 13 — Frontend API Client dan Integrasi Auth (UI dipertahankan)

```text
Sekarang integrasikan frontend yang sudah ada dengan backend nyata, mulai dari fondasi API client dan authentication. Baca ulang audit frontend serta seluruh halaman terkait. Ini bukan tugas redesign.

Aturan keras:
- Jangan mengganti React/Vite/Tailwind, jangan menghapus komponen, jangan merombak desain, jangan membuat aplikasi frontend kedua.
- Jangan melakukan query Supabase/PostgreSQL langsung dari browser. Semua data melalui Go REST API.
- Pertahankan route, visual identity, layout, komponen reusable, responsive behavior, dan bahasa yang telah digunakan kecuali ada bug nyata.

Implementasikan:
1. Konfigurasi base URL API via environment variable frontend (`VITE_API_BASE_URL` atau nama yang konsisten dengan repo); `.env.example` hanya placeholder.
2. API client terpusat dengan JSON/multipart handling, cookie credentials sesuai arsitektur, request ID/error parsing, timeout jika sesuai, dan tipe TypeScript untuk DTO berdasarkan `docs/api-contract.md`.
3. Integrasikan login, logout, session restoration dengan `/auth/me`, role guard, dan redirect yang memang sesuai router existing.
4. Tangani loading, error, empty state, session expired, dan validation error menggunakan komponen/pola UI existing.
5. Jangan simpan token dalam localStorage/sessionStorage. Cookie HttpOnly dikelola browser; client tidak membaca token.
6. Periksa CSRF/CORS/cookie behavior lintas origin dan dokumentasikan env dev/prod. Jangan membuka CORS wildcard dengan credentials.
7. Tambahkan test API client/auth memakai tooling yang memang sudah ada.

Selesaikan hanya fondasi API/auth. Jangan integrasikan semua modul sekaligus. Tampilkan file yang berubah dan yakinkan tidak ada perubahan desain yang tidak perlu.
```

**Selesai jika:** frontend dapat login/logout/restore auth lewat API tanpa token JS, dan UI yang ada tetap utuh.

# PROMPT 14 — Integrasi Frontend Mahasiswa

```text
Integrasikan layar mahasiswa yang benar-benar ada dengan API backend. Gunakan `docs/api-contract.md` dan komponen existing; jangan membangun halaman duplikat dan jangan redesign.

Integrasikan yang tersedia dalam frontend:
1. profile dan own data;
2. daftar/detail session yang tersedia;
3. multi-angle face enrollment melalui panduan kamera FRONT/LEFT/RIGHT dan satu multipart upload ke Go API;
4. presensi/check-in dan check-out dengan flow verification;
5. status dan histori attendance milik sendiri;
6. informasi allocation dan asset history milik sendiri.

Aturan:
- Browser menangani izin kamera dengan `getUserMedia()` hanya pada layar yang memerlukan kamera. Jangan merekam/menyimpan video.
- Kirim image ke Go API, bukan langsung ke Python face service atau Supabase Storage.
- UI harus membedakan camera permission error, image invalid/no face, verification failed, server unavailable, dan sukses.
- Jangan menganggap check-in sukses sampai backend mengonfirmasi.
- Jangan menampilkan embedding, raw face image permanen, similarity score sensitif atau detail internal.
- Tangani retry dan loading state agar double-click tidak menghasilkan request berulang; backend tetap menjadi sumber pencegahan duplikasi.

Tambahkan/ubah tests yang sesuai, jalankan lint/build/test frontend. Laporkan layar yang sudah/ belum terintegrasi karena endpoint/komponen belum tersedia. Jangan menyentuh halaman admin.
```

**Selesai jika:** alur mahasiswa terhubung ke API tanpa akses database langsung dan UI masih konsisten dengan desain semula.

# PROMPT 15 — Integrasi Frontend Admin

```text
Integrasikan layar admin yang benar-benar ada ke REST API tanpa redesign dan tanpa membuat halaman duplikat. Ikuti API contract, RBAC, komponen tabel/form/modal dan pola UI existing.

Integrasikan sesuai halaman yang ditemukan:
1. dashboard admin;
2. users dan status akun;
3. lab sessions;
4. asset master;
5. allocation dan return;
6. maintenance/status transition;
7. attendance monitoring;
8. reports/filter dan CSV bila endpoint CSV sudah selesai.

Persyaratan:
- Semua mutasi memakai Go REST API; jangan query Supabase langsung.
- Setiap aksi menampilkan feedback sukses/gagal dari response backend.
- Tangani loading, error, empty, pagination/filter jika API mendukungnya, serta validasi form.
- Jangan menampilkan field `password_hash`, embedding, raw biometric image, token, dan metadata internal sensitif.
- Hide/disable control sesuai role untuk UX, tetapi backend tetap menjadi enforcement authority.
- Perhatikan status transitions: jangan membuat UI menawarkan transisi yang backend larang.
- Tambahkan tests yang cocok dengan tooling existing, lalu jalankan lint/build/test.

Jangan mengubah halaman mahasiswa atau melakukan redesign global.
```

**Selesai jika:** layar admin yang ada memakai data backend dan tidak ada akses database dari browser.

# PROMPT 16 — End-to-End Testing, Security Review, API Docs, dan Deployment Readiness

```text
Lakukan hardening dan verifikasi akhir terhadap seluruh sistem, jangan menambah fitur baru kecuali perbaikan blocker.

Cakupan:
1. Jalankan format/lint/test/build untuk Go, Python, dan frontend menggunakan command yang tercatat di repo. Perbaiki regresi yang disebabkan implementasi proyek.
2. Tambahkan atau lengkapi integration/E2E tests untuk alur utama:
   - login -> `/auth/me` -> logout;
   - admin create user/session;
   - student enroll face;
   - verification berhasil -> check-in -> check-out;
   - duplicate attendance ditolak;
   - admin create asset -> allocate -> return;
   - maintenance -> asset unavailable -> completion/status restore;
   - dashboard/report role access.
3. Pastikan semua data test bersifat fiktif dan database test tidak merusak database pengguna/production.
4. Tinjau security: cookies, CSRF/CORS, RBAC dan object ownership, SQL parameterization, validation, request size, rate limits, upload security, secrets, log redaction, audit log, error response, database constraints, service-to-service auth, dependency/license concerns face model.
5. Verifikasi tidak ada Supabase key/service-role/DB password di frontend bundle, git history baru, logs, test fixtures, screenshot, atau docs.
6. Pastikan raw camera images tidak disimpan permanen sesuai keputusan privasi; pastikan embedding dilindungi; jangan membuat bucket foto wajah tanpa persetujuan tertulis.
7. Lengkapi `docs/openapi.yaml` dan Swagger UI sesuai endpoint aktual; dokumentasikan cara menjalankan semua services lokal, environment variables, migration, seed, test, deployment, CORS, HTTPS, dan rollback.
8. Jangan deploy atau menjalankan migration pada Supabase remote tanpa konfirmasi eksplisit. Jangan mengubah data production.
9. Periksa kesesuaian route frontend dengan OpenAPI; cari mock data yang masih digunakan dan laporkan mana yang sengaja dipertahankan.
10. Beri status setiap acceptance criterion: PASS, FAIL, atau NOT RUN. Jangan mengklaim test lulus jika tidak dijalankan.

Buat `docs/implementation-status.md` berisi ringkasan modul, test command + hasil sebenarnya, blocker, environment yang belum disiapkan, serta daftar langkah manual deployment. Akhiri dengan ringkasan perubahan dan sisa pekerjaan.
```

**Selesai jika:** semua hasil test yang bisa dijalankan tercatat secara jujur, API docs cocok dengan implementasi, isu keamanan kritis tertangani, dan deployment remote tetap menunggu persetujuan.

---

## Checklist sebelum memulai Supabase sungguhan

- [ ] Audit dan kontrak API disetujui.
- [ ] Nama tabel/kolom/status serta migration final telah direview.
- [ ] Kebijakan biometric, consent, revoke/re-enrollment, dan retensi jelas; Supabase Storage tidak digunakan untuk raw face image pada MVP.
- [ ] Project Supabase development dibuat manual oleh pemilik project.
- [ ] Connection string tersedia untuk backend lokal di `.env`, bukan di frontend dan bukan di Git.
- [ ] `.env` masuk `.gitignore`; periksa `git status` sebelum commit.
- [ ] Migration pertama diuji di local/dev database sebelum ke remote.
- [ ] Tidak menggunakan data mahasiswa nyata untuk seed/test.
- [ ] Secret service-to-service untuk Go ↔ Python dibuat melalui environment variable.

## Catatan prioritas

Untuk MVP, selesaikan prompt 00–12 secara berurutan, lalu integrasikan frontend pada prompt 13–15 dan validasi keseluruhan melalui prompt 16. Fitur Should Have/Nice to Have ditunda bila mengancam login, keamanan, face verification, attendance, asset allocation/return, maintenance, dan API documentation.
