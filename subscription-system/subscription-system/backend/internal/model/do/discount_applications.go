// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DiscountApplications is the golang structure of table discount_applications for DAO operations like Where/Data.
type DiscountApplications struct {
	g.Meta         `orm:"table:discount_applications, do:true"`
	Id             interface{} //
	DiscountRuleId interface{} //
	InvoiceId      interface{} //
	CustomerId     interface{} //
	SubscriptionId interface{} //
	Amount         interface{} //
	ReversedAt     *gtime.Time //
	AppliedAt      *gtime.Time //
}
