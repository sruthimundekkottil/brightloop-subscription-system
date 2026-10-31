package billingdate

import (
	"testing"
	"time"
)

func d(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestAddDateNormalisesOverflow(t *testing.T) {
	// Documents the trap: the standard library does not clamp.
	if got := d("2026-01-31").AddDate(0, 1, 0); !got.Equal(d("2026-03-03")) {
		t.Fatalf("expected 2026-03-03, got %s", got.Format("2006-01-02"))
	}
}

func TestPeriodStart(t *testing.T) {
	cases := []struct {
		anchor, period string
		every, k       int
		want           string
	}{
		{"2026-01-31", "monthly", 1, 0, "2026-01-31"},
		{"2026-01-31", "monthly", 1, 1, "2026-02-28"},
		{"2026-01-31", "monthly", 1, 2, "2026-03-31"},
		{"2026-01-31", "monthly", 1, 3, "2026-04-30"},
		{"2026-01-31", "monthly", 1, 4, "2026-05-31"},
		{"2028-01-31", "monthly", 1, 1, "2028-02-29"}, // leap year
		{"2026-01-31", "monthly", 3, 1, "2026-04-30"}, // quarterly
		{"2024-02-29", "yearly", 1, 1, "2025-02-28"},
		{"2024-02-29", "yearly", 1, 4, "2028-02-29"},
		{"2026-03-01", "weekly", 1, 2, "2026-03-15"},
		{"2026-03-01", "daily", 1, 31, "2026-04-01"},
	}
	for _, c := range cases {
		got := PeriodStart(d(c.anchor), c.period, c.every, c.k).Format("2006-01-02")
		if got != c.want {
			t.Errorf("PeriodStart(%s,%s,%d,%d) = %s, want %s", c.anchor, c.period, c.every, c.k, got, c.want)
		}
	}
}

func TestPeriodEnd(t *testing.T) {
	if got := PeriodEnd(d("2026-01-31"), "monthly", 1, 0).Format("2006-01-02"); got != "2026-02-27" {
		t.Errorf("got %s, want 2026-02-27", got)
	}
	if got := PeriodEnd(d("2026-03-01"), "monthly", 1, 0).Format("2006-01-02"); got != "2026-03-31" {
		t.Errorf("got %s, want 2026-03-31", got)
	}
}

func TestNextPeriodOnOrAfter(t *testing.T) {
	k, s := NextPeriodOnOrAfter(d("2026-01-31"), "monthly", 1, d("2026-04-15"))
	if k != 3 || s.Format("2006-01-02") != "2026-04-30" {
		t.Errorf("got k=%d start=%s, want k=3 start=2026-04-30", k, s.Format("2006-01-02"))
	}
	k, s = NextPeriodOnOrAfter(d("2026-01-31"), "monthly", 1, d("2026-03-31"))
	if k != 2 || s.Format("2006-01-02") != "2026-03-31" {
		t.Errorf("got k=%d start=%s, want k=2 start=2026-03-31", k, s.Format("2006-01-02"))
	}
}
