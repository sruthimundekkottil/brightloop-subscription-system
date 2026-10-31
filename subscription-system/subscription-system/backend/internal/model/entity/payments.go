// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

// Payments is the golang structure for table payments.
type Payments struct {
	Id             uint64          `json:"id"             ` //
	InvoiceId      uint64          `json:"invoiceId"      ` //
	Method         string          `json:"method"         ` //
	Amount         decimal.Decimal `json:"amount"         ` //
	PaymentDate    *gtime.Time     `json:"paymentDate"    ` //
	Note           string          `json:"note"           ` //
	IdempotencyKey string          `json:"idempotencyKey" ` //
	RecordedBy     uint64          `json:"recordedBy"     ` //
	CreatedAt      *gtime.Time     `json:"createdAt"      ` //
}
