// Mirror of backend/internal/logic/auth/password.go: same character classes, same limits, same messages.
// Characters are counted as code points ([...p]) like utf8.RuneCountInString; bytes as UTF-8 like len(p).

export type PasswordRule = 'length' | 'maxLength' | 'upper' | 'lower' | 'special'

export const PASSWORD_RULES: { key: PasswordRule; label: string }[] = [
  { key: 'length', label: 'Use more than 8 characters.' },
  { key: 'maxLength', label: 'Use at most 72 bytes.' },
  { key: 'upper', label: 'Add an uppercase letter.' },
  { key: 'lower', label: 'Add a lowercase letter.' },
  { key: 'special', label: 'Add a special character.' },
]

export function passwordProblems(p: string): Set<PasswordRule> {
  const probs = new Set<PasswordRule>()
  if ([...p].length < 9) probs.add('length')
  if (new TextEncoder().encode(p).length > 72) probs.add('maxLength')
  if (!/[A-Z]/.test(p)) probs.add('upper')
  if (!/[a-z]/.test(p)) probs.add('lower')
  if (!/[^A-Za-z0-9]/.test(p)) probs.add('special')
  return probs
}

// AntD Form rule: the same messages the server would send.
export const passwordRule = {
  validator: (_: unknown, value?: string) => {
    const probs = passwordProblems(value ?? '')
    return probs.size === 0
      ? Promise.resolve()
      : Promise.reject(new Error(PASSWORD_RULES.filter((r) => probs.has(r.key)).map((r) => r.label).join(' ')))
  },
}
