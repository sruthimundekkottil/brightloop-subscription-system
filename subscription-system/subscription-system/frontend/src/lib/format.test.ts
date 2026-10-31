import { describe, expect, it } from 'vitest'
import { formatDate, formatMoney } from './format'

describe('formatMoney', () => {
  it('shows two decimals whether or not the API trimmed zeros', () => {
    expect(formatMoney('1180')).toBe(formatMoney('1180.00'))
    expect(formatMoney('1180.5')).toContain('1,180.50')
  })
  it('handles the DECIMAL(12,2) maximum exactly', () => {
    expect(formatMoney('9999999999.99')).toContain('9,99,99,99,999.99') // en-IN grouping
  })
  it('shows a dash for missing values', () => {
    expect(formatMoney(null)).toBe('—')
  })
})

describe('formatDate', () => {
  it('does not shift business dates across time zones', () => {
    expect(formatDate('2026-02-28')).toBe('28 Feb 2026')
    expect(formatDate('2026-02-28 23:59:59')).toBe('28 Feb 2026')
  })
})
