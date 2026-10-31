import { Spin } from 'antd'
import { Navigate, Outlet, useLocation } from 'react-router-dom'
import type { Role } from '../requests/api'
import { homeFor, useSession } from '../stores/session'

// Route guard: no session → login (remembering where the user was going); wrong role → that role's home.
// The server enforces the same rules; this only keeps users off screens they cannot use.
export function RequireRole({ roles }: { roles: Role[] }) {
  const { me, loading } = useSession()
  const location = useLocation()
  if (loading) return <Spin fullscreen />
  if (!me) return <Navigate to="/login" replace state={{ from: location.pathname + location.search }} />
  if (!roles.includes(me.role)) return <Navigate to={homeFor(me.role)} replace />
  return <Outlet />
}

// GuestOnly keeps logged-in users off the login and register screens.
export function GuestOnly() {
  const { me, loading } = useSession()
  if (loading) return <Spin fullscreen />
  return me ? <Navigate to={homeFor(me.role)} replace /> : <Outlet />
}
