// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// BillingRecords is the golang structure for table billing_records.
type BillingRecords struct {
	Id             uint64      `json:"id"             ` //
	SubscriptionId uint64      `json:"subscriptionId" ` //
	PeriodStart    *gtime.Time `json:"periodStart"    ` //
	PeriodEnd      *gtime.Time `json:"periodEnd"      ` //
	InvoiceId      uint64      `json:"invoiceId"      ` //
	Inputs         string      `json:"inputs"         ` //
	CreatedAt      *gtime.Time `json:"createdAt"      ` //
}
