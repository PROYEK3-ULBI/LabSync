import { useState } from 'react'
import { AdminApp } from './apps/AdminApp'
import { StudentApp } from './apps/StudentApp'
import { Login } from './components/Login'
import { ToastProvider } from './components/ui'

type Role = 'student' | 'admin'
type Screen = { role: Role; authed: boolean }

export default function App() {
  const [screen, setScreen] = useState<Screen>({ role: 'student', authed: false })

  return (
    <ToastProvider>
      {!screen.authed ? (
        <Login
          role={screen.role}
          onLogin={() => setScreen((s) => ({ ...s, authed: true }))}
          onSwitchRole={() => setScreen((s) => ({ role: s.role === 'student' ? 'admin' : 'student', authed: false }))}
        />
      ) : screen.role === 'student' ? (
        <StudentApp onLogout={() => setScreen({ role: 'student', authed: false })} />
      ) : (
        <AdminApp onLogout={() => setScreen({ role: 'admin', authed: false })} />
      )}
    </ToastProvider>
  )
}
