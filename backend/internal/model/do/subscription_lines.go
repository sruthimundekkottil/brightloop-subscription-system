// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// SubscriptionLines is the golang structure of table subscription_lines for DAO operations like Where/Data.
type SubscriptionLines struct {
	g.Meta         `orm:"table:subscription_lines, do:true"`
	Id             interface{} //
	SubscriptionId interface{} //
	ProductId      interface{} //
	VariantId      interface{} //
	Quantity       interface{} //
	UnitPrice      interface{} //
	Amount         interface{} //
}
