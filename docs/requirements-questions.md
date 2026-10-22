# Requirements & Questions

Analysis of `prd_app_proyek3_revisi_face_enrollment_multi_angle.md` and `diagram-proyek3_plantuml_revisi_face_multi_angle.md`.

## Business Rules & Prioritization
- **Must Have:** Login, RBAC (Admin/Student), Face Enrollment (3-pose: FRONT, LEFT, RIGHT), Face Verification 1:1, Sessions, Attendance (Check-in/out), Asset Mgmt, Allocations, Maintenance, Reports, API, DB.
- **Should Have:** CSV Export, Report Filters, Audit Log, Re-enrollment workflow.
- **Security:** Passwords hashed. Raw biometrics dropped. Only encrypted embeddings stored. 

## Conflicts & Ambiguities

| Issue | Source of Conflict | Impact | Recommendation / Default | Human Decision Needed? |
| :--- | :--- | :--- | :--- | :--- |
| **Token Delivery** | PRD 5.1 mentions "cookie aman" (HttpOnly), but also mentions JWT struct. | Frontend needs to handle auth state correctly. If HttpOnly cookie, frontend can't read JWT payload directly. | Use HttpOnly cookie for token. Add a `GET /auth/me` endpoint to fetch user profile/roles. | No. Standard practice. |
| **Biometric Consent Tracking** | PRD 6.4: `face_templates` has "metadata persetujuan". ERD Diagram 1: `FACE_TEMPLATES` lacks a `consent_version` field. | Cannot definitively track which terms the user agreed to. | Add `consent_version` (string) to `FACE_TEMPLATES` schema. | Disetujui. |
| **Multi-Sample Verif Strategy** | PRD 6.5: "Aturan pembandingan yang menggabungkan beberapa sampel... harus ditetapkan". | How does the Python service decide `verified = true` against 3 samples? (e.g., 1 of 3 passes threshold, or average distance). | Aturan "Minimum Distance" (1 dari 3 lolos). Jika ada 1 sampel < threshold, hasil true. | Disetujui. |
| **Timezone/Timestamps** | PRD 5.7: `checked_in_at = Server Time`. | If server time is UTC, UI might display wrong time if not localized. | DB stores standard `timestamp with time zone` (UTC). Frontend formats to `Asia/Jakarta`. | No. Standard practice. |
| **API Framework (Backend)** | PRD Architecture: Go + Fiber. Diagram 7 (Deployment): "FastAPI Backend". | Which framework handles the main API? Go or Python? | Follow PRD 8.5/`AGENTS.md`: Go+Fiber for main API, Python+FastAPI for Face Service. Diagram 7 is slightly inaccurate. | No. Trust `AGENTS.md` and PRD text. |