import { useState } from 'react'

type Status = 'Active' | 'Maintenance' | 'Broken'
type Asset = {
  id: string; component: string; brand: string; lab: string
  status: Status; category: string; addedDate: string; repairHistory: string[]
}

const ASSETS: Asset[] = [
  { id: 'AST-0001', component: 'Monitor',        brand: 'Dell P2422H',          lab: 'Lab A', status: 'Active',      category: 'Display',    addedDate: 'Jan 10, 2024', repairHistory: [] },
  { id: 'AST-0002', component: 'CPU',             brand: 'HP ProDesk 400 G7',    lab: 'Lab A', status: 'Active',      category: 'Computer',   addedDate: 'Jan 10, 2024', repairHistory: [] },
  { id: 'AST-0003', component: 'Keyboard',        brand: 'Logitech K120',        lab: 'Lab B', status: 'Maintenance', category: 'Peripheral', addedDate: 'Feb 3, 2024',  repairHistory: ['Mar 12 — Key replacement', 'Apr 1 — Deep cleaning'] },
  { id: 'AST-0004', component: 'Mouse',           brand: 'Logitech M100',        lab: 'Lab A', status: 'Broken',      category: 'Peripheral', addedDate: 'Jan 10, 2024', repairHistory: ['Aug 20 — Scroll wheel replaced', 'Sep 15 — Reported broken'] },
  { id: 'AST-0005', component: 'Monitor',         brand: 'Acer V246HL',          lab: 'Lab B', status: 'Active',      category: 'Display',    addedDate: 'Mar 5, 2024',  repairHistory: [] },
  { id: 'AST-0006', component: 'UPS',             brand: 'APC BX700U',           lab: 'Lab C', status: 'Maintenance', category: 'Power',      addedDate: 'Jan 15, 2024', repairHistory: ['Jun 3 — Battery replacement'] },
  { id: 'AST-0007', component: 'CPU',             brand: 'Lenovo ThinkCentre M70q', lab: 'Lab C', status: 'Active',  category: 'Computer',   addedDate: 'Apr 20, 2024', repairHistory: [] },
  { id: 'AST-0008', component: 'Headset',         brand: 'JBL Quantum 100',      lab: 'Lab B', status: 'Active',      category: 'Peripheral', addedDate: 'May 1, 2024',  repairHistory: [] },
  { id: 'AST-0009', component: 'Network Switch',  brand: 'TP-Link TL-SG108',     lab: 'Lab A', status: 'Active',      category: 'Network',    addedDate: 'Jan 10, 2024', repairHistory: [] },
  { id: 'AST-0010', component: 'Webcam',          brand: 'Logitech C270',        lab: 'Lab C', status: 'Broken',      category: 'Peripheral', addedDate: 'Jun 10, 2024', repairHistory: ['Sep 20 — Lens cracked, parts ordered'] },
]

const CATEGORIES = ['All', 'Computer', 'Display', 'Peripheral', 'Power', 'Network']

const STATUS_STYLE: Record<Status, { bg: string; color: string; border: string; dot: string }> = {
  Active:      { bg: '#ecfdf5', color: '#059669', border: '#a7f3d0', dot: '#10b981' },
  Maintenance: { bg: '#fffbeb', color: '#d97706', border: '#fde68a', dot: '#f59e0b' },
  Broken:      { bg: '#fef2f2', color: '#dc2626', border: '#fecaca', dot: '#ef4444' },
}

export default function Assets() {
  const [category, setCategory] = useState('All')
  const [search, setSearch]     = useState('')
  const [selected, setSelected] = useState<Asset | null>(null)
  const [showModal, setShowModal] = useState(false)
  const [form, setForm] = useState({ name: '', brand: '', lab: '', category: 'Computer' })

  const filtered = ASSETS.filter(a =>
    (category === 'All' || a.category === category) &&
    (a.component.toLowerCase().includes(search.toLowerCase()) ||
     a.brand.toLowerCase().includes(search.toLowerCase()) ||
     a.id.toLowerCase().includes(search.toLowerCase()))
  )

  const stats = {
    active:      ASSETS.filter(a => a.status === 'Active').length,
    maintenance: ASSETS.filter(a => a.status === 'Maintenance').length,
    broken:      ASSETS.filter(a => a.status === 'Broken').length,
  }

  return (
    <div className="space-y-5 max-w-6xl">
      {/* Summary */}
      <div className="flex items-center gap-3">
        {[
          { label: 'Active',      val: stats.active,      ...STATUS_STYLE.Active },
          { label: 'Maintenance', val: stats.maintenance, ...STATUS_STYLE.Maintenance },
          { label: 'Broken',      val: stats.broken,      ...STATUS_STYLE.Broken },
        ].map(s => (
          <div key={s.label} className="flex items-center gap-2 px-4 py-2 rounded-xl text-[13px] font-semibold cursor-pointer hover:opacity-90 transition-opacity"
            style={{ background: s.bg, color: s.color, border: `1px solid ${s.border}` }}
            onClick={() => {}}
          >
            <span className="w-2 h-2 rounded-full" style={{ background: s.dot }}></span>
            <span style={{ fontFamily: 'JetBrains Mono', fontSize: 15, fontWeight: 700 }}>{s.val}</span>
            <span className="font-medium text-slate-500">{s.label}</span>
          </div>
        ))}
      </div>

      {/* Toolbar */}
      <div className="bg-white rounded-2xl p-4 flex items-center gap-3 flex-wrap" style={{ border: '1px solid #E8EDF7' }}>
        <div className="flex items-center gap-2 flex-1 min-w-48 max-w-sm px-3.5 py-2.5 rounded-xl" style={{ background: '#F4F6FB', border: '1px solid #E8EDF7' }}>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="11" cy="11" r="8"/><path d="M21 21l-4.35-4.35"/>
          </svg>
          <input type="text" value={search} onChange={e => setSearch(e.target.value)}
            placeholder="Search assets..." className="bg-transparent text-[13px] text-slate-600 placeholder-slate-400 outline-none flex-1" />
        </div>

        <div className="flex items-center gap-2 px-3.5 py-2.5 rounded-xl" style={{ background: '#F4F6FB', border: '1px solid #E8EDF7' }}>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/>
          </svg>
          <select value={category} onChange={e => setCategory(e.target.value)}
            className="bg-transparent text-[13px] text-slate-600 outline-none cursor-pointer">
            {CATEGORIES.map(c => <option key={c}>{c}</option>)}
          </select>
        </div>

        <button onClick={() => setShowModal(true)}
          className="ml-auto flex items-center gap-2 px-4 py-2.5 rounded-xl text-[13px] font-semibold text-white cursor-pointer hover:opacity-90 active:scale-95 transition-all"
          style={{ background: '#3B5998' }}>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="2.5" strokeLinecap="round">
            <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          Add New Asset
        </button>
      </div>

      {/* Table + Detail Panel */}
      <div className="flex gap-4 items-start">
        <div className="bg-white rounded-2xl overflow-hidden flex-1" style={{ border: '1px solid #E8EDF7' }}>
          <table className="w-full">
            <thead>
              <tr style={{ background: '#F8FAFD', borderBottom: '1px solid #E8EDF7' }}>
                {['Asset ID', 'Component', 'Brand', 'Location', 'Status', ''].map(col => (
                  <th key={col} className="text-left px-5 py-3.5 text-[11px] font-bold text-slate-400 uppercase tracking-widest">{col}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {filtered.map(a => {
                const ss = STATUS_STYLE[a.status]
                const isSelected = selected?.id === a.id
                return (
                  <tr key={a.id}
                    onClick={() => setSelected(isSelected ? null : a)}
                    className="cursor-pointer transition-colors"
                    style={{ borderBottom: '1px solid #F4F6FB', background: isSelected ? '#F0F4FF' : 'white' }}
                    onMouseEnter={e => { if (!isSelected) (e.currentTarget as HTMLElement).style.background = '#F8FAFD' }}
                    onMouseLeave={e => { if (!isSelected) (e.currentTarget as HTMLElement).style.background = 'white' }}
                  >
                    <td className="px-5 py-4 text-[12.5px] font-bold" style={{ color: '#3B5998', fontFamily: 'JetBrains Mono' }}>{a.id}</td>
                    <td className="px-5 py-4 text-[13.5px] font-semibold text-slate-700">{a.component}</td>
                    <td className="px-5 py-4 text-[13px] text-slate-500">{a.brand}</td>
                    <td className="px-5 py-4">
                      <span className="text-[12px] font-semibold px-2.5 py-1 rounded-lg" style={{ background: '#EEF2FB', color: '#3B5998' }}>{a.lab}</span>
                    </td>
                    <td className="px-5 py-4">
                      <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[12px] font-semibold"
                        style={{ background: ss.bg, color: ss.color, border: `1px solid ${ss.border}` }}>
                        <span className="w-1.5 h-1.5 rounded-full" style={{ background: ss.dot }}></span>
                        {a.status}
                      </span>
                    </td>
                    <td className="px-5 py-4">
                      <div className="flex items-center gap-1 justify-end">
                        <button className="w-7 h-7 rounded-lg flex items-center justify-center hover:bg-slate-100 transition-colors cursor-pointer"
                          onClick={e => { e.stopPropagation(); setSelected(isSelected ? null : a) }}>
                          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/>
                          </svg>
                        </button>
                      </div>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
          {filtered.length === 0 && (
            <div className="py-14 text-center text-[13px] text-slate-400">No assets found</div>
          )}
        </div>

        {/* Detail Panel */}
        {selected && (
          <div className="w-72 shrink-0 bg-white rounded-2xl overflow-hidden" style={{ border: '1px solid #E8EDF7' }}>
            {/* Header */}
            <div className="px-5 py-4 flex items-center justify-between" style={{ borderBottom: '1px solid #E8EDF7' }}>
              <p className="text-[13.5px] font-bold text-slate-700">Asset Details</p>
              <button onClick={() => setSelected(null)} className="w-7 h-7 rounded-lg flex items-center justify-center text-slate-400 hover:bg-slate-100 cursor-pointer transition-colors">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
                  <path d="M18 6L6 18M6 6l12 12"/>
                </svg>
              </button>
            </div>

            {/* Card */}
            <div className="p-5 space-y-4">
              <div className="rounded-xl p-4" style={{ background: '#EEF2FB' }}>
                <p className="text-[11px] font-bold" style={{ color: '#6B85C4', fontFamily: 'JetBrains Mono' }}>{selected.id}</p>
                <p className="text-[17px] font-bold text-slate-800 mt-1">{selected.component}</p>
                <p className="text-[13px] text-slate-500 mt-0.5">{selected.brand}</p>
                <div className="mt-3">
                  <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11.5px] font-semibold"
                    style={{ background: STATUS_STYLE[selected.status].bg, color: STATUS_STYLE[selected.status].color, border: `1px solid ${STATUS_STYLE[selected.status].border}` }}>
                    <span className="w-1.5 h-1.5 rounded-full" style={{ background: STATUS_STYLE[selected.status].dot }}></span>
                    {selected.status}
                  </span>
                </div>
              </div>

              <div className="space-y-2.5">
                {[
                  { label: 'Location', val: selected.lab },
                  { label: 'Category', val: selected.category },
                  { label: 'Added',    val: selected.addedDate },
                ].map(f => (
                  <div key={f.label} className="flex justify-between items-center">
                    <span className="text-[12px] text-slate-400">{f.label}</span>
                    <span className="text-[12.5px] font-semibold text-slate-700">{f.val}</span>
                  </div>
                ))}
              </div>

              <div style={{ borderTop: '1px solid #E8EDF7', paddingTop: 14 }}>
                <p className="text-[11px] font-bold text-slate-400 uppercase tracking-widest mb-3">Repair History</p>
                {selected.repairHistory.length === 0 ? (
                  <p className="text-[12px] text-slate-400 italic">No history</p>
                ) : (
                  <div className="space-y-2">
                    {selected.repairHistory.map((r, i) => (
                      <div key={i} className="flex items-start gap-2 p-2.5 rounded-lg" style={{ background: '#fffbeb' }}>
                        <span className="w-1.5 h-1.5 rounded-full mt-1.5 shrink-0" style={{ background: '#f59e0b' }}></span>
                        <span className="text-[12px] text-slate-600">{r}</span>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              <div className="flex gap-2 pt-1">
                <button className="flex-1 py-2.5 rounded-xl text-[12.5px] font-semibold cursor-pointer hover:opacity-90 transition-opacity"
                  style={{ background: '#fffbeb', color: '#d97706', border: '1px solid #fde68a' }}>
                  Mark Repair
                </button>
                <button className="flex-1 py-2.5 rounded-xl text-[12.5px] font-semibold cursor-pointer hover:opacity-90 transition-opacity"
                  style={{ background: '#fef2f2', color: '#dc2626', border: '1px solid #fecaca' }}>
                  Report Broken
                </button>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Add Asset Modal */}
      {showModal && (
        <div className="fixed inset-0 flex items-center justify-center z-50" style={{ background: 'rgba(15,23,42,0.35)', backdropFilter: 'blur(4px)' }}
          onClick={() => setShowModal(false)}>
          <div className="bg-white rounded-2xl w-full max-w-md shadow-2xl overflow-hidden" onClick={e => e.stopPropagation()}
            style={{ border: '1px solid #E8EDF7' }}>
            <div className="px-6 py-5 flex items-center justify-between" style={{ borderBottom: '1px solid #E8EDF7' }}>
              <div>
                <h2 className="text-[15px] font-bold text-slate-800">Add New Asset</h2>
                <p className="text-[12px] text-slate-400 mt-0.5">Fill in the asset details below</p>
              </div>
              <button onClick={() => setShowModal(false)} className="w-8 h-8 rounded-xl flex items-center justify-center text-slate-400 hover:bg-slate-100 cursor-pointer transition-colors">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
                  <path d="M18 6L6 18M6 6l12 12"/>
                </svg>
              </button>
            </div>
            <div className="p-6 space-y-4">
              {[
                { key: 'name',  label: 'Component Name', placeholder: 'e.g. Monitor, CPU, Keyboard' },
                { key: 'brand', label: 'Brand / Model',   placeholder: 'e.g. Dell P2422H' },
                { key: 'lab',   label: 'Lab Location',    placeholder: 'e.g. Lab A' },
              ].map(f => (
                <div key={f.key}>
                  <label className="text-[11.5px] font-semibold text-slate-500 uppercase tracking-wider">{f.label}</label>
                  <input type="text" placeholder={f.placeholder}
                    value={(form as any)[f.key]}
                    onChange={e => setForm(prev => ({ ...prev, [f.key]: e.target.value }))}
                    className="mt-1.5 w-full rounded-xl px-4 py-2.5 text-[13.5px] text-slate-700 placeholder-slate-400 outline-none transition-all"
                    style={{ border: '1.5px solid #E8EDF7', background: '#F8FAFD' }}
                    onFocus={e => (e.target.style.borderColor = '#3B5998')}
                    onBlur={e => (e.target.style.borderColor = '#E8EDF7')}
                  />
                </div>
              ))}
              <div>
                <label className="text-[11.5px] font-semibold text-slate-500 uppercase tracking-wider">Category</label>
                <select value={form.category} onChange={e => setForm(prev => ({ ...prev, category: e.target.value }))}
                  className="mt-1.5 w-full rounded-xl px-4 py-2.5 text-[13.5px] text-slate-700 outline-none cursor-pointer"
                  style={{ border: '1.5px solid #E8EDF7', background: '#F8FAFD' }}>
                  {CATEGORIES.filter(c => c !== 'All').map(c => <option key={c}>{c}</option>)}
                </select>
              </div>
            </div>
            <div className="px-6 pb-5 flex gap-3">
              <button onClick={() => setShowModal(false)}
                className="flex-1 py-3 rounded-xl text-[13px] font-semibold cursor-pointer hover:bg-slate-100 transition-colors text-slate-600"
                style={{ border: '1px solid #E8EDF7' }}>
                Cancel
              </button>
              <button onClick={() => setShowModal(false)}
                className="flex-1 py-3 rounded-xl text-[13px] font-semibold text-white cursor-pointer hover:opacity-90 transition-opacity active:scale-95"
                style={{ background: '#3B5998' }}>
                Add Asset
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
