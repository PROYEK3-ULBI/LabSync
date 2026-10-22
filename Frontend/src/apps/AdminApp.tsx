import { useEffect, useState } from 'react'
import { api, ApiError, unwrapList } from '../lib/api'
import type { AssetAllocation, AssetRecord, AttendanceRecord, AuthUser, LabSession, MaintenanceRecord } from '../lib/api'
import { Icon } from '../components/icons'
import { Shell, type NavItem } from '../components/Shell'
import {
  Alert,
  Avatar,
  Badge,
  Button,
  Card,
  Code,
  EmptyState,
  Field,
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
  useToast,
} from '../components/ui'

const NAV: NavItem[] = [
  { key: 'dashboard', label: 'Dasbor', icon: 'dashboard' },
  { key: 'users', label: 'Pengguna', icon: 'users' },
  { key: 'sessions', label: 'Sesi Praktikum', icon: 'calendar' },
  { key: 'attendance', label: 'Presensi', icon: 'checkCircle' },
  { key: 'assets', label: 'Aset', icon: 'box' },
  { key: 'allocations', label: 'Alokasi', icon: 'swap' },
  { key: 'maintenance', label: 'Pemeliharaan', icon: 'wrench' },
  { key: 'reports', label: 'Laporan', icon: 'report' },
]

type AdminState = {
  users: AuthUser[]
  sessions: LabSession[]
  attendance: AttendanceRecord[]
  assets: AssetRecord[]
  allocations: AssetAllocation[]
  maintenance: MaintenanceRecord[]
  loading: boolean
  error: string
}

const emptyState: AdminState = {
  users: [],
  sessions: [],
  attendance: [],
  assets: [],
  allocations: [],
  maintenance: [],
  loading: true,
  error: '',
}

function formatDate(value?: string) {
  if (!value) return '—'
  const d = new Date(value)
  return Number.isNaN(d.getTime()) ? value : d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
}

function formatTime(value?: string) {
  if (!value) return '—'
  const d = new Date(value)
  return Number.isNaN(d.getTime()) ? value : d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

function formatDateTime(value?: string) {
  if (!value) return '—'
  return `${formatDate(value)}, ${formatTime(value)}`
}

function errMsg(error: unknown) {
  return error instanceof ApiError ? error.message : 'Operasi gagal diproses.'
}

function roleLabel(role?: string) {
  return role === 'ADMIN' ? 'Admin' : 'Mahasiswa'
}

function assetStatusLabel(status?: string) {
  return status === 'AVAILABLE' ? 'Available' : status === 'ALLOCATED' ? 'Allocated' : status === 'MAINTENANCE' ? 'Maintenance' : status === 'RETIRED' ? 'Retired' : status || '—'
}

function sessionStatusLabel(status?: string) {
  return status === 'ACTIVE' ? 'Aktif' : status === 'CLOSED' ? 'Ditutup' : status === 'CANCELLED' ? 'Dibatalkan' : 'Akan Datang'
}

// Extract time HH:mm from ISO string
function extractTime(iso?: string) {
  if (!iso) return ''
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

// Extract date YYYY-MM-DD from ISO or date string
function extractDate(value?: string) {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  return d.toISOString().slice(0, 10)
}

export function AdminApp({ onLogout, user }: { onLogout: () => void; user?: AuthUser | null }) {
  const [route, setRoute] = useState('dashboard')
  const [state, setState] = useState<AdminState>(emptyState)

  const refresh = async () => {
    setState((s) => ({ ...s, loading: true, error: '' }))
    try {
      const [users, sessions, attendance, assets, allocations, maintenance] = await Promise.all([
        api.users.list(),
        api.sessions.list(),
        api.attendance.list(),
        api.assets.list(),
        api.allocations.list(),
        api.maintenance.list(),
      ])
      setState({
        users: unwrapList<AuthUser>(users),
        sessions: unwrapList<LabSession>(sessions),
        attendance: unwrapList<AttendanceRecord>(attendance),
        assets: unwrapList<AssetRecord>(assets),
        allocations: unwrapList<AssetAllocation>(allocations),
        maintenance: unwrapList<MaintenanceRecord>(maintenance),
        loading: false,
        error: '',
      })
    } catch (error) {
      setState((s) => ({ ...s, loading: false, error: errMsg(error) }))
    }
  }

  useEffect(() => { void refresh() }, [route])

  return (
    <Shell
      brand="LabPresensi"
      roleLabel="Admin / Laboran"
      roleTone="slate"
      nav={NAV}
      active={route}
      onNavigate={setRoute}
      user={{ name: user?.full_name || 'Admin Demo LabSync', sub: user?.identifier || 'Admin' }}
      onLogout={onLogout}
    >
      {state.error && <div className="mb-4"><Alert tone="error" title="Gagal memuat data">{state.error}</Alert></div>}
      {route === 'dashboard' && <Dashboard data={state} go={setRoute} />}
      {route === 'users' && <Users data={state} refresh={refresh} />}
      {route === 'sessions' && <Sessions data={state} refresh={refresh} />}
      {route === 'attendance' && <AttendanceAdmin data={state} />}
      {route === 'assets' && <Assets data={state} refresh={refresh} />}
      {route === 'allocations' && <Allocations data={state} refresh={refresh} />}
      {route === 'maintenance' && <Maintenance data={state} refresh={refresh} />}
      {route === 'reports' && <Reports data={state} />}
    </Shell>
  )
}

/* --------------------------------- Dashboard --------------------------------- */
function Dashboard({ data, go }: { data: AdminState; go: (r: string) => void }) {
  const activeSessions = data.sessions.filter((s) => s.status === 'ACTIVE')
  const available = data.assets.filter((a) => a.status === 'AVAILABLE').length
  const allocated = data.assets.filter((a) => a.status === 'ALLOCATED').length
  const maintenance = data.assets.filter((a) => a.status === 'MAINTENANCE').length
  return (
    <>
      <PageHeader title="Dasbor Operasional" subtitle="Ringkasan data backend untuk presensi dan aset laboratorium." />
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-5">
        <StatCard label="Sesi Aktif" value={activeSessions.length} icon={<Icon.calendar />} tone="primary" />
        <StatCard label="Presensi" value={data.attendance.length} icon={<Icon.checkCircle />} tone="success" />
        <StatCard label="Aset Tersedia" value={available} icon={<Icon.box />} tone="info" />
        <StatCard label="Dialokasikan" value={allocated} icon={<Icon.laptop />} tone="warning" />
        <StatCard label="Pemeliharaan" value={maintenance} icon={<Icon.wrench />} tone="error" />
      </div>
      <div className="grid gap-6 lg:grid-cols-3 mt-6">
        <Card className="lg:col-span-2" padded={false}>
          <div className="flex items-center justify-between p-5 border-b border-[var(--color-border)]">
            <h3 className="font-semibold">Sesi Terbaru</h3>
            <Button variant="ghost" size="sm" onClick={() => go('sessions')}>Kelola sesi</Button>
          </div>
          <SimpleSessionsTable rows={data.sessions.slice(0, 5)} loading={data.loading} />
        </Card>
        <Card>
          <h3 className="font-semibold mb-4">Status Aset</h3>
          {[['Available', available], ['Allocated', allocated], ['Maintenance', maintenance]].map(([label, value]) => (
            <div key={label} className="flex items-center justify-between py-2 text-sm border-b border-[var(--color-border)] last:border-0">
              <span>{label}</span><Badge tone="neutral">{value}</Badge>
            </div>
          ))}
          <Button variant="secondary" className="w-full mt-5" onClick={() => go('assets')}>Kelola Aset</Button>
        </Card>
      </div>
    </>
  )
}

function SimpleSessionsTable({ rows, loading }: { rows: LabSession[]; loading: boolean }) {
  return <Table><THead cols={['Kode', 'Judul', 'Tanggal', 'Waktu', 'Lokasi', 'Status']} />{loading ? <SkeletonRows cols={6} /> : <tbody>{rows.map((s) => <TR key={s.id}><TD><Code>{s.code}</Code></TD><TD className="font-medium">{s.title}</TD><TD>{formatDate(s.session_date)}</TD><TD>{formatTime(s.start_at)}-{formatTime(s.end_at)}</TD><TD>{s.location}</TD><TD><StatusBadge status={sessionStatusLabel(s.status)} /></TD></TR>)}</tbody>}</Table>
}

/* --------------------------------- Users --------------------------------- */
function Users({ data, refresh }: { data: AdminState; refresh: () => Promise<void> }) {
  const toast = useToast()
  const [createOpen, setCreateOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<AuthUser | null>(null)
  const [form, setForm] = useState({ identifier: '', email: '', full_name: '', password: '', role: 'STUDENT' as 'STUDENT' | 'ADMIN' })
  const [editForm, setEditForm] = useState({ identifier: '', email: '', full_name: '', role: 'STUDENT' as string, status: 'ACTIVE' })

  const resetCreateForm = () => setForm({ identifier: '', email: '', full_name: '', password: '', role: 'STUDENT' })

  const submitCreate = async () => {
    try {
      await api.users.create({ ...form, status: 'ACTIVE' })
      toast('Pengguna berhasil dibuat.')
      setCreateOpen(false)
      resetCreateForm()
      await refresh()
    } catch (error) { toast(errMsg(error), 'error') }
  }

  const openEdit = (u: AuthUser) => {
    setEditTarget(u)
    setEditForm({ identifier: u.identifier, email: u.email, full_name: u.full_name, role: u.role, status: u.status })
  }

  const submitEdit = async () => {
    if (!editTarget) return
    try {
      // Update data via PUT
      await api.users.update(editTarget.id, {
        identifier: editForm.identifier,
        email: editForm.email,
        full_name: editForm.full_name,
        role: editForm.role as 'STUDENT' | 'ADMIN',
      })
      // Update status if changed
      if (editForm.status !== editTarget.status) {
        await api.users.updateStatus(editTarget.id, editForm.status)
      }
      toast('Pengguna berhasil diperbarui.')
      setEditTarget(null)
      await refresh()
    } catch (error) { toast(errMsg(error), 'error') }
  }

  return (
    <>
      <PageHeader title="Manajemen Pengguna" subtitle="Data akun dari backend." actions={<Button icon={<Icon.plus width={16} />} onClick={() => { resetCreateForm(); setCreateOpen(true) }}>Tambah Pengguna</Button>} />
      <Card padded={false}>
        <Table>
          <THead cols={['Nama', 'Identifier', 'Email', 'Peran', 'Status', 'Wajah', 'Dibuat', 'Aksi']} />
          {data.loading ? <SkeletonRows cols={8} /> : data.users.length === 0 ? <tbody><tr><td colSpan={8}><EmptyState title="Belum ada pengguna" desc="Tambahkan pengguna admin atau mahasiswa." icon={<Icon.users />} /></td></tr></tbody> : (
            <tbody>{data.users.map((u) => <TR key={u.id}>
              <TD><div className="flex items-center gap-2.5"><Avatar name={u.full_name} size={32} /><span className="font-medium">{u.full_name}</span></div></TD>
              <TD><Code>{u.identifier}</Code></TD><TD>{u.email}</TD><TD>{roleLabel(u.role)}</TD><TD><StatusBadge status={u.status} /></TD>
              <TD>{u.face_enrolled ? <StatusBadge status="Terdaftar" /> : <Badge tone="warning">Belum</Badge>}</TD>
              <TD>{formatDate(u.created_at)}</TD>
              <TD><Button size="sm" variant="secondary" onClick={() => openEdit(u)}>Edit</Button></TD>
            </TR>)}</tbody>
          )}
        </Table>
      </Card>

      {/* Create Modal */}
      <Modal open={createOpen} onClose={() => setCreateOpen(false)} title="Tambah Pengguna" footer={<><Button variant="secondary" onClick={() => setCreateOpen(false)}>Batal</Button><Button onClick={submitCreate}>Simpan</Button></>}>
        <div className="grid sm:grid-cols-2 gap-4">
          <Field label="Nama Lengkap" required><Input value={form.full_name} onChange={(e) => setForm({ ...form, full_name: e.target.value })} /></Field>
          <Field label="Identifier" required><Input value={form.identifier} onChange={(e) => setForm({ ...form, identifier: e.target.value })} /></Field>
          <Field label="Email" required><Input type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} /></Field>
          <Field label="Password" required><Input type="password" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} /></Field>
          <Field label="Role"><Select value={form.role} onChange={(e) => setForm({ ...form, role: e.target.value as 'STUDENT' | 'ADMIN' })}><option value="STUDENT">Mahasiswa</option><option value="ADMIN">Admin</option></Select></Field>
        </div>
      </Modal>

      {/* Edit Modal */}
      <Modal open={!!editTarget} onClose={() => setEditTarget(null)} title="Edit Pengguna" footer={<><Button variant="secondary" onClick={() => setEditTarget(null)}>Batal</Button><Button onClick={submitEdit}>Simpan</Button></>}>
        <div className="grid sm:grid-cols-2 gap-4">
          <Field label="Nama Lengkap"><Input value={editForm.full_name} onChange={(e) => setEditForm({ ...editForm, full_name: e.target.value })} /></Field>
          <Field label="Identifier"><Input value={editForm.identifier} onChange={(e) => setEditForm({ ...editForm, identifier: e.target.value })} /></Field>
          <Field label="Email"><Input type="email" value={editForm.email} onChange={(e) => setEditForm({ ...editForm, email: e.target.value })} /></Field>
          <Field label="Role"><Select value={editForm.role} onChange={(e) => setEditForm({ ...editForm, role: e.target.value })}><option value="STUDENT">Mahasiswa</option><option value="ADMIN">Admin</option></Select></Field>
          <Field label="Status"><Select value={editForm.status} onChange={(e) => setEditForm({ ...editForm, status: e.target.value })}><option value="ACTIVE">ACTIVE</option><option value="INACTIVE">INACTIVE</option></Select></Field>
        </div>
      </Modal>
    </>
  )
}

/* --------------------------------- Sessions --------------------------------- */
function Sessions({ data, refresh }: { data: AdminState; refresh: () => Promise<void> }) {
  const toast = useToast()
  const [createOpen, setCreateOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<LabSession | null>(null)
  const [form, setForm] = useState({ code: '', title: '', course_code: '', location: 'Lab Komputer 1', session_date: '', start: '', end: '', status: 'SCHEDULED' })
  const [editForm, setEditForm] = useState({ code: '', title: '', course_code: '', location: '', session_date: '', start: '', end: '', status: '' })

  const resetForm = () => setForm({ code: '', title: '', course_code: '', location: 'Lab Komputer 1', session_date: '', start: '', end: '', status: 'SCHEDULED' })

  const submitCreate = async () => {
    try {
      await api.sessions.create({ code: form.code, title: form.title, course_code: form.course_code, location: form.location, session_date: form.session_date, start_at: `${form.session_date}T${form.start}:00+07:00`, end_at: `${form.session_date}T${form.end}:00+07:00`, status: form.status })
      toast('Sesi berhasil dibuat.')
      setCreateOpen(false)
      resetForm()
      await refresh()
    } catch (error) { toast(errMsg(error), 'error') }
  }

  const openEdit = (s: LabSession) => {
    setEditTarget(s)
    setEditForm({
      code: s.code,
      title: s.title,
      course_code: s.course_code || '',
      location: s.location,
      session_date: extractDate(s.session_date),
      start: extractTime(s.start_at),
      end: extractTime(s.end_at),
      status: s.status,
    })
  }

  const submitEdit = async () => {
    if (!editTarget) return
    try {
      // Update fields via PUT
      await api.sessions.update(editTarget.id, {
        code: editForm.code,
        title: editForm.title,
        course_code: editForm.course_code,
        location: editForm.location,
        session_date: editForm.session_date,
        start_at: `${editForm.session_date}T${editForm.start}:00+07:00`,
        end_at: `${editForm.session_date}T${editForm.end}:00+07:00`,
      })
      // Update status if changed
      if (editForm.status !== editTarget.status) {
        await api.sessions.updateStatus(editTarget.id, editForm.status)
      }
      toast('Sesi berhasil diperbarui.')
      setEditTarget(null)
      await refresh()
    } catch (error) { toast(errMsg(error), 'error') }
  }

  return <>
    <PageHeader title="Sesi Praktikum" subtitle="Kelola jadwal dan status sesi." actions={<Button icon={<Icon.plus width={16} />} onClick={() => { resetForm(); setCreateOpen(true) }}>Tambah Sesi</Button>} />
    <Card padded={false}>
      <Table>
        <THead cols={['Kode', 'Judul', 'Tanggal', 'Waktu', 'Lokasi', 'Status', 'Aksi']} />
        {data.loading ? <SkeletonRows cols={7} /> : <tbody>{data.sessions.map((s) => <TR key={s.id}>
          <TD><Code>{s.code}</Code></TD><TD className="font-medium">{s.title}</TD><TD>{formatDate(s.session_date)}</TD><TD>{formatTime(s.start_at)}-{formatTime(s.end_at)}</TD><TD>{s.location}</TD><TD><StatusBadge status={sessionStatusLabel(s.status)} /></TD>
          <TD><Button size="sm" variant="secondary" onClick={() => openEdit(s)}>Edit</Button></TD>
        </TR>)}</tbody>}
      </Table>
    </Card>

    {/* Create Modal */}
    <Modal open={createOpen} onClose={() => setCreateOpen(false)} title="Tambah Sesi" size="lg" footer={<><Button variant="secondary" onClick={() => setCreateOpen(false)}>Batal</Button><Button onClick={submitCreate}>Simpan</Button></>}>
      <div className="grid sm:grid-cols-2 gap-4">
        <Field label="Kode" required><Input value={form.code} onChange={(e) => setForm({ ...form, code: e.target.value })} /></Field>
        <Field label="Judul" required><Input value={form.title} onChange={(e) => setForm({ ...form, title: e.target.value })} /></Field>
        <Field label="Kode Mata Kuliah"><Input value={form.course_code} onChange={(e) => setForm({ ...form, course_code: e.target.value })} /></Field>
        <Field label="Lokasi"><Input value={form.location} onChange={(e) => setForm({ ...form, location: e.target.value })} /></Field>
        <Field label="Tanggal" required><Input type="date" value={form.session_date} onChange={(e) => setForm({ ...form, session_date: e.target.value })} /></Field>
        <div className="grid grid-cols-2 gap-3"><Field label="Mulai"><Input type="time" value={form.start} onChange={(e) => setForm({ ...form, start: e.target.value })} /></Field><Field label="Selesai"><Input type="time" value={form.end} onChange={(e) => setForm({ ...form, end: e.target.value })} /></Field></div>
        <Field label="Status"><Select value={form.status} onChange={(e) => setForm({ ...form, status: e.target.value })}><option value="SCHEDULED">SCHEDULED</option><option value="ACTIVE">ACTIVE</option><option value="CLOSED">CLOSED</option><option value="CANCELLED">CANCELLED</option></Select></Field>
      </div>
    </Modal>

    {/* Edit Modal */}
    <Modal open={!!editTarget} onClose={() => setEditTarget(null)} title="Edit Sesi" size="lg" footer={<><Button variant="secondary" onClick={() => setEditTarget(null)}>Batal</Button><Button onClick={submitEdit}>Simpan</Button></>}>
      <div className="grid sm:grid-cols-2 gap-4">
        <Field label="Kode"><Input value={editForm.code} onChange={(e) => setEditForm({ ...editForm, code: e.target.value })} /></Field>
        <Field label="Judul"><Input value={editForm.title} onChange={(e) => setEditForm({ ...editForm, title: e.target.value })} /></Field>
        <Field label="Kode Mata Kuliah"><Input value={editForm.course_code} onChange={(e) => setEditForm({ ...editForm, course_code: e.target.value })} /></Field>
        <Field label="Lokasi"><Input value={editForm.location} onChange={(e) => setEditForm({ ...editForm, location: e.target.value })} /></Field>
        <Field label="Tanggal"><Input type="date" value={editForm.session_date} onChange={(e) => setEditForm({ ...editForm, session_date: e.target.value })} /></Field>
        <div className="grid grid-cols-2 gap-3"><Field label="Mulai"><Input type="time" value={editForm.start} onChange={(e) => setEditForm({ ...editForm, start: e.target.value })} /></Field><Field label="Selesai"><Input type="time" value={editForm.end} onChange={(e) => setEditForm({ ...editForm, end: e.target.value })} /></Field></div>
        <Field label="Status"><Select value={editForm.status} onChange={(e) => setEditForm({ ...editForm, status: e.target.value })}><option value="SCHEDULED">SCHEDULED</option><option value="ACTIVE">ACTIVE</option><option value="CLOSED">CLOSED</option><option value="CANCELLED">CANCELLED</option></Select></Field>
      </div>
    </Modal>
  </>
}

/* --------------------------------- Attendance --------------------------------- */
function AttendanceAdmin({ data }: { data: AdminState }) {
  const userMap = new Map(data.users.map((u) => [u.id, u]))
  return <>
    <PageHeader title="Pemantauan Presensi" subtitle="Data presensi seluruh mahasiswa dari backend." />
    <Card padded={false}>
      <Table>
        <THead cols={['Sesi', 'Mahasiswa', 'NIM', 'Check-in', 'Check-out', 'Verifikasi']} />
        {data.loading ? <SkeletonRows cols={6} /> : <tbody>{data.attendance.map((a) => {
          const u = userMap.get(a.user_id || '')
          return <TR key={a.id}>
            <TD className="font-medium">{a.session?.title || a.session_id}</TD>
            <TD>{u?.full_name || a.user_id || '—'}</TD>
            <TD><Code>{u?.identifier || '—'}</Code></TD>
            <TD>{formatDateTime(a.checked_in_at)}</TD>
            <TD>{formatDateTime(a.checked_out_at)}</TD>
            <TD><StatusBadge status={a.verification_status || '—'} /></TD>
          </TR>
        })}</tbody>}
      </Table>
      <Pagination page={1} pages={1} total={data.attendance.length} onPage={() => {}} />
    </Card>
  </>
}

/* --------------------------------- Assets --------------------------------- */
function Assets({ data, refresh }: { data: AdminState; refresh: () => Promise<void> }) {
  const toast = useToast()
  const [createOpen, setCreateOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<AssetRecord | null>(null)
  const [form, setForm] = useState({ asset_code: '', asset_name: '', category: 'Laptop', brand: '', model: '', serial_number: '', location: 'Lab Komputer 1', status: 'AVAILABLE' })
  const [editForm, setEditForm] = useState({ asset_code: '', asset_name: '', category: '', brand: '', model: '', serial_number: '', location: '', status: '' })

  const resetForm = () => setForm({ asset_code: '', asset_name: '', category: 'Laptop', brand: '', model: '', serial_number: '', location: 'Lab Komputer 1', status: 'AVAILABLE' })

  const submitCreate = async () => {
    try { await api.assets.create(form); toast('Aset berhasil dibuat.'); setCreateOpen(false); resetForm(); await refresh() }
    catch (error) { toast(errMsg(error), 'error') }
  }

  const openEdit = (a: AssetRecord) => {
    setEditTarget(a)
    setEditForm({
      asset_code: a.asset_code,
      asset_name: a.asset_name,
      category: a.category || '',
      brand: a.brand || '',
      model: a.model || '',
      serial_number: a.serial_number || '',
      location: a.location || '',
      status: a.status || 'AVAILABLE',
    })
  }

  const submitEdit = async () => {
    if (!editTarget) return
    try {
      await api.assets.update(editTarget.id, {
        asset_code: editForm.asset_code,
        asset_name: editForm.asset_name,
        category: editForm.category,
        brand: editForm.brand,
        model: editForm.model,
        serial_number: editForm.serial_number,
        location: editForm.location,
      })
      if (editForm.status !== editTarget.status) {
        await api.assets.updateStatus(editTarget.id, editForm.status)
      }
      toast('Aset berhasil diperbarui.')
      setEditTarget(null)
      await refresh()
    } catch (error) { toast(errMsg(error), 'error') }
  }

  return <>
    <PageHeader title="Master Aset" subtitle="Inventaris aset dari backend." actions={<Button icon={<Icon.plus width={16} />} onClick={() => { resetForm(); setCreateOpen(true) }}>Tambah Aset</Button>} />
    <Card padded={false}>
      <Table>
        <THead cols={['Kode', 'Nama', 'Kategori', 'Merek/Model', 'Serial', 'Lokasi', 'Status', 'Aksi']} />
        {data.loading ? <SkeletonRows cols={8} /> : <tbody>{data.assets.map((a) => <TR key={a.id}>
          <TD><Code>{a.asset_code}</Code></TD><TD className="font-medium">{a.asset_name}</TD><TD>{a.category}</TD><TD>{a.brand} {a.model}</TD><TD>{a.serial_number || '—'}</TD><TD>{a.location}</TD><TD><StatusBadge status={assetStatusLabel(a.status)} /></TD>
          <TD><Button size="sm" variant="secondary" onClick={() => openEdit(a)}>Edit</Button></TD>
        </TR>)}</tbody>}
      </Table>
    </Card>

    {/* Create Modal */}
    <Modal open={createOpen} onClose={() => setCreateOpen(false)} title="Tambah Aset" footer={<><Button variant="secondary" onClick={() => setCreateOpen(false)}>Batal</Button><Button onClick={submitCreate}>Simpan</Button></>}>
      <div className="grid sm:grid-cols-2 gap-4">
        <Field label="Kode" required><Input value={form.asset_code} onChange={(e) => setForm({ ...form, asset_code: e.target.value })} /></Field>
        <Field label="Nama" required><Input value={form.asset_name} onChange={(e) => setForm({ ...form, asset_name: e.target.value })} /></Field>
        <Field label="Kategori"><Input value={form.category} onChange={(e) => setForm({ ...form, category: e.target.value })} /></Field>
        <Field label="Brand"><Input value={form.brand} onChange={(e) => setForm({ ...form, brand: e.target.value })} /></Field>
        <Field label="Model"><Input value={form.model} onChange={(e) => setForm({ ...form, model: e.target.value })} /></Field>
        <Field label="Serial"><Input value={form.serial_number} onChange={(e) => setForm({ ...form, serial_number: e.target.value })} /></Field>
        <Field label="Lokasi"><Input value={form.location} onChange={(e) => setForm({ ...form, location: e.target.value })} /></Field>
      </div>
    </Modal>

    {/* Edit Modal */}
    <Modal open={!!editTarget} onClose={() => setEditTarget(null)} title="Edit Aset" footer={<><Button variant="secondary" onClick={() => setEditTarget(null)}>Batal</Button><Button onClick={submitEdit}>Simpan</Button></>}>
      <div className="grid sm:grid-cols-2 gap-4">
        <Field label="Kode"><Input value={editForm.asset_code} onChange={(e) => setEditForm({ ...editForm, asset_code: e.target.value })} /></Field>
        <Field label="Nama"><Input value={editForm.asset_name} onChange={(e) => setEditForm({ ...editForm, asset_name: e.target.value })} /></Field>
        <Field label="Kategori"><Input value={editForm.category} onChange={(e) => setEditForm({ ...editForm, category: e.target.value })} /></Field>
        <Field label="Brand"><Input value={editForm.brand} onChange={(e) => setEditForm({ ...editForm, brand: e.target.value })} /></Field>
        <Field label="Model"><Input value={editForm.model} onChange={(e) => setEditForm({ ...editForm, model: e.target.value })} /></Field>
        <Field label="Serial"><Input value={editForm.serial_number} onChange={(e) => setEditForm({ ...editForm, serial_number: e.target.value })} /></Field>
        <Field label="Lokasi"><Input value={editForm.location} onChange={(e) => setEditForm({ ...editForm, location: e.target.value })} /></Field>
        <Field label="Status"><Select value={editForm.status} onChange={(e) => setEditForm({ ...editForm, status: e.target.value })}><option value="AVAILABLE">AVAILABLE</option><option value="ALLOCATED">ALLOCATED</option><option value="MAINTENANCE">MAINTENANCE</option><option value="RETIRED">RETIRED</option></Select></Field>
      </div>
    </Modal>
  </>
}

/* --------------------------------- Allocations --------------------------------- */
function Allocations({ data, refresh }: { data: AdminState; refresh: () => Promise<void> }) {
  const toast = useToast()
  const [open, setOpen] = useState(false)
  const [form, setForm] = useState({ user_id: '', session_id: '', asset_id: '' })
  const students = data.users.filter((u) => u.role === 'STUDENT')
  const activeSessions = data.sessions.filter((s) => s.status === 'ACTIVE')
  const availableAssets = data.assets.filter((a) => a.status === 'AVAILABLE')
  const userMap = new Map(data.users.map((u) => [u.id, u]))

  const submit = async () => { try { await api.allocations.create(form); toast('Aset berhasil dialokasikan.'); setOpen(false); setForm({ user_id: '', session_id: '', asset_id: '' }); await refresh() } catch (error) { toast(errMsg(error), 'error') } }
  const returnAsset = async (id: string) => { try { await api.allocations.returnAsset(id, 'GOOD'); toast('Aset berhasil dikembalikan.'); await refresh() } catch (error) { toast(errMsg(error), 'error') } }

  return <>
    <PageHeader title="Alokasi Aset" subtitle="Alokasi aset memakai ID backend asli." actions={<Button icon={<Icon.plus width={16} />} onClick={() => setOpen(true)}>Alokasi Baru</Button>} />
    <Card padded={false}>
      <Table>
        <THead cols={['Mahasiswa', 'Sesi', 'Aset', 'Waktu', 'Status', 'Aksi']} />
        {data.loading ? <SkeletonRows cols={6} /> : <tbody>{data.allocations.map((a) => {
          const u = userMap.get(a.user_id || '')
          return <TR key={a.id}>
            <TD>{u?.full_name || a.user_id || '—'}</TD>
            <TD>{a.session?.title || a.session_id}</TD>
            <TD><Code>{a.asset?.asset_code || a.asset_id || '—'}</Code></TD>
            <TD>{formatDateTime(a.allocated_at)}</TD>
            <TD><StatusBadge status={a.returned_at ? 'Selesai' : 'Aktif'} /></TD>
            <TD>{!a.returned_at && <Button size="sm" variant="secondary" onClick={() => returnAsset(a.id)}>Kembalikan</Button>}</TD>
          </TR>
        })}</tbody>}
      </Table>
    </Card>
    <Modal open={open} onClose={() => setOpen(false)} title="Alokasi Baru" footer={<><Button variant="secondary" onClick={() => setOpen(false)}>Batal</Button><Button onClick={submit} disabled={!form.user_id || !form.session_id || !form.asset_id}>Alokasikan</Button></>}>
      <div className="grid gap-4">
        <Field label="Mahasiswa"><Select value={form.user_id} onChange={(e) => setForm({ ...form, user_id: e.target.value })}><option value="">Pilih mahasiswa</option>{students.map((u) => <option key={u.id} value={u.id}>{u.full_name} - {u.identifier}</option>)}</Select></Field>
        <Field label="Sesi Aktif"><Select value={form.session_id} onChange={(e) => setForm({ ...form, session_id: e.target.value })}><option value="">Pilih sesi</option>{activeSessions.map((s) => <option key={s.id} value={s.id}>{s.title}</option>)}</Select></Field>
        <Field label="Aset Tersedia"><Select value={form.asset_id} onChange={(e) => setForm({ ...form, asset_id: e.target.value })}><option value="">Pilih aset</option>{availableAssets.map((a) => <option key={a.id} value={a.id}>{a.asset_code} - {a.asset_name}</option>)}</Select></Field>
      </div>
    </Modal>
  </>
}

/* --------------------------------- Maintenance --------------------------------- */
function Maintenance({ data, refresh }: { data: AdminState; refresh: () => Promise<void> }) {
  const toast = useToast()
  const [createOpen, setCreateOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<MaintenanceRecord | null>(null)
  const [form, setForm] = useState({ asset_id: '', issue_type: 'HARDWARE', description: '', priority: 'MEDIUM' })
  const [editForm, setEditForm] = useState({ description: '', priority: '', status: '', resolution: '' })

  const resetForm = () => setForm({ asset_id: '', issue_type: 'HARDWARE', description: '', priority: 'MEDIUM' })

  const submitCreate = async () => {
    try { await api.maintenance.create(form); toast('Tiket pemeliharaan dibuat.'); setCreateOpen(false); resetForm(); await refresh() }
    catch (error) { toast(errMsg(error), 'error') }
  }

  const openEdit = (m: MaintenanceRecord) => {
    setEditTarget(m)
    setEditForm({
      description: m.description || '',
      priority: m.priority || 'MEDIUM',
      status: m.status || 'OPEN',
      resolution: m.resolution || '',
    })
  }

  const submitEdit = async () => {
    if (!editTarget) return
    try {
      const body: Record<string, string> = {
        description: editForm.description,
        priority: editForm.priority,
      }
      // Status transition
      if (editForm.status !== editTarget.status) {
        body.status = editForm.status
      }
      // Resolution wajib saat COMPLETED
      if (editForm.resolution) {
        body.resolution = editForm.resolution
      }
      await api.maintenance.update(editTarget.id, body)
      toast('Tiket pemeliharaan diperbarui.')
      setEditTarget(null)
      await refresh()
    } catch (error) { toast(errMsg(error), 'error') }
  }

  // Valid next statuses per state machine
  const getNextStatuses = (current?: string) => {
    switch (current) {
      case 'OPEN': return ['OPEN', 'IN_PROGRESS', 'CANCELLED']
      case 'IN_PROGRESS': return ['IN_PROGRESS', 'COMPLETED', 'CANCELLED']
      case 'COMPLETED': return ['COMPLETED']
      case 'CANCELLED': return ['CANCELLED']
      default: return ['OPEN']
    }
  }

  return <>
    <PageHeader title="Pemeliharaan" subtitle="Tiket perbaikan aset." actions={<Button icon={<Icon.plus width={16} />} onClick={() => { resetForm(); setCreateOpen(true) }}>Buat Tiket</Button>} />
    <Card padded={false}>
      <Table>
        <THead cols={['ID', 'Aset', 'Keluhan', 'Prioritas', 'Status', 'Dibuat', 'Aksi']} />
        {data.loading ? <SkeletonRows cols={7} /> : <tbody>{data.maintenance.map((m) => <TR key={m.id}>
          <TD><Code>{m.id.slice(0, 8)}</Code></TD>
          <TD><Code>{m.asset?.asset_code || m.asset_id || '—'}</Code></TD>
          <TD>{m.description || m.issue_type}</TD>
          <TD><StatusBadge status={m.priority || 'MEDIUM'} /></TD>
          <TD><StatusBadge status={m.status || 'OPEN'} /></TD>
          <TD>{formatDate(m.created_at)}</TD>
          <TD>
            {(m.status !== 'COMPLETED' && m.status !== 'CANCELLED') && (
              <Button size="sm" variant="secondary" onClick={() => openEdit(m)}>Edit</Button>
            )}
          </TD>
        </TR>)}</tbody>}
      </Table>
    </Card>

    {/* Create Modal */}
    <Modal open={createOpen} onClose={() => setCreateOpen(false)} title="Buat Tiket Pemeliharaan" footer={<><Button variant="secondary" onClick={() => setCreateOpen(false)}>Batal</Button><Button onClick={submitCreate} disabled={!form.asset_id || !form.description}>Simpan</Button></>}>
      <div className="grid gap-4">
        <Field label="Aset"><Select value={form.asset_id} onChange={(e) => setForm({ ...form, asset_id: e.target.value })}><option value="">Pilih aset</option>{data.assets.map((a) => <option key={a.id} value={a.id}>{a.asset_code} - {a.asset_name}</option>)}</Select></Field>
        <Field label="Jenis"><Select value={form.issue_type} onChange={(e) => setForm({ ...form, issue_type: e.target.value })}><option value="HARDWARE">Hardware</option><option value="SOFTWARE">Software</option><option value="NETWORK">Jaringan</option></Select></Field>
        <Field label="Prioritas"><Select value={form.priority} onChange={(e) => setForm({ ...form, priority: e.target.value })}><option value="LOW">LOW</option><option value="MEDIUM">MEDIUM</option><option value="HIGH">HIGH</option></Select></Field>
        <Field label="Deskripsi" required><Input value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} /></Field>
      </div>
    </Modal>

    {/* Edit Modal */}
    <Modal open={!!editTarget} onClose={() => setEditTarget(null)} title="Edit Tiket Pemeliharaan" footer={<><Button variant="secondary" onClick={() => setEditTarget(null)}>Batal</Button><Button onClick={submitEdit}>Simpan</Button></>}>
      <div className="grid gap-4">
        <Field label="Deskripsi"><Input value={editForm.description} onChange={(e) => setEditForm({ ...editForm, description: e.target.value })} /></Field>
        <Field label="Prioritas"><Select value={editForm.priority} onChange={(e) => setEditForm({ ...editForm, priority: e.target.value })}><option value="LOW">LOW</option><option value="MEDIUM">MEDIUM</option><option value="HIGH">HIGH</option></Select></Field>
        <Field label="Status">
          <Select value={editForm.status} onChange={(e) => setEditForm({ ...editForm, status: e.target.value })}>
            {getNextStatuses(editTarget?.status).map((s) => <option key={s} value={s}>{s}</option>)}
          </Select>
        </Field>
        {(editForm.status === 'COMPLETED' || editForm.status === 'IN_PROGRESS') && (
          <Field label="Resolusi" required={editForm.status === 'COMPLETED'}>
            <Input value={editForm.resolution} onChange={(e) => setEditForm({ ...editForm, resolution: e.target.value })} placeholder="Catatan perbaikan yang dilakukan" />
          </Field>
        )}
        {editForm.status === 'COMPLETED' && !editForm.resolution && (
          <Alert tone="warning" title="Resolusi wajib">Isi resolusi perbaikan sebelum menyelesaikan tiket.</Alert>
        )}
      </div>
    </Modal>
  </>
}

/* --------------------------------- Reports --------------------------------- */
function Reports({ data }: { data: AdminState }) {
  const [tab, setTab] = useState('attendance')
  return <>
    <PageHeader title="Laporan" subtitle="Rekap dan export CSV dari endpoint backend." />
    <Tabs tabs={['attendance', 'assets', 'allocations', 'maintenance']} active={tab} onChange={setTab} />
    <Card>
      <div className="grid gap-4 sm:grid-cols-4">
        <StatCard label="Users" value={data.users.length} icon={<Icon.users />} tone="primary" />
        <StatCard label="Sessions" value={data.sessions.length} icon={<Icon.calendar />} tone="info" />
        <StatCard label="Assets" value={data.assets.length} icon={<Icon.box />} tone="success" />
        <StatCard label="Maintenance" value={data.maintenance.length} icon={<Icon.wrench />} tone="warning" />
      </div>
      <Button className="mt-5" icon={<Icon.download width={16} />} onClick={() => window.open(api.reports.csvUrl(tab), '_blank')}>Export CSV {tab}</Button>
    </Card>
  </>
}
