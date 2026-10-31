import { PlusOutlined } from '@ant-design/icons'
import { Breadcrumb, Button, Card, Descriptions, Empty, Input, Modal, Result, Space, Spin, Table, Tabs, Tag, Typography, message } from 'antd'
import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { CustomerForm } from '../../components/CustomerForm'
import { usePaged } from '../../hooks/usePaged'
import { formatDateTime } from '../../lib/format'
import { customerApi, type Customer } from '../../requests/api'
import { ApiError } from '../../requests/client'

export function CustomersPage() {
  const [q, setQ] = useState('')
  const { rows, loading, error, reload, pagination } = usePaged(customerApi.list, { q })
  const [adding, setAdding] = useState(false)
  const navigate = useNavigate()

  return (
    <>
      <Space style={{ width: '100%', justifyContent: 'space-between', marginBottom: 16 }} wrap>
        <Typography.Title level={3} style={{ margin: 0 }}>
          Customers
        </Typography.Title>
        <Space wrap>
          <Input.Search placeholder="Search name, company or email" allowClear onSearch={setQ} style={{ width: 280 }} />
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setAdding(true)}>
            New customer
          </Button>
        </Space>
      </Space>
      {error ? (
        <Result status="error" title="Customers could not be loaded." extra={<Button onClick={reload}>Try again</Button>} />
      ) : (
        <Table<Customer>
          rowKey="id"
          loading={loading}
          dataSource={rows}
          pagination={pagination}
          scroll={{ x: 'max-content' }}
          locale={{ emptyText: <Empty description={q ? 'No customers match your search.' : 'No customers yet.'} /> }}
          onRow={(c) => ({ onClick: () => navigate(`/customers/${c.id}`), style: { cursor: 'pointer' } })}
          columns={[
            { title: 'Name', dataIndex: 'name', render: (v, c) => <Link to={`/customers/${c.id}`}>{v}</Link> },
            { title: 'Company', dataIndex: 'companyName', render: (v) => v ?? '—' },
            { title: 'Login email', dataIndex: 'email', render: (v) => v ?? <Tag>No login</Tag> },
            { title: 'Phone', dataIndex: 'phone', render: (v) => v ?? '—' },
            { title: 'Added', dataIndex: 'createdAt', render: formatDateTime },
          ]}
        />
      )}
      <Modal title="New customer" open={adding} onCancel={() => setAdding(false)} footer={null} destroyOnClose>
        <CustomerForm
          submitLabel="Create customer"
          onSubmit={async (v) => {
            const c = await customerApi.create(v)
            message.success(`${c.name} added.`)
            navigate(`/customers/${c.id}`)
          }}
        />
      </Modal>
    </>
  )
}

export function CustomerDetailPage() {
  const id = Number(useParams().id)
  const [customer, setCustomer] = useState<Customer>()
  const [missing, setMissing] = useState(false)

  useEffect(() => {
    customerApi.get(id).then(setCustomer, (e) => setMissing(e instanceof ApiError && e.status === 404))
  }, [id])

  if (missing) return <Result status="404" title="Customer not found" extra={<Link to="/customers">Back to customers</Link>} />
  if (!customer) return <Spin />

  return (
    <>
      <Breadcrumb items={[{ title: <Link to="/customers">Customers</Link> }, { title: customer.name }]} style={{ marginBottom: 16 }} />
      <Card>
        <Descriptions title={customer.name} column={{ xs: 1, md: 2 }}>
          <Descriptions.Item label="Company">{customer.companyName ?? '—'}</Descriptions.Item>
          <Descriptions.Item label="Login email">{customer.email ?? <Tag>No login</Tag>}</Descriptions.Item>
          <Descriptions.Item label="Phone">{customer.phone ?? '—'}</Descriptions.Item>
          <Descriptions.Item label="Added">{formatDateTime(customer.createdAt)}</Descriptions.Item>
          <Descriptions.Item label="Billing address" span={2}>
            <span style={{ whiteSpace: 'pre-line' }}>{customer.billingAddress ?? '—'}</span>
          </Descriptions.Item>
        </Descriptions>
      </Card>
      {/* Members C and D replace the two placeholders with their subscription and invoice tables. */}
      <Tabs
        style={{ marginTop: 16 }}
        items={[
          { key: 'subscriptions', label: 'Subscriptions', children: <Empty description="Subscriptions appear here." /> },
          { key: 'invoices', label: 'Invoices', children: <Empty description="Invoices appear here." /> },
          {
            key: 'profile',
            label: 'Edit profile',
            children: (
              <CustomerForm
                initial={customer}
                submitLabel="Save changes"
                onSubmit={async (v) => {
                  setCustomer(await customerApi.update(id, v))
                  message.success('Customer saved.')
                }}
              />
            ),
          },
        ]}
      />
    </>
  )
}
