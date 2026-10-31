package billing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shopspring/decimal"

	"brightloop/internal/logic/platform"
	"brightloop/internal/logic/subscription/billingdate"
)

// MaxPeriodsPerRun is the safety bound: a subscription missing more periods than this is refused, not billed.
const MaxPeriodsPerRun = 31

const batchSize = 100

var errDryRun = errors.New("dry run: roll back")

// testHookBeforeCommit runs inside each subscription's transaction after all its writes. Tests use it to crash mid-run.
var testHookBeforeCommit func(subscriptionID uint64)

type Invoice struct {
	SubscriptionID         uint64
	PeriodStart, PeriodEnd time.Time
	Items                  []Item
	Total                  decimal.Decimal
}

type Result struct {
	RunID                         uint64
	Due, Created, Skipped, Failed int
	Errors                        []RunError
	Invoices                      []Invoice // what was (or, in a dry run, would be) invoiced
}

type RunError struct {
	SubscriptionID uint64 `json:"subscription_id"`
	Message        string `json:"message"`
}

type settings struct {
	prefix      string
	autoConfirm bool
}

// Run bills every active subscription due on or before asOf: for each one it computes the periods that should be
// billed, removes those already in billing_records, and invoices the rest in one transaction.
// A dry run computes the same invoices and rolls every subscription transaction back; only the billing_runs row stays.
func Run(ctx context.Context, asOf time.Time, dryRun bool) (*Result, error) {
	asOfDate := asOf.Format(platform.DateLayout)
	db := g.DB()
	r, err := db.Exec(ctx, "INSERT INTO billing_runs (as_of_date, dry_run) VALUES (?, ?)", asOfDate, dryRun)
	if err != nil {
		return nil, err
	}
	id, _ := r.LastInsertId()
	res := &Result{RunID: uint64(id)}

	cfg, err := db.GetOne(ctx, "SELECT invoice_prefix, auto_confirm_recurring_invoices FROM company_settings WHERE id = 1")
	if err != nil {
		return res, err
	}
	if cfg.IsEmpty() {
		return res, errors.New("company_settings row missing")
	}
	st := settings{prefix: cfg["invoice_prefix"].String(), autoConfirm: cfg["auto_confirm_recurring_invoices"].Bool()}

	var lastID uint64
	for ctx.Err() == nil {
		ids, err := db.GetArray(ctx, `SELECT id FROM subscriptions WHERE status = 'active' AND next_billing_date <= ? AND id > ?
			ORDER BY id LIMIT ?`, asOfDate, lastID, batchSize)
		if err != nil {
			return res, err
		}
		if len(ids) == 0 {
			break
		}
		for _, v := range ids {
			subID := v.Uint64()
			lastID = subID
			res.Due++
			invs, err := billWithRetry(ctx, subID, asOf, dryRun, st, res.RunID)
			switch {
			case err != nil:
				res.Failed++
				res.Errors = append(res.Errors, RunError{subID, err.Error()})
			case len(invs) == 0:
				res.Skipped++
			default:
				res.Created += len(invs)
				res.Invoices = append(res.Invoices, invs...)
			}
		}
	}

	var errs any // NULL when nothing failed
	if res.Errors != nil {
		js, _ := json.Marshal(res.Errors)
		errs = string(js)
	}
	_, err = db.Exec(context.WithoutCancel(ctx), `UPDATE billing_runs SET finished_at = NOW(3), due_count = ?, created_count = ?,
		skipped_count = ?, failed_count = ?, errors = ? WHERE id = ?`, res.Due, res.Created, res.Skipped, res.Failed, errs, res.RunID)
	if err == nil {
		err = ctx.Err()
	}
	return res, err
}

// billWithRetry retries once on deadlock (1213) or lock wait timeout (1205).
func billWithRetry(ctx context.Context, subID uint64, asOf time.Time, dryRun bool, st settings, runID uint64) ([]Invoice, error) {
	invs, err := billOne(ctx, subID, asOf, dryRun, st, runID)
	var me *mysql.MySQLError
	if errors.As(err, &me) && (me.Number == 1213 || me.Number == 1205) {
		invs, err = billOne(ctx, subID, asOf, dryRun, st, runID)
	}
	return invs, err
}

func billOne(ctx context.Context, subID uint64, asOf time.Time, dryRun bool, st settings, runID uint64) ([]Invoice, error) {
	var out []Invoice
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		out = nil
		// Lock only the subscription row; another worker holding it means it is being billed right now.
		sub, err := tx.GetOne(`SELECT id, customer_id, plan_id, payment_terms_id,
				DATE_FORMAT(COALESCE(billing_anchor_date, start_date), '%Y-%m-%d') AS anchor,
				DATE_FORMAT(next_billing_date, '%Y-%m-%d') AS next_billing,
				DATE_FORMAT(expiry_date, '%Y-%m-%d') AS expiry,
				DATE_FORMAT(scheduled_close_date, '%Y-%m-%d') AS close_on
			FROM subscriptions WHERE id = ? AND status = 'active' AND next_billing_date <= ? FOR UPDATE SKIP LOCKED`,
			subID, asOf.Format(platform.DateLayout))
		if err != nil || sub.IsEmpty() {
			return err
		}
		plan, err := tx.GetOne("SELECT id, name, price, billing_period, billing_interval_count FROM plans WHERE id = ?", sub["plan_id"])
		if err != nil {
			return err
		}
		anchor, next := mustDate(sub["anchor"].String()), mustDate(sub["next_billing"].String())
		period, every := plan["billing_period"].String(), plan["billing_interval_count"].Int()

		// Compute what should be billed up to asOf, starting at the period that contains next_billing_date ...
		kAfter, _ := billingdate.NextPeriodOnOrAfter(anchor, period, every, next.AddDate(0, 0, 1))
		k := max(kAfter-1, 0)
		start := billingdate.PeriodStart(anchor, period, every, k)
		var due []int
		for ; !start.After(asOf); k, start = k+1, billingdate.PeriodStart(anchor, period, every, k+1) {
			if e := sub["expiry"].String(); e != "" && start.After(mustDate(e)) {
				break
			}
			if c := sub["close_on"].String(); c != "" && !start.Before(mustDate(c)) {
				break
			}
			due = append(due, k)
		}
		// ... minus what billing_records says is already billed.
		billed, err := tx.GetAll("SELECT DATE_FORMAT(period_start, '%Y-%m-%d') AS ps FROM billing_records WHERE subscription_id = ?", subID)
		if err != nil {
			return err
		}
		done := map[string]bool{}
		for _, b := range billed {
			done[b["ps"].String()] = true
		}
		missing := due[:0]
		for _, k := range due {
			if !done[billingdate.PeriodStart(anchor, period, every, k).Format(platform.DateLayout)] {
				missing = append(missing, k)
			}
		}
		if len(missing) > MaxPeriodsPerRun {
			return fmt.Errorf("safety bound: %d periods missing, more than %d; nothing billed", len(missing), MaxPeriodsPerRun)
		}

		bases, err := baseItems(tx, sub["id"].Uint64(), plan)
		if err != nil {
			return err
		}
		customer, err := tx.GetOne(`SELECT c.name, c.billing_address, t.due_days FROM customers c, payment_terms t
			WHERE c.id = ? AND t.id = ?`, sub["customer_id"], sub["payment_terms_id"])
		if err != nil {
			return err
		}
		for _, k := range missing {
			ps, pe := billingdate.PeriodStart(anchor, period, every, k), billingdate.PeriodEnd(anchor, period, every, k)
			inv, err := writeInvoice(tx, sub, customer, bases, ps, pe, st, runID)
			if err != nil {
				return err
			}
			out = append(out, inv)
		}

		_, after := billingdate.NextPeriodOnOrAfter(anchor, period, every, asOf.AddDate(0, 0, 1))
		if _, err = tx.Exec("UPDATE subscriptions SET next_billing_date = ? WHERE id = ?", after.Format(platform.DateLayout), subID); err != nil {
			return err
		}
		if testHookBeforeCommit != nil {
			testHookBeforeCommit(subID)
		}
		if dryRun {
			return errDryRun
		}
		return nil
	})
	if errors.Is(err, errDryRun) {
		err = nil
	}
	return out, err
}

// baseItems is the plan line plus the subscription's recurring product lines, each with its tax rules.
func baseItems(tx gdb.TX, subID uint64, plan gdb.Record) ([]Base, error) {
	defaultTax, err := taxRules(tx, "SELECT id, name, rate_percent FROM tax_rules WHERE is_default AND is_active")
	if err != nil {
		return nil, err
	}
	price, err := decimal.NewFromString(plan["price"].String())
	if err != nil {
		return nil, err
	}
	bases := []Base{{PlanID: plan["id"].Uint64(), Description: plan["name"].String(), Quantity: 1, UnitPrice: price, Taxes: defaultTax}}

	// shortcut: one-time products (is_recurring = FALSE) are not billed; add them as 'fixed' items when that is in scope.
	lines, err := tx.GetAll(`SELECT l.product_id, p.name, l.quantity, l.unit_price FROM subscription_lines l
		JOIN products p ON p.id = l.product_id WHERE l.subscription_id = ? AND p.is_recurring ORDER BY l.id`, subID)
	if err != nil {
		return nil, err
	}
	for _, l := range lines {
		unit, err := decimal.NewFromString(l["unit_price"].String())
		if err != nil {
			return nil, err
		}
		taxes, err := taxRules(tx, `SELECT t.id, t.name, t.rate_percent FROM product_tax_rules pt
			JOIN tax_rules t ON t.id = pt.tax_rule_id WHERE pt.product_id = ? AND t.is_active ORDER BY t.id`, l["product_id"])
		if err != nil {
			return nil, err
		}
		if len(taxes) == 0 {
			taxes = defaultTax
		}
		bases = append(bases, Base{ProductID: l["product_id"].Uint64(), Description: l["name"].String(),
			Quantity: l["quantity"].Int64(), UnitPrice: unit, Taxes: taxes})
	}
	return bases, nil
}

func taxRules(tx gdb.TX, sql string, args ...any) ([]TaxRule, error) {
	rows, err := tx.GetAll(sql, args...)
	if err != nil {
		return nil, err
	}
	var out []TaxRule
	for _, r := range rows {
		rate, err := decimal.NewFromString(r["rate_percent"].String())
		if err != nil {
			return nil, err
		}
		out = append(out, TaxRule{ID: r["id"].Uint64(), Name: r["name"].String(), Rate: rate})
	}
	return out, nil
}

// applicableDiscounts locks eligible rules in ascending id order, so concurrent workers cannot deadlock on them
// or both take the last use. shortcut: replace with member B's code (pricing.ApplicableDiscounts).
func applicableDiscounts(tx gdb.TX, customerID, subID uint64, on string, bases []Base) ([]Discount, error) {
	rows, err := tx.GetAll(`SELECT id, name, discount_type, value, scope, min_purchase_amount, min_quantity,
			max_uses_per_customer, max_cycles_per_subscription, max_discount_amount
		FROM discount_rules
		WHERE is_active AND ? BETWEEN starts_on AND ends_on
		  AND (max_uses IS NULL OR times_used < max_uses)
		  AND (customer_id IS NULL OR customer_id = ?)
		ORDER BY id FOR UPDATE`, on, customerID)
	if err != nil {
		return nil, err
	}
	var out []Discount
	for _, r := range rows {
		ruleID := r["id"].Uint64()
		// Usage limits are read after the rule lock, which every writer of discount_applications holds.
		if !r["max_uses_per_customer"].IsNil() {
			n, err := tx.GetValue("SELECT COUNT(*) FROM discount_applications WHERE discount_rule_id = ? AND customer_id = ? AND reversed_at IS NULL", ruleID, customerID)
			if err != nil {
				return nil, err
			}
			if n.Int() >= r["max_uses_per_customer"].Int() {
				continue
			}
		}
		if !r["max_cycles_per_subscription"].IsNil() {
			n, err := tx.GetValue("SELECT COUNT(*) FROM discount_applications WHERE discount_rule_id = ? AND subscription_id = ? AND reversed_at IS NULL", ruleID, subID)
			if err != nil {
				return nil, err
			}
			if n.Int() >= r["max_cycles_per_subscription"].Int() {
				continue
			}
		}
		d := Discount{RuleID: ruleID, Name: r["name"].String(), Percentage: r["discount_type"].String() == "percentage",
			Value: decimal.RequireFromString(r["value"].String())}
		if !r["max_discount_amount"].IsNil() {
			d.Cap = decimal.NewNullDecimal(decimal.RequireFromString(r["max_discount_amount"].String()))
		}
		if r["scope"].String() == "products" {
			ids, err := tx.GetAll("SELECT product_id FROM discount_rule_products WHERE discount_rule_id = ?", ruleID)
			if err != nil {
				return nil, err
			}
			d.ProductIDs = map[uint64]bool{}
			for _, id := range ids {
				d.ProductIDs[id["product_id"].Uint64()] = true
			}
		}
		amount, qty := decimal.Zero, int64(0)
		for _, b := range bases {
			if d.ProductIDs == nil || d.ProductIDs[b.ProductID] {
				amount = amount.Add(round2(b.UnitPrice.Mul(decimal.NewFromInt(b.Quantity))))
				qty += b.Quantity
			}
		}
		if !r["min_purchase_amount"].IsNil() && amount.LessThan(decimal.RequireFromString(r["min_purchase_amount"].String())) {
			continue
		}
		if !r["min_quantity"].IsNil() && qty < r["min_quantity"].Int64() {
			continue
		}
		out = append(out, d)
	}
	return out, nil
}

func writeInvoice(tx gdb.TX, sub, customer gdb.Record, bases []Base, ps, pe time.Time, st settings, runID uint64) (Invoice, error) {
	subID, customerID := sub["id"].Uint64(), sub["customer_id"].Uint64()
	invDate := ps.Format(platform.DateLayout)
	discounts, err := applicableDiscounts(tx, customerID, subID, invDate, bases)
	if err != nil {
		return Invoice{}, err
	}
	items, err := PriceItems(bases, discounts)
	if err != nil {
		return Invoice{}, err
	}
	subtotal, discTotal, taxTotal, total := Totals(items)
	inv := Invoice{SubscriptionID: subID, PeriodStart: ps, PeriodEnd: pe, Items: items, Total: total}

	r, err := tx.Exec(`INSERT INTO invoices (customer_id, subscription_id, invoice_date, due_date, period_start, period_end,
			billing_name, billing_address, subtotal, discount_total, tax_total, total)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		customerID, subID, invDate, ps.AddDate(0, 0, customer["due_days"].Int()).Format(platform.DateLayout), invDate,
		pe.Format(platform.DateLayout), customer["name"].String(), customer["billing_address"].Val(),
		subtotal.StringFixed(2), discTotal.StringFixed(2), taxTotal.StringFixed(2), total.StringFixed(2))
	if err != nil {
		return inv, err
	}
	id, _ := r.LastInsertId()
	invoiceID := uint64(id)

	itemIDs := make([]int64, len(items)) // base items come first, so links always point at an inserted row
	for i, it := range items {
		var linked any
		if it.Link >= 0 {
			linked = itemIDs[it.Link]
		}
		var svcStart, svcEnd any
		if it.Type == "recurring" {
			svcStart, svcEnd = invDate, pe.Format(platform.DateLayout)
		}
		r, err := tx.Exec(`INSERT INTO invoice_items (invoice_id, item_type, linked_item_id, product_id, plan_id,
				discount_rule_id, tax_rule_id, description, service_start, service_end, quantity, unit_price, rate_percent, amount)
			VALUES (?, ?, ?, NULLIF(?, 0), NULLIF(?, 0), NULLIF(?, 0), NULLIF(?, 0), ?, ?, ?, ?, ?, ?, ?)`,
			invoiceID, it.Type, linked, it.ProductID, it.PlanID, it.DiscountRuleID, it.TaxRule, it.Description,
			svcStart, svcEnd, it.Quantity, nullDec(it.UnitPrice), nullDec(it.RatePercent), it.Amount.StringFixed(2))
		if err != nil {
			return inv, err
		}
		itemIDs[i], _ = r.LastInsertId()
	}

	inputs, _ := json.Marshal(g.Map{"run_id": runID, "plan_id": sub["plan_id"].Uint64(), "bases": bases, "discounts": discounts})
	if _, err = tx.Exec(`INSERT INTO billing_records (subscription_id, period_start, period_end, invoice_id, inputs)
		VALUES (?, ?, ?, ?, ?)`, subID, invDate, pe.Format(platform.DateLayout), invoiceID, string(inputs)); err != nil {
		return inv, err // a duplicate period fails here on uq_billing_period, rolling the whole subscription back
	}

	for _, it := range items {
		if it.Type == "discount" {
			if _, err = tx.Exec(`INSERT INTO discount_applications (discount_rule_id, invoice_id, customer_id, subscription_id, amount)
				VALUES (?, ?, ?, ?, ?)`, it.DiscountRuleID, invoiceID, customerID, subID, discTotal.StringFixed(2)); err != nil {
				return inv, err
			}
			if _, err = tx.Exec("UPDATE discount_rules SET times_used = times_used + 1 WHERE id = ?", it.DiscountRuleID); err != nil {
				return inv, err
			}
			break // one rule per invoice
		}
	}

	if st.autoConfirm {
		year := ps.Format("2006")
		n, err := platform.NextNumber(tx, "invoice-"+year)
		if err != nil {
			return inv, err
		}
		if _, err = tx.Exec("UPDATE invoices SET status = 'confirmed', number = ?, confirmed_at = NOW(3) WHERE id = ?",
			fmt.Sprintf("%s-%s-%05d", st.prefix, year, n), invoiceID); err != nil {
			return inv, err
		}
	}
	return inv, platform.Audit(tx, "invoices", invoiceID, "create", fmt.Sprintf("billing run %d", runID),
		g.Map{"subscription_id": subID, "period_start": invDate, "total": total.StringFixed(2)})
}

func nullDec(d decimal.NullDecimal) any {
	if !d.Valid {
		return nil
	}
	return d.Decimal.StringFixed(2)
}

func mustDate(s string) time.Time {
	t, err := time.Parse(platform.DateLayout, s)
	if err != nil {
		panic(fmt.Sprintf("bad date %q from MySQL: %v", s, err)) // DATE_FORMAT output; cannot happen for valid rows
	}
	return t
}
