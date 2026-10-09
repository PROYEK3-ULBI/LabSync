# Frontend

React 19 + Vite 8 + Tailwind CSS v4 prototype.

## Commands

Run from `Frontend/`:

```bash
pnpm dev       # Vite dev server on port 8443
pnpm build     # production build
pnpm preview   # serve production build
pnpm format    # oxfmt formatter
```

## Toolchain

- Node 22, pnpm 10.34.3 (pinned in `.mise.toml`)
- Both `pnpm-lock.yaml` and `package-lock.json` exist — use **pnpm**
- No ESLint, no Prettier, no test framework configured
- No CI/CD

## Path alias

`@/*` → `./src/*` (configured in both `tsconfig.json` and `vite.config.ts`)

## Key files

- `src/main.tsx` — React entrypoint, mounts `App.tsx`
- `src/App.tsx` — root component: auth state + role (`student`/`admin`) drives which app renders
- `src/apps/AdminApp.tsx` — admin dashboard (all screens inlined, ~1155 lines)
- `src/apps/StudentApp.tsx` — student dashboard (all screens inlined, ~674 lines)
- `src/components/ui.tsx` — full design system (Button, Badge, Card, Table, Modal, Toast, etc.)
- `src/components/FaceFlow.tsx` — simulated face capture (no real biometric)
- `src/components/Login.tsx` — demo login (password "salah" = error, anything else = success)
- `src/lib/data.ts` — all mock/fixture data
- `src/index.css` — global styles + Tailwind v4 theme (fonts: Inter, IBM Plex Sans, JetBrains Mono)

## Architecture notes

- **No router** — navigation is state-based via string keys, no react-router
- **No API calls** — everything uses hardcoded mock data from `src/lib/data.ts`
- `src/pages/` directory exists (Dashboard, Attendance, Assets, Kiosk) but is **unused** — the `src/apps/` files contain all active screens
- `FaceFlow` component fakes camera capture with SVG silhouette animation

## Styling

Tailwind CSS v4 via `@tailwindcss/vite` plugin. No `tailwind.config.*` or PostCSS config needed. Theme customization goes in `src/index.css` after the `@import 'tailwindcss'` line.

## Conventions

- All UI text: Bahasa Indonesia
- Use Tailwind utility classes in JSX; global/theme CSS in `src/index.css`
