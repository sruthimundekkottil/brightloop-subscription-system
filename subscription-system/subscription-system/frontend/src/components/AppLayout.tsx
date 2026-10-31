import {
  AuditOutlined,
  FileTextOutlined,
  LogoutOutlined,
  ProfileOutlined,
  SettingOutlined,
  TeamOutlined,
  UserOutlined,
} from '@ant-design/icons'
import { Button, Layout, Menu, Space, Tag, Typography } from 'antd'
import type { MenuProps } from 'antd'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import type { Role } from '../requests/api'
import { useSession } from '../stores/session'

type Item = Required<MenuProps>['items'][number] & { key: string }

// Admin sees configuration, staff sees operations, customers see the portal only. Members B-D add their entries here.
const MENU: Record<Role, Item[]> = {
  admin: [
    { key: '/customers', icon: <UserOutlined />, label: 'Customers' },
    { key: '/staff', icon: <TeamOutlined />, label: 'Staff' },
    { key: '/settings', icon: <SettingOutlined />, label: 'Settings' },
    { key: '/audit-log', icon: <AuditOutlined />, label: 'Audit log' },
  ],
  staff: [
    { key: '/customers', icon: <UserOutlined />, label: 'Customers' },
    { key: '/settings', icon: <SettingOutlined />, label: 'Settings' },
  ],
  customer: [
    { key: '/portal', icon: <FileTextOutlined />, label: 'My Subscriptions' },
    { key: '/portal/invoices', icon: <FileTextOutlined />, label: 'My Invoices' },
    { key: '/portal/profile', icon: <ProfileOutlined />, label: 'My Profile' },
  ],
}

export function AppLayout() {
  const { me, logout } = useSession()
  const navigate = useNavigate()
  const { pathname } = useLocation()
  if (!me) return null
  const items = MENU[me.role]
  // Highlight the longest menu key that prefixes the path (/customers/7 → Customers).
  const selected = items.map((i) => i.key).filter((k) => pathname === k || pathname.startsWith(k + '/')).sort((a, b) => b.length - a.length)[0]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Layout.Sider breakpoint="lg" collapsedWidth={0} theme="light">
        <Typography.Title level={4} style={{ padding: '16px 24px', margin: 0 }}>
          Brightloop
        </Typography.Title>
        <Menu mode="inline" selectedKeys={selected ? [selected] : []} items={items} onClick={(e) => navigate(e.key)} />
      </Layout.Sider>
      <Layout>
        <Layout.Header style={{ background: 'transparent', display: 'flex', justifyContent: 'flex-end', paddingInline: 24 }}>
          <Space>
            <Typography.Text>{me.email}</Typography.Text>
            <Tag>{me.role}</Tag>
            <Button icon={<LogoutOutlined />} onClick={() => logout().then(() => navigate('/login'))}>
              Log out
            </Button>
          </Space>
        </Layout.Header>
        <Layout.Content style={{ padding: 24 }}>
          <Outlet />
        </Layout.Content>
      </Layout>
    </Layout>
  )
}
