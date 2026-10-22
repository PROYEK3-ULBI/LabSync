# Repo Audit

**Date:** 2026-10-09
**Status:** Frontend prototype only. No backend exists.

## 1. Stack & Toolchain
- **Package Manager:** `pnpm` (ignoring root `package-lock.json`).
- **Framework:** React 19 + Vite 8.
- **Styling:** Tailwind CSS v4 (`@tailwindcss/vite`).
- **Routing:** None. State-based navigation via string keys in `App.tsx`.
- **State Management:** Raw React `useState`. No Redux/Zustand.
- **Data Fetching:** None. Uses hardcoded mock data from `src/lib/data.ts`.
- **Backend/DB:** None found.

## 2. Directory Structure
```
/
├── Frontend/
│   ├── src/
│   │   ├── apps/
│   │   │   ├── AdminApp.tsx (Inlines 15+ admin screens)
│   │   │   └── StudentApp.tsx (Inlines 10+ student screens)
│   │   ├── components/
│   │   │   ├── FaceFlow.tsx (SVG animation mockup)
│   │   │   ├── Login.tsx
│   │   │   └── ui.tsx (Design system)
│   │   ├── lib/
│   │   │   └── data.ts (Mock data fixtures)
│   │   ├── pages/ (Unused)
│   │   ├── index.css (Tailwind theme)
│   │   └── main.tsx
├── AGENTS.md (Agent instructions)
├── prd_app_proyek3_revisi_face_enrollment_multi_angle.md (PRD)
└── diagram-proyek3_plantuml_revisi_face_multi_angle.md (Diagrams)
```

## 3. Existing UI Components (Mocked)
- **Student Screens (`StudentApp.tsx`):** Dashboard, Sessions, Attendance, MyAssets, History, AttendanceDetail, Checkout, Profile, Enroll.
- **Admin Screens (`AdminApp.tsx`):** Dashboard, Users, SessionMgmt, AttendanceMonitoring, AssetMgmt, AssetDetail, UserDetail, SessionDetail, Allocations, AllocationDetail, MaintenanceMgmt, MaintenanceDetail, Reports, AuditLog, Profile.

## 4. Key Data Read/Write Points
- **Current State:** All reads/writes happen strictly in-memory against local component state or `lib/data.ts`. No HTTP requests are made.

## 5. Summary
Frontend is structurally a single-page prototype heavily inlining components in two massive files (`AdminApp.tsx`, `StudentApp.tsx`). To connect to a real backend, the state-based router needs HTTP fetching (e.g., React Query or native `fetch`), and the mocked functions (`go`, `onView`) must be wired to real API endpoints.