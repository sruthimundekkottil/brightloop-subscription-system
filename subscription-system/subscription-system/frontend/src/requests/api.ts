// Typed calls for member A's endpoints. Other members add their own files next to this one.
import { get, list, send } from './client'

export type Role = 'admin' | 'staff' | 'customer'

export interface Me {
  id: number
  email: string
  role: Role
  customerId: number | null
}

export interface User {
  id: number
  email: string
  role: Role
  isActive: boolean
  fullName: string | null
  jobTitle: string | null
  customerId: number | null
  createdAt: string
}

export interface Customer {
  id: number
  userId: number | null
  email: string | null
  name: string
  companyName: string | null
  phone: string | null
  billingAddress: string | null
  createdAt: string
}
export type CustomerFields = Partial<Pick<Customer, 'name' | 'companyName' | 'phone' | 'billingAddress'>>

export interface Settings {
  companyName: string
  address: string | null
  currencyCode: string
  invoicePrefix: string
  defaultQuotationValidityDays: number
  autoConfirmRecurringInvoices: boolean
  testClockDate: string | null
  today: string
  devMode: boolean
}

export interface AuditEntry {
  id: number
  tableName: string
  recordId: number
  action: string
  changedBy: number | null
  changedByEmail: string | null
  reason: string | null
  changes: unknown
  createdAt: string
}

export interface PageParams {
  page?: number
  pageSize?: number
}

export const authApi = {
  me: () => get<Me>('/auth/me'),
  login: (email: string, password: string) => send<Me>('post', '/auth/login', { email, password }),
  register: (b: { email: string; password: string; name: string; companyName?: string; phone?: string }) => send<Me>('post', '/auth/register', b),
  logout: () => send('post', '/auth/logout'),
  forgotPassword: (email: string) => send<{ message: string }>('post', '/auth/forgot-password', { email }),
  resetPassword: (token: string, password: string) => send<{ message: string }>('post', '/auth/reset-password', { token, password }),
}

export const userApi = {
  list: (p: PageParams & { role?: Role; active?: string; q?: string }) => list<User>('/users', p),
  createStaff: (b: { email: string; password: string; fullName: string; jobTitle?: string }) => send<User>('post', '/staff', b),
  setStatus: (id: number, isActive: boolean) => send<User>('patch', `/users/${id}/status`, { isActive }),
}

export const customerApi = {
  list: (p: PageParams & { q?: string }) => list<Customer>('/customers', p),
  get: (id: number) => get<Customer>(`/customers/${id}`),
  create: (b: CustomerFields) => send<Customer>('post', '/customers', b),
  update: (id: number, b: CustomerFields) => send<Customer>('patch', `/customers/${id}`, b),
  myProfile: () => get<Customer>('/portal/profile'),
  updateMyProfile: (b: CustomerFields) => send<Customer>('patch', '/portal/profile', b),
}

export const settingsApi = {
  get: () => get<Settings>('/settings'),
  update: (b: Partial<Settings>) => send<Settings>('put', '/settings', b),
  setClock: (date: string) => send<Settings>('put', '/dev/clock', { date }),
}

export const auditApi = {
  list: (p: PageParams & { table?: string; recordId?: number; userId?: number; from?: string; to?: string }) => list<AuditEntry>('/audit-log', p),
}
