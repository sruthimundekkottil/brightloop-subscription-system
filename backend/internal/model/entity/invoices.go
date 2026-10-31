// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

// Invoices is the golang structure for table invoices.
type Invoices struct {
	Id             uint64          `json:"id"             ` //
	Number         string          `json:"number"         ` //
	CustomerId     uint64          `json:"customerId"     ` //
	SubscriptionId uint64          `json:"subscriptionId" ` //
	Status         string          `json:"status"         ` //
	InvoiceDate    *gtime.Time     `json:"invoiceDate"    ` //
	DueDate        *gtime.Time     `json:"dueDate"        ` //
	PeriodStart    *gtime.Time     `json:"periodStart"    ` //
	PeriodEnd      *gtime.Time     `json:"periodEnd"      ` //
	BillingName    string          `json:"billingName"    ` //
	BillingAddress string          `json:"billingAddress" ` //
	Subtotal       decimal.Decimal `json:"subtotal"       ` //
	DiscountTotal  decimal.Decimal `json:"discountTotal"  ` //
	TaxTotal       decimal.Decimal `json:"taxTotal"       ` //
	Total          decimal.Decimal `json:"total"          ` //
	AmountPaid     decimal.Decimal `json:"amountPaid"     ` //
	PaymentStatus  string          `json:"paymentStatus"  ` //
	ConfirmedAt    *gtime.Time     `json:"confirmedAt"    ` //
	SentAt         *gtime.Time     `json:"sentAt"         ` //
	CancelledAt    *gtime.Time     `json:"cancelledAt"    ` //
	CancelReason   string          `json:"cancelReason"   ` //
	CreatedAt      *gtime.Time     `json:"createdAt"      ` //
	UpdatedAt      *gtime.Time     `json:"updatedAt"      ` //
}
