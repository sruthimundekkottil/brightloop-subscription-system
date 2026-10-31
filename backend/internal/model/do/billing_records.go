// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// BillingRecords is the golang structure of table billing_records for DAO operations like Where/Data.
type BillingRecords struct {
	g.Meta         `orm:"table:billing_records, do:true"`
	Id             interface{} //
	SubscriptionId interface{} //
	PeriodStart    *gtime.Time //
	PeriodEnd      *gtime.Time //
	InvoiceId      interface{} //
	Inputs         interface{} //
	CreatedAt      *gtime.Time //
}
