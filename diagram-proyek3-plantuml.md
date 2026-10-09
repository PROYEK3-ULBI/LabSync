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

Sembilan tabel inti beserta relasinya (PRD Part 10).

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
  password_hash : string
  status : string
}
entity FACE_TEMPLATES {
  * id : uuid <<PK>>
  --
  user_id : uuid <<FK>>
  model_name : string
  model_version : string
  embedding_ciphertext : bytea
  status : string
}
entity LAB_SESSIONS {
  * id : uuid <<PK>>
  --
  code : string
  title : string
  session_date : date
  start_at : timestamp
  end_at : timestamp
  status : string
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
}
entity MAINTENANCE {
  * id : uuid <<PK>>
  --
  asset_id : uuid <<FK>>
  reported_by : uuid <<FK>>
  status : string
  priority : string
}
entity AUDIT_LOGS {
  * id : uuid <<PK>>
  --
  user_id : uuid <<FK>>
  action : string
  entity_type : string
}

ROLES ||--o{ USERS : has
USERS ||--o| FACE_TEMPLATES : owns
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

## 2. Sequence Diagram — Face Verification 1:1

PRD 6.6: alur verifikasi wajah saat presensi.

```plantuml
@startuml
actor Mahasiswa as M
participant Browser as B
participant "REST API" as A
participant "Face Service" as F
database Database as D

M -> B : Buka halaman presensi
B -> A : GET status session
A --> B : Session valid
B -> M : Tampilkan kamera
M -> B : Arahkan wajah, capture frame
B -> A : POST /face/verify (image, session_id)
A -> F : Deteksi wajah & buat embedding
F --> A : Embedding
A -> D : Ambil template user login
D --> A : Template wajah
A -> F : Hitung similarity vs threshold
F --> A : verified = true/false
alt Verified
  A -> D : Buat attendance (checked_in_at)
  A --> B : 200 verified, attendance dibuat
else Tidak verified
  A --> B : 200 verified = false
end
@enduml
```

---

## 3. Sequence Diagram — Attendance (Check-in & Check-out)

PRD 7.3–7.5.

```plantuml
@startuml
actor Mahasiswa as M
participant "REST API" as A
participant "Session Service" as S
participant "Face Service" as F
database Database as D

== Check-in ==
M -> A : Pilih session
A -> S : Validasi session aktif
S --> A : Session aktif
A -> D : Cek attendance existing
D --> A : Belum ada
A -> F : Jalankan face verification 1:1
F --> A : verified = true
A -> D : Insert attendance, checked_in_at = server time
D --> A : Tersimpan
A --> M : Check-in berhasil

== Check-out ==
M -> A : Check-out
A -> D : Ambil attendance aktif
D --> A : Ada, belum check-out
A -> D : Set checked_out_at = server time
D --> A : Diperbarui
A --> M : Check-out berhasil
@enduml
```

---

## 4. Activity Diagram — Login

PRD 7.1. Ditulis sebagai activity diagram (setara flowchart di PlantUML).

```plantuml
@startuml
start
:Buka halaman login;
:Isi kredensial;
:POST /auth/login;
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
:Buat access token;
:Set httponly cookie;
:Dashboard;
stop
@enduml
```

---

## 5. Activity Diagram — Face Enrollment

PRD 7.2. Ada loop retry kalau jumlah wajah terdeteksi tidak tepat satu.

```plantuml
@startuml
start
:Buka face enrollment;
:Request akses kamera;
if (Izin kamera?) then (tidak)
  :Error kamera;
  stop
else (ya)
endif
while (Tepat satu wajah terdeteksi?) is (tidak, ulangi)
  :Capture citra wajah;
endwhile (ya)
:Face embedding (model pretrained);
:Simpan template;
:Enrollment berhasil;
stop
@enduml
```

---

## 6. Component Diagram — Architecture (Layered)

PRD Part 8: enam lapisan backend.

```plantuml
@startuml
component "Web Mahasiswa" as WebM
component "Web Admin" as WebA
component "API Layer\n(routing, autentikasi, validasi)" as API
component "Business Logic Layer\n(aturan domain & validasi bisnis)" as BLL
component "Face Verification Layer\n(deteksi, embedding, threshold)" as FVL
component "Data Access Layer\n(query, transaction, mapping)" as DAL
database "PostgreSQL Database" as DB

WebM --> API
WebA --> API
API --> BLL
BLL --> FVL
FVL --> DAL
DAL --> DB
@enduml
```

---

## 7. Deployment Diagram

PRD Part 22.

```plantuml
@startuml
actor "Browser Mahasiswa" as BM
actor "Browser Admin" as BA
node "Static Web Hosting" as Hosting
node "API Gateway (HTTPS)" as Gateway
node "FastAPI Backend" as Backend
node "Face Verification Service\n(pretrained model lokal)" as FaceSvc
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

PRD 10.7, 5.10–5.12.

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

PRD 10.9.

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

PRD Part 9: struktur service backend dan ketergantungannya.

```plantuml
@startuml
class AuthenticationService {
  +login()
  +validateToken()
  +checkRole()
}
class FaceVerificationService {
  +enrollFace()
  +verifyFace()
  +computeEmbedding()
}
class AttendanceService {
  +checkIn()
  +checkOut()
}
class AssetService {
  +allocateAsset()
  +returnAsset()
  +getAssetHistory()
}
class MaintenanceService {
  +reportDamage()
  +updateStatus()
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
MaintenanceService ..> AssetService : update status aset
ReportingService ..> AttendanceService : ambil data
ReportingService ..> AssetService : ambil data
ReportingService ..> MaintenanceService : ambil data
AuthenticationService --> Repository
FaceVerificationService --> Repository
AttendanceService --> Repository
AssetService --> Repository
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
:Face verification 1:1;
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
