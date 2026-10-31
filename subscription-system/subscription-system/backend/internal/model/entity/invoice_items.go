// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

// InvoiceItems is the golang structure for table invoice_items.
type InvoiceItems struct {
	Id             uint64          `json:"id"             ` //
	InvoiceId      uint64          `json:"invoiceId"      ` //
	ItemType       string          `json:"itemType"       ` //
	LinkedItemId   uint64          `json:"linkedItemId"   ` //
	ProductId      uint64          `json:"productId"      ` //
	PlanId         uint64          `json:"planId"         ` //
	DiscountRuleId uint64          `json:"discountRuleId" ` //
	TaxRuleId      uint64          `json:"taxRuleId"      ` //
	Description    string          `json:"description"    ` //
	ServiceStart   *gtime.Time     `json:"serviceStart"   ` //
	ServiceEnd     *gtime.Time     `json:"serviceEnd"     ` //
	Quantity       uint            `json:"quantity"       ` //
	UnitPrice      decimal.Decimal `json:"unitPrice"      ` //
	RatePercent    decimal.Decimal `json:"ratePercent"    ` //
	Amount         decimal.Decimal `json:"amount"         ` //
	CreatedAt      *gtime.Time     `json:"createdAt"      ` //
}
