import { App as AntApp, ConfigProvider } from 'antd'
import enGB from 'antd/locale/en_GB'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import { router } from './routes'
import { useSession } from './stores/session'

// Theme tokens are set once here; every page inherits them.
const theme = { token: { colorPrimary: '#3b5bdb', borderRadius: 8, fontFamily: 'system-ui, -apple-system, "Segoe UI", Roboto, sans-serif' } }

useSession.getState().load()

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ConfigProvider theme={theme} locale={enGB}>
      <AntApp>
        <RouterProvider router={router} />
      </AntApp>
    </ConfigProvider>
  </StrictMode>,
)
