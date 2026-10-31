// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

// DiscountApplications is the golang structure for table discount_applications.
type DiscountApplications struct {
	Id             uint64          `json:"id"             ` //
	DiscountRuleId uint64          `json:"discountRuleId" ` //
	InvoiceId      uint64          `json:"invoiceId"      ` //
	CustomerId     uint64          `json:"customerId"     ` //
	SubscriptionId uint64          `json:"subscriptionId" ` //
	Amount         decimal.Decimal `json:"amount"         ` //
	ReversedAt     *gtime.Time     `json:"reversedAt"     ` //
	AppliedAt      *gtime.Time     `json:"appliedAt"      ` //
}
