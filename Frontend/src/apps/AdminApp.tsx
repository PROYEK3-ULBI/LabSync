import { useMemo, useState } from 'react'
import { Icon } from '../components/icons'
import { Shell, type NavItem } from '../components/Shell'
import {
  allocations,
  assets,
  attendanceMonitoring,
  auditLog,
  maintenance,
  recentActivity,
  sessions,
  users,
  weeklyAttendance,
} from '../lib/data'
import {
  Alert,
  Avatar,
  Badge,
  Button,
  Card,
  Code,
  EmptyState,
  Field,
  IconButton,
  Input,
  Modal,
  PageHeader,
  Pagination,
  Select,
  StatCard,
  StatusBadge,
  SkeletonRows,
  Table,
  TD,
  THead,
  TR,
  Tabs,
  Timeline,
  useSimulatedLoading,
  useToast,
} from '../components/ui'
import type { Allocation, Asset, Maintenance, Session, User } from '../lib/data'

const NAV: NavItem[] = [
  { key: 'dashboard', label: 'Dasbor', icon: 'dashboard' },
  { key: 'users', label: 'Pengguna', icon: 'users' },
  { key: 'sessions', label: 'Sesi Praktikum', icon: 'calendar' },
  { key: 'attendance', label: 'Presensi', icon: 'checkCircle' },
  { key: 'assets', label: 'Aset', icon: 'box' },
  { key: 'allocations', label: 'Alokasi', icon: 'swap' },
  { key: 'maintenance', label: 'Pemeliharaan', icon: 'wrench' },
  { key: 'reports', label: 'Laporan', icon: 'report' },
  { key: 'audit', label: 'Log Audit', icon: 'history' },
  { key: 'profile', label: 'Profil', icon: 'user' },
]

export function AdminApp({ onLogout }: { onLogout: () => void }) {
  const [route, setRoute] = useState('dashboard')
  const [asset, setAsset] = useState<Asset | null>(null)
  const [user, setUser] = useState<User | null>(null)
  const [session, setSession] = useState<Session | null>(null)
  const openAsset = (a: Asset) => {
    setAsset(a)
    setRoute('assetDetail')
  }
  const openUser = (u: User) => {
    setUser(u)
    setRoute('userDetail')
  }
  const openSession = (s: Session) => {
    setSession(s)
    setRoute('sessionDetail')
  }
  const [alloc, setAlloc] = useState<Allocation | null>(null)
  const [mt, setMt] = useState<Maintenance | null>(null)
  const openAlloc = (a: Allocation) => {
    setAlloc(a)
    setRoute('allocationDetail')
  }
  const openMt = (m: Maintenance) => {
    setMt(m)
    setRoute('maintenanceDetail')
  }
  return (
    <Shell
      brand="LabPresensi"
      roleLabel="Admin / Laboran"
      roleTone="slate"
      nav={NAV}
      active={route}
      onNavigate={setRoute}
      user={{ name: 'Dewi Anggraini', sub: 'Laboran' }}
      onLogout={onLogout}
    >
      {route === 'dashboard' && <Dashboard go={setRoute} />}
      {route === 'users' && <Users onView={openUser} />}
      {route === 'userDetail' && <UserDetail user={user} go={setRoute} />}
      {route === 'sessions' && <SessionMgmt onView={openSession} />}
      {route === 'sessionDetail' && <SessionDetail session={session} go={setRoute} />}
      {route === 'attendance' && <AttendanceMonitoring />}
      {route === 'assets' && <AssetMgmt onView={openAsset} />}
      {route === 'assetDetail' && <AssetDetail asset={asset} go={setRoute} />}
      {route === 'allocations' && <Allocations onView={openAlloc} />}
      {route === 'allocationDetail' && <AllocationDetail alloc={alloc} go={setRoute} />}
      {route === 'maintenance' && <MaintenanceMgmt onView={openMt} />}
      {route === 'maintenanceDetail' && <MaintenanceDetail mt={mt} go={setRoute} />}
      {route === 'reports' && <Reports />}
      {route === 'audit' && <AuditLog />}
      {route === 'profile' && <Profile />}
    </Shell>
  )
}

/* --------------------------------- Dashboard --------------------------------- */
function Dashboard({ go }: { go: (r: string) => void }) {
  const assetStatus = useMemo(() => {
    const c = { Available: 0, Allocated: 0, Maintenance: 0, Retired: 0 } as Record<string, number>
    assets.forEach((a) => (c[a.status] += 1))
    return c
  }, [])
  const barTones: Record<string, string> = {
    Available: 'var(--color-success)',
    Allocated: 'var(--color-info)',
    Maintenance: 'var(--color-warning)',
    Retired: 'var(--color-text-3)',
  }
  const max = Math.max(...weeklyAttendance.map((d) => d.total))

  return (
    <>
      <PageHeader title="Dasbor Operasional" subtitle="Ringkasan presensi, aset, dan pemeliharaan laboratorium." />

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-5">
        <StatCard label="Sesi Aktif" value="1" hint="Pemrograman Web" icon={<Icon.calendar />} tone="primary" />
        <StatCard label="Presensi Hari Ini" value="22/30" hint="73% kehadiran" icon={<Icon.checkCircle />} tone="success" />
        <StatCard label="Aset Tersedia" value={assetStatus.Available} icon={<Icon.box />} tone="info" />
        <StatCard label="Aset Dialokasikan" value={assetStatus.Allocated} icon={<Icon.laptop />} tone="warning" />
        <StatCard label="Dalam Pemeliharaan" value={assetStatus.Maintenance} icon={<Icon.wrench />} tone="error" />
      </div>

      <div className="grid gap-6 lg:grid-cols-3 mt-6">
        {/* attendance chart */}
        <Card className="lg:col-span-2">
          <div className="flex items-center justify-between mb-5">
            <div>
              <h3 className="font-semibold">Ringkasan Kehadiran</h3>
              <p className="text-sm text-[var(--color-text-2)]">Lima hari terakhir</p>
            </div>
            <Badge tone="success">Rata-rata 86%</Badge>
          </div>
          <div className="flex items-end gap-4 h-48">
            {weeklyAttendance.map((d) => (
              <div key={d.day} className="flex-1 flex flex-col items-center gap-2">
                <div className="w-full flex-1 flex items-end">
                  <div className="w-full rounded-t-md bg-[var(--color-surface-2)] relative overflow-hidden" style={{ height: '100%' }}>
                    <div
                      className="absolute bottom-0 left-0 right-0 rounded-t-md bg-[var(--color-primary)] transition-all"
                      style={{ height: `${(d.hadir / max) * 100}%` }}
                      title={`${d.hadir}/${d.total}`}
                    />
                  </div>
                </div>
                <span className="text-xs text-[var(--color-text-2)]">{d.day}</span>
                <span className="text-[11px] text-[var(--color-text-3)]">{d.hadir}/{d.total}</span>
              </div>
            ))}
          </div>
        </Card>

        {/* asset status */}
        <Card>
          <h3 className="font-semibold mb-4">Status Aset</h3>
          <div className="space-y-3">
            {Object.entries(assetStatus).map(([k, v]) => (
              <div key={k}>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-[var(--color-text-2)]">{k}</span>
                  <span className="font-medium">{v}</span>
                </div>
                <div className="h-2 rounded-full bg-[var(--color-surface-2)] overflow-hidden">
                  <div className="h-full rounded-full" style={{ width: `${(v / assets.length) * 100}%`, background: barTones[k] }} />
                </div>
              </div>
            ))}
          </div>
          <Button variant="secondary" className="w-full mt-5" onClick={() => go('assets')}>
            Kelola Aset
          </Button>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-3 mt-6">
        <Card className="lg:col-span-2" padded={false}>
          <div className="flex items-center justify-between p-5 border-b border-[var(--color-border)]">
            <h3 className="font-semibold">Sesi Aktif</h3>
            <Button variant="ghost" size="sm" onClick={() => go('sessions')}>Semua sesi</Button>
          </div>
          <Table>
            <THead cols={['Sesi', 'Ruang', 'Waktu', 'Kehadiran', 'Status']} />
            <tbody>
              {sessions.slice(0, 3).map((s) => (
                <TR key={s.code}>
                  <TD className="font-medium">{s.name}<div className="text-xs text-[var(--color-text-3)] font-mono">{s.code}</div></TD>
                  <TD>{s.room}</TD>
                  <TD>{s.start}–{s.end}</TD>
                  <TD>{s.status === 'Aktif' ? '22/30' : '—'}</TD>
                  <TD><StatusBadge status={s.status} /></TD>
                </TR>
              ))}
            </tbody>
          </Table>
        </Card>

        <Card padded={false}>
          <div className="p-5 border-b border-[var(--color-border)]">
            <h3 className="font-semibold">Aktivitas Terbaru</h3>
          </div>
          <ul className="divide-y divide-[var(--color-border)]">
            {recentActivity.map((a, i) => (
              <li key={i} className="flex items-start gap-3 px-5 py-3">
                <span className="mt-0.5 inline-flex items-center justify-center w-8 h-8 rounded-full bg-[var(--color-surface-2)] text-[var(--color-text-2)] shrink-0">
                  <Icon.history width={15} height={15} />
                </span>
                <div className="min-w-0">
                  <p className="text-sm font-medium truncate">{a.action}</p>
                  <p className="text-xs text-[var(--color-text-2)]">{a.user} · {a.entity}</p>
                  <p className="text-[11px] text-[var(--color-text-3)]">{a.time}</p>
                </div>
              </li>
            ))}
          </ul>
        </Card>
      </div>
    </>
  )
}

/* --------------------------------- Users --------------------------------- */
function Users({ onView }: { onView: (u: User) => void }) {
  const toast = useToast()
  const loading = useSimulatedLoading()
  const [q, setQ] = useState('')
  const [role, setRole] = useState('Semua')
  const [open, setOpen] = useState(false)
  const rows = users.filter(
    (u) =>
      (role === 'Semua' || u.role === role) &&
      (u.name.toLowerCase().includes(q.toLowerCase()) || u.nim.toLowerCase().includes(q.toLowerCase())),
  )
  return (
    <>
      <PageHeader
        title="Manajemen Pengguna"
        subtitle="Kelola mahasiswa dan staf laboratorium."
        actions={<Button icon={<Icon.plus width={16} />} onClick={() => setOpen(true)}>Tambah Pengguna</Button>}
      />
      <Card padded={false}>
        <div className="flex flex-wrap items-center gap-3 p-4 border-b border-[var(--color-border)]">
          <div className="flex-1 min-w-[220px]"><Input placeholder="Cari nama atau NIM…" value={q} onChange={(e) => setQ(e.target.value)} /></div>
          <Select value={role} onChange={(e) => setRole(e.target.value)} className="w-44">
            {['Semua', 'Mahasiswa', 'Laboran', 'Tata Usaha'].map((r) => <option key={r}>{r}</option>)}
          </Select>
        </div>
        {loading ? (
          <Table>
            <THead cols={['Nama', 'Identifier', 'Email', 'Peran', 'Status', 'Wajah', 'Dibuat', 'Aksi']} />
            <SkeletonRows cols={8} />
          </Table>
        ) : rows.length === 0 ? (
          <EmptyState title="Pengguna tidak ditemukan" desc="Coba ubah kata kunci atau filter." icon={<Icon.users />} />
        ) : (
          <>
            <Table>
              <THead cols={['Nama', 'Identifier', 'Email', 'Peran', 'Status', 'Wajah', 'Dibuat', 'Aksi']} />
              <tbody>
                {rows.map((u) => (
                  <TR key={u.nim} onClick={() => onView(u)}>
                    <TD>
                      <div className="flex items-center gap-2.5">
                        <Avatar name={u.name} size={32} />
                        <span className="font-medium">{u.name}</span>
                      </div>
                    </TD>
                    <TD><Code>{u.nim}</Code></TD>
                    <TD className="text-[var(--color-text-2)]">{u.email}</TD>
                    <TD>{u.role}</TD>
                    <TD><StatusBadge status={u.status} /></TD>
                    <TD>{u.face === '—' ? <span className="text-[var(--color-text-3)]">—</span> : <StatusBadge status={u.face} />}</TD>
                    <TD className="text-[var(--color-text-2)]">{u.created}</TD>
                    <TD>
                      <div className="flex gap-0.5" onClick={(e) => e.stopPropagation()}>
                        <IconButton label="Lihat" onClick={() => onView(u)}><Icon.eye width={17} /></IconButton>
                        <IconButton label="Ubah" onClick={() => onView(u)}><Icon.edit width={17} /></IconButton>
                      </div>
                    </TD>
                  </TR>
                ))}
              </tbody>
            </Table>
            <Pagination page={1} pages={1} total={rows.length} onPage={() => {}} />
          </>
        )}
      </Card>

      <Modal
        open={open}
        onClose={() => setOpen(false)}
        title="Tambah Pengguna"
        footer={
          <>
            <Button variant="secondary" onClick={() => setOpen(false)}>Batal</Button>
            <Button onClick={() => { setOpen(false); toast('Pengguna baru berhasil ditambahkan.') }}>Simpan</Button>
          </>
        }
      >
        <div className="grid sm:grid-cols-2 gap-4">
          <Field label="Nama Lengkap" required><Input placeholder="cth. Aulia Rahman" /></Field>
          <Field label="Identifier / NIM" required><Input placeholder="cth. TI240010" /></Field>
          <Field label="Email" required><Input type="email" placeholder="nama@univ.ac.id" /></Field>
          <Field label="Peran" required>
            <Select><option>Mahasiswa</option><option>Laboran</option><option>Tata Usaha</option></Select>
          </Field>
          <Field label="Status">
            <Select><option>Aktif</option><option>Nonaktif</option></Select>
          </Field>
        </div>
      </Modal>
    </>
  )
}

/* --------------------------------- Sessions --------------------------------- */
function SessionMgmt({ onView }: { onView: (s: Session) => void }) {
  const toast = useToast()
  const [open, setOpen] = useState(false)
  return (
    <>
      <PageHeader
        title="Manajemen Sesi Praktikum"
        subtitle="Buat dan kelola jadwal sesi praktikum."
        actions={<Button icon={<Icon.plus width={16} />} onClick={() => setOpen(true)}>Tambah Sesi</Button>}
      />
      <Card padded={false}>
        <Table>
          <THead cols={['Kode', 'Nama Sesi', 'Mata Kuliah', 'Tanggal', 'Waktu', 'Laboratorium', 'Status', 'Aksi']} />
          <tbody>
            {sessions.map((s) => (
              <TR key={s.code} onClick={() => onView(s)}>
                <TD><Code>{s.code}</Code></TD>
                <TD className="font-medium">{s.name}</TD>
                <TD>{s.course}</TD>
                <TD>{s.date}</TD>
                <TD>{s.start}–{s.end}</TD>
                <TD>{s.room}</TD>
                <TD><StatusBadge status={s.status} /></TD>
                <TD>
                  <div className="flex gap-0.5" onClick={(e) => e.stopPropagation()}>
                    <IconButton label="Lihat" onClick={() => onView(s)}><Icon.eye width={17} /></IconButton>
                    <IconButton label="Ubah" onClick={() => onView(s)}><Icon.edit width={17} /></IconButton>
                  </div>
                </TD>
              </TR>
            ))}
          </tbody>
        </Table>
      </Card>

      <Modal
        open={open}
        onClose={() => setOpen(false)}
        title="Tambah Sesi Praktikum"
        size="lg"
        footer={
          <>
            <Button variant="secondary" onClick={() => setOpen(false)}>Batal</Button>
            <Button onClick={() => { setOpen(false); toast('Sesi praktikum berhasil dibuat.') }}>Simpan Sesi</Button>
          </>
        }
      >
        <div className="grid sm:grid-cols-2 gap-4">
          <Field label="Kode Sesi" required><Input placeholder="cth. PWB-08" /></Field>
          <Field label="Nama Sesi" required><Input placeholder="Praktikum Pemrograman Web" /></Field>
          <Field label="Mata Kuliah" required><Input placeholder="Pemrograman Web" /></Field>
          <Field label="Laboratorium" required>
            <Select><option>Lab Komputer 1</option><option>Lab Komputer 2</option></Select>
          </Field>
          <Field label="Tanggal" required><Input type="date" /></Field>
          <div className="grid grid-cols-2 gap-4">
            <Field label="Mulai" required><Input type="time" /></Field>
            <Field label="Selesai" required><Input type="time" /></Field>
          </div>
          <Field label="Status">
            <Select><option>Draft</option><option>Aktif</option><option>Ditutup</option><option>Dibatalkan</option></Select>
          </Field>
        </div>
      </Modal>
    </>
  )
}

/* --------------------------------- Attendance monitoring --------------------------------- */
function AttendanceMonitoring() {
  const loading = useSimulatedLoading()
  const [session, setSession] = useState('Pemrograman Web')
  return (
    <>
      <PageHeader title="Pemantauan Presensi" subtitle="Pantau kehadiran dan status verifikasi mahasiswa per sesi." />
      <Card padded={false}>
        <div className="flex flex-wrap items-center gap-3 p-4 border-b border-[var(--color-border)]">
          <Select value={session} onChange={(e) => setSession(e.target.value)} className="w-56">
            <option>Pemrograman Web</option><option>Basis Data</option><option>Jaringan Komputer</option>
          </Select>
          <Input type="date" defaultValue="2026-09-18" className="w-44" />
          <div className="flex-1 min-w-[180px]"><Input placeholder="Cari mahasiswa…" /></div>
          <Select className="w-40"><option>Semua Status</option><option>Hadir</option><option>Tidak Hadir</option></Select>
        </div>
        <Table>
          <THead cols={['Mahasiswa', 'NIM', 'Sesi', 'Check-in', 'Check-out', 'Verifikasi', 'Presensi']} />
          {loading ? <SkeletonRows cols={7} /> : (
          <tbody>
            {attendanceMonitoring.map((r) => (
              <TR key={r.nim}>
                <TD>
                  <div className="flex items-center gap-2.5"><Avatar name={r.student} size={30} /><span className="font-medium">{r.student}</span></div>
                </TD>
                <TD><Code>{r.nim}</Code></TD>
                <TD>{r.session}</TD>
                <TD>{r.checkIn}</TD>
                <TD>{r.checkOut}</TD>
                <TD>{r.verify === '—' ? <span className="text-[var(--color-text-3)]">—</span> : <StatusBadge status={r.verify} />}</TD>
                <TD><StatusBadge status={r.status} /></TD>
              </TR>
            ))}
          </tbody>
          )}
        </Table>
        <Pagination page={1} pages={1} total={attendanceMonitoring.length} onPage={() => {}} />
      </Card>
    </>
  )
}

/* --------------------------------- Assets --------------------------------- */
function AssetMgmt({ onView }: { onView: (a: Asset) => void }) {
  const toast = useToast()
  const loading = useSimulatedLoading()
  const [q, setQ] = useState('')
  const [status, setStatus] = useState('Semua')
  const [open, setOpen] = useState(false)
  const rows = assets.filter(
    (a) => (status === 'Semua' || a.status === status) && (a.name.toLowerCase().includes(q.toLowerCase()) || a.code.toLowerCase().includes(q.toLowerCase())),
  )
  const counts = { Available: 0, Allocated: 0, Maintenance: 0, Retired: 0 } as Record<string, number>
  assets.forEach((a) => (counts[a.status] += 1))
  return (
    <>
      <PageHeader
        title="Manajemen Aset"
        subtitle="Inventaris perangkat laboratorium komputer."
        actions={<Button icon={<Icon.plus width={16} />} onClick={() => setOpen(true)}>Tambah Aset</Button>}
      />
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-5 mb-6">
        <StatCard label="Total Aset" value={assets.length} icon={<Icon.box />} tone="primary" />
        <StatCard label="Tersedia" value={counts.Available} icon={<Icon.checkCircle />} tone="success" />
        <StatCard label="Dialokasikan" value={counts.Allocated} icon={<Icon.laptop />} tone="info" />
        <StatCard label="Pemeliharaan" value={counts.Maintenance} icon={<Icon.wrench />} tone="warning" />
        <StatCard label="Ditarik" value={counts.Retired} icon={<Icon.box />} tone="neutral" />
      </div>
      <Card padded={false}>
        <div className="flex flex-wrap items-center gap-3 p-4 border-b border-[var(--color-border)]">
          <div className="flex-1 min-w-[220px]"><Input placeholder="Cari kode atau nama aset…" value={q} onChange={(e) => setQ(e.target.value)} /></div>
          <Select value={status} onChange={(e) => setStatus(e.target.value)} className="w-44">
            {['Semua', 'Available', 'Allocated', 'Maintenance', 'Retired'].map((s) => <option key={s}>{s}</option>)}
          </Select>
        </div>
        {loading ? (
          <Table>
            <THead cols={['Kode', 'Nama Aset', 'Kategori', 'Merek', 'Serial', 'Lokasi', 'Status', 'Aksi']} />
            <SkeletonRows cols={8} />
          </Table>
        ) : rows.length === 0 ? (
          <EmptyState title="Aset tidak ditemukan" desc="Tidak ada aset yang cocok dengan filter." icon={<Icon.box />} />
        ) : (
          <Table>
            <THead cols={['Kode', 'Nama Aset', 'Kategori', 'Merek', 'Serial', 'Lokasi', 'Status', 'Aksi']} />
            <tbody>
              {rows.map((a) => (
                <TR key={a.code} onClick={() => onView(a)}>
                  <TD><Code>{a.code}</Code></TD>
                  <TD className="font-medium">{a.name}</TD>
                  <TD>{a.category}</TD>
                  <TD>{a.brand}</TD>
                  <TD className="font-mono text-xs text-[var(--color-text-2)]">{a.serial}</TD>
                  <TD>{a.location}</TD>
                  <TD><StatusBadge status={a.status} /></TD>
                  <TD>
                    <div className="flex gap-0.5" onClick={(e) => e.stopPropagation()}>
                      <IconButton label="Lihat" onClick={() => onView(a)}><Icon.eye width={17} /></IconButton>
                      <IconButton label="Ubah"><Icon.edit width={17} /></IconButton>
                      <IconButton label="Riwayat" onClick={() => onView(a)}><Icon.history width={17} /></IconButton>
                    </div>
                  </TD>
                </TR>
              ))}
            </tbody>
          </Table>
        )}
      </Card>

      <Modal
        open={open}
        onClose={() => setOpen(false)}
        title="Tambah Aset"
        size="lg"
        footer={
          <>
            <Button variant="secondary" onClick={() => setOpen(false)}>Batal</Button>
            <Button onClick={() => { setOpen(false); toast('Aset baru berhasil ditambahkan.') }}>Simpan Aset</Button>
          </>
        }
      >
        <div className="grid sm:grid-cols-2 gap-4">
          <Field label="Kode Aset" required><Input placeholder="LAB-LT-014" /></Field>
          <Field label="Nama Aset" required><Input placeholder="Dell Latitude 5490" /></Field>
          <Field label="Kategori" required>
            <Select><option>Laptop</option><option>Desktop</option><option>Perangkat Jaringan</option></Select>
          </Field>
          <Field label="Merek"><Input placeholder="Dell" /></Field>
          <Field label="Model"><Input placeholder="Latitude 5490" /></Field>
          <Field label="Serial Number"><Input placeholder="DL5490-0100" /></Field>
          <Field label="Lokasi">
            <Select><option>Lab Komputer 1</option><option>Lab Komputer 2</option><option>Gudang</option></Select>
          </Field>
          <Field label="Status">
            <Select><option>Available</option><option>Maintenance</option><option>Retired</option></Select>
          </Field>
        </div>
      </Modal>
    </>
  )
}

/* --------------------------------- Asset detail (tabbed) --------------------------------- */
function AssetDetail({ asset, go }: { asset: Asset | null; go: (r: string) => void }) {
  const a = asset ?? assets[0]
  const toast = useToast()
  const [tab, setTab] = useState('Ringkasan')
  return (
    <>
      <PageHeader
        title={a.name}
        breadcrumb={['Aset', a.code]}
        subtitle={a.category}
        actions={
          <>
            <Button variant="secondary" icon={<Icon.arrowLeft width={16} />} onClick={() => go('assets')}>Kembali</Button>
            <Button variant="secondary" icon={<Icon.wrench width={16} />} onClick={() => toast('Tiket pemeliharaan dibuat.')}>Buat Pemeliharaan</Button>
            <Button icon={<Icon.edit width={16} />}>Ubah</Button>
          </>
        }
      />

      <Card className="mb-6 flex flex-wrap items-center gap-4">
        <span className="inline-flex items-center justify-center w-14 h-14 rounded-xl bg-[var(--color-primary-soft)] text-[var(--color-primary)]">
          <Icon.laptop width={28} height={28} />
        </span>
        <div className="flex-1 min-w-[200px]">
          <div className="flex items-center gap-2">
            <h3 className="font-semibold text-lg">{a.name}</h3>
            <StatusBadge status={a.status} />
          </div>
          <p className="text-sm text-[var(--color-text-2)] mt-0.5">Kode Aset: <Code>{a.code}</Code> · Serial <span className="font-mono">{a.serial}</span></p>
        </div>
      </Card>

      <Tabs tabs={['Ringkasan', 'Riwayat Alokasi', 'Riwayat Pemeliharaan', 'Log Aktivitas']} active={tab} onChange={setTab} />

      {tab === 'Ringkasan' && (
        <div className="grid gap-6 lg:grid-cols-2">
          <Card>
            <h4 className="font-semibold mb-4">Informasi Umum</h4>
            <div className="rounded-lg border border-[var(--color-border)] divide-y divide-[var(--color-border)]">
              <RowLine k="Kode Aset" v={<Code>{a.code}</Code>} />
              <RowLine k="Kategori" v={a.category} />
              <RowLine k="Merek" v={a.brand} />
              <RowLine k="Model" v={a.model} />
              <RowLine k="Lokasi" v={a.location} />
              <RowLine k="Status" v={<StatusBadge status={a.status} />} />
            </div>
          </Card>
          <Card>
            <h4 className="font-semibold mb-4">Spesifikasi Teknis</h4>
            <div className="rounded-lg border border-[var(--color-border)] divide-y divide-[var(--color-border)]">
              <RowLine k="Prosesor" v="Intel Core i7-8650U" />
              <RowLine k="RAM" v="16 GB DDR4" />
              <RowLine k="Penyimpanan" v="512 GB SSD" />
              <RowLine k="Sistem Operasi" v="Windows 11 Pro" />
              <RowLine k="Serial Number" v={<span className="font-mono">{a.serial}</span>} />
            </div>
            <div className="mt-4">
              {a.status === 'Allocated' ? (
                <Alert tone="info" title="Alokasi Saat Ini">Dialokasikan kepada Ahmad Fauzan · Sesi Pemrograman Web.</Alert>
              ) : a.status === 'Maintenance' ? (
                <Alert tone="warning" title="Dalam Pemeliharaan">Aset tidak tersedia untuk dialokasikan.</Alert>
              ) : a.status === 'Available' ? (
                <Alert tone="success" title="Tersedia">Aset siap untuk dialokasikan.</Alert>
              ) : (
                <Alert tone="neutral" title="Ditarik">Aset tidak lagi digunakan.</Alert>
              )}
            </div>
          </Card>
        </div>
      )}

      {tab === 'Riwayat Alokasi' && (
        <Card padded={false}>
          <Table>
            <THead cols={['Tanggal', 'Mahasiswa', 'Sesi', 'Keluar', 'Kembali', 'Kondisi', 'Status']} />
            <tbody>
              {[
                ['18 Sep 2026', 'Ahmad Fauzan', 'Pemrograman Web', '13.04', '—', 'Baik', 'Aktif'],
                ['11 Sep 2026', 'Siti Rahma', 'Pemrograman Web', '13.05', '15.29', 'Baik', 'Selesai'],
                ['28 Agu 2026', 'Nabila Putri', 'Basis Data', '08.06', '10.33', 'Perlu Perbaikan', 'Selesai'],
              ].map((r, i) => (
                <TR key={i}>
                  <TD>{r[0]}</TD><TD className="font-medium">{r[1]}</TD><TD>{r[2]}</TD>
                  <TD>{r[3]}</TD><TD>{r[4]}</TD>
                  <TD><Badge tone={r[5] === 'Baik' ? 'success' : 'warning'}>{r[5]}</Badge></TD>
                  <TD><StatusBadge status={r[6]} /></TD>
                </TR>
              ))}
            </tbody>
          </Table>
        </Card>
      )}

      {tab === 'Riwayat Pemeliharaan' && (
        <Card padded={false}>
          <Table>
            <THead cols={['ID', 'Keluhan', 'Prioritas', 'Status', 'Dilaporkan', 'Selesai']} />
            <tbody>
              {[
                ['MT-190', 'Baterai lemah', 'Sedang', 'Selesai', '20 Agu 2026', '23 Agu 2026'],
                ['MT-142', 'Pembersihan rutin', 'Rendah', 'Selesai', '02 Agu 2026', '02 Agu 2026'],
              ].map((r, i) => (
                <TR key={i}>
                  <TD><Code>{r[0]}</Code></TD><TD>{r[1]}</TD>
                  <TD><StatusBadge status={r[2]} /></TD>
                  <TD><StatusBadge status={r[3]} /></TD>
                  <TD>{r[4]}</TD><TD>{r[5]}</TD>
                </TR>
              ))}
            </tbody>
          </Table>
        </Card>
      )}

      {tab === 'Log Aktivitas' && (
        <Card>
          <Timeline
            items={[
              { title: 'Dialokasikan ke Ahmad Fauzan', time: '18 Sep 2026 13.04', tone: 'info' },
              { title: 'Dikembalikan oleh Siti Rahma', time: '11 Sep 2026 15.29', tone: 'success' },
              { title: 'Pemeliharaan selesai (MT-190)', time: '23 Agu 2026', tone: 'success' },
              { title: 'Aset didaftarkan', time: '20 Agu 2026', tone: 'neutral' },
            ]}
          />
        </Card>
      )}
    </>
  )
}

/* --------------------------------- User detail --------------------------------- */
function UserDetail({ user, go }: { user: User | null; go: (r: string) => void }) {
  const toast = useToast()
  const u = user ?? users[0]
  const isStudent = u.role === 'Mahasiswa'
  return (
    <>
      <PageHeader
        title={u.name}
        breadcrumb={['Pengguna', u.nim]}
        subtitle={u.role}
        actions={
          <>
            <Button variant="secondary" icon={<Icon.arrowLeft width={16} />} onClick={() => go('users')}>Kembali</Button>
            <Button
              variant={u.status === 'Aktif' ? 'danger' : 'success'}
              onClick={() => toast(u.status === 'Aktif' ? 'Pengguna dinonaktifkan.' : 'Pengguna diaktifkan.', u.status === 'Aktif' ? 'error' : 'success')}
            >
              {u.status === 'Aktif' ? 'Nonaktifkan' : 'Aktifkan'}
            </Button>
            <Button icon={<Icon.edit width={16} />}>Ubah</Button>
          </>
        }
      />
      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="text-center">
          <Avatar name={u.name} size={72} />
          <h3 className="mt-3 font-semibold text-lg">{u.name}</h3>
          <p className="text-sm text-[var(--color-text-2)] font-mono">{u.nim}</p>
          <div className="mt-3 flex justify-center gap-2 flex-wrap">
            <StatusBadge status={u.status} />
            {isStudent && u.face !== '—' && <StatusBadge status={u.face} />}
          </div>
        </Card>

        <div className="lg:col-span-2 space-y-6">
          <Card>
            <h4 className="font-semibold mb-4">Informasi Akun</h4>
            <div className="rounded-lg border border-[var(--color-border)] divide-y divide-[var(--color-border)]">
              <RowLine k="Nama Lengkap" v={u.name} />
              <RowLine k="Identifier" v={<Code>{u.nim}</Code>} />
              <RowLine k="Email" v={u.email} />
              <RowLine k="Peran" v={u.role} />
              <RowLine k="Status" v={<StatusBadge status={u.status} />} />
              <RowLine k="Dibuat" v={u.created} />
            </div>
          </Card>

          {isStudent && (
            <Card>
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="font-semibold">Pendaftaran Wajah</h4>
                  <p className="text-sm text-[var(--color-text-2)] mt-0.5">
                    {u.face === 'Terdaftar' ? 'Templat wajah aktif untuk verifikasi presensi.' : 'Mahasiswa belum mendaftarkan wajah.'}
                  </p>
                </div>
                <StatusBadge status={u.face} />
              </div>
              {u.face !== 'Terdaftar' && (
                <div className="mt-4">
                  <Alert tone="warning">Mahasiswa tidak dapat melakukan presensi hingga wajah didaftarkan.</Alert>
                </div>
              )}
            </Card>
          )}
        </div>
      </div>
    </>
  )
}

/* --------------------------------- Session detail --------------------------------- */
function SessionDetail({ session, go }: { session: Session | null; go: (r: string) => void }) {
  const toast = useToast()
  const s = session ?? sessions[0]
  const [tab, setTab] = useState('Informasi')
  const present = attendanceMonitoring.filter((r) => r.status === 'Hadir').length
  return (
    <>
      <PageHeader
        title={s.name}
        breadcrumb={['Sesi Praktikum', s.code]}
        subtitle={`${s.course} · ${s.room}`}
        actions={
          <>
            <Button variant="secondary" icon={<Icon.arrowLeft width={16} />} onClick={() => go('sessions')}>Kembali</Button>
            <Button variant="secondary" onClick={() => toast('Status sesi diperbarui.')}>Ubah Status</Button>
            <Button icon={<Icon.edit width={16} />}>Ubah Sesi</Button>
          </>
        }
      />

      <div className="grid gap-4 sm:grid-cols-4 mb-6">
        <StatCard label="Status" value={<StatusBadge status={s.status} />} icon={<Icon.calendar />} tone="primary" />
        <StatCard label="Hadir" value={`${present}/${attendanceMonitoring.length}`} icon={<Icon.checkCircle />} tone="success" />
        <StatCard label="Waktu" value={`${s.start}–${s.end}`} hint={s.date} icon={<Icon.clock />} tone="info" />
        <StatCard label="Laboratorium" value={s.room} icon={<Icon.box />} tone="neutral" />
      </div>

      <Tabs tabs={['Informasi', 'Peserta & Kehadiran']} active={tab} onChange={setTab} />

      {tab === 'Informasi' ? (
        <Card className="max-w-2xl">
          <h4 className="font-semibold mb-4">Detail Sesi</h4>
          <div className="rounded-lg border border-[var(--color-border)] divide-y divide-[var(--color-border)]">
            <RowLine k="Kode Sesi" v={<Code>{s.code}</Code>} />
            <RowLine k="Nama Sesi" v={s.name} />
            <RowLine k="Mata Kuliah" v={s.course} />
            <RowLine k="Tanggal" v={s.date} />
            <RowLine k="Waktu" v={`${s.start}–${s.end}`} />
            <RowLine k="Laboratorium" v={s.room} />
            <RowLine k="Status" v={<StatusBadge status={s.status} />} />
          </div>
        </Card>
      ) : (
        <Card padded={false}>
          <Table>
            <THead cols={['Mahasiswa', 'NIM', 'Check-in', 'Verifikasi', 'Presensi']} />
            <tbody>
              {attendanceMonitoring.map((r) => (
                <TR key={r.nim}>
                  <TD>
                    <div className="flex items-center gap-2.5"><Avatar name={r.student} size={30} /><span className="font-medium">{r.student}</span></div>
                  </TD>
                  <TD><Code>{r.nim}</Code></TD>
                  <TD>{r.checkIn}</TD>
                  <TD>{r.verify === '—' ? <span className="text-[var(--color-text-3)]">—</span> : <StatusBadge status={r.verify} />}</TD>
                  <TD><StatusBadge status={r.status} /></TD>
                </TR>
              ))}
            </tbody>
          </Table>
        </Card>
      )}
    </>
  )
}

/* --------------------------------- Allocations (create + return) --------------------------------- */
function Allocations({ onView }: { onView: (a: Allocation) => void }) {
  const toast = useToast()
  const [createOpen, setCreateOpen] = useState(false)
  const [returnRow, setReturnRow] = useState<(typeof allocations)[number] | null>(null)
  const [student, setStudent] = useState('Nabila Putri')
  const [session, setSession] = useState('Pemrograman Web')
  const [asset, setAsset] = useState('LAB-LT-008')
  const available = assets.filter((a) => a.status === 'Available')

  return (
    <>
      <PageHeader
        title="Alokasi Aset"
        subtitle="Kelola alokasi dan pengembalian aset per mahasiswa dan sesi."
        actions={<Button icon={<Icon.plus width={16} />} onClick={() => setCreateOpen(true)}>Alokasi Baru</Button>}
      />
      <Card padded={false}>
        <Table>
          <THead cols={['ID', 'Mahasiswa', 'Sesi', 'Aset', 'Waktu', 'Kondisi', 'Status', 'Aksi']} />
          <tbody>
            {allocations.map((a) => (
              <TR key={a.id} onClick={() => onView(a)}>
                <TD><Code>{a.id}</Code></TD>
                <TD>
                  <div className="font-medium">{a.student}</div>
                  <div className="text-xs text-[var(--color-text-3)] font-mono">{a.nim}</div>
                </TD>
                <TD>{a.session}</TD>
                <TD><Code>{a.asset}</Code><div className="text-xs text-[var(--color-text-3)]">{a.assetName}</div></TD>
                <TD>{a.time}</TD>
                <TD><Badge tone="success">{a.condition}</Badge></TD>
                <TD><StatusBadge status={a.status} /></TD>
                <TD>
                  <div onClick={(e) => e.stopPropagation()}>
                    {a.status === 'Aktif' ? (
                      <Button size="sm" variant="secondary" icon={<Icon.swap width={15} />} onClick={() => setReturnRow(a)}>
                        Kembalikan
                      </Button>
                    ) : (
                      <IconButton label="Lihat" onClick={() => onView(a)}><Icon.eye width={17} /></IconButton>
                    )}
                  </div>
                </TD>
              </TR>
            ))}
          </tbody>
        </Table>
      </Card>

      {/* Create allocation */}
      <Modal
        open={createOpen}
        onClose={() => setCreateOpen(false)}
        title="Buat Alokasi Aset"
        size="lg"
        footer={
          <>
            <Button variant="secondary" onClick={() => setCreateOpen(false)}>Batal</Button>
            <Button
              icon={<Icon.swap width={16} />}
              onClick={() => { setCreateOpen(false); toast('Aset berhasil dialokasikan.') }}
            >
              Alokasikan Aset
            </Button>
          </>
        }
      >
        <div className="grid sm:grid-cols-3 gap-4">
          <Field label="Mahasiswa" required>
            <Select value={student} onChange={(e) => setStudent(e.target.value)}>
              {users.filter((u) => u.role === 'Mahasiswa').map((u) => <option key={u.nim}>{u.name}</option>)}
            </Select>
          </Field>
          <Field label="Sesi" required>
            <Select value={session} onChange={(e) => setSession(e.target.value)}>
              <option>Pemrograman Web</option><option>Basis Data</option><option>Jaringan Komputer</option>
            </Select>
          </Field>
          <Field label="Aset Tersedia" required>
            <Select value={asset} onChange={(e) => setAsset(e.target.value)}>
              {available.map((a) => <option key={a.code} value={a.code}>{a.code}</option>)}
            </Select>
          </Field>
        </div>
        <div className="grid sm:grid-cols-2 gap-4 mt-4">
          <Field label="Waktu Alokasi"><Input type="datetime-local" defaultValue="2026-09-18T13:10" /></Field>
          <Field label="Kondisi Sebelum">
            <Select><option>Baik</option><option>Perlu Perhatian</option></Select>
          </Field>
        </div>
        <Field label="Catatan"><Input placeholder="Opsional" /></Field>

        <div className="mt-5 rounded-lg border border-[var(--color-border)] bg-[var(--color-surface-2)] p-4">
          <p className="text-xs font-semibold uppercase tracking-wide text-[var(--color-text-3)] mb-3">Ringkasan Alokasi</p>
          <div className="grid grid-cols-3 gap-3 text-sm">
            <div><p className="text-[var(--color-text-2)] text-xs">Mahasiswa</p><p className="font-medium">{student}</p></div>
            <div><p className="text-[var(--color-text-2)] text-xs">Sesi</p><p className="font-medium">{session}</p></div>
            <div><p className="text-[var(--color-text-2)] text-xs">Aset</p><p className="font-medium font-mono">{asset}</p></div>
          </div>
        </div>
        {available.length === 0 && (
          <div className="mt-4"><Alert tone="warning" title="Aset tidak tersedia">Aset tidak tersedia untuk dialokasikan.</Alert></div>
        )}
      </Modal>

      {/* Return allocation */}
      <Modal
        open={!!returnRow}
        onClose={() => setReturnRow(null)}
        title="Pengembalian Aset"
        footer={
          <>
            <Button variant="secondary" onClick={() => setReturnRow(null)}>Batal</Button>
            <Button
              variant="success"
              icon={<Icon.check width={16} />}
              onClick={() => { toast(`Aset ${returnRow?.asset} berhasil dikembalikan. Status: Available.`); setReturnRow(null) }}
            >
              Konfirmasi Pengembalian
            </Button>
          </>
        }
      >
        {returnRow && (
          <>
            <div className="rounded-lg border border-[var(--color-border)] divide-y divide-[var(--color-border)] mb-4">
              <RowLine k="Mahasiswa" v={returnRow.student} />
              <RowLine k="Sesi" v={returnRow.session} />
              <RowLine k="Aset" v={<Code>{returnRow.asset}</Code>} />
              <RowLine k="Waktu Alokasi" v={returnRow.time} />
              <RowLine k="Kondisi Sebelum" v={<Badge tone="success">{returnRow.condition}</Badge>} />
            </div>
            <div className="grid sm:grid-cols-2 gap-4">
              <Field label="Kondisi Setelah" required>
                <Select><option>Baik</option><option>Perlu Perbaikan</option></Select>
              </Field>
              <Field label="Catatan"><Input placeholder="Opsional" /></Field>
            </div>
            <div className="mt-4"><Alert tone="info">Jika kondisi memerlukan perbaikan, status aset akan otomatis menjadi <b>Maintenance</b>.</Alert></div>
          </>
        )}
      </Modal>
    </>
  )
}
function RowLine({ k, v }: { k: string; v: React.ReactNode }) {
  return (
    <div className="flex items-center justify-between gap-3 px-4 py-2.5 text-sm">
      <span className="text-[var(--color-text-2)]">{k}</span>
      <span className="font-medium text-right">{v}</span>
    </div>
  )
}

/* --------------------------------- Maintenance --------------------------------- */
function MaintenanceMgmt({ onView }: { onView: (m: Maintenance) => void }) {
  const toast = useToast()
  const [open, setOpen] = useState(false)
  const counts = { Terbuka: 0, 'Sedang Dikerjakan': 0, Selesai: 0 } as Record<string, number>
  maintenance.forEach((m) => (counts[m.status] !== undefined && (counts[m.status] += 1)))
  return (
    <>
      <PageHeader
        title="Pemeliharaan Aset"
        subtitle="Kelola perbaikan dan pemeliharaan perangkat laboratorium."
        actions={<Button icon={<Icon.plus width={16} />} onClick={() => setOpen(true)}>Buat Pemeliharaan</Button>}
      />
      <div className="grid gap-4 sm:grid-cols-3 mb-6">
        <StatCard label="Terbuka" value={counts.Terbuka} icon={<Icon.alert />} tone="error" />
        <StatCard label="Sedang Dikerjakan" value={counts['Sedang Dikerjakan']} icon={<Icon.wrench />} tone="warning" />
        <StatCard label="Selesai" value={counts.Selesai} icon={<Icon.checkCircle />} tone="success" />
      </div>
      <Card padded={false}>
        <Table>
          <THead cols={['ID', 'Kode Aset', 'Nama', 'Keluhan', 'Prioritas', 'Status', 'Dilaporkan', 'Selesai', 'Aksi']} />
          <tbody>
            {maintenance.map((m) => (
              <TR key={m.id} onClick={() => onView(m)}>
                <TD><Code>{m.id}</Code></TD>
                <TD><Code>{m.code}</Code></TD>
                <TD>{m.name}</TD>
                <TD>{m.issue}</TD>
                <TD><StatusBadge status={m.priority} /></TD>
                <TD><StatusBadge status={m.status} /></TD>
                <TD>{m.reported}</TD>
                <TD>{m.completed}</TD>
                <TD>
                  <div onClick={(e) => e.stopPropagation()}>
                    {m.status !== 'Selesai' ? (
                      <Button size="sm" variant="secondary" onClick={() => toast(`Pemeliharaan ${m.id} diselesaikan. Aset kembali tersedia.`)}>
                        Selesaikan
                      </Button>
                    ) : (
                      <IconButton label="Lihat" onClick={() => onView(m)}><Icon.eye width={17} /></IconButton>
                    )}
                  </div>
                </TD>
              </TR>
            ))}
          </tbody>
        </Table>
      </Card>

      <Modal
        open={open}
        onClose={() => setOpen(false)}
        title="Buat Tiket Pemeliharaan"
        footer={
          <>
            <Button variant="secondary" onClick={() => setOpen(false)}>Batal</Button>
            <Button onClick={() => { setOpen(false); toast('Tiket pemeliharaan dibuat. Aset ditandai Maintenance.') }}>Simpan</Button>
          </>
        }
      >
        <div className="space-y-4">
          <Field label="Aset" required>
            <Select>{assets.map((a) => <option key={a.code}>{a.code} — {a.name}</option>)}</Select>
          </Field>
          <div className="grid sm:grid-cols-2 gap-4">
            <Field label="Jenis Keluhan" required>
              <Select><option>Perangkat Keras</option><option>Perangkat Lunak</option><option>Jaringan</option></Select>
            </Field>
            <Field label="Prioritas" required>
              <Select><option>Rendah</option><option>Sedang</option><option>Tinggi</option></Select>
            </Field>
          </div>
          <Field label="Deskripsi" required><Input placeholder="Jelaskan permasalahan…" /></Field>
        </div>
      </Modal>
    </>
  )
}

/* --------------------------------- Allocation detail --------------------------------- */
function AllocationDetail({ alloc, go }: { alloc: Allocation | null; go: (r: string) => void }) {
  const toast = useToast()
  const a = alloc ?? allocations[0]
  const active = a.status === 'Aktif'
  return (
    <>
      <PageHeader
        title={`Alokasi ${a.id}`}
        breadcrumb={['Alokasi', a.id]}
        subtitle={`${a.student} · ${a.session}`}
        actions={
          <>
            <Button variant="secondary" icon={<Icon.arrowLeft width={16} />} onClick={() => go('allocations')}>Kembali</Button>
            {active && (
              <Button icon={<Icon.swap width={16} />} onClick={() => toast('Buka alur pengembalian aset.')}>Kembalikan Aset</Button>
            )}
          </>
        }
      />
      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <div className="flex items-center justify-between mb-4">
            <h4 className="font-semibold">Ringkasan Alokasi</h4>
            <StatusBadge status={a.status} />
          </div>
          <div className="grid sm:grid-cols-3 gap-3 mb-4">
            <div className="rounded-lg border border-[var(--color-border)] p-3">
              <p className="text-xs text-[var(--color-text-2)] mb-1">Mahasiswa</p>
              <div className="flex items-center gap-2"><Avatar name={a.student} size={28} /><span className="font-medium text-sm">{a.student}</span></div>
            </div>
            <div className="rounded-lg border border-[var(--color-border)] p-3">
              <p className="text-xs text-[var(--color-text-2)] mb-1">Sesi</p>
              <p className="font-medium text-sm">{a.session}</p>
            </div>
            <div className="rounded-lg border border-[var(--color-border)] p-3">
              <p className="text-xs text-[var(--color-text-2)] mb-1">Aset</p>
              <p className="font-medium text-sm font-mono">{a.asset}</p>
              <p className="text-xs text-[var(--color-text-3)]">{a.assetName}</p>
            </div>
          </div>
          <div className="rounded-lg border border-[var(--color-border)] divide-y divide-[var(--color-border)]">
            <RowLine k="ID Alokasi" v={<Code>{a.id}</Code>} />
            <RowLine k="NIM" v={<Code>{a.nim}</Code>} />
            <RowLine k="Waktu Alokasi" v={a.time} />
            <RowLine k="Kondisi Sebelum" v={<Badge tone="success">{a.condition}</Badge>} />
            <RowLine k="Waktu Pengembalian" v={active ? '—' : '15.29'} />
            <RowLine k="Kondisi Setelah" v={active ? '—' : <Badge tone="success">Baik</Badge>} />
          </div>
        </Card>
        <Card>
          <h4 className="font-semibold mb-4">Linimasa</h4>
          <Timeline
            items={
              active
                ? [
                    { title: 'Alokasi dibuat', time: a.time, tone: 'info' },
                    { title: 'Sedang digunakan', time: '—', tone: 'warning' },
                  ]
                : [
                    { title: 'Alokasi dibuat', time: a.time, tone: 'info' },
                    { title: 'Aset dikembalikan', time: '15.29', tone: 'success' },
                    { title: 'Status: Available', time: '15.29', tone: 'success' },
                  ]
            }
          />
        </Card>
      </div>
    </>
  )
}

/* --------------------------------- Maintenance detail --------------------------------- */
function MaintenanceDetail({ mt, go }: { mt: Maintenance | null; go: (r: string) => void }) {
  const toast = useToast()
  const m = mt ?? maintenance[0]
  const done = m.status === 'Selesai'
  return (
    <>
      <PageHeader
        title={`Pemeliharaan ${m.id}`}
        breadcrumb={['Pemeliharaan', m.id]}
        subtitle={`${m.code} · ${m.name}`}
        actions={
          <>
            <Button variant="secondary" icon={<Icon.arrowLeft width={16} />} onClick={() => go('maintenance')}>Kembali</Button>
            {!done && (
              <Button variant="success" icon={<Icon.check width={16} />} onClick={() => toast(`Pemeliharaan ${m.id} diselesaikan. Aset kembali tersedia.`)}>
                Selesaikan
              </Button>
            )}
          </>
        }
      />
      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <div className="flex items-center justify-between mb-4">
            <h4 className="font-semibold">Detail Tiket</h4>
            <div className="flex gap-2"><StatusBadge status={m.priority} /><StatusBadge status={m.status} /></div>
          </div>
          <div className="rounded-lg border border-[var(--color-border)] divide-y divide-[var(--color-border)]">
            <RowLine k="ID Tiket" v={<Code>{m.id}</Code>} />
            <RowLine k="Aset" v={<><Code>{m.code}</Code> <span className="text-[var(--color-text-2)]">{m.name}</span></>} />
            <RowLine k="Keluhan" v={m.issue} />
            <RowLine k="Prioritas" v={<StatusBadge status={m.priority} />} />
            <RowLine k="Status" v={<StatusBadge status={m.status} />} />
            <RowLine k="Dilaporkan" v={m.reported} />
            <RowLine k="Selesai" v={m.completed} />
          </div>
          <div className="mt-4">
            {done ? (
              <Alert tone="success" title="Pemeliharaan selesai">Aset telah dikembalikan ke status tersedia.</Alert>
            ) : (
              <Alert tone="warning" title="Aset tidak tersedia">Selama pemeliharaan, aset ini tidak dapat dialokasikan.</Alert>
            )}
          </div>
        </Card>
        <Card>
          <h4 className="font-semibold mb-4">Linimasa</h4>
          <Timeline
            items={
              done
                ? [
                    { title: 'Tiket dibuat', time: m.reported, tone: 'error' },
                    { title: 'Sedang dikerjakan', time: m.reported, tone: 'warning' },
                    { title: 'Selesai', time: m.completed, tone: 'success' },
                  ]
                : [
                    { title: 'Tiket dibuat', time: m.reported, tone: 'error' },
                    { title: 'Sedang dikerjakan', time: m.reported, tone: 'warning' },
                    { title: 'Menunggu penyelesaian', time: '—', tone: 'neutral' },
                  ]
            }
          />
        </Card>
      </div>
    </>
  )
}

/* --------------------------------- Reports --------------------------------- */
type ReportDef = {
  metrics: { label: string; value: string }[]
  cols: string[]
  rows: React.ReactNode[][]
  extraFilter?: React.ReactNode
}

function Reports() {
  const toast = useToast()
  const [tab, setTab] = useState('Presensi')
  const loading = useSimulatedLoading()

  const pct = (v: string) => <Badge tone={Number(v.replace('%', '')) >= 85 ? 'success' : 'warning'}>{v}</Badge>

  const defs: Record<string, ReportDef> = {
    Presensi: {
      metrics: [
        { label: 'Total Kehadiran', value: '124' },
        { label: 'Rata-rata Kehadiran', value: '86%' },
        { label: 'Verifikasi Gagal', value: '3' },
        { label: 'Sesi Tercatat', value: '18' },
      ],
      cols: ['Tanggal', 'Sesi', 'Laboratorium', 'Hadir', 'Tidak Hadir', 'Kehadiran'],
      rows: [
        ['18 Sep 2026', 'Pemrograman Web', 'Lab Komputer 1', '22', '8', pct('73%')],
        ['12 Sep 2026', 'Basis Data', 'Lab Komputer 2', '28', '2', pct('93%')],
        ['11 Sep 2026', 'Pemrograman Web', 'Lab Komputer 1', '26', '4', pct('87%')],
        ['05 Sep 2026', 'Basis Data', 'Lab Komputer 2', '27', '3', pct('90%')],
      ],
    },
    Aset: {
      metrics: [
        { label: 'Total Aset', value: String(assets.length) },
        { label: 'Tersedia', value: String(assets.filter((a) => a.status === 'Available').length) },
        { label: 'Dialokasikan', value: String(assets.filter((a) => a.status === 'Allocated').length) },
        { label: 'Pemeliharaan', value: String(assets.filter((a) => a.status === 'Maintenance').length) },
      ],
      cols: ['Kode', 'Nama Aset', 'Kategori', 'Lokasi', 'Status'],
      rows: assets.map((a) => [<Code>{a.code}</Code>, a.name, a.category, a.location, <StatusBadge status={a.status} />]),
      extraFilter: (
        <Field label="Kategori">
          <Select className="w-44"><option>Semua Kategori</option><option>Laptop</option><option>Desktop</option><option>Perangkat Jaringan</option></Select>
        </Field>
      ),
    },
    Alokasi: {
      metrics: [
        { label: 'Total Alokasi', value: String(allocations.length) },
        { label: 'Aktif', value: String(allocations.filter((a) => a.status === 'Aktif').length) },
        { label: 'Selesai', value: String(allocations.filter((a) => a.status === 'Selesai').length) },
        { label: 'Aset Terpakai', value: String(new Set(allocations.map((a) => a.asset)).size) },
      ],
      cols: ['ID', 'Mahasiswa', 'Sesi', 'Aset', 'Waktu', 'Status'],
      rows: allocations.map((a) => [<Code>{a.id}</Code>, a.student, a.session, <Code>{a.asset}</Code>, a.time, <StatusBadge status={a.status} />]),
    },
    Pemeliharaan: {
      metrics: [
        { label: 'Total Tiket', value: String(maintenance.length) },
        { label: 'Terbuka', value: String(maintenance.filter((m) => m.status === 'Terbuka').length) },
        { label: 'Dikerjakan', value: String(maintenance.filter((m) => m.status === 'Sedang Dikerjakan').length) },
        { label: 'Selesai', value: String(maintenance.filter((m) => m.status === 'Selesai').length) },
      ],
      cols: ['ID', 'Aset', 'Keluhan', 'Prioritas', 'Status', 'Dilaporkan'],
      rows: maintenance.map((m) => [<Code>{m.id}</Code>, <><Code>{m.code}</Code> <span className="text-[var(--color-text-2)]">{m.name}</span></>, m.issue, <StatusBadge status={m.priority} />, <StatusBadge status={m.status} />, m.reported]),
      extraFilter: (
        <Field label="Prioritas">
          <Select className="w-40"><option>Semua Prioritas</option><option>Tinggi</option><option>Sedang</option><option>Rendah</option></Select>
        </Field>
      ),
    },
  }

  const def = defs[tab]

  return (
    <>
      <PageHeader title="Laporan" subtitle="Rekapitulasi presensi, aset, alokasi, dan pemeliharaan." />
      <Tabs tabs={['Presensi', 'Aset', 'Alokasi', 'Pemeliharaan']} active={tab} onChange={setTab} />

      <Card padded={false}>
        <div className="flex flex-wrap items-end gap-3 p-4 border-b border-[var(--color-border)]">
          <Field label="Dari Tanggal"><Input type="date" defaultValue="2026-09-01" className="w-40" /></Field>
          <Field label="Sampai Tanggal"><Input type="date" defaultValue="2026-09-18" className="w-40" /></Field>
          {def.extraFilter ?? (
            <Field label="Sesi">
              <Select className="w-48"><option>Semua Sesi</option><option>Pemrograman Web</option><option>Basis Data</option></Select>
            </Field>
          )}
          <div className="ml-auto flex gap-2">
            <Button variant="ghost">Reset</Button>
            <Button icon={<Icon.download width={16} />} onClick={() => toast(`Laporan ${tab} diekspor ke CSV.`)}>Ekspor CSV</Button>
          </div>
        </div>

        <div className="grid gap-4 sm:grid-cols-4 p-4 border-b border-[var(--color-border)]">
          {def.metrics.map((m) => (
            <Metric key={m.label} label={m.label} value={m.value} />
          ))}
        </div>

        <Table>
          <THead cols={def.cols} />
          {loading ? (
            <SkeletonRows cols={def.cols.length} />
          ) : (
            <tbody>
              {def.rows.map((r, i) => (
                <TR key={i}>
                  {r.map((cell, j) => (
                    <TD key={j} className={j === 1 ? 'font-medium' : ''}>{cell}</TD>
                  ))}
                </TR>
              ))}
            </tbody>
          )}
        </Table>
        <Pagination page={1} pages={1} total={def.rows.length} onPage={() => {}} />
      </Card>
    </>
  )
}
function Metric({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-lg border border-[var(--color-border)] p-3">
      <p className="text-xs text-[var(--color-text-2)]">{label}</p>
      <p className="mt-1 text-xl font-semibold tracking-tight">{value}</p>
    </div>
  )
}

/* --------------------------------- Audit log --------------------------------- */
function AuditLog() {
  return (
    <>
      <PageHeader title="Log Audit" subtitle="Jejak aktivitas operasional sistem." />
      <Card padded={false}>
        <div className="flex flex-wrap items-center gap-3 p-4 border-b border-[var(--color-border)]">
          <div className="flex-1 min-w-[200px]"><Input placeholder="Cari pengguna atau aksi…" /></div>
          <Select className="w-48"><option>Semua Aksi</option><option>Login</option><option>Face Verification</option><option>Allocation Created</option></Select>
        </div>
        <Table>
          <THead cols={['Waktu', 'Pengguna', 'Aksi', 'Entitas', 'Status', 'IP / Perangkat']} />
          <tbody>
            {auditLog.map((a, i) => (
              <TR key={i}>
                <TD className="whitespace-nowrap text-[var(--color-text-2)]">{a.time}</TD>
                <TD className="font-medium">{a.user}</TD>
                <TD>{a.action}</TD>
                <TD className="text-[var(--color-text-2)]">{a.entity}</TD>
                <TD><StatusBadge status={a.status} /></TD>
                <TD className="font-mono text-xs text-[var(--color-text-2)]">{a.ip}</TD>
              </TR>
            ))}
          </tbody>
        </Table>
        <div className="p-4 border-t border-[var(--color-border)]">
          <Alert tone="neutral">
            Informasi sensitif seperti kata sandi, token, gambar mentah, atau data biometrik tidak ditampilkan pada log.
          </Alert>
        </div>
      </Card>
    </>
  )
}

/* --------------------------------- Profile --------------------------------- */
function Profile() {
  const toast = useToast()
  return (
    <>
      <PageHeader title="Profil Admin" subtitle="Informasi akun staf laboratorium." />
      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="text-center">
          <Avatar name="Dewi Anggraini" size={72} />
          <h3 className="mt-3 font-semibold text-lg">Dewi Anggraini</h3>
          <p className="text-sm text-[var(--color-text-2)]">Laboran</p>
          <div className="mt-3 flex justify-center"><Badge tone="success">Aktif</Badge></div>
        </Card>
        <Card className="lg:col-span-2">
          <h4 className="font-semibold mb-4">Informasi Akun</h4>
          <div className="grid sm:grid-cols-2 gap-4">
            <Field label="Nama Lengkap"><Input defaultValue="Dewi Anggraini" /></Field>
            <Field label="Identifier"><Input defaultValue="LAB0001" disabled /></Field>
            <Field label="Email"><Input defaultValue="dewi.laboran@univ.ac.id" /></Field>
            <Field label="Peran"><Input defaultValue="Laboran" disabled /></Field>
          </div>
          <div className="mt-5 flex justify-end">
            <Button onClick={() => toast('Perubahan profil berhasil disimpan.')}>Simpan Perubahan</Button>
          </div>
        </Card>
      </div>
    </>
  )
}
