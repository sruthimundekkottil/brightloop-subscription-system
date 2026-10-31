// Package billing turns due subscription periods into invoices.
package billing

import (
	"fmt"
	"sort"

	"github.com/shopspring/decimal"
)

// shortcut: replace with member B's code — this file stands in for pricing.PriceItems
// (one best discount per invoice, taxes on the discounted amount, item by item).

var hundred = decimal.NewFromInt(100)

// round2 rounds half away from zero (shopspring documents 5.45 -> 5.5).
func round2(d decimal.Decimal) decimal.Decimal { return d.Round(2) }

type TaxRule struct {
	ID   uint64
	Name string
	Rate decimal.Decimal // percent
}

type Discount struct {
	RuleID     uint64
	Name       string
	Percentage bool
	Value      decimal.Decimal     // percent, or a fixed amount
	Cap        decimal.NullDecimal // max_discount_amount per invoice
	ProductIDs map[uint64]bool     // nil = subscription scope (every base item)
}

// Base is one billable line before discounts and taxes.
type Base struct {
	PlanID, ProductID uint64 // exactly one is set
	Description       string
	Quantity          int64
	UnitPrice         decimal.Decimal
	Taxes             []TaxRule
}

// Item is one invoice_items row. Link is the index (in the returned slice) of the base item a discount or tax applies to.
type Item struct {
	Type                                       string // "recurring", "discount", "tax"
	Link                                       int    // -1 for base items
	PlanID, ProductID, DiscountRuleID, TaxRule uint64
	Description                                string
	Quantity                                   int64
	UnitPrice                                  decimal.NullDecimal
	RatePercent                                decimal.NullDecimal
	Amount                                     decimal.Decimal
}

// PriceItems returns base items, then a discount item per discounted base item (best rule only), then tax items.
func PriceItems(bases []Base, discounts []Discount) ([]Item, error) {
	items := make([]Item, 0, len(bases)*3)
	for _, b := range bases {
		items = append(items, Item{Type: "recurring", Link: -1, PlanID: b.PlanID, ProductID: b.ProductID,
			Description: b.Description, Quantity: b.Quantity, UnitPrice: decimal.NewNullDecimal(b.UnitPrice),
			Amount: round2(b.UnitPrice.Mul(decimal.NewFromInt(b.Quantity)))})
	}

	var best []decimal.Decimal // per base item
	var bestRule *Discount
	bestTotal := decimal.Zero
	for i := range discounts {
		per, err := allocate(items, &discounts[i])
		if err != nil {
			return nil, err
		}
		if t := sum(per); t.GreaterThan(bestTotal) {
			best, bestRule, bestTotal = per, &discounts[i], t
		}
	}
	for i := range bases {
		if bestRule != nil && best[i].IsPositive() {
			d := Item{Type: "discount", Link: i, DiscountRuleID: bestRule.RuleID, Description: bestRule.Name, Quantity: 1, Amount: best[i].Neg()}
			if bestRule.Percentage {
				d.RatePercent = decimal.NewNullDecimal(bestRule.Value)
			}
			items = append(items, d)
		}
	}

	for i, b := range bases {
		net := items[i].Amount
		if best != nil {
			net = net.Sub(best[i])
		}
		for _, t := range b.Taxes {
			items = append(items, Item{Type: "tax", Link: i, TaxRule: t.ID, Description: t.Name, Quantity: 1,
				RatePercent: decimal.NewNullDecimal(t.Rate), Amount: round2(net.Mul(t.Rate).Div(hundred))})
		}
	}
	return items, nil
}

// allocate returns the positive discount for each base item (zero when out of scope).
func allocate(items []Item, d *Discount) ([]decimal.Decimal, error) {
	n := 0
	for n < len(items) && items[n].Link == -1 {
		n++
	}
	per := make([]decimal.Decimal, n)
	var scope []int
	for i := 0; i < n; i++ {
		if items[i].Amount.IsPositive() && (d.ProductIDs == nil || d.ProductIDs[items[i].ProductID]) {
			scope = append(scope, i)
		}
	}
	if len(scope) == 0 {
		return per, nil
	}
	if d.Percentage {
		for _, i := range scope {
			per[i] = round2(items[i].Amount.Mul(d.Value).Div(hundred))
		}
		if !d.Cap.Valid || sum(per).LessThanOrEqual(d.Cap.Decimal) {
			return per, nil
		}
		return per, spread(per, items, scope, d.Cap.Decimal)
	}
	target := d.Value
	if d.Cap.Valid {
		target = decimal.Min(target, d.Cap.Decimal)
	}
	return per, spread(per, items, scope, target)
}

// spread splits target across the scoped items in proportion to their amounts; the rounding remainder goes on the largest item.
func spread(per []decimal.Decimal, items []Item, scope []int, target decimal.Decimal) error {
	total := decimal.Zero
	for _, i := range scope {
		total = total.Add(items[i].Amount)
	}
	target = decimal.Min(target, total)
	sort.SliceStable(scope, func(a, b int) bool { return items[scope[a]].Amount.LessThan(items[scope[b]].Amount) })
	given := decimal.Zero
	for _, i := range scope[:len(scope)-1] {
		per[i] = round2(target.Mul(items[i].Amount).Div(total))
		given = given.Add(per[i])
	}
	last := scope[len(scope)-1]
	per[last] = target.Sub(given)
	if per[last].IsNegative() || per[last].GreaterThan(items[last].Amount) {
		return fmt.Errorf("discount %s cannot be split within item amounts", target)
	}
	return nil
}

func sum(ds []decimal.Decimal) decimal.Decimal {
	t := decimal.Zero
	for _, d := range ds {
		t = t.Add(d)
	}
	return t
}

// Totals gives the invoice header columns: subtotal, discount_total (positive), tax_total, total.
func Totals(items []Item) (subtotal, discount, tax, total decimal.Decimal) {
	for _, it := range items {
		switch it.Type {
		case "discount":
			discount = discount.Sub(it.Amount)
		case "tax":
			tax = tax.Add(it.Amount)
		default:
			subtotal = subtotal.Add(it.Amount)
		}
		total = total.Add(it.Amount)
	}
	return
}
