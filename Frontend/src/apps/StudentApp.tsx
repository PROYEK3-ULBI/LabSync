import { useState } from 'react'
import { Icon } from '../components/icons'
import { FaceFlow } from '../components/FaceFlow'
import { Shell, type NavItem } from '../components/Shell'
import {
  activeSession,
  assetHistory,
  attendanceHistory,
  currentStudent,
  myAsset,
  sessions,
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
  Input,
  PageHeader,
  Pagination,
  Select,
  StatCard,
  StatusBadge,
  SkeletonRows,
  StepIndicator,
  Table,
  TD,
  THead,
  TR,
  Tabs,
  Timeline,
  useSimulatedLoading,
  useToast,
  Modal,
} from '../components/ui'
import type { AttendanceRow } from '../lib/data'

const NAV: NavItem[] = [
  { key: 'dashboard', label: 'Dasbor', icon: 'dashboard' },
  { key: 'sessions', label: 'Sesi Praktikum', icon: 'calendar' },
  { key: 'attendance', label: 'Presensi', icon: 'checkCircle' },
  { key: 'assets', label: 'Aset Saya', icon: 'laptop' },
  { key: 'history', label: 'Riwayat', icon: 'history' },
  { key: 'profile', label: 'Profil', icon: 'user' },
]

export function StudentApp({ onLogout }: { onLogout: () => void }) {
  const [route, setRoute] = useState('dashboard')
  const [detail, setDetail] = useState<AttendanceRow | null>(null)
  const [enrolled, setEnrolled] = useState(currentStudent.faceEnrolled)
  const openDetail = (r: AttendanceRow) => {
    setDetail(r)
    setRoute('attendanceDetail')
  }

  return (
    <Shell
      brand="LabPresensi"
      roleLabel="Mahasiswa"
      roleTone="primary"
      nav={NAV}
      active={route}
      onNavigate={setRoute}
      user={{ name: currentStudent.name, sub: currentStudent.nim }}
      onLogout={onLogout}
    >
      {route === 'dashboard' && <Dashboard go={setRoute} enrolled={enrolled} />}
      {route === 'sessions' && <Sessions go={setRoute} />}
      {route === 'attendance' && <Attendance go={setRoute} enrolled={enrolled} />}
      {route === 'assets' && <MyAssets />}
      {route === 'history' && <History onOpen={openDetail} />}
      {route === 'attendanceDetail' && <AttendanceDetail row={detail} go={setRoute} />}
      {route === 'checkout' && <Checkout go={setRoute} />}
      {route === 'profile' && <Profile go={setRoute} enrolled={enrolled} onReset={() => setEnrolled(false)} />}
      {route === 'enroll' && <Enroll go={setRoute} onDone={() => setEnrolled(true)} />}
    </Shell>
  )
}

/* --------------------------------- Dashboard --------------------------------- */
function Dashboard({ go, enrolled }: { go: (r: string) => void; enrolled: boolean }) {
  return (
    <>
      <PageHeader
        title={`Halo, ${currentStudent.name.split(' ')[0]} 👋`}
        subtitle="Berikut ringkasan praktikum dan presensi Anda hari ini."
      />

      {!enrolled && (
        <div className="mb-6">
          <Alert tone="warning" title="Wajah belum terdaftar">
            <div className="flex flex-wrap items-center gap-3">
              <span>Daftarkan wajah Anda terlebih dahulu agar dapat melakukan presensi.</span>
              <Button size="sm" icon={<Icon.camera width={15} />} onClick={() => go('enroll')}>Daftar Wajah Sekarang</Button>
            </div>
          </Alert>
        </div>
      )}

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <StatCard label="Sesi Hari Ini" value="1 Sesi" hint="Pemrograman Web · 13.00" icon={<Icon.calendar />} tone="primary" />
        <StatCard label="Status Presensi" value="Belum Hadir" hint="Sesi aktif menunggu presensi" icon={<Icon.clock />} tone="warning" />
        <StatCard label="Aset Aktif" value="LAB-LT-007" hint="Dell Latitude 5490" icon={<Icon.laptop />} tone="info" />
        <StatCard
          label="Status Verifikasi"
          value={enrolled ? 'Terdaftar' : 'Belum Terdaftar'}
          hint={enrolled ? 'Templat wajah aktif' : 'Perlu pendaftaran wajah'}
          icon={<Icon.shield />}
          tone={enrolled ? 'success' : 'warning'}
        />
      </div>

      <div className="grid gap-6 lg:grid-cols-3 mt-6">
        {/* Active session */}
        <Card className="lg:col-span-2 border-[var(--color-primary)]/30 ring-1 ring-[var(--color-primary)]/10">
          <div className="flex items-start justify-between gap-3">
            <div>
              <Badge tone="success">Sesi Aktif</Badge>
              <h3 className="mt-2 text-lg font-semibold tracking-tight">{activeSession.name}</h3>
              <p className="text-sm text-[var(--color-text-2)]">{activeSession.course}</p>
            </div>
            <Code>{activeSession.code}</Code>
          </div>

          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-5">
            <Info icon={<Icon.calendar width={16} />} label="Tanggal" value={activeSession.date} />
            <Info icon={<Icon.clock width={16} />} label="Waktu" value={`${activeSession.start}–${activeSession.end}`} />
            <Info icon={<Icon.box width={16} />} label="Ruang" value={activeSession.room} />
            <Info icon={<Icon.users width={16} />} label="Kehadiran" value={`${activeSession.present}/${activeSession.total}`} />
          </div>

          <div className="mt-6 flex flex-wrap gap-2">
            <Button size="lg" icon={<Icon.camera width={18} height={18} />} onClick={() => go('attendance')}>
              Mulai Presensi
            </Button>
            <Button size="lg" variant="secondary" onClick={() => go('sessions')}>
              Lihat Semua Sesi
            </Button>
            <Button size="lg" variant="ghost" icon={<Icon.logout width={18} height={18} />} onClick={() => go('checkout')}>
              Check-out Sesi
            </Button>
          </div>
        </Card>

        {/* Current asset */}
        <Card>
          <div className="flex items-center justify-between">
            <h3 className="font-semibold">Alokasi Aset Saat Ini</h3>
            <StatusBadge status={myAsset.status} />
          </div>
          <div className="mt-4 flex items-center gap-3 p-3 rounded-lg bg-[var(--color-surface-2)] border border-[var(--color-border)]">
            <span className="inline-flex items-center justify-center w-11 h-11 rounded-lg bg-white border border-[var(--color-border)] text-[var(--color-primary)]">
              <Icon.laptop />
            </span>
            <div className="min-w-0">
              <p className="font-medium text-sm truncate">{myAsset.name}</p>
              <p className="text-xs text-[var(--color-text-2)] font-mono">{myAsset.code}</p>
            </div>
          </div>
          <dl className="mt-4 space-y-2.5 text-sm">
            <Row k="Sesi" v={myAsset.session} />
            <Row k="Dialokasikan" v={myAsset.allocatedAt} />
            <Row k="Kondisi" v={<Badge tone="success">{myAsset.condition}</Badge>} />
          </dl>
          <Button variant="secondary" className="w-full mt-4" onClick={() => go('assets')}>
            Detail Aset
          </Button>
        </Card>
      </div>
    </>
  )
}

function Info({ icon, label, value }: { icon: React.ReactNode; label: string; value: string }) {
  return (
    <div>
      <div className="flex items-center gap-1.5 text-[var(--color-text-3)] text-xs">{icon}{label}</div>
      <p className="mt-1 text-sm font-medium text-[var(--color-text)]">{value}</p>
    </div>
  )
}
function Row({ k, v }: { k: string; v: React.ReactNode }) {
  return (
    <div className="flex items-center justify-between gap-3">
      <dt className="text-[var(--color-text-2)]">{k}</dt>
      <dd className="font-medium text-right">{v}</dd>
    </div>
  )
}

/* --------------------------------- Sessions --------------------------------- */
function Sessions({ go }: { go: (r: string) => void }) {
  const [status, setStatus] = useState('Semua')
  const filtered = sessions.filter((s) => status === 'Semua' || s.status === status)
  return (
    <>
      <PageHeader
        title="Sesi Praktikum"
        subtitle="Jadwal dan status sesi praktikum Anda."
        actions={
          <Select value={status} onChange={(e) => setStatus(e.target.value)} className="w-44">
            {['Semua', 'Aktif', 'Akan Datang', 'Ditutup', 'Dibatalkan'].map((s) => (
              <option key={s}>{s}</option>
            ))}
          </Select>
        }
      />
      {filtered.length === 0 ? (
        <Card>
          <EmptyState title="Tidak ada sesi" desc="Tidak ada sesi praktikum untuk filter yang dipilih." icon={<Icon.calendar />} />
        </Card>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {filtered.map((s) => {
            const active = s.status === 'Aktif'
            return (
              <Card key={s.code} className={active ? 'ring-1 ring-[var(--color-primary)]/20 border-[var(--color-primary)]/30' : ''}>
                <div className="flex items-start justify-between gap-2">
                  <StatusBadge status={s.status} />
                  <Code>{s.code}</Code>
                </div>
                <h3 className="mt-3 font-semibold">{s.name}</h3>
                <p className="text-sm text-[var(--color-text-2)]">{s.course}</p>
                <div className="mt-4 space-y-2 text-sm text-[var(--color-text-2)]">
                  <p className="flex items-center gap-2"><Icon.calendar width={15} /> {s.date} · {s.start}–{s.end}</p>
                  <p className="flex items-center gap-2"><Icon.box width={15} /> {s.room}</p>
                  <p className="flex items-center gap-2">
                    <Icon.checkCircle width={15} /> Presensi: <StatusBadge status={s.attendance ?? 'Belum Hadir'} />
                  </p>
                </div>
                <div className="mt-4">
                  {active ? (
                    <Button className="w-full" onClick={() => go('attendance')} icon={<Icon.camera width={16} />}>
                      Mulai Presensi
                    </Button>
                  ) : (
                    <Button variant="secondary" className="w-full" disabled={s.status === 'Dibatalkan'}>
                      {s.status === 'Ditutup' ? 'Lihat Detail' : s.status === 'Dibatalkan' ? 'Dibatalkan' : 'Belum Dibuka'}
                    </Button>
                  )}
                </div>
              </Card>
            )
          })}
        </div>
      )}
    </>
  )
}

/* --------------------------------- Attendance / verify flow --------------------------------- */
function Attendance({ go, enrolled }: { go: (r: string) => void; enrolled: boolean }) {
  const [step, setStep] = useState<'verify' | 'success'>('verify')

  if (!enrolled) {
    return (
      <>
        <PageHeader title="Presensi — Verifikasi Wajah" subtitle="Verifikasi identitas Anda untuk mencatat kehadiran." />
        <div className="max-w-xl">
          <Alert tone="warning" title="Wajah belum terdaftar">
            Anda perlu mendaftarkan wajah terlebih dahulu sebelum dapat melakukan verifikasi presensi.
          </Alert>
          <div className="mt-4 flex gap-2">
            <Button icon={<Icon.camera width={16} />} onClick={() => go('enroll')}>Daftar Wajah</Button>
            <Button variant="secondary" onClick={() => go('dashboard')}>Kembali</Button>
          </div>
        </div>
      </>
    )
  }

  if (step === 'success') {
    return (
      <div className="max-w-xl mx-auto">
        <Card className="text-center animate-fade">
          <span className="mx-auto inline-flex items-center justify-center w-16 h-16 rounded-full bg-[var(--color-success-soft)] text-[var(--color-success)]">
            <Icon.checkCircle width={34} height={34} />
          </span>
          <h2 className="mt-4 text-xl font-semibold tracking-tight">Presensi berhasil</h2>
          <p className="text-sm text-[var(--color-text-2)] mt-1">
            Kehadiran Anda telah tercatat untuk sesi ini.
          </p>
          <div className="mt-5 text-left rounded-lg border border-[var(--color-border)] divide-y divide-[var(--color-border)]">
            <Line k="Mahasiswa" v={currentStudent.name} />
            <Line k="Sesi" v={activeSession.name} />
            <Line k="Waktu Check-in" v="18 Sep 2026, 13.04" />
            <Line k="Ruang" v={activeSession.room} />
            <Line k="Status" v={<Badge tone="success">Hadir</Badge>} />
          </div>
          <div className="mt-5 flex flex-col sm:flex-row gap-2 justify-center">
            <Button onClick={() => go('history')} icon={<Icon.history width={16} />}>Lihat Detail Presensi</Button>
            <Button variant="secondary" onClick={() => go('assets')}>Lihat Alokasi Aset</Button>
          </div>
        </Card>
      </div>
    )
  }

  return (
    <>
      <PageHeader
        title="Presensi — Verifikasi Wajah"
        subtitle="Verifikasi identitas Anda untuk mencatat kehadiran pada sesi aktif."
        breadcrumb={['Presensi', 'Verifikasi']}
      />
      <div className="mb-5">
        <Alert tone="info" title="Anda memverifikasi identitas untuk presensi">
          Sistem membandingkan wajah di depan kamera dengan templat wajah akun <b>{currentStudent.name}</b> yang sedang login.
        </Alert>
      </div>
      <FaceFlow
        mode="verify"
        studentName={currentStudent.name}
        sessionName={activeSession.name}
        onComplete={() => setStep('success')}
        onCancel={() => go('dashboard')}
      />
    </>
  )
}
function Line({ k, v }: { k: string; v: React.ReactNode }) {
  return (
    <div className="flex items-center justify-between gap-3 px-4 py-2.5 text-sm">
      <span className="text-[var(--color-text-2)]">{k}</span>
      <span className="font-medium text-right">{v}</span>
    </div>
  )
}

/* --------------------------------- My assets --------------------------------- */
function MyAssets() {
  const [tab, setTab] = useState('Alokasi Aktif')
  return (
    <>
      <PageHeader title="Aset Saya" subtitle="Aset yang sedang dan pernah dialokasikan kepada Anda." />
      <Tabs tabs={['Alokasi Aktif', 'Riwayat Aset']} active={tab} onChange={setTab} />

      {tab === 'Alokasi Aktif' ? (
        <Card>
          <div className="flex flex-col sm:flex-row sm:items-center gap-4">
            <span className="inline-flex items-center justify-center w-16 h-16 rounded-xl bg-[var(--color-primary-soft)] text-[var(--color-primary)] shrink-0">
              <Icon.laptop width={30} height={30} />
            </span>
            <div className="flex-1">
              <div className="flex items-center gap-2">
                <h3 className="font-semibold text-lg">{myAsset.name}</h3>
                <StatusBadge status={myAsset.status} />
              </div>
              <p className="text-sm text-[var(--color-text-2)] mt-0.5">Kode Aset: <Code>{myAsset.code}</Code></p>
            </div>
          </div>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4 mt-6 pt-5 border-t border-[var(--color-border)]">
            <Info icon={<Icon.box width={15} />} label="Merek / Model" value={`${myAsset.brand} ${myAsset.model}`} />
            <Info icon={<Icon.calendar width={15} />} label="Sesi" value={myAsset.session} />
            <Info icon={<Icon.clock width={15} />} label="Waktu Alokasi" value={myAsset.allocatedAt} />
            <Info icon={<Icon.shield width={15} />} label="Kondisi" value={myAsset.condition} />
            <Info icon={<Icon.checkCircle width={15} />} label="Status" value="Sedang Digunakan" />
          </div>
          <p className="mt-4 text-xs text-[var(--color-text-3)]">
            Pengembalian aset dilakukan oleh laboran saat check-out sesi.
          </p>
        </Card>
      ) : (
        <Card padded={false}>
          <Table>
            <THead cols={['Tanggal', 'Kode Aset', 'Sesi', 'Keluar', 'Kembali', 'Kondisi', 'Status']} />
            <tbody>
              {assetHistory.map((a, i) => (
                <TR key={i}>
                  <TD>{a.date}</TD>
                  <TD><Code>{a.code}</Code></TD>
                  <TD>{a.session}</TD>
                  <TD>{a.out}</TD>
                  <TD>{a.in}</TD>
                  <TD>
                    <Badge tone={a.condition === 'Baik' ? 'success' : 'warning'}>{a.condition}</Badge>
                  </TD>
                  <TD><StatusBadge status={a.status} /></TD>
                </TR>
              ))}
            </tbody>
          </Table>
        </Card>
      )}
    </>
  )
}

/* --------------------------------- History --------------------------------- */
function History({ onOpen }: { onOpen: (r: AttendanceRow) => void }) {
  const loading = useSimulatedLoading()
  const [status, setStatus] = useState('Semua')
  const [q, setQ] = useState('')
  const rows = attendanceHistory.filter(
    (r) => (status === 'Semua' || r.status === status) && r.session.toLowerCase().includes(q.toLowerCase()),
  )
  return (
    <>
      <PageHeader title="Riwayat Presensi" subtitle="Seluruh catatan kehadiran praktikum Anda." />
      <Card padded={false}>
        <div className="flex flex-wrap items-center gap-3 p-4 border-b border-[var(--color-border)]">
          <div className="flex-1 min-w-[200px]">
            <Input placeholder="Cari sesi…" value={q} onChange={(e) => setQ(e.target.value)} />
          </div>
          <Select value={status} onChange={(e) => setStatus(e.target.value)} className="w-40">
            {['Semua', 'Selesai', 'Pending', 'Gagal'].map((s) => (
              <option key={s}>{s}</option>
            ))}
          </Select>
        </div>
        {loading ? (
          <Table>
            <THead cols={['Tanggal', 'Sesi', 'Laboratorium', 'Check-in', 'Check-out', 'Status']} />
            <SkeletonRows cols={6} />
          </Table>
        ) : rows.length === 0 ? (
          <EmptyState title="Belum ada catatan presensi" desc="Tidak ada catatan yang cocok dengan filter." icon={<Icon.history />} />
        ) : (
          <>
            <Table>
              <THead cols={['Tanggal', 'Sesi', 'Laboratorium', 'Check-in', 'Check-out', 'Status']} />
              <tbody>
                {rows.map((r, i) => (
                  <TR key={i} onClick={() => onOpen(r)}>
                    <TD>{r.date}</TD>
                    <TD className="font-medium">{r.session}</TD>
                    <TD>{r.room}</TD>
                    <TD>{r.checkIn}</TD>
                    <TD>{r.checkOut}</TD>
                    <TD><StatusBadge status={r.status} /></TD>
                  </TR>
                ))}
              </tbody>
            </Table>
            <Pagination page={1} pages={1} total={rows.length} onPage={() => {}} />
          </>
        )}
      </Card>
    </>
  )
}

/* --------------------------------- Attendance detail --------------------------------- */
function AttendanceDetail({ row, go }: { row: AttendanceRow | null; go: (r: string) => void }) {
  const r = row ?? attendanceHistory[0]
  const pending = r.status === 'Pending'
  const failed = r.status === 'Gagal'
  return (
    <>
      <PageHeader
        title="Detail Presensi"
        subtitle={`${r.session} · ${r.date}`}
        breadcrumb={['Riwayat', 'Detail Presensi']}
        actions={
          <Button variant="secondary" icon={<Icon.arrowLeft width={16} />} onClick={() => go('history')}>
            Kembali
          </Button>
        }
      />
      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-semibold">Ringkasan Kehadiran</h3>
            <StatusBadge status={r.status} />
          </div>
          <div className="rounded-lg border border-[var(--color-border)] divide-y divide-[var(--color-border)]">
            <Line k="Mahasiswa" v={currentStudent.name} />
            <Line k="Sesi" v={r.session} />
            <Line k="Laboratorium" v={r.room} />
            <Line k="Tanggal" v={r.date} />
            <Line k="Check-in" v={r.checkIn} />
            <Line k="Check-out" v={r.checkOut} />
            <Line k="Status Verifikasi" v={failed ? <Badge tone="error">Gagal</Badge> : <Badge tone="success">Terverifikasi</Badge>} />
          </div>
          {pending && (
            <div className="mt-4 flex gap-2">
              <Button icon={<Icon.logout width={16} />} onClick={() => go('checkout')}>Lakukan Check-out</Button>
            </div>
          )}
          {failed && (
            <div className="mt-4">
              <Alert tone="error" title="Verifikasi wajah gagal">
                Presensi tidak tercatat karena verifikasi gagal. Silakan ulangi presensi pada sesi yang tersedia.
              </Alert>
            </div>
          )}
        </Card>

        <Card>
          <h3 className="font-semibold mb-4">Linimasa</h3>
          <Timeline
            items={
              failed
                ? [
                    { title: 'Sesi dibuka', time: r.date, tone: 'neutral' },
                    { title: 'Verifikasi wajah gagal', time: '—', tone: 'error', desc: 'Wajah tidak sesuai / tidak jelas.' },
                  ]
                : [
                    { title: 'Sesi dibuka', time: r.date, tone: 'neutral' },
                    { title: 'Verifikasi wajah berhasil', time: r.checkIn, tone: 'success' },
                    { title: 'Check-in tercatat', time: r.checkIn, tone: 'success' },
                    r.checkOut !== '—'
                      ? { title: 'Check-out tercatat', time: r.checkOut, tone: 'info' }
                      : { title: 'Menunggu check-out', time: '—', tone: 'warning' },
                  ]
            }
          />
        </Card>
      </div>
    </>
  )
}

/* --------------------------------- Check-out --------------------------------- */
function Checkout({ go }: { go: (r: string) => void }) {
  const toast = useToast()
  const [confirm, setConfirm] = useState(false)
  const [done, setDone] = useState(false)

  if (done) {
    return (
      <div className="max-w-xl mx-auto">
        <Card className="text-center animate-fade">
          <span className="mx-auto inline-flex items-center justify-center w-16 h-16 rounded-full bg-[var(--color-success-soft)] text-[var(--color-success)]">
            <Icon.checkCircle width={34} height={34} />
          </span>
          <h2 className="mt-4 text-xl font-semibold tracking-tight">Check-out berhasil</h2>
          <p className="text-sm text-[var(--color-text-2)] mt-1">Sesi praktikum Anda telah selesai dan tercatat.</p>
          <div className="mt-5 text-left rounded-lg border border-[var(--color-border)] divide-y divide-[var(--color-border)]">
            <Line k="Sesi" v={activeSession.name} />
            <Line k="Waktu Check-in" v="13.04" />
            <Line k="Waktu Check-out" v="15.28" />
            <Line k="Durasi Sesi" v="2 jam 24 menit" />
            <Line k="Status" v={<Badge tone="success">Selesai</Badge>} />
          </div>
          <div className="mt-5 flex justify-center gap-2">
            <Button onClick={() => go('history')}>Lihat Riwayat</Button>
            <Button variant="secondary" onClick={() => go('dashboard')}>Kembali ke Dasbor</Button>
          </div>
        </Card>
      </div>
    )
  }

  return (
    <>
      <PageHeader title="Check-out Sesi" subtitle="Akhiri sesi praktikum aktif Anda." breadcrumb={['Presensi', 'Check-out']} />
      <div className="max-w-xl">
        <Card>
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-semibold">Presensi Aktif</h3>
            <Badge tone="success">Sedang Berlangsung</Badge>
          </div>
          <div className="rounded-lg border border-[var(--color-border)] divide-y divide-[var(--color-border)]">
            <Line k="Sesi" v={activeSession.name} />
            <Line k="Laboratorium" v={activeSession.room} />
            <Line k="Waktu Check-in" v="13.04" />
            <Line k="Waktu Saat Ini" v="15.28" />
          </div>
          <div className="mt-5 flex gap-2">
            <Button variant="secondary" onClick={() => go('dashboard')}>Batal</Button>
            <Button icon={<Icon.logout width={16} />} onClick={() => setConfirm(true)}>Check-out Sekarang</Button>
          </div>
        </Card>
      </div>

      <Modal
        open={confirm}
        onClose={() => setConfirm(false)}
        title="Konfirmasi Check-out"
        footer={
          <>
            <Button variant="secondary" onClick={() => setConfirm(false)}>Batal</Button>
            <Button variant="success" icon={<Icon.check width={16} />} onClick={() => { setConfirm(false); setDone(true); toast('Check-out berhasil.') }}>
              Ya, Check-out
            </Button>
          </>
        }
      >
        <p className="text-sm text-[var(--color-text-2)]">
          Anda akan mengakhiri sesi <b>{activeSession.name}</b>. Tindakan ini akan mencatat waktu check-out Anda. Lanjutkan?
        </p>
      </Modal>
    </>
  )
}

/* --------------------------------- Profile --------------------------------- */
function Profile({ go, enrolled, onReset }: { go: (r: string) => void; enrolled: boolean; onReset: () => void }) {
  const toast = useToast()
  return (
    <>
      <PageHeader title="Profil Saya" subtitle="Informasi akun dan status pendaftaran wajah." />
      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-1 text-center">
          <Avatar name={currentStudent.name} size={72} />
          <h3 className="mt-3 font-semibold text-lg">{currentStudent.name}</h3>
          <p className="text-sm text-[var(--color-text-2)] font-mono">{currentStudent.nim}</p>
          <div className="mt-3 flex justify-center gap-2">
            <StatusBadge status={currentStudent.status} />
            <Badge tone={enrolled ? 'success' : 'warning'}>{enrolled ? 'Wajah Terdaftar' : 'Belum Terdaftar'}</Badge>
          </div>
        </Card>

        <div className="lg:col-span-2 space-y-6">
          <Card>
            <h4 className="font-semibold mb-4">Informasi Akun</h4>
            <div className="grid sm:grid-cols-2 gap-4">
              <Field label="Nama Lengkap"><Input defaultValue={currentStudent.name} /></Field>
              <Field label="NIM"><Input defaultValue={currentStudent.nim} disabled /></Field>
              <Field label="Email"><Input defaultValue={currentStudent.email} /></Field>
              <Field label="Program Studi"><Input defaultValue={currentStudent.prodi} disabled /></Field>
            </div>
            <div className="mt-5 flex justify-end">
              <Button onClick={() => toast('Perubahan profil berhasil disimpan.')}>Simpan Perubahan</Button>
            </div>
          </Card>

          <Card>
            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-semibold">Pendaftaran Wajah</h4>
                <p className="text-sm text-[var(--color-text-2)] mt-0.5">
                  {enrolled
                    ? 'Templat wajah aktif dan digunakan untuk verifikasi presensi.'
                    : 'Wajah Anda belum terdaftar. Daftarkan untuk mengaktifkan presensi.'}
                </p>
              </div>
              <Badge tone={enrolled ? 'success' : 'warning'}>{enrolled ? 'Terdaftar' : 'Belum Terdaftar'}</Badge>
            </div>
            <div className="mt-4 flex flex-wrap gap-2">
              <Button icon={<Icon.camera width={16} />} onClick={() => go('enroll')}>
                {enrolled ? 'Daftar Ulang Wajah' : 'Daftar Wajah'}
              </Button>
              {enrolled && (
                <Button variant="ghost" onClick={() => { onReset(); toast('Status pendaftaran wajah direset (demo).', 'warning') }}>
                  Reset (demo)
                </Button>
              )}
            </div>
          </Card>
        </div>
      </div>
    </>
  )
}

/* --------------------------------- Enroll --------------------------------- */
function Enroll({ go, onDone }: { go: (r: string) => void; onDone: () => void }) {
  const toast = useToast()
  return (
    <>
      <PageHeader
        title="Pendaftaran Wajah"
        subtitle="Ikuti panduan untuk mengambil beberapa sudut wajah Anda."
        breadcrumb={['Profil', 'Pendaftaran Wajah']}
      />
      <FaceFlow
        mode="enroll"
        studentName={currentStudent.name}
        onComplete={() => {
          onDone()
          toast('Wajah berhasil didaftarkan.')
          go('profile')
        }}
        onCancel={() => go('profile')}
      />
    </>
  )
}
