import { Alert, Button, Card, DatePicker, Form, Input, InputNumber, Popconfirm, Result, Space, Spin, Switch, Typography, message } from 'antd'
import dayjs, { type Dayjs } from 'dayjs'
import { useEffect, useState } from 'react'
import { formatDate } from '../../lib/format'
import { settingsApi, type Settings } from '../../requests/api'
import { applyFieldErrors } from '../../requests/client'
import { useSession } from '../../stores/session'

export function SettingsPage() {
  const isAdmin = useSession((s) => s.me?.role === 'admin')
  const [settings, setSettings] = useState<Settings>()
  const [failed, setFailed] = useState(false)
  const [form] = Form.useForm<Settings>()
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    settingsApi.get().then(setSettings, () => setFailed(true))
  }, [])
  useEffect(() => settings && form.setFieldsValue(settings), [form, settings])

  if (failed) return <Result status="error" title="Settings could not be loaded." extra={<Button onClick={() => location.reload()}>Try again</Button>} />
  if (!settings) return <Spin />

  const save = async (v: Settings) => {
    setSaving(true)
    try {
      setSettings(await settingsApi.update({ ...v, address: v.address ?? '' }))
      message.success('Settings saved.')
    } catch (e) {
      applyFieldErrors(form, e)
    } finally {
      setSaving(false)
    }
  }

  return (
    <Space direction="vertical" size="large" style={{ width: '100%', maxWidth: 640 }}>
      <Typography.Title level={3} style={{ margin: 0 }}>
        Settings
      </Typography.Title>
      {!isAdmin && <Alert type="info" showIcon message="Only an admin can change these settings." />}
      <Card title="Company">
        <Form form={form} layout="vertical" onFinish={save} disabled={!isAdmin}>
          <Form.Item name="companyName" label="Company name" rules={[{ required: true, whitespace: true }, { max: 150 }]}>
            <Input />
          </Form.Item>
          <Form.Item name="address" label="Address" rules={[{ max: 1000 }]}>
            <Input.TextArea rows={3} />
          </Form.Item>
          <Space wrap size="large">
            <Form.Item name="currencyCode" label="Currency" rules={[{ required: true, pattern: /^[A-Z]{3}$/, message: 'Three capital letters, e.g. INR.' }]}>
              <Input style={{ width: 100 }} maxLength={3} />
            </Form.Item>
            <Form.Item name="invoicePrefix" label="Invoice prefix" rules={[{ required: true, pattern: /^[A-Z0-9]{1,10}$/, message: 'Up to 10 capital letters or digits.' }]}>
              <Input style={{ width: 140 }} maxLength={10} />
            </Form.Item>
            <Form.Item name="defaultQuotationValidityDays" label="Quotation valid for (days)" rules={[{ required: true, type: 'integer', min: 1, max: 365 }]}>
              <InputNumber min={1} max={365} />
            </Form.Item>
          </Space>
          <Form.Item name="autoConfirmRecurringInvoices" label="Confirm recurring invoices automatically" valuePropName="checked" extra="When off, the nightly billing run creates drafts for staff to review.">
            <Switch />
          </Form.Item>
          {isAdmin && (
            <Button type="primary" htmlType="submit" loading={saving}>
              Save settings
            </Button>
          )}
        </Form>
      </Card>
      {settings.devMode && isAdmin && <TestClockCard settings={settings} onChange={setSettings} />}
    </Space>
  )
}

// Development only: move the business date forward to demo billing, renewals and overdue reminders.
function TestClockCard({ settings, onChange }: { settings: Settings; onChange: (s: Settings) => void }) {
  const today = dayjs(settings.today)
  const [date, setDate] = useState<Dayjs | null>(null)
  const [busy, setBusy] = useState(false)
  const move = async () => {
    if (!date) return
    setBusy(true)
    try {
      onChange(await settingsApi.setClock(date.format('YYYY-MM-DD')))
      message.success(`Today is now ${date.format('DD MMM YYYY')}.`)
      setDate(null)
    } catch {
      // toast shown
    } finally {
      setBusy(false)
    }
  }
  return (
    <Card title="Test clock (development only)">
      <Typography.Paragraph>
        Business date: <strong>{formatDate(settings.today)}</strong>
        {settings.testClockDate ? ' (test clock)' : ' (real date)'}
      </Typography.Paragraph>
      <Typography.Paragraph type="secondary">The clock only moves forward. The nightly jobs use this date the next time they run.</Typography.Paragraph>
      <Space wrap>
        <DatePicker value={date} onChange={setDate} disabledDate={(d) => !d.isAfter(today, 'day')} aria-label="New business date" />
        <Popconfirm title={`Move the clock to ${date?.format('DD MMM YYYY')}?`} description="This cannot be undone." onConfirm={move} disabled={!date}>
          <Button type="primary" disabled={!date} loading={busy}>
            Move clock forward
          </Button>
        </Popconfirm>
      </Space>
    </Card>
  )
}
