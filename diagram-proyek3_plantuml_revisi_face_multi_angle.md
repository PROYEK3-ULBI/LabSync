# Dokumentasi Diagram (PlantUML) — Proyek 3

Semua 12 diagram di bawah ini ditulis dalam kode **PlantUML**. Satu situs saja untuk semuanya, tidak perlu Claude atau software tambahan.

## Cara membuka

1. Buka **https://www.planttext.com** (gratis, tanpa akun).
2. Hapus contoh kode bawaan di kotak kiri.
3. Copy salah satu blok kode di bawah (mulai dari `@startuml` sampai `@enduml`), paste ke kotak kiri.
4. Diagram otomatis muncul di kanan.
5. Klik tombol **Refresh** kalau tidak otomatis render, lalu klik ikon **download** di kanan atas hasil gambar untuk simpan sebagai PNG/SVG.

Alternatif resmi: **https://www.plantuml.com/plantuml/uml/** — fungsinya sama, tampilannya saja yang sedikit berbeda.

> **Catatan khusus untuk diagram #10 (Gantt chart):** sintaksnya memakai `@startgantt` (bukan `@startuml`), fitur yang relatif baru. Kalau planttext.com menolak/error, coba buka di **https://www.plantuml.com/plantuml/uml/** — mesinnya biasanya versi lebih baru dan lebih mendukung gantt.

---

## 1. ERD Database

Sepuluh tabel sesuai migration `database/migrations/000001_init_schema.up.sql`.

```plantuml
@startuml
entity ROLES {
  * id : uuid <<PK>>
  --
  code : string
  name : string
}
entity USERS {
  * id : uuid <<PK>>
  --
  role_id : uuid <<FK>>
  identifier : string
  email : string
  full_name : string
  password_hash : string
  status : string
  created_at : timestamp
  updated_at : timestamp
}
entity FACE_TEMPLATES {
  * id : uuid <<PK>>
  --
  user_id : uuid <<FK>>
  model_name : string
  model_version : string
  consent_version : string
  consented_at : timestamp
  status : string
  enrolled_at : timestamp
  updated_at : timestamp
}
entity FACE_TEMPLATE_SAMPLES {
  * id : uuid <<PK>>
  --
  face_template_id : uuid <<FK>>
  pose : string
  embedding_ciphertext : bytea
  quality_status : string
  created_at : timestamp
}

entity LAB_SESSIONS {
  * id : uuid <<PK>>
  --
  code : string
  title : string
  course_code : string
  location : string
  session_date : date
  start_at : timestamp
  end_at : timestamp
  status : string
  created_by : uuid <<FK>>
}
entity ATTENDANCE {
  * id : uuid <<PK>>
  --
  user_id : uuid <<FK>>
  session_id : uuid <<FK>>
  checked_in_at : timestamp
  checked_out_at : timestamp
  verification_status : string
}
entity ASSETS {
  * id : uuid <<PK>>
  --
  asset_code : string
  asset_name : string
  category : string
  serial_number : string
  location : string
  status : string
}
entity ASSET_ALLOCATIONS {
  * id : uuid <<PK>>
  --
  asset_id : uuid <<FK>>
  user_id : uuid <<FK>>
  session_id : uuid <<FK>>
  allocated_at : timestamp
  returned_at : timestamp
  condition_before : string
  condition_after : string
}
entity MAINTENANCE {
  * id : uuid <<PK>>
  --
  asset_id : uuid <<FK>>
  reported_by : uuid <<FK>>
  issue_type : string
  description : text
  status : string
  priority : string
  started_at : timestamp
  completed_at : timestamp
}
entity AUDIT_LOGS {
  * id : uuid <<PK>>
  --
  user_id : uuid <<FK>>
  action : string
  entity_type : string
  entity_id : uuid
  status : string
  metadata : jsonb
}

ROLES ||--o{ USERS : has
USERS ||--o{ FACE_TEMPLATES : owns
FACE_TEMPLATES ||--o{ FACE_TEMPLATE_SAMPLES : contains
USERS ||--o{ ATTENDANCE : records
LAB_SESSIONS ||--o{ ATTENDANCE : contains
USERS ||--o{ ASSET_ALLOCATIONS : allocated
LAB_SESSIONS ||--o{ ASSET_ALLOCATIONS : contains
ASSETS ||--o{ ASSET_ALLOCATIONS : allocated_to
ASSETS ||--o{ MAINTENANCE : undergoes
USERS ||--o{ MAINTENANCE : reports
USERS ||--o{ AUDIT_LOGS : generates
USERS ||--o{ LAB_SESSIONS : creates
@enduml
```

---

## 2. Sequence Diagram — Multi-Angle Enrollment dan Face Verification 1:1

```plantuml
@startuml
actor Mahasiswa as M
participant Browser as B
participant "Go REST API" as A
participant "Python Face Service" as F
database "Supabase PostgreSQL" as D

== Enrollment ==
M -> B : Mulai enrollment + setuju pemberitahuan
B -> M : Panduan pose FRONT
M -> B : Capture FRONT
B -> M : Panduan pose LEFT
M -> B : Capture LEFT
B -> M : Panduan pose RIGHT
M -> B : Capture RIGHT
B -> A : POST enroll dengan 3 file pose
A -> F : Validasi dan ekstrak embedding FRONT
F --> A : FRONT embedding + quality result
A -> F : Validasi dan ekstrak embedding LEFT
F --> A : LEFT embedding + quality result
A -> F : Validasi dan ekstrak embedding RIGHT
F --> A : RIGHT embedding + quality result
alt Semua pose valid
  A -> D : Transaction: simpan enrollment set + 3 encrypted samples
  D --> A : Commit berhasil
  A --> B : Enrollment aktif; raw images dibuang
else Ada pose tidak valid
  A --> B : Enrollment ditolak; minta capture ulang pose terkait
end

== Verification ==
M -> B : Buka presensi dan pilih session
B -> A : Request check-in dengan live frame
A -> D : Ambil enrollment set user login
D --> A : FRONT/LEFT/RIGHT samples
A -> F : Ekstrak live embedding + compare 1:1
F --> A : Decision result
alt Verified
  A -> D : Simpan attendance dengan waktu server
  A --> B : Verification dan check-in berhasil
else Not verified
  A --> B : Verification gagal; tanpa attendance
end
@enduml
```


## 3. Sequence Diagram — Attendance (Check-in & Check-out)

Flow aktual endpoint `/api/v1/attendance/check-in` dan `/api/v1/attendance/check-out`.

```plantuml
@startuml
actor Mahasiswa as M
participant Browser as B
participant "Go REST API" as A
participant "Attendance Service" as AS
participant "Python Face Service" as F
database "PostgreSQL" as D

== Check-in ==
M -> B : Pilih session + capture wajah
B -> A : POST /attendance/check-in (session_id + image)
A -> AS : Validasi role STUDENT + session aktif
AS -> D : Ambil attendance + active face template user
D --> AS : Data sesi, attendance, 3 encrypted samples
AS -> F : Verify 1:1 live image vs 3 samples user
F --> AS : verified / distance / decision_rule
alt verified
  AS -> D : Insert attendance, checked_in_at = server time
  A --> B : Check-in berhasil
else not verified
  A --> B : Verifikasi wajah gagal; attendance tidak dibuat
end

== Check-out ==
M -> B : Check-out
B -> A : POST /api/v1/attendance/check-out (session_id)
A -> AS : Cari attendance aktif milik user
AS -> D : Set checked_out_at = server time
A --> B : Check-out berhasil
@enduml
```

---

## 4. Activity Diagram — Login

Flow aktual `/api/v1/auth/login` dengan JWT HttpOnly cookie.

```plantuml
@startuml
start
:Buka halaman login;
:Isi kredensial;
:POST /api/v1/auth/login;
:Validasi user (cek password hash);
if (Password benar?) then (tidak)
  :401 Unauthorized;
  stop
else (ya)
endif
if (Akun aktif?) then (tidak)
  :403 Forbidden;
  stop
else (ya)
endif
:Buat JWT;
:Set HttpOnly cookie token;
:Dashboard;
stop
@enduml
```

---

## 5. Activity Diagram — Face Enrollment Multi-Angle

```plantuml
@startuml
start
:Login dan setujui pemrosesan biometrik;
:Request akses kamera;
if (Izin kamera?) then (tidak)
  :Tampilkan error kamera;
  stop
else (ya)
endif
:Capture pose FRONT;
:Capture pose LEFT;
:Capture pose RIGHT;
:POST /api/v1/face/enroll dengan 3 file;
:Backend kirim tiap pose ke face_service internal;
if (Semua pose valid dan 1 wajah?) then (ya)
  :Ekstrak embedding tiap pose;
  :Enkripsi dan simpan 3 sampel + metadata model;
  :Buang raw image dari memori;
  :Enrollment aktif;
else (tidak)
  :Tampilkan error pose/citra;
endif
stop
@enduml
```

## 6. Component Diagram — Architecture (Layered)

Arsitektur aktual: React frontend memanggil Go Fiber API; Go API memanggil face_service internal.

```plantuml
@startuml
component "Web Mahasiswa" as WebM
component "Web Admin" as WebA
component "Go Fiber API\n(router, middleware, handler)" as API
component "Service Layer\n(auth, user, face, attendance, asset)" as BLL
component "FastAPI face_service\n(DeepFace ArcFace, detector)" as FVL
component "Repository Layer\n(pgx, transaction, mapping)" as DAL
database "PostgreSQL Database" as DB

WebM --> API
WebA --> API
API --> BLL
BLL --> FVL : internal HTTP + X-Internal-API-Key
BLL --> DAL
DAL --> DB
@enduml
```

---

## 7. Deployment Diagram

Deployment target/logical runtime.

```plantuml
@startuml
actor "Browser Mahasiswa" as BM
actor "Browser Admin" as BA
node "Static Web Hosting" as Hosting
node "API Gateway (HTTPS)" as Gateway
node "Go Fiber Backend" as Backend
node "FastAPI Face Service\n(private, DeepFace model lokal)" as FaceSvc
database "PostgreSQL Database\n(managed)" as DB

BM --> Hosting
BA --> Hosting
Hosting --> Gateway
Gateway --> Backend
Backend --> FaceSvc
Backend --> DB
@enduml
```

---

## 8. State Diagram — Status Asset

Status asset sesuai constraint database dan service allocation.

```plantuml
@startuml
[*] --> AVAILABLE : asset didaftarkan
AVAILABLE --> ALLOCATED : allocation dibuat
ALLOCATED --> AVAILABLE : return, kondisi baik
ALLOCATED --> MAINTENANCE : return, kondisi rusak
AVAILABLE --> MAINTENANCE : dilaporkan rusak
MAINTENANCE --> AVAILABLE : maintenance selesai
MAINTENANCE --> RETIRED : tidak layak digunakan
RETIRED --> [*]
@enduml
```

---

## 9. State Diagram — Status Maintenance

Status maintenance sesuai constraint database.

```plantuml
@startuml
[*] --> OPEN : laporan dibuat
OPEN --> IN_PROGRESS : mulai diproses
IN_PROGRESS --> COMPLETED : perbaikan selesai
OPEN --> CANCELLED : dibatalkan
IN_PROGRESS --> CANCELLED : dibatalkan
COMPLETED --> [*]
CANCELLED --> [*]
@enduml
```

---

## 10. Gantt Chart — Rencana 7 Sprint

PRD 14.9. Tanggal hanya contoh — sesuaikan dengan kalender akademik. **Pakai `@startgantt`, bukan `@startuml`** (lihat catatan di atas kalau error di planttext.com).

```plantuml
@startgantt
Project starts 2026-10-01
[Sprint 1: Requirement & arsitektur] lasts 7 days
[Sprint 2: Database, backend, auth] lasts 7 days
[Sprint 2: Database, backend, auth] starts at [Sprint 1: Requirement & arsitektur]'s end
[Sprint 3: Face enrollment & verification] lasts 7 days
[Sprint 3: Face enrollment & verification] starts at [Sprint 2: Database, backend, auth]'s end
[Sprint 4: Manajemen aset & alokasi] lasts 7 days
[Sprint 4: Manajemen aset & alokasi] starts at [Sprint 3: Face enrollment & verification]'s end
[Sprint 5: Histori, maintenance, laporan] lasts 7 days
[Sprint 5: Histori, maintenance, laporan] starts at [Sprint 4: Manajemen aset & alokasi]'s end
[Sprint 6: Web Mahasiswa & Web Admin] lasts 7 days
[Sprint 6: Web Mahasiswa & Web Admin] starts at [Sprint 5: Histori, maintenance, laporan]'s end
[Sprint 7: Pengujian & deployment] lasts 7 days
[Sprint 7: Pengujian & deployment] starts at [Sprint 6: Web Mahasiswa & Web Admin]'s end
@endgantt
```

---

## 11. Class Diagram — Backend Service

Struktur service backend aktual (ringkas).

```plantuml
@startuml
class AuthenticationService {
  +login()
  +ValidateToken()
  +logout()
}
class FaceVerificationService {
  +Enroll()
  +GetStatus()
  +Revoke()
}
class AttendanceService {
  +CheckIn()
  +CheckOut()
  +History()
}
class AssetService {
  +List()
  +Create()
  +UpdateStatus()
}
class AllocationService {
  +Create()
  +Return()
  +AssetHistory()
}
class MaintenanceService {
  +Create()
  +Update()
  +List()
}
class ReportingService {
  +generateReport()
  +exportCSV()
}
class Repository {
  +query()
  +transaction()
}

AttendanceService ..> AuthenticationService : gunakan sesi login
AttendanceService ..> FaceVerificationService : verifikasi wajah
AssetService ..> AuthenticationService : cek otorisasi
AllocationService ..> AssetService : update status aset
MaintenanceService ..> AssetService : cek/update status aset
ReportingService ..> AttendanceService : ambil data
ReportingService ..> AssetService : ambil data
ReportingService ..> MaintenanceService : ambil data
AuthenticationService --> Repository
FaceVerificationService --> Repository
AttendanceService --> Repository
AssetService --> Repository
AllocationService --> Repository
MaintenanceService --> Repository
@enduml
```

---

## 12. Activity Diagram — Swimlane Presensi

Pembagian tanggung jawab Mahasiswa vs Sistem. PlantUML punya sintaks swimlane native (`|NamaLane|`) — lebih rapi daripada versi gambar tangan sebelumnya.

```plantuml
@startuml
|Mahasiswa|
start
:Login;
|Sistem|
:Validasi kredensial;
|Mahasiswa|
:Pilih sesi & buka kamera;
|Sistem|
:Validasi sesi aktif;
|Mahasiswa|
:Arahkan wajah, capture;
|Sistem|
:Ambil 3 sampel enrollment user;
:Face verification 1:1 via face_service;
if (Verified?) then (tidak)
  :Presensi ditolak;
  stop
else (ya)
endif
:Simpan attendance (checked_in_at = waktu server);
|Mahasiswa|
:Terima notifikasi berhasil;
stop
@enduml
```

---

*Dibuat untuk Proyek 3 — Keyla S. & Ridwan H.*
