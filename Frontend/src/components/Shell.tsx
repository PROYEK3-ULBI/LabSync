import { useState, type ReactNode } from 'react'
import { Icon, type IconName } from './icons'
import { Avatar, IconButton } from './ui'

export type NavItem = { key: string; label: string; icon: IconName }

export function Shell({
  brand,
  roleLabel,
  roleTone,
  nav,
  active,
  onNavigate,
  user,
  onLogout,
  children,
}: {
  brand: string
  roleLabel: string
  roleTone: 'primary' | 'slate'
  nav: NavItem[]
  active: string
  onNavigate: (key: string) => void
  user: { name: string; sub: string }
  onLogout: () => void
  children: ReactNode
}) {
  const [mobileOpen, setMobileOpen] = useState(false)
  const [menuOpen, setMenuOpen] = useState(false)
  const activeItem = nav.find((n) => n.key === active)

  const SidebarInner = (
    <>
      <div className="flex items-center gap-2.5 px-5 h-16 border-b border-[var(--color-border)]">
        <span className="inline-flex items-center justify-center w-9 h-9 rounded-lg bg-[var(--color-primary)] text-white shrink-0">
          <Icon.shield width={20} height={20} />
        </span>
        <div className="min-w-0">
          <p className="font-semibold text-[15px] leading-tight text-[var(--color-text)] truncate">{brand}</p>
          <p className="text-[11px] text-[var(--color-text-3)] leading-tight">Sistem Presensi &amp; Aset Lab</p>
        </div>
      </div>

      <div className="px-3 pt-4 pb-2">
        <span
          className={`inline-flex items-center gap-1.5 text-[11px] font-semibold uppercase tracking-wide px-2.5 py-1 rounded-md ${
            roleTone === 'primary'
              ? 'bg-[var(--color-primary-soft)] text-[var(--color-primary)]'
              : 'bg-[#1d2939] text-white'
          }`}
        >
          {roleLabel}
        </span>
      </div>

      <nav className="flex-1 overflow-y-auto px-3 py-1 space-y-0.5">
        {nav.map((item) => {
          const I = Icon[item.icon]
          const on = item.key === active
          return (
            <button
              key={item.key}
              onClick={() => {
                onNavigate(item.key)
                setMobileOpen(false)
              }}
              className={`w-full flex items-center gap-3 px-3 h-10 rounded-lg text-sm font-medium transition-colors ${
                on
                  ? 'bg-[var(--color-primary-soft)] text-[var(--color-primary)]'
                  : 'text-[var(--color-text-2)] hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text)]'
              }`}
            >
              <I width={18} height={18} className={on ? '' : 'text-[var(--color-text-3)]'} />
              {item.label}
            </button>
          )
        })}
      </nav>

      <div className="p-3 border-t border-[var(--color-border)]">
        <button
          onClick={onLogout}
          className="w-full flex items-center gap-3 px-3 h-10 rounded-lg text-sm font-medium text-[var(--color-text-2)] hover:bg-[var(--color-error-soft)] hover:text-[var(--color-error)] transition-colors"
        >
          <Icon.logout width={18} height={18} />
          Keluar
        </button>
      </div>
    </>
  )

  return (
    <div className="min-h-screen flex bg-[var(--color-bg)]">
      {/* Sidebar — desktop */}
      <aside className="hidden lg:flex w-64 shrink-0 flex-col bg-white border-r border-[var(--color-border)] sticky top-0 h-screen">
        {SidebarInner}
      </aside>

      {/* Sidebar — mobile drawer */}
      {mobileOpen && (
        <div className="lg:hidden fixed inset-0 z-40">
          <div className="absolute inset-0 bg-[#101828]/40" onClick={() => setMobileOpen(false)} />
          <aside className="absolute left-0 top-0 h-full w-64 flex flex-col bg-white border-r border-[var(--color-border)] animate-fade">
            {SidebarInner}
          </aside>
        </div>
      )}

      <div className="flex-1 min-w-0 flex flex-col">
        {/* Topbar */}
        <header className="h-16 shrink-0 bg-white border-b border-[var(--color-border)] flex items-center gap-3 px-4 lg:px-6 sticky top-0 z-30">
          <IconButton label="Menu" className="lg:hidden" onClick={() => setMobileOpen(true)}>
            <Icon.dashboard width={20} height={20} />
          </IconButton>
          <div className="min-w-0">
            <h2 className="font-semibold text-[var(--color-text)] truncate">{activeItem?.label ?? brand}</h2>
          </div>
          <div className="ml-auto flex items-center gap-1.5">
            <IconButton label="Notifikasi">
              <span className="relative">
                <Icon.bell width={20} height={20} />
                <span className="absolute -top-0.5 -right-0.5 w-2 h-2 rounded-full bg-[var(--color-error)] ring-2 ring-white" />
              </span>
            </IconButton>
            <div className="relative">
              <button
                onClick={() => setMenuOpen((v) => !v)}
                className="flex items-center gap-2.5 pl-1.5 pr-2.5 h-11 rounded-lg hover:bg-[var(--color-surface-2)] transition-colors"
              >
                <Avatar name={user.name} size={32} />
                <span className="hidden sm:block text-left">
                  <span className="block text-[13px] font-medium leading-tight text-[var(--color-text)]">{user.name}</span>
                  <span className="block text-[11px] leading-tight text-[var(--color-text-3)]">{user.sub}</span>
                </span>
                <Icon.chevronDown width={16} height={16} className="text-[var(--color-text-3)] hidden sm:block" />
              </button>
              {menuOpen && (
                <>
                  <div className="fixed inset-0 z-10" onClick={() => setMenuOpen(false)} />
                  <div className="absolute right-0 mt-1.5 w-56 bg-white rounded-xl border border-[var(--color-border)] shadow-[var(--shadow-pop)] z-20 py-1.5 animate-fade">
                    <div className="px-3 py-2 border-b border-[var(--color-border)]">
                      <p className="text-sm font-medium text-[var(--color-text)]">{user.name}</p>
                      <p className="text-xs text-[var(--color-text-3)]">{user.sub}</p>
                    </div>
                    <button
                      onClick={() => {
                        onNavigate('profile')
                        setMenuOpen(false)
                      }}
                      className="w-full flex items-center gap-2.5 px-3 py-2 text-sm text-[var(--color-text-2)] hover:bg-[var(--color-surface-2)]"
                    >
                      <Icon.user width={16} height={16} /> Profil Saya
                    </button>
                    <button
                      onClick={onLogout}
                      className="w-full flex items-center gap-2.5 px-3 py-2 text-sm text-[var(--color-error)] hover:bg-[var(--color-error-soft)]"
                    >
                      <Icon.logout width={16} height={16} /> Keluar
                    </button>
                  </div>
                </>
              )}
            </div>
          </div>
        </header>

        <main className="flex-1 p-4 lg:p-6 max-w-[1400px] w-full mx-auto animate-fade" key={active}>
          {children}
        </main>
      </div>
    </div>
  )
}
