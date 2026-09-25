import { useState, useEffect } from 'react'

const RECENT = [
  { name: 'Maria Santos',    id: 'STU-001', time: '07:58 AM', pc: 'PC-A-04' },
  { name: 'Juan dela Cruz',  id: 'STU-002', time: '08:01 AM', pc: 'PC-A-07' },
  { name: 'Carlos Bautista', id: 'STU-004', time: '07:55 AM', pc: 'PC-B-02' },
]

export default function Kiosk({ onExit }: { onExit: () => void }) {
  const [now, setNow]         = useState(new Date())
  const [phase, setPhase]     = useState<'idle' | 'scanning' | 'success'>('idle')
  const [lastEntry, setLastEntry] = useState<typeof RECENT[0] | null>(null)
  const [recentLog, setRecentLog] = useState(RECENT)

  useEffect(() => {
    const t = setInterval(() => setNow(new Date()), 1000)
    return () => clearInterval(t)
  }, [])

  function handleScan() {
    if (phase !== 'idle') return
    setPhase('scanning')
    setTimeout(() => {
      const entry = { name: 'Liza Mendoza', id: 'STU-005', time: new Date().toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: true }), pc: 'PC-A-11' }
      setLastEntry(entry)
      setRecentLog(prev => [entry, ...prev.slice(0, 2)])
      setPhase('success')
      setTimeout(() => {
        setPhase('idle')
        setLastEntry(null)
      }, 3500)
    }, 2000)
  }

  const timeStr = now.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: true })
  const dateStr = now.toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric', year: 'numeric' })

  return (
    <div className="min-h-screen flex flex-col items-center justify-center relative overflow-hidden select-none"
      style={{ background: 'linear-gradient(160deg, #03091A 0%, #060D22 50%, #040B1C 100%)' }}>

      {/* Subtle grid */}
      <div className="absolute inset-0 pointer-events-none" style={{
        backgroundImage: 'linear-gradient(rgba(0,212,255,0.04) 1px, transparent 1px), linear-gradient(90deg, rgba(0,212,255,0.04) 1px, transparent 1px)',
        backgroundSize: '48px 48px',
      }} />

      {/* Center glow */}
      <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
        <div style={{ width: 700, height: 700, borderRadius: '50%', background: 'radial-gradient(circle, rgba(0,212,255,0.06) 0%, transparent 65%)' }} />
      </div>

      {/* Admin exit */}
      <button onClick={onExit}
        className="absolute top-5 left-6 flex items-center gap-2 cursor-pointer transition-all hover:opacity-80"
        style={{ color: '#334155', fontSize: 12, fontFamily: 'JetBrains Mono', border: '1px solid #1e293b', borderRadius: 10, padding: '6px 12px', background: 'rgba(255,255,255,0.03)' }}>
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="M19 12H5M12 19l-7-7 7-7"/>
        </svg>
        Admin
      </button>

      {/* Live badge */}
      <div className="absolute top-5 right-6 flex items-center gap-2"
        style={{ color: '#00D4FF', fontSize: 11, fontFamily: 'JetBrains Mono', letterSpacing: '0.15em', border: '1px solid rgba(0,212,255,0.2)', borderRadius: 8, padding: '5px 12px', background: 'rgba(0,212,255,0.05)' }}>
        <span className="w-1.5 h-1.5 rounded-full" style={{ background: '#00D4FF', animation: 'pulse 2s infinite' }}></span>
        SYSTEM ONLINE
      </div>

      {/* Clock */}
      <div className="text-center mb-2 z-10">
        <div style={{ fontSize: 68, fontWeight: 700, fontFamily: 'JetBrains Mono', color: 'white', letterSpacing: '-0.02em', lineHeight: 1, textShadow: '0 0 60px rgba(0,212,255,0.2)' }}>
          {timeStr}
        </div>
        <div style={{ fontSize: 14, fontFamily: 'JetBrains Mono', color: '#4B6CB7', marginTop: 8, letterSpacing: '0.1em' }}>
          {dateStr}
        </div>
      </div>

      <div className="z-10 my-7" style={{ width: 160, height: 1, background: 'linear-gradient(90deg, transparent, rgba(0,212,255,0.3), transparent)' }} />

      {/* Scanner */}
      <button onClick={handleScan} disabled={phase !== 'idle'}
        className="z-10 relative group"
        style={{ cursor: phase === 'idle' ? 'pointer' : 'default' }}>
        <div className="flex flex-col items-center justify-center transition-all duration-300"
          style={{
            width: 224, height: 224, borderRadius: 20,
            border: `2px solid ${phase === 'success' ? '#10b981' : phase === 'scanning' ? '#00D4FF' : 'rgba(0,212,255,0.35)'}`,
            background: phase === 'success' ? 'rgba(16,185,129,0.07)' : phase === 'scanning' ? 'rgba(0,212,255,0.05)' : 'rgba(0,212,255,0.03)',
            boxShadow: phase === 'success' ? '0 0 40px rgba(16,185,129,0.15)' : phase === 'scanning' ? '0 0 40px rgba(0,212,255,0.12)' : 'none',
            position: 'relative', overflow: 'hidden',
          }}>

          {/* Corner marks */}
          {[
            { top: -1, left: -1, bt: 'border-t-2', bl: 'border-l-2' },
            { top: -1, right: -1, bt: 'border-t-2', bl: 'border-r-2' },
            { bottom: -1, left: -1, bt: 'border-b-2', bl: 'border-l-2' },
            { bottom: -1, right: -1, bt: 'border-b-2', bl: 'border-r-2' },
          ].map((c, i) => (
            <div key={i} style={{ position: 'absolute', width: 22, height: 22, ...c }} className={`${c.bt} ${c.bl}`}
              style2={{ borderColor: phase === 'success' ? '#10b981' : '#00D4FF' }}
            />
          ))}

          {/* Scan line */}
          {phase === 'scanning' && (
            <div style={{
              position: 'absolute', left: 0, right: 0, height: 2,
              background: 'linear-gradient(90deg, transparent, #00D4FF, transparent)',
              animation: 'scanLine 1.4s ease-in-out infinite',
            }} />
          )}

          {phase === 'success' ? (
            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 10 }}>
              <div style={{ width: 52, height: 52, borderRadius: '50%', background: 'rgba(16,185,129,0.15)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#10b981" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M5 13l4 4L19 7"/>
                </svg>
              </div>
              <div style={{ textAlign: 'center' }}>
                <p style={{ color: '#10b981', fontFamily: 'JetBrains Mono', fontSize: 13, fontWeight: 600 }}>Logged In</p>
                <p style={{ color: '#64748b', fontSize: 11, fontFamily: 'JetBrains Mono', marginTop: 2 }}>{lastEntry?.pc}</p>
              </div>
            </div>
          ) : phase === 'scanning' ? (
            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 12 }}>
              <div style={{ width: 44, height: 44, position: 'relative' }}>
                <div style={{ position: 'absolute', inset: 0, borderRadius: '50%', border: '2px solid rgba(0,212,255,0.15)' }} />
                <div style={{ position: 'absolute', inset: 0, borderRadius: '50%', border: '2px solid #00D4FF', borderTopColor: 'transparent', animation: 'spin 0.8s linear infinite' }} />
              </div>
              <p style={{ color: '#00D4FF', fontFamily: 'JetBrains Mono', fontSize: 12 }}>Reading ID...</p>
            </div>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 14 }}>
              {/* QR icon */}
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3,1fr)', gap: 5, opacity: 0.5 }}>
                {[1,1,0,1,0,1,0,1,1].map((v, i) => (
                  <div key={i} style={{ width: 16, height: 16, borderRadius: 3, background: v ? '#00D4FF' : 'transparent', border: v ? 'none' : '1px solid rgba(0,212,255,0.2)' }} />
                ))}
              </div>
              <p style={{ color: 'rgba(0,212,255,0.5)', fontFamily: 'JetBrains Mono', fontSize: 10, letterSpacing: '0.2em', textTransform: 'uppercase' }}>
                Tap to demo
              </p>
            </div>
          )}
        </div>
      </button>

      {/* Instruction text */}
      <div className="z-10 mt-8 text-center">
        <p style={{ color: '#e2e8f0', fontSize: 17, fontWeight: 600, letterSpacing: '0.01em' }}>
          {phase === 'success' && lastEntry ? `Welcome, ${lastEntry.name}!` : 'Please scan your ID card'}
        </p>
        <p style={{ color: '#3B5998', fontFamily: 'JetBrains Mono', fontSize: 12, marginTop: 6, letterSpacing: '0.15em', textTransform: 'uppercase' }}>
          {phase === 'success' ? 'Attendance recorded' : 'to mark attendance'}
        </p>
      </div>

      {/* Recent log */}
      <div className="z-10 mt-8 w-full max-w-md px-4">
        <p style={{ color: '#1e3a5f', fontFamily: 'JetBrains Mono', fontSize: 10, letterSpacing: '0.2em', textTransform: 'uppercase', textAlign: 'center', marginBottom: 10 }}>
          Recent Check-ins
        </p>
        <div className="flex flex-col gap-2">
          {recentLog.map((s, i) => (
            <div key={i} className="flex items-center justify-between px-4 py-3 rounded-xl"
              style={{ background: 'rgba(255,255,255,0.03)', border: '1px solid rgba(255,255,255,0.05)' }}>
              <div className="flex items-center gap-3">
                <div style={{ width: 30, height: 30, borderRadius: 8, background: 'rgba(0,212,255,0.1)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                  <span style={{ color: '#00D4FF', fontFamily: 'JetBrains Mono', fontSize: 11, fontWeight: 700 }}>
                    {s.name.split(' ').map(n => n[0]).join('')}
                  </span>
                </div>
                <div>
                  <p style={{ color: '#e2e8f0', fontSize: 12.5, fontWeight: 600 }}>{s.name}</p>
                  <p style={{ color: '#334155', fontFamily: 'JetBrains Mono', fontSize: 10.5, marginTop: 1 }}>{s.id} · {s.pc}</p>
                </div>
              </div>
              <span style={{ color: '#00D4FF', fontFamily: 'JetBrains Mono', fontSize: 11.5 }}>{s.time}</span>
            </div>
          ))}
        </div>
      </div>

      {/* Bottom session info */}
      <div className="absolute bottom-5" style={{ display: 'flex', gap: 20, color: '#1e3a5f', fontFamily: 'JetBrains Mono', fontSize: 11, letterSpacing: '0.1em' }}>
        <span>AM Block · 08:00 – 11:00</span>
        <span style={{ color: '#0f2240' }}>·</span>
        <span>Lab A — 45 seats</span>
        <span style={{ color: '#0f2240' }}>·</span>
        <span>44 / 45 present</span>
      </div>

      <style>{`
        @keyframes scanLine {
          0%   { top: -2px; opacity: 0; }
          10%  { opacity: 1; }
          90%  { opacity: 1; }
          100% { top: calc(100% + 2px); opacity: 0; }
        }
        @keyframes spin { to { transform: rotate(360deg); } }
        @keyframes pulse {
          0%, 100% { opacity: 1; }
          50% { opacity: 0.3; }
        }
      `}</style>
    </div>
  )
}
