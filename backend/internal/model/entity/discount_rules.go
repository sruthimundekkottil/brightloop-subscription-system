// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

// DiscountRules is the golang structure for table discount_rules.
type DiscountRules struct {
	Id                       uint64          `json:"id"                       ` //
	Name                     string          `json:"name"                     ` //
	DiscountType             string          `json:"discountType"             ` //
	Value                    decimal.Decimal `json:"value"                    ` //
	Scope                    string          `json:"scope"                    ` //
	CustomerId               uint64          `json:"customerId"               ` //
	MinPurchaseAmount        decimal.Decimal `json:"minPurchaseAmount"        ` //
	MinQuantity              uint            `json:"minQuantity"              ` //
	StartsOn                 *gtime.Time     `json:"startsOn"                 ` //
	EndsOn                   *gtime.Time     `json:"endsOn"                   ` //
	MaxUses                  uint            `json:"maxUses"                  ` //
	MaxUsesPerCustomer       uint            `json:"maxUsesPerCustomer"       ` //
	MaxCyclesPerSubscription uint            `json:"maxCyclesPerSubscription" ` //
	MaxDiscountAmount        decimal.Decimal `json:"maxDiscountAmount"        ` //
	TimesUsed                uint            `json:"timesUsed"                ` //
	IsActive                 int             `json:"isActive"                 ` //
	CreatedBy                uint64          `json:"createdBy"                ` //
	CreatedAt                *gtime.Time     `json:"createdAt"                ` //
}
