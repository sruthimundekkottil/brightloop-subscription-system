// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DiscountRules is the golang structure of table discount_rules for DAO operations like Where/Data.
type DiscountRules struct {
	g.Meta                   `orm:"table:discount_rules, do:true"`
	Id                       interface{} //
	Name                     interface{} //
	DiscountType             interface{} //
	Value                    interface{} //
	Scope                    interface{} //
	CustomerId               interface{} //
	MinPurchaseAmount        interface{} //
	MinQuantity              interface{} //
	StartsOn                 *gtime.Time //
	EndsOn                   *gtime.Time //
	MaxUses                  interface{} //
	MaxUsesPerCustomer       interface{} //
	MaxCyclesPerSubscription interface{} //
	MaxDiscountAmount        interface{} //
	TimesUsed                interface{} //
	IsActive                 interface{} //
	CreatedBy                interface{} //
	CreatedAt                *gtime.Time //
}
