# Implementation Status

Last checked: 2026-10-12

## Current Scope

- Roadmap source: `opencode_prompt_roadmap_backend_lab_revisi_face_multi_angle.md`.
- Active PRD/diagram sources: `prd_app_proyek3_revisi_face_enrollment_multi_angle.md`, `diagram-proyek3_plantuml_revisi_face_multi_angle.md`.
- Old sources `prd_app_proyek3.md` and `diagram-proyek3-plantuml.md` are deleted in the current worktree and appear replaced by the revised files.

## Prompt 14 Checkpoint

- Student frontend now has API client coverage for auth/session restoration, face status/enrollment, sessions, attendance check-in/out/history, and own allocations.
- Student face flow uses `navigator.mediaDevices.getUserMedia()` only inside camera flow, captures still frames via canvas, and sends images to the Go API.
- Enrollment sends one multipart request with `front_image`, `left_image`, `right_image`, `consent_agreed`, and `re_enroll` to `/api/v1/face/enroll`.
- Check-in sends `session_id` plus `image` to `/api/v1/attendance/check-in`; UI waits for backend success before showing success.
- Check-out sends `{ "session_id": "..." }` to `/api/v1/attendance/check-out` and disables duplicate submit while pending.
- UI distinguishes camera permission/device problems, invalid face/no face/multi-face backend errors, verification failure, server/network failure, and success messages.
- Frontend does not display embeddings, permanent raw face images, similarity score, or face-service internals.

## Prompt 15 Readiness

- Prompt 15 admin frontend finalization is in progress/completed for the current checkpoint.
- Admin screens now use Go REST API data directly for users, sessions, attendance, assets, allocations, maintenance, and reports CSV export.
- Admin create flows are wired for users, sessions, assets, allocations, and maintenance tickets using backend payloads/IDs.
- Allocation return and maintenance completion call backend mutation endpoints.
- Admin UI still uses the existing Shell and shared UI components; no second admin app was created.
- Supabase transaction pooler compatibility fixed in `backend/internal/repository/db.go` by disabling pgx prepared statement cache via simple protocol.
- Face service now loads `face_service/.env` in `face_service/app/config.py`; previously it only used process environment, causing backend-to-face-service `401 Unauthorized` when `INTERNAL_API_KEY` differed from runtime env.
- `INTERNAL_API_KEY` was rotated and synced in both `backend/.env` and `face_service/.env` after accidental local terminal exposure during debugging.

## Verification Run

- Supabase dev migration `psql "$DATABASE_URL" -f database/migrations/000001_init_schema.up.sql`: PASS.
- Supabase dev schema verification: PASS, 10 expected public tables and `ADMIN`/`STUDENT` roles present.
- Backend health against Supabase dev: PASS, `{"status":"healthy","services":{"database":"connected"}}`.
- Development seed `cd backend && go run cmd/seed/main.go`: PASS. Seeded 2 users, 2 sessions, 3 assets, and 1 maintenance ticket.
- `cd backend && go test -v ./...`: PASS.
- `cd Frontend && pnpm test`: PASS, 7 tests.
- `cd Frontend && pnpm test`: PASS, 8 tests after admin API endpoint coverage was added.
- `cd Frontend && pnpm build`: ran with existing `vite.config.ts` warning: `__dirname` should use `import.meta.dirname`.
- `cd Frontend && pnpm exec tsc --noEmit`: NOT CLEAN because unused/dead `src/pages/*` prototype code imports missing `recharts` and invalid old prototype props. Active `src/apps/*` build path is covered by Vite build.
- Admin API smoke with seeded admin cookie: PASS for `auth/me`, `dashboard/admin`, `users`, `sessions`, `attendance`, `assets`, `allocations`, `maintenance`, all report endpoints, and `audit-logs`.
- Student API smoke with seeded student cookie: PASS for `auth/me`, `dashboard/student`, `sessions`, `sessions/active`, `face/status`, `attendance/history`, and `assets/my-allocations`.
- `cd Frontend && pnpm format`: ran earlier, but global formatter touched many unrelated files; unrelated changes were restored. Run with care because `oxfmt` formats all frontend files.

## Not Run

- Face service `pytest -v`: NOT RUN in this checkpoint.
- Local browser E2E with real DB + backend + face service: NOT RUN.
- Face service health at `http://localhost:8001/health`: NOT RUN/PASS because service was not running during smoke check.
- Face service health later verified: `GET http://localhost:8001/health` returned `200 OK` with ArcFace metadata.
- Internal face auth smoke after `.env` loader fix: direct extract request with synced key no longer returned `401`; dummy text file returned expected `400 INVALID_IMAGE`.
- Backend-to-face enrollment smoke after restart: login as `MHS-DEMO-001`, `POST /api/v1/face/enroll` with three dummy `image/jpeg` files now reaches face service and returns image validation `400`; no `401`, `503`, or multipart `application/octet-stream` regression.
- Browser face enrollment with real camera still needs manual retest; services are running on backend `8000`, face service `8001`, frontend `8443`.
- PostgreSQL migration apply: RUN against configured Supabase dev database.
- Supabase remote setup/deploy: NOT RUN and must remain manual/explicit.

## Manual Steps Before Prompt 15

1. Start services in separate terminals: face service on `8001`, backend on `8000`, frontend on `8443`.
2. Use the dev credentials printed by `go run cmd/seed/main.go`; do not commit or paste them into docs.
3. Browser-check admin login and at least one list endpoint before starting Prompt 15.

## Risks / Notes

- The backend can start in development without DB, but admin integration needs real data to verify behavior.
- Prompt 15 must not query Supabase/PostgreSQL directly from React; use Go REST only.
- Do not run migrations against Supabase remote without explicit confirmation.
- Do not commit `.env`, DB passwords, Supabase service-role keys, raw face images, embeddings, or test images containing real people.
