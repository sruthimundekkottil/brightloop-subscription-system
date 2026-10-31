// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

// Subscriptions is the golang structure for table subscriptions.
type Subscriptions struct {
	Id                  uint64          `json:"id"                  ` //
	Number              string          `json:"number"              ` //
	CustomerId          uint64          `json:"customerId"          ` //
	PlanId              uint64          `json:"planId"              ` //
	TemplateId          uint64          `json:"templateId"          ` //
	PaymentTermsId      uint64          `json:"paymentTermsId"      ` //
	Status              string          `json:"status"              ` //
	StartDate           *gtime.Time     `json:"startDate"           ` //
	ExpiryDate          *gtime.Time     `json:"expiryDate"          ` //
	BillingAnchorDate   *gtime.Time     `json:"billingAnchorDate"   ` //
	NextBillingDate     *gtime.Time     `json:"nextBillingDate"     ` //
	QuotationValidUntil *gtime.Time     `json:"quotationValidUntil" ` //
	ScheduledCloseDate  *gtime.Time     `json:"scheduledCloseDate"  ` //
	RecurringTotal      decimal.Decimal `json:"recurringTotal"      ` //
	RecurringMonthly    decimal.Decimal `json:"recurringMonthly"    ` //
	RenewedFromId       uint64          `json:"renewedFromId"       ` //
	Terms               string          `json:"terms"               ` //
	SalespersonId       uint64          `json:"salespersonId"       ` //
	PausedAt            *gtime.Time     `json:"pausedAt"            ` //
	ClosedAt            *gtime.Time     `json:"closedAt"            ` //
	CloseReasonId       uint64          `json:"closeReasonId"       ` //
	CloseNote           string          `json:"closeNote"           ` //
	CreatedAt           *gtime.Time     `json:"createdAt"           ` //
	UpdatedAt           *gtime.Time     `json:"updatedAt"           ` //
	OpenRenewalMarker   uint64          `json:"openRenewalMarker"   ` //
}
