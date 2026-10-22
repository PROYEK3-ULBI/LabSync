import { useEffect, useRef, useState } from "react"
import { Icon } from "./icons"
import { Alert, Badge, Button, Card, StepIndicator, type Tone } from "./ui"

const ENROLL_STEPS = [
  { key: "front_image", label: "Hadap depan" },
  { key: "left_image", label: "Hadap kiri" },
  { key: "right_image", label: "Hadap kanan" },
] as const

type Phase = "intro" | "ready" | "processing" | "success" | "failed" | "nocamera"

type PhaseMeta = {
  label: string
  tone: Tone
}

type FlowError = {
  title: string
  message: string
}

type SubmitError = {
  code?: string
  message?: string
}

const phaseMeta: Record<Phase, PhaseMeta> = {
  intro: { label: "Siap", tone: "neutral" },
  ready: { label: "Kamera aktif", tone: "info" },
  processing: { label: "Mengirim...", tone: "info" },
  success: { label: "Berhasil", tone: "success" },
  failed: { label: "Gagal", tone: "error" },
  nocamera: { label: "Kamera bermasalah", tone: "error" },
}

function errorTitle(code?: string) {
  if (code === "CAMERA_PERMISSION") return "Izin kamera ditolak"
  if (code === "NO_FACE_DETECTED" || code === "MULTIPLE_FACES_DETECTED") {
    return "Gambar wajah tidak valid"
  }
  if (code === "VERIFICATION_FAILED") return "Verifikasi wajah gagal"
  if (code === "NETWORK_ERROR" || code === "SERVICE_UNAVAILABLE") {
    return "Server belum tersedia"
  }
  return "Proses gagal"
}

export function FaceFlow({
  mode,
  studentName,
  sessionName,
  onComplete,
  onCancel,
  onSubmit,
}: {
  mode: "verify" | "enroll"
  studentName: string
  sessionName?: string
  onComplete: () => void
  onCancel: () => void
  onSubmit: (payload: Blob | Record<string, Blob>) => Promise<void>
}) {
  const [phase, setPhase] = useState<Phase>("intro")
  const [cap, setCap] = useState(0)
  const [samples, setSamples] = useState<Partial<Record<string, Blob>>>({})
  const [error, setError] = useState<FlowError | null>(null)
  const streamRef = useRef<MediaStream | null>(null)
  const videoRef = useRef<HTMLVideoElement | null>(null)
  const canvasRef = useRef<HTMLCanvasElement | null>(null)
  const isVerify = mode === "verify"
  const meta = phaseMeta[phase]

  const stopCamera = () => {
    streamRef.current?.getTracks().forEach((track) => track.stop())
    streamRef.current = null
  }

  useEffect(() => stopCamera, [])

  const startCamera = async () => {
    setError(null)
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: "user" },
        audio: false,
      })
      streamRef.current = stream
      if (videoRef.current) videoRef.current.srcObject = stream
      setPhase("ready")
    } catch {
      setPhase("nocamera")
      setError({
        title: errorTitle("CAMERA_PERMISSION"),
        message:
          "Browser menolak atau tidak menemukan kamera. Izinkan kamera lalu coba lagi.",
      })
    }
  }

  const captureFrame = () =>
    new Promise<Blob>((resolve, reject) => {
      const video = videoRef.current
      const canvas = canvasRef.current
      if (!video || !canvas || video.videoWidth === 0) {
        reject(new Error("Kamera belum siap."))
        return
      }
      canvas.width = video.videoWidth
      canvas.height = video.videoHeight
      canvas.getContext("2d")?.drawImage(video, 0, 0)
      canvas.toBlob(
        (blob) =>
          blob ? resolve(blob) : reject(new Error("Gagal mengambil gambar.")),
        "image/jpeg",
        0.9,
      )
    })

  const submit = async (payload: Blob | Record<string, Blob>) => {
    setPhase("processing")
    setError(null)
    try {
      await onSubmit(payload)
      stopCamera()
      setPhase("success")
    } catch (err) {
      const maybe = err as SubmitError
      setPhase("failed")
      setError({
        title: errorTitle(maybe.code),
        message:
          maybe.message || "Permintaan gagal diproses. Silakan coba lagi.",
      })
    }
  }

  const capture = async () => {
    if (phase === "processing") return
    try {
      const blob = await captureFrame()
      if (isVerify) {
        await submit(blob)
        return
      }
      const key = ENROLL_STEPS[cap].key
      const next = { ...samples, [key]: blob }
      setSamples(next)
      if (cap < ENROLL_STEPS.length - 1) {
        setCap((c) => c + 1)
        return
      }
      await submit(next as Record<string, Blob>)
    } catch (err) {
      setPhase("failed")
      setError({
        title: "Gambar tidak bisa diambil",
        message: err instanceof Error ? err.message : "Kamera belum siap.",
      })
    }
  }

  return (
    <div className="grid lg:grid-cols-[1.4fr_1fr] gap-6 items-start">
      <Card padded={false} className="overflow-hidden">
        <div className="p-5 border-b border-[var(--color-border)] flex items-center justify-between">
          <div>
            <h3 className="font-semibold text-[var(--color-text)]">
              {isVerify ? "Verifikasi Wajah" : "Pendaftaran Wajah"}
            </h3>
            <p className="text-sm text-[var(--color-text-2)]">
              {isVerify
                ? "Kirim satu gambar live ke backend untuk presensi."
                : "Ambil FRONT, LEFT, RIGHT lalu kirim sekali ke backend."}
            </p>
          </div>
          <Badge tone={meta.tone}>{meta.label}</Badge>
        </div>
        <div className="p-5">
          {!isVerify && phase !== "intro" && (
            <div className="mb-4">
              <StepIndicator
                steps={ENROLL_STEPS.map((step) => step.label)}
                current={phase === "success" ? ENROLL_STEPS.length : cap}
              />
              <p className="mt-2 text-sm text-[var(--color-text-2)]">
                Langkah {cap + 1} dari {ENROLL_STEPS.length} -{" "}
                <span className="font-medium text-[var(--color-text)]">
                  {ENROLL_STEPS[cap]?.label || "Selesai"}
                </span>
              </p>
            </div>
          )}

          <div className="relative aspect-[4/3] w-full rounded-xl overflow-hidden bg-[#0b1220]">
            <video
              ref={videoRef}
              autoPlay
              playsInline
              muted
              className={
                phase === "ready" || phase === "processing"
                  ? "h-full w-full object-cover"
                  : "hidden"
              }
            />
            {phase !== "ready" && phase !== "processing" && (
              <div className="absolute inset-0 flex flex-col items-center justify-center text-white/45 gap-2">
                <Icon.camera width={40} height={40} />
                <span className="text-xs">Pratinjau kamera nonaktif</span>
              </div>
            )}
            {(phase === "ready" || phase === "processing") && (
              <>
                <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
                  <div className="rounded-[50%] border-2 border-dashed border-white/80 w-[46%] h-[78%] shadow-[0_0_0_9999px_rgba(11,18,32,0.35)]" />
                </div>
                <div className="absolute top-3 left-3 flex items-center gap-1.5 rounded-full bg-black/50 px-2.5 py-1 text-[11px] font-medium text-white">
                  <span className="w-1.5 h-1.5 rounded-full bg-[#f04438] animate-pulse" />
                  LIVE
                </div>
              </>
            )}
            {phase === "processing" && (
              <div className="absolute inset-0 flex items-center justify-center bg-black/25">
                <span className="w-10 h-10 rounded-full border-[3px] border-white/40 border-t-white animate-spin-slow" />
              </div>
            )}
          </div>
          <canvas ref={canvasRef} className="hidden" />

          <div className="mt-5 flex flex-wrap items-center gap-2">
            {phase === "intro" || phase === "nocamera" ? (
              <Button
                size="lg"
                icon={<Icon.camera width={18} />}
                onClick={startCamera}
              >
                Aktifkan Kamera
              </Button>
            ) : null}
            {phase === "ready" && (
              <Button
                size="lg"
                onClick={capture}
                icon={<Icon.checkCircle width={18} />}
              >
                {isVerify
                  ? "Verifikasi Sekarang"
                  : cap < ENROLL_STEPS.length - 1
                    ? `Ambil (${cap + 1}/${ENROLL_STEPS.length})`
                    : "Ambil & Daftarkan"}
              </Button>
            )}
            {phase === "failed" && (
              <Button
                size="lg"
                onClick={startCamera}
                icon={<Icon.history width={18} />}
              >
                Coba Lagi
              </Button>
            )}
            {phase === "success" && (
              <Button
                size="lg"
                variant="success"
                onClick={onComplete}
                icon={<Icon.arrowRight width={18} />}
              >
                {isVerify ? "Lanjutkan" : "Selesai"}
              </Button>
            )}
            <Button
              variant="ghost"
              size="lg"
              onClick={onCancel}
              disabled={phase === "processing"}
            >
              Batal
            </Button>
          </div>
        </div>
      </Card>

      <div className="space-y-4">
        <Card>
          <h4 className="text-sm font-semibold text-[var(--color-text)] mb-3">
            Informasi
          </h4>
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

        {error && (
          <Alert tone="error" title={error.title}>
            {error.message}
          </Alert>
        )}
        {phase === "success" && (
          <Alert
            tone="success"
            title={
              isVerify ? "Presensi dikonfirmasi" : "Wajah berhasil didaftarkan"
            }
          >
            Backend telah mengonfirmasi proses berhasil.
          </Alert>
        )}
        <Card className="bg-[var(--color-surface-2)]">
          <div className="flex gap-2.5">
            <Icon.lock
              width={18}
              height={18}
              className="text-[var(--color-text-2)] shrink-0 mt-0.5"
            />
            <p className="text-xs text-[var(--color-text-2)] leading-relaxed">
              Browser hanya mengambil frame saat tombol ditekan. Gambar dikirim
              ke API Go, tidak disimpan permanen di frontend, dan tidak dikirim
              ke face service langsung.
            </p>
          </div>
        </Card>
      </div>
    </div>
  )
}
