import type { SVGProps } from 'react'

// Consistent outline icon set (24x24, 1.75 stroke)
function Base({ children, ...p }: SVGProps<SVGSVGElement> & { children: React.ReactNode }) {
  return (
    <svg
      width="20"
      height="20"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.75"
      strokeLinecap="round"
      strokeLinejoin="round"
      {...p}
    >
      {children}
    </svg>
  )
}

export const Icon = {
  dashboard: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <rect x="3" y="3" width="7" height="9" rx="1.5" />
      <rect x="14" y="3" width="7" height="5" rx="1.5" />
      <rect x="14" y="12" width="7" height="9" rx="1.5" />
      <rect x="3" y="16" width="7" height="5" rx="1.5" />
    </Base>
  ),
  calendar: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <rect x="3" y="4.5" width="18" height="16" rx="2" />
      <path d="M3 9h18M8 2.5v4M16 2.5v4" />
    </Base>
  ),
  check: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M20 6 9 17l-5-5" />
    </Base>
  ),
  checkCircle: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <circle cx="12" cy="12" r="9" />
      <path d="m8.5 12 2.5 2.5 4.5-5" />
    </Base>
  ),
  xCircle: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <circle cx="12" cy="12" r="9" />
      <path d="m9 9 6 6M15 9l-6 6" />
    </Base>
  ),
  alert: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M10.3 3.9 1.8 18a2 2 0 0 0 1.7 3h17a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0Z" />
      <path d="M12 9v4M12 17h.01" />
    </Base>
  ),
  clock: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <circle cx="12" cy="12" r="9" />
      <path d="M12 7v5l3 2" />
    </Base>
  ),
  camera: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M4 8a2 2 0 0 1 2-2h2l1.2-1.6a1 1 0 0 1 .8-.4h4a1 1 0 0 1 .8.4L18 6h2a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2Z" />
      <circle cx="12" cy="12.5" r="3.2" />
    </Base>
  ),
  user: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <circle cx="12" cy="8" r="4" />
      <path d="M4 21c0-4 3.6-6 8-6s8 2 8 6" />
    </Base>
  ),
  users: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <circle cx="9" cy="8" r="3.5" />
      <path d="M2.5 20c0-3.3 3-5 6.5-5s6.5 1.7 6.5 5" />
      <path d="M16 5.2A3.5 3.5 0 0 1 16 12M21.5 20c0-2.6-1.6-4.2-4-4.8" />
    </Base>
  ),
  box: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M12 2.5 21 7v10l-9 4.5L3 17V7Z" />
      <path d="M3 7l9 4.5L21 7M12 11.5V21" />
    </Base>
  ),
  laptop: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <rect x="4" y="5" width="16" height="11" rx="1.5" />
      <path d="M2 20h20" />
    </Base>
  ),
  swap: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M7 4 3 8l4 4M3 8h13M17 20l4-4-4-4M21 16H8" />
    </Base>
  ),
  wrench: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M14.7 6.3a4 4 0 0 0-5.2 5.2L3 18l3 3 6.5-6.5a4 4 0 0 0 5.2-5.2l-2.5 2.5-2.5-.5-.5-2.5Z" />
    </Base>
  ),
  report: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M6 2.5h8l4 4V21a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1V3.5a1 1 0 0 1 1-1Z" />
      <path d="M13 2.5V7h5M8.5 13v4M12 11v6M15.5 15v2" />
    </Base>
  ),
  history: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M3.5 9a9 9 0 1 1-1 5" />
      <path d="M3 4.5V9h4.5M12 8v4.5l3 2" />
    </Base>
  ),
  shield: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M12 2.5 20 5.5V11c0 5-3.4 8.7-8 10.5C7.4 19.7 4 16 4 11V5.5Z" />
      <path d="m9 12 2 2 4-4.5" />
    </Base>
  ),
  logout: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M15 4h3a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2h-3" />
      <path d="M10 8 6 12l4 4M6 12h11" />
    </Base>
  ),
  bell: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M6 9a6 6 0 0 1 12 0c0 5 2 6 2 6H4s2-1 2-6Z" />
      <path d="M10 20a2 2 0 0 0 4 0" />
    </Base>
  ),
  search: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <circle cx="11" cy="11" r="7" />
      <path d="m21 21-4.3-4.3" />
    </Base>
  ),
  plus: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M12 5v14M5 12h14" />
    </Base>
  ),
  filter: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M3 5h18l-7 8v6l-4-2v-4Z" />
    </Base>
  ),
  eye: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7-10-7-10-7Z" />
      <circle cx="12" cy="12" r="3" />
    </Base>
  ),
  eyeOff: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M3 3l18 18M10.6 6.2A9.8 9.8 0 0 1 12 6c6.5 0 10 6 10 6a17 17 0 0 1-3.3 3.9M6.5 8.2A17 17 0 0 0 2 12s3.5 6 10 6a9.5 9.5 0 0 0 3.4-.6" />
      <path d="M9.9 9.9a3 3 0 0 0 4.2 4.2" />
    </Base>
  ),
  edit: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M4 20h4L18.5 9.5a2 2 0 0 0-2.8-2.8L5 17Z" />
      <path d="M14.5 8.5 17 11" />
    </Base>
  ),
  chevronDown: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="m6 9 6 6 6-6" />
    </Base>
  ),
  chevronRight: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="m9 6 6 6-6 6" />
    </Base>
  ),
  arrowRight: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M4 12h16M14 6l6 6-6 6" />
    </Base>
  ),
  arrowLeft: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M20 12H4M10 6l-6 6 6 6" />
    </Base>
  ),
  download: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M12 3v12M7 10l5 5 5-5M4 21h16" />
    </Base>
  ),
  lock: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <rect x="4.5" y="10.5" width="15" height="10" rx="2" />
      <path d="M8 10.5V7a4 4 0 0 1 8 0v3.5" />
    </Base>
  ),
  dots: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <circle cx="5" cy="12" r="1.4" />
      <circle cx="12" cy="12" r="1.4" />
      <circle cx="19" cy="12" r="1.4" />
    </Base>
  ),
  inbox: (p: SVGProps<SVGSVGElement>) => (
    <Base {...p}>
      <path d="M3 13h5l2 3h4l2-3h5" />
      <path d="M5 5h14l2 8v6a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1v-6Z" />
    </Base>
  ),
}

export type IconName = keyof typeof Icon
