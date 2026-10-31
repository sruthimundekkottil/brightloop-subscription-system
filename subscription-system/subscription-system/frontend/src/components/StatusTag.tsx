import { Tag } from 'antd'

// One colour per status across subscriptions, invoices, discounts and users.
const COLORS: Record<string, string> = {
  draft: 'default',
  quotation: 'blue',
  confirmed: 'cyan',
  active: 'green',
  paused: 'orange',
  closed: 'default',
  cancelled: 'red',
  paid: 'green',
  partial: 'gold',
  not_paid: 'orange',
  overdue: 'red',
  expired: 'volcano',
  inactive: 'default',
}

export function StatusTag({ status }: { status: string }) {
  const label = status.replace(/_/g, ' ')
  return <Tag color={COLORS[status] ?? 'default'}>{label.charAt(0).toUpperCase() + label.slice(1)}</Tag>
}
