import { useState } from 'react'
import Dashboard from './pages/Dashboard'
import Attendance from './pages/Attendance'
import Assets from './pages/Assets'
import Kiosk from './pages/Kiosk'

export type Page = 'dashboard' | 'attendance' | 'assets' | 'kiosk'

export default function App() {
  const [page, setPage] = useState<Page>('dashboard')

  if (page === 'kiosk') return <Kiosk onExit={() => setPage('dashboard')} />

  return (
    <div className="flex h-screen overflow-hidden" style={{ background: '#F4F6FB' }}>
      <Sidebar currentPage={page} onNavigate={setPage} />
      <div className="flex-1 flex flex-col overflow-hidden">
        <Header currentPage={page} />
        <main className="flex-1 overflow-y-auto p-7">
          {page === 'dashboard' && <Dashboard onNavigate={setPage} />}
          {page === 'attendance' && <Attendance />}
          {page === 'assets' && <Assets />}
        </main>
      </div>
    </div>
  )
}

const NAV = [
  {
    id: 'dashboard' as Page, label: 'Dashboard',
    icon: (active: boolean) => (
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke={active ? '#3B5998' : '#94a3b8'} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <rect x="3" y="3" width="7" height="7" rx="1" /><rect x="14" y="3" width="7" height="7" rx="1" />
        <rect x="3" y="14" width="7" height="7" rx="1" /><rect x="14" y="14" width="7" height="7" rx="1" />
      </svg>
    ),
  },
  {
    id: 'attendance' as Page, label: 'Attendance',
    icon: (active: boolean) => (
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke={active ? '#3B5998' : '#94a3b8'} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M9 11l3 3L22 4" /><path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11" />
      </svg>
    ),
  },
  {
    id: 'assets' as Page, label: 'Assets',
    icon: (active: boolean) => (
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke={active ? '#3B5998' : '#94a3b8'} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <rect x="2" y="3" width="20" height="14" rx="2" /><path d="M8 21h8m-4-4v4" />
      </svg>
    ),
  },
  {
    id: 'kiosk' as Page, label: 'Kiosk Mode',
    icon: (active: boolean) => (
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke={active ? '#3B5998' : '#94a3b8'} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <rect x="5" y="2" width="14" height="20" rx="2" /><circle cx="12" cy="17" r="1" fill={active ? '#3B5998' : '#94a3b8'} />
      </svg>
    ),
  },
]

function Sidebar({ currentPage, onNavigate }: { currentPage: Page; onNavigate: (p: Page) => void }) {
  return (
    <aside className="w-60 bg-white flex flex-col shrink-0" style={{ borderRight: '1px solid #E8EDF7' }}>
      {/* Logo */}
      <div className="px-6 py-5" style={{ borderBottom: '1px solid #E8EDF7' }}>
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-xl flex items-center justify-center shrink-0" style={{ background: '#3B5998' }}>
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
              <rect x="2" y="3" width="20" height="14" rx="2" /><path d="M8 21h8m-4-4v4" />
            </svg>
          </div>
          <div>
            <div className="text-[15px] font-bold text-slate-800 leading-tight">CompuLab</div>
            <div className="text-[11px] text-slate-400 tracking-wide">Management System</div>
          </div>
        </div>
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 py-5 space-y-1">
        <p className="text-[10px] font-semibold text-slate-400 uppercase tracking-widest px-3 mb-3">Menu</p>
        {NAV.map(item => {
          const active = currentPage === item.id
          return (
            <button
              key={item.id}
              onClick={() => onNavigate(item.id)}
              className="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-[13.5px] font-medium transition-all duration-150 cursor-pointer group"
              style={{
                background: active ? '#EEF2FB' : 'transparent',
                color: active ? '#3B5998' : '#64748b',
              }}
            >
              <span className="shrink-0">{item.icon(active)}</span>
              <span>{item.label}</span>
              {item.id === 'kiosk' && (
                <span className="ml-auto text-[9px] font-bold px-2 py-0.5 rounded-full" style={{ background: '#3B5998', color: 'white', letterSpacing: '0.05em' }}>LIVE</span>
              )}
              {active && (
                <span className="ml-auto w-1.5 h-5 rounded-full" style={{ background: '#3B5998' }}></span>
              )}
            </button>
          )
        })}
      </nav>

      {/* User */}
      <div className="px-4 py-4 mx-3 mb-4 rounded-xl cursor-pointer hover:bg-slate-50 transition-colors" style={{ border: '1px solid #E8EDF7' }}>
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-full flex items-center justify-center text-[13px] font-bold text-white shrink-0" style={{ background: '#3B5998' }}>AD</div>
          <div className="flex-1 min-w-0">
            <div className="text-[13px] font-semibold text-slate-700 truncate">Admin User</div>
            <div className="text-[11px] text-slate-400 truncate">Lab Supervisor</div>
          </div>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="M9 18l6-6-6-6" />
          </svg>
        </div>
      </div>
    </aside>
  )
}

const PAGE_TITLES: Record<Page, string> = {
  dashboard: 'Dashboard',
  attendance: 'Attendance',
  assets: 'Asset Inventory',
  kiosk: 'Kiosk',
}

function Header({ currentPage }: { currentPage: Page }) {
  const [notifOpen, setNotifOpen] = useState(false)
  return (
    <header className="h-16 bg-white flex items-center px-7 gap-5 shrink-0" style={{ borderBottom: '1px solid #E8EDF7' }}>
      <div>
        <h1 className="text-[17px] font-bold text-slate-800">{PAGE_TITLES[currentPage]}</h1>
        <p className="text-[11px] text-slate-400">Thursday, September 25, 2026</p>
      </div>

      <div className="ml-auto flex items-center gap-3">
        <div className="flex items-center gap-2 rounded-xl px-3.5 py-2 text-[13px]" style={{ background: '#F4F6FB', border: '1px solid #E8EDF7' }}>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="11" cy="11" r="8" /><path d="M21 21l-4.35-4.35" />
          </svg>
          <input type="text" placeholder="Search..." className="bg-transparent text-slate-600 placeholder-slate-400 outline-none w-36" style={{ fontSize: 13 }} />
        </div>

        <div className="relative">
          <button
            onClick={() => setNotifOpen(v => !v)}
            className="relative w-9 h-9 rounded-xl flex items-center justify-center cursor-pointer transition-colors hover:bg-slate-50"
            style={{ border: '1px solid #E8EDF7' }}
          >
            <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="#64748b" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9" /><path d="M13.73 21a2 2 0 01-3.46 0" />
            </svg>
            <span className="absolute top-1.5 right-1.5 w-2 h-2 rounded-full" style={{ background: '#f87171', border: '2px solid white' }}></span>
          </button>
          {notifOpen && (
            <div className="absolute right-0 top-11 w-72 bg-white rounded-2xl shadow-xl z-50 overflow-hidden" style={{ border: '1px solid #E8EDF7' }}>
              <div className="px-4 py-3 border-b border-slate-100">
                <p className="text-[13px] font-semibold text-slate-700">Notifications</p>
              </div>
              {[
                { title: 'AST-0041 marked for repair', time: '9:14 AM', color: '#f59e0b' },
                { title: '3 students checked in late', time: '8:20 AM', color: '#3B5998' },
                { title: 'Webcam AST-0010 is broken', time: '8:05 AM', color: '#ef4444' },
              ].map((n, i) => (
                <div key={i} className="flex items-start gap-3 px-4 py-3 hover:bg-slate-50 cursor-pointer transition-colors border-b border-slate-50">
                  <span className="w-2 h-2 rounded-full mt-1.5 shrink-0" style={{ background: n.color }}></span>
                  <div>
                    <p className="text-[12.5px] text-slate-700 font-medium">{n.title}</p>
                    <p className="text-[11px] text-slate-400 mt-0.5">{n.time}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <button className="w-9 h-9 rounded-full flex items-center justify-center text-[13px] font-bold text-white cursor-pointer hover:opacity-90 transition-opacity" style={{ background: '#3B5998' }}>
          AD
        </button>
      </div>
    </header>
  )
}
