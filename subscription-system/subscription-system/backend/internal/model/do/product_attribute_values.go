// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ProductAttributeValues is the golang structure of table product_attribute_values for DAO operations like Where/Data.
type ProductAttributeValues struct {
	g.Meta           `orm:"table:product_attribute_values, do:true"`
	Id               interface{} //
	ProductId        interface{} //
	AttributeValueId interface{} //
	PriceExtra       interface{} //
	IsActive         interface{} //
}
