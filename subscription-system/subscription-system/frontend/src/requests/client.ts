import axios, { AxiosError } from 'axios'
import { message } from 'antd'
import type { FormInstance } from 'antd'

export interface PageMeta {
  page: number
  pageSize: number
  total: number
}
export interface Page<T> {
  data: T[]
  meta: PageMeta
}

// ApiError is what every failed request rejects with: the server's { error: { code, message, fields } }.
export class ApiError extends Error {
  constructor(
    public status: number,
    public code: string,
    message: string,
    public fields: Record<string, string> = {},
  ) {
    super(message)
  }
}

export const http = axios.create({ baseURL: '/api', withCredentials: true })

let onUnauthenticated = () => {}
export function setOnUnauthenticated(f: () => void) {
  onUnauthenticated = f
}

// One interceptor for the whole app: field errors are left for the form, everything else is a toast.
http.interceptors.response.use(undefined, (err: AxiosError<{ error?: { code: string; message: string; fields?: Record<string, string> } }>) => {
  const body = err.response?.data?.error
  const e = body
    ? new ApiError(err.response!.status, body.code, body.message, body.fields ?? {})
    : new ApiError(0, 'NETWORK_ERROR', 'Cannot reach the server. Check your connection and try again.')
  const isSessionProbe = err.config?.url === '/auth/me'
  if (e.status === 401 && !isSessionProbe) onUnauthenticated()
  if (!isSessionProbe && Object.keys(e.fields).length === 0) message.error(e.message)
  return Promise.reject(e)
})

export async function get<T>(url: string, params?: object): Promise<T> {
  return (await http.get(url, { params })).data.data
}
export async function list<T>(url: string, params?: object): Promise<Page<T>> {
  return (await http.get(url, { params })).data
}
export async function send<T>(method: 'post' | 'put' | 'patch', url: string, body?: object, headers?: Record<string, string>): Promise<T> {
  return (await http.request({ method, url, data: body ?? {}, headers })).data.data
}

// applyFieldErrors shows server field errors under the matching form fields; returns false if there were none.
export function applyFieldErrors(form: FormInstance, err: unknown): boolean {
  if (!(err instanceof ApiError) || Object.keys(err.fields).length === 0) return false
  form.setFields(Object.entries(err.fields).map(([name, msg]) => ({ name, errors: [msg] })))
  return true
}
