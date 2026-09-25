import { useState } from 'react'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import type { Page } from '../App'

const weeklyData = [
  { day: 'Mon', present: 38, absent: 7 },
  { day: 'Tue', present: 42, absent: 3 },
  { day: 'Wed', present: 35, absent: 10 },
  { day: 'Thu', present: 44, absent: 1 },
  { day: 'Fri', present: 29, absent: 16 },
  { day: 'Sat', present: 18, absent: 5 },
]

const activityLog = [
  { id: 'AST-0041', action: 'Marked for repair', component: 'Monitor — Dell P2422H', lab: 'Lab A', time: '09:14 AM', type: 'warn' },
  { id: 'AST-0028', action: 'Returned to active', component: 'Keyboard — Logitech K120', lab: 'Lab B', time: '08:55 AM', type: 'ok' },
  { id: 'AST-0015', action: 'Asset added', component: 'CPU — HP ProDesk 400', lab: 'Lab C', time: '08:30 AM', type: 'info' },
  { id: 'AST-0036', action: 'Status: Broken', component: 'Mouse — Logitech M100', lab: 'Lab A', time: '08:10 AM', type: 'danger' },
  { id: 'AST-0009', action: 'Maintenance done', component: 'Monitor — Acer V246HL', lab: 'Lab B', time: '07:50 AM', type: 'ok' },
]

const typeStyle: Record<string, { dot: string; bg: string }> = {
  ok:     { dot: '#10b981', bg: '#ecfdf5' },
  warn:   { dot: '#f59e0b', bg: '#fffbeb' },
  danger: { dot: '#ef4444', bg: '#fef2f2' },
  info:   { dot: '#3B5998', bg: '#EEF2FB' },
}

export default function Dashboard({ onNavigate }: { onNavigate: (p: Page) => void }) {
  const [activeBar, setActiveBar] = useState<string | null>(null)

  const metrics = [
    {
      label: 'Students Present', value: '44', sub: 'of 45 today',
      icon: <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#10b981" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 00-3-3.87M16 3.13a4 4 0 010 7.75"/></svg>,
      color: '#10b981', bg: '#ecfdf5', border: '#a7f3d0', page: 'attendance' as Page,
    },
    {
      label: 'Students Absent', value: '1', sub: 'below average',
      icon: <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#ef4444" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>,
      color: '#ef4444', bg: '#fef2f2', border: '#fecaca', page: 'attendance' as Page,
    },
    {
      label: 'Total Computers', value: '120', sub: 'across 3 labs',
      icon: <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#3B5998" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="2" y="3" width="20" height="14" rx="2"/><path d="M8 21h8m-4-4v4"/></svg>,
      color: '#3B5998', bg: '#EEF2FB', border: '#c7d3ee', page: 'assets' as Page,
    },
    {
      label: 'Under Maintenance', value: '7', sub: '5.8% of inventory',
      icon: <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M14.7 6.3a1 1 0 000 1.4l1.6 1.6a1 1 0 001.4 0l3.77-3.77a6 6 0 01-7.94 7.94l-6.91 6.91a2.12 2.12 0 01-3-3l6.91-6.91a6 6 0 017.94-7.94l-3.76 3.76z"/></svg>,
      color: '#f59e0b', bg: '#fffbeb', border: '#fde68a', page: 'assets' as Page,
    },
  ]

  return (
    <div className="space-y-6 max-w-6xl">
      {/* Metric Cards */}
      <div className="grid grid-cols-4 gap-4">
        {metrics.map(m => (
          <button
            key={m.label}
            onClick={() => onNavigate(m.page)}
            className="text-left rounded-2xl p-5 transition-all duration-200 hover:shadow-md hover:-translate-y-0.5 cursor-pointer group"
            style={{ background: m.bg, border: `1px solid ${m.border}` }}
          >
            <div className="flex items-start justify-between mb-3">
              <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ background: 'white' }}>
                {m.icon}
              </div>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke={m.color} strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" className="mt-1 opacity-0 group-hover:opacity-100 transition-opacity">
                <path d="M7 17L17 7M17 7H7M17 7v10" />
              </svg>
            </div>
            <div className="text-[30px] font-bold leading-none" style={{ color: m.color, fontFamily: 'JetBrains Mono, monospace' }}>{m.value}</div>
            <div className="text-[13px] font-semibold text-slate-700 mt-2">{m.label}</div>
            <div className="text-[11.5px] text-slate-400 mt-0.5">{m.sub}</div>
          </button>
        ))}
      </div>

      {/* Chart + Log */}
      <div className="grid grid-cols-5 gap-5">
        {/* Bar Chart */}
        <div className="col-span-3 bg-white rounded-2xl p-6" style={{ border: '1px solid #E8EDF7' }}>
          <div className="flex items-center justify-between mb-5">
            <div>
              <h2 className="text-[15px] font-bold text-slate-800">Weekly Attendance</h2>
              <p className="text-[12px] text-slate-400 mt-0.5">Sep 20 – Sep 25, 2026</p>
            </div>
            <div className="flex items-center gap-4 text-[12px] text-slate-500">
              <span className="flex items-center gap-1.5">
                <span className="w-3 h-3 rounded-sm inline-block" style={{ background: '#3B5998' }}></span>Present
              </span>
              <span className="flex items-center gap-1.5">
                <span className="w-3 h-3 rounded-sm inline-block" style={{ background: '#E8EDF7', border: '1px solid #c7d3ee' }}></span>Absent
              </span>
            </div>
          </div>
          <ResponsiveContainer width="100%" height={210}>
            <BarChart data={weeklyData} barSize={22} barGap={4}
              onMouseMove={e => e.activeLabel ? setActiveBar(e.activeLabel as string) : null}
              onMouseLeave={() => setActiveBar(null)}
            >
              <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" vertical={false} />
              <XAxis dataKey="day" tick={{ fontSize: 12, fill: '#94a3b8', fontFamily: 'JetBrains Mono' }} axisLine={false} tickLine={false} />
              <YAxis tick={{ fontSize: 12, fill: '#94a3b8', fontFamily: 'JetBrains Mono' }} axisLine={false} tickLine={false} width={28} />
              <Tooltip
                contentStyle={{ fontSize: 13, border: '1px solid #E8EDF7', borderRadius: 12, boxShadow: '0 8px 24px rgba(0,0,0,0.08)', fontFamily: 'Inter' }}
                cursor={{ fill: '#F4F6FB', radius: 8 }}
              />
              <Bar dataKey="present" fill="#3B5998" radius={[5, 5, 0, 0]} />
              <Bar dataKey="absent" fill="#E8EDF7" stroke="#c7d3ee" strokeWidth={1} radius={[5, 5, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>

        {/* Activity Log */}
        <div className="col-span-2 bg-white rounded-2xl p-6" style={{ border: '1px solid #E8EDF7' }}>
          <div className="flex items-center justify-between mb-5">
            <div>
              <h2 className="text-[15px] font-bold text-slate-800">Asset Activity</h2>
              <p className="text-[12px] text-slate-400 mt-0.5">Latest 5 events</p>
            </div>
            <button
              onClick={() => onNavigate('assets')}
              className="text-[11.5px] font-semibold px-3 py-1.5 rounded-lg cursor-pointer transition-colors hover:opacity-80"
              style={{ color: '#3B5998', background: '#EEF2FB' }}
            >
              View all
            </button>
          </div>
          <div className="space-y-3">
            {activityLog.map((log, i) => {
              const s = typeStyle[log.type]
              return (
                <div key={i} className="flex items-start gap-3 p-3 rounded-xl cursor-pointer hover:bg-slate-50 transition-colors group">
                  <div className="w-8 h-8 rounded-lg flex items-center justify-center shrink-0 mt-0.5" style={{ background: s.bg }}>
                    <span className="w-2 h-2 rounded-full" style={{ background: s.dot }}></span>
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-baseline justify-between">
                      <span className="text-[11.5px] font-bold" style={{ color: '#3B5998', fontFamily: 'JetBrains Mono' }}>{log.id}</span>
                      <span className="text-[10.5px] text-slate-400 shrink-0 ml-2">{log.time}</span>
                    </div>
                    <p className="text-[12.5px] font-semibold text-slate-700 leading-tight mt-0.5">{log.action}</p>
                    <p className="text-[11.5px] text-slate-400 truncate">{log.component}</p>
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      </div>
    </div>
  )
}
