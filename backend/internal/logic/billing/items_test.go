package billing

import (
	"math/rand"
	"testing"

	"github.com/shopspring/decimal"
)

func dec(s string) decimal.Decimal { return decimal.RequireFromString(s) }

var gst = TaxRule{ID: 1, Name: "GST 18%", Rate: dec("18")}

func amounts(items []Item) []string {
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.Type + ":" + it.Amount.StringFixed(2)
	}
	return out
}

func TestPriceItemsLinksDiscountAndTaxToBase(t *testing.T) {
	bases := []Base{
		{PlanID: 1, Description: "Monthly plan", Quantity: 1, UnitPrice: dec("500"), Taxes: []TaxRule{gst}},
		{ProductID: 9, Description: "CRM seats", Quantity: 3, UnitPrice: dec("33.33"), Taxes: []TaxRule{gst}},
	}
	items, err := PriceItems(bases, []Discount{{RuleID: 4, Name: "Seats 10%", Percentage: true, Value: dec("10"), ProductIDs: map[uint64]bool{9: true}}})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"recurring:500.00", "recurring:99.99", "discount:-10.00", "tax:90.00", "tax:16.20"}
	if got := amounts(items); len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	} else {
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("got %v, want %v", got, want)
			}
		}
	}
	// discount and second tax point at the seat line; first tax at the plan line
	if items[2].Link != 1 || items[2].DiscountRuleID != 4 || items[3].Link != 0 || items[4].Link != 1 || items[0].Link != -1 {
		t.Fatalf("links wrong: %+v", items)
	}
	sub, disc, tax, total := Totals(items)
	if sub.StringFixed(2) != "599.99" || disc.StringFixed(2) != "10.00" || tax.StringFixed(2) != "106.20" || total.StringFixed(2) != "696.19" {
		t.Fatalf("totals %s %s %s %s", sub, disc, tax, total)
	}
}

func TestRoundingIsHalfAwayFromZero(t *testing.T) {
	// 2.50 * 1% = 0.025: half-away gives 0.03 (banker's rounding would give 0.02)
	items, _ := PriceItems([]Base{{PlanID: 1, Quantity: 1, UnitPrice: dec("2.50"), Taxes: []TaxRule{{ID: 2, Name: "Cess", Rate: dec("1")}}}}, nil)
	if got := items[1].Amount.StringFixed(2); got != "0.03" {
		t.Fatalf("tax = %s, want 0.03", got)
	}
}

func TestBestDiscountWinsAndCapSplitsProportionally(t *testing.T) {
	bases := []Base{
		{PlanID: 1, Quantity: 1, UnitPrice: dec("100")},
		{ProductID: 2, Quantity: 1, UnitPrice: dec("200")},
	}
	items, err := PriceItems(bases, []Discount{
		{RuleID: 1, Name: "Flat 20", Value: dec("20")},
		{RuleID: 2, Name: "50% capped at 100", Percentage: true, Value: dec("50"), Cap: decimal.NewNullDecimal(dec("100"))},
	})
	if err != nil {
		t.Fatal(err)
	}
	// rule 2 wins (100 > 20); cap 100 split 1:2 -> 33.33 + 66.67
	got := amounts(items)
	if len(got) != 4 || got[2] != "discount:-33.33" || got[3] != "discount:-66.67" || items[2].DiscountRuleID != 2 {
		t.Fatalf("got %v", got)
	}
}

func TestFixedDiscountNeverExceedsItems(t *testing.T) {
	items, _ := PriceItems([]Base{{PlanID: 1, Quantity: 1, UnitPrice: dec("15")}}, []Discount{{RuleID: 1, Value: dec("50")}})
	if _, d, _, total := Totals(items); d.StringFixed(2) != "15.00" || !total.IsZero() {
		t.Fatalf("discount %s total %s", d, total)
	}
}

// Self-check: random fixed discounts always split exactly, stay within each item, and keep every discount item negative.
func TestSpreadInvariant(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	for n := 0; n < 5000; n++ {
		var bases []Base
		subtotal := decimal.Zero
		for i := 0; i < 1+r.Intn(5); i++ {
			p := decimal.New(int64(1+r.Intn(100000)), -2)
			bases = append(bases, Base{ProductID: uint64(i + 1), Quantity: 1, UnitPrice: p})
			subtotal = subtotal.Add(p)
		}
		target := decimal.New(int64(1+r.Intn(int(subtotal.Shift(2).IntPart()))), -2)
		items, err := PriceItems(bases, []Discount{{RuleID: 1, Value: target}})
		if err != nil {
			t.Fatalf("case %d: %v", n, err)
		}
		_, d, _, _ := Totals(items)
		if !d.Equal(target) {
			t.Fatalf("case %d: split %s, want %s", n, d, target)
		}
		for _, it := range items {
			if it.Type == "discount" && (!it.Amount.IsNegative() || it.Amount.Neg().GreaterThan(items[it.Link].Amount)) {
				t.Fatalf("case %d: bad discount item %+v", n, it)
			}
		}
	}
}
