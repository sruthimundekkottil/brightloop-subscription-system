// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Subscriptions is the golang structure of table subscriptions for DAO operations like Where/Data.
type Subscriptions struct {
	g.Meta              `orm:"table:subscriptions, do:true"`
	Id                  interface{} //
	Number              interface{} //
	CustomerId          interface{} //
	PlanId              interface{} //
	TemplateId          interface{} //
	PaymentTermsId      interface{} //
	Status              interface{} //
	StartDate           *gtime.Time //
	ExpiryDate          *gtime.Time //
	BillingAnchorDate   *gtime.Time //
	NextBillingDate     *gtime.Time //
	QuotationValidUntil *gtime.Time //
	ScheduledCloseDate  *gtime.Time //
	RecurringTotal      interface{} //
	RecurringMonthly    interface{} //
	RenewedFromId       interface{} //
	Terms               interface{} //
	SalespersonId       interface{} //
	PausedAt            *gtime.Time //
	ClosedAt            *gtime.Time //
	CloseReasonId       interface{} //
	CloseNote           interface{} //
	CreatedAt           *gtime.Time //
	UpdatedAt           *gtime.Time //
	OpenRenewalMarker   interface{} //
}
