import { useEffect, useRef, useState } from 'react'
import { Icon } from './icons'
import { Alert, Badge, Button, Card, StepIndicator, type Tone } from './ui'

const ENROLL_STEPS = ['Hadap depan', 'Hadap kiri', 'Hadap kanan']

type Phase =
  | 'intro'
  | 'permission'
  | 'ready'
  | 'detected'
  | 'processing'
  | 'success'
  | 'failed'
  | 'multiple'
  | 'noface'
  | 'nocamera'

const phaseMeta: Record<Phase, { label: string; tone: Tone }> = {
  intro: { label: 'Siap', tone: 'neutral' },
  permission: { label: 'Menunggu izin kamera', tone: 'warning' },
  ready: { label: 'Kamera aktif — mendeteksi', tone: 'info' },
  detected: { label: 'Wajah terdeteksi', tone: 'success' },
  processing: { label: 'Memverifikasi…', tone: 'info' },
  success: { label: 'Verifikasi berhasil', tone: 'success' },
  failed: { label: 'Verifikasi gagal', tone: 'error' },
  multiple: { label: 'Terdeteksi lebih dari satu wajah', tone: 'error' },
  noface: { label: 'Wajah belum terdeteksi', tone: 'warning' },
  nocamera: { label: 'Kamera tidak tersedia', tone: 'error' },
}

// Simulated camera preview — abstract silhouette so no real biometric data is implied.
function CameraCanvas({ phase }: { phase: Phase }) {
  const guideColor =
    phase === 'success'
      ? '#12b76a'
      : phase === 'failed' || phase === 'multiple'
        ? '#f04438'
        : phase === 'detected'
          ? '#12b76a'
          : phase === 'ready' || phase === 'processing'
            ? '#53b1fd'
            : '#98a2b3'

  const live = ['ready', 'detected', 'processing'].includes(phase)
  const off = ['intro', 'permission', 'nocamera'].includes(phase)

  return (
    <div className="relative aspect-[4/3] w-full rounded-xl overflow-hidden bg-[#0b1220]">
      {/* faux camera feed */}
      {!off && (
        <div
          className="absolute inset-0"
          style={{
            background: 'radial-gradient(120% 90% at 50% 30%, #24344d 0%, #131c2b 55%, #0b1220 100%)',
          }}
        >
          {/* silhouette */}
          <svg viewBox="0 0 200 150" className="absolute inset-0 w-full h-full opacity-70">
            <ellipse cx="100" cy="70" rx="34" ry="42" fill="#1f2c42" />
            <path d="M55 150 Q100 96 145 150 Z" fill="#1f2c42" />
            {phase === 'multiple' && (
              <>
                <ellipse cx="150" cy="80" rx="24" ry="30" fill="#28374f" opacity="0.9" />
              </>
            )}
          </svg>
        </div>
      )}

      {off && (
        <div className="absolute inset-0 flex flex-col items-center justify-center text-white/40 gap-2">
          <Icon.camera width={40} height={40} />
          <span className="text-xs">{phase === 'nocamera' ? 'Kamera tidak terdeteksi' : 'Pratinjau kamera nonaktif'}</span>
        </div>
      )}

      {/* face guide oval */}
      {!off && (
        <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
          <div
            className="relative rounded-[50%] transition-colors duration-300"
            style={{
              width: '46%',
              height: '78%',
              border: `2.5px dashed ${guideColor}`,
              boxShadow: `0 0 0 9999px rgba(11,18,32,0.45)`,
            }}
          >
            {/* corner ticks */}
            {['-top-px -left-px', '-top-px -right-px', '-bottom-px -left-px', '-bottom-px -right-px'].map((c) => (
              <span key={c} className={`absolute ${c} w-4 h-4`} style={{ borderColor: guideColor }} />
            ))}
          </div>
        </div>
      )}

      {/* scan line while processing/ready */}
      {live && (
        <div
          className="absolute left-[27%] right-[27%] h-0.5 animate-scan"
          style={{ background: `linear-gradient(90deg,transparent,${guideColor},transparent)` }}
        />
      )}

      {/* live indicator */}
      {live && (
        <div className="absolute top-3 left-3 flex items-center gap-1.5 rounded-full bg-black/50 px-2.5 py-1 text-[11px] font-medium text-white">
          <span className="w-1.5 h-1.5 rounded-full bg-[#f04438] animate-pulse" /> LIVE
        </div>
      )}

      {/* instruction pill */}
      {!off && (
        <div className="absolute bottom-3 left-1/2 -translate-x-1/2 rounded-full bg-black/55 px-3 py-1.5 text-[12px] text-white/90 backdrop-blur whitespace-nowrap">
          {phase === 'multiple'
            ? 'Pastikan hanya satu wajah di depan kamera'
            : phase === 'noface'
              ? 'Wajah belum terdeteksi'
              : phase === 'success'
                ? 'Identitas sesuai'
                : 'Posisikan wajah di dalam area'}
        </div>
      )}

      {processingOverlay(phase)}
    </div>
  )
}

function processingOverlay(phase: Phase) {
  if (phase !== 'processing') return null
  return (
    <div className="absolute inset-0 flex items-center justify-center bg-black/25">
      <span className="w-10 h-10 rounded-full border-[3px] border-white/40 border-t-white animate-spin-slow" />
    </div>
  )
}

export function FaceFlow({
  mode,
  studentName,
  sessionName,
  onComplete,
  onCancel,
}: {
  mode: 'verify' | 'enroll'
  studentName: string
  sessionName?: string
  onComplete: () => void
  onCancel: () => void
}) {
  const [phase, setPhase] = useState<Phase>('intro')
  const [cap, setCap] = useState(0) // langkah pengambilan (enrollment)
  const timers = useRef<number[]>([])
  const clear = () => {
    timers.current.forEach(clearTimeout)
    timers.current = []
  }
  useEffect(() => () => clear(), [])

  const meta = phaseMeta[phase]
  const isVerify = mode === 'verify'

  const grant = () => {
    setPhase('permission')
    timers.current.push(
      window.setTimeout(() => setPhase('ready'), 900),
      window.setTimeout(() => setPhase('detected'), 2600),
    )
  }

  const finish = () => {
    clear()
    setPhase('processing')
    timers.current.push(window.setTimeout(() => setPhase('success'), 2200))
  }

  // Verifikasi: sekali proses. Enrollment: kumpulkan beberapa sudut wajah.
  const capture = () => {
    if (isVerify) return finish()
    if (cap < ENROLL_STEPS.length - 1) {
      clear()
      setCap((c) => c + 1)
      setPhase('ready')
      timers.current.push(window.setTimeout(() => setPhase('detected'), 1500))
    } else {
      finish()
    }
  }

  const simulate = (p: Phase) => {
    clear()
    setPhase(p)
  }

  const retry = () => {
    clear()
    setPhase('ready')
    timers.current.push(window.setTimeout(() => setPhase('detected'), 1800))
  }

  const title = isVerify ? 'Verifikasi Wajah' : 'Pendaftaran Wajah'

  return (
    <div className="grid lg:grid-cols-[1.4fr_1fr] gap-6 items-start">
      <Card padded={false} className="overflow-hidden">
        <div className="p-5 border-b border-[var(--color-border)] flex items-center justify-between">
          <div>
            <h3 className="font-semibold text-[var(--color-text)]">{title}</h3>
            <p className="text-sm text-[var(--color-text-2)]">
              {isVerify
                ? 'Verifikasi identitas untuk pencatatan presensi.'
                : 'Daftarkan wajah Anda untuk digunakan pada verifikasi presensi.'}
            </p>
          </div>
          <Badge tone={meta.tone}>{meta.label}</Badge>
        </div>
        <div className="p-5">
          {!isVerify && phase !== 'intro' && (
            <div className="mb-4">
              <StepIndicator steps={ENROLL_STEPS} current={phase === 'success' ? ENROLL_STEPS.length : cap} />
              {['ready', 'detected'].includes(phase) && (
                <p className="mt-2 text-sm text-[var(--color-text-2)]">
                  Langkah {cap + 1} dari {ENROLL_STEPS.length} — <span className="font-medium text-[var(--color-text)]">{ENROLL_STEPS[cap]}</span>
                </p>
              )}
            </div>
          )}

          <CameraCanvas phase={phase} />

          {/* controls */}
          <div className="mt-5 flex flex-wrap items-center gap-2">
            {phase === 'intro' && (
              <Button size="lg" icon={<Icon.camera width={18} height={18} />} onClick={grant}>
                Aktifkan Kamera
              </Button>
            )}
            {(phase === 'ready' || phase === 'detected') && (
              <Button size="lg" disabled={phase !== 'detected'} onClick={capture} icon={<Icon.checkCircle width={18} height={18} />}>
                {isVerify
                  ? 'Verifikasi Sekarang'
                  : cap < ENROLL_STEPS.length - 1
                    ? `Ambil (${cap + 1}/${ENROLL_STEPS.length})`
                    : 'Ambil & Daftarkan'}
              </Button>
            )}
            {(phase === 'failed' || phase === 'multiple' || phase === 'noface') && (
              <Button size="lg" onClick={retry} icon={<Icon.history width={18} height={18} />}>
                Coba Lagi
              </Button>
            )}
            {phase === 'nocamera' && (
              <Button size="lg" onClick={() => simulate('intro')} icon={<Icon.history width={18} height={18} />}>
                Muat Ulang Kamera
              </Button>
            )}
            {phase === 'success' && (
              <Button size="lg" variant="success" onClick={onComplete} icon={<Icon.arrowRight width={18} height={18} />}>
                {isVerify ? 'Lanjutkan ke Presensi' : 'Selesai'}
              </Button>
            )}
            <Button variant="ghost" size="lg" onClick={onCancel}>
              Batal
            </Button>

            {/* demo: simulate edge-case states */}
            {['ready', 'detected'].includes(phase) && (
              <div className="ml-auto flex items-center gap-1">
                <span className="text-[11px] text-[var(--color-text-3)] mr-1">Simulasi:</span>
                <Button variant="secondary" size="sm" onClick={() => simulate('multiple')}>
                  Banyak wajah
                </Button>
                <Button variant="secondary" size="sm" onClick={() => simulate('failed')}>
                  Gagal
                </Button>
                <Button variant="secondary" size="sm" onClick={() => simulate('nocamera')}>
                  Kamera mati
                </Button>
              </div>
            )}
          </div>
        </div>
      </Card>

      {/* side info */}
      <div className="space-y-4">
        <Card>
          <h4 className="text-sm font-semibold text-[var(--color-text)] mb-3">Informasi</h4>
          <dl className="space-y-3 text-sm">
            <div className="flex justify-between gap-3">
              <dt className="text-[var(--color-text-2)]">Mahasiswa</dt>
              <dd className="font-medium text-right">{studentName}</dd>
            </div>
            {sessionName && (
              <div className="flex justify-between gap-3">
                <dt className="text-[var(--color-text-2)]">Sesi</dt>
                <dd className="font-medium text-right">{sessionName}</dd>
              </div>
            )}
            <div className="flex justify-between gap-3">
              <dt className="text-[var(--color-text-2)]">Status</dt>
              <dd>
                <Badge tone={meta.tone}>{meta.label}</Badge>
              </dd>
            </div>
          </dl>
        </Card>

        {phase === 'permission' && (
          <Alert tone="warning" title="Izin kamera diperlukan">
            Akses kamera diperlukan untuk verifikasi wajah. Izinkan penggunaan kamera pada peramban Anda.
          </Alert>
        )}
        {phase === 'failed' && (
          <Alert tone="error" title="Verifikasi wajah gagal">
            Pastikan wajah terlihat jelas dan pencahayaan cukup, lalu coba lagi.
          </Alert>
        )}
        {phase === 'multiple' && (
          <Alert tone="error" title="Terdeteksi lebih dari satu wajah">
            Pastikan hanya satu wajah berada di depan kamera.
          </Alert>
        )}
        {phase === 'noface' && (
          <Alert tone="warning" title="Wajah belum terdeteksi">
            Posisikan wajah Anda di dalam area panduan.
          </Alert>
        )}
        {phase === 'nocamera' && (
          <Alert tone="error" title="Kamera tidak tersedia">
            Kamera tidak terdeteksi. Periksa perangkat kamera atau izin peramban.
          </Alert>
        )}
        {phase === 'success' && (
          <Alert tone="success" title={isVerify ? 'Identitas terverifikasi' : 'Wajah berhasil didaftarkan'}>
            {isVerify
              ? 'Identitas Anda sesuai. Anda dapat melanjutkan pencatatan presensi.'
              : 'Templat wajah Anda telah terdaftar dan siap digunakan untuk verifikasi presensi.'}
          </Alert>
        )}

        <Card className="bg-[var(--color-surface-2)]">
          <div className="flex gap-2.5">
            <Icon.lock width={18} height={18} className="text-[var(--color-text-2)] shrink-0 mt-0.5" />
            <p className="text-xs text-[var(--color-text-2)] leading-relaxed">
              Verifikasi wajah digunakan untuk membantu memastikan kesesuaian identitas pengguna yang sedang login. Sistem
              membandingkan wajah dengan templat milik akun Anda sendiri — bukan mencari identitas. Data biometrik tidak
              ditampilkan.
            </p>
          </div>
        </Card>
      </div>
    </div>
  )
}
