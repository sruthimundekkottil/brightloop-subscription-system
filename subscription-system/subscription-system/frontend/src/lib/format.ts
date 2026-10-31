import dayjs from 'dayjs'
import utc from 'dayjs/plugin/utc'

dayjs.extend(utc)

// Money arrives as a decimal string ("1180.00", or "1180" from shopspring/decimal). Display only: never do
// arithmetic on the result. shortcut: Number() is exact for display up to DECIMAL(12,2)'s 9,999,999,999.99.
export function formatMoney(amount: string | null | undefined, currency = 'INR'): string {
  if (amount == null || amount === '') return '—'
  const n = Number(amount)
  if (Number.isNaN(n)) return amount
  return new Intl.NumberFormat('en-IN', { style: 'currency', currency, minimumFractionDigits: 2 }).format(n)
}

// Business dates are YYYY-MM-DD with no time zone: format them as they are.
export function formatDate(d: string | null | undefined): string {
  return d ? dayjs(d.slice(0, 10)).format('DD MMM YYYY') : '—'
}

// Timestamps from the API are UTC ("2026-10-31 08:00:00"); show them in the viewer's time zone.
export function formatDateTime(ts: string | null | undefined): string {
  return ts ? dayjs.utc(ts).local().format('DD MMM YYYY, HH:mm') : '—'
}
