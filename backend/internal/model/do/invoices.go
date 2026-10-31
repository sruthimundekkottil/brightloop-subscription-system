// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Invoices is the golang structure of table invoices for DAO operations like Where/Data.
type Invoices struct {
	g.Meta         `orm:"table:invoices, do:true"`
	Id             interface{} //
	Number         interface{} //
	CustomerId     interface{} //
	SubscriptionId interface{} //
	Status         interface{} //
	InvoiceDate    *gtime.Time //
	DueDate        *gtime.Time //
	PeriodStart    *gtime.Time //
	PeriodEnd      *gtime.Time //
	BillingName    interface{} //
	BillingAddress interface{} //
	Subtotal       interface{} //
	DiscountTotal  interface{} //
	TaxTotal       interface{} //
	Total          interface{} //
	AmountPaid     interface{} //
	PaymentStatus  interface{} //
	ConfirmedAt    *gtime.Time //
	SentAt         *gtime.Time //
	CancelledAt    *gtime.Time //
	CancelReason   interface{} //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
