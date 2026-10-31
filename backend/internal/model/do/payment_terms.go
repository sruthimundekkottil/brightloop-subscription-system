// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentTerms is the golang structure of table payment_terms for DAO operations like Where/Data.
type PaymentTerms struct {
	g.Meta  `orm:"table:payment_terms, do:true"`
	Id      interface{} //
	Name    interface{} //
	DueDays interface{} //
}
