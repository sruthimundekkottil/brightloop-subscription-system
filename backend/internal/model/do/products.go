// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Products is the golang structure of table products for DAO operations like Where/Data.
type Products struct {
	g.Meta      `orm:"table:products, do:true"`
	Id          interface{} //
	Name        interface{} //
	ProductType interface{} //
	SalesPrice  interface{} //
	CostPrice   interface{} //
	Notes       interface{} //
	IsRecurring interface{} //
	IsActive    interface{} //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
