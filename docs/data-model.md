# Data Model & Relational Schema

Implementasi untuk Supabase PostgreSQL (10 Tabel).

## 1. `roles`
Menyimpan definisi peran.
- `id` (UUID, PK)
- `code` (VARCHAR, Unique) - e.g., 'STUDENT', 'ADMIN'
- `name` (VARCHAR)

## 2. `users`
- `id` (UUID, PK)
- `role_id` (UUID, FK ke `roles`)
- `identifier` (VARCHAR, Unique) - NIM/NIP
- `email` (VARCHAR, Unique)
- `password_hash` (VARCHAR)
- `status` (VARCHAR) - 'ACTIVE', 'INACTIVE'
- `created_at`, `updated_at` (TIMESTAMP)

## 3. `face_templates`
Menyimpan metadata enrollment aktif.
- `id` (UUID, PK)
- `user_id` (UUID, FK ke `users`)
- `model_name` (VARCHAR) - e.g., 'Facenet512'
- `model_version` (VARCHAR)
- `consent_version` (VARCHAR) - Menyimpan jejak persetujuan biometrik saat enrollment.
- `status` (VARCHAR) - 'ACTIVE', 'REVOKED'
- `enrolled_at`, `updated_at` (TIMESTAMP)
- **Constraint:** Partial Unique Index pada `(user_id)` WHERE `status = 'ACTIVE'` (Satu user hanya punya maksimal 1 set aktif).

## 4. `face_template_samples`
Menyimpan sampel pose individu.
- `id` (UUID, PK)
- `face_template_id` (UUID, FK ke `face_templates` ON DELETE CASCADE)
- `pose` (VARCHAR) - 'FRONT', 'LEFT', 'RIGHT'
- `embedding_ciphertext` (BYTEA) - Hasil enkripsi dari array of floats.
- `quality_status` (VARCHAR) - 'VALID'
- `created_at` (TIMESTAMP)
- **Constraint:** Unique Index pada `(face_template_id, pose)` (Satu set hanya boleh memiliki satu FRONT, satu LEFT, satu RIGHT).

## 5. `lab_sessions`
- `id` (UUID, PK)
- `code` (VARCHAR, Unique)
- `title` (VARCHAR)
- `session_date` (DATE)
- `start_at` (TIMESTAMP WITH TIME ZONE)
- `end_at` (TIMESTAMP WITH TIME ZONE)
- `status` (VARCHAR) - 'SCHEDULED', 'ACTIVE', 'CLOSED', 'CANCELLED'

## 6. `attendance`
- `id` (UUID, PK)
- `user_id` (UUID, FK ke `users`)
- `session_id` (UUID, FK ke `lab_sessions`)
- `checked_in_at` (TIMESTAMP WITH TIME ZONE)
- `checked_out_at` (TIMESTAMP WITH TIME ZONE, Nullable)
- `verification_status` (VARCHAR) - 'VERIFIED_FACE'
- **Constraint:** Unique Index pada `(user_id, session_id)` (Satu user hanya 1 presensi per sesi).

## 7. `assets`
- `id` (UUID, PK)
- `asset_code` (VARCHAR, Unique)
- `asset_name` (VARCHAR)
- `status` (VARCHAR) - 'AVAILABLE', 'ALLOCATED', 'MAINTENANCE', 'RETIRED'

## 8. `asset_allocations`
- `id` (UUID, PK)
- `asset_id` (UUID, FK ke `assets`)
- `user_id` (UUID, FK ke `users`)
- `session_id` (UUID, FK ke `lab_sessions`)
- `allocated_at` (TIMESTAMP WITH TIME ZONE)
- `returned_at` (TIMESTAMP WITH TIME ZONE, Nullable)
- **Constraint:** Partial Unique Index pada `(asset_id)` WHERE `returned_at IS NULL` (Mencegah double allocation aktif).
- **Operasi:** Insert allocation dan Update status `assets` menjadi 'ALLOCATED' dibungkus dalam **Satu Database Transaction**.

## 9. `maintenance`
- `id` (UUID, PK)
- `asset_id` (UUID, FK ke `assets`)
- `reported_by` (UUID, FK ke `users`)
- `status` (VARCHAR) - 'OPEN', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED'
- `priority` (VARCHAR) - 'LOW', 'MEDIUM', 'HIGH'
- `created_at`, `updated_at` (TIMESTAMP)

## 10. `audit_logs`
- `id` (UUID, PK)
- `user_id` (UUID, FK ke `users`)
- `action` (VARCHAR) - e.g., 'ASSET_ALLOCATED', 'MAINTENANCE_CREATED'
- `entity_type` (VARCHAR)
- `entity_id` (UUID)
- `created_at` (TIMESTAMP)

## Catatan Aturan Transaksional
- Modul aset (alokasi dan pengembalian) **WAJIB** menggunakan SQL Transactions untuk mengunci baris (Row Lock) demi menghindari *race condition* jika dua admin mengalokasikan aset yang sama secara bersamaan.