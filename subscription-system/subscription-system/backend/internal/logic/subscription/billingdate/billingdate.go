// Package billingdate computes subscription billing periods from an anchor date.
// Periods are always derived from the anchor (never from the previous period),
// so month-end dates don't drift. Dates are UTC midnights.
package billingdate

import "time"

// addMonthsClamped adds months to t, clamping the day to the target month's last day.
// time.AddDate is NOT used for months because it normalises overflow
// (31 Jan + 1 month = 3 Mar in a non-leap year).
func addMonthsClamped(t time.Time, months int) time.Time {
	y, m, d := t.Date()
	first := time.Date(y, m+time.Month(months), 1, 0, 0, 0, 0, time.UTC)
	last := first.AddDate(0, 1, -1).Day()
	if d > last {
		d = last
	}
	return time.Date(first.Year(), first.Month(), d, 0, 0, 0, 0, time.UTC)
}

// PeriodStart returns the start of billing period k (k = 0 is the anchor).
// period is "daily", "weekly", "monthly" or "yearly"; every is the interval count (>= 1).
func PeriodStart(anchor time.Time, period string, every, k int) time.Time {
	switch period {
	case "daily":
		return anchor.AddDate(0, 0, every*k)
	case "weekly":
		return anchor.AddDate(0, 0, 7*every*k)
	case "monthly":
		return addMonthsClamped(anchor, every*k)
	default: // "yearly"
		return addMonthsClamped(anchor, 12*every*k)
	}
}

// PeriodEnd returns the last day of billing period k.
func PeriodEnd(anchor time.Time, period string, every, k int) time.Time {
	return PeriodStart(anchor, period, every, k+1).AddDate(0, 0, -1)
}

// NextPeriodOnOrAfter returns the index and start of the first period starting on or after date.
func NextPeriodOnOrAfter(anchor time.Time, period string, every int, date time.Time) (int, time.Time) {
	k := 0
	for {
		s := PeriodStart(anchor, period, every, k)
		if !s.Before(date) {
			return k, s
		}
		k++
	}
}
