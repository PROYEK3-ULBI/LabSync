import { useEffect, useState } from "react"
import { AdminApp } from "./apps/AdminApp"
import { StudentApp } from "./apps/StudentApp"
import { Login } from "./components/Login"
import { ToastProvider } from "./components/ui"
import { api, type AuthUser } from "./lib/api"

type Role = "student" | "admin"
type Screen = {
  role: Role
  authed: boolean
  user: AuthUser | null
}

export default function App() {
  const [screen, setScreen] = useState<Screen>({
    role: "student",
    authed: false,
    user: null,
  })
  const [checking, setChecking] = useState(true)

  // Session restoration saat aplikasi dibuka
  useEffect(() => {
    let active = true
    async function restoreSession() {
      try {
        const user = await api.auth.me()
        if (active && user) {
          const role: Role = user.role === "ADMIN" ? "admin" : "student"
          setScreen({ role, authed: true, user })
        }
      } catch {
        // Belum terautentikasi atau cookie tidak valid/expired
        if (active) {
          setScreen((s) => ({ ...s, authed: false, user: null }))
        }
      } finally {
        if (active) setChecking(false)
      }
    }
    restoreSession()
    return () => {
      active = false
    }
  }, [])

  const handleLogin = (user: AuthUser) => {
    const role: Role = user.role === "ADMIN" ? "admin" : "student"
    setScreen({ role, authed: true, user })
  }

  const handleLogout = async () => {
    try {
      await api.auth.logout()
    } catch {
      // Abaikan kegagalan jaringan saat logout
    } finally {
      setScreen((s) => ({ ...s, authed: false, user: null }))
    }
  }

  if (checking) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-[var(--color-bg)]">
        <div className="flex flex-col items-center gap-3">
          <div className="w-8 h-8 border-3 border-[var(--color-primary)] border-t-transparent rounded-full animate-spin" />
          <p className="text-sm font-medium text-[var(--color-text-2)]">
            Memeriksa sesi pengguna…
          </p>
        </div>
      </div>
    )
  }

  return (
    <ToastProvider>
      {!screen.authed ? (
        <Login
          role={screen.role}
          onLogin={handleLogin}
          onSwitchRole={() =>
            setScreen((s) => ({
              role: s.role === "student" ? "admin" : "student",
              authed: false,
              user: null,
            }))
          }
        />
      ) : screen.role === "student" ? (
        <StudentApp onLogout={handleLogout} user={screen.user} />
      ) : (
        <AdminApp onLogout={handleLogout} user={screen.user} />
      )}
    </ToastProvider>
  )
}
