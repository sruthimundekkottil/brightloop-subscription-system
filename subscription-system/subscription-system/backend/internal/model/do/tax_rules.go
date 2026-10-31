// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// TaxRules is the golang structure of table tax_rules for DAO operations like Where/Data.
type TaxRules struct {
	g.Meta        `orm:"table:tax_rules, do:true"`
	Id            interface{} //
	Name          interface{} //
	RatePercent   interface{} //
	IsDefault     interface{} //
	IsActive      interface{} //
	DefaultMarker interface{} //
}
