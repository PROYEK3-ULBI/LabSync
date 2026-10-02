import { useState } from 'react'
import { Icon } from './icons'
import { Alert, Button, Field, Input } from './ui'

export function Login({
  role,
  onLogin,
  onSwitchRole,
}: {
  role: 'student' | 'admin'
  onLogin: () => void
  onSwitchRole: () => void
}) {
  const isStudent = role === 'student'
  const [email, setEmail] = useState(isStudent ? 'ahmad.fauzan@student.univ.ac.id' : 'dewi.laboran@univ.ac.id')
  const [pw, setPw] = useState('')
  const [show, setShow] = useState(false)
  const [loading, setLoading] = useState(false)
  const [err, setErr] = useState<string | null>(null)
  const [touched, setTouched] = useState(false)

  const submit = (e: React.FormEvent) => {
    e.preventDefault()
    setTouched(true)
    setErr(null)
    if (!email || !pw) return
    setLoading(true)
    setTimeout(() => {
      setLoading(false)
      // demo: password "demo" fails to show error state, anything else succeeds
      if (pw.toLowerCase() === 'salah') {
        setErr('Email atau kata sandi tidak sesuai. Silakan periksa kembali.')
        return
      }
      onLogin()
    }, 1100)
  }

  const accent = isStudent ? 'var(--color-primary)' : '#1d2939'

  return (
    <div className="min-h-screen grid lg:grid-cols-2">
      {/* Brand panel */}
      <div
        className="hidden lg:flex flex-col justify-between p-12 text-white relative overflow-hidden"
        style={{ background: isStudent ? 'linear-gradient(150deg,#1e40af,#1d4ed8 55%,#2563eb)' : 'linear-gradient(150deg,#0f172a,#1d2939 60%,#334155)' }}
      >
        <div className="relative z-10 flex items-center gap-3">
          <span className="inline-flex items-center justify-center w-11 h-11 rounded-xl bg-white/15 backdrop-blur">
            <Icon.shield width={24} height={24} />
          </span>
          <span className="font-semibold text-lg">LabPresensi</span>
        </div>

        <div className="relative z-10 max-w-md">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-white/60 mb-4">
            {isStudent ? 'Portal Mahasiswa' : 'Portal Admin / Laboran'}
          </p>
          <h1 className="text-3xl font-semibold leading-tight tracking-tight">
            Presensi &amp; Manajemen Aset Laboratorium Komputer
          </h1>
          <p className="mt-4 text-white/70 leading-relaxed">
            {isStudent
              ? 'Verifikasi wajah untuk presensi praktikum, pantau sesi, dan kelola aset yang dialokasikan kepada Anda.'
              : 'Kelola pengguna, sesi praktikum, presensi, alokasi aset, pemeliharaan, dan laporan operasional laboratorium.'}
          </p>
          <div className="mt-8 space-y-3">
            {[
              isStudent ? 'Verifikasi wajah 1:1 melalui kamera peramban' : 'Pemantauan presensi & aset secara real-time',
              isStudent ? 'Riwayat presensi & alokasi aset' : 'Alur alokasi, pengembalian, dan pemeliharaan aset',
              'Antarmuka aman, sesuai kebutuhan akademik',
            ].map((t) => (
              <div key={t} className="flex items-center gap-2.5 text-sm text-white/85">
                <Icon.checkCircle width={18} height={18} className="text-white/80" /> {t}
              </div>
            ))}
          </div>
        </div>

        <p className="relative z-10 text-xs text-white/40">
          Data yang ditampilkan bersifat fiktif untuk keperluan demonstrasi antarmuka.
        </p>
        <div className="absolute -right-24 -bottom-24 w-96 h-96 rounded-full bg-white/5" />
        <div className="absolute right-16 top-16 w-40 h-40 rounded-full bg-white/5" />
      </div>

      {/* Form panel */}
      <div className="flex items-center justify-center p-6 sm:p-12 bg-[var(--color-bg)]">
        <div className="w-full max-w-sm">
          <div className="lg:hidden flex items-center gap-2.5 mb-8">
            <span className="inline-flex items-center justify-center w-10 h-10 rounded-lg bg-[var(--color-primary)] text-white">
              <Icon.shield width={20} height={20} />
            </span>
            <span className="font-semibold text-lg">LabPresensi</span>
          </div>

          <span
            className="inline-flex items-center gap-1.5 text-[11px] font-semibold uppercase tracking-wide px-2.5 py-1 rounded-md mb-3"
            style={{ background: isStudent ? 'var(--color-primary-soft)' : '#eef2f6', color: accent }}
          >
            {isStudent ? 'Portal Mahasiswa' : 'Portal Admin'}
          </span>
          <h2 className="text-2xl font-semibold tracking-tight text-[var(--color-text)]">Selamat datang kembali</h2>
          <p className="text-sm text-[var(--color-text-2)] mt-1.5">
            Masuk untuk melanjutkan ke sistem presensi dan manajemen aset.
          </p>

          {err && (
            <div className="mt-5">
              <Alert tone="error" title="Gagal masuk">
                {err}
              </Alert>
            </div>
          )}

          <form onSubmit={submit} className="mt-5 space-y-4">
            <Field
              label={isStudent ? 'Email / NIM' : 'Email'}
              required
              error={touched && !email ? 'Wajib diisi.' : undefined}
            >
              <Input
                type="text"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                invalid={touched && !email}
                placeholder="nama@univ.ac.id"
                autoComplete="username"
              />
            </Field>

            <Field label="Kata Sandi" required error={touched && !pw ? 'Wajib diisi.' : undefined}>
              <div className="relative">
                <Input
                  type={show ? 'text' : 'password'}
                  value={pw}
                  onChange={(e) => setPw(e.target.value)}
                  invalid={touched && !pw}
                  placeholder="Masukkan kata sandi"
                  autoComplete="current-password"
                  className="pr-10"
                />
                <button
                  type="button"
                  onClick={() => setShow((v) => !v)}
                  aria-label={show ? 'Sembunyikan kata sandi' : 'Tampilkan kata sandi'}
                  className="absolute right-2.5 top-1/2 -translate-y-1/2 text-[var(--color-text-3)] hover:text-[var(--color-text)]"
                >
                  {show ? <Icon.eyeOff width={18} height={18} /> : <Icon.eye width={18} height={18} />}
                </button>
              </div>
            </Field>

            <Button type="submit" size="lg" loading={loading} className="w-full">
              {loading ? 'Memverifikasi…' : 'Masuk'}
            </Button>
          </form>

          <p className="mt-4 text-xs text-[var(--color-text-3)] leading-relaxed">
            Petunjuk demo: gunakan kata sandi apa pun untuk masuk, atau ketik <span className="font-mono">salah</span> untuk
            melihat status kesalahan kredensial.
          </p>

          <div className="mt-6 pt-5 border-t border-[var(--color-border)] flex items-center justify-between text-sm">
            <span className="text-[var(--color-text-2)]">
              {isStudent ? 'Anda seorang admin/laboran?' : 'Anda seorang mahasiswa?'}
            </span>
            <button onClick={onSwitchRole} className="font-medium text-[var(--color-primary)] hover:underline flex items-center gap-1">
              {isStudent ? 'Portal Admin' : 'Portal Mahasiswa'} <Icon.arrowRight width={14} height={14} />
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
