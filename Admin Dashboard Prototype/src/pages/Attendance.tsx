import { useState } from 'react'

const records = [
  { id: 'STU-001', name: 'Maria Santos',     timeIn: '07:58 AM', timeOut: '11:02 AM', pc: 'PC-A-04', status: 'Present' },
  { id: 'STU-002', name: 'Juan dela Cruz',   timeIn: '08:01 AM', timeOut: '11:00 AM', pc: 'PC-A-07', status: 'Present' },
  { id: 'STU-003', name: 'Ana Reyes',        timeIn: '—',        timeOut: '—',        pc: '—',       status: 'Absent'  },
  { id: 'STU-004', name: 'Carlos Bautista',  timeIn: '07:55 AM', timeOut: '10:58 AM', pc: 'PC-B-02', status: 'Present' },
  { id: 'STU-005', name: 'Liza Mendoza',     timeIn: '08:15 AM', timeOut: '11:00 AM', pc: 'PC-A-11', status: 'Present' },
  { id: 'STU-006', name: 'Rico Villanueva',  timeIn: '—',        timeOut: '—',        pc: '—',       status: 'Absent'  },
  { id: 'STU-007', name: 'Trisha Gomez',     timeIn: '08:03 AM', timeOut: '11:01 AM', pc: 'PC-C-06', status: 'Present' },
  { id: 'STU-008', name: 'Paolo Aquino',     timeIn: '07:59 AM', timeOut: '—',        pc: 'PC-B-09', status: 'Present' },
  { id: 'STU-009', name: 'Jenny Ong',        timeIn: '—',        timeOut: '—',        pc: '—',       status: 'Absent'  },
  { id: 'STU-010', name: 'Marco Tan',        timeIn: '08:00 AM', timeOut: '10:55 AM', pc: 'PC-A-03', status: 'Present' },
  { id: 'STU-011', name: 'Rina Castillo',    timeIn: '08:07 AM', timeOut: '11:00 AM', pc: 'PC-C-01', status: 'Present' },
  { id: 'STU-012', name: 'Dino Ramos',       timeIn: '07:52 AM', timeOut: '10:59 AM', pc: 'PC-B-05', status: 'Present' },
]

type Filter = 'All' | 'Present' | 'Absent'

export default function Attendance() {
  const [search, setSearch]   = useState('')
  const [date, setDate]       = useState('2026-09-25')
  const [filter, setFilter]   = useState<Filter>('All')
  const [selected, setSelected] = useState<string | null>(null)
  const [exported, setExported] = useState(false)

  const filtered = records.filter(r => {
    const q = search.toLowerCase()
    return (
      (filter === 'All' || r.status === filter) &&
      (r.name.toLowerCase().includes(q) || r.id.toLowerCase().includes(q) || r.pc.toLowerCase().includes(q))
    )
  })

  const present = records.filter(r => r.status === 'Present').length
  const absent  = records.filter(r => r.status === 'Absent').length

  function handleExport() {
    setExported(true)
    setTimeout(() => setExported(false), 2500)
  }

  return (
    <div className="space-y-5 max-w-5xl">
      {/* Summary pills */}
      <div className="flex items-center gap-3">
        {[
          { label: 'Total', val: records.length, color: '#3B5998', bg: '#EEF2FB' },
          { label: 'Present', val: present, color: '#10b981', bg: '#ecfdf5' },
          { label: 'Absent',  val: absent,  color: '#ef4444', bg: '#fef2f2' },
        ].map(s => (
          <div key={s.label} className="flex items-center gap-2 px-4 py-2 rounded-xl text-[13px] font-semibold" style={{ background: s.bg, color: s.color }}>
            <span style={{ fontFamily: 'JetBrains Mono', fontSize: 16 }}>{s.val}</span>
            <span className="font-medium text-slate-500">{s.label}</span>
          </div>
        ))}
        <div className="ml-auto text-[12px] text-slate-400 font-medium">
          Session: AM Block · 8:00 – 11:00
        </div>
      </div>

      {/* Toolbar */}
      <div className="bg-white rounded-2xl p-4 flex items-center gap-3 flex-wrap" style={{ border: '1px solid #E8EDF7' }}>
        {/* Search */}
        <div className="flex items-center gap-2 flex-1 min-w-48 max-w-sm px-3.5 py-2.5 rounded-xl" style={{ background: '#F4F6FB', border: '1px solid #E8EDF7' }}>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="11" cy="11" r="8"/><path d="M21 21l-4.35-4.35"/>
          </svg>
          <input
            type="text" value={search} onChange={e => setSearch(e.target.value)}
            placeholder="Search name or ID..."
            className="bg-transparent text-[13px] text-slate-600 placeholder-slate-400 outline-none flex-1"
          />
          {search && (
            <button onClick={() => setSearch('')} className="text-slate-300 hover:text-slate-500 cursor-pointer">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
                <path d="M18 6L6 18M6 6l12 12"/>
              </svg>
            </button>
          )}
        </div>

        {/* Date */}
        <div className="flex items-center gap-2 px-3.5 py-2.5 rounded-xl cursor-pointer" style={{ background: '#F4F6FB', border: '1px solid #E8EDF7' }}>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/>
          </svg>
          <input type="date" value={date} onChange={e => setDate(e.target.value)}
            className="bg-transparent text-[13px] text-slate-600 outline-none cursor-pointer" />
        </div>

        {/* Filter tabs */}
        <div className="flex rounded-xl overflow-hidden" style={{ border: '1px solid #E8EDF7', background: '#F4F6FB' }}>
          {(['All', 'Present', 'Absent'] as Filter[]).map(f => (
            <button
              key={f}
              onClick={() => setFilter(f)}
              className="px-4 py-2.5 text-[12.5px] font-semibold cursor-pointer transition-all"
              style={filter === f
                ? { background: 'white', color: '#3B5998', boxShadow: '0 1px 3px rgba(0,0,0,0.08)' }
                : { background: 'transparent', color: '#94a3b8' }
              }
            >
              {f}
            </button>
          ))}
        </div>

        {/* Export */}
        <button
          onClick={handleExport}
          className="ml-auto flex items-center gap-2 px-4 py-2.5 rounded-xl text-[13px] font-semibold text-white cursor-pointer transition-all hover:opacity-90 active:scale-95"
          style={{ background: exported ? '#10b981' : '#3B5998' }}
        >
          {exported ? (
            <>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                <path d="M5 13l4 4L19 7"/>
              </svg>
              Exported!
            </>
          ) : (
            <>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/>
              </svg>
              Export Data
            </>
          )}
        </button>
      </div>

      {/* Table */}
      <div className="bg-white rounded-2xl overflow-hidden" style={{ border: '1px solid #E8EDF7' }}>
        <table className="w-full">
          <thead>
            <tr style={{ background: '#F8FAFD', borderBottom: '1px solid #E8EDF7' }}>
              {['Student ID', 'Full Name', 'Time In', 'Time Out', 'PC Used', 'Status'].map(col => (
                <th key={col} className="text-left px-5 py-3.5 text-[11px] font-bold text-slate-400 uppercase tracking-widest">{col}</th>
              ))}
              <th className="px-5 py-3.5"></th>
            </tr>
          </thead>
          <tbody>
            {filtered.map(r => (
              <tr
                key={r.id}
                onClick={() => setSelected(selected === r.id ? null : r.id)}
                className="cursor-pointer transition-colors"
                style={{
                  borderBottom: '1px solid #F4F6FB',
                  background: selected === r.id ? '#F0F4FF' : 'white',
                }}
                onMouseEnter={e => { if (selected !== r.id) (e.currentTarget as HTMLElement).style.background = '#F8FAFD' }}
                onMouseLeave={e => { if (selected !== r.id) (e.currentTarget as HTMLElement).style.background = 'white' }}
              >
                <td className="px-5 py-4 text-[12.5px] font-bold" style={{ color: '#3B5998', fontFamily: 'JetBrains Mono' }}>{r.id}</td>
                <td className="px-5 py-4 text-[13.5px] font-semibold text-slate-700">{r.name}</td>
                <td className="px-5 py-4 text-[12.5px] text-slate-500" style={{ fontFamily: 'JetBrains Mono' }}>{r.timeIn}</td>
                <td className="px-5 py-4 text-[12.5px] text-slate-500" style={{ fontFamily: 'JetBrains Mono' }}>{r.timeOut}</td>
                <td className="px-5 py-4 text-[12.5px] font-medium text-slate-600" style={{ fontFamily: 'JetBrains Mono' }}>{r.pc}</td>
                <td className="px-5 py-4">
                  <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[12px] font-semibold"
                    style={r.status === 'Present'
                      ? { background: '#ecfdf5', color: '#059669', border: '1px solid #a7f3d0' }
                      : { background: '#fef2f2', color: '#dc2626', border: '1px solid #fecaca' }
                    }
                  >
                    <span className="w-1.5 h-1.5 rounded-full" style={{ background: r.status === 'Present' ? '#10b981' : '#ef4444' }}></span>
                    {r.status}
                  </span>
                </td>
                <td className="px-5 py-4">
                  <button className="w-7 h-7 rounded-lg flex items-center justify-center transition-colors cursor-pointer hover:bg-slate-100">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <circle cx="12" cy="5" r="1" fill="#94a3b8"/><circle cx="12" cy="12" r="1" fill="#94a3b8"/><circle cx="12" cy="19" r="1" fill="#94a3b8"/>
                    </svg>
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {filtered.length === 0 && (
          <div className="py-16 text-center">
            <div className="text-slate-300 mb-2">
              <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" className="mx-auto">
                <circle cx="11" cy="11" r="8"/><path d="M21 21l-4.35-4.35"/>
              </svg>
            </div>
            <p className="text-[13px] text-slate-400 font-medium">No records found</p>
            <button onClick={() => { setSearch(''); setFilter('All') }} className="mt-2 text-[12px] cursor-pointer" style={{ color: '#3B5998' }}>Clear filters</button>
          </div>
        )}

        {/* Footer */}
        <div className="px-5 py-3.5 flex items-center justify-between" style={{ borderTop: '1px solid #F4F6FB' }}>
          <span className="text-[12px] text-slate-400">Showing <span className="font-semibold text-slate-600">{filtered.length}</span> of {records.length} records</span>
          <div className="flex gap-1">
            {[1, 2, 3].map(p => (
              <button key={p}
                className="w-8 h-8 rounded-lg text-[12.5px] font-semibold cursor-pointer transition-colors"
                style={p === 1 ? { background: '#3B5998', color: 'white' } : { color: '#64748b' }}
              >{p}</button>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
