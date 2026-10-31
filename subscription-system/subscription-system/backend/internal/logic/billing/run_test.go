package billing

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"brightloop/internal/testdb"
)

const dbName = "bl_test_billing"

var asOf = time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC) // anchor 31 Jan: periods 31 Jan, 28 Feb, 31 Mar are due

func TestMain(m *testing.M) {
	if os.Getenv("BILLING_CRASH_AFTER") != "" {
		testdb.Connect(dbName) // crash child: use the parent's seeded database as is
	} else if err := testdb.Setup(dbName); err != nil {
		fmt.Println("testdb:", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func exec1(t *testing.T, sql string, args ...any) {
	t.Helper()
	if _, err := g.DB().Exec(context.Background(), sql, args...); err != nil {
		t.Fatalf("%v\n%s", err, sql)
	}
}

func count(t *testing.T, sql string, args ...any) int {
	t.Helper()
	v, err := g.DB().GetValue(context.Background(), sql, args...)
	if err != nil {
		t.Fatal(err)
	}
	return v.Int()
}

// reset gives the test a fresh database with n active monthly subscriptions anchored on 31 Jan 2026.
// Each bills: plan 500.00 + GST 18%; 3 CRM seats at 33.33 with GST 18% + Cess 1%, 10% off for the first 2 cycles;
// and a one-time setup product, which is out of scope and must not be billed.
func reset(t *testing.T, n int) {
	t.Helper()
	if err := testdb.Setup(dbName); err != nil {
		t.Fatal(err)
	}
	exec1(t, `INSERT INTO users (id, email, password_hash, role) VALUES (1, 'admin@brightloop.test', 'x', 'admin');
		INSERT INTO company_settings (id) VALUES (1);
		INSERT INTO payment_terms (id, name, due_days) VALUES (1, 'Net 15', 15);
		INSERT INTO tax_rules (id, name, rate_percent, is_default) VALUES (1, 'GST 18%', 18, TRUE), (2, 'Cess 1%', 1, FALSE);
		INSERT INTO plans (id, name, price, billing_period, available_from) VALUES (1, 'Monthly', 500.00, 'monthly', '2026-01-01');
		INSERT INTO products (id, name, product_type, sales_price, cost_price, is_recurring) VALUES
			(1, 'CRM seat', 'software', 33.33, 10, TRUE), (2, 'Setup', 'service', 100, 0, FALSE);
		INSERT INTO product_tax_rules (product_id, tax_rule_id) VALUES (1, 1), (1, 2);
		INSERT INTO discount_rules (id, name, discount_type, value, scope, starts_on, ends_on, max_cycles_per_subscription, created_by)
			VALUES (1, 'Seats 10%', 'percentage', 10, 'products', '2026-01-01', '2026-12-31', 2, 1);
		INSERT INTO discount_rule_products (discount_rule_id, product_id) VALUES (1, 1);`)
	var cust, subs, lines []string
	for i := 1; i <= n; i++ {
		cust = append(cust, fmt.Sprintf("(%d, 'Customer %d')", i, i))
		subs = append(subs, fmt.Sprintf("(%d, 'SUB-%05d', %d, 1, 1, 'active', '2026-01-31', '2026-01-31')", i, i, i))
		lines = append(lines, fmt.Sprintf("(%d, 1, 3, 33.33), (%d, 2, 1, 100.00)", i, i))
	}
	exec1(t, "INSERT INTO customers (id, name) VALUES "+strings.Join(cust, ","))
	exec1(t, "INSERT INTO subscriptions (id, number, customer_id, plan_id, payment_terms_id, status, start_date, next_billing_date) VALUES "+strings.Join(subs, ","))
	exec1(t, "INSERT INTO subscription_lines (subscription_id, product_id, quantity, unit_price) VALUES "+strings.Join(lines, ","))
}

// snapshot lists every invoice with its items, independent of ids and run order.
func snapshot(t *testing.T) []string {
	t.Helper()
	rows, err := g.DB().GetAll(context.Background(), `SELECT CONCAT_WS(' ', inv.subscription_id, inv.period_start, inv.period_end,
			inv.status, inv.subtotal, inv.discount_total, inv.tax_total, inv.total,
			(SELECT GROUP_CONCAT(CONCAT(i.item_type, ':', i.amount, '->', COALESCE(l.description, '')) ORDER BY i.id)
			 FROM invoice_items i LEFT JOIN invoice_items l ON l.id = i.linked_item_id WHERE i.invoice_id = inv.id)) AS s
		FROM invoices inv ORDER BY inv.subscription_id, inv.period_start`)
	if err != nil {
		t.Fatal(err)
	}
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = r["s"].String()
	}
	return out
}

// assertEveryPeriodOnce checks that each of n subscriptions has exactly the three due periods, each invoiced once.
func assertEveryPeriodOnce(t *testing.T, n int) {
	t.Helper()
	if got := count(t, "SELECT COUNT(*) FROM invoices"); got != 3*n {
		t.Fatalf("invoices = %d, want %d", got, 3*n)
	}
	if got := count(t, `SELECT COUNT(*) FROM (SELECT subscription_id FROM billing_records GROUP BY subscription_id
		HAVING COUNT(*) = 3 AND GROUP_CONCAT(period_start ORDER BY period_start) = '2026-01-31,2026-02-28,2026-03-31') ok`); got != n {
		t.Fatalf("%d of %d subscriptions have exactly periods 31 Jan, 28 Feb, 31 Mar", got, n)
	}
	if got := count(t, "SELECT COUNT(*) FROM subscriptions WHERE next_billing_date = '2026-04-30'"); got != n {
		t.Fatalf("%d of %d subscriptions moved next_billing_date to 30 Apr", got, n)
	}
}

func TestInvoiceContents(t *testing.T) {
	reset(t, 1)
	if _, err := Run(context.Background(), asOf, false); err != nil {
		t.Fatal(err)
	}
	got := snapshot(t)
	want := []string{
		"1 2026-01-31 2026-02-27 draft 599.99 10.00 107.10 697.09 recurring:500.00->,recurring:99.99->,discount:-10.00->CRM seat,tax:90.00->Monthly,tax:16.20->CRM seat,tax:0.90->CRM seat",
		"1 2026-02-28 2026-03-30 draft 599.99 10.00 107.10 697.09 recurring:500.00->,recurring:99.99->,discount:-10.00->CRM seat,tax:90.00->Monthly,tax:16.20->CRM seat,tax:0.90->CRM seat",
		"1 2026-03-31 2026-04-29 draft 599.99 0.00 109.00 708.99 recurring:500.00->,recurring:99.99->,tax:90.00->Monthly,tax:18.00->CRM seat,tax:1.00->CRM seat",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if n := count(t, "SELECT COUNT(*) FROM audit_log WHERE table_name = 'invoices' AND action = 'create'"); n != 3 {
		t.Fatalf("audit rows = %d, want 3", n)
	}
	if n := count(t, "SELECT times_used FROM discount_rules WHERE id = 1"); n != 2 {
		t.Fatalf("times_used = %d, want 2", n)
	}
}

func TestRunTwiceCreatesNoDuplicates(t *testing.T) {
	reset(t, 5)
	exec1(t, "UPDATE company_settings SET auto_confirm_recurring_invoices = TRUE")
	ctx := context.Background()
	first, err := Run(ctx, asOf, false)
	if err != nil || first.Created != 15 {
		t.Fatalf("first run: %+v %v", first, err)
	}
	// Force the second run to look at every subscription again, as if next_billing_date had not moved.
	exec1(t, "UPDATE subscriptions SET next_billing_date = '2026-01-31'")
	second, err := Run(ctx, asOf, false)
	if err != nil || second.Created != 0 || second.Skipped != 5 || second.Failed != 0 {
		t.Fatalf("second run: %+v %v", second, err)
	}
	assertEveryPeriodOnce(t, 5)
	// Confirmed numbers are gapless.
	if n := count(t, `SELECT COUNT(*) FROM invoices WHERE status = 'confirmed'
		AND number BETWEEN 'INV-2026-00001' AND 'INV-2026-00015'`); n != 15 {
		t.Fatalf("confirmed+numbered = %d, want 15", n)
	}
}

func TestDryRunWritesNothing(t *testing.T) {
	reset(t, 3)
	res, err := Run(context.Background(), asOf, true)
	if err != nil || res.Created != 9 || len(res.Invoices) != 9 {
		t.Fatalf("dry run: %+v %v", res, err)
	}
	if !res.Invoices[0].Total.Equal(dec("697.09")) {
		t.Fatalf("preview total %s", res.Invoices[0].Total)
	}
	for table, want := range map[string]int{"invoices": 0, "invoice_items": 0, "billing_records": 0, "discount_applications": 0, "audit_log": 0} {
		if n := count(t, "SELECT COUNT(*) FROM "+table); n != want {
			t.Fatalf("%s has %d rows after dry run", table, n)
		}
	}
	if n := count(t, "SELECT COUNT(*) FROM subscriptions WHERE next_billing_date = '2026-01-31'"); n != 3 {
		t.Fatal("dry run moved next_billing_date")
	}
	if n := count(t, "SELECT COUNT(*) FROM billing_runs WHERE dry_run AND created_count = 9 AND finished_at IS NOT NULL"); n != 1 {
		t.Fatal("dry run not recorded in billing_runs")
	}
}

func TestSafetyBoundRefusesRunaway(t *testing.T) {
	reset(t, 2)
	exec1(t, "UPDATE subscriptions SET start_date = '2023-01-31', next_billing_date = '2023-01-31' WHERE id = 2") // 39 periods missing
	res, err := Run(context.Background(), asOf, false)
	if err != nil || res.Created != 3 || res.Failed != 1 || !strings.Contains(res.Errors[0].Message, "safety bound") {
		t.Fatalf("got %+v %v", res, err)
	}
	if n := count(t, "SELECT COUNT(*) FROM invoices WHERE subscription_id = 2"); n != 0 {
		t.Fatalf("runaway subscription got %d invoices", n)
	}
	if n := count(t, "SELECT COUNT(*) FROM billing_runs WHERE failed_count = 1 AND JSON_LENGTH(errors) = 1"); n != 1 {
		t.Fatal("failure not recorded in billing_runs.errors")
	}
}

// TestCrashChild is the process TestCrashHalfwayThenRerun kills. It does nothing in a normal test run.
func TestCrashChild(t *testing.T) {
	after, _ := strconv.Atoi(os.Getenv("BILLING_CRASH_AFTER"))
	if after == 0 {
		t.Skip("only runs as the crash child")
	}
	seen := 0
	testHookBeforeCommit = func(uint64) {
		if seen++; seen == after {
			os.Exit(3) // die mid-transaction: MySQL rolls back this subscription when the connection drops
		}
	}
	Run(context.Background(), asOf, false)
	t.Fatal("child was meant to crash")
}

func TestCrashHalfwayThenRerun(t *testing.T) {
	const subs, crashAt = 10, 4
	reset(t, subs)
	cmd := exec.Command(os.Args[0], "-test.run=^TestCrashChild$", "-test.count=1")
	cmd.Env = append(os.Environ(), fmt.Sprintf("BILLING_CRASH_AFTER=%d", crashAt))
	out, err := cmd.CombinedOutput()
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 3 {
		t.Fatalf("child did not crash as planned: %v\n%s", err, out)
	}
	// Subscriptions 1-3 committed; 4 was mid-transaction and must have left nothing behind.
	if n := count(t, "SELECT COUNT(*) FROM invoices"); n != 3*(crashAt-1) {
		t.Fatalf("after crash: %d invoices, want %d", n, 3*(crashAt-1))
	}
	if n := count(t, "SELECT COUNT(*) FROM invoices WHERE subscription_id = ?", crashAt); n != 0 {
		t.Fatalf("crashed subscription left %d invoices", n)
	}
	if n := count(t, "SELECT COUNT(*) FROM billing_runs WHERE finished_at IS NULL"); n != 1 {
		t.Fatal("crashed run should stay unfinished in billing_runs")
	}

	res, err := Run(context.Background(), asOf, false)
	if err != nil || res.Created != 3*(subs-crashAt+1) || res.Failed != 0 {
		t.Fatalf("rerun: %+v %v", res, err)
	}
	assertEveryPeriodOnce(t, subs)
}

func TestTwoWorkersMatchOneWorker(t *testing.T) {
	const subs = 40
	ctx := context.Background()

	reset(t, subs)
	if _, err := Run(ctx, asOf, false); err != nil {
		t.Fatal(err)
	}
	one := snapshot(t)

	reset(t, subs)
	var wg sync.WaitGroup
	results := make([]*Result, 2)
	for w := range results {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			var err error
			if results[w], err = Run(ctx, asOf, false); err != nil {
				t.Error(err)
			}
		}(w)
	}
	wg.Wait()
	if t.Failed() {
		return
	}
	created, failed := results[0].Created+results[1].Created, results[0].Failed+results[1].Failed
	if created != 3*subs || failed != 0 {
		t.Fatalf("two workers created %d (want %d), failed %d: %+v %+v", created, 3*subs, failed, results[0].Errors, results[1].Errors)
	}
	t.Logf("worker split: %d / %d invoices", results[0].Created, results[1].Created)
	two := snapshot(t)
	if strings.Join(one, "\n") != strings.Join(two, "\n") {
		t.Fatalf("results differ\none worker:\n%s\ntwo workers:\n%s", strings.Join(one, "\n"), strings.Join(two, "\n"))
	}
	assertEveryPeriodOnce(t, subs)
}
