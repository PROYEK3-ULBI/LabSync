# Implementation Plan

Based on the Agile MVP scope (Must Haves) and current frontend prototype.

## Phase 1: Database & Go Backend Skeleton
**Goal:** Initialize database schema and basic Go Fiber API.
- **Dependencies:** None.
- **Target Files:**
  - `database/init.sql` (Supabase schema, 10 tables per PlantUML)
  - `backend/cmd/server/main.go`
  - `backend/internal/config/`
  - `backend/internal/router/`
- **DoD:** Database running in Supabase. Go Fiber server boots, connects to DB, and returns `200 OK` on `/health`.

## Phase 2: Auth & User Management API
**Goal:** Implement login, JWT HttpOnly cookies, and RBAC.
- **Dependencies:** Phase 1.
- **Target Files:**
  - `backend/internal/model/user.go`
  - `backend/internal/repository/user_repo.go`
  - `backend/internal/service/auth_service.go`
  - `backend/internal/handler/auth_handler.go`
  - `backend/internal/middleware/auth.go`
- **DoD:** Users can login. APIs are protected by role-based middleware.

## Phase 3: Face Verification Microservice (Python)
**Goal:** Standalone FastAPI service for DeepFace extraction and comparison.
- **Dependencies:** None.
- **Target Files:**
  - `face_service/main.py`
  - `face_service/routes/extract.py`
  - `face_service/routes/verify.py`
  - `face_service/services/deepface_service.py`
  - `face_service/requirements.txt`
- **DoD:** Service accepts image, returns embedding. Service accepts live image + 3 embeddings, returns boolean verif result. Raw images are dropped.

## Phase 4: Face Enrollment & Verification Integration (Go)
**Goal:** Go API orchestrates calls to Face Service and encrypts/stores embeddings.
- **Dependencies:** Phase 2, Phase 3.
- **Target Files:**
  - `backend/internal/repository/face_repo.go` (Handles `face_templates` and `face_template_samples`)
  - `backend/internal/service/face_orchestrator.go`
- **DoD:** Go API securely stores 3 encrypted embeddings (FRONT, LEFT, RIGHT). Verif endpoint orchestrates check against stored embeddings.

## Phase 5: Session & Attendance API
**Goal:** Core business logic for sessions and check-in/out.
- **Dependencies:** Phase 4.
- **Target Files:**
  - `backend/internal/repository/session_repo.go`
  - `backend/internal/repository/attendance_repo.go`
  - `backend/internal/service/attendance_service.go`
- **DoD:** Check-in requires successful face verification (from Phase 4). Check-out records timestamp. Prevents duplicates.

## Phase 6: Asset & Maintenance API
**Goal:** Asset lifecycle management.
- **Dependencies:** Phase 2.
- **Target Files:**
  - `backend/internal/repository/asset_repo.go` (Assets, Allocations, Maintenance)
  - `backend/internal/service/asset_service.go`
- **DoD:** Admin can allocate assets to users. State machine transitions enforced (AVAILABLE -> ALLOCATED -> MAINTENANCE).

## Phase 7: Frontend Integration
**Goal:** Wire `Frontend/src/` components to the real Go API.
- **Dependencies:** Phases 1-6.
- **Target Files:**
  - `Frontend/src/lib/api.ts` (New: Axios/fetch wrapper)
  - `Frontend/src/components/Login.tsx` (Wire POST /auth/login)
  - `Frontend/src/apps/StudentApp.tsx` (Wire Attendance, Enroll flows)
  - `Frontend/src/apps/AdminApp.tsx` (Wire Assets, Users)
- **DoD:** `data.ts` mock is deleted. All UI data flows from HTTP requests. Multi-angle face flow correctly POSTs 3 images to backend.