import { DatePicker, Input, InputNumber, Result, Space, Table, Typography } from 'antd'
import type { Dayjs } from 'dayjs'
import { useState } from 'react'
import { usePaged } from '../../hooks/usePaged'
import { formatDateTime } from '../../lib/format'
import { auditApi, type AuditEntry } from '../../requests/api'

// Renders {"field": [old, new]} as "field: old → new"; anything else as JSON.
function Changes({ value }: { value: unknown }) {
  if (!value || typeof value !== 'object') return <>—</>
  const entries = Object.entries(value as Record<string, unknown>)
  return (
    <>
      {entries.map(([k, v]) => (
        <div key={k}>
          <Typography.Text type="secondary">{k}:</Typography.Text>{' '}
          {Array.isArray(v) && v.length === 2 ? `${JSON.stringify(v[0])} → ${JSON.stringify(v[1])}` : JSON.stringify(v)}
        </div>
      ))}
    </>
  )
}

export function AuditLogPage() {
  const [table, setTable] = useState('')
  const [recordId, setRecordId] = useState<number>()
  const [range, setRange] = useState<[Dayjs | null, Dayjs | null] | null>(null)
  const filters = {
    table: table || undefined,
    recordId,
    from: range?.[0]?.format('YYYY-MM-DD'),
    to: range?.[1]?.format('YYYY-MM-DD'),
  }
  const { rows, loading, error, reload, pagination } = usePaged(auditApi.list, filters)

  return (
    <>
      <Typography.Title level={3}>Audit log</Typography.Title>
      <Space wrap style={{ marginBottom: 16 }}>
        <Input.Search placeholder="Table, e.g. customers" allowClear onSearch={(v) => setTable(v.trim())} style={{ width: 220 }} />
        <InputNumber placeholder="Record id" min={1} value={recordId} onChange={(v) => setRecordId(v ?? undefined)} />
        <DatePicker.RangePicker value={range} onChange={setRange} allowEmpty={[true, true]} />
      </Space>
      {error ? (
        <Result status="error" title="The audit log could not be loaded." extra={<a onClick={reload}>Try again</a>} />
      ) : (
        <Table<AuditEntry>
          rowKey="id"
          loading={loading}
          dataSource={rows}
          pagination={pagination}
          scroll={{ x: 'max-content' }}
          columns={[
            { title: 'When', dataIndex: 'createdAt', render: formatDateTime },
            { title: 'Who', dataIndex: 'changedByEmail', render: (v) => v ?? <Typography.Text type="secondary">System</Typography.Text> },
            { title: 'Record', render: (_, r) => `${r.tableName} #${r.recordId}` },
            { title: 'Action', dataIndex: 'action' },
            { title: 'Reason', dataIndex: 'reason', render: (v) => v ?? '—' },
            { title: 'Changes', dataIndex: 'changes', render: (v) => <Changes value={v} /> },
          ]}
        />
      )}
    </>
  )
}
