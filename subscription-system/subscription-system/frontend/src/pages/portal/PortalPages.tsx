import { Card, Empty, Result, Spin, Typography, message } from 'antd'
import { useEffect, useState } from 'react'
import { CustomerForm } from '../../components/CustomerForm'
import { customerApi, type Customer } from '../../requests/api'

// Portal shell. Member C fills "My Subscriptions", member D "My Invoices"; both read only the logged-in customer's data.
export function MySubscriptionsPage() {
  return (
    <>
      <Typography.Title level={3}>My Subscriptions</Typography.Title>
      <Card>
        <Empty description="Your subscriptions will appear here." />
      </Card>
    </>
  )
}

export function MyInvoicesPage() {
  return (
    <>
      <Typography.Title level={3}>My Invoices</Typography.Title>
      <Card>
        <Empty description="Your invoices will appear here." />
      </Card>
    </>
  )
}

export function MyProfilePage() {
  const [profile, setProfile] = useState<Customer>()
  const [failed, setFailed] = useState(false)
  useEffect(() => {
    customerApi.myProfile().then(setProfile, () => setFailed(true))
  }, [])
  if (failed) return <Result status="error" title="Your profile could not be loaded." />
  if (!profile) return <Spin />
  return (
    <>
      <Typography.Title level={3}>My Profile</Typography.Title>
      <Card>
        <CustomerForm
          initial={profile}
          submitLabel="Save changes"
          onSubmit={async (v) => {
            setProfile(await customerApi.updateMyProfile(v))
            message.success('Profile saved.')
          }}
        />
      </Card>
    </>
  )
}
