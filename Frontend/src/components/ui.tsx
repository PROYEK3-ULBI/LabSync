import { createContext, useContext, useEffect, useState, type ReactNode } from 'react'
import { Icon } from './icons'

/* ------------------------------------------------------------------ */
/* Button                                                              */
/* ------------------------------------------------------------------ */
type BtnVariant = 'primary' | 'secondary' | 'ghost' | 'danger' | 'success'
export function Button({
  children,
  variant = 'primary',
  size = 'md',
  loading = false,
  disabled = false,
  icon,
  className = '',
  ...rest
}: {
  children?: ReactNode
  variant?: BtnVariant
  size?: 'sm' | 'md' | 'lg'
  loading?: boolean
  icon?: ReactNode
} & React.ButtonHTMLAttributes<HTMLButtonElement>) {
  const variants: Record<BtnVariant, string> = {
    primary:
      'bg-[var(--color-primary)] text-white hover:bg-[var(--color-primary-hover)] shadow-sm',
    secondary:
      'bg-white text-[var(--color-text)] border border-[var(--color-border-strong)] hover:bg-[var(--color-surface-2)]',
    ghost: 'text-[var(--color-text-2)] hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text)]',
    danger: 'bg-[var(--color-error)] text-white hover:brightness-95 shadow-sm',
    success: 'bg-[var(--color-success)] text-white hover:brightness-95 shadow-sm',
  }
  const sizes = {
    sm: 'h-8 px-3 text-[13px] gap-1.5',
    md: 'h-10 px-4 text-sm gap-2',
    lg: 'h-12 px-6 text-[15px] gap-2',
  }
  return (
    <button
      disabled={disabled || loading}
      className={`inline-flex items-center justify-center rounded-lg font-medium transition-colors duration-150 disabled:opacity-50 disabled:pointer-events-none ${variants[variant]} ${sizes[size]} ${className}`}
      {...rest}
    >
      {loading && (
        <span className="w-4 h-4 rounded-full border-2 border-current border-t-transparent animate-spin-slow" />
      )}
      {!loading && icon}
      {children}
    </button>
  )
}

export function IconButton({
  children,
  label,
  className = '',
  ...rest
}: { children: ReactNode; label: string } & React.ButtonHTMLAttributes<HTMLButtonElement>) {
  return (
    <button
      aria-label={label}
      title={label}
      className={`inline-flex items-center justify-center w-9 h-9 rounded-lg text-[var(--color-text-2)] hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text)] transition-colors ${className}`}
      {...rest}
    >
      {children}
    </button>
  )
}

/* ------------------------------------------------------------------ */
/* Badge / Status                                                      */
/* ------------------------------------------------------------------ */
export type Tone = 'neutral' | 'primary' | 'success' | 'warning' | 'error' | 'info'
const toneMap: Record<Tone, string> = {
  neutral: 'bg-[var(--color-surface-2)] text-[var(--color-text-2)] border-[var(--color-border)]',
  primary: 'bg-[var(--color-primary-soft)] text-[var(--color-primary)] border-[#c7d7fe]',
  success: 'bg-[var(--color-success-soft)] text-[var(--color-success)] border-[#abefc6]',
  warning: 'bg-[var(--color-warning-soft)] text-[var(--color-warning)] border-[#fedf89]',
  error: 'bg-[var(--color-error-soft)] text-[var(--color-error)] border-[#fecdca]',
  info: 'bg-[var(--color-info-soft)] text-[var(--color-info)] border-[#b2ddff]',
}

export function Badge({ children, tone = 'neutral', dot = true }: { children: ReactNode; tone?: Tone; dot?: boolean }) {
  const dotColor: Record<Tone, string> = {
    neutral: 'bg-[var(--color-text-3)]',
    primary: 'bg-[var(--color-primary)]',
    success: 'bg-[var(--color-success)]',
    warning: 'bg-[var(--color-warning)]',
    error: 'bg-[var(--color-error)]',
    info: 'bg-[var(--color-info)]',
  }
  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium whitespace-nowrap ${toneMap[tone]}`}
    >
      {dot && <span className={`w-1.5 h-1.5 rounded-full ${dotColor[tone]}`} />}
      {children}
    </span>
  )
}

/* Map arbitrary status strings to tones + label consistently across apps */
const statusTone: Record<string, Tone> = {
  Available: 'success',
  Allocated: 'info',
  Maintenance: 'warning',
  Retired: 'neutral',
  Active: 'success',
  Aktif: 'success',
  Upcoming: 'info',
  'Akan Datang': 'info',
  Closed: 'neutral',
  Ditutup: 'neutral',
  Cancelled: 'error',
  Dibatalkan: 'error',
  Draft: 'neutral',
  Attended: 'success',
  Hadir: 'success',
  'Not Attended': 'error',
  'Tidak Hadir': 'error',
  Present: 'success',
  Completed: 'success',
  Selesai: 'success',
  Pending: 'warning',
  Failed: 'error',
  Gagal: 'error',
  Open: 'error',
  Terbuka: 'error',
  'In Progress': 'warning',
  'Sedang Dikerjakan': 'warning',
  High: 'error',
  Tinggi: 'error',
  Medium: 'warning',
  Sedang: 'warning',
  Low: 'neutral',
  Rendah: 'neutral',
  Verified: 'success',
  Terverifikasi: 'success',
  Terdaftar: 'success',
  'Belum Terdaftar': 'warning',
  Nonaktif: 'error',
}
export function StatusBadge({ status }: { status: string }) {
  return <Badge tone={statusTone[status] ?? 'neutral'}>{status}</Badge>
}

/* ------------------------------------------------------------------ */
/* Card                                                                */
/* ------------------------------------------------------------------ */
export function Card({
  children,
  className = '',
  padded = true,
}: {
  children: ReactNode
  className?: string
  padded?: boolean
}) {
  return (
    <div
      className={`bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-card)] shadow-[var(--shadow-card)] ${padded ? 'p-5' : ''} ${className}`}
    >
      {children}
    </div>
  )
}

export function StatCard({
  label,
  value,
  icon,
  tone = 'primary',
  hint,
}: {
  label: string
  value: ReactNode
  icon: ReactNode
  tone?: Tone
  hint?: string
}) {
  return (
    <Card className="flex items-start justify-between gap-3">
      <div className="min-w-0">
        <p className="text-[13px] text-[var(--color-text-2)] font-medium">{label}</p>
        <p className="mt-2 text-2xl font-semibold font-[var(--font-display)] text-[var(--color-text)] tracking-tight">
          {value}
        </p>
        {hint && <p className="mt-1 text-xs text-[var(--color-text-3)]">{hint}</p>}
      </div>
      <span className={`shrink-0 inline-flex items-center justify-center w-10 h-10 rounded-lg border ${toneMap[tone]}`}>
        {icon}
      </span>
    </Card>
  )
}

/* ------------------------------------------------------------------ */
/* Section header                                                      */
/* ------------------------------------------------------------------ */
export function PageHeader({
  title,
  subtitle,
  actions,
  breadcrumb,
}: {
  title: string
  subtitle?: string
  actions?: ReactNode
  breadcrumb?: string[]
}) {
  return (
    <div className="flex flex-wrap items-end justify-between gap-4 mb-6">
      <div>
        {breadcrumb && (
          <nav className="flex items-center gap-1.5 text-xs text-[var(--color-text-3)] mb-1.5">
            {breadcrumb.map((b, i) => (
              <span key={i} className="flex items-center gap-1.5">
                {i > 0 && <Icon.chevronRight width={12} height={12} />}
                <span className={i === breadcrumb.length - 1 ? 'text-[var(--color-text-2)]' : ''}>{b}</span>
              </span>
            ))}
          </nav>
        )}
        <h1 className="text-xl font-semibold tracking-tight text-[var(--color-text)]">{title}</h1>
        {subtitle && <p className="text-sm text-[var(--color-text-2)] mt-1">{subtitle}</p>}
      </div>
      {actions && <div className="flex items-center gap-2">{actions}</div>}
    </div>
  )
}

/* ------------------------------------------------------------------ */
/* Inputs                                                              */
/* ------------------------------------------------------------------ */
export function Field({
  label,
  required,
  hint,
  error,
  children,
}: {
  label?: string
  required?: boolean
  hint?: string
  error?: string
  children: ReactNode
}) {
  return (
    <label className="block">
      {label && (
        <span className="block text-[13px] font-medium text-[var(--color-text)] mb-1.5">
          {label} {required && <span className="text-[var(--color-error)]">*</span>}
        </span>
      )}
      {children}
      {error ? (
        <span className="mt-1.5 flex items-center gap-1 text-xs text-[var(--color-error)]">
          <Icon.alert width={13} height={13} /> {error}
        </span>
      ) : (
        hint && <span className="block mt-1.5 text-xs text-[var(--color-text-3)]">{hint}</span>
      )}
    </label>
  )
}

const inputBase =
  'w-full h-10 rounded-lg border bg-white px-3 text-sm text-[var(--color-text)] placeholder:text-[var(--color-text-3)] transition-colors focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/25 focus:border-[var(--color-primary)] disabled:bg-[var(--color-surface-2)] disabled:text-[var(--color-text-3)]'

export function Input({
  invalid,
  className = '',
  ...rest
}: { invalid?: boolean } & React.InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      className={`${inputBase} ${invalid ? 'border-[var(--color-error)]' : 'border-[var(--color-border-strong)]'} ${className}`}
      {...rest}
    />
  )
}

export function Select({
  children,
  className = '',
  ...rest
}: React.SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <div className="relative">
      <select
        className={`${inputBase} border-[var(--color-border-strong)] appearance-none pr-9 ${className}`}
        {...rest}
      >
        {children}
      </select>
      <Icon.chevronDown
        width={16}
        height={16}
        className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-[var(--color-text-3)]"
      />
    </div>
  )
}

export function SearchInput({ className = '', ...rest }: React.InputHTMLAttributes<HTMLInputElement>) {
  return (
    <div className={`relative ${className}`}>
      <Icon.search
        width={16}
        height={16}
        className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--color-text-3)]"
      />
      <input
        className={`${inputBase} border-[var(--color-border-strong)] pl-9`}
        placeholder="Cari…"
        {...rest}
      />
    </div>
  )
}

/* ------------------------------------------------------------------ */
/* Table                                                               */
/* ------------------------------------------------------------------ */
export function Table({ children }: { children: ReactNode }) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full border-collapse text-sm">{children}</table>
    </div>
  )
}
export function THead({ cols }: { cols: string[] }) {
  return (
    <thead>
      <tr className="border-b border-[var(--color-border)] bg-[var(--color-surface-2)]">
        {cols.map((c) => (
          <th
            key={c}
            className="text-left font-medium text-xs uppercase tracking-wide text-[var(--color-text-2)] px-4 py-3 whitespace-nowrap"
          >
            {c}
          </th>
        ))}
      </tr>
    </thead>
  )
}
export function TR({ children, onClick }: { children: ReactNode; onClick?: () => void }) {
  return (
    <tr
      onClick={onClick}
      className={`border-b border-[var(--color-border)] last:border-0 ${onClick ? 'cursor-pointer' : ''} hover:bg-[var(--color-surface-2)] transition-colors`}
    >
      {children}
    </tr>
  )
}
export function TD({ children, className = '' }: { children?: ReactNode; className?: string }) {
  return <td className={`px-4 py-3 align-middle text-[var(--color-text)] ${className}`}>{children}</td>
}

export function EmptyState({ title, desc, icon }: { title: string; desc?: string; icon?: ReactNode }) {
  return (
    <div className="flex flex-col items-center justify-center text-center py-16 px-6">
      <span className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-[var(--color-surface-2)] text-[var(--color-text-3)] mb-3">
        {icon ?? <Icon.inbox />}
      </span>
      <p className="font-medium text-[var(--color-text)]">{title}</p>
      {desc && <p className="text-sm text-[var(--color-text-2)] mt-1 max-w-sm">{desc}</p>}
    </div>
  )
}

export function Skeleton({ className = '' }: { className?: string }) {
  return <div className={`animate-pulse bg-[var(--color-surface-2)] rounded ${className}`} />
}

/* Simulated initial-load flag for demo loading states */
export function useSimulatedLoading(ms = 750) {
  const [loading, setLoading] = useState(true)
  useEffect(() => {
    const t = setTimeout(() => setLoading(false), ms)
    return () => clearTimeout(t)
  }, [ms])
  return loading
}

/* Skeleton rows for a table with `cols` columns */
export function SkeletonRows({ rows = 5, cols }: { rows?: number; cols: number }) {
  return (
    <tbody>
      {Array.from({ length: rows }).map((_, r) => (
        <tr key={r} className="border-b border-[var(--color-border)] last:border-0">
          {Array.from({ length: cols }).map((_, c) => (
            <td key={c} className="px-4 py-3.5">
              <Skeleton className={`h-4 ${c === 0 ? 'w-2/3' : 'w-4/5'}`} />
            </td>
          ))}
        </tr>
      ))}
    </tbody>
  )
}

export function Pagination({
  page,
  pages,
  onPage,
  total,
}: {
  page: number
  pages: number
  onPage: (p: number) => void
  total: number
}) {
  return (
    <div className="flex items-center justify-between px-4 py-3 border-t border-[var(--color-border)] text-sm text-[var(--color-text-2)]">
      <span>{total} data</span>
      <div className="flex items-center gap-1">
        <Button variant="secondary" size="sm" disabled={page <= 1} onClick={() => onPage(page - 1)}>
          Sebelumnya
        </Button>
        <span className="px-3 text-[var(--color-text)]">
          {page} / {pages}
        </span>
        <Button variant="secondary" size="sm" disabled={page >= pages} onClick={() => onPage(page + 1)}>
          Berikutnya
        </Button>
      </div>
    </div>
  )
}

/* ------------------------------------------------------------------ */
/* Modal                                                               */
/* ------------------------------------------------------------------ */
export function Modal({
  open,
  onClose,
  title,
  children,
  footer,
  size = 'md',
}: {
  open: boolean
  onClose: () => void
  title: string
  children: ReactNode
  footer?: ReactNode
  size?: 'sm' | 'md' | 'lg'
}) {
  useEffect(() => {
    if (!open) return
    const h = (e: KeyboardEvent) => e.key === 'Escape' && onClose()
    window.addEventListener('keydown', h)
    return () => window.removeEventListener('keydown', h)
  }, [open, onClose])
  if (!open) return null
  const w = { sm: 'max-w-md', md: 'max-w-lg', lg: 'max-w-2xl' }[size]
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-[#101828]/40 backdrop-blur-[1px]" onClick={onClose} />
      <div className={`relative w-full ${w} bg-white rounded-xl shadow-[var(--shadow-pop)] animate-fade`}>
        <div className="flex items-center justify-between px-5 py-4 border-b border-[var(--color-border)]">
          <h3 className="font-semibold text-[var(--color-text)]">{title}</h3>
          <IconButton label="Tutup" onClick={onClose}>
            <Icon.xCircle width={18} height={18} />
          </IconButton>
        </div>
        <div className="px-5 py-4 max-h-[70vh] overflow-y-auto">{children}</div>
        {footer && (
          <div className="flex items-center justify-end gap-2 px-5 py-4 border-t border-[var(--color-border)] bg-[var(--color-surface-2)] rounded-b-xl">
            {footer}
          </div>
        )}
      </div>
    </div>
  )
}

/* ------------------------------------------------------------------ */
/* Toast                                                               */
/* ------------------------------------------------------------------ */
type Toast = { id: number; msg: string; tone: Tone }
const ToastCtx = createContext<(msg: string, tone?: Tone) => void>(() => {})
export const useToast = () => useContext(ToastCtx)

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([])
  const push = (msg: string, tone: Tone = 'success') => {
    const id = Date.now() + Math.random()
    setToasts((t) => [...t, { id, msg, tone }])
    setTimeout(() => setToasts((t) => t.filter((x) => x.id !== id)), 3500)
  }
  return (
    <ToastCtx.Provider value={push}>
      {children}
      <div className="fixed bottom-6 right-6 z-[60] flex flex-col gap-2 w-80">
        {toasts.map((t) => (
          <div
            key={t.id}
            className={`flex items-start gap-2.5 rounded-lg border px-4 py-3 shadow-[var(--shadow-pop)] text-sm animate-fade ${toneMap[t.tone]}`}
          >
            {t.tone === 'error' ? <Icon.xCircle width={18} height={18} /> : <Icon.checkCircle width={18} height={18} />}
            <span className="text-[var(--color-text)] leading-snug">{t.msg}</span>
          </div>
        ))}
      </div>
    </ToastCtx.Provider>
  )
}

/* ------------------------------------------------------------------ */
/* Alert                                                               */
/* ------------------------------------------------------------------ */
export function Alert({ tone = 'info', title, children }: { tone?: Tone; title?: string; children: ReactNode }) {
  const I =
    tone === 'error' ? Icon.xCircle : tone === 'warning' ? Icon.alert : tone === 'success' ? Icon.checkCircle : Icon.shield
  return (
    <div className={`flex gap-3 rounded-lg border px-4 py-3 text-sm ${toneMap[tone]}`}>
      <I width={18} height={18} className="shrink-0 mt-0.5" />
      <div>
        {title && <p className="font-medium text-[var(--color-text)]">{title}</p>}
        <div className="text-[var(--color-text-2)]">{children}</div>
      </div>
    </div>
  )
}

/* Avatar with initials */
export function Avatar({ name, size = 36 }: { name: string; size?: number }) {
  const initials = name
    .split(' ')
    .slice(0, 2)
    .map((n) => n[0])
    .join('')
    .toUpperCase()
  return (
    <span
      className="inline-flex items-center justify-center rounded-full bg-[var(--color-primary-soft)] text-[var(--color-primary)] font-semibold shrink-0"
      style={{ width: size, height: size, fontSize: size * 0.38 }}
    >
      {initials}
    </span>
  )
}

/* Mono code chip for asset codes / identifiers */
export function Code({ children }: { children: ReactNode }) {
  return (
    <span className="font-mono text-[13px] font-medium text-[var(--color-text)] bg-[var(--color-surface-2)] border border-[var(--color-border)] rounded px-1.5 py-0.5">
      {children}
    </span>
  )
}

export function Tabs({ tabs, active, onChange }: { tabs: string[]; active: string; onChange: (t: string) => void }) {
  return (
    <div className="flex items-center gap-1 border-b border-[var(--color-border)] mb-5 overflow-x-auto">
      {tabs.map((t) => (
        <button
          key={t}
          onClick={() => onChange(t)}
          className={`px-3.5 py-2.5 text-sm font-medium border-b-2 -mb-px whitespace-nowrap transition-colors ${
            active === t
              ? 'border-[var(--color-primary)] text-[var(--color-primary)]'
              : 'border-transparent text-[var(--color-text-2)] hover:text-[var(--color-text)]'
          }`}
        >
          {t}
        </button>
      ))}
    </div>
  )
}

/* ------------------------------------------------------------------ */
/* Step indicator                                                      */
/* ------------------------------------------------------------------ */
export function StepIndicator({ steps, current }: { steps: string[]; current: number }) {
  return (
    <ol className="flex items-center w-full">
      {steps.map((s, i) => {
        const done = i < current
        const active = i === current
        return (
          <li key={s} className={`flex items-center ${i < steps.length - 1 ? 'flex-1' : ''}`}>
            <div className="flex items-center gap-2 shrink-0">
              <span
                className={`inline-flex items-center justify-center w-7 h-7 rounded-full text-xs font-semibold border transition-colors ${
                  done
                    ? 'bg-[var(--color-success)] border-[var(--color-success)] text-white'
                    : active
                      ? 'bg-[var(--color-primary)] border-[var(--color-primary)] text-white'
                      : 'bg-white border-[var(--color-border-strong)] text-[var(--color-text-3)]'
                }`}
              >
                {done ? <Icon.check width={14} height={14} /> : i + 1}
              </span>
              <span className={`text-[13px] font-medium hidden sm:block ${active ? 'text-[var(--color-text)]' : 'text-[var(--color-text-2)]'}`}>
                {s}
              </span>
            </div>
            {i < steps.length - 1 && (
              <span className={`flex-1 h-0.5 mx-3 rounded ${done ? 'bg-[var(--color-success)]' : 'bg-[var(--color-border)]'}`} />
            )}
          </li>
        )
      })}
    </ol>
  )
}

/* ------------------------------------------------------------------ */
/* Timeline                                                            */
/* ------------------------------------------------------------------ */
export function Timeline({
  items,
}: {
  items: { title: string; time: string; desc?: string; tone?: Tone }[]
}) {
  return (
    <ol className="relative">
      {items.map((it, i) => (
        <li key={i} className="relative pl-8 pb-5 last:pb-0">
          {i < items.length - 1 && <span className="absolute left-[9px] top-4 bottom-0 w-px bg-[var(--color-border)]" />}
          <span
            className={`absolute left-0 top-1 w-[18px] h-[18px] rounded-full border-2 border-white ${
              toneMap[it.tone ?? 'primary'].split(' ')[0]
            } ring-1 ring-[var(--color-border)]`}
          />
          <div className="flex flex-wrap items-center gap-2">
            <p className="text-sm font-medium text-[var(--color-text)]">{it.title}</p>
            <span className="text-xs text-[var(--color-text-3)]">{it.time}</span>
          </div>
          {it.desc && <p className="text-sm text-[var(--color-text-2)] mt-0.5">{it.desc}</p>}
        </li>
      ))}
    </ol>
  )
}
