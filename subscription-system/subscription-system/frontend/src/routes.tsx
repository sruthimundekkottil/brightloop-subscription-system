import { Result } from 'antd'
import { Link, Navigate, createBrowserRouter } from 'react-router-dom'
import { AppLayout } from './components/AppLayout'
import { GuestOnly, RequireRole } from './components/RequireRole'
import { AuditLogPage } from './pages/admin/AuditLogPage'
import { SettingsPage } from './pages/admin/SettingsPage'
import { StaffPage } from './pages/admin/StaffPage'
import { ForgotPasswordPage, LoginPage, RegisterPage, ResetPasswordPage } from './pages/auth/AuthPages'
import { CustomerDetailPage, CustomersPage } from './pages/customers/CustomerPages'
import { MyInvoicesPage, MyProfilePage, MySubscriptionsPage } from './pages/portal/PortalPages'
import { useSession, homeFor } from './stores/session'

function Home() {
  const me = useSession((s) => s.me)
  return <Navigate to={me ? homeFor(me.role) : '/login'} replace />
}

// Role guards mirror the API's route groups (backend/internal/cmd/routes.go).
export const router = createBrowserRouter([
  {
    element: <GuestOnly />,
    children: [
      { path: '/login', element: <LoginPage /> },
      { path: '/register', element: <RegisterPage /> },
    ],
  },
  { path: '/forgot-password', element: <ForgotPasswordPage /> },
  { path: '/reset-password', element: <ResetPasswordPage /> },
  {
    element: <RequireRole roles={['admin', 'staff', 'customer']} />,
    children: [
      { path: '/', element: <Home /> },
      {
        element: <AppLayout />,
        children: [
          {
            element: <RequireRole roles={['admin']} />,
            children: [
              { path: '/staff', element: <StaffPage /> },
              { path: '/audit-log', element: <AuditLogPage /> },
            ],
          },
          {
            element: <RequireRole roles={['admin', 'staff']} />,
            children: [
              { path: '/customers', element: <CustomersPage /> },
              { path: '/customers/:id', element: <CustomerDetailPage /> },
              { path: '/settings', element: <SettingsPage /> },
            ],
          },
          {
            element: <RequireRole roles={['customer']} />,
            children: [
              { path: '/portal', element: <MySubscriptionsPage /> },
              { path: '/portal/invoices', element: <MyInvoicesPage /> },
              { path: '/portal/profile', element: <MyProfilePage /> },
            ],
          },
        ],
      },
    ],
  },
  { path: '*', element: <Result status="404" title="Page not found" extra={<Link to="/">Go home</Link>} /> },
])
