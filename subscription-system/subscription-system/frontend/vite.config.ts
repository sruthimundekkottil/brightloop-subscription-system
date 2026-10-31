import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

// The dev server proxies /api, so the session cookie is same-origin. In docker-compose the target is http://api:8000.
export default defineConfig({
  plugins: [react()],
  server: { port: 5173, host: true, proxy: { '/api': process.env.VITE_API_TARGET ?? 'http://localhost:8000' } },
})
