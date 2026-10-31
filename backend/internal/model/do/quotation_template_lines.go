// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// QuotationTemplateLines is the golang structure of table quotation_template_lines for DAO operations like Where/Data.
type QuotationTemplateLines struct {
	g.Meta     `orm:"table:quotation_template_lines, do:true"`
	Id         interface{} //
	TemplateId interface{} //
	ProductId  interface{} //
	VariantId  interface{} //
	Quantity   interface{} //
}
