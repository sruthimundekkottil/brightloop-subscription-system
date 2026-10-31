import { describe, expect, it } from 'vitest'
import { passwordProblems } from './password'

// Same table as backend/internal/logic/auth/password_test.go; keep the two in step.
const cases: [string, string[]][] = [
  ['Abcdefg!9', []],
  ['Abcdef!9', ['length']],
  ['abcdefgh!', ['upper']],
  ['ABCDEFGH!', ['lower']],
  ['Abcdefghi', ['special']],
  ['Abcdefgh ', []],
  ['Aa!' + 'x'.repeat(69), []],
  ['Aa!' + 'x'.repeat(70), ['maxLength']],
  ['Ab!' + 'é'.repeat(35), ['maxLength']],
  ['Äbcdéfgh', ['length', 'upper']],
  ['', ['length', 'lower', 'special', 'upper']],
]

describe('passwordProblems matches the server rule', () => {
  it.each(cases)('%j', (pw, want) => {
    expect([...passwordProblems(pw)].sort()).toEqual([...want].sort())
  })
})
