// Package seed loads the reference data and the three demo logins. Safe to re-run: existing rows are kept.
package seed

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/crypto/bcrypt"
)

// DevPassword is the password of every seeded login. Development data only.
const DevPassword = "Brightloop#2026"

const (
	AdminEmail    = "admin@brightloop.test"
	StaffEmail    = "staff@brightloop.test"
	CustomerEmail = "customer@brightloop.test"
)

// Run inserts what is missing. It writes no audit rows: this is bootstrap data, not a user action.
func Run(ctx context.Context) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(DevPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	stmts := []struct {
		sql  string
		args []any
	}{
		{"INSERT INTO company_settings (id) VALUES (1) ON DUPLICATE KEY UPDATE id = id", nil},
		{`INSERT INTO payment_terms (name, due_days) VALUES ('Immediate', 0), ('Net 15', 15), ('Net 30', 30)
			ON DUPLICATE KEY UPDATE id = id`, nil},
		// is_churn = FALSE for closes that are not a lost customer (dashboard churn excludes them)
		{`INSERT INTO close_reasons (name, is_churn) VALUES ('Renewed', FALSE), ('Quotation expired', FALSE),
			('Ended', TRUE), ('Cancelled at end of term', TRUE), ('Too expensive', TRUE), ('Switched to a competitor', TRUE),
			('No longer needed', TRUE), ('Other', TRUE) ON DUPLICATE KEY UPDATE id = id`, nil},
		{`INSERT INTO users (email, password_hash, role) VALUES (?, ?, 'admin'), (?, ?, 'staff'), (?, ?, 'customer')
			ON DUPLICATE KEY UPDATE id = id`, []any{AdminEmail, hash, StaffEmail, hash, CustomerEmail, hash}},
		{`INSERT INTO staff_profiles (user_id, full_name, job_title, created_by)
			SELECT s.id, 'Sam Staff', 'Account manager', a.id FROM users s JOIN users a ON a.email = ? WHERE s.email = ?
			ON DUPLICATE KEY UPDATE staff_profiles.id = staff_profiles.id`, []any{AdminEmail, StaffEmail}},
		{`INSERT INTO customers (user_id, name, company_name, billing_address)
			SELECT id, 'Casey Customer', 'Small Shop', '12 Market Road, Kochi' FROM users WHERE email = ?
			ON DUPLICATE KEY UPDATE customers.id = customers.id`, []any{CustomerEmail}},
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, s := range stmts {
			if _, err := tx.Exec(s.sql, s.args...); err != nil {
				return err
			}
		}
		return nil
	})
}
