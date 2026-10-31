// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InvoiceItems is the golang structure of table invoice_items for DAO operations like Where/Data.
type InvoiceItems struct {
	g.Meta         `orm:"table:invoice_items, do:true"`
	Id             interface{} //
	InvoiceId      interface{} //
	ItemType       interface{} //
	LinkedItemId   interface{} //
	ProductId      interface{} //
	PlanId         interface{} //
	DiscountRuleId interface{} //
	TaxRuleId      interface{} //
	Description    interface{} //
	ServiceStart   *gtime.Time //
	ServiceEnd     *gtime.Time //
	Quantity       interface{} //
	UnitPrice      interface{} //
	RatePercent    interface{} //
	Amount         interface{} //
	CreatedAt      *gtime.Time //
}
