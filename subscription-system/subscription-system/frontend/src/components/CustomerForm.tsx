import { Button, Form, Input } from 'antd'
import { useEffect, useState } from 'react'
import type { Customer, CustomerFields } from '../requests/api'
import { applyFieldErrors } from '../requests/client'

// The customer profile form, shared by the back office (create and edit) and the portal profile.
export function CustomerForm({ initial, submitLabel, onSubmit }: { initial?: Customer; submitLabel: string; onSubmit: (v: CustomerFields) => Promise<unknown> }) {
  const [form] = Form.useForm<CustomerFields>()
  const [saving, setSaving] = useState(false)
  useEffect(() => form.setFieldsValue(initial ?? {}), [form, initial])

  const finish = async (v: CustomerFields) => {
    setSaving(true)
    try {
      // "" clears an optional field on the server
      await onSubmit({ name: v.name?.trim(), companyName: v.companyName ?? '', phone: v.phone ?? '', billingAddress: v.billingAddress ?? '' })
    } catch (e) {
      applyFieldErrors(form, e)
    } finally {
      setSaving(false)
    }
  }

  return (
    <Form form={form} layout="vertical" onFinish={finish} style={{ maxWidth: 520 }}>
      <Form.Item name="name" label="Name" rules={[{ required: true, whitespace: true, message: 'Enter a name.' }, { max: 150 }]}>
        <Input autoComplete="name" />
      </Form.Item>
      <Form.Item name="companyName" label="Company" rules={[{ max: 150 }]}>
        <Input autoComplete="organization" />
      </Form.Item>
      <Form.Item name="phone" label="Phone" rules={[{ max: 30 }]}>
        <Input autoComplete="tel" />
      </Form.Item>
      <Form.Item name="billingAddress" label="Billing address" rules={[{ max: 1000 }]}>
        <Input.TextArea rows={3} autoComplete="street-address" />
      </Form.Item>
      <Button type="primary" htmlType="submit" loading={saving}>
        {submitLabel}
      </Button>
    </Form>
  )
}
