import { create } from 'zustand'
import { authApi, type Me, type Role } from '../requests/api'
import { setOnUnauthenticated } from '../requests/client'

interface Session {
  me: Me | null
  loading: boolean // true until the first /auth/me answer
  load: () => Promise<void>
  login: (email: string, password: string) => Promise<Me>
  setMe: (me: Me) => void
  logout: () => Promise<void>
}

export const useSession = create<Session>((set) => ({
  me: null,
  loading: true,
  load: async () => {
    try {
      set({ me: await authApi.me(), loading: false })
    } catch {
      set({ me: null, loading: false })
    }
  },
  login: async (email, password) => {
    const me = await authApi.login(email, password)
    set({ me })
    return me
  },
  setMe: (me) => set({ me }),
  logout: async () => {
    try {
      await authApi.logout()
    } finally {
      set({ me: null })
    }
  },
}))

// An expired or revoked session anywhere sends the user back to the login screen (via the route guard).
setOnUnauthenticated(() => useSession.setState({ me: null }))

export function homeFor(role: Role): string {
  return role === 'customer' ? '/portal' : '/customers'
}
