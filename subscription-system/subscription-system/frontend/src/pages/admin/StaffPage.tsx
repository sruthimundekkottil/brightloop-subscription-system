import { PlusOutlined } from '@ant-design/icons'
import { Button, Form, Input, Modal, Popconfirm, Space, Switch, Table, Typography, message } from 'antd'
import { useState } from 'react'
import { PasswordChecklist } from '../../components/PasswordChecklist'
import { StatusTag } from '../../components/StatusTag'
import { usePaged } from '../../hooks/usePaged'
import { formatDateTime } from '../../lib/format'
import { passwordRule } from '../../lib/password'
import { userApi, type User } from '../../requests/api'
import { applyFieldErrors } from '../../requests/client'
import { useSession } from '../../stores/session'

export function StaffPage() {
  const me = useSession((s) => s.me)
  const [q, setQ] = useState('')
  const { rows, loading, reload, pagination } = usePaged(userApi.list, { role: 'staff' as const, q })
  const [adding, setAdding] = useState(false)

  const toggle = async (u: User) => {
    await userApi.setStatus(u.id, !u.isActive)
    message.success(u.isActive ? `${u.email} can no longer log in.` : `${u.email} can log in again.`)
    reload()
  }

  return (
    <>
      <Space style={{ width: '100%', justifyContent: 'space-between', marginBottom: 16 }} wrap>
        <Typography.Title level={3} style={{ margin: 0 }}>
          Staff
        </Typography.Title>
        <Space wrap>
          <Input.Search placeholder="Search name or email" allowClear onSearch={setQ} style={{ width: 260 }} />
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setAdding(true)}>
            Add staff
          </Button>
        </Space>
      </Space>
      <Table<User>
        rowKey="id"
        loading={loading}
        dataSource={rows}
        pagination={pagination}
        scroll={{ x: 'max-content' }}
        columns={[
          { title: 'Name', dataIndex: 'fullName' },
          { title: 'Email', dataIndex: 'email' },
          { title: 'Job title', dataIndex: 'jobTitle', render: (v) => v ?? '—' },
          { title: 'Status', dataIndex: 'isActive', render: (v) => <StatusTag status={v ? 'active' : 'inactive'} /> },
          { title: 'Added', dataIndex: 'createdAt', render: formatDateTime },
          {
            title: 'Can log in',
            render: (_, u) => (
              <Popconfirm
                title={u.isActive ? 'Deactivate this staff member?' : 'Activate this staff member?'}
                description={u.isActive ? 'They are logged out at once and cannot log in.' : undefined}
                onConfirm={() => toggle(u)}
                disabled={u.id === me?.id}
              >
                <Switch checked={u.isActive} disabled={u.id === me?.id} aria-label={`Can ${u.email} log in`} />
              </Popconfirm>
            ),
          },
        ]}
      />
      <AddStaffModal
        open={adding}
        onClose={() => setAdding(false)}
        onAdded={() => {
          setAdding(false)
          reload()
        }}
      />
    </>
  )
}

function AddStaffModal({ open, onClose, onAdded }: { open: boolean; onClose: () => void; onAdded: () => void }) {
  const [form] = Form.useForm()
  const password = Form.useWatch('password', form)
  const [busy, setBusy] = useState(false)
  const submit = async () => {
    const v = await form.validateFields()
    setBusy(true)
    try {
      const u = await userApi.createStaff(v)
      message.success(`${u.email} can now log in as staff.`)
      form.resetFields()
      onAdded()
    } catch (e) {
      applyFieldErrors(form, e)
    } finally {
      setBusy(false)
    }
  }
  return (
    <Modal title="Add staff" open={open} onCancel={onClose} onOk={submit} okText="Add staff" confirmLoading={busy} destroyOnClose>
      <Form form={form} layout="vertical" preserve={false}>
        <Form.Item name="fullName" label="Full name" rules={[{ required: true, whitespace: true, message: 'Enter a name.' }, { max: 150 }]}>
          <Input />
        </Form.Item>
        <Form.Item name="jobTitle" label="Job title" rules={[{ max: 100 }]}>
          <Input />
        </Form.Item>
        <Form.Item name="email" label="Email" normalize={(v: string) => v.trim()} rules={[{ required: true, type: 'email', message: 'Enter a valid email address.' }]}>
          <Input autoComplete="off" />
        </Form.Item>
        <Form.Item name="password" label="First password" rules={[passwordRule]} style={{ marginBottom: 8 }} extra="Share it with them privately; they can change it with Forgot password.">
          <Input.Password autoComplete="new-password" />
        </Form.Item>
        <PasswordChecklist value={password} />
      </Form>
    </Modal>
  )
}
