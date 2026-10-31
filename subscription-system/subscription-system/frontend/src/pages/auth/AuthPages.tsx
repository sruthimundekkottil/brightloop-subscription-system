import { Alert, Button, Card, Form, Input, Result, Typography, message } from 'antd'
import type { ReactNode } from 'react'
import { useState } from 'react'
import { Link, useLocation, useNavigate, useSearchParams } from 'react-router-dom'
import { PasswordChecklist } from '../../components/PasswordChecklist'
import { passwordRule } from '../../lib/password'
import { authApi } from '../../requests/api'
import { applyFieldErrors } from '../../requests/client'
import { homeFor, useSession } from '../../stores/session'

function AuthCard({ title, children }: { title: string; children: ReactNode }) {
  return (
    <div style={{ minHeight: '100vh', display: 'grid', placeItems: 'center', padding: 16 }}>
      <Card style={{ width: '100%', maxWidth: 420 }}>
        <Typography.Title level={3}>{title}</Typography.Title>
        {children}
      </Card>
    </div>
  )
}

const emailRules = [{ required: true, message: 'Enter your email.' }, { type: 'email' as const, message: 'Enter a valid email address.' }]
const trimEmail = (e: string) => e.trim()

export function LoginPage() {
  const login = useSession((s) => s.login)
  const navigate = useNavigate()
  const from = (useLocation().state as { from?: string } | null)?.from
  const [busy, setBusy] = useState(false)

  const finish = async (v: { email: string; password: string }) => {
    setBusy(true)
    try {
      const me = await login(v.email, v.password)
      navigate(from ?? homeFor(me.role), { replace: true })
    } catch {
      // the interceptor already showed the server's message
    } finally {
      setBusy(false)
    }
  }

  return (
    <AuthCard title="Log in">
      <Form layout="vertical" onFinish={finish} requiredMark={false}>
        <Form.Item name="email" label="Email" rules={emailRules} normalize={trimEmail}>
          <Input autoComplete="email" autoFocus />
        </Form.Item>
        <Form.Item name="password" label="Password" rules={[{ required: true, message: 'Enter your password.' }]}>
          <Input.Password autoComplete="current-password" />
        </Form.Item>
        <Button type="primary" htmlType="submit" block loading={busy}>
          Log in
        </Button>
      </Form>
      <p style={{ marginTop: 16 }}>
        <Link to="/forgot-password">Forgot password?</Link> · <Link to="/register">Create an account</Link>
      </p>
    </AuthCard>
  )
}

export function RegisterPage() {
  const [form] = Form.useForm()
  const password = Form.useWatch('password', form)
  const setMe = useSession((s) => s.setMe)
  const navigate = useNavigate()
  const [busy, setBusy] = useState(false)

  const finish = async (v: { email: string; password: string; name: string; companyName?: string; phone?: string }) => {
    setBusy(true)
    try {
      const me = await authApi.register(v)
      setMe(me)
      message.success('Welcome to Brightloop!')
      navigate(homeFor(me.role), { replace: true })
    } catch (e) {
      applyFieldErrors(form, e)
    } finally {
      setBusy(false)
    }
  }

  return (
    <AuthCard title="Create an account">
      <Form form={form} layout="vertical" onFinish={finish} requiredMark="optional">
        <Form.Item name="name" label="Your name" rules={[{ required: true, whitespace: true, message: 'Enter your name.' }, { max: 150 }]}>
          <Input autoComplete="name" />
        </Form.Item>
        <Form.Item name="companyName" label="Company" rules={[{ max: 150 }]}>
          <Input autoComplete="organization" />
        </Form.Item>
        <Form.Item name="email" label="Email" rules={emailRules} normalize={trimEmail}>
          <Input autoComplete="email" />
        </Form.Item>
        <Form.Item name="password" label="Password" rules={[passwordRule]} style={{ marginBottom: 8 }}>
          <Input.Password autoComplete="new-password" />
        </Form.Item>
        <PasswordChecklist value={password} />
        <Button type="primary" htmlType="submit" block loading={busy}>
          Create account
        </Button>
      </Form>
      <p style={{ marginTop: 16 }}>
        Already registered? <Link to="/login">Log in</Link>
      </p>
    </AuthCard>
  )
}

export function ForgotPasswordPage() {
  const [sent, setSent] = useState<string>()
  const [busy, setBusy] = useState(false)
  const finish = async (v: { email: string }) => {
    setBusy(true)
    try {
      setSent((await authApi.forgotPassword(v.email)).message)
    } catch {
      // toast shown
    } finally {
      setBusy(false)
    }
  }
  if (sent) {
    return (
      <AuthCard title="Check your email">
        <Result status="success" subTitle={sent} extra={<Link to="/login">Back to log in</Link>} />
      </AuthCard>
    )
  }
  return (
    <AuthCard title="Forgot password">
      <Typography.Paragraph type="secondary">Enter your email and we'll send you a link to set a new password.</Typography.Paragraph>
      <Form layout="vertical" onFinish={finish} requiredMark={false}>
        <Form.Item name="email" label="Email" rules={emailRules} normalize={trimEmail}>
          <Input autoComplete="email" autoFocus />
        </Form.Item>
        <Button type="primary" htmlType="submit" block loading={busy}>
          Send reset link
        </Button>
      </Form>
      <p style={{ marginTop: 16 }}>
        <Link to="/login">Back to log in</Link>
      </p>
    </AuthCard>
  )
}

export function ResetPasswordPage() {
  const [params] = useSearchParams()
  const token = params.get('token') ?? ''
  const [form] = Form.useForm()
  const password = Form.useWatch('password', form)
  const navigate = useNavigate()
  const [tokenError, setTokenError] = useState<string>()
  const [busy, setBusy] = useState(false)

  const finish = async (v: { password: string }) => {
    setBusy(true)
    try {
      message.success((await authApi.resetPassword(token, v.password)).message)
      navigate('/login', { replace: true })
    } catch (e) {
      applyFieldErrors(form, e)
      const fields = (e as { fields?: Record<string, string> }).fields
      if (fields?.token) setTokenError(fields.token)
    } finally {
      setBusy(false)
    }
  }

  return (
    <AuthCard title="Set a new password">
      {!token || tokenError ? (
        <Alert
          type="error"
          showIcon
          message={tokenError ?? 'This reset link is incomplete.'}
          description={<Link to="/forgot-password">Request a new link</Link>}
          style={{ marginBottom: 16 }}
        />
      ) : null}
      <Form form={form} layout="vertical" onFinish={finish} requiredMark={false} disabled={!token}>
        <Form.Item name="password" label="New password" rules={[passwordRule]} style={{ marginBottom: 8 }}>
          <Input.Password autoComplete="new-password" autoFocus />
        </Form.Item>
        <PasswordChecklist value={password} />
        <Button type="primary" htmlType="submit" block loading={busy}>
          Change password
        </Button>
      </Form>
    </AuthCard>
  )
}
