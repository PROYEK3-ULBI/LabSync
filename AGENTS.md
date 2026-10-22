# LabSync

Web lab attendance + asset management. User-facing UI/API text must be Bahasa Indonesia.

## Boundaries

- `Frontend/`: React 19 + Vite 8 + Tailwind v4 prototype; read `Frontend/AGENTS.md` before frontend edits.
- `backend/`: Go 1.27 Fiber API. Entrypoint `cmd/server/main.go`; routes wired in `internal/router/router.go`; config loader reads `.env` from current working directory, so run commands from `backend/` for `backend/.env`.
- `face_service/`: private FastAPI biometric service. Frontend must never call it; backend calls it with `X-Internal-API-Key`.
- `database/migrations/000001_init_schema.{up,down}.sql`: canonical PostgreSQL schema for 10 tables.
- Check `docs/api-contract.md`, `docs/openapi.yaml`, and `docs/security-decisions.md` before API/security changes.

## Commands

- Frontend, from `Frontend/`: `pnpm dev`, `pnpm build`, `pnpm test`, `pnpm format`. Use pnpm; both `pnpm-lock.yaml` and `package-lock.json` exist.
- Frontend single tests: `pnpm test -- src/lib/name.test.ts` does not apply; script is fixed to `node --test src/lib/*.test.ts`.
- Backend, from `backend/`: `go run cmd/server/main.go`, `go test -v ./...`. Focused package/test: `go test -v ./internal/router -run TestName`.
- Seed dev data, from `backend/` after `DATABASE_URL` is set: `go run cmd/seed/main.go`. Seed is idempotent and blocked when `APP_ENV=production`.
- Face service, from `face_service/`: `python3.11 -m venv venv`, `source venv/bin/activate`, `pip install -r requirements.txt`, `uvicorn main:app --host 0.0.0.0 --port 8001 --reload`, `pytest -v`.
- Database, from repo root: `psql -U postgres -d labsync -f database/migrations/000001_init_schema.up.sql`; rollback with `.down.sql`.

## Gotchas

- Frontend dev server uses strict port `8443` and proxies `/api` to `http://localhost:8000`.
- Frontend active screens live in `src/apps/AdminApp.tsx` and `src/apps/StudentApp.tsx`; `src/pages/` is unused dead prototype code.
- Frontend has no router; screen changes are state keys in `src/App.tsx` / app components.
- Frontend alias is `@/* -> ./src/*`; Tailwind v4 theme lives in `src/index.css`, no `tailwind.config.*`.
- `pnpm exec tsc --noEmit` is known dirty because unused `src/pages/*` imports missing old prototype deps; use `pnpm build` for active app verification unless fixing dead pages.
- `pnpm format` runs `oxfmt` across frontend files; avoid it when only checking unrelated work.
- Backend can boot without DB in development but reports degraded health; production requires `DATABASE_URL`, changed `JWT_SECRET`, and `BIOMETRIC_ENCRYPTION_KEY` length >= 32.
- Supabase transaction pooler requires pgx simple protocol; `internal/repository/db.go` already sets `DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol`.
- Face service falls back to `MockFaceEngine` if DeepFace fails at startup; integration tests may skip when model/libs/weights are missing. Use Python 3.10-3.12, not 3.14+.

## Security Invariants

- Never store raw face images to disk, object storage, or DB. Store only encrypted embeddings in `face_template_samples`.
- Face enrollment requires exactly 3 poses: `FRONT`, `LEFT`, `RIGHT`.
- Face verification is 1:1 against user-owned enrolled samples, not 1:N identification.
- Do not set `IS_CALIBRATED=true` until empirical calibration on target population is done; uncalibrated face service forces `verified=false`.
- Backend owns auth. JWT uses HttpOnly cookie; frontend API calls must include credentials and use `/api/v1/auth/me` for auth status.
