import { useEffect, useState } from "react"
import { api, ApiError, unwrapList } from "@/lib/api"
import type {
  AssetAllocation,
  AttendanceRecord,
  AuthUser,
  LabSession,
} from "@/lib/api"
import { Icon } from "../components/icons"
import { FaceFlow } from "../components/FaceFlow"
import { Shell, type NavItem } from "../components/Shell"
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
  StepIndicator,
  Table,
  TD,
  THead,
  TR,
  Tabs,
  Timeline,
  useToast,
} from "../components/ui"

type AttendanceRow = {
  date: string
  session: string
  room: string
  checkIn: string
  checkOut: string
  status: string
}

const NAV: NavItem[] = [
  { key: "dashboard", label: "Dasbor", icon: "dashboard" },
  { key: "sessions", label: "Sesi Praktikum", icon: "calendar" },
  { key: "attendance", label: "Presensi", icon: "checkCircle" },
  { key: "assets", label: "Aset Saya", icon: "laptop" },
  { key: "history", label: "Riwayat", icon: "history" },
  { key: "profile", label: "Profil", icon: "user" },
]

function fmtDate(value?: string) {
  if (!value) return "-"
  const date = new Date(value)
  return Number.isNaN(date.getTime())
    ? value
    : date.toLocaleDateString("id-ID", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      })
}

function fmtTime(value?: string) {
  if (!value) return "-"
  const date = new Date(value)
  return Number.isNaN(date.getTime())
    ? value
    : date.toLocaleTimeString("id-ID", { hour: "2-digit", minute: "2-digit" })
}

function sessionStatus(status?: string) {
  return status === "ACTIVE"
    ? "Aktif"
    : status === "CLOSED"
      ? "Ditutup"
      : status === "CANCELLED"
        ? "Dibatalkan"
        : "Akan Datang"
}

function errorCode(err: unknown) {
  return err instanceof ApiError ? err.code : "CLIENT_ERROR"
}

export function StudentApp({
  onLogout,
  user,
}: {
  onLogout: () => void
  user?: AuthUser | null
}) {
  const [route, setRoute] = useState("dashboard")
  const [detail, setDetail] = useState<AttendanceRow | null>(null)
  const [enrolled, setEnrolled] = useState(user?.face_enrolled ?? false)
  const [apiState, setApiState] = useState({
    loading: true,
    error: "",
    sessions: [] as LabSession[],
    activeSessions: [] as LabSession[],
    attendance: [] as AttendanceRecord[],
    allocations: [] as AssetAllocation[],
  })

  const refreshStudentData = async () => {
    setApiState((s) => ({ ...s, loading: true, error: "" }))
    try {
      const [face, allSessions, activeSessions, attendance, allocations] =
        await Promise.all([
          api.face.status(),
          api.sessions.list(),
          api.sessions.active(),
          api.attendance.history(),
          api.assets.myAllocations(),
        ])
      setEnrolled(face.enrolled)
      setApiState({
        loading: false,
        error: "",
        sessions: unwrapList(allSessions),
        activeSessions: unwrapList(activeSessions),
        attendance: unwrapList(attendance),
        allocations: unwrapList(allocations),
      })
    } catch (err) {
      setApiState((s) => ({
        ...s,
        loading: false,
        error:
          err instanceof ApiError
            ? err.message
            : "Gagal memuat data mahasiswa dari server.",
      }))
    }
  }

  // Auto-refresh on route change for realtime feel
  useEffect(() => {
    void refreshStudentData()
  }, [route])

  const active =
    apiState.activeSessions[0] ||
    apiState.sessions.find((s) => s.status === "ACTIVE")
  const openDetail = (r: AttendanceRow) => {
    setDetail(r)
    setRoute("attendanceDetail")
  }

  // Cek apakah user sudah check-in di sesi tertentu
  const isCheckedIn = (sessionId?: string) =>
    sessionId ? apiState.attendance.some((a) => a.session_id === sessionId) : false
  const getAttendance = (sessionId?: string) =>
    sessionId ? apiState.attendance.find((a) => a.session_id === sessionId) : undefined

  return (
    <Shell
      brand="LabPresensi"
      roleLabel="Mahasiswa"
      roleTone="primary"
      nav={NAV}
      active={route}
      onNavigate={setRoute}
      user={{
        name: user?.full_name || "Mahasiswa",
        sub: user?.identifier || "-",
      }}
      onLogout={onLogout}
    >
      {route === "dashboard" && (
        <Dashboard
          go={setRoute}
          enrolled={enrolled}
          user={user}
          active={active}
          allocation={apiState.allocations[0]}
          error={apiState.error}
          isCheckedIn={isCheckedIn(active?.id)}
          activeAttendance={getAttendance(active?.id)}
        />
      )}
      {route === "sessions" && (
        <Sessions
          go={setRoute}
          rows={apiState.sessions}
          loading={apiState.loading}
          attendance={apiState.attendance}
        />
      )}
      {route === "attendance" && (
        <AttendanceView
          go={setRoute}
          enrolled={enrolled}
          user={user}
          activeSessions={apiState.activeSessions.length > 0 ? apiState.activeSessions : apiState.sessions.filter(s => s.status === "ACTIVE")}
          onRefresh={refreshStudentData}
          attendance={apiState.attendance}
        />
      )}
      {route === "assets" && <MyAssets rows={apiState.allocations} />}
      {route === "history" && (
        <History
          onOpen={openDetail}
          rows={apiState.attendance}
          loading={apiState.loading}
        />
      )}
      {route === "attendanceDetail" && (
        <AttendanceDetail row={detail} go={setRoute} user={user} />
      )}
      {route === "checkout" && (
        <Checkout
          go={setRoute}
          activeSessions={apiState.activeSessions.length > 0 ? apiState.activeSessions : apiState.sessions.filter(s => s.status === "ACTIVE")}
          attendance={apiState.attendance}
          onRefresh={refreshStudentData}
        />
      )}
      {route === "profile" && (
        <Profile
          go={setRoute}
          user={user}
          enrolled={enrolled}
          onReset={() => setEnrolled(false)}
        />
      )}
      {route === "enroll" && (
        <Enroll
          go={setRoute}
          user={user}
          onDone={() => {
            setEnrolled(true)
            void refreshStudentData()
          }}
        />
      )}
    </Shell>
  )
}

/* --------------------------------- Dashboard --------------------------------- */
function Dashboard({
  go,
  enrolled,
  user,
  active,
  allocation,
  error,
  isCheckedIn,
  activeAttendance,
}: {
  go: (r: string) => void
  enrolled: boolean
  user?: AuthUser | null
  active?: LabSession
  allocation?: AssetAllocation
  error?: string
  isCheckedIn: boolean
  activeAttendance?: AttendanceRecord
}) {
  const name = user?.full_name || "Mahasiswa"
  const asset = allocation?.asset
  const hasCheckedOut = !!activeAttendance?.checked_out_at

  // Status presensi realtime
  const presenceLabel = !active
    ? "Tidak Ada Sesi"
    : isCheckedIn
      ? hasCheckedOut
        ? "Selesai"
        : "Hadir"
      : "Belum Hadir"
  const presenceTone = !active
    ? "warning"
    : isCheckedIn
      ? hasCheckedOut
        ? "info"
        : ("success" as const)
      : "warning"

  return (
    <>
      <PageHeader
        title={`Halo, ${name.split(" ")[0]}`}
        subtitle="Berikut ringkasan praktikum dan presensi Anda hari ini."
      />

      {!enrolled && (
        <div className="mb-6">
          <Alert tone="warning" title="Wajah belum terdaftar">
            <div className="flex flex-wrap items-center gap-3">
              <span>
                Daftarkan wajah Anda terlebih dahulu agar dapat melakukan
                presensi.
              </span>
              <Button
                size="sm"
                icon={<Icon.camera width={15} />}
                onClick={() => go("enroll")}
              >
                Daftar Wajah Sekarang
              </Button>
            </div>
          </Alert>
        </div>
      )}

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="Sesi Hari Ini"
          value={active ? "1 Sesi" : "0 Sesi"}
          hint={
            active
              ? `${active.title} · ${fmtTime(active.start_at)}`
              : "Tidak ada sesi aktif"
          }
          icon={<Icon.calendar />}
          tone="primary"
        />
        <StatCard
          label="Status Presensi"
          value={presenceLabel}
          hint={
            isCheckedIn
              ? `Check-in: ${fmtTime(activeAttendance?.checked_in_at)}`
              : active
                ? "Sesi aktif menunggu presensi"
                : "Menunggu sesi dibuka"
          }
          icon={<Icon.clock />}
          tone={presenceTone}
        />
        <StatCard
          label="Aset Aktif"
          value={asset?.asset_code || "-"}
          hint={asset?.asset_name || "Tidak ada alokasi aktif"}
          icon={<Icon.laptop />}
          tone="info"
        />
        <StatCard
          label="Status Verifikasi"
          value={enrolled ? "Terdaftar" : "Belum Terdaftar"}
          hint={enrolled ? "Templat wajah aktif" : "Perlu pendaftaran wajah"}
          icon={<Icon.shield />}
          tone={enrolled ? "success" : "warning"}
        />
      </div>

      <div className="grid gap-6 lg:grid-cols-3 mt-6">
        {/* Active session */}
        <Card className="lg:col-span-2 border-[var(--color-primary)]/30 ring-1 ring-[var(--color-primary)]/10">
          <div className="flex items-start justify-between gap-3">
            <div>
              <Badge tone="success">Sesi Aktif</Badge>
              <h3 className="mt-2 text-lg font-semibold tracking-tight">
                {active?.title || "Tidak ada sesi aktif"}
              </h3>
              <p className="text-sm text-[var(--color-text-2)]">
                {active?.course_code || "-"}
              </p>
            </div>
            {active && <Code>{active.code}</Code>}
          </div>

          {active && (
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-5">
              <Info
                icon={<Icon.calendar width={16} />}
                label="Tanggal"
                value={fmtDate(active.session_date)}
              />
              <Info
                icon={<Icon.clock width={16} />}
                label="Waktu"
                value={`${fmtTime(active.start_at)}-${fmtTime(active.end_at)}`}
              />
              <Info
                icon={<Icon.box width={16} />}
                label="Ruang"
                value={active.location}
              />
              <Info
                icon={<Icon.checkCircle width={16} />}
                label="Status"
                value={presenceLabel}
              />
            </div>
          )}

          <div className="mt-6 flex flex-wrap gap-2">
            {!isCheckedIn ? (
              <Button
                size="lg"
                icon={<Icon.camera width={18} height={18} />}
                disabled={!active}
                onClick={() => go("attendance")}
              >
                Mulai Presensi
              </Button>
            ) : !hasCheckedOut ? (
              <Button
                size="lg"
                icon={<Icon.logout width={18} height={18} />}
                onClick={() => go("checkout")}
              >
                Check-out Sesi
              </Button>
            ) : (
              <Button size="lg" variant="secondary" disabled>
                Sesi Selesai ✓
              </Button>
            )}
            <Button
              size="lg"
              variant="secondary"
              onClick={() => go("sessions")}
            >
              Lihat Semua Sesi
            </Button>
          </div>
        </Card>

        {/* Current asset */}
        <Card>
          <div className="flex items-center justify-between">
            <h3 className="font-semibold">Alokasi Aset Saat Ini</h3>
            <StatusBadge status={allocation ? (allocation.returned_at ? "Selesai" : "Aktif") : "-"} />
          </div>
          <div className="mt-4 flex items-center gap-3 p-3 rounded-lg bg-[var(--color-surface-2)] border border-[var(--color-border)]">
            <span className="inline-flex items-center justify-center w-11 h-11 rounded-lg bg-white border border-[var(--color-border)] text-[var(--color-primary)]">
              <Icon.laptop />
            </span>
            <div className="min-w-0">
              <p className="font-medium text-sm truncate">
                {asset?.asset_name || "Tidak ada aset aktif"}
              </p>
              <p className="text-xs text-[var(--color-text-2)] font-mono">
                {asset?.asset_code || "-"}
              </p>
            </div>
          </div>
          <dl className="mt-4 space-y-2.5 text-sm">
            <Row k="Sesi" v={allocation?.session?.title || "-"} />
            <Row
              k="Dialokasikan"
              v={fmtTime(allocation?.allocated_at)}
            />
            <Row
              k="Kondisi"
              v={<Badge tone="success">{allocation?.condition || "-"}</Badge>}
            />
          </dl>
          <Button
            variant="secondary"
            className="w-full mt-4"
            onClick={() => go("assets")}
          >
            Detail Aset
          </Button>
        </Card>
      </div>
    </>
  )
}

function Info({
  icon,
  label,
  value,
}: {
  icon: React.ReactNode
  label: string
  value: string
}) {
  return (
    <div>
      <div className="flex items-center gap-1.5 text-[var(--color-text-3)] text-xs">
        {icon}
        {label}
      </div>
      <p className="mt-1 text-sm font-medium text-[var(--color-text)]">
        {value}
      </p>
    </div>
  )
}
type KeyValueLineProps = {
  k: string
  v: React.ReactNode
}

function Row({ k, v }: KeyValueLineProps) {
  return (
    <div className="flex items-center justify-between gap-3">
      <dt className="text-[var(--color-text-2)]">{k}</dt>
      <dd className="font-medium text-right">{v}</dd>
    </div>
  )
}

/* --------------------------------- Sessions --------------------------------- */
function Sessions({
  go,
  rows,
  loading,
  attendance,
}: {
  go: (r: string) => void
  rows: LabSession[]
  loading: boolean
  attendance: AttendanceRecord[]
}) {
  const [status, setStatus] = useState("Semua")
  const source = rows.length ? rows : []
  const filtered = source.filter(
    (s) => status === "Semua" || sessionStatus(s.status) === status,
  )
  // Lookup set attendance session IDs
  const attendedSessionIds = new Set(attendance.map((a) => a.session_id))

  return (
    <>
      <PageHeader
        title="Sesi Praktikum"
        subtitle="Jadwal dan status sesi praktikum Anda."
        actions={
          <Select
            value={status}
            onChange={(e) => setStatus(e.target.value)}
            className="w-44"
          >
            {["Semua", "Aktif", "Akan Datang", "Ditutup", "Dibatalkan"].map(
              (s) => (
                <option key={s}>{s}</option>
              ),
            )}
          </Select>
        }
      />
      {loading ? (
        <Card padded={false}>
          <Table>
            <SkeletonRows cols={4} rows={4} />
          </Table>
        </Card>
      ) : filtered.length === 0 ? (
        <Card>
          <EmptyState
            title="Tidak ada sesi"
            desc="Tidak ada sesi praktikum untuk filter yang dipilih."
            icon={<Icon.calendar />}
          />
        </Card>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {filtered.map((s) => {
            const isActive = s.status === "ACTIVE"
            const attended = attendedSessionIds.has(s.id)
            return (
              <Card
                key={s.id || s.code}
                className={
                  isActive
                    ? "ring-1 ring-[var(--color-primary)]/20 border-[var(--color-primary)]/30"
                    : ""
                }
              >
                <div className="flex items-start justify-between gap-2">
                  <StatusBadge status={sessionStatus(s.status)} />
                  <Code>{s.code}</Code>
                </div>
                <h3 className="mt-3 font-semibold">{s.title}</h3>
                <p className="text-sm text-[var(--color-text-2)]">
                  {s.course_code}
                </p>
                <div className="mt-4 space-y-2 text-sm text-[var(--color-text-2)]">
                  <p className="flex items-center gap-2">
                    <Icon.calendar width={15} /> {fmtDate(s.session_date)} ·{" "}
                    {fmtTime(s.start_at)}-{fmtTime(s.end_at)}
                  </p>
                  <p className="flex items-center gap-2">
                    <Icon.box width={15} /> {s.location}
                  </p>
                  <p className="flex items-center gap-2">
                    <Icon.checkCircle width={15} /> Presensi:{" "}
                    <StatusBadge
                      status={attended ? "Hadir" : "Belum Hadir"}
                    />
                  </p>
                </div>
                <div className="mt-4">
                  {isActive && !attended ? (
                    <Button
                      className="w-full"
                      onClick={() => go("attendance")}
                      icon={<Icon.camera width={16} />}
                    >
                      Mulai Presensi
                    </Button>
                  ) : isActive && attended ? (
                    <Button variant="secondary" className="w-full" disabled>
                      Sudah Hadir ✓
                    </Button>
                  ) : (
                    <Button
                      variant="secondary"
                      className="w-full"
                      disabled={sessionStatus(s.status) === "Dibatalkan"}
                    >
                      {sessionStatus(s.status) === "Ditutup"
                        ? "Lihat Detail"
                        : sessionStatus(s.status) === "Dibatalkan"
                          ? "Dibatalkan"
                          : "Belum Dibuka"}
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
function AttendanceView({
  go,
  enrolled,
  user,
  activeSessions,
  onRefresh,
  attendance,
}: {
  go: (r: string) => void
  enrolled: boolean
  user?: AuthUser | null
  activeSessions: LabSession[]
  onRefresh: () => Promise<void>
  attendance: AttendanceRecord[]
}) {
  const [step, setStep] = useState<"verify" | "success">("verify")
  const [saved, setSaved] = useState<AttendanceRecord | null>(null)
  // Pilih sesi yang belum di-check-in
  const attendedIds = new Set(attendance.map((a) => a.session_id))
  const availableSessions = activeSessions.filter((s) => !attendedIds.has(s.id))
  const [selectedIdx, setSelectedIdx] = useState(0)
  const active = availableSessions[selectedIdx] || availableSessions[0]

  const submitCheckIn = async (payload: Blob | Record<string, Blob>) => {
    if (!active || !(payload instanceof Blob)) return
    const form = new FormData()
    form.append("session_id", active.id)
    form.append("image", payload, "check-in.jpg")
    const res = await api.attendance.checkIn(form)
    setSaved(res.attendance)
    await onRefresh()
  }

  if (!enrolled) {
    return (
      <>
        <PageHeader
          title="Presensi — Verifikasi Wajah"
          subtitle="Verifikasi identitas Anda untuk mencatat kehadiran."
        />
        <div className="max-w-xl">
          <Alert tone="warning" title="Wajah belum terdaftar">
            Anda perlu mendaftarkan wajah terlebih dahulu sebelum dapat
            melakukan verifikasi presensi.
          </Alert>
          <div className="mt-4 flex gap-2">
            <Button
              icon={<Icon.camera width={16} />}
              onClick={() => go("enroll")}
            >
              Daftar Wajah
            </Button>
            <Button variant="secondary" onClick={() => go("dashboard")}>
              Kembali
            </Button>
          </div>
        </div>
      </>
    )
  }

  if (availableSessions.length === 0 && step === "verify") {
    return (
      <>
        <PageHeader
          title="Presensi — Verifikasi Wajah"
          subtitle="Verifikasi identitas Anda untuk mencatat kehadiran."
        />
        <div className="max-w-xl">
          <Alert tone="info" title="Semua sesi sudah dihadiri">
            Anda sudah melakukan presensi pada semua sesi aktif yang tersedia.
          </Alert>
          <div className="mt-4">
            <Button variant="secondary" onClick={() => go("dashboard")}>
              Kembali ke Dasbor
            </Button>
          </div>
        </div>
      </>
    )
  }

  if (step === "success") {
    return (
      <div className="max-w-xl mx-auto">
        <Card className="text-center animate-fade">
          <span className="mx-auto inline-flex items-center justify-center w-16 h-16 rounded-full bg-[var(--color-success-soft)] text-[var(--color-success)]">
            <Icon.checkCircle width={34} height={34} />
          </span>
          <h2 className="mt-4 text-xl font-semibold tracking-tight">
            Presensi berhasil
          </h2>
          <p className="text-sm text-[var(--color-text-2)] mt-1">
            Kehadiran Anda telah tercatat untuk sesi ini.
          </p>
          <div className="mt-5 text-left rounded-lg border border-[var(--color-border)] divide-y divide-[var(--color-border)]">
            <Line k="Mahasiswa" v={user?.full_name || "Mahasiswa"} />
            <Line k="Sesi" v={active?.title || "-"} />
            <Line k="Waktu Check-in" v={fmtTime(saved?.checked_in_at)} />
            <Line k="Ruang" v={active?.location || "-"} />
            <Line k="Status" v={<Badge tone="success">Hadir</Badge>} />
          </div>
          <div className="mt-5 flex flex-col sm:flex-row gap-2 justify-center">
            <Button
              onClick={() => go("history")}
              icon={<Icon.history width={16} />}
            >
              Lihat Detail Presensi
            </Button>
            <Button variant="secondary" onClick={() => go("assets")}>
              Lihat Alokasi Aset
            </Button>
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
        breadcrumb={["Presensi", "Verifikasi"]}
      />
      {/* Pemilih sesi jika ada lebih dari satu sesi aktif */}
      {availableSessions.length > 1 && (
        <div className="mb-5">
          <Field label="Pilih Sesi">
            <Select value={String(selectedIdx)} onChange={(e) => setSelectedIdx(Number(e.target.value))}>
              {availableSessions.map((s, i) => (
                <option key={s.id} value={String(i)}>{s.code} — {s.title}</option>
              ))}
            </Select>
          </Field>
        </div>
      )}
      <div className="mb-5">
        <Alert tone="info" title={`Presensi untuk: ${active?.title || "-"}`}>
          Sistem membandingkan wajah di depan kamera dengan templat wajah akun{" "}
          <b>{user?.full_name || "Mahasiswa"}</b> yang sedang login.
        </Alert>
      </div>
      {active && (
        <FaceFlow
          mode="verify"
          studentName={user?.full_name || "Mahasiswa"}
          sessionName={active.title}
          onSubmit={submitCheckIn}
          onComplete={() => setStep("success")}
          onCancel={() => go("dashboard")}
        />
      )}
    </>
  )
}
function Line({ k, v }: KeyValueLineProps) {
  return (
    <div className="flex items-center justify-between gap-3 px-4 py-2.5 text-sm">
      <span className="text-[var(--color-text-2)]">{k}</span>
      <span className="font-medium text-right">{v}</span>
    </div>
  )
}

/* --------------------------------- My assets --------------------------------- */
function MyAssets({ rows }: { rows: AssetAllocation[] }) {
  const [tab, setTab] = useState("Alokasi Aktif")
  const current = rows.find((r) => !r.returned_at) || rows[0]
  const asset = current?.asset
  return (
    <>
      <PageHeader
        title="Aset Saya"
        subtitle="Aset yang sedang dan pernah dialokasikan kepada Anda."
      />
      <Tabs
        tabs={["Alokasi Aktif", "Riwayat Aset"]}
        active={tab}
        onChange={setTab}
      />

      {tab === "Alokasi Aktif" ? (
        <Card>
          <div className="flex flex-col sm:flex-row sm:items-center gap-4">
            <span className="inline-flex items-center justify-center w-16 h-16 rounded-xl bg-[var(--color-primary-soft)] text-[var(--color-primary)] shrink-0">
              <Icon.laptop width={30} height={30} />
            </span>
            <div className="flex-1">
              <div className="flex items-center gap-2">
                <h3 className="font-semibold text-lg">
                  {asset?.asset_name || "Tidak ada aset aktif"}
                </h3>
                <StatusBadge status={current?.status || "-"} />
              </div>
              <p className="text-sm text-[var(--color-text-2)] mt-0.5">
                Kode Aset: <Code>{asset?.asset_code || "-"}</Code>
              </p>
            </div>
          </div>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4 mt-6 pt-5 border-t border-[var(--color-border)]">
            <Info
              icon={<Icon.box width={15} />}
              label="Merek / Model"
              value={`${asset?.brand || "-"} ${asset?.model || ""}`}
            />
            <Info
              icon={<Icon.calendar width={15} />}
              label="Sesi"
              value={current?.session?.title || "-"}
            />
            <Info
              icon={<Icon.clock width={15} />}
              label="Waktu Alokasi"
              value={fmtTime(current?.allocated_at)}
            />
            <Info
              icon={<Icon.shield width={15} />}
              label="Kondisi"
              value={current?.condition || "-"}
            />
            <Info
              icon={<Icon.checkCircle width={15} />}
              label="Status"
              value={current ? "Sedang Digunakan" : "Tidak Ada"}
            />
          </div>
          <p className="mt-4 text-xs text-[var(--color-text-3)]">
            Pengembalian aset dilakukan oleh laboran saat check-out sesi.
          </p>
        </Card>
      ) : (
        <Card padded={false}>
          <Table>
            <THead
              cols={[
                "Tanggal",
                "Kode Aset",
                "Sesi",
                "Keluar",
                "Kembali",
                "Kondisi",
                "Status",
              ]}
            />
            <tbody>
              {rows.map((a, i) => (
                <TR key={i}>
                  <TD>{fmtDate(a.allocated_at)}</TD>
                  <TD>
                    <Code>{a.asset?.asset_code || "-"}</Code>
                  </TD>
                  <TD>{a.session?.title || "-"}</TD>
                  <TD>{fmtTime(a.allocated_at)}</TD>
                  <TD>{fmtTime(a.returned_at)}</TD>
                  <TD>
                    <Badge
                      tone={
                        a.condition === "GOOD" || a.condition === "Baik"
                          ? "success"
                          : "warning"
                      }
                    >
                      {a.condition || "-"}
                    </Badge>
                  </TD>
                  <TD>
                    <StatusBadge status={a.returned_at ? "Selesai" : "Aktif"} />
                  </TD>
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
function toAttendanceRow(r: AttendanceRecord): AttendanceRow {
  const title = r.session?.title || r.session_id
  return {
    date: fmtDate(r.session?.session_date || r.checked_in_at),
    session: title,
    room: r.session?.location || "-",
    checkIn: fmtTime(r.checked_in_at),
    checkOut: fmtTime(r.checked_out_at),
    status: r.checked_out_at
      ? "Selesai"
      : r.checked_in_at
        ? "Pending"
        : "Gagal",
  }
}

function History({
  onOpen,
  rows: apiRows,
  loading,
}: {
  onOpen: (r: AttendanceRow) => void
  rows: AttendanceRecord[]
  loading: boolean
}) {
  const [status, setStatus] = useState("Semua")
  const [q, setQ] = useState("")
  const rows = apiRows
    .map(toAttendanceRow)
    .filter(
      (r) =>
        (status === "Semua" || r.status === status) &&
        r.session.toLowerCase().includes(q.toLowerCase()),
    )
  return (
    <>
      <PageHeader
        title="Riwayat Presensi"
        subtitle="Seluruh catatan kehadiran praktikum Anda."
      />
      <Card padded={false}>
        <div className="flex flex-wrap items-center gap-3 p-4 border-b border-[var(--color-border)]">
          <div className="flex-1 min-w-[200px]">
            <Input
              placeholder="Cari sesi…"
              value={q}
              onChange={(e) => setQ(e.target.value)}
            />
          </div>
          <Select
            value={status}
            onChange={(e) => setStatus(e.target.value)}
            className="w-40"
          >
            {["Semua", "Selesai", "Pending", "Gagal"].map((s) => (
              <option key={s}>{s}</option>
            ))}
          </Select>
        </div>
        {loading ? (
          <Table>
            <THead
              cols={[
                "Tanggal",
                "Sesi",
                "Laboratorium",
                "Check-in",
                "Check-out",
                "Status",
              ]}
            />
            <SkeletonRows cols={6} />
          </Table>
        ) : rows.length === 0 ? (
          <EmptyState
            title="Belum ada catatan presensi"
            desc="Tidak ada catatan yang cocok dengan filter."
            icon={<Icon.history />}
          />
        ) : (
          <>
            <Table>
              <THead
                cols={[
                  "Tanggal",
                  "Sesi",
                  "Laboratorium",
                  "Check-in",
                  "Check-out",
                  "Status",
                ]}
              />
              <tbody>
                {rows.map((r, i) => (
                  <TR key={i} onClick={() => onOpen(r)}>
                    <TD>{r.date}</TD>
                    <TD className="font-medium">{r.session}</TD>
                    <TD>{r.room}</TD>
                    <TD>{r.checkIn}</TD>
                    <TD>{r.checkOut}</TD>
                    <TD>
                      <StatusBadge status={r.status} />
                    </TD>
                  </TR>
                ))}
              </tbody>
            </Table>
            <Pagination
              page={1}
              pages={1}
              total={rows.length}
              onPage={() => {}}
            />
          </>
        )}
      </Card>
    </>
  )
}

/* --------------------------------- Attendance detail --------------------------------- */
function AttendanceDetail({
  row,
  go,
  user,
}: {
  row: AttendanceRow | null
  go: (r: string) => void
  user?: AuthUser | null
}) {
  if (!row) {
    go("history")
    return null
  }
  const r = row
  const pending = r.status === "Pending"
  const failed = r.status === "Gagal"
  return (
    <>
      <PageHeader
        title="Detail Presensi"
        subtitle={`${r.session} · ${r.date}`}
        breadcrumb={["Riwayat", "Detail Presensi"]}
        actions={
          <Button
            variant="secondary"
            icon={<Icon.arrowLeft width={16} />}
            onClick={() => go("history")}
          >
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
            <Line k="Mahasiswa" v={user?.full_name || "Mahasiswa"} />
            <Line k="Sesi" v={r.session} />
            <Line k="Laboratorium" v={r.room} />
            <Line k="Tanggal" v={r.date} />
            <Line k="Check-in" v={r.checkIn} />
            <Line k="Check-out" v={r.checkOut} />
            <Line
              k="Status Verifikasi"
              v={
                failed ? (
                  <Badge tone="error">Gagal</Badge>
                ) : (
                  <Badge tone="success">Terverifikasi</Badge>
                )
              }
            />
          </div>
          {pending && (
            <div className="mt-4 flex gap-2">
              <Button
                icon={<Icon.logout width={16} />}
                onClick={() => go("checkout")}
              >
                Lakukan Check-out
              </Button>
            </div>
          )}
          {failed && (
            <div className="mt-4">
              <Alert tone="error" title="Verifikasi wajah gagal">
                Presensi tidak tercatat karena verifikasi gagal. Silakan ulangi
                presensi pada sesi yang tersedia.
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
                    { title: "Sesi dibuka", time: r.date, tone: "neutral" },
                    {
                      title: "Verifikasi wajah gagal",
                      time: "—",
                      tone: "error",
                      desc: "Wajah tidak sesuai / tidak jelas.",
                    },
                  ]
                : [
                    { title: "Sesi dibuka", time: r.date, tone: "neutral" },
                    {
                      title: "Verifikasi wajah berhasil",
                      time: r.checkIn,
                      tone: "success",
                    },
                    {
                      title: "Check-in tercatat",
                      time: r.checkIn,
                      tone: "success",
                    },
                    r.checkOut !== "-"
                      ? {
                          title: "Check-out tercatat",
                          time: r.checkOut,
                          tone: "info",
                        }
                      : {
                          title: "Menunggu check-out",
                          time: "—",
                          tone: "warning",
                        },
                  ]
            }
          />
        </Card>
      </div>
    </>
  )
}

/* --------------------------------- Check-out --------------------------------- */
function Checkout({
  go,
  activeSessions,
  attendance,
  onRefresh,
}: {
  go: (r: string) => void
  activeSessions: LabSession[]
  attendance: AttendanceRecord[]
  onRefresh: () => Promise<void>
}) {
  const toast = useToast()
  const [confirm, setConfirm] = useState(false)
  const [done, setDone] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState("")
  const [saved, setSaved] = useState<AttendanceRecord | null>(null)

  // Sesi yang bisa di-checkout: sudah check-in tapi belum check-out
  const checkedInIds = new Set(
    attendance.filter((a) => a.checked_in_at && !a.checked_out_at).map((a) => a.session_id)
  )
  const checkoutableSessions = activeSessions.filter((s) => checkedInIds.has(s.id))
  const [selectedIdx, setSelectedIdx] = useState(0)
  const active = checkoutableSessions[selectedIdx] || checkoutableSessions[0]

  const submit = async () => {
    if (!active || submitting) return
    setSubmitting(true)
    setError("")
    try {
      const res = await api.attendance.checkOut(active.id)
      setSaved(res.attendance)
      await onRefresh()
      setConfirm(false)
      setDone(true)
      toast("Check-out berhasil.")
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Check-out gagal.")
    } finally {
      setSubmitting(false)
    }
  }

  if (done) {
    return (
      <div className="max-w-xl mx-auto">
        <Card className="text-center animate-fade">
          <span className="mx-auto inline-flex items-center justify-center w-16 h-16 rounded-full bg-[var(--color-success-soft)] text-[var(--color-success)]">
            <Icon.checkCircle width={34} height={34} />
          </span>
          <h2 className="mt-4 text-xl font-semibold tracking-tight">
            Check-out berhasil
          </h2>
          <p className="text-sm text-[var(--color-text-2)] mt-1">
            Sesi praktikum Anda telah selesai dan tercatat.
          </p>
          <div className="mt-5 text-left rounded-lg border border-[var(--color-border)] divide-y divide-[var(--color-border)]">
            <Line k="Sesi" v={active?.title || "-"} />
            <Line k="Waktu Check-in" v={fmtTime(saved?.checked_in_at)} />
            <Line k="Waktu Check-out" v={fmtTime(saved?.checked_out_at)} />
            <Line k="Status" v={<Badge tone="success">Selesai</Badge>} />
          </div>
          <div className="mt-5 flex justify-center gap-2">
            <Button onClick={() => go("history")}>Lihat Riwayat</Button>
            <Button variant="secondary" onClick={() => go("dashboard")}>
              Kembali ke Dasbor
            </Button>
          </div>
        </Card>
      </div>
    )
  }

  return (
    <>
      <PageHeader
        title="Check-out Sesi"
        subtitle="Akhiri sesi praktikum aktif Anda."
        breadcrumb={["Presensi", "Check-out"]}
      />
      <div className="max-w-xl">
        {checkoutableSessions.length > 1 && (
          <div className="mb-4">
            <Field label="Pilih Sesi">
              <Select value={String(selectedIdx)} onChange={(e) => setSelectedIdx(Number(e.target.value))}>
                {checkoutableSessions.map((s, i) => (
                  <option key={s.id} value={String(i)}>{s.code} — {s.title}</option>
                ))}
              </Select>
            </Field>
          </div>
        )}
        <Card>
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-semibold">Presensi Aktif</h3>
            <Badge tone="success">Sedang Berlangsung</Badge>
          </div>
          <div className="rounded-lg border border-[var(--color-border)] divide-y divide-[var(--color-border)]">
            <Line k="Sesi" v={active?.title || "-"} />
            <Line k="Laboratorium" v={active?.location || "-"} />
          </div>
          {!active && (
            <div className="mt-4">
              <Alert tone="warning" title="Tidak ada sesi untuk di-checkout">
                Anda belum melakukan check-in pada sesi aktif manapun.
              </Alert>
            </div>
          )}
          {error && (
            <div className="mt-4">
              <Alert tone="error" title="Check-out gagal">
                {error}
              </Alert>
            </div>
          )}
          <div className="mt-5 flex gap-2">
            <Button variant="secondary" onClick={() => go("dashboard")}>
              Batal
            </Button>
            <Button
              icon={<Icon.logout width={16} />}
              disabled={!active || submitting}
              onClick={() => setConfirm(true)}
            >
              Check-out Sekarang
            </Button>
          </div>
        </Card>
      </div>

      <Modal
        open={confirm}
        onClose={() => setConfirm(false)}
        title="Konfirmasi Check-out"
        footer={
          <>
            <Button variant="secondary" onClick={() => setConfirm(false)}>
              Batal
            </Button>
            <Button
              variant="success"
              icon={<Icon.check width={16} />}
              disabled={submitting}
              onClick={submit}
            >
              {submitting ? "Memproses..." : "Ya, Check-out"}
            </Button>
          </>
        }
      >
        <p className="text-sm text-[var(--color-text-2)]">
          Anda akan mengakhiri sesi <b>{active?.title || "-"}</b>. Tindakan ini
          akan mencatat waktu check-out Anda. Lanjutkan?
        </p>
      </Modal>
    </>
  )
}

/* --------------------------------- Profile --------------------------------- */
function Profile({
  go,
  user,
  enrolled,
  onReset,
}: {
  go: (r: string) => void
  user?: AuthUser | null
  enrolled: boolean
  onReset: () => void
}) {
  const toast = useToast()
  const name = user?.full_name || "Mahasiswa"
  const nim = user?.identifier || "-"
  return (
    <>
      <PageHeader
        title="Profil Saya"
        subtitle="Informasi akun dan status pendaftaran wajah."
      />
      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-1 text-center">
          <Avatar name={name} size={72} />
          <h3 className="mt-3 font-semibold text-lg">{name}</h3>
          <p className="text-sm text-[var(--color-text-2)] font-mono">{nim}</p>
          <div className="mt-3 flex justify-center gap-2">
            <StatusBadge status={user?.status || "ACTIVE"} />
            <Badge tone={enrolled ? "success" : "warning"}>
              {enrolled ? "Wajah Terdaftar" : "Belum Terdaftar"}
            </Badge>
          </div>
        </Card>

        <div className="lg:col-span-2 space-y-6">
          <Card>
            <h4 className="font-semibold mb-4">Informasi Akun</h4>
            <div className="grid sm:grid-cols-2 gap-4">
              <Field label="Nama Lengkap">
                <Input value={name} readOnly />
              </Field>
              <Field label="NIM">
                <Input value={nim} disabled />
              </Field>
              <Field label="Email">
                <Input value={user?.email || "-"} readOnly />
              </Field>
            </div>
          </Card>

          <Card>
            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-semibold">Pendaftaran Wajah</h4>
                <p className="text-sm text-[var(--color-text-2)] mt-0.5">
                  {enrolled
                    ? "Templat wajah aktif dan digunakan untuk verifikasi presensi."
                    : "Wajah Anda belum terdaftar. Daftarkan untuk mengaktifkan presensi."}
                </p>
              </div>
              <Badge tone={enrolled ? "success" : "warning"}>
                {enrolled ? "Terdaftar" : "Belum Terdaftar"}
              </Badge>
            </div>
            <div className="mt-4 flex flex-wrap gap-2">
              <Button
                icon={<Icon.camera width={16} />}
                onClick={() => go("enroll")}
              >
                {enrolled ? "Daftar Ulang Wajah" : "Daftar Wajah"}
              </Button>
            </div>
          </Card>
        </div>
      </div>
    </>
  )
}

/* --------------------------------- Enroll --------------------------------- */
function Enroll({
  go,
  user,
  onDone,
}: {
  go: (r: string) => void
  user?: AuthUser | null
  onDone: () => void
}) {
  const toast = useToast()
  const submitEnrollment = async (payload: Blob | Record<string, Blob>) => {
    if (payload instanceof Blob) return
    const form = new FormData()
    form.append("front_image", payload.front_image, "front.jpg")
    form.append("left_image", payload.left_image, "left.jpg")
    form.append("right_image", payload.right_image, "right.jpg")
    form.append("consent_agreed", "true")
    form.append("re_enroll", "true")
    await api.face.enroll(form)
  }
  return (
    <>
      <PageHeader
        title="Pendaftaran Wajah"
        subtitle="Ikuti panduan untuk mengambil beberapa sudut wajah Anda."
        breadcrumb={["Profil", "Pendaftaran Wajah"]}
      />
      <FaceFlow
        mode="enroll"
        studentName={user?.full_name || "Mahasiswa"}
        onSubmit={submitEnrollment}
        onComplete={() => {
          onDone()
          toast("Wajah berhasil didaftarkan.")
          go("profile")
        }}
        onCancel={() => go("profile")}
      />
    </>
  )
}
