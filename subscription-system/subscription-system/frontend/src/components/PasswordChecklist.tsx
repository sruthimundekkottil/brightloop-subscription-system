import { CheckCircleFilled, CloseCircleOutlined } from '@ant-design/icons'
import { Space, Typography } from 'antd'
import { PASSWORD_RULES, passwordProblems } from '../lib/password'

// Live checklist under a password field; each rule turns green as it is met.
export function PasswordChecklist({ value }: { value?: string }) {
  const probs = passwordProblems(value ?? '')
  return (
    <Space direction="vertical" size={0} style={{ marginBottom: 16 }} aria-live="polite">
      {PASSWORD_RULES.map((r) => {
        const ok = !probs.has(r.key)
        return (
          <Typography.Text key={r.key} type={ok ? 'success' : 'secondary'}>
            {ok ? <CheckCircleFilled /> : <CloseCircleOutlined />} {r.label}
            <span style={{ position: 'absolute', left: -9999 }}>{ok ? ' (met)' : ' (not met)'}</span>
          </Typography.Text>
        )
      })}
    </Space>
  )
}
