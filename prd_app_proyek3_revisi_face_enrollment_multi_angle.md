op# PRODUCT REQUIREMENT DOCUMENT (PRD)

## RANCANG BANGUN WEB SERVICE PRESENSI DAN MANAJEMEN ASET LABORATORIUM KOMPUTER MENGGUNAKAN METODE AGILE

**Program Studi:** D4 Teknik Informatika
**Jenis Proyek:** Proyek 3 – Backend Developer
**Jumlah Anggota:** 2 mahasiswa
**Platform:** Web-based / Web-only
**Metode Pengembangan:** Agile
**Versi Dokumen:** 1.1 (draft revisi face enrollment multi-angle)

---

# ========================================

# PART 1 — IDENTITAS DAN RINGKASAN PROYEK

# ========================================

## 1.1 Nama Proyek

**RANCANG BANGUN WEB SERVICE PRESENSI DAN MANAJEMEN ASET LABORATORIUM KOMPUTER MENGGUNAKAN METODE AGILE**

## 1.2 Deskripsi Singkat Sistem

Sistem yang dirancang merupakan web service terintegrasi untuk mendukung proses presensi kegiatan praktikum dan pengelolaan aset laboratorium komputer melalui pendekatan arsitektur berbasis RESTful Web Service/API. Sistem terdiri atas dua aplikasi web, yaitu **Web Mahasiswa/User** dan **Web Admin/Laboran/Tata Usaha**, yang seluruh akses datanya dilakukan melalui Backend REST API tanpa akses langsung ke database.

Sistem mengintegrasikan data pengguna, sesi praktikum, presensi, face verification melalui kamera browser, data aset laboratorium, alokasi perangkat kepada mahasiswa, histori penggunaan perangkat, proses maintenance, serta pelaporan. Setiap unit perangkat dicatat secara individual menggunakan kode aset sehingga histori penggunaan dan kondisi perangkat dapat ditelusuri.

Proses presensi menggunakan **face verification 1:1**. Mahasiswa melakukan login terlebih dahulu sehingga sistem telah mengetahui identitas pengguna yang sedang aktif. Sistem kemudian menerima citra wajah dari kamera browser dan membandingkannya dengan template wajah milik pengguna tersebut. Dengan demikian, sistem tidak melakukan pencarian identitas berdasarkan seluruh pengguna atau face identification 1:N sebagai mekanisme utama.

## 1.3 Tujuan Produk

Tujuan produk adalah menghasilkan sistem web yang mampu menyediakan layanan presensi dan manajemen aset laboratorium secara terintegrasi, terstruktur, terdokumentasi, dan dapat dikembangkan dalam waktu satu semester oleh dua mahasiswa.

Secara khusus, sistem bertujuan menyediakan mekanisme autentikasi pengguna, pendaftaran template wajah, face verification 1:1, pencatatan kehadiran pada sesi praktikum, pencatatan penggunaan unit perangkat secara individual, pengelolaan aset, pencatatan maintenance, serta penyediaan laporan berbasis REST API.

Produk juga bertujuan menjadikan REST API sebagai pusat komunikasi antara kedua aplikasi web dengan database sehingga pemisahan antara antarmuka, logika bisnis, dan pengelolaan data dapat diterapkan secara jelas.

## 1.4 Visi Produk

Visi produk adalah membangun layanan digital laboratorium komputer yang terintegrasi, terdokumentasi, dan mudah ditelusuri, dengan memanfaatkan RESTful Web Service sebagai fondasi utama pertukaran data dan face verification sebagai mekanisme verifikasi identitas pada proses presensi.

Sistem diarahkan untuk menghasilkan satu sumber data terpusat mengenai kegiatan praktikum, presensi mahasiswa, penggunaan unit aset, dan kondisi perangkat sehingga informasi operasional laboratorium dapat dikelola secara lebih konsisten.

## 1.5 Ringkasan Masalah

Kegiatan laboratorium komputer melibatkan beberapa jenis data yang saling berkaitan, antara lain data mahasiswa, jadwal atau sesi praktikum, kehadiran, penggunaan komputer atau laptop, serta kondisi perangkat. Apabila data tersebut dicatat secara terpisah, proses penelusuran hubungan antara mahasiswa, sesi praktikum, perangkat yang digunakan, dan histori kondisi perangkat menjadi lebih sulit.

Pada proses presensi, identitas mahasiswa perlu diverifikasi sebelum kehadiran dinyatakan valid. Penggunaan face verification 1:1 melalui kamera browser memberikan mekanisme verifikasi yang terhubung langsung dengan akun pengguna yang telah login.

Pada sisi aset, setiap perangkat membutuhkan identitas individual sehingga informasi seperti asset code, kondisi perangkat, pengguna yang menggunakan perangkat, waktu alokasi, pengembalian, dan riwayat maintenance dapat ditelusuri melalui sistem.

## 1.6 Solusi yang Ditawarkan

Solusi yang ditawarkan adalah satu ekosistem web yang terdiri dari Web Mahasiswa/User dan Web Admin/Laboran/Tata Usaha. Kedua aplikasi menggunakan REST API sebagai satu-satunya jalur akses ke data aplikasi.

Pada sisi mahasiswa, sistem menangani login, enrollment wajah, verifikasi wajah, presensi pada sesi praktikum, pemantauan status presensi, informasi alokasi perangkat, serta histori pribadi yang diizinkan.

Pada sisi admin, sistem menangani pengelolaan pengguna, sesi praktikum, aset, alokasi perangkat, maintenance, dashboard, dan laporan.

Face verification menggunakan model face recognition yang telah dilatih sebelumnya. Sistem tidak melakukan training model face recognition dari nol. Model digunakan untuk menghasilkan representasi atau embedding wajah yang dibandingkan terhadap embedding milik pengguna yang telah terdaftar.

## 1.7 Target Pengguna

Target pengguna sistem terdiri atas:

1. **Mahasiswa**, sebagai pengguna yang melakukan autentikasi, enrollment wajah, verifikasi wajah, presensi, dan menggunakan perangkat laboratorium.
2. **Admin/Laboran/Tata Usaha**, sebagai pengguna operasional yang mengelola data pengguna, sesi praktikum, aset, alokasi, maintenance, dan laporan.

Role lain tidak menjadi bagian dari scope utama karena kebutuhan pengelolaan dapat ditangani oleh satu role operasional admin.

## 1.8 Nilai dan Manfaat Sistem

Sistem memberikan nilai berupa integrasi proses presensi dan pengelolaan aset dalam satu layanan data. Setiap kehadiran dapat dikaitkan dengan sesi praktikum, sedangkan setiap penggunaan aset dapat dikaitkan dengan mahasiswa dan sesi tertentu.

Bagi mahasiswa, sistem menyediakan proses presensi yang lebih terintegrasi dengan identitas akun dan menyediakan informasi penggunaan aset. Bagi admin, sistem menyediakan kemampuan pelacakan aset, riwayat penggunaan, kondisi perangkat, maintenance, serta laporan operasional.

Bagi pengembangan perangkat lunak, sistem menyediakan implementasi nyata RESTful Web Service, database relasional, autentikasi, authorization, pengelolaan data, integrasi face verification, dokumentasi API, dan penerapan Agile pada proyek dua mahasiswa.

## 1.9 Konsep Inovasi

Konsep inovasi sistem terletak pada integrasi **presensi berbasis face verification 1:1** dengan **manajemen aset individual laboratorium** melalui satu RESTful Web Service.

Face verification digunakan setelah mahasiswa login sehingga identitas pengguna telah diketahui terlebih dahulu. Pendekatan ini berbeda dari sistem identifikasi 1:N yang harus mencari identitas pengguna dari seluruh basis data wajah.

Inovasi lainnya adalah keterkaitan data presensi dengan konteks sesi praktikum dan keterkaitan penggunaan perangkat dengan identitas mahasiswa serta sesi tersebut. Integrasi tersebut memungkinkan histori aset ditelusuri dari sisi penggunaan maupun maintenance.

Sistem dirancang untuk **mengurangi potensi penyalahgunaan atau titip absen**, bukan untuk menyatakan bahwa penyalahgunaan dapat dicegah secara absolut.

---

# ========================================

# PART 2 — LATAR BELAKANG DAN IDENTIFIKASI MASALAH

# ========================================

## 2.1 Latar Belakang

Laboratorium komputer merupakan fasilitas pendukung kegiatan akademik yang melibatkan penggunaan perangkat komputasi oleh mahasiswa pada waktu dan sesi tertentu. Aktivitas laboratorium tidak hanya membutuhkan pencatatan kehadiran, tetapi juga membutuhkan pengelolaan penggunaan perangkat, identitas unit aset, kondisi perangkat, dan histori maintenance.

Presensi merupakan salah satu data penting dalam kegiatan praktikum. Sistem presensi harus mampu menghubungkan identitas mahasiswa dengan sesi praktikum tertentu dan memastikan bahwa data kehadiran tercatat dalam konteks sesi yang benar. Karena presensi menggunakan pendekatan biometrik, proses verifikasi identitas harus dilakukan dengan mempertimbangkan aspek keamanan, privasi, kualitas citra, dan keterbatasan teknologi.

Selain presensi, penggunaan perangkat laboratorium perlu dicatat secara individual. Perangkat seperti laptop atau komputer harus memiliki asset code dan identitas unik agar sistem dapat membedakan setiap unit. Contoh data perangkat adalah Dell Latitude 5490 Core i7 dengan Asset Code `LAB-LT-007`. Contoh tersebut hanya merupakan format representasi unit, bukan data aktual laboratorium.

Kebutuhan pelacakan aset semakin penting ketika perangkat mengalami kerusakan atau membutuhkan maintenance. Histori penggunaan dan maintenance perlu tersimpan sehingga admin dapat mengetahui hubungan antara aset, mahasiswa, sesi praktikum, dan kejadian maintenance.

Pengelolaan seluruh proses tersebut membutuhkan integrasi data. Sistem yang dibangun menggunakan Backend REST API sebagai pusat komunikasi antara aplikasi web dan database. Web Mahasiswa dan Web Admin tidak melakukan koneksi database secara langsung. Seluruh proses pengambilan, perubahan, validasi, dan penyimpanan data dilakukan melalui backend.

Arsitektur tersebut memberikan pemisahan tanggung jawab yang jelas antara antarmuka pengguna, business logic, pengelolaan face verification, dan database. Pendekatan ini juga menghasilkan luaran Web Services/Doc API yang dapat digunakan sebagai dasar implementasi serta dokumentasi proyek.

## 2.2 Identifikasi Masalah

Permasalahan utama yang diidentifikasi adalah:

1. Data presensi membutuhkan mekanisme verifikasi identitas mahasiswa yang terintegrasi dengan akun pengguna.
2. Penggunaan aset laboratorium perlu dicatat secara individual berdasarkan unit perangkat.
3. Data presensi, sesi praktikum, aset, alokasi, dan maintenance perlu saling terhubung.
4. Pengelolaan data melalui akses database langsung dari frontend meningkatkan ketergantungan antarmuka terhadap struktur database.
5. Diperlukan API terpusat yang mengatur validasi, authorization, business logic, dan akses data.
6. Data biometrik membutuhkan mekanisme penyimpanan dan akses yang lebih ketat.
7. Histori penggunaan perangkat dan maintenance perlu tersedia dalam bentuk data terstruktur dan dapat dilaporkan.
8. Proyek harus tetap realistis untuk dikerjakan oleh dua mahasiswa dalam satu semester.

## 2.3 Rumusan Masalah

Rumusan masalah dalam proyek ini adalah:

1. Bagaimana membangun RESTful Web Service yang menjadi pusat komunikasi antara dua aplikasi web dan database?
2. Bagaimana menerapkan face verification 1:1 melalui kamera browser setelah mahasiswa melakukan login?
3. Bagaimana menghubungkan hasil verifikasi dengan pencatatan presensi pada sesi praktikum?
4. Bagaimana mengelola aset secara individual berdasarkan asset code dan histori penggunaan?
5. Bagaimana mencatat alokasi perangkat, pengembalian perangkat, serta maintenance secara terintegrasi?
6. Bagaimana menerapkan authentication, authorization, validation, dan perlindungan data biometrik secara realistis untuk proyek mahasiswa?
7. Bagaimana mengembangkan sistem secara bertahap menggunakan metode Agile?

## 2.4 Tujuan Pengembangan

Pengembangan bertujuan:

1. Menghasilkan Backend REST API yang menangani seluruh akses data.
2. Menghasilkan Web Mahasiswa.
3. Menghasilkan Web Admin/Laboran/Tata Usaha.
4. Menghasilkan mekanisme enrollment dan face verification 1:1.
5. Menghasilkan modul presensi yang tervalidasi berdasarkan akun dan sesi praktikum.
6. Menghasilkan modul manajemen aset dengan identitas individual setiap unit.
7. Menghasilkan modul alokasi, pengembalian, dan histori perangkat.
8. Menghasilkan modul maintenance dan pelaporan.
9. Menghasilkan dokumentasi API berbasis OpenAPI/Swagger.
10. Menghasilkan sistem yang dapat diuji secara fungsional dan terintegrasi.

## 2.5 Manfaat Pengembangan

Manfaat pengembangan bagi operasional laboratorium adalah tersedianya satu sistem terintegrasi untuk presensi dan pengelolaan aset. Bagi mahasiswa, sistem menyediakan proses identifikasi yang terhubung dengan akun dan sesi praktikum. Bagi admin, sistem membantu pengelolaan data aset, alokasi, maintenance, dan laporan.

Manfaat akademis proyek adalah tersedianya implementasi nyata konsep backend development, REST API, database relational, authentication, authorization, face verification, Agile, testing, dan dokumentasi teknis.

## 2.6 Batasan Pengembangan

Batasan sistem adalah sebagai berikut:

1. Sistem hanya berbasis web dan tidak menyediakan aplikasi Android atau iOS.
2. Sistem terdiri dari dua aplikasi web, yaitu Web Mahasiswa/User dan Web Admin/Laboran/Tata Usaha.
3. Seluruh akses data dilakukan melalui REST API.
4. Face verification menggunakan pendekatan 1:1.
5. Mahasiswa harus login sebelum proses face verification.
6. Sistem tidak menggunakan face identification 1:N sebagai fitur utama.
7. Sistem tidak melakukan training model face recognition dari nol.
8. Sistem menggunakan pretrained face model/library.
9. Sistem menggunakan kamera browser untuk memperoleh citra wajah.
10. Sistem tidak menggunakan QR Code sebagai metode utama presensi.
11. Sistem tidak membutuhkan perangkat keras khusus.
12. Sistem tidak menggunakan IoT, RFID, atau scanner biometrik khusus.
13. Anti-spoofing tingkat lanjut bukan bagian wajib MVP.
14. Threshold verifikasi tidak ditentukan secara arbitrer dan ditentukan berdasarkan pengujian.
15. Data mahasiswa dan data aset aktual belum ditentukan dan diberi penanda `[PERLU VALIDASI]`.
16. Integrasi dengan sistem akademik eksternal belum menjadi bagian scope utama.
17. Pembayaran, inventarisasi keuangan, pengadaan, dan pemusnahan aset tidak termasuk scope.
18. Sistem ditargetkan untuk skala laboratorium/proyek akademik, bukan sistem enterprise dengan multi-campus tenant.

---

# ========================================

# PART 3 — USER DAN ROLE

# ========================================

## 3.1 Aktor Sistem

### 3.1.1 Mahasiswa

Mahasiswa merupakan pengguna yang melakukan autentikasi, enrollment wajah, face verification, presensi, melihat informasi sesi, dan melihat informasi penggunaan aset yang berkaitan dengan dirinya.

### 3.1.2 Admin/Laboran/Tata Usaha

Admin/Laboran/Tata Usaha merupakan role operasional yang mengelola pengguna, sesi praktikum, aset, alokasi perangkat, maintenance, dashboard, dan laporan.

Role tambahan tidak digunakan dalam scope utama karena dua role telah mencukupi kebutuhan sistem.

## 3.2 Hak Akses Mahasiswa

Mahasiswa memiliki hak untuk:

* Login dan logout.
* Melihat profil sendiri.
* Melakukan face enrollment.
* Melakukan face verification.
* Melakukan attendance check-in pada sesi aktif.
* Melakukan attendance check-out.
* Melihat status presensi sendiri.
* Melihat histori presensi sendiri.
* Melihat sesi praktikum yang tersedia.
* Melihat alokasi aset miliknya.
* Melihat histori penggunaan aset miliknya.
* Tidak dapat mengelola data mahasiswa lain.
* Tidak dapat mengubah data aset.
* Tidak dapat membuat atau mengubah maintenance.
* Tidak dapat mengakses data embedding milik dirinya secara mentah.
* Tidak dapat mengakses database secara langsung.

## 3.3 Hak Akses Admin/Laboran/Tata Usaha

Admin memiliki hak untuk:

* Login dan logout.
* Mengelola data pengguna.
* Mengaktifkan/nonaktifkan pengguna.
* Mengelola sesi praktikum.
* Melihat data presensi.
* Mengelola aset.
* Mengelola alokasi aset.
* Mencatat pengembalian aset.
* Membuat dan memperbarui data maintenance.
* Melihat histori aset.
* Melihat dashboard.
* Mengakses laporan.
* Mengakses audit trail sesuai kebutuhan operasional.
* Tidak dapat mengakses password dalam bentuk plaintext.
* Tidak dapat mengakses raw biometric image yang telah dihapus setelah pemrosesan.
* Akses terhadap metadata verifikasi dibatasi sesuai authorization.

## 3.4 Matriks Hak Akses

| Role      | Modul             | Akses                                                           |
| --------- | ----------------- | --------------------------------------------------------------- |
| Mahasiswa | Authentication    | Login, logout                                                   |
| Mahasiswa | Profile           | Read profil sendiri                                             |
| Mahasiswa | Face Enrollment   | Create/replace template sesuai aturan                           |
| Mahasiswa | Face Verification | Verify diri sendiri                                             |
| Mahasiswa | Session           | Read sesi yang tersedia                                         |
| Mahasiswa | Attendance        | Check-in, check-out, read histori sendiri                       |
| Mahasiswa | Asset             | Read aset yang dialokasikan                                     |
| Mahasiswa | Allocation        | Read alokasi sendiri                                            |
| Mahasiswa | Maintenance       | Tidak ada akses pengelolaan                                     |
| Mahasiswa | Reporting         | Read laporan pribadi terbatas                                   |
| Admin     | Authentication    | Login, logout                                                   |
| Admin     | User Management   | CRUD terbatas                                                   |
| Admin     | Face Enrollment   | Melihat status enrollment, re-enrollment melalui prosedur admin |
| Admin     | Face Verification | Read status/metadata verifikasi                                 |
| Admin     | Session           | CRUD                                                            |
| Admin     | Attendance        | Read dan monitoring                                             |
| Admin     | Asset             | CRUD                                                            |
| Admin     | Allocation        | CRUD operasional                                                |
| Admin     | Maintenance       | CRUD operasional                                                |
| Admin     | Reporting         | Read/export                                                     |
| Admin     | Dashboard         | Read                                                            |
| Admin     | Audit Log         | Read terbatas                                                   |

---

# ========================================

# PART 4 — PRODUCT SCOPE

# ========================================

## 4.1 In Scope

Fitur yang berada dalam scope utama adalah:

1. Authentication.
2. User Management.
3. Face Enrollment.
4. Face Verification 1:1.
5. Session Praktikum.
6. Attendance.
7. Attendance Check-in.
8. Attendance Check-out.
9. Asset Management.
10. Asset Allocation.
11. Asset Return.
12. Asset History.
13. Maintenance.
14. Dashboard.
15. Reporting.
16. REST API.
17. API Documentation.
18. Audit Log dasar.
19. Security dasar aplikasi.
20. Testing.
21. Deployment web.

## 4.2 Out of Scope

Fitur berikut berada di luar scope:

1. Mobile application.
2. Face identification 1:N.
3. Training model face recognition dari nol.
4. IoT.
5. RFID.
6. Hardware biometric scanner.
7. Integrasi absensi hardware.
8. Sistem pengadaan barang.
9. Sistem keuangan aset.
10. Advanced AI analytics.
11. Predictive maintenance berbasis Machine Learning.
12. Chatbot.
13. Sistem akademik lengkap.
14. Integrasi SSO kampus yang belum tersedia spesifik `[PERLU VALIDASI]`.
15. Multi-campus enterprise architecture.
16. Advanced anti-spoofing yang membutuhkan riset khusus; basic liveness challenge hanya dapat ditambahkan sebagai peningkatan setelah MVP dan harus diuji terpisah.
17. Video recording mahasiswa.
18. Penyimpanan rekaman kamera.
19. Sinkronisasi dengan perangkat jaringan laboratorium.

## 4.3 Must Have

Fitur wajib selesai adalah:

* Login.
* Role-based authorization.
* User management dasar.
* Face enrollment.
* Face verification 1:1.
* Session management.
* Attendance.
* Asset management.
* Asset allocation.
* Asset return.
* Maintenance.
* Reporting dasar.
* Dashboard dasar.
* REST API.
* Swagger/OpenAPI.
* Database.
* Web Mahasiswa.
* Web Admin.
* Security dasar.
* Testing integrasi.

## 4.4 Should Have

Fitur yang termasuk prioritas setelah fungsi utama stabil adalah:

* CSV export.
* Filter laporan.
* Audit log yang lebih lengkap.
* Re-enrollment dengan approval admin.
* Ringkasan histori aset.
* Metadata similarity score untuk kebutuhan audit/pengujian dengan akses terbatas.

## 4.5 Nice to Have

Fitur tambahan yang dapat dikerjakan hanya setelah MVP selesai:

* Tampilan grafik dashboard.
* Rekap penggunaan aset per periode.
* Ringkasan maintenance per kategori aset.
* Pengaturan threshold melalui konfigurasi environment/admin secara terbatas.
* Anti-spoofing tambahan jika hasil evaluasi menunjukkan kebutuhan.

## 4.6 Alasan Pembatasan Scope

Pembatasan scope dilakukan karena proyek dikerjakan oleh dua mahasiswa dalam satu semester. Fokus utama proyek adalah implementasi backend service, REST API, database, face verification, dan dua aplikasi web. Fitur tambahan yang membutuhkan hardware, training model, integrasi eksternal kompleks, atau machine learning tambahan tidak dimasukkan karena berpotensi mengurangi keberhasilan penyelesaian sistem inti.

---

# ========================================

# PART 5 — FITUR PRODUK

# ========================================

## 5.1 Authentication

**Tujuan:** Menyediakan mekanisme autentikasi aman bagi mahasiswa dan admin.

**Aktor:** Mahasiswa, Admin/Laboran/Tata Usaha.

**Deskripsi:** Sistem menerima kredensial login, melakukan validasi terhadap data pengguna, memeriksa status akun, membuat token autentikasi, dan menetapkan sesi autentikasi.

**Precondition:**

* Pengguna telah terdaftar.
* Akun aktif.
* Password telah disimpan dalam bentuk hash.

**Main Flow:**

1. Pengguna membuka halaman login.
2. Pengguna mengisi username/email dan password.
3. Frontend mengirim request ke REST API.
4. Backend memvalidasi format request.
5. Backend mengambil user berdasarkan identifier.
6. Backend membandingkan password dengan password hash.
7. Backend membuat access token.
8. Token dikirim melalui mekanisme cookie aman.
9. Backend mengembalikan data profil dan role yang diizinkan.

**Alternative Flow:**

* Kredensial salah → `401 Unauthorized`.
* Akun nonaktif → `403 Forbidden`.
* Request tidak valid → `422 Unprocessable Entity`.
* Rate limit terlampaui → `429 Too Many Requests`.

**Output:** Sesi login aktif.

**Business Rule:**

* Password tidak pernah disimpan plaintext.
* Token tidak boleh dicatat dalam log.
* Akun nonaktif tidak dapat login.

**Acceptance Criteria:**

* Login dengan kredensial valid berhasil.
* Login dengan password salah ditolak.
* User mendapatkan role yang benar.
* Endpoint terlindungi tidak dapat diakses tanpa authentication.

## 5.2 User Management

**Tujuan:** Mengelola data pengguna dan role.

**Aktor:** Admin.

**Deskripsi:** Admin dapat membuat, membaca, memperbarui, mengaktifkan, dan menonaktifkan data pengguna.

**Precondition:** Admin telah login.

**Main Flow:** Admin membuka modul pengguna, membuat atau memperbarui data, backend melakukan validation, data disimpan, dan API mengembalikan data terbaru.

**Alternative Flow:** Email/identifier duplikat ditolak.

**Output:** Data pengguna tersimpan.

**Business Rule:**

* Role hanya berasal dari enum yang ditentukan.
* User nonaktif tidak dapat login.
* Password default, apabila digunakan, harus segera diganti dan tidak ditampilkan kembali.

**Acceptance Criteria:**

* Admin dapat membuat pengguna.
* Admin dapat menonaktifkan pengguna.
* Duplikasi identifier ditolak.
* Mahasiswa tidak dapat mengakses modul admin.

## 5.3 Face Enrollment Multi-Angle

**Tujuan:** Mendaftarkan representasi wajah pengguna dari beberapa sudut agar verifikasi berikutnya lebih toleran terhadap variasi posisi wajah.

**Aktor:** Mahasiswa.

**Deskripsi:** Setelah login, mahasiswa menyetujui pemberitahuan penggunaan data biometrik dan mengikuti panduan kamera untuk mengambil wajah dari tiga pose: **depan, sedikit menoleh ke kiri, dan sedikit menoleh ke kanan**. Sistem memproses setiap citra menjadi embedding menggunakan pretrained face model. Citra mentah hanya diproses sementara dan tidak disimpan permanen pada MVP. Sistem menyimpan tiga sampel embedding terenkripsi beserta label pose dan versi model.

**Precondition:**

* Pengguna telah login dan akun aktif.
* Pengguna telah membaca dan menyetujui pemberitahuan pemrosesan biometrik.
* Kamera browser tersedia dan izin kamera diberikan.
* Backend face service tersedia.

**Main Flow:**

1. User membuka halaman Face Enrollment.
2. Frontend menampilkan penjelasan tujuan, data yang diproses, masa retensi, dan cara meminta penghapusan atau pendaftaran ulang.
3. Browser meminta akses kamera.
4. UI memandu pengambilan pose depan, kiri, dan kanan satu per satu.
5. Untuk setiap pose, backend memvalidasi file, kualitas citra, dan bahwa tepat satu wajah terdeteksi serta pose cukup sesuai.
6. Jika citra valid, face service menghasilkan embedding untuk pose tersebut.
7. Setelah ketiga pose valid, backend menyimpan satu enrollment set aktif dan tiga sampel embedding terenkripsi dengan label `FRONT`, `LEFT`, dan `RIGHT`.
8. Backend mencatat model name/version, timestamp enrollment, dan status persetujuan; raw image dibuang setelah pemrosesan.
9. UI menampilkan enrollment berhasil.

**Alternative Flow:**

* Tidak ada wajah, lebih dari satu wajah, citra buram, terlalu gelap, atau wajah terlalu miring → pose terkait diminta diambil ulang.
* Kamera tidak diizinkan atau tidak tersedia → enrollment berhenti tanpa menyimpan template parsial.
* Face service gagal → tampilkan error yang aman; jangan menandai enrollment berhasil.
* Enrollment berhenti di tengah proses → sampel sementara dibuang atau enrollment ditandai belum aktif; hanya enrollment lengkap yang dapat dipakai untuk verifikasi.
* Template sudah aktif → pendaftaran ulang dilakukan dengan alur replace/revoke yang jelas; jangan menghapus template lama sebelum enrollment pengganti valid.

**Output:** Satu enrollment set aktif dengan sampel embedding terenkripsi untuk tiga pose. Tidak ada foto mentah yang disimpan permanen.

**Business Rule:**

* Enrollment menggunakan pretrained model; sistem tidak melatih model dari nol.
* Satu user memiliki maksimal satu enrollment set aktif pada MVP, yang dapat berisi beberapa sampel embedding.
* Setiap sampel menyimpan label pose, versi model, dan embedding terenkripsi.
* Jangan menyimpan foto mentah ke database atau Supabase Storage pada MVP.
* Kunci enkripsi disimpan sebagai secret di lingkungan backend, bukan di database atau repository.
* Perubahan model/preprocessing yang tidak kompatibel memerlukan rencana migrasi atau enrollment ulang.

**Acceptance Criteria:**

* Pengguna dapat menyelesaikan tiga pose dengan panduan yang jelas.
* Enrollment gagal jika satu atau lebih pose wajib belum valid.
* Sistem menolak citra tanpa wajah atau dengan lebih dari satu wajah.
* Enrollment set aktif hanya tersedia setelah ketiga sampel tersimpan berhasil.
* Status enrollment dapat dibaca melalui API, tetapi embedding tidak pernah dikirim ke frontend.
* Foto mentah tidak tersimpan pada database, Storage, log, atau folder permanen.

## 5.4 Face Verification 1:1

**Tujuan:** Memastikan bahwa wajah live dari kamera sesuai dengan enrollment set milik user yang sedang login.

**Aktor:** Mahasiswa.

**Deskripsi:** Saat presensi, user tidak perlu melakukan pendaftaran wajah ulang. Sistem mengambil satu frame frontal yang valid atau burst singkat beberapa frame, membentuk embedding, lalu membandingkannya **hanya dengan sampel embedding pada enrollment set milik user yang sedang login**. Aturan keputusan dan threshold ditentukan berdasarkan pengujian, bukan dipilih secara arbitrer.

**Precondition:**

* User login dan akun aktif.
* Enrollment set aktif dan lengkap.
* Session praktikum valid sesuai aturan presensi.
* Kamera tersedia.

**Main Flow:**

1. User membuka halaman presensi dan memilih session.
2. Backend memvalidasi session, status akun, dan keberadaan enrollment set aktif.
3. Browser membuka kamera dan menampilkan preview.
4. Frontend menangkap frame frontal atau burst singkat sesuai kontrak API.
5. Backend memvalidasi file, kualitas citra, dan tepat satu wajah.
6. Face service menghasilkan embedding dari citra live.
7. Backend mengambil semua sampel embedding aktif yang terkait dengan `user_id` dari sesi login; tidak ada pencarian identitas 1:N.
8. Face service menghitung distance/similarity antara embedding live dan sampel tersimpan menggunakan model dan metrik yang sama.
9. Backend menerapkan aturan keputusan yang telah ditetapkan melalui pengujian dan mencatat hasil verifikasi tanpa menyimpan citra mentah.
10. Jika terverifikasi, proses check-in boleh dilanjutkan; jika gagal, attendance tidak dibuat.

**Alternative Flow:**

* Enrollment belum tersedia atau tidak lengkap → arahkan user ke Face Enrollment.
* Tidak ada wajah, lebih dari satu wajah, atau kualitas citra buruk → minta capture ulang dalam batas percobaan.
* Wajah tidak cocok → kembalikan `VERIFICATION_FAILED`; jangan membuat attendance.
* Face service atau database tidak tersedia → kembalikan error sistem; jangan menganggapnya sebagai hasil verifikasi gagal dan jangan membuat attendance.
* Percobaan terlalu banyak → rate limit.

**Output:** Hasil `verified=true/false` dan kode hasil yang aman. Embedding, threshold internal, dan sampel biometrik tidak dikirim ke frontend.

**Business Rule:**

* Verifikasi selalu 1:1 menggunakan identitas user yang telah login.
* Threshold dan aturan pembandingan antarsampel harus ditetapkan melalui pengujian pasangan wajah positif/negatif dan dicatat versinya.
* Strategi keputusan tidak boleh sekadar menerima nilai maksimum similarity tanpa pengujian false acceptance dan false rejection.
* Foto mentah tidak disimpan permanen.
* Verifikasi kamera web 2D bukan pengganti keamanan Face ID/face unlock pada perangkat yang memiliki sensor kedalaman, pemrosesan aman, dan anti-spoofing khusus.

**Acceptance Criteria:**

* Pemilik enrollment set dapat diverifikasi pada kondisi kamera yang masuk kriteria pengujian.
* Wajah orang berbeda ditolak sesuai hasil evaluasi threshold.
* Backend hanya membaca enrollment set milik user login.
* Verifikasi gagal atau error tidak membuat attendance.
* Respons API tidak memuat embedding atau data biometrik mentah.

## 5.5 Session Praktikum

**Tujuan:** Mengelola konteks waktu dan kegiatan presensi.

**Aktor:** Admin.

**Deskripsi:** Session menyimpan informasi kegiatan praktikum, waktu mulai, waktu selesai, lokasi, dan status.

**Precondition:** Admin login.

**Main Flow:**

1. Admin membuat session.
2. Admin menetapkan tanggal, waktu, nama kegiatan, dan lokasi.
3. Backend memvalidasi rentang waktu.
4. Session disimpan.
5. Session berubah menjadi `ACTIVE` berdasarkan waktu atau status yang ditentukan.

**Alternative Flow:**

* Waktu selesai sebelum waktu mulai → ditolak.
* Session duplikat sesuai aturan bisnis → ditolak.

**Output:** Session tersimpan.

**Business Rule:**

* Attendance hanya dapat terjadi pada session yang valid dan aktif sesuai aturan.
* Session yang telah memiliki attendance tidak dihapus secara fisik; status dapat diubah menjadi `CANCELLED` atau `CLOSED`.

**Acceptance Criteria:**

* Session dapat dibuat.
* Status session dapat dipantau.
* Attendance terikat pada satu session.

## 5.6 Attendance

**Tujuan:** Mencatat dan mengelola kehadiran mahasiswa.

**Aktor:** Mahasiswa, Admin.

**Deskripsi:** Attendance menyimpan keterkaitan antara user dan session serta waktu check-in/check-out.

**Precondition:**

* User login.
* Session valid.
* Verification berhasil.

**Main Flow:**

1. User memilih session.
2. Sistem memvalidasi eligibility.
3. Sistem menjalankan face verification.
4. Sistem memeriksa duplicate attendance.
5. Sistem membuat attendance record.
6. API mengembalikan status hadir.

**Alternative Flow:**

* User telah memiliki attendance → duplicate ditolak.
* Session tidak aktif → presensi ditolak.
* Verification failed → attendance tidak dibuat.

**Output:** Record attendance.

**Business Rule:**

* Satu user hanya memiliki satu attendance per session.
* Duplicate attendance ditolak.
* Attendance harus memiliki user dan session valid.

**Acceptance Criteria:**

* Kehadiran hanya tersimpan setelah verification berhasil.
* Duplicate attendance tidak dapat dibuat.

## 5.7 Check-in

**Tujuan:** Menandai awal kehadiran mahasiswa.

**Aktor:** Mahasiswa.

**Precondition:** User login, session aktif, dan face verification berhasil.

**Main Flow:** User menjalankan verification → backend memvalidasi session → backend membuat attendance check-in.

**Alternative Flow:** Duplicate, session tidak aktif, atau verification failed.

**Output:** `checked_in_at`.

**Business Rule:** Waktu check-in dibuat oleh server, bukan dari waktu yang dikirim client.

**Acceptance Criteria:** Timestamp check-in berasal dari server dan record tersimpan satu kali.

## 5.8 Check-out

**Tujuan:** Mencatat waktu akhir kehadiran.

**Aktor:** Mahasiswa.

**Precondition:** Attendance check-in tersedia dan belum check-out.

**Main Flow:** User memilih check-out → backend memeriksa record → backend menyimpan `checked_out_at`.

**Alternative Flow:** Belum check-in → ditolak; telah check-out → duplicate operation ditolak.

**Output:** Attendance lengkap.

**Business Rule:** Check-out tidak dapat dilakukan sebelum check-in.

**Acceptance Criteria:** Record memiliki timestamp check-out yang berasal dari server.

## 5.9 Asset Management

**Tujuan:** Mengelola master data setiap unit perangkat.

**Aktor:** Admin.

**Precondition:** Admin login.

**Main Flow:** Admin membuat atau memperbarui data unit → API validasi asset code → data tersimpan.

**Output:** Asset master record.

**Business Rule:**

* `asset_code` unik.
* Asset dapat memiliki serial number unik jika tersedia.
* Status asset terkontrol.

**Acceptance Criteria:** Asset code tidak dapat diduplikasi dan detail unit dapat ditampilkan.

## 5.10 Asset Allocation

**Tujuan:** Mencatat unit perangkat yang diberikan kepada mahasiswa pada session tertentu.

**Aktor:** Admin.

**Precondition:**

* User valid.
* Session valid.
* Asset tersedia.
* Attendance user telah aktif jika aturan tersebut diterapkan sebagai prerequisite.

**Main Flow:**

1. Admin memilih session.
2. Admin memilih mahasiswa.
3. Admin memilih asset.
4. Backend memeriksa availability.
5. Backend membuat allocation.
6. Status asset menjadi `ALLOCATED`.

**Alternative Flow:** Asset telah dialokasikan → request ditolak.

**Output:** Allocation record.

**Business Rule:**

* Satu asset tidak boleh memiliki dua allocation aktif.
* Allocation harus memiliki session dan user.
* Perubahan status asset dan insert allocation diproses dalam satu transaction.

**Acceptance Criteria:** Double allocation ditolak dan allocation tersimpan.

## 5.11 Asset History

**Tujuan:** Menelusuri histori penggunaan dan kondisi suatu aset.

**Aktor:** Admin, Mahasiswa untuk data yang berkaitan dengan dirinya.

**Deskripsi:** Histori dibentuk dari data allocation, return, maintenance, dan audit log yang relevan.

**Precondition:** Asset tersedia di database.

**Output:** Timeline atau daftar histori.

**Business Rule:** Histori tidak dihapus secara fisik untuk record operasional yang telah terjadi.

**Acceptance Criteria:** Sistem mampu menampilkan hubungan asset dengan user/session/allocation/maintenance.

## 5.12 Maintenance

**Tujuan:** Mencatat kerusakan dan proses perbaikan aset.

**Aktor:** Admin.

**Precondition:** Asset tersedia.

**Main Flow:**

1. Admin membuat laporan maintenance.
2. Admin mengisi issue dan description.
3. Status asset menjadi `MAINTENANCE`.
4. Maintenance diproses.
5. Admin mengisi hasil perbaikan.
6. Maintenance ditutup.
7. Asset dikembalikan menjadi `AVAILABLE` atau status lain sesuai kondisi.

**Alternative Flow:** Asset ditemukan tidak layak digunakan → status menjadi `RETIRED`.

**Output:** Maintenance record.

**Business Rule:** Asset berstatus maintenance tidak dapat dialokasikan.

**Acceptance Criteria:** Asset maintenance tidak dapat dipilih dalam allocation aktif.

## 5.13 Reporting

**Tujuan:** Menyediakan laporan operasional.

**Aktor:** Admin.

**Jenis laporan:**

* Attendance report.
* Asset report.
* Allocation report.
* Maintenance report.

**Precondition:** Admin login.

**Main Flow:** Admin memilih jenis laporan → memasukkan filter → API mengambil data → sistem menampilkan hasil → data dapat diekspor dalam format CSV jika fitur Should Have telah selesai.

**Output:** Data laporan.

**Business Rule:** Filter hanya menghasilkan data sesuai hak akses admin.

**Acceptance Criteria:** Laporan dapat difilter berdasarkan periode dan entity yang relevan.

## 5.14 Dashboard

**Tujuan:** Memberikan ringkasan status sistem.

**Aktor:** Mahasiswa, Admin.

**Deskripsi:**

* Mahasiswa melihat status session, presensi, dan allocation.
* Admin melihat jumlah session aktif, attendance, aset available, allocated, maintenance, dan ringkasan aktivitas.

**Output:** Dashboard teragregasi.

**Business Rule:** Dashboard mengambil data dari endpoint backend dan tidak menghitung data bisnis secara langsung di frontend.

**Acceptance Criteria:** Dashboard menampilkan angka yang konsisten dengan data backend.

---

# ========================================

# PART 6 — FACE ENROLLMENT MULTI-ANGLE DAN FACE VERIFICATION 1:1

# ========================================

## 6.1 Pengalaman yang Ditargetkan

Pengalaman pengguna terinspirasi dari alur pendaftaran wajah pada perangkat: pengguna mengikuti panduan kamera untuk mendaftarkan wajah dari beberapa sudut, kemudian saat presensi cukup membuka kamera untuk diperiksa kembali. Untuk proyek ini istilah yang digunakan adalah **multi-angle face enrollment dan 1:1 face verification**, bukan klaim bahwa sistem memiliki keamanan setara Face ID pada ponsel.

Implementasi web memakai kamera browser 2D dan model pretrained. Face ID pada perangkat tertentu memanfaatkan kamera/sensor kedalaman, pemrosesan aman, dan mekanisme anti-spoofing khusus; karena itu implementasi web ini tidak boleh diklaim mempunyai tingkat keamanan yang sama. Dokumentasi Apple menerangkan penggunaan data kedalaman TrueDepth dan Secure Enclave pada Face ID: https://support.apple.com/id-id/102381.

## 6.2 Konsep Face Verification 1:1

Sistem terlebih dahulu mengetahui identitas dari akun yang login. Sistem hanya menjawab apakah wajah live yang sedang diuji cocok dengan enrollment set milik akun tersebut. Sistem **tidak** mencari identitas wajah dari seluruh pengguna (1:N).

Face model mengubah citra wajah menjadi embedding, yaitu representasi numerik. Embedding live dibandingkan dengan sampel enrollment yang tersimpan untuk user tersebut menggunakan metrik distance/similarity yang sama.

## 6.3 Enrollment Multi-Angle

Enrollment dilakukan melalui Web Mahasiswa setelah login dan persetujuan pemrosesan biometrik.

Pose MVP:

1. `FRONT`: wajah menghadap kamera.
2. `LEFT`: kepala menoleh sedikit ke kiri, dengan wajah masih terdeteksi jelas.
3. `RIGHT`: kepala menoleh sedikit ke kanan, dengan wajah masih terdeteksi jelas.

UI memberi instruksi visual per pose. Untuk setiap pose, sistem dapat mengambil beberapa frame sementara dan memilih satu frame yang memenuhi kualitas minimum. Jika pose tidak valid, pengguna mengulang pose tersebut. Enrollment baru berstatus aktif setelah ketiga sampel lolos validasi dan berhasil disimpan.

Alur:

```text
Login + persetujuan biometrik
  ↓
Buka Face Enrollment
  ↓
Kamera aktif dan preview
  ↓
Capture pose FRONT
  ↓
Validasi kualitas + satu wajah
  ↓
Capture pose LEFT
  ↓
Validasi kualitas + satu wajah
  ↓
Capture pose RIGHT
  ↓
Validasi kualitas + satu wajah
  ↓
Ekstraksi 3 embedding (pretrained model)
  ↓
Simpan enrollment set + sampel terenkripsi
  ↓
Buang citra mentah dari memori
  ↓
Enrollment berhasil
```

Satu user memiliki satu enrollment set aktif pada MVP, tetapi set tersebut berisi beberapa sampel. Sistem tidak menyimpan foto wajah mentah secara permanen.

## 6.4 Data yang Disimpan

`face_templates` menjadi metadata enrollment set; sampel embedding disimpan pada tabel anak `face_template_samples`.

Metadata enrollment set minimal:

* `id`, `user_id`, `model_name`, `model_version`;
* `status` (`ACTIVE` atau `REVOKED`);
* `enrolled_at`, `updated_at`;
* metadata persetujuan/versi pemberitahuan seperlunya, tanpa memasukkan data sensitif berlebihan.

Sampel minimal:

* `id`, `face_template_id`;
* `pose` (`FRONT`, `LEFT`, `RIGHT`);
* `embedding_ciphertext`;
* `quality_status` dan `created_at`.

Kunci enkripsi disimpan sebagai secret backend dan tidak boleh disimpan dalam database atau repository. Embedding tidak sama dengan foto, tetapi tetap merupakan data biometrik sensitif sehingga harus dilindungi dengan kontrol akses ketat.

## 6.5 Verifikasi Saat Presensi

Pada penggunaan normal, pengguna tidak perlu mengulang tiga pose setiap kali presensi. Alur MVP menangkap frame frontal atau burst singkat dari kamera, lalu membandingkannya dengan enrollment set milik akun yang login.

```text
Login
  ↓
Pilih session aktif
  ↓
Camera preview
  ↓
Capture frame frontal / burst singkat
  ↓
Validasi citra + tepat satu wajah
  ↓
Buat embedding live
  ↓
Ambil sampel FRONT/LEFT/RIGHT milik user login
  ↓
Hitung distance/similarity
  ↓
Aturan keputusan dan threshold hasil evaluasi
  ↓
Verified / Verification failed
  ↓
Jika verified → lanjut check-in
```

Aturan pembandingan yang menggabungkan beberapa sampel enrollment harus ditetapkan melalui pengujian. Jangan menetapkan threshold secara arbitrer atau sekadar menggunakan kecocokan tertinggi tanpa mengukur dampak false acceptance.

## 6.6 Face Model dan Versi

Sistem memakai pretrained model; tidak melatih model dari nol. DeepFace digunakan sebagai orchestration layer dan ArcFace menjadi kandidat awal yang harus divalidasi. Model final ditetapkan setelah pengujian kompatibilitas, lisensi, akurasi lokal, serta waktu inference.

Semua sampel dalam satu enrollment set harus dibuat dengan versi model dan preprocessing yang kompatibel. Perubahan model/detector/preprocessing yang mengubah dimensi atau distribusi embedding memerlukan strategi migrasi atau enrollment ulang.

## 6.7 Validasi Citra

Validasi minimum pada tiap pose dan verifikasi:

1. File benar-benar berupa citra yang dapat dibaca.
2. MIME type dan ukuran file sesuai batas.
3. Tepat satu wajah terdeteksi.
4. Area wajah cukup besar dan tidak terlalu blur.
5. Pencahayaan cukup dan tidak ekstrem.
6. Pose sesuai dengan pose yang sedang diminta untuk enrollment.
7. Proses detection/alignment/embedding selesai tanpa error.

Jika syarat tidak terpenuhi, jangan menyimpan sampel enrollment dan jangan membuat attendance. Tampilkan pesan yang mudah dipahami, tanpa membocorkan detail model internal.

## 6.8 Threshold dan Evaluasi

Threshold bergantung pada model, detector, preprocessing, kamera, dan data pengujian. Threshold final belum dapat ditentukan hanya dari PRD.

Pengujian minimum meliputi pasangan wajah orang yang sama dan orang berbeda, dengan kondisi kamera/pencahayaan bervariasi. Evaluasi false acceptance dan false rejection, tentukan aturan keputusan, lalu catat versi threshold/configuration yang digunakan. Hasil pengujian harus menjadi dasar acceptance criteria.

## 6.9 Spoofing dan Liveness

Kamera web 2D dan pencocokan embedding saja tidak menjamin bahwa input berasal dari wajah manusia yang hadir langsung. Foto pada layar atau media lain masih merupakan risiko. Multi-angle enrollment meningkatkan variasi sampel rujukan tetapi **tidak dengan sendirinya menjadi liveness detection**.

Untuk MVP, anti-spoofing tingkat lanjut tetap di luar scope. Basic active liveness challenge, misalnya meminta gerakan kepala secara acak pada sesi tertentu, dapat dipertimbangkan sebagai peningkatan terpisah setelah pipeline verifikasi stabil. Jika belum diimplementasikan dan diuji, dokumentasikan batasan ini secara jelas.

## 6.10 Privasi dan Retensi

* Beri penjelasan dan minta persetujuan pengguna sebelum enrollment.
* Simpan hanya embedding terenkripsi dan metadata minimum yang diperlukan.
* Jangan menyimpan raw image di PostgreSQL, Supabase Storage, log, atau folder permanen pada MVP.
* Raw image diproses sementara di memori dan dibuang setelah ekstraksi embedding.
* Jangan pernah mengirim embedding ke frontend atau menampilkannya kepada admin.
* Batasi akses ke tabel biometrik hanya untuk backend/face service yang membutuhkannya.
* Kunci enkripsi berada di environment secret/secret manager, bukan di database.
* Sediakan proses revoke dan re-enrollment; definisikan penghapusan seluruh sampel saat enrollment dicabut atau akun dihapus.
* Kebijakan persetujuan, retensi dan penghapusan tetap perlu divalidasi dengan kebijakan institusi.

## 6.11 Keterbatasan Sistem

Sistem dipengaruhi kualitas kamera, pencahayaan, pose, perubahan penampilan, perangkat, dan performa model. Karena memakai kamera browser 2D, sistem bukan pengganti autentikasi biometrik perangkat yang memiliki sensor kedalaman dan hardware keamanan khusus. Jangan menyatakan sistem mampu mencegah seluruh bentuk titip absen atau spoofing.

## 6.12 Acceptance Criteria untuk Modul Wajah

* Enrollment mengharuskan pose FRONT, LEFT, dan RIGHT valid.
* Set enrollment yang belum lengkap tidak dapat digunakan untuk presensi.
* Verifikasi hanya membandingkan wajah live dengan enrollment set user yang login.
* Enrollment menyimpan tiga sampel embedding terenkripsi beserta pose dan versi model.
* Raw image tidak tersimpan permanen.
* Threshold dan aturan keputusan terdokumentasi berdasarkan pengujian.
* Jika verifikasi gagal/error, attendance tidak dibuat.
* Sistem menyediakan status enrollment, revoke, dan re-enrollment yang terkontrol.

---

# ========================================

# PART 7 — USER FLOW DAN BUSINESS PROCESS

# ========================================

## 7.1 Flow Login

```text
Mulai
  ↓
Buka Halaman Login
  ↓
Input Email/Username + Password
  ↓
POST /api/v1/auth/login
  ↓
Validasi Request
  ↓
Validasi User
  ↓
Verifikasi Password
  ↓
Akun Aktif?
  ├── Tidak → Error 403
  ↓ Ya
Buat Access Token
  ↓
Set Secure HttpOnly Cookie
  ↓
Response Profile + Role
  ↓
Dashboard
```

## 7.2 Flow Face Enrollment Multi-Angle

```text
Login + persetujuan biometrik
  ↓
Buka Face Enrollment
  ↓
Minta izin kamera
  ↓
Capture FRONT → validasi
  ↓ valid
Capture LEFT → validasi
  ↓ valid
Capture RIGHT → validasi
  ↓ valid
Ekstraksi embedding untuk tiap pose
  ↓
Simpan 3 sampel terenkripsi + metadata model
  ↓
Buang raw image dari memori
  ↓
Enrollment aktif
```

Jika salah satu pose gagal divalidasi, ulangi pose tersebut. Jika proses berhenti sebelum ketiga sampel tersimpan, enrollment tidak boleh berstatus aktif.

## 7.3 Flow Presensi dengan Face Verification

```text
Login
  ↓
Pilih Session
  ↓
Session Valid dan Enrollment Aktif?
  ├── Tidak → Error / arahkan enrollment
  ↓ Ya
Buka kamera + capture frame frontal/burst singkat
  ↓
Validasi citra dan satu wajah
  ↓
Ekstraksi embedding live
  ↓
Ambil 3 sampel enrollment user login
  ↓
Similarity/distance + aturan keputusan hasil evaluasi
  ↓
Verified?
  ├── Tidak → Verification failed, attendance tidak dibuat
  ↓ Ya
Cek duplicate attendance
  ↓
Buat attendance dalam transaksi + waktu server
  ↓
Attendance Success
```

## 7.4 Flow Check-in

```text
Login
  ↓
Session Validation
  ↓
Face Verification
  ↓
Duplicate Check
  ↓
INSERT Attendance
  ↓
checked_in_at = Server Time
  ↓
Response Success
```

## 7.5 Flow Check-out

```text
Login
  ↓
Ambil Attendance Aktif
  ↓
Ada Check-in?
  ├── Tidak → Error
  ↓ Ya
Sudah Check-out?
  ├── Ya → Error Duplicate
  ↓ Tidak
Set checked_out_at
  ↓
Response Success
```

## 7.6 Flow Asset Allocation

```text
Admin Login
  ↓
Pilih Session
  ↓
Pilih Mahasiswa
  ↓
Pilih Asset
  ↓
Check Asset Status
  ↓
Available?
  ├── Tidak → Allocation Failed
  ↓ Ya
Check Active Allocation
  ↓
Create Allocation
  ↓
Update Asset = ALLOCATED
  ↓
Commit Transaction
  ↓
Allocation Success
```

## 7.7 Flow Asset Return

```text
Admin Login
  ↓
Pilih Active Allocation
  ↓
Input Condition After Return
  ↓
Update Allocation
  ↓
returned_at = Server Time
  ↓
Tentukan Asset Status
  ├── Good → AVAILABLE
  └── Damaged → MAINTENANCE
  ↓
Commit Transaction
  ↓
Return Success
```

## 7.8 Flow Maintenance

```text
Admin Login
  ↓
Pilih Asset
  ↓
Create Maintenance
  ↓
Asset = MAINTENANCE
  ↓
Diagnosis / Repair
  ↓
Update Maintenance
  ↓
Maintenance Completed?
  ├── Belum → Tetap MAINTENANCE
  └── Ya
       ↓
Condition Evaluation
       ↓
AVAILABLE / RETIRED
       ↓
Save
```

## 7.9 Flow Reporting

```text
Admin Login
  ↓
Pilih Report
  ↓
Set Filter
  ↓
GET Report API
  ↓
Backend Query
  ↓
Validation Access
  ↓
Generate Result
  ↓
Display Report
  ↓
Optional CSV Export
```

---

# ========================================

# PART 8 — ARSITEKTUR SISTEM

# ========================================

## 8.1 Presentation Layer

Presentation Layer terdiri dari:

1. Web Mahasiswa/User.
2. Web Admin/Laboran/Tata Usaha.

Frontend menggunakan React 19 dengan TypeScript, Vite 8 sebagai build tool, dan Tailwind CSS v4 sebagai utility-first CSS framework. Frontend tidak memiliki kredensial database.

Frontend hanya berkomunikasi dengan endpoint REST API melalui HTTP/HTTPS.

## 8.2 API Layer

API Layer merupakan pintu masuk seluruh request aplikasi.

Tanggung jawab:

* routing endpoint;
* authentication middleware;
* authorization;
* request validation;
* response formatting;
* status code;
* exception handling;
* rate limiting;
* audit context.

## 8.3 Business Logic Layer

Business Logic Layer menangani aturan domain, misalnya:

* session harus aktif;
* attendance hanya satu kali;
* asset tidak boleh double allocation;
* asset maintenance tidak dapat dialokasikan;
* check-out harus memiliki check-in;
* verification harus menggunakan user yang sedang login.

## 8.4 Face Verification Layer

Face Verification Layer menangani:

* image preprocessing;
* face detection;
* face alignment;
* embedding extraction;
* normalization;
* similarity/distance calculation;
* threshold evaluation;
* verification result.

Layer ini menggunakan pretrained model.

## 8.5 Data Access Layer

Data Access Layer mengelola komunikasi backend Go dengan PostgreSQL (Supabase).

Tanggung jawab:

* query;
* insert;
* update;
* transaction;
* constraint handling;
* mapping database record ke response schema.

## 8.6 Database Layer

Database menggunakan PostgreSQL (Supabase) sebagai penyimpanan:

* user;
* role;
* face template;
* session;
* attendance;
* asset;
* allocation;
* maintenance;
* audit log.

## 8.7 External Service

Tidak ada external functional API yang wajib digunakan untuk proses bisnis utama.

Komponen eksternal hanya berupa layanan deployment/infrastruktur dan package/model yang dipasang pada backend.

---

# ========================================

# PART 9 — BACKEND DESIGN

# ========================================

## 9.1 Authentication Service

**Fungsi:** Login, logout, token validation, dan password verification.

**Input:** identifier, password, token.

**Proses:** validation → user lookup → password verification → token generation.

**Output:** authentication result, user identity, role.

**Dependency:** User Repository, Password Hashing Library, JWT Library.

## 9.2 Authorization

**Fungsi:** Memastikan user memiliki role dan permission yang sesuai.

**Input:** authenticated user, endpoint requirement.

**Proses:** decode token → identify role → check permission.

**Output:** allow/deny.

**Dependency:** Authentication Service.

## 9.3 User Service

**Fungsi:** CRUD user dan user status.

**Input:** user payload.

**Proses:** validation → duplicate check → persistence.

**Output:** sanitized user data.

**Dependency:** User Repository, Audit Log.

## 9.4 Face Verification Service

**Fungsi:** Enrollment dan verification.

**Input:** image, user_id, model configuration.

**Proses:** face detection → alignment → embedding → comparison.

**Output:** verification status dan metadata.

**Dependency:** DeepFace/model runtime, Face Template Repository.

## 9.5 Attendance Service

**Fungsi:** Check-in, check-out, histori.

**Input:** user_id, session_id, verification result.

**Proses:** session validation → duplicate check → insert/update.

**Output:** attendance record.

**Dependency:** Session Service, Face Verification Service, Attendance Repository.

## 9.6 Session Service

**Fungsi:** CRUD session dan status session.

**Input:** session data.

**Proses:** validation → persistence → status calculation.

**Output:** session record.

**Dependency:** Session Repository, Audit Log.

## 9.7 Asset Service

**Fungsi:** CRUD asset dan status asset.

**Input:** asset data.

**Proses:** validation → persistence → status transition.

**Output:** asset record.

**Dependency:** Asset Repository.

## 9.8 Allocation Service

**Fungsi:** Allocation dan return.

**Input:** asset_id, user_id, session_id, conditions.

**Proses:** availability check → transaction → allocation → asset status update.

**Output:** allocation record.

**Dependency:** Asset Service, User Service, Session Service, Allocation Repository.

## 9.9 Maintenance Service

**Fungsi:** Menangani laporan kerusakan dan proses maintenance.

**Input:** asset_id, issue, description, status, resolution.

**Proses:** create maintenance → update asset status → close maintenance → restore status.

**Output:** maintenance record.

**Dependency:** Asset Service, Maintenance Repository, Audit Log.

## 9.10 Reporting Service

**Fungsi:** Membentuk attendance report, asset report, allocation report, maintenance report.

**Input:** filter parameter.

**Proses:** authorization → query → aggregation → formatting.

**Output:** JSON report dan CSV jika tersedia.

**Dependency:** seluruh repository terkait.

## 9.11 Audit Log

**Fungsi:** mencatat event penting.

**Event minimal:**

* login success/failure;
* user create/update;
* face enrollment;
* face verification result;
* attendance create;
* asset create/update/delete;
* allocation;
* return;
* maintenance create/update;
* perubahan status penting.

**Ketentuan:** Audit log tidak menyimpan password, token, atau raw image.

---

# ========================================

# PART 10 — DATABASE DESIGN

# ========================================

## 10.1 Keputusan Struktur Database

Database final menggunakan sembilan tabel inti bisnis dan satu tabel tambahan untuk sampel embedding multi-angle:

1. `roles`
2. `users`
3. `face_templates`
4. `lab_sessions`
5. `attendance`
6. `assets`
7. `asset_allocations`
8. `maintenance`
9. `audit_logs`
10. `face_template_samples`

Tabel `face_template_samples` menyimpan tiga sampel embedding per enrollment set (`FRONT`, `LEFT`, `RIGHT`); tabel `face_templates` menyimpan metadata enrollment set. Tabel `maintenance_history` tidak dipisahkan dalam MVP karena histori maintenance dapat direpresentasikan sebagai record append/update pada tabel `maintenance` dengan timestamp dan status. Pemisahan `maintenance_history` tidak diperlukan selama kebutuhan hanya mencatat satu alur maintenance per kejadian.

## 10.2 Tabel `roles`

**Tujuan:** Menyimpan role aplikasi.

| Field      | Tipe         | Constraint       |
| ---------- | ------------ | ---------------- |
| id         | UUID         | PK               |
| code       | VARCHAR(30)  | UNIQUE, NOT NULL |
| name       | VARCHAR(100) | NOT NULL         |
| created_at | TIMESTAMP    | NOT NULL         |

Nilai utama:

* `STUDENT`
* `ADMIN`

## 10.3 Tabel `users`

| Field         | Tipe         | Constraint            |
| ------------- | ------------ | --------------------- |
| id            | UUID         | PK                    |
| role_id       | UUID         | FK roles.id, NOT NULL |
| identifier    | VARCHAR(100) | UNIQUE, NOT NULL      |
| email         | VARCHAR(150) | UNIQUE, NOT NULL      |
| full_name     | VARCHAR(150) | NOT NULL              |
| password_hash | TEXT         | NOT NULL              |
| status        | VARCHAR(20)  | NOT NULL              |
| created_at    | TIMESTAMP    | NOT NULL              |
| updated_at    | TIMESTAMP    | NOT NULL              |

Constraint:

* `status IN ('ACTIVE','INACTIVE')`.

## 10.4 Tabel `face_templates`

Tabel ini menyimpan metadata satu enrollment set untuk satu user. Data embedding dipisahkan ke tabel anak agar beberapa pose dapat disimpan secara terstruktur.

| Field       | Tipe         | Constraint                      |
| ----------- | ------------ | ------------------------------- |
| id          | UUID         | PK                              |
| user_id     | UUID         | FK users.id, NOT NULL           |
| model_name  | VARCHAR(100) | NOT NULL                        |
| model_version | VARCHAR(100) | NOT NULL                     |
| consent_version | VARCHAR(50) | NOT NULL                     |
| consented_at | TIMESTAMP    | NOT NULL                      |
| status      | VARCHAR(20)  | NOT NULL                        |
| enrolled_at | TIMESTAMP    | NOT NULL                        |
| updated_at  | TIMESTAMP    | NOT NULL                        |

Constraint:

* Maksimal satu enrollment set `ACTIVE` per user pada MVP; enforce dengan partial unique index pada `user_id` untuk baris `status = 'ACTIVE'`.
* `status IN ('ACTIVE','REVOKED')`.
* `consent_version` dan `consented_at` mencatat versi pemberitahuan dan waktu persetujuan.
* Enrollment baru hanya diaktifkan setelah seluruh pose wajib valid dan seluruh sampel berhasil disimpan.

## 10.4.1 Tabel `face_template_samples`

| Field                 | Tipe         | Constraint                                  |
| --------------------- | ------------ | ------------------------------------------- |
| id                    | UUID         | PK                                          |
| face_template_id      | UUID         | FK face_templates.id, NOT NULL              |
| pose                  | VARCHAR(20)  | NOT NULL                                    |
| embedding_ciphertext  | BYTEA        | NOT NULL                                    |
| quality_status        | VARCHAR(20)  | NOT NULL                                    |
| created_at            | TIMESTAMP    | NOT NULL                                    |

Constraint:

* `pose IN ('FRONT','LEFT','RIGHT')`.
* `UNIQUE (face_template_id, pose)` untuk enrollment set MVP.
* Sampel harus dibuat oleh model/preprocessing yang kompatibel dengan metadata enrollment set.
* Tidak ada raw image column.
* Enrollment set dianggap lengkap bila memiliki sampel valid untuk ketiga pose.

## 10.5 Tabel `lab_sessions`

| Field        | Tipe         | Constraint            |
| ------------ | ------------ | --------------------- |
| id           | UUID         | PK                    |
| code         | VARCHAR(50)  | UNIQUE, NOT NULL      |
| title        | VARCHAR(150) | NOT NULL              |
| course_code  | VARCHAR(50)  | NULL                  |
| location     | VARCHAR(150) | NOT NULL              |
| session_date | DATE         | NOT NULL              |
| start_at     | TIMESTAMP    | NOT NULL              |
| end_at       | TIMESTAMP    | NOT NULL              |
| status       | VARCHAR(20)  | NOT NULL              |
| created_by   | UUID         | FK users.id, NOT NULL |
| created_at   | TIMESTAMP    | NOT NULL              |
| updated_at   | TIMESTAMP    | NOT NULL              |

Constraint:

* `end_at > start_at`;
* `status IN ('DRAFT','ACTIVE','CLOSED','CANCELLED')`.

## 10.6 Tabel `attendance`

| Field                      | Tipe         | Constraint                   |
| -------------------------- | ------------ | ---------------------------- |
| id                         | UUID         | PK                           |
| user_id                    | UUID         | FK users.id, NOT NULL        |
| session_id                 | UUID         | FK lab_sessions.id, NOT NULL |
| checked_in_at              | TIMESTAMP    | NOT NULL                     |
| checked_out_at             | TIMESTAMP    | NULL                         |
| verification_status        | VARCHAR(20)  | NOT NULL                     |
| verification_score         | DECIMAL(8,6) | NULL                         |
| verification_model_version | VARCHAR(100) | NULL                         |
| created_at                 | TIMESTAMP    | NOT NULL                     |
| updated_at                 | TIMESTAMP    | NOT NULL                     |

Constraint:

* UNIQUE `(user_id, session_id)`;
* `checked_out_at IS NULL OR checked_out_at >= checked_in_at`.

## 10.7 Tabel `assets`

| Field            | Tipe         | Constraint       |
| ---------------- | ------------ | ---------------- |
| id               | UUID         | PK               |
| asset_code       | VARCHAR(50)  | UNIQUE, NOT NULL |
| asset_name       | VARCHAR(150) | NOT NULL         |
| category         | VARCHAR(50)  | NOT NULL         |
| brand            | VARCHAR(100) | NULL             |
| model            | VARCHAR(100) | NULL             |
| serial_number    | VARCHAR(100) | UNIQUE, NULL     |
| specification    | JSONB        | NULL             |
| location         | VARCHAR(150) | NOT NULL         |
| status           | VARCHAR(30)  | NOT NULL         |
| acquisition_year | SMALLINT     | NULL             |
| notes            | TEXT         | NULL             |
| created_at       | TIMESTAMP    | NOT NULL         |
| updated_at       | TIMESTAMP    | NOT NULL         |

Status:

* `AVAILABLE`
* `ALLOCATED`
* `MAINTENANCE`
* `RETIRED`

## 10.8 Tabel `asset_allocations`

| Field            | Tipe        | Constraint                   |
| ---------------- | ----------- | ---------------------------- |
| id               | UUID        | PK                           |
| asset_id         | UUID        | FK assets.id, NOT NULL       |
| user_id          | UUID        | FK users.id, NOT NULL        |
| session_id       | UUID        | FK lab_sessions.id, NOT NULL |
| allocated_at     | TIMESTAMP   | NOT NULL                     |
| returned_at      | TIMESTAMP   | NULL                         |
| condition_before | VARCHAR(30) | NULL                         |
| condition_after  | VARCHAR(30) | NULL                         |
| notes            | TEXT        | NULL                         |
| created_at       | TIMESTAMP   | NOT NULL                     |
| updated_at       | TIMESTAMP   | NOT NULL                     |

Constraint:

* allocation aktif tidak boleh lebih dari satu untuk asset yang sama;
* `returned_at >= allocated_at` jika tidak null.

## 10.9 Tabel `maintenance`

| Field        | Tipe          | Constraint             |
| ------------ | ------------- | ---------------------- |
| id           | UUID          | PK                     |
| asset_id     | UUID          | FK assets.id, NOT NULL |
| reported_by  | UUID          | FK users.id, NOT NULL  |
| issue_type   | VARCHAR(50)   | NOT NULL               |
| description  | TEXT          | NOT NULL               |
| priority     | VARCHAR(20)   | NOT NULL               |
| status       | VARCHAR(20)   | NOT NULL               |
| reported_at  | TIMESTAMP     | NOT NULL               |
| started_at   | TIMESTAMP     | NULL                   |
| completed_at | TIMESTAMP     | NULL                   |
| resolution   | TEXT          | NULL                   |
| cost         | NUMERIC(14,2) | NULL                   |
| created_at   | TIMESTAMP     | NOT NULL               |
| updated_at   | TIMESTAMP     | NOT NULL               |

Status:

* `OPEN`
* `IN_PROGRESS`
* `COMPLETED`
* `CANCELLED`

## 10.10 Tabel `audit_logs`

| Field       | Tipe         | Constraint        |
| ----------- | ------------ | ----------------- |
| id          | UUID         | PK                |
| user_id     | UUID         | FK users.id, NULL |
| action      | VARCHAR(100) | NOT NULL          |
| entity_type | VARCHAR(50)  | NOT NULL          |
| entity_id   | UUID         | NULL              |
| status      | VARCHAR(20)  | NOT NULL          |
| ip_address  | INET         | NULL              |
| user_agent  | TEXT         | NULL              |
| metadata    | JSONB        | NULL              |
| created_at  | TIMESTAMP    | NOT NULL          |

`metadata` tidak boleh berisi password, token, raw biometric image, atau embedding plaintext.

## 10.11 Relasi Utama

```text
roles
  │
  └────< users
           │
           ├────< face_templates ────< face_template_samples
           │
           ├────< attendance >──── lab_sessions
           │
           ├────< asset_allocations >──── assets
           │
           └────< audit_logs

lab_sessions
  │
  ├────< attendance
  └────< asset_allocations

assets
  │
  ├────< asset_allocations
  └────< maintenance

users
  │
  ├────< lab_sessions (created_by)
  ├────< maintenance (reported_by)
  └────< audit_logs
```

## 10.12 ERD Teks

```text
ROLE
 └──< USER
       ├──< FACE_TEMPLATE
       │      └──< FACE_TEMPLATE_SAMPLE (FRONT / LEFT / RIGHT)
       ├──< ATTENDANCE >── LAB_SESSION
       ├──< ASSET_ALLOCATION >── ASSET
       ├──< AUDIT_LOG
       └──< MAINTENANCE >── ASSET

LAB_SESSION
 ├──< ATTENDANCE
 └──< ASSET_ALLOCATION

ASSET
 ├──< ASSET_ALLOCATION
 └──< MAINTENANCE
```

---

# ========================================

# PART 11 — REST API / WEB SERVICE

# ========================================

## 11.1 Authentication API

| Method | Endpoint              | Role  | Auth | Request              | Response             | Validation         | Status | Error              |
| ------ | --------------------- | ----- | ---- | -------------------- | -------------------- | ------------------ | ------ | ------------------ |
| POST   | `/api/v1/auth/login`  | Semua | No   | identifier, password | token/session + user | credential         | 200    | 401, 403, 422, 429 |
| POST   | `/api/v1/auth/logout` | Semua | Yes  | Cookie/token         | success              | authenticated user | 200    | 401                |
| GET    | `/api/v1/auth/me`     | Semua | Yes  | none                 | profile              | token valid        | 200    | 401                |

## 11.2 User API

| Method | Endpoint                    | Role  | Auth | Request        | Response       | Validation              | Status | Error       |
| ------ | --------------------------- | ----- | ---- | -------------- | -------------- | ----------------------- | ------ | ----------- |
| GET    | `/api/v1/users`             | Admin | Yes  | query/filter   | user list      | role admin              | 200    | 401,403     |
| POST   | `/api/v1/users`             | Admin | Yes  | user payload   | created user   | unique identifier/email | 201    | 400,409,422 |
| GET    | `/api/v1/users/{id}`        | Admin | Yes  | path id        | user detail    | valid UUID              | 200    | 404         |
| PUT    | `/api/v1/users/{id}`        | Admin | Yes  | update payload | updated user   | valid fields            | 200    | 404,409,422 |
| PATCH  | `/api/v1/users/{id}/status` | Admin | Yes  | status         | updated status | enum                    | 200    | 404,422     |

## 11.3 Face API

| Method | Endpoint                | Role      | Auth | Request                                                    | Response                          | Validation                                 | Status | Error        |
| ------ | ----------------------- | --------- | ---- | ---------------------------------------------------------- | --------------------------------- | ------------------------------------------ | ------ | ------------ |
| POST   | `/api/v1/face/enroll`   | Mahasiswa | Yes  | multipart: `front_image`, `left_image`, `right_image`     | enrollment status + sample count  | consent, one face/pose, quality, size, MIME | 201    | 400,422,503  |
| GET    | `/api/v1/face/status`   | Mahasiswa | Yes  | none                                                       | enrolled/status/pose completeness | auth + ownership                            | 200    | 401          |
| POST   | `/api/v1/face/verify`   | Mahasiswa | Yes  | multipart: live `image` or short frame burst             | verified/status code              | active enrollment, one face, quality         | 200    | 400,422,429,503 |
| DELETE | `/api/v1/face/template` | Admin     | Yes  | user identifier                                            | revoke status                     | admin access; revoke all samples             | 204    | 403,404      |

API tidak pernah mengembalikan embedding. File input diproses sementara dan tidak disimpan sebagai file permanen. Endpoint `/api/v1/face/verify` hanya mengembalikan hasil verifikasi dan tidak membuat attendance. Endpoint aktual check-in wajib menjalankan verifikasi di server dari citra yang dikirim saat check-in; backend tidak boleh mempercayai flag `verified=true` dari frontend. Enrollment tidak boleh dinyatakan berhasil sebelum ketiga sampel pose tersimpan secara konsisten.

## 11.4 Session API

| Method | Endpoint                       | Role  | Auth | Request         | Response        | Validation        | Status | Error   |
| ------ | ------------------------------ | ----- | ---- | --------------- | --------------- | ----------------- | ------ | ------- |
| POST   | `/api/v1/sessions`             | Admin | Yes  | session payload | created session | datetime          | 201    | 409,422 |
| GET    | `/api/v1/sessions`             | Semua | Yes  | filters         | session list    | date/status       | 200    | 401     |
| GET    | `/api/v1/sessions/{id}`        | Semua | Yes  | path id         | detail          | valid UUID        | 200    | 404     |
| PUT    | `/api/v1/sessions/{id}`        | Admin | Yes  | payload         | updated session | status/time rules | 200    | 404,422 |
| PATCH  | `/api/v1/sessions/{id}/status` | Admin | Yes  | status          | updated status  | enum              | 200    | 404,422 |

## 11.5 Attendance API

| Method | Endpoint                       | Role      | Auth | Request                                      | Response        | Validation               | Status | Error       |
| ------ | ------------------------------ | --------- | ---- | -------------------------------------------- | --------------- | ------------------------ | ------ | ----------- |
| POST   | `/api/v1/attendance/check-in`  | Mahasiswa | Yes  | multipart: `session_id` + live `image`/short frame burst | attendance      | active session, server-side face verification | 201    | 400,409,422,429,503 |
| POST   | `/api/v1/attendance/check-out` | Mahasiswa | Yes  | session_id                                   | attendance      | existing check-in        | 200    | 409,422     |
| GET    | `/api/v1/attendance/history`   | Mahasiswa | Yes  | filters                                      | own history     | ownership                | 200    | 401         |
| GET    | `/api/v1/attendance`           | Admin     | Yes  | filters                                      | attendance list | admin access             | 200    | 403         |
| GET    | `/api/v1/attendance/{id}`      | Admin     | Yes  | path id                                      | detail          | valid UUID               | 200    | 404         |

## 11.6 Asset API

| Method | Endpoint                      | Role  | Auth | Request       | Response            | Validation             | Status | Error   |
| ------ | ----------------------------- | ----- | ---- | ------------- | ------------------- | ---------------------- | ------ | ------- |
| GET    | `/api/v1/assets`              | Semua | Yes  | filters       | asset list          | filter                 | 200    | 401     |
| POST   | `/api/v1/assets`              | Admin | Yes  | asset payload | created asset       | asset code unique      | 201    | 409,422 |
| GET    | `/api/v1/assets/{id}`         | Admin | Yes  | path id       | asset detail        | UUID                   | 200    | 404     |
| PUT    | `/api/v1/assets/{id}`         | Admin | Yes  | payload       | updated asset       | status rules           | 200    | 404,422 |
| DELETE | `/api/v1/assets/{id}`         | Admin | Yes  | path id       | deleted/deactivated | no active dependencies | 204    | 409,404 |
| GET    | `/api/v1/assets/{id}/history` | Admin | Yes  | path id       | asset timeline      | UUID                   | 200    | 404     |

## 11.7 Allocation API

| Method | Endpoint                          | Role      | Auth | Request              | Response            | Validation        | Status | Error   |
| ------ | --------------------------------- | --------- | ---- | -------------------- | ------------------- | ----------------- | ------ | ------- |
| POST   | `/api/v1/allocations`             | Admin     | Yes  | asset, user, session | allocation          | asset available   | 201    | 409,422 |
| GET    | `/api/v1/allocations`             | Admin     | Yes  | filters              | list                | filters           | 200    | 403     |
| GET    | `/api/v1/allocations/my`          | Mahasiswa | Yes  | filters              | own allocations     | ownership         | 200    | 401     |
| GET    | `/api/v1/allocations/{id}`        | Admin     | Yes  | path id              | detail              | UUID              | 200    | 404     |
| POST   | `/api/v1/allocations/{id}/return` | Admin     | Yes  | condition_after      | returned allocation | active allocation | 200    | 409,422 |

## 11.8 Maintenance API

| Method | Endpoint                   | Role  | Auth | Request                | Response            | Validation       | Status | Error   |
| ------ | -------------------------- | ----- | ---- | ---------------------- | ------------------- | ---------------- | ------ | ------- |
| POST   | `/api/v1/maintenance`      | Admin | Yes  | asset, issue, priority | maintenance         | asset valid      | 201    | 404,422 |
| GET    | `/api/v1/maintenance`      | Admin | Yes  | filters                | list                | filters          | 200    | 403     |
| GET    | `/api/v1/maintenance/{id}` | Admin | Yes  | path id                | detail              | UUID             | 200    | 404     |
| PUT    | `/api/v1/maintenance/{id}` | Admin | Yes  | status/resolution      | updated maintenance | transition rules | 200    | 409,422 |

## 11.9 Report API

| Method | Endpoint                      | Role  | Auth | Request              | Response           | Validation | Status | Error   |
| ------ | ----------------------------- | ----- | ---- | -------------------- | ------------------ | ---------- | ------ | ------- |
| GET    | `/api/v1/reports/attendance`  | Admin | Yes  | date/session filters | attendance report  | date range | 200    | 403,422 |
| GET    | `/api/v1/reports/assets`      | Admin | Yes  | asset filters        | asset report       | filters    | 200    | 403     |
| GET    | `/api/v1/reports/maintenance` | Admin | Yes  | date/status          | maintenance report | filters    | 200    | 403     |
| GET    | `/api/v1/reports/allocations` | Admin | Yes  | session/date         | allocation report  | filters    | 200    | 403     |

## 11.10 Dashboard API

| Method | Endpoint                    | Role      | Auth | Request | Response         | Validation | Status | Error |
| ------ | --------------------------- | --------- | ---- | ------- | ---------------- | ---------- | ------ | ----- |
| GET    | `/api/v1/dashboard/student` | Mahasiswa | Yes  | none    | personal summary | ownership  | 200    | 401   |
| GET    | `/api/v1/dashboard/admin`   | Admin     | Yes  | none    | system summary   | role admin | 200    | 403   |

## 11.11 Aturan API Umum

Format request menggunakan `application/json`, sedangkan endpoint face menggunakan `multipart/form-data`.

Format response sukses:

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {}
}
```

Format response error:

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": [
    {
      "field": "email",
      "code": "INVALID_EMAIL"
    }
  ],
  "request_id": "generated-request-id"
}
```

Semua timestamp menggunakan format ISO 8601 dan ditentukan oleh server untuk event yang bersifat audit/transactional.

---

# ========================================

# PART 12 — API DOCUMENTATION / DOC API

# ========================================

## 12.1 Struktur Dokumentasi API

Dokumentasi API terdiri dari:

1. API Overview.
2. Base URL.
3. Authentication.
4. Authorization.
5. API Tags.
6. Endpoint list.
7. Request schema.
8. Response schema.
9. Error schema.
10. Example request.
11. Example response.
12. HTTP status.
13. Validation rule.
14. Authentication requirement.

## 12.2 Base URL

Development:

`http://localhost:3000`

Production:

`[PERLU VALIDASI]`

Base API:

`/api/v1`

## 12.3 Authentication

Authentication menggunakan access token berbasis JWT yang disimpan melalui cookie `HttpOnly`, `Secure`, dan `SameSite` sesuai arsitektur deployment.

Token berisi minimal:

```json
{
  "sub": "user-id",
  "role": "STUDENT",
  "exp": "expiration"
}
```

## 12.4 Authorization

Authorization menggunakan role:

* `STUDENT`
* `ADMIN`

Setiap endpoint menetapkan role yang diizinkan.

## 12.5 API Tags

Swagger/OpenAPI menggunakan tag:

* Authentication
* Users
* Face
* Sessions
* Attendance
* Assets
* Allocations
* Maintenance
* Reports
* Dashboard

## 12.6 Request Schema

Contoh login:

```json
{
  "identifier": "student@example.com",
  "password": "example-password"
}
```

Contoh session:

```json
{
  "code": "LAB-SESSION-001",
  "title": "Praktikum Pemrograman Web",
  "course_code": "[PERLU VALIDASI]",
  "location": "[PERLU VALIDASI]",
  "session_date": "2026-09-18",
  "start_at": "2026-09-18T08:00:00+07:00",
  "end_at": "2026-09-18T10:00:00+07:00"
}
```

## 12.7 Response Schema

```json
{
  "success": true,
  "message": "Session created",
  "data": {
    "id": "uuid",
    "code": "LAB-SESSION-001",
    "status": "DRAFT"
  }
}
```

## 12.8 Error Schema

```json
{
  "success": false,
  "message": "Asset is not available",
  "errors": [
    {
      "field": "asset_id",
      "code": "ASSET_NOT_AVAILABLE"
    }
  ],
  "request_id": "uuid"
}
```

## 12.9 Example Request

```http
POST /api/v1/attendance/check-in
Content-Type: application/json
Authorization: authenticated session

{
  "session_id": "uuid"
}
```

Face sample dikirim pada endpoint verification:

```http
POST /api/v1/face/verify
Content-Type: multipart/form-data
```

Field:

```text
session_id=uuid
image=<captured-frame>
```

## 12.10 Example Response

```json
{
  "success": true,
  "message": "Face verified",
  "data": {
    "verified": true,
    "verification_id": "uuid"
  }
}
```

Data internal seperti embedding tidak dikembalikan.

## 12.11 Swagger/OpenAPI

Dokumentasi API menggunakan OpenAPI. Backend menyediakan endpoint dokumentasi:

```text
/swagger
/api/v1/openapi.json
```

Dokumentasi menjadi salah satu luaran utama proyek karena menunjukkan seluruh kontrak antara frontend dan backend.

---

# ========================================

# PART 13 — SECURITY

# ========================================

## 13.1 Authentication

Authentication menggunakan JWT access token dengan lifetime terbatas.

Ketentuan:

* Password tidak dikirim kecuali melalui HTTPS.
* Token memiliki expiry.
* Token tidak dimasukkan ke URL.
* Logout menghapus cookie/session token sisi browser.

## 13.2 Authorization

Authorization dilakukan pada backend, bukan hanya frontend.

Frontend boleh menyembunyikan tombol, tetapi backend tetap memeriksa role setiap request.

## 13.3 Password Security

Password disimpan menggunakan password hashing adaptif, dengan algoritma yang mendukung Argon2id menjadi pilihan implementasi.

Tidak terdapat endpoint yang mengembalikan password hash kepada frontend.

## 13.4 Token Security

Token:

* tidak disimpan dalam localStorage pada desain utama;
* menggunakan HttpOnly cookie;
* menggunakan Secure pada production;
* menggunakan SameSite yang sesuai;
* memiliki expiry;
* tidak dicatat dalam application log.

## 13.5 API Security

API menerapkan:

* HTTPS pada production;
* CORS terbatas;
* rate limiting untuk login dan face verification;
* request size limit;
* validation;
* authentication middleware;
* authorization middleware;
* structured error response;
* request ID.

## 13.6 Input Validation

Validation dilakukan pada:

* type;
* required field;
* format email;
* enum;
* length;
* UUID;
* tanggal/waktu;
* file MIME type;
* image size;
* asset code;
* duplicate rule.

## 13.7 Database Security

Database:

* tidak diakses langsung frontend;
* memiliki credential terpisah dari source code;
* menggunakan environment variable;
* menggunakan foreign key;
* menggunakan unique constraint;
* menggunakan transaction untuk operasi terkait status;
* backup dilakukan sesuai kemampuan deployment `[PERLU VALIDASI]`.

## 13.8 Biometric Data Security

Data biometric:

* tidak ditampilkan dalam UI;
* tidak disimpan dalam log;
* raw image tidak disimpan permanen pada MVP;
* embedding dilindungi;
* endpoint template dibatasi;
* model version dicatat;
* akses dibatasi backend service.

## 13.9 Access Control

Contoh:

```text
STUDENT
  ├── own profile
  ├── own face template process
  ├── own attendance
  └── own allocation

ADMIN
  ├── users
  ├── sessions
  ├── attendance monitoring
  ├── assets
  ├── allocations
  ├── maintenance
  └── reports
```

## 13.10 Audit Trail

Audit trail mencatat event penting untuk kebutuhan penelusuran.

Contoh:

```text
LOGIN_SUCCESS
LOGIN_FAILED
FACE_ENROLL
FACE_VERIFY_SUCCESS
FACE_VERIFY_FAILED
ATTENDANCE_CREATE
ATTENDANCE_CHECKOUT
ASSET_CREATE
ASSET_UPDATE
ALLOCATION_CREATE
ASSET_RETURN
MAINTENANCE_CREATE
MAINTENANCE_COMPLETE
```

## 13.11 Privacy Consideration

Privacy diperhatikan melalui prinsip minimization, least privilege, dan pembatasan data biometric.

Kebijakan retensi biometric dan data akun harus mengikuti kebijakan institusi:

**[PERLU VALIDASI]**

## 13.12 Error Handling

Error tidak boleh mengungkap:

* password;
* secret;
* token;
* SQL query;
* file path internal;
* embedding;
* detail sistem internal.

Error response bersifat terstruktur dan memiliki request ID untuk troubleshooting.

---

# ========================================

# PART 14 — AGILE DEVELOPMENT

# ========================================

## 14.1 Alasan Menggunakan Agile

Agile digunakan karena proyek memiliki banyak modul yang saling bergantung dan membutuhkan validasi bertahap. Face verification, REST API, database, dan frontend memiliki risiko teknis yang perlu diuji lebih awal.

Pendekatan Agile memungkinkan sistem dibangun sebagai increment sehingga risiko kegagalan dapat diketahui sebelum akhir semester.

## 14.2 Product Backlog

Product backlog berisi kebutuhan:

* authentication;
* user management;
* multi-angle face enrollment (FRONT/LEFT/RIGHT);
* face verification 1:1;
* session;
* attendance;
* assets;
* allocation;
* maintenance;
* reporting;
* dashboard;
* documentation;
* testing;
* deployment.

## 14.3 Sprint

Sprint memiliki durasi tetap sesuai kalender akademik:

**[DURASI SPRINT PERLU VALIDASI]**

Pembagian menggunakan tujuh sprint pengembangan.

## 14.4 Sprint Planning

Setiap awal sprint:

1. meninjau backlog;
2. menentukan sprint goal;
3. memilih user story;
4. memeriksa dependency;
5. membagi pekerjaan dua mahasiswa;
6. menentukan acceptance criteria.

## 14.5 Daily/Progress Monitoring

Monitoring dilakukan secara rutin melalui:

* status task;
* blocker;
* progress backend;
* progress frontend;
* integration status;
* testing status.

## 14.6 Sprint Review

Setiap akhir sprint, increment ditinjau berdasarkan acceptance criteria.

## 14.7 Retrospective

Tim mengevaluasi:

* pekerjaan yang berjalan baik;
* masalah teknis;
* dependency;
* kualitas implementasi;
* scope;
* distribusi beban kerja.

## 14.8 Increment

Increment harus berupa bagian sistem yang berjalan dan dapat diuji, bukan hanya file kode yang belum terintegrasi.

## 14.9 Sprint Plan

### Sprint 1 — Requirement + Architecture

Output:

* PRD;
* use case;
* arsitektur;
* database draft;
* API convention;
* repository structure;
* setup development environment.

### Sprint 2 — Database + Backend + Authentication

Output:

* PostgreSQL;
* migration;
* models;
* repository;
* authentication;
* authorization;
* user management API;
* Swagger initial.

### Sprint 3 — Face Enrollment + Face Verification + Attendance

Output:

* browser camera;
* guided multi-angle face enrollment (FRONT/LEFT/RIGHT);
* multiple encrypted face embeddings per enrollment set;
* verification service;
* attendance;
* session validation;
* verification test awal.

### Sprint 4 — Asset Management + Allocation

Output:

* asset CRUD;
* allocation;
* return;
* asset status;
* transaction;
* asset history.

### Sprint 5 — History + Maintenance

Output:

* maintenance;
* status transitions;
* asset timeline;
* audit log;
* reporting query awal.

### Sprint 6 — Web Mahasiswa + Web Admin + Integration

Output:

* Web Mahasiswa;
* Web Admin;
* API integration;
* end-to-end flow.

### Sprint 7 — Testing + Documentation + Finalization

Output:

* unit test;
* API test;
* integration test;
* face verification testing;
* Swagger;
* tutorial;
* journal draft;
* poster;
* deployment;
* final bug fixing.

---

# ========================================

# PART 15 — PRODUCT BACKLOG

# ========================================

| ID    | Epic           | User Story                                                                                                    | Priority | Acceptance Criteria                            | Dependency         |
| ----- | -------------- | ------------------------------------------------------------------------------------------------------------- | -------- | ---------------------------------------------- | ------------------ |
| PB-01 | Authentication | As a user, I want to login so that I can access the system securely.                                          | Must     | Login valid berhasil dan login invalid ditolak | User DB            |
| PB-02 | Authentication | As a user, I want to logout so that my session is terminated.                                                 | Must     | Cookie/token session dihapus                   | PB-01              |
| PB-03 | User           | As an admin, I want to manage users so that account data remains current.                                     | Must     | CRUD user berjalan                             | PB-01              |
| PB-04 | Face           | As a student, I want to enroll my face so that the system has my verification template.                       | Must     | Template aktif tersimpan                       | PB-01              |
| PB-05 | Face           | As a student, I want to verify my face so that my identity can be checked during attendance.                  | Must     | 1:1 verification menghasilkan status           | PB-04              |
| PB-06 | Session        | As an admin, I want to manage lab sessions so that attendance is tied to a specific activity.                 | Must     | Session CRUD berjalan                          | PB-03              |
| PB-07 | Attendance     | As a student, I want to check in so that my attendance is recorded.                                           | Must     | Attendance tersimpan setelah verification      | PB-05, PB-06       |
| PB-08 | Attendance     | As a student, I want to check out so that my attendance end time is recorded.                                 | Must     | `checked_out_at` tersimpan                     | PB-07              |
| PB-09 | Asset          | As an admin, I want to register assets so that each unit can be tracked individually.                         | Must     | Asset code unique                              | PB-03              |
| PB-10 | Allocation     | As an admin, I want to allocate an asset so that device usage is recorded.                                    | Must     | Allocation aktif tersimpan                     | PB-07, PB-09       |
| PB-11 | Allocation     | As an admin, I want to return an asset so that the device status is updated.                                  | Must     | Return tercatat                                | PB-10              |
| PB-12 | History        | As an admin, I want to view asset history so that I can trace usage.                                          | Must     | Timeline tampil                                | PB-10, PB-11       |
| PB-13 | Maintenance    | As an admin, I want to record maintenance so that damaged assets can be managed.                              | Must     | Maintenance tersimpan                          | PB-09              |
| PB-14 | Dashboard      | As an admin, I want a dashboard so that I can see operational summaries.                                      | Must     | Summary tampil                                 | PB-06–PB-13        |
| PB-15 | Reporting      | As an admin, I want attendance reports so that I can monitor attendance.                                      | Must     | Report dapat difilter                          | PB-07              |
| PB-16 | Reporting      | As an admin, I want asset reports so that I can monitor asset status.                                         | Must     | Asset report tersedia                          | PB-09              |
| PB-17 | Reporting      | As an admin, I want maintenance reports so that I can monitor repair history.                                 | Must     | Maintenance report tersedia                    | PB-13              |
| PB-18 | API Docs       | As a developer, I want API documentation so that frontend and backend integration has a clear contract.       | Must     | Swagger/OpenAPI tersedia                       | API                |
| PB-19 | Security       | As an admin, I want audit records so that important activity can be traced.                                   | Should   | Event penting tercatat                         | Authentication/API |
| PB-20 | Reporting      | As an admin, I want CSV export so that reports can be reused.                                                 | Should   | CSV berhasil diunduh                           | PB-15–PB-17        |
| PB-21 | UI             | As a student, I want a dashboard so that I can see my attendance and allocation status.                       | Must     | Data pribadi tampil                            | PB-07, PB-10       |
| PB-22 | Integration    | As a team, we want two web applications to consume one REST API so that the architecture remains centralized. | Must     | Tidak ada direct DB access dari frontend       | API + Frontend     |
| PB-23 | Testing        | As a developer, I want API and integration tests so that system behavior can be validated.                    | Must     | Test suite tersedia                            | All core modules   |
| PB-24 | Deployment     | As a team, we want the web service to be deployed so that the system can be demonstrated.                     | Must     | Production/demo environment aktif              | Final build        |

---

# ========================================

# PART 16 — PEMBAGIAN TUGAS 2 MAHASISWA

# ========================================

## 16.1 Prinsip Pembagian

Pembagian kerja mempertimbangkan dependency dan beban integrasi. Kedua mahasiswa wajib memiliki kontribusi backend, API, database, frontend, serta dokumentasi.

Face verification dan attendance diletakkan pada Mahasiswa A. Asset, allocation, maintenance, dan reporting diletakkan pada Mahasiswa B. Kedua modul bertemu pada session dan API contract.

## 16.2 Pembagian Final

### Mahasiswa A

Fokus:

* Authentication.
* Authorization.
* User Management.
* Face Enrollment.
* Face Verification.
* Session Service.
* Attendance Service.
* Web Mahasiswa.
* API Face.
* API Attendance.
* Database terkait user, face, session, attendance.

### Mahasiswa B

Fokus:

* Asset Management.
* Asset Allocation.
* Asset Return.
* Asset History.
* Maintenance.
* Reporting.
* Dashboard Admin.
* Web Admin.
* API Asset.
* API Allocation.
* API Maintenance.
* API Report.
* Database terkait asset, allocation, maintenance.

### Tanggung Jawab Bersama

* API contract.
* Database migration integration.
* Testing.
* Security review.
* Deployment.
* Swagger/OpenAPI.
* Dokumentasi.
* Journal.
* Tutorial.
* Poster.
* Bug fixing.

## 16.3 Matriks Tugas

| Mahasiswa | Modul             | Backend | API     | Frontend      | Database        | Dokumentasi |
| --------- | ----------------- | ------- | ------- | ------------- | --------------- | ----------- |
| A         | Authentication    | Utama   | Utama   | Web Mahasiswa | User/Role       | Utama       |
| A         | Face Enrollment   | Utama   | Utama   | Web Mahasiswa | Face Template   | Utama       |
| A         | Face Verification | Utama   | Utama   | Web Mahasiswa | Face/Attendance | Utama       |
| A         | Session           | Utama   | Utama   | Web Mahasiswa | Session         | Bersama     |
| A         | Attendance        | Utama   | Utama   | Web Mahasiswa | Attendance      | Utama       |
| B         | Asset             | Utama   | Utama   | Web Admin     | Assets          | Utama       |
| B         | Allocation        | Utama   | Utama   | Web Admin     | Allocation      | Utama       |
| B         | Maintenance       | Utama   | Utama   | Web Admin     | Maintenance     | Utama       |
| B         | Reporting         | Utama   | Utama   | Web Admin     | Query/Report    | Utama       |
| B         | Dashboard Admin   | Utama   | Utama   | Web Admin     | Aggregation     | Utama       |
| A+B       | Integration       | Bersama | Bersama | Bersama       | Bersama         | Bersama     |
| A+B       | Testing           | Bersama | Bersama | Bersama       | Bersama         | Bersama     |
| A+B       | Deployment        | Bersama | Bersama | Bersama       | Bersama         | Bersama     |

---

# ========================================

# PART 17 — NON-FUNCTIONAL REQUIREMENTS

# ========================================

## 17.1 Performance

Target performance proyek:

* API CRUD umum memiliki target response normal di bawah 2 detik pada lingkungan uji dengan beban ringan.
* Face verification memiliki target waktu respons di bawah 5 detik pada lingkungan pengujian normal.
* Target tersebut merupakan **target desain**, bukan hasil pengujian yang telah terbukti.

Hasil aktual:

**[DIISI SETELAH PENGUJIAN]**

## 17.2 Reliability

Sistem harus:

* menggunakan transaction pada operasi multi-tabel;
* memiliki validation;
* mencegah duplicate attendance;
* mencegah double allocation;
* menangani error service;
* menjaga konsistensi status asset.

## 17.3 Security

Sistem harus:

* menggunakan HTTPS pada deployment;
* menyimpan password dalam bentuk hash;
* menggunakan authentication dan authorization;
* membatasi endpoint;
* melindungi biometric data;
* melakukan input validation;
* tidak menyimpan secret di source code.

## 17.4 Usability

Antarmuka harus:

* memiliki navigasi jelas;
* memberikan pesan error yang dapat dipahami;
* menunjukkan status presensi;
* menunjukkan status camera;
* menyediakan feedback verification;
* membedakan role mahasiswa dan admin.

## 17.5 Maintainability

Backend harus dipisahkan berdasarkan domain/service.

Struktur konseptual backend Go:

```text
cmd/
├── server/
│   └── main.go
internal/
├── handler/
├── middleware/
├── model/
├── repository/
├── service/
├── config/
└── router/
```

Struktur konseptual microservice face verification (Python):

```text
face_service/
├── main.py
├── routes/
├── services/
└── requirements.txt
```

## 17.6 Scalability

Arsitektur mendukung penambahan:

* session;
* user;
* asset;
* endpoint;
* frontend component.

Namun sistem belum ditujukan untuk high-scale enterprise deployment.

## 17.7 Availability

Target availability:

* sistem dapat diakses pada lingkungan demo;
* downtime akibat deployment/development masih dimungkinkan;
* target uptime formal enterprise tidak ditetapkan.

Production uptime:

**[PERLU VALIDASI]**

---

# ========================================

# PART 18 — TESTING

# ========================================

## 18.1 Unit Testing

Unit testing dilakukan pada:

* validation;
* password verification;
* token generation;
* authorization;
* session validation;
* duplicate attendance;
* asset availability;
* allocation rule;
* maintenance status transition;
* report filter.

## 18.2 API Testing

API testing memvalidasi:

* HTTP method;
* endpoint;
* status code;
* schema;
* validation;
* authentication;
* authorization;
* error response.

## 18.3 Integration Testing

Integration testing mencakup:

```text
Frontend
   ↓
REST API
   ↓
Service
   ↓
Repository
   ↓
PostgreSQL
```

Flow yang diuji:

* login;
* enrollment;
* verification;
* attendance;
* allocation;
* return;
* maintenance;
* report.

## 18.4 Functional Testing

Functional testing membandingkan:

**Input → Proses → Expected Output**

berdasarkan acceptance criteria setiap fitur.

## 18.5 Security Testing

Security testing meliputi:

* unauthorized request;
* wrong role;
* invalid token;
* expired token;
* invalid input;
* oversized image;
* invalid image format;
* duplicate operation;
* sensitive data exposure.

## 18.6 Face Verification Testing

Face verification diuji menggunakan:

1. Wajah user yang sama.
2. Wajah user berbeda.
3. Kamera buruk.
4. Pencahayaan rendah.
5. Wajah tidak terdeteksi.
6. Dua wajah.
7. Citra blur.
8. Pose berbeda.
9. Kondisi perangkat berbeda.

Threshold diuji pada dataset validasi internal:

**[DATASET PENGUJIAN PERLU DITENTUKAN]**

## 18.7 Test Case

| ID     | Test Case                    | Input               | Expected Result     | Actual Result             | Status                    |
| ------ | ---------------------------- | ------------------- | ------------------- | ------------------------- | ------------------------- |
| TC-001 | Login valid                  | Credential valid    | Login berhasil      | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-002 | Login password salah         | Password salah      | 401                 | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-003 | User inactive                | Akun inactive       | 403                 | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-004 | Face enrollment valid        | Image satu wajah    | Template tersimpan  | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-005 | Enrollment tanpa wajah       | Image tanpa wajah   | Enrollment ditolak  | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-006 | Enrollment dua wajah         | Image dua wajah     | Enrollment ditolak  | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-007 | Verification user sama       | Wajah user          | Verified            | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-008 | Verification user berbeda    | Wajah orang lain    | Failed              | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-009 | Attendance valid             | Session + verified  | Attendance created  | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-010 | Duplicate attendance         | Attendance existing | Rejected            | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-011 | Check-out valid              | Existing attendance | Checkout saved      | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-012 | Asset create                 | Asset code unique   | Created             | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-013 | Duplicate asset code         | Existing code       | Rejected            | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-014 | Allocation valid             | Available asset     | Allocated           | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-015 | Double allocation            | Allocated asset     | Rejected            | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-016 | Asset return                 | Active allocation   | Returned            | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-017 | Maintenance create           | Damaged asset       | Maintenance created | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-018 | Maintenance asset allocation | Asset MAINTENANCE   | Allocation rejected | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-019 | Unauthorized asset update    | Student token       | 403                 | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |
| TC-020 | Report attendance            | Valid date filter   | Report generated    | [DIISI SETELAH PENGUJIAN] | [DIISI SETELAH PENGUJIAN] |

---

# ========================================

# PART 19 — ACCEPTANCE CRITERIA SISTEM

# ========================================

Sistem dianggap memenuhi acceptance criteria apabila:

1. Mahasiswa dapat login menggunakan akun aktif.
2. Admin dapat login menggunakan akun operasional.
3. Role mahasiswa dan admin diterapkan pada backend.
4. Password disimpan dalam bentuk hash.
5. Mahasiswa dapat melakukan face enrollment.
6. Face verification berjalan menggunakan 1:1 verification.
7. Sistem menggunakan template milik user yang sedang login.
8. Sistem tidak menggunakan face identification 1:N sebagai metode presensi utama.
9. Presensi hanya dapat dilakukan pada session yang valid.
10. Presensi hanya dapat dibuat setelah verification berhasil.
11. Duplicate attendance ditolak.
12. Check-out dapat dilakukan setelah check-in.
13. Admin dapat membuat session.
14. Admin dapat mendaftarkan aset.
15. Asset code bersifat unik.
16. Admin dapat mengalokasikan asset.
17. Asset yang sedang dialokasikan tidak dapat dialokasikan kembali.
18. Asset return dapat dicatat.
19. Histori penggunaan asset tersimpan.
20. Maintenance dapat dibuat.
21. Asset maintenance tidak dapat dialokasikan.
22. Maintenance dapat diselesaikan dan status asset diperbarui.
23. Report attendance tersedia.
24. Report asset tersedia.
25. Report maintenance tersedia.
26. Dashboard tersedia.
27. Web Mahasiswa menggunakan REST API.
28. Web Admin menggunakan REST API.
29. Tidak ada akses database langsung dari frontend.
30. REST API terdokumentasi menggunakan Swagger/OpenAPI.
31. API memiliki error handling terstruktur.
32. Audit event penting tercatat.
33. Biometric data tidak ditampilkan sebagai raw data pada frontend.
34. Testing dilakukan terhadap core business flow.
35. Sistem dapat dideploy pada environment demo.
36. Seluruh fitur MVP memiliki acceptance criteria yang telah diuji.

---

# ========================================

# PART 20 — RISIKO DAN MITIGASI

# ========================================

| Risk                       | Probability | Impact | Mitigation                                                  | Contingency                                      |
| -------------------------- | ----------- | ------ | ----------------------------------------------------------- | ------------------------------------------------ |
| Face verification gagal    | Medium      | High   | Threshold diuji, validasi citra, retry terbatas             | Verification diulang sesuai aturan               |
| Kualitas kamera rendah     | High        | Medium | Image quality validation, instruksi posisi                  | Gunakan perangkat dengan kamera memadai          |
| Pencahayaan buruk          | High        | Medium | Panduan pencahayaan, face quality check                     | Capture ulang                                    |
| Data biometrik bocor       | Low         | High   | Encryption/protection, access control, no raw image storage | Revoke template, evaluasi incident               |
| Spoofing                   | Medium      | High   | Login + 1:1 + audit + rate limit                            | Review manual/admin bila diperlukan              |
| Duplicate attendance       | Medium      | High   | Unique constraint user-session                              | Reject transaction                               |
| Asset double allocation    | Medium      | High   | Transaction + unique active allocation                      | Reject allocation                                |
| Database inconsistency     | Medium      | High   | FK, transaction, constraints                                | Rollback transaction                             |
| API failure                | Medium      | High   | Error handling, logging, health check                       | Retry/restart service                            |
| Scope creep                | High        | High   | Prioritas MVP                                               | Fitur Nice to Have ditunda                       |
| Deployment failure         | Medium      | Medium | Environment config, deployment checklist                    | Rollback/redeploy                                |
| Keterlambatan sprint       | Medium      | High   | Sprint backlog, dependency mapping                          | Reprioritize Must Have                           |
| Model terlalu berat        | Medium      | Medium | Benchmark local CPU                                         | Gunakan model yang lebih ringan setelah validasi |
| Lisensi model tidak sesuai | Low         | High   | License review                                              | Ganti model/library sebelum deployment           |
| Perbedaan browser kamera   | Medium      | Medium | Uji browser utama                                           | Dokumentasikan browser yang didukung             |
| Data aktual belum tersedia | High        | Medium | Seed/sample data                                            | `[PERLU VALIDASI]` saat integrasi data nyata     |

---

# ========================================

# PART 21 — TECHNOLOGY STACK

# ========================================

## 21.1 Frontend

### Alternatif A — HTML/CSS/JavaScript

| Aspek        | Penilaian                                             |
| ------------ | ----------------------------------------------------- |
| Kelebihan    | Sederhana, ringan, dependency sedikit, mudah dipahami |
| Kekurangan   | Struktur komponen dan state management harus disiplin |
| Kompleksitas | Rendah                                                |
| Keputusan    | Tidak dipilih                                         |

### Alternatif B — React + Vite + Tailwind CSS

| Aspek        | Penilaian                                                                      |
| ------------ | ------------------------------------------------------------------------------ |
| Kelebihan    | Component-based, state management terstruktur, hot reload, utility-first CSS   |
| Kekurangan   | Dependency dan build tooling lebih banyak                                      |
| Kompleksitas | Sedang                                                                         |
| Keputusan    | **FINAL**                                                                      |

**Keputusan final frontend:** **React 19 + TypeScript + Vite 8 + Tailwind CSS v4**.

Alasan keputusan adalah kebutuhan frontend memiliki banyak komponen interaktif (face capture, dashboard, tabel data, form kompleks) yang lebih efisien dikelola menggunakan component-based architecture dengan state management React. Tailwind CSS v4 menyediakan utility-first styling yang konsisten. TypeScript memberikan type safety pada integrasi API.

## 21.2 Backend

### Go + Fiber

| Aspek        | Penilaian                                                                           |
| ------------ | ----------------------------------------------------------------------------------- |
| Kelebihan    | Performa tinggi, compiled binary, type-safe, concurrency bawaan, dependency minimal |
| Kekurangan   | Face verification library (DeepFace) tidak tersedia di Go                           |
| Kompleksitas | Sedang                                                                              |
| Keputusan    | **FINAL** (backend utama)                                                           |

### FastAPI

| Aspek        | Penilaian                                                                  |
| ------------ | -------------------------------------------------------------------------- |
| Kelebihan    | Python, cocok untuk API, type hints, validation, integrasi library AI/face |
| Kekurangan   | Struktur aplikasi perlu dirancang sendiri                                  |
| Kompleksitas | Sedang                                                                     |
| Keputusan    | **FINAL** (microservice face verification)                                 |

### Flask

| Aspek        | Penilaian                                               |
| ------------ | ------------------------------------------------------- |
| Kelebihan    | Ringan dan sederhana                                    |
| Kekurangan   | Banyak komponen harus dipilih dan diintegrasikan manual |
| Kompleksitas | Sedang                                                  |
| Keputusan    | Tidak dipilih                                           |

### Django

| Aspek        | Penilaian                      |
| ------------ | ------------------------------ |
| Kelebihan    | Banyak fitur bawaan            |
| Kekurangan   | Lebih besar dari kebutuhan MVP |
| Kompleksitas | Sedang–tinggi                  |
| Keputusan    | Tidak dipilih                  |

### Node.js/Express

| Aspek        | Penilaian                                                  |
| ------------ | ---------------------------------------------------------- |
| Kelebihan    | Cocok untuk REST API                                       |
| Kekurangan   | Face processing tetap membutuhkan stack Python/AI tambahan |
| Kompleksitas | Sedang–tinggi untuk proyek ini                             |
| Keputusan    | Tidak dipilih                                              |

**Keputusan final backend:** **Go + Fiber sebagai backend utama + Python FastAPI sebagai microservice face verification**.

Go + Fiber dipilih sebagai backend utama untuk menangani seluruh REST API (auth, CRUD, business logic) karena performa tinggi dan type safety. Python FastAPI digunakan sebagai microservice terpisah khusus face verification karena DeepFace merupakan library Python. Arsitektur microservice ini memisahkan concern face verification dari business logic utama.

## 21.3 Database

### PostgreSQL

| Aspek        | Penilaian                                               |
| ------------ | ------------------------------------------------------- |
| Kelebihan    | Relational, constraint kuat, transaction, UUID, JSONB   |
| Kekurangan   | Setup sedikit lebih kompleks daripada database embedded |
| Kompleksitas | Sedang                                                  |
| Keputusan    | **FINAL**                                               |

### MySQL

| Aspek        | Penilaian                                                                    |
| ------------ | ---------------------------------------------------------------------------- |
| Kelebihan    | Familiar dan umum digunakan                                                  |
| Kekurangan   | Tidak memberikan keuntungan signifikan dibanding PostgreSQL untuk desain ini |
| Kompleksitas | Sedang                                                                       |
| Keputusan    | Tidak dipilih                                                                |

**Keputusan final database:** **PostgreSQL (Supabase)**.

Database menggunakan Supabase sebagai managed PostgreSQL yang dapat diakses melalui connection string dari backend Go maupun microservice Python. Pada MVP, Supabase Storage tidak digunakan untuk menyimpan foto wajah mentah. Foto diproses sementara di memori dan yang disimpan hanya sampel embedding terenkripsi di PostgreSQL.

## 21.4 Face Verification

### DeepFace

| Aspek        | Penilaian                                                                        |
| ------------ | -------------------------------------------------------------------------------- |
| Kelebihan    | Python, API verification/embedding tersedia, mendukung beberapa pretrained model |
| Kekurangan   | Dependensi model dan lisensi model perlu diperiksa                               |
| Kompleksitas | Sedang                                                                           |
| Keputusan    | **FINAL**                                                                        |

DeepFace menyediakan fungsi verification serta representasi embedding dan mendukung model seperti ArcFace.

### InsightFace

| Aspek        | Penilaian                                                         |
| ------------ | ----------------------------------------------------------------- |
| Kelebihan    | Ekosistem face analysis dan ArcFace                               |
| Kekurangan   | Model/runtime dapat memiliki kompleksitas deployment lebih tinggi |
| Kompleksitas | Sedang–tinggi                                                     |
| Keputusan    | Alternatif teknis untuk evaluasi                                  |

**Keputusan final:** **DeepFace dengan model pretrained yang divalidasi, dengan ArcFace sebagai kandidat utama model embedding.**

## 21.5 API Documentation

**Teknologi:** Swagger/OpenAPI.

Keputusan:

* OpenAPI sebagai contract.
* Swagger UI sebagai interactive documentation.
* OpenAPI schema menjadi sumber dokumentasi.

## 21.6 Version Control

**Git + GitHub**

Penggunaan:

```text
main
develop
feature/*
```

Branching final dapat disederhanakan jika workload dua mahasiswa membutuhkan workflow yang lebih ringan.

## 21.7 Deployment

**Platform final:** Railway, Render, atau layanan PaaS setara yang mendukung Go, Python, dan static hosting.

Komponen:

* Static Web Site untuk frontend React;
* Web Service untuk Go + Fiber (backend utama);
* Web Service untuk Python FastAPI (microservice face verification);
* Supabase managed PostgreSQL database;
* Tidak menggunakan persistent face-image storage pada MVP; sampel embedding disimpan terenkripsi di PostgreSQL.

Nama environment production:

**[PERLU VALIDASI]**

Deployment platform dipilih berdasarkan kemudahan deployment mahasiswa dan kemampuan menyediakan environment web backend/database tanpa membangun server infrastructure sendiri.

---

# ========================================

# PART 22 — DEPLOYMENT ARCHITECTURE

# ========================================

## 22.1 Arsitektur Deployment

```text
                    INTERNET
                       │
                       ▼
                Browser Mahasiswa
                       │
                       ▼
                 Web Mahasiswa
                       │
                       │ HTTPS / REST API
                       ▼
                ┌──────────────┐
                │ Reverse/API  │
                │   Gateway    │
                └──────┬───────┘
                       │
                       ▼
                ┌──────────────┐
                │  Go + Fiber  │
                │   Backend    │
                └──────┬───────┘
                       │
          ┌────────────┼─────────────┐
          │            │             │
          ▼            ▼             ▼
   Python FastAPI   Services    Supabase
   Face Service                 PostgreSQL
          │
          ▼
   Pretrained Model
```

Web Admin mengikuti jalur yang sama:

```text
Browser Admin
    ↓
Web Admin
    ↓
REST API
    ↓
Go + Fiber Backend
    ↓
Supabase PostgreSQL
```

## 22.2 Lokasi Komponen

**Frontend:** Static web hosting/PaaS (React build output).

**Backend:** Go + Fiber Web Service.

**Face Verification:** Python FastAPI microservice (internal, tidak diakses langsung oleh frontend).

**Face Model:** Berjalan pada server microservice Python sebagai local dependency/model artifact.

**Database:** Supabase managed PostgreSQL.

**Storage:** Tidak menggunakan persistent face-image storage pada MVP; sampel embedding disimpan terenkripsi di PostgreSQL.

**Source Code:** GitHub.

## 22.3 Environment

Environment minimal:

* Development.
* Testing/Staging.
* Production/Demo.

Secret disimpan sebagai environment variable.

---

# ========================================

# PART 23 — PROJECT FEASIBILITY

# ========================================

## 23.1 Feasibility Teknis

Proyek feasible secara teknis karena seluruh komponen utama dapat dibangun menggunakan teknologi yang tersedia untuk proyek mahasiswa. REST API, relational database, frontend web, dan pretrained face verification tidak memerlukan pembuatan framework baru.

Kompleksitas tertinggi berada pada face verification, terutama pada preprocessing, kualitas citra, threshold, resource penggunaan model, dan privacy.

## 23.2 Feasibility Waktu

Proyek feasible apabila tujuh sprint difokuskan pada Must Have. Fitur Nice to Have tidak menjadi dependency MVP.

Urutan pengembangan memastikan dependency penting diselesaikan terlebih dahulu.

## 23.3 Feasibility Tim

Dua mahasiswa dapat membangun sistem dengan pembagian domain:

* Mahasiswa A: identity, face, attendance, Web Mahasiswa.
* Mahasiswa B: asset, allocation, maintenance, reporting, Web Admin.

Integrasi API dan database menjadi tanggung jawab bersama.

## 23.4 Feasibility Data

Data minimum berupa:

* user;
* session;
* asset;
* attendance;
* allocation;
* maintenance.

Data aktual institusi:

**[PERLU VALIDASI]**

Pengembangan awal dapat menggunakan data uji yang tidak mewakili data resmi institusi.

## 23.5 Feasibility Face Verification

Face verification feasible sebagai fitur proyek apabila menggunakan pretrained model.

Sistem tidak melatih model dari nol. Fokus implementasi adalah pipeline:

```text
Camera
→ Detection
→ Embedding
→ Similarity
→ Threshold
→ Verification
```

Keberhasilan akhir harus ditentukan melalui pengujian nyata.

## 23.6 Feasibility API

REST API feasible dan menjadi fokus utama proyek. API memiliki domain yang jelas dan dapat terdokumentasi melalui OpenAPI.

## 23.7 Feasibility Database

Relational database dengan sembilan tabel inti cukup untuk seluruh core business flow. Foreign key dan transaction mendukung konsistensi data.

## 23.8 Scope Feasibility

Scope MVP dapat diselesaikan dengan dua mahasiswa karena:

* hanya dua role;
* hanya dua web app;
* satu backend;
* satu database;
* 1:1 verification;
* tidak ada mobile app;
* tidak ada hardware;
* tidak ada training model;
* tidak ada IoT.

## 23.9 Risiko Utama

Risiko utama:

1. Threshold face verification belum tervalidasi.
2. Performa model pada server deployment.
3. Kualitas kamera.
4. Privacy biometric.
5. Integrasi dua frontend dengan API.
6. Scope creep.

## 23.10 Kesimpulan Feasibility

Proyek dinilai **feasible untuk dikerjakan sebagai proyek satu semester oleh dua mahasiswa** dengan syarat MVP menjadi prioritas, face verification menggunakan pretrained model, threshold ditentukan melalui pengujian, dan fitur tambahan tidak menggeser penyelesaian komponen inti.

---

# ========================================

# PART 24 — KETERKAITAN DENGAN LUARAN PROYEK 3

# ========================================

## 24.1 Aplikasi Web Services / Doc API

**Yang dihasilkan:**

* Go + Fiber backend;
* Python FastAPI microservice face verification;
* REST API;
* OpenAPI schema;
* Swagger UI;
* endpoint authentication, face, attendance, asset, allocation, maintenance, report.

**Bagian proyek yang digunakan:**

* Backend Service;
* Database;
* API contract;
* Authentication;
* Business Logic.

**Materi yang dapat dimasukkan:**

* arsitektur;
* endpoint;
* request/response;
* validation;
* authentication;
* testing API.

**Hubungan dengan sistem:** REST API menjadi pusat komunikasi dua web application dan database.

## 24.2 Draft Jurnal

**Yang dihasilkan:**

* draft penelitian/pengembangan sistem.

**Bagian proyek yang digunakan:**

* masalah presensi;
* asset management;
* face verification;
* Agile;
* REST API;
* testing.

**Materi yang dapat dimasukkan:**

* latar belakang;
* metodologi;
* arsitektur;
* implementation;
* hasil pengujian;
* pembahasan.

## 24.3 Buku Tutorial

**Yang dihasilkan:**

* panduan pembangunan sistem.

**Bagian proyek yang digunakan:**

* setup environment;
* database;
* backend;
* API;
* face verification;
* frontend;
* deployment.

**Hubungan:** Buku tutorial mendokumentasikan proses teknis pembangunan dari awal sampai deployment.

## 24.4 Poster

**Yang dihasilkan:**

* visual ringkas sistem.

**Bagian proyek yang digunakan:**

* masalah;
* solusi;
* inovasi;
* arsitektur;
* face verification flow;
* attendance;
* asset flow;
* teknologi.

**Hubungan:** Poster menjadi media komunikasi visual dari sistem yang dibangun.

---

# ========================================

# PART 25 — RANCANGAN DRAFT JURNAL

# ========================================

## 25.1 Judul Jurnal

**Implementasi Face Verification 1:1 pada Web Service Presensi dan Manajemen Aset Laboratorium Komputer Menggunakan Metode Agile**

## 25.2 Abstrak

**Draft:**

Kegiatan laboratorium komputer membutuhkan pengelolaan presensi mahasiswa dan aset perangkat secara terintegrasi agar data kehadiran, penggunaan perangkat, dan kondisi aset dapat ditelusuri secara konsisten. Penelitian/pengembangan ini merancang dan mengimplementasikan web service presensi dan manajemen aset laboratorium komputer dengan menggunakan RESTful Web Service sebagai pusat komunikasi antara aplikasi web dan database. Sistem terdiri atas Web Mahasiswa dan Web Admin/Laboran/Tata Usaha. Presensi menggunakan pendekatan face verification 1:1 melalui kamera browser. Mahasiswa melakukan login terlebih dahulu sehingga sistem memperoleh identitas pengguna, kemudian citra wajah dibandingkan dengan template wajah milik pengguna tersebut menggunakan pretrained face model. Sistem juga mengelola sesi praktikum, attendance, asset allocation, histori penggunaan perangkat, maintenance, dan reporting. Pengembangan dilakukan dengan metode Agile melalui beberapa sprint yang mencakup analisis, arsitektur, backend, face verification, asset management, integrasi, pengujian, dan dokumentasi. Evaluasi sistem mencakup pengujian API, integrasi, fungsionalitas, keamanan, dan face verification. Hasil pengujian final mengenai performa dan tingkat keberhasilan verifikasi diisi setelah implementasi dan pengujian aktual.

## 25.3 Pendahuluan

Pendahuluan membahas:

1. Kebutuhan digitalisasi operasional laboratorium.
2. Integrasi presensi dan aset.
3. Keterbatasan pencatatan yang tidak terintegrasi.
4. Kebutuhan REST API.
5. Pemanfaatan face verification 1:1.
6. Kebutuhan data asset yang dapat ditelusuri.
7. Agile sebagai metode pengembangan.
8. Tujuan proyek.

## 25.4 Identifikasi Masalah

Masalah yang dibahas:

* verifikasi identitas mahasiswa;
* pencatatan kehadiran;
* pelacakan unit perangkat;
* histori penggunaan;
* maintenance;
* integrasi data;
* kebutuhan backend web service.

## 25.5 Tujuan

Tujuan jurnal:

1. Mendesain RESTful Web Service.
2. Mengimplementasikan face verification 1:1.
3. Mengimplementasikan attendance.
4. Mengimplementasikan asset management.
5. Mengevaluasi sistem melalui testing.

## 25.6 Tinjauan Pustaka

Subtopik:

1. Web Service.
2. REST API.
3. Backend development.
4. Relational database.
5. Face verification.
6. Face embedding.
7. Pretrained face model.
8. Attendance system.
9. Asset management.
10. Agile software development.

**Referensi ilmiah final:**

**[REFERENSI PERLU DICARI]**

## 25.7 Metodologi

Metodologi mencakup:

* requirement analysis;
* system design;
* Agile sprint;
* database design;
* REST API development;
* face verification implementation;
* frontend integration;
* testing.

## 25.8 Implementasi

Implementasi membahas:

1. Database PostgreSQL (Supabase).
2. Go + Fiber backend.
3. Python FastAPI face verification microservice.
3. Authentication.
4. User management.
5. Face enrollment.
6. Face verification.
7. Attendance.
8. Asset management.
9. Allocation.
10. Maintenance.
11. Reporting.
12. Web applications.
13. Swagger/OpenAPI.

## 25.9 Pengujian

Pengujian meliputi:

* Unit testing.
* API testing.
* Integration testing.
* Functional testing.
* Security testing.
* Face verification testing.

## 25.10 Hasil

**[DIISI SETELAH IMPLEMENTASI DAN PENGUJIAN]**

Minimal data yang akan dimasukkan:

* jumlah test case;
* pass/fail;
* response performance;
* verification result;
* false acceptance;
* false rejection;
* API availability pada demo environment.

## 25.11 Pembahasan

Pembahasan membandingkan hasil pengujian dengan requirement dan acceptance criteria.

Aspek yang dibahas:

* keberhasilan REST API;
* konsistensi database;
* keberhasilan face verification;
* usability;
* risiko spoofing;
* performa;
* keterbatasan sistem.

## 25.12 Kesimpulan

**[DIISI SETELAH IMPLEMENTASI DAN PENGUJIAN]**

Kesimpulan harus menjawab:

1. Apakah REST API berhasil diimplementasikan?
2. Apakah dua web application berhasil terintegrasi?
3. Apakah face verification 1:1 dapat berjalan?
4. Apakah attendance berhasil dicatat?
5. Apakah asset lifecycle berhasil dikelola?
6. Apa keterbatasan sistem berdasarkan hasil pengujian?

---

# ========================================

# PART 26 — RANCANGAN BUKU TUTORIAL

# ========================================

## Bab 1 — Pengenalan Proyek

Membahas tujuan, ruang lingkup, arsitektur, dan hasil akhir sistem.

## Bab 2 — Analisis Masalah

Membahas permasalahan presensi, asset tracking, dan maintenance.

## Bab 3 — Requirement

Membahas functional dan non-functional requirements.

## Bab 4 — Arsitektur Sistem

Membahas Web Mahasiswa, Web Admin, REST API, backend, face service, dan database.

## Bab 5 — Database

Membahas:

* schema;
* table;
* primary key;
* foreign key;
* relation;
* migration.

## Bab 6 — Backend

Membahas:

* setup Go + Fiber;
* project structure;
* routing;
* schema;
* service;
* repository.

## Bab 7 — Authentication

Membahas:

* login;
* password hashing;
* JWT;
* cookie;
* role.

## Bab 8 — Face Enrollment

Membahas:

* browser camera;
* image capture;
* face detection;
* embedding;
* template storage.

## Bab 9 — Face Verification

Membahas:

* verification 1:1;
* similarity;
* threshold;
* error handling;
* limitation.

## Bab 10 — REST API

Membahas endpoint dan Swagger/OpenAPI.

## Bab 11 — Attendance

Membahas session, check-in, duplicate prevention, dan check-out.

## Bab 12 — Asset Management

Membahas asset master data dan status perangkat.

## Bab 13 — Allocation

Membahas allocation, return, dan consistency.

## Bab 14 — Maintenance

Membahas issue, maintenance status, dan asset status.

## Bab 15 — Web Mahasiswa

Membahas UI:

* login;
* enrollment;
* camera;
* attendance;
* allocation history.

## Bab 16 — Web Admin

Membahas UI:

* dashboard;
* users;
* sessions;
* assets;
* allocations;
* maintenance;
* reports.

## Bab 17 — Testing

Membahas:

* unit;
* API;
* integration;
* functional;
* security;
* face verification.

## Bab 18 — Swagger/OpenAPI

Membahas penggunaan dokumentasi API sebagai kontrak frontend-backend.

## Bab 19 — Deployment

Membahas:

* environment variable;
* backend;
* frontend;
* database;
* HTTPS;
* production configuration.

## Bab 20 — Troubleshooting

Membahas:

* camera permission;
* model loading;
* database connection;
* API error;
* CORS;
* authentication;
* deployment;
* face verification.

---

# ========================================

# PART 27 — RANCANGAN POSTER

# ========================================

## 27.1 Judul

**RANCANG BANGUN WEB SERVICE PRESENSI DAN MANAJEMEN ASET LABORATORIUM KOMPUTER MENGGUNAKAN METODE AGILE**

## 27.2 Masalah

* Presensi memerlukan verifikasi identitas mahasiswa.
* Penggunaan unit perangkat perlu dilacak secara individual.
* Data presensi, session, asset, dan maintenance perlu terintegrasi.
* Diperlukan REST API sebagai pusat komunikasi aplikasi dan database.

## 27.3 Solusi

Web service terintegrasi yang menyediakan:

```text
Web Mahasiswa
        +
Web Admin
        ↓
    REST API
        ↓
    PostgreSQL
```

dengan face verification 1:1 untuk presensi.

## 27.4 Inovasi

**Face Verification 1:1 + Asset Tracking Individual + REST API Integration**

Konsep:

```text
Login
 ↓
Known Identity
 ↓
Face Verification
 ↓
Attendance
 ↓
Asset Allocation
 ↓
Asset History
 ↓
Maintenance
```

## 27.5 Arsitektur

```text
Web Mahasiswa ─┐
               ├── REST API ── Go+Fiber Backend ── Supabase PostgreSQL
Web Admin ─────┘
                    │
                    └── Python FastAPI Face Verification Service
```

## 27.6 Face Verification Flow

```text
Login
 ↓
Camera
 ↓
Face Detection
 ↓
Embedding
 ↓
Compare with User Template
 ↓
Threshold
 ↓
Verified / Failed
```

## 27.7 Attendance Flow

```text
Session
 ↓
Login
 ↓
Face Verification
 ↓
Check-in
 ↓
Attendance Record
```

## 27.8 Asset Flow

```text
Asset Registration
 ↓
Available
 ↓
Allocation
 ↓
In Use
 ↓
Return
 ↓
Available / Maintenance
```

## 27.9 Teknologi

* React 19
* TypeScript
* Tailwind CSS v4
* Go + Fiber
* Python FastAPI
* Supabase PostgreSQL
* DeepFace
* Pretrained face model
* Swagger/OpenAPI
* Git/GitHub
* PaaS deployment

## 27.10 Hasil yang Ditargetkan

1. Dua aplikasi web terintegrasi.
2. REST API terdokumentasi.
3. Database terstruktur.
4. Face enrollment.
5. Face verification 1:1.
6. Attendance.
7. Asset management.
8. Allocation dan return.
9. Maintenance.
10. Reporting.

## 27.11 QR/Link Demo

**QR Demo:** `[DIISI SETELAH DEPLOYMENT]`

**API Documentation:** `[DIISI SETELAH DEPLOYMENT]`

**Repository:** `[DIISI SETELAH REPOSITORY DITETAPKAN]`

## 27.12 Layout Poster

```text
┌─────────────────────────────────────────────────────────┐
│                     JUDUL PROYEK                        │
│             Nama Anggota / Program Studi               │
├───────────────────────┬─────────────────────────────────┤
│ MASALAH               │ SOLUSI                          │
│ - Presensi            │ Web User + Web Admin            │
│ - Asset tracking      │ REST API + PostgreSQL           │
│ - Maintenance         │ Face Verification 1:1           │
├───────────────────────┴─────────────────────────────────┤
│                  ARSITEKTUR SISTEM                     │
│ Web → REST API → Backend → DB + Face Verification      │
├───────────────────────┬─────────────────────────────────┤
│ FACE VERIFICATION     │ ATTENDANCE FLOW                 │
│ Camera → Embedding    │ Login → Verify → Check-in      │
│ → Similarity → Result │ → Database                      │
├───────────────────────┼─────────────────────────────────┤
│ ASSET FLOW            │ TEKNOLOGI                       │
│ Register → Allocate   │ Go+Fiber, React, Supabase,      │
│ → Return → Maintenance│ DeepFace, OpenAPI               │
├───────────────────────┴─────────────────────────────────┤
│              HASIL YANG DITARGETKAN                     │
│  Web Service | Presensi | Asset | Maintenance | Report │
├─────────────────────────────────────────────────────────┤
│                 QR DEMO / REPOSITORY                    │
└─────────────────────────────────────────────────────────┘
```

---

# ========================================

# PART 28 — FINAL SYSTEM SPECIFICATION

# ========================================

## 28.1 Produk

Produk adalah sistem web terintegrasi untuk presensi dan manajemen aset laboratorium komputer.

## 28.2 Dua Web Application

1. **Web Mahasiswa/User**
2. **Web Admin/Laboran/Tata Usaha**

## 28.3 Backend

**Go + Fiber** sebagai backend service utama. **Python FastAPI** sebagai microservice face verification.

## 28.4 REST API

REST API menjadi pusat komunikasi antara frontend dan database.

Base path:

`/api/v1`

## 28.5 Database

**PostgreSQL (Supabase)**

Tabel inti:

* roles
* users
* face_templates
* lab_sessions
* attendance
* assets
* asset_allocations
* maintenance
* audit_logs

## 28.6 Face Verification

* Browser camera.
* Login terlebih dahulu.
* 1:1 verification.
* Pretrained face model.
* Embedding.
* Similarity/distance.
* Threshold berdasarkan pengujian.
* Tidak ada training dari nol.
* Tidak menggunakan 1:N sebagai presensi utama.

## 28.7 Asset Management

* Asset code individual.
* Master asset.
* Status.
* Allocation.
* Return.
* History.

## 28.8 Maintenance

* Issue.
* Priority.
* Status.
* Resolution.
* Asset status transition.

## 28.9 Security

* Password hashing.
* JWT.
* HttpOnly cookie.
* Authorization.
* HTTPS.
* Validation.
* Transaction.
* Biometric protection.
* Audit log.

## 28.10 Agile

Pengembangan menggunakan tujuh sprint:

1. Requirement + Architecture.
2. Database + Backend + Authentication.
3. Face + Attendance.
4. Asset + Allocation.
5. History + Maintenance.
6. Web Applications + Integration.
7. Testing + Documentation + Finalization.

## 28.11 Pembagian Kerja

**Mahasiswa A:** Identity, Face, Session, Attendance, Web Mahasiswa.

**Mahasiswa B:** Asset, Allocation, Maintenance, Reporting, Web Admin.

Backend dan integration menjadi tanggung jawab bersama.

## 28.12 Luaran

1. Web Services / Doc API.
2. Draft Jurnal.
3. Buku Tutorial.
4. Poster.

---

# ========================================

# PART 29 — FINAL MVP

# ========================================

| Fitur                 | Tujuan                            | Output                               |
| --------------------- | --------------------------------- | ------------------------------------ |
| Authentication        | Memastikan identitas user         | Session login                        |
| Authorization         | Membatasi akses role              | Access control                       |
| User Management       | Mengelola user                    | User data                            |
| Face Enrollment       | Menyimpan template wajah          | Active face template                 |
| Face Verification 1:1 | Memverifikasi user                | Verification result                  |
| Session Praktikum     | Menentukan konteks presensi       | Session                              |
| Attendance Check-in   | Mencatat hadir                    | Attendance                           |
| Attendance Check-out  | Mencatat akhir kehadiran          | Checked-out attendance               |
| Asset Management      | Mengidentifikasi setiap unit      | Asset master                         |
| Asset Allocation      | Mencatat perangkat yang digunakan | Allocation                           |
| Asset Return          | Mencatat pengembalian             | Returned allocation                  |
| Asset History         | Melacak histori                   | Asset timeline                       |
| Maintenance           | Mengelola perangkat rusak         | Maintenance record                   |
| Dashboard             | Menampilkan ringkasan             | Dashboard data                       |
| Reporting             | Menyediakan laporan               | Attendance/asset/maintenance reports |
| REST API              | Menjadi pusat komunikasi          | API                                  |
| Swagger/OpenAPI       | Mendokumentasikan API             | API documentation                    |
| Database              | Menyimpan data terintegrasi       | PostgreSQL database                  |
| Web Mahasiswa         | Menyediakan antarmuka user        | Student web                          |
| Web Admin             | Menyediakan antarmuka operasional | Admin web                            |
| Security              | Melindungi sistem                 | Secured application                  |
| Testing               | Memvalidasi fungsi                | Test evidence                        |
| Deployment            | Menyediakan demo sistem           | Deployed application                 |

MVP dianggap utuh apabila seluruh fitur di atas terintegrasi dalam satu alur:

```text
Login
 ↓
Face Enrollment
 ↓
Session
 ↓
Face Verification
 ↓
Attendance
 ↓
Asset Allocation
 ↓
Asset Return
 ↓
Maintenance jika diperlukan
 ↓
History
 ↓
Reporting
```

---

# ========================================

# PART 30 — FINAL OUT OF SCOPE

# ========================================

| Out of Scope                               | Alasan                                             |
| ------------------------------------------ | -------------------------------------------------- |
| Mobile application                         | Platform proyek ditetapkan Web-only                |
| Face identification 1:N                    | Tidak diperlukan karena user login terlebih dahulu |
| Training model sendiri                     | Tidak realistis untuk scope satu semester          |
| IoT                                        | Tidak diperlukan untuk pencatatan inti             |
| RFID                                       | Tidak diperlukan untuk MVP                         |
| Hardware scanner                           | Menambah biaya dan kompleksitas                    |
| Advanced AI                                | Tidak diperlukan untuk kebutuhan inti              |
| Predictive maintenance                     | Membutuhkan data historis dan model tambahan       |
| QR Code sebagai presensi utama             | Metode presensi ditetapkan face verification       |
| Video recording                            | Meningkatkan risiko privacy dan storage            |
| Sistem keuangan aset                       | Di luar kebutuhan inti                             |
| Pengadaan aset                             | Di luar kebutuhan inti                             |
| Penghapusan aset tingkat enterprise        | Tidak diperlukan untuk MVP                         |
| Integrasi SSO kampus                       | Detail sistem eksternal belum tersedia             |
| Integrasi sistem akademik penuh            | Dependency eksternal tinggi                        |
| Multi-campus                               | Di luar skala proyek                               |
| Liveness research khusus                   | Kompleksitas dan scope terlalu tinggi untuk MVP    |
| Training face model dari nol               | Tidak sesuai kebutuhan dan batas waktu             |
| Face identification 1:N sebagai attendance | Tidak sesuai keputusan desain 1:1                  |

---

# ========================================

# PART 31 — GLOSSARY

# ========================================

## Web Service

Layanan perangkat lunak yang menyediakan fungsi atau pertukaran data melalui jaringan sehingga sistem lain dapat berkomunikasi dengannya.

## REST API

Application Programming Interface yang menggunakan prinsip REST dan HTTP method seperti GET, POST, PUT, PATCH, dan DELETE untuk pertukaran data.

## Backend

Bagian sistem yang menangani business logic, authentication, validation, pemrosesan data, dan komunikasi dengan database.

## Frontend

Bagian aplikasi yang berinteraksi langsung dengan pengguna melalui antarmuka web.

## Face Detection

Proses mendeteksi keberadaan dan lokasi wajah pada citra.

## Face Verification

Proses membandingkan dua representasi wajah untuk menentukan apakah keduanya sesuai dengan identitas yang sama.

## Face Embedding

Representasi numerik wajah yang dihasilkan oleh model face recognition.

## Template

Data terstruktur yang digunakan sebagai referensi untuk proses verifikasi. Dalam sistem ini template berupa embedding yang dilindungi dan dikaitkan dengan user tertentu.

## Attendance

Data kehadiran mahasiswa pada suatu sesi praktikum.

## Asset Allocation

Proses mencatat bahwa unit aset tertentu digunakan oleh pengguna tertentu pada session tertentu.

## Maintenance

Proses pencatatan, pemeriksaan, perbaikan, dan penyelesaian masalah pada aset.

## Agile

Metode pengembangan perangkat lunak iteratif yang menghasilkan increment melalui siklus pengembangan singkat.

## Sprint

Periode pengembangan dengan tujuan dan backlog yang telah ditentukan.

## Product Backlog

Daftar seluruh kebutuhan, fitur, dan pekerjaan produk yang diprioritaskan.

## API Documentation

Dokumentasi yang menjelaskan endpoint, method, request, response, authentication, authorization, validation, dan error API.

## Authentication

Proses memastikan identitas pengguna.

## Authorization

Proses menentukan apakah pengguna yang telah terautentikasi memiliki hak untuk menjalankan suatu operasi.

## JWT

JSON Web Token yang digunakan untuk merepresentasikan informasi autentikasi secara terstruktur.

## RESTful Web Service

Web service yang menggunakan prinsip REST dan HTTP sebagai mekanisme komunikasi antar aplikasi.

## Database

Sistem penyimpanan data terstruktur yang digunakan untuk menyimpan user, session, attendance, asset, allocation, maintenance, dan audit.

## Asset Code

Kode unik yang digunakan untuk mengidentifikasi setiap unit perangkat secara individual.

## Asset History

Riwayat penggunaan, alokasi, pengembalian, dan maintenance suatu unit perangkat.

## Check-in

Proses pencatatan waktu awal kehadiran mahasiswa.

## Check-out

Proses pencatatan waktu akhir kehadiran mahasiswa.

## Pretrained Model

Model machine learning yang telah dilatih sebelumnya dan digunakan tanpa melakukan training dari nol pada proyek.

## Similarity

Nilai yang digunakan untuk menunjukkan kedekatan dua embedding sesuai metode perbandingan yang digunakan.

## Threshold

Nilai batas yang digunakan untuk menentukan apakah hasil similarity/distance dianggap memenuhi kriteria verifikasi.

## False Acceptance

Kondisi ketika sistem menerima pasangan wajah yang seharusnya tidak diterima.

## False Rejection

Kondisi ketika sistem menolak pasangan wajah yang seharusnya diterima.

## Spoofing

Upaya memberikan input palsu untuk menipu mekanisme verifikasi.

## Audit Log

Catatan aktivitas penting sistem yang digunakan untuk penelusuran dan pemeriksaan.

## MVP

Minimum Viable Product, yaitu versi minimum sistem yang telah mencakup seluruh fungsi inti dan menghasilkan satu alur sistem yang utuh.

## OpenAPI

Spesifikasi standar untuk mendeskripsikan kontrak REST API.

## Swagger UI

Antarmuka dokumentasi interaktif untuk API berbasis OpenAPI.

---

# DOKUMEN PENUTUP

Sistem **RANCANG BANGUN WEB SERVICE PRESENSI DAN MANAJEMEN ASET LABORATORIUM KOMPUTER MENGGUNAKAN METODE AGILE** ditetapkan sebagai sistem berbasis web yang terdiri atas Web Mahasiswa/User dan Web Admin/Laboran/Tata Usaha, dengan RESTful Web Service sebagai pusat komunikasi antara frontend React dan Supabase PostgreSQL database.

Presensi menggunakan **face verification 1:1 melalui kamera browser**, dengan mahasiswa melakukan login terlebih dahulu sehingga sistem telah mengetahui identitas pengguna. Face verification menggunakan pretrained face model untuk menghasilkan embedding, sedangkan backend bertanggung jawab terhadap enrollment, penyimpanan template, verification, threshold, dan integrasi terhadap attendance.

Manajemen aset menggunakan identitas individual melalui asset code, allocation, return, histori, dan maintenance. Seluruh operasi utama dilakukan melalui REST API dan dilindungi oleh authentication, authorization, validation, transaction, serta pengamanan biometric data.

Pengembangan dilakukan menggunakan Agile dalam tujuh sprint dengan pembagian tanggung jawab dua mahasiswa yang tetap memberikan kontribusi backend pada masing-masing domain. Luaran proyek berupa Web Services/Doc API, Draft Jurnal, Buku Tutorial, dan Poster diturunkan langsung dari sistem yang dibangun.

Parameter yang belum tersedia secara empiris atau institusional tetap dinyatakan secara eksplisit menggunakan penanda `[PERLU VALIDASI]` dan `[DIISI SETELAH PENGUJIAN]` tanpa mengarang data, statistik, atau hasil pengujian.
