// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/shopspring/decimal"
)

// SubscriptionLines is the golang structure for table subscription_lines.
type SubscriptionLines struct {
	Id             uint64          `json:"id"             ` //
	SubscriptionId uint64          `json:"subscriptionId" ` //
	ProductId      uint64          `json:"productId"      ` //
	VariantId      uint64          `json:"variantId"      ` //
	Quantity       uint            `json:"quantity"       ` //
	UnitPrice      decimal.Decimal `json:"unitPrice"      ` //
	Amount         decimal.Decimal `json:"amount"         ` //
}
