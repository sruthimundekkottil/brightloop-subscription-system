// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ProductVariants is the golang structure of table product_variants for DAO operations like Where/Data.
type ProductVariants struct {
	g.Meta    `orm:"table:product_variants, do:true"`
	Id        interface{} //
	ProductId interface{} //
	Name      interface{} //
	IsActive  interface{} //
}
