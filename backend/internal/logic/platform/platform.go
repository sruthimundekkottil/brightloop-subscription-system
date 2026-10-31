// Package platform holds the shared infrastructure billing relies on:
// the test clock, the audit log, gapless document numbers and the job queue.
package platform

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

const DateLayout = "2006-01-02"

var (
	ErrClockDisabled  = errors.New("test clock is only available when APP_ENV=development")
	ErrClockBackwards = errors.New("test clock can only move forward")
)

func devMode() bool { return os.Getenv("APP_ENV") == "development" }

// Today is the UTC business date. In development it is the test clock when that is ahead of real time.
func Today(ctx context.Context) (time.Time, error) {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if !devMode() {
		return today, nil
	}
	v, err := g.DB().GetValue(ctx, "SELECT DATE_FORMAT(test_clock_date, '%Y-%m-%d') FROM company_settings WHERE id = 1")
	if err != nil || v.IsEmpty() {
		return today, err
	}
	clock, err := time.Parse(DateLayout, v.String())
	if err != nil {
		return today, err
	}
	if clock.After(today) {
		return clock, nil
	}
	return today, nil
}

// SetTestClock moves the development clock to date; it never moves back.
func SetTestClock(ctx context.Context, date time.Time) error {
	if !devMode() {
		return ErrClockDisabled
	}
	d := date.Format(DateLayout)
	res, err := g.DB().Exec(ctx,
		"UPDATE company_settings SET test_clock_date = ? WHERE id = 1 AND (test_clock_date IS NULL OR test_clock_date < ?)", d, d)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrClockBackwards
	}
	return nil
}

// Audit writes one audit_log row inside the caller's transaction. changed_by is the Actor in the transaction's
// context, or NULL for system jobs.
func Audit(tx gdb.TX, table string, recordID uint64, action, reason string, changes any) error {
	js, err := json.Marshal(changes)
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO audit_log (table_name, record_id, action, changed_by, reason, changes) VALUES (?, ?, ?, NULLIF(?, 0), NULLIF(?, ''), ?)",
		table, recordID, action, ActorFrom(tx.GetCtx()).UserID, reason, string(js))
	return err
}

// NextNumber takes the next value of a gapless counter. The row lock is held until tx ends,
// and a rollback returns the number, so committed numbers have no gaps.
func NextNumber(tx gdb.TX, name string) (uint64, error) {
	if _, err := tx.Exec("INSERT INTO number_sequences (name) VALUES (?) ON DUPLICATE KEY UPDATE name = name", name); err != nil {
		return 0, err
	}
	v, err := tx.GetValue("SELECT next_value FROM number_sequences WHERE name = ? FOR UPDATE", name)
	if err != nil {
		return 0, err
	}
	if _, err = tx.Exec("UPDATE number_sequences SET next_value = next_value + 1 WHERE name = ?", name); err != nil {
		return 0, err
	}
	return v.Uint64(), nil
}
