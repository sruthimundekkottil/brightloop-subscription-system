// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ProductVariantValues is the golang structure of table product_variant_values for DAO operations like Where/Data.
type ProductVariantValues struct {
	g.Meta                  `orm:"table:product_variant_values, do:true"`
	VariantId               interface{} //
	ProductAttributeValueId interface{} //
}
